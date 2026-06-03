package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/smol-utils/actiondoc/internal/model"
)

// parseString writes src to a temp .yml file and parses it as a workflow.
func parseString(t *testing.T, src string) (*model.Workflow, error) {
	t.Helper()
	p := filepath.Join(t.TempDir(), "wf.yml")
	if err := os.WriteFile(p, []byte(src), 0o644); err != nil {
		t.Fatalf("write temp: %v", err)
	}
	return ParseFile(p)
}

// TestParseStepFields covers step `with:` blocks, `continue-on-error:`, and the trailing
// version comment on SHA-pinned `uses:` refs.
func TestParseStepFields(t *testing.T) {
	w, err := ParseFile(testdataPath("s3/steps.yml"))
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}
	if len(w.Jobs) != 3 {
		t.Fatalf("got %d jobs, want 3", len(w.Jobs))
	}
	build := w.Jobs[0]
	if len(build.Steps) != 5 {
		t.Fatalf("build job: got %d steps, want 5", len(build.Steps))
	}

	// Step 1: named, SHA-pinned uses with a version comment, two with: entries.
	checkout := build.Steps[0]
	if checkout.UsesVersion != "v4.1.1" {
		t.Errorf("UsesVersion = %q, want v4.1.1", checkout.UsesVersion)
	}
	wantWith := []model.KV{
		{Key: "fetch-depth", Value: "0"},
		{Key: "token", Value: "${{ secrets.CHECKOUT_TOKEN }}"},
	}
	if len(checkout.With) != len(wantWith) {
		t.Fatalf("With = %+v, want %+v", checkout.With, wantWith)
	}
	for i, kv := range wantWith {
		if checkout.With[i] != kv {
			t.Errorf("With[%d] = %+v, want %+v", i, checkout.With[i], kv)
		}
	}

	// Step 3: tag-pinned uses (not a SHA) -> no version captured.
	if v := build.Steps[2].UsesVersion; v != "" {
		t.Errorf("tag-pinned uses got UsesVersion %q, want empty", v)
	}

	// Step 4: run-only with continue-on-error.
	if !build.Steps[3].ContinueOnError {
		t.Error("expected ContinueOnError on run-only step")
	}
	if build.Steps[0].ContinueOnError {
		t.Error("unexpected ContinueOnError on checkout step")
	}

	// Deploy job step 1: env: block parsed in source order.
	deploy := w.Jobs[1]
	wantEnv := []model.KV{
		{Key: "DOCKER_BUILDKIT", Value: "1"},
		{Key: "IMAGE_SIGNING_KEY", Value: "${{ secrets.IMAGE_SIGNING_KEY }}"},
	}
	push := deploy.Steps[0]
	if len(push.Env) != len(wantEnv) {
		t.Fatalf("Env = %+v, want %+v", push.Env, wantEnv)
	}
	for i, kv := range wantEnv {
		if push.Env[i] != kv {
			t.Errorf("Env[%d] = %+v, want %+v", i, push.Env[i], kv)
		}
	}
}

// TestParseMultilineNames verifies that block-scalar name: values (which carry embedded
// or trailing newlines) are normalized to a single line at parse time. A multi-line name
// would otherwise break the Markdown structures built around names: bold step titles,
// ASCII call-graph tree labels, and heading anchors.
func TestParseMultilineNames(t *testing.T) {
	src := `name: >
  Continuous
  Integration
on: push
jobs:
  build:
    name: |
      Build and
      Test
    runs-on: ubuntu-latest
    steps:
      - name: >-
          Migration Tests: ${{ matrix.python-version }}:
          ${{ env.PARALLEL_TEST_TYPES }}
        run: ./run.sh
`
	w, err := parseString(t, src)
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}
	for _, tc := range []struct{ got, want, what string }{
		{w.Name, "Continuous Integration", "workflow name"},
		{w.Jobs[0].Name, "Build and Test", "job name"},
		{w.Jobs[0].Steps[0].Name, "Migration Tests: ${{ matrix.python-version }}: ${{ env.PARALLEL_TEST_TYPES }}", "step name"},
	} {
		if tc.got != tc.want {
			t.Errorf("%s = %q, want %q", tc.what, tc.got, tc.want)
		}
		if strings.Contains(tc.got, "\n") {
			t.Errorf("%s still contains a newline: %q", tc.what, tc.got)
		}
	}
}

// TestParseContinueOnErrorExpression verifies that an expression-valued
// continue-on-error (e.g. matrix-driven) is captured as an expression rather than
// silently coerced to false, so the step is still flagged as failure-tolerant.
func TestParseContinueOnErrorExpression(t *testing.T) {
	src := `name: CI
on: push
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - name: maybe-flaky
        run: ./flaky.sh
        continue-on-error: ${{ matrix.experimental }}
      - name: hard-tolerant
        run: ./x.sh
        continue-on-error: true
`
	w, err := parseString(t, src)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	steps := w.Jobs[0].Steps
	if steps[0].ContinueOnError {
		t.Error("expression continue-on-error must not be coerced to the literal-true bool")
	}
	if steps[0].ContinueOnErrorExpr != "${{ matrix.experimental }}" {
		t.Errorf("ContinueOnErrorExpr = %q, want the raw expression", steps[0].ContinueOnErrorExpr)
	}
	if !steps[1].ContinueOnError || steps[1].ContinueOnErrorExpr != "" {
		t.Errorf("literal true: ContinueOnError=%v Expr=%q, want true/empty", steps[1].ContinueOnError, steps[1].ContinueOnErrorExpr)
	}
}

// TestParseRunsOnVariants covers scalar, list, and map (group/labels) runs-on forms.
func TestParseRunsOnVariants(t *testing.T) {
	w, err := ParseFile(testdataPath("s3/steps.yml"))
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}

	tests := []struct {
		job  int
		want string
	}{
		{0, "self-hosted, linux, x64"},                     // list
		{1, "group: deploy-runners, labels: linux, arm64"}, // map with nested list
		{2, "ubuntu-latest"},                               // scalar
	}
	for _, tt := range tests {
		if got := w.Jobs[tt.job].RunsOn; got != tt.want {
			t.Errorf("job %d RunsOn = %q, want %q", tt.job, got, tt.want)
		}
	}
}

// TestParseMatrix covers literal scalar axes, list-of-objects axes (flattened to dotted
// names), and include: entries merging into the axes with the adjusted flag set.
func TestParseMatrix(t *testing.T) {
	w, err := ParseFile(testdataPath("s3/steps.yml"))
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}

	// build: scalar axis java: [17, 21, 24]
	build := w.Jobs[0]
	if vals, ok := build.MatrixValues("java"); !ok || strings.Join(vals, ",") != "17,21,24" {
		t.Errorf("build matrix java = %v (found=%v), want [17 21 24]", vals, ok)
	}

	// deploy: list-of-objects axis target -> dotted axes target.env / target.url
	deploy := w.Jobs[1]
	if vals, ok := deploy.MatrixValues("target.env"); !ok || strings.Join(vals, ",") != "staging,production" {
		t.Errorf("deploy matrix target.env = %v (found=%v), want [staging production]", vals, ok)
	}
	if vals, ok := deploy.MatrixValues("target.url"); !ok || len(vals) != 2 {
		t.Errorf("deploy matrix target.url = %v (found=%v), want 2 values", vals, ok)
	}

	// verify: literal axis case: [a, b] plus include: adding case: c -> values merge and
	// the matrix is flagged as adjusted.
	verify := w.Jobs[2]
	if vals, ok := verify.MatrixValues("case"); !ok || strings.Join(vals, ",") != "a,b,c" {
		t.Errorf("verify matrix case = %v (found=%v), want [a b c] (include values merged)", vals, ok)
	}
	if !verify.MatrixAdjusted {
		t.Error("verify matrix must be flagged adjusted (include: present)")
	}
	// build and deploy have no include/exclude -> not adjusted.
	if build.MatrixAdjusted || deploy.MatrixAdjusted {
		t.Error("matrices without include/exclude must not be flagged adjusted")
	}

	// deploy job has a multi-line if: -- make sure it survives parsing with its newline.
	if !strings.Contains(deploy.If, "\n") {
		t.Errorf("deploy If = %q, want multi-line", deploy.If)
	}
}

// TestParseMatrixMixedDynamic verifies that a matrix mixing a literal axis with an
// expression-valued axis keeps both: the literal axis lists its values and the
// expression axis shows the expression itself, since it cannot be enumerated statically.
func TestParseMatrixMixedDynamic(t *testing.T) {
	src := `name: CI
on: push
jobs:
  build:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        os: [ubuntu-latest]
        node: ${{ fromJSON(needs.detect.outputs.versions) }}
    steps:
      - run: build
`
	w, err := parseString(t, src)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	job := w.Jobs[0]
	if vals, ok := job.MatrixValues("os"); !ok || strings.Join(vals, ",") != "ubuntu-latest" {
		t.Errorf("os axis = %v (found=%v), want [ubuntu-latest]", vals, ok)
	}
	if vals, ok := job.MatrixValues("node"); !ok || strings.Join(vals, ",") != "${{ fromJSON(needs.detect.outputs.versions) }}" {
		t.Errorf("node axis = %v (found=%v), want the raw expression", vals, ok)
	}
	if job.MatrixAdjusted {
		t.Error("expression axes alone must not flag the matrix as adjusted")
	}
}

// TestParseMatrixIncludeOnly verifies that a matrix declared purely via include: derives
// its axes from the entries' keys and values, and is flagged as adjusted.
func TestParseMatrixIncludeOnly(t *testing.T) {
	src := `name: CI
on: push
jobs:
  build:
    runs-on: ${{ matrix.os }}
    strategy:
      matrix:
        include:
          - os: windows-latest
            jdk: 25
          - os: macos-latest
            jdk: 25
          - os: ubuntu-latest
            jdk: 17
    steps:
      - run: build
`
	w, err := parseString(t, src)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	job := w.Jobs[0]
	if vals, ok := job.MatrixValues("os"); !ok || strings.Join(vals, ",") != "windows-latest,macos-latest,ubuntu-latest" {
		t.Errorf("os axis = %v (found=%v), want the include entries' values deduped in order", vals, ok)
	}
	if vals, ok := job.MatrixValues("jdk"); !ok || strings.Join(vals, ",") != "25,17" {
		t.Errorf("jdk axis = %v (found=%v), want [25 17] (deduped)", vals, ok)
	}
	if !job.MatrixAdjusted {
		t.Error("include-only matrix must be flagged adjusted")
	}
}

// TestVersionComment covers the trailing-comment-to-version extraction rules.
func TestVersionComment(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"v4.1.1", "v4.1.1"},
		{"v4", "v4"},
		{"4.2.0", "4.2.0"},
		{"v4.1.1 pinned for reproducibility", "v4.1.1"},
		{"renovate: tag=v4", ""}, // not a version-leading comment
		{"pin to latest", ""},
		{"", ""},
	}
	for _, tt := range tests {
		if got := versionComment(tt.in); got != tt.want {
			t.Errorf("versionComment(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

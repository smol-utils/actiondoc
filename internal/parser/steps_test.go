package parser

import (
	"strings"
	"testing"

	"github.com/smol-utils/actiondoc/internal/model"
)

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

// TestParseMatrix covers static scalar axes, list-of-objects axes (flattened to dotted
// names), and the include/exclude fall-through that disables static resolution.
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

	// verify: include: present -> no static axes at all
	verify := w.Jobs[2]
	if len(verify.Matrix) != 0 {
		t.Errorf("verify matrix = %+v, want empty (include/exclude disables resolution)", verify.Matrix)
	}

	// deploy job has a multi-line if: -- make sure it survives parsing with its newline.
	if !strings.Contains(deploy.If, "\n") {
		t.Errorf("deploy If = %q, want multi-line", deploy.If)
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

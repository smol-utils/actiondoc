package renderer

import (
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/smol-utils/actiondoc/internal/model"
)

const shaPin = "8f4b7f84864484a7bf31766abe9204da3cbe65b3"

// TestStepTitleFallback covers the title chain: name -> id -> friendly uses -> first run
// line -> positional fallback.
func TestStepTitleFallback(t *testing.T) {
	tests := []struct {
		name string
		step model.Step
		want string
	}{
		{"explicit name wins", model.Step{Name: "Build", ID: "b", Uses: "actions/checkout@v4"}, "Build"},
		{"id over uses", model.Step{ID: "build-step", Uses: "actions/checkout@v4"}, "build-step"},
		{"sha-pinned uses with version", model.Step{Uses: "actions/checkout@" + shaPin, UsesVersion: "v4.1.1"}, "actions/checkout@v4.1.1"},
		{"sha-pinned uses without version", model.Step{Uses: "actions/checkout@" + shaPin}, "actions/checkout"},
		{"tag-pinned uses verbatim", model.Step{Uses: "actions/cache@v4"}, "actions/cache@v4"},
		{"first run line", model.Step{Run: "# setup\n\nmake build\nmake test"}, "make build"},
		{"long run line truncated", model.Step{Run: strings.Repeat("x", 80)}, strings.Repeat("x", 57) + "..."},
		{"positional fallback", model.Step{}, "Step 3"},
	}
	for _, tt := range tests {
		if got := stepTitle(&tt.step, 3); got != tt.want {
			t.Errorf("%s: stepTitle = %q, want %q", tt.name, got, tt.want)
		}
	}
}

// TestResolveJobName covers static matrix expansion in job names, dotted axes, and the
// verbatim fall-through for unresolvable or non-matrix expressions.
func TestResolveJobName(t *testing.T) {
	tests := []struct {
		name string
		job  model.Job
		want string
	}{
		{
			"static axis",
			model.Job{Name: "Java ${{ matrix.java }}", Matrix: []model.MatrixAxis{{Name: "java", Values: []string{"17", "21", "24"}}}},
			"Java 17, 21, 24",
		},
		{
			"dotted axis",
			model.Job{Name: "Deploy ${{ matrix.target.env }}", Matrix: []model.MatrixAxis{{Name: "target.env", Values: []string{"staging", "production"}}}},
			"Deploy staging, production",
		},
		{
			"two axes",
			model.Job{Name: "${{ matrix.os }} / ${{ matrix.arch }}", Matrix: []model.MatrixAxis{
				{Name: "os", Values: []string{"linux", "macos"}},
				{Name: "arch", Values: []string{"amd64"}},
			}},
			"linux, macos / amd64",
		},
		{
			"unresolvable axis stays verbatim",
			model.Job{Name: "Test ${{ matrix.case }}"},
			"Test ${{ matrix.case }}",
		},
		{
			"non-matrix expression stays verbatim",
			model.Job{Name: "Run for ${{ github.event.number }}"},
			"Run for ${{ github.event.number }}",
		},
		{
			"no expression",
			model.Job{Name: "Build"},
			"Build",
		},
	}
	for _, tt := range tests {
		if got := resolveJobName(&tt.job); got != tt.want {
			t.Errorf("%s: resolveJobName = %q, want %q", tt.name, got, tt.want)
		}
	}
}

// TestJobConditionMultilineEscaped verifies a multi-line job if: renders inside its table
// cell with <br> instead of raw newlines, so the table does not break.
func TestJobConditionMultilineEscaped(t *testing.T) {
	w := &model.Workflow{
		File: "test.yml",
		Name: "Test",
		On:   []string{"push"},
		Jobs: []model.Job{{
			ID:     "check",
			Name:   "check",
			RunsOn: "ubuntu-latest",
			If:     "github.event_name == 'push' &&\ncontains(github.event.pull_request.labels.*.name, 'ci')\n",
		}},
	}

	md := RenderMarkdown(w)

	if !strings.Contains(md, "'push' &&<br>contains(") {
		t.Errorf("multi-line condition not joined with <br>:\n%s", md)
	}
	// No table row may contain a raw newline: every line starting with | must also end with |.
	for _, line := range strings.Split(md, "\n") {
		if strings.HasPrefix(line, "|") && !strings.HasSuffix(line, "|") {
			t.Errorf("broken table row: %q", line)
		}
	}
	// The literal-block trailing newline must not leave a dangling <br> at the cell end.
	if strings.Contains(md, "<br>` |") {
		t.Errorf("dangling <br> at end of condition cell:\n%s", md)
	}
}

// TestRenderStepDetails covers the continue-on-error badge, with: rendering, the version
// annotation on SHA-pinned uses, and composite-action input doc enrichment.
func TestRenderStepDetails(t *testing.T) {
	action := &model.Action{
		Name: "Deploy",
		Inputs: []model.ActionInput{
			{Name: "environment", Description: "Target environment name", Required: true},
			{Name: "token", Description: "Deployment token"},
		},
	}
	step := model.Step{
		Name:            "Deploy to staging",
		Uses:            "./.github/actions/deploy",
		ContinueOnError: true,
		With: []model.KV{
			{Key: "environment", Value: "staging"},
			{Key: "token", Value: "${{ secrets.DEPLOY_TOKEN }}"},
			{Key: "undeclared", Value: "x"},
		},
		UsesAction: action,
	}

	var b strings.Builder
	renderStep(&b, &step, 1)
	got := b.String()

	checks := []string{
		"`[continue-on-error]`",
		"- `environment`: `staging` - Target environment name (required)",
		"- `token`: `${{ secrets.DEPLOY_TOKEN }}` - Deployment token",
		"- `undeclared`: `x`", // no doc suffix for keys the action does not declare
	}
	for _, want := range checks {
		if !strings.Contains(got, want) {
			t.Errorf("output missing %q\n\nFull output:\n%s", want, got)
		}
	}

	// SHA-pinned uses gets its version annotation on the detail line.
	var b2 strings.Builder
	renderStep(&b2, &model.Step{Uses: "actions/checkout@" + shaPin, UsesVersion: "v4.1.1"}, 1)
	if !strings.Contains(b2.String(), "`actions/checkout@"+shaPin+"` (v4.1.1)") {
		t.Errorf("missing version annotation on uses detail line:\n%s", b2.String())
	}
}

// TestRenderTOC covers the contents listing, the under-two-entries suppression, and
// duplicate-title anchor disambiguation.
func TestRenderTOC(t *testing.T) {
	if got := RenderTOC([]string{"Only One"}); got != "" {
		t.Errorf("single-entry TOC = %q, want empty", got)
	}
	if got := RenderTOC(nil); got != "" {
		t.Errorf("empty TOC = %q, want empty", got)
	}

	got := RenderTOC([]string{"CI Pipeline", "Release", "CI Pipeline"})
	checks := []string{
		"# Contents",
		"- [CI Pipeline](#ci-pipeline)\n",
		"- [Release](#release)\n",
		"- [CI Pipeline](#ci-pipeline-1)\n", // duplicate title gets the -1 anchor suffix
	}
	for _, want := range checks {
		if !strings.Contains(got, want) {
			t.Errorf("TOC missing %q\n\nFull output:\n%s", want, got)
		}
	}
}

// TestRunsOnEscapedInTable verifies a normalized runs-on list renders inside its table cell.
func TestRunsOnEscapedInTable(t *testing.T) {
	w := &model.Workflow{
		File: "test.yml",
		Name: "Test",
		On:   []string{"push"},
		Jobs: []model.Job{{ID: "build", Name: "build", RunsOn: "self-hosted, linux, x64"}},
	}
	md := RenderMarkdown(w)
	if !strings.Contains(md, "| Runs on | `self-hosted, linux, x64` |") {
		t.Errorf("runs-on row missing:\n%s", md)
	}
}

// TestTruncateRuneSafe verifies truncate counts and slices by rune, never splitting a
// multi-byte UTF-8 character (which would emit invalid UTF-8).
func TestTruncateRuneSafe(t *testing.T) {
	if got := truncate("hello", 10); got != "hello" {
		t.Errorf("short ASCII: got %q", got)
	}
	if got := truncate("abcdefghij", 8); got != "abcde..." {
		t.Errorf("ASCII cut: got %q, want %q", got, "abcde...")
	}
	// 10 multi-byte runes; truncating to 8 must yield 5 runes + "..." and stay valid UTF-8.
	s := "日本語のテストです字" // 10 runes
	got := truncate(s, 8)
	if !utf8.ValidString(got) {
		t.Errorf("truncate produced invalid UTF-8: %q", got)
	}
	if r := []rune(got); len(r) != 8 || string(r[5:]) != "..." {
		t.Errorf("rune cut: got %q (%d runes), want 5 runes + ...", got, len(r))
	}
}

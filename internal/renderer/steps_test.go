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

// TestJobNameDeTemplated locks the rule that a job name embedding ${{ ... }} expressions
// is rendered as a readable label rather than leaking the raw expression into the heading:
// each placeholder becomes its trailing identifier in parentheses. Matrix placeholders are
// still never expanded into joined value lists (GitHub creates one job per combination;
// "Java 17, 21" is a job name that never exists); the Matrix property row carries the axis
// values instead.
func TestJobNameDeTemplated(t *testing.T) {
	w := &model.Workflow{
		File: "test.yml",
		Name: "Test",
		On:   []string{"push"},
		Jobs: []model.Job{{
			ID:     "build",
			Name:   "Java ${{ matrix.java }} on ${{ matrix.os }}",
			RunsOn: "${{ matrix.os }}",
			Matrix: []model.MatrixAxis{
				{Name: "java", Values: []string{"17", "21", "24"}},
				{Name: "os", Values: []string{"ubuntu-latest", "macos-14"}},
			},
		}},
	}

	md := RenderMarkdown(w)

	// Heading: the de-templated label, never the raw expression or expanded values.
	heading := "### Java (java) on (os) (`build`)"
	if !strings.Contains(md, heading) {
		t.Errorf("job heading must show the de-templated label:\n%s", md)
	}
	// The heading line itself must not leak a raw ${{ ... }} expression (the runs-on cell
	// below it legitimately keeps the expression, so this is scoped to the heading line).
	for _, line := range strings.Split(md, "\n") {
		if strings.HasPrefix(line, "### ") && strings.Contains(line, "${{") {
			t.Errorf("job heading leaked a raw ${{ ... }} expression: %q", line)
		}
	}
	if strings.Contains(md, "Java 17, 21, 24") {
		t.Errorf("job heading must not expand matrix values:\n%s", md)
	}
	// The Matrix property row carries the axis values.
	if !strings.Contains(md, "| Matrix | `java`: 17, 21, 24; `os`: ubuntu-latest, macos-14 |") {
		t.Errorf("Matrix property row missing or wrong:\n%s", md)
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
		Env: []model.KV{
			{Key: "DEPLOY_REGION", Value: "us-east-1"},
			{Key: "API_KEY", Value: "${{ secrets.DEPLOY_API_KEY }}"},
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
		"   - Env:",
		"- `DEPLOY_REGION`: `us-east-1`",
		"- `API_KEY`: `${{ secrets.DEPLOY_API_KEY }}`",
	}
	for _, want := range checks {
		if !strings.Contains(got, want) {
			t.Errorf("output missing %q\n\nFull output:\n%s", want, got)
		}
	}

	// An unnamed SHA-pinned step is titled by its collapsed ref@version, so the redundant
	// full-SHA Uses: detail line is dropped entirely (the title already carries the version).
	var b2 strings.Builder
	renderStep(&b2, &model.Step{Uses: "actions/checkout@" + shaPin, UsesVersion: "v4.1.1"}, 1)
	if got := b2.String(); !strings.Contains(got, "**actions/checkout@v4.1.1**") || strings.Contains(got, shaPin) {
		t.Errorf("unnamed SHA pin should title as ref@version with no redundant Uses line:\n%s", got)
	}

	// A named SHA-pinned step keeps the Uses: line (so the action it runs is visible) but
	// collapses to ref@version -- the 40-char SHA is dropped when a version is known.
	var b3 strings.Builder
	renderStep(&b3, &model.Step{Name: "Checkout", Uses: "actions/checkout@" + shaPin, UsesVersion: "v4.1.1"}, 1)
	if got := b3.String(); !strings.Contains(got, "   - Uses: `actions/checkout@v4.1.1`\n") || strings.Contains(got, shaPin) {
		t.Errorf("named SHA pin should show collapsed ref@version on the Uses line:\n%s", got)
	}

	// A bare SHA pin (no known version) keeps its full ref on the Uses line so nothing is
	// lost: the title collapses to the bare ref, so the SHA only survives on the detail line.
	var b4 strings.Builder
	renderStep(&b4, &model.Step{Uses: "actions/checkout@" + shaPin}, 1)
	if got := b4.String(); !strings.Contains(got, "   - Uses: `actions/checkout@"+shaPin+"`\n") {
		t.Errorf("bare SHA pin should keep its full SHA on the Uses line:\n%s", got)
	}
}

// TestRenderStepOmitsEmptyWithEnv verifies that a step's `with:`/`env:` entries with an
// empty/unset value are omitted (rather than rendering as a bare dash), while entries with a
// real value are kept.
func TestRenderStepOmitsEmptyWithEnv(t *testing.T) {
	step := model.Step{
		Name: "Lock threads",
		Uses: "dessant/lock-threads@v5",
		With: []model.KV{
			{Key: "issue-inactive-days", Value: "30"},
			{Key: "add-issue-labels", Value: ""},
			{Key: "exclude-issue-created-before", Value: "  "},
		},
		Env: []model.KV{
			{Key: "REGION", Value: "us-east-1"},
			{Key: "UNSET", Value: ""},
		},
	}

	var b strings.Builder
	renderStep(&b, &step, 1)
	got := b.String()

	for _, want := range []string{
		"- `issue-inactive-days`: `30`",
		"- `REGION`: `us-east-1`",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("output should keep real value %q\n\nFull output:\n%s", want, got)
		}
	}
	for _, banned := range []string{
		"add-issue-labels",             // empty value omitted entirely
		"exclude-issue-created-before", // whitespace-only value omitted entirely
		"UNSET",                        // empty env value omitted entirely
		": -",                          // no bare-dash placeholder for empty values
	} {
		if strings.Contains(got, banned) {
			t.Errorf("output should omit empty entry %q\n\nFull output:\n%s", banned, got)
		}
	}
}

// TestRenderStepDropsEmptyWithHeader verifies that when every `with:` value is empty, the
// "With:" header itself is suppressed (no empty block).
func TestRenderStepDropsEmptyWithHeader(t *testing.T) {
	step := model.Step{
		Name: "All empty",
		Uses: "some/action@v1",
		With: []model.KV{
			{Key: "a", Value: ""},
			{Key: "b", Value: ""},
		},
	}

	var b strings.Builder
	renderStep(&b, &step, 1)
	if got := b.String(); strings.Contains(got, "With:") {
		t.Errorf("an all-empty with: block should drop its header:\n%s", got)
	}
}

// TestRenderTOC covers the grouped contents listing, the under-two-entries suppression,
// group headings, and that only non-empty groups render.
func TestRenderTOC(t *testing.T) {
	if got := RenderTOC([]TOCGroup{{Heading: "Workflows", Entries: []TOCEntry{{Label: "Only One", Anchor: "only-one"}}}}); got != "" {
		t.Errorf("single-entry TOC = %q, want empty", got)
	}
	if got := RenderTOC(nil); got != "" {
		t.Errorf("empty TOC = %q, want empty", got)
	}

	got := RenderTOC([]TOCGroup{
		{Heading: "Workflows", Entries: []TOCEntry{
			{Label: "CI Pipeline", Anchor: "ci-pipeline"},
			{Label: "Release - workflow_dispatch, push", Anchor: "release"},
		}},
		{Heading: "Reusable workflows"}, // empty: must not render
		{Heading: "Composite actions", Entries: []TOCEntry{
			{Label: "Setup", Anchor: "setup"},
		}},
	})
	checks := []string{
		"## Contents",
		"**Workflows**\n\n",
		"- [CI Pipeline](#ci-pipeline)\n",
		"- [Release - workflow_dispatch, push](#release)\n",
		"**Composite actions**\n\n",
		"- [Setup](#setup)\n",
	}
	for _, want := range checks {
		if !strings.Contains(got, want) {
			t.Errorf("TOC missing %q\n\nFull output:\n%s", want, got)
		}
	}
	if strings.Contains(got, "Reusable workflows") {
		t.Errorf("empty group should not render:\n%s", got)
	}
}

// TestRenderDocumentHeader covers title emission, zero-count omission, and pluralization.
func TestRenderDocumentHeader(t *testing.T) {
	got := RenderDocumentHeader("airflow", 12, 0, 1)
	if !strings.Contains(got, "# airflow\n\n") {
		t.Errorf("missing title:\n%s", got)
	}
	if !strings.Contains(got, "12 workflows, 1 composite action\n") {
		t.Errorf("inventory wrong (zero omitted, plural/singular):\n%s", got)
	}
	if strings.Contains(got, "reusable") {
		t.Errorf("zero count should be omitted:\n%s", got)
	}
	if got := RenderDocumentHeader("", 1, 0, 0); got != "" {
		t.Errorf("empty title should suppress header, got %q", got)
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

// TestStepTitleMarkupEscaped locks the seam discipline for step titles: backticks and
// asterisks in a step name or run-derived title must render literally, never as Markdown
// markup that breaks or restyles the surrounding bold.
func TestStepTitleMarkupEscaped(t *testing.T) {
	var b strings.Builder
	renderStep(&b, &model.Step{Run: "echo \"a `b` c\" | grep d"}, 1)
	got := b.String()
	if !strings.Contains(got, "1. **echo \"a \\`b\\` c\" | grep d**") {
		t.Errorf("run-derived title not escaped:\n%s", got)
	}

	var b2 strings.Builder
	renderStep(&b2, &model.Step{Name: "Run **everything** now"}, 1)
	if !strings.Contains(b2.String(), "**Run \\*\\*everything\\*\\* now**") {
		t.Errorf("step name with asterisks not escaped:\n%s", b2.String())
	}
}

// TestStepTitleSkipsPunctuationLines verifies the run-derived title skips lines with no
// letters or digits (a shell group's opening brace), and that unnamed uses: steps title
// with the collapsed pin form everywhere they are referenced.
func TestStepTitleSkipsPunctuationLines(t *testing.T) {
	if got := stepTitle(&model.Step{Run: "{\n  echo hello\n} > out.txt"}, 1); got != "echo hello" {
		t.Errorf("stepTitle = %q, want %q (brace-only line skipped)", got, "echo hello")
	}
	if got := stepTitle(&model.Step{Run: "{\n}\n"}, 4); got != "Step 4" {
		t.Errorf("stepTitle = %q, want positional fallback for punctuation-only script", got)
	}
	// Letters and digits in any script count as meaningful, not just ASCII.
	if got := stepTitle(&model.Step{Run: "{\n  输出结果\n}"}, 1); got != "输出结果" {
		t.Errorf("stepTitle = %q, want the non-Latin run line used as the title", got)
	}
	if got := stepTitle(&model.Step{Run: "(λ)"}, 1); got != "(λ)" {
		t.Errorf("stepTitle = %q, want a line with a non-ASCII letter kept", got)
	}
}

// TestUsedByCellGroupsByJob covers the "Used by" cell layout: usage sites group one line per
// job (joined with <br>), in first-seen job order; within a job the sites keep source order.
// A secret used across two jobs and two steps each must produce exactly two lines, each with
// its two `<step> (<name>)` entries. Job-level (no step), run/if, and workflow-env variants
// must render in their reduced forms.
func TestUsedByCellGroupsByJob(t *testing.T) {
	// A secret used by 2 jobs x 2 steps -> 2 lines, each with 2 entries in source order.
	sites := []model.Site{
		{Job: "build-cli", Step: "Setup Graal", Name: "github-token"},
		{Job: "build-cli", Step: "Checkout smoketests repository", Name: "token"},
		{Job: "build-tool", Step: "Setup Graal", Name: "github-token"},
		{Job: "build-tool", Step: "Checkout smoketests repository", Name: "token"},
	}
	want := "`build-cli`: Setup Graal (`github-token`), Checkout smoketests repository (`token`)" +
		"<br>`build-tool`: Setup Graal (`github-token`), Checkout smoketests repository (`token`)"
	if got := usedByCell(sites); got != want {
		t.Errorf("two-job/two-step cell:\n got = %q\nwant = %q", got, want)
	}

	// A job-level env site (no step) renders as just `(<name>)` under the job.
	if got := usedByCell([]model.Site{{Job: "release", Name: "API_KEY"}}); got != "`release`: (`API_KEY`)" {
		t.Errorf("job-level site = %q, want `release`: (`API_KEY`)", got)
	}

	// run/if sites carry the verb as the parenthesized name.
	runIf := []model.Site{
		{Job: "deploy", Step: "Sign", Name: "run"},
		{Job: "deploy", Step: "Gate", Name: "if"},
	}
	if got := usedByCell(runIf); got != "`deploy`: Sign (`run`), Gate (`if`)" {
		t.Errorf("run/if cell = %q", got)
	}

	// A workflow-level env site groups under a leading `workflow env` line.
	if got := usedByCell([]model.Site{{Name: "GLOBAL_TOKEN"}}); got != "workflow env: (`GLOBAL_TOKEN`)" {
		t.Errorf("workflow-env site = %q, want workflow env: (`GLOBAL_TOKEN`)", got)
	}
}

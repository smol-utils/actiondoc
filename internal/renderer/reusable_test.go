package renderer

import (
	"strings"
	"testing"

	"github.com/smol-utils/actiondoc/internal/callgraph"
	"github.com/smol-utils/actiondoc/internal/model"
)

// chainGraph builds the canonical three-level test graph:
//
//	release.yml [workflow_dispatch] -> middle.yml [workflow_call] -> leaf.yml [workflow_call]
//
// release.yml's caller job forwards an input and a secret; leaf.yml declares a @secret
// tag. Returns the graph plus the three workflows keyed by node id.
func chainGraph() (*callgraph.Graph, map[string]*model.Workflow) {
	release := &model.Workflow{
		File: "release.yml", Name: "Release", On: []string{"workflow_dispatch"},
		Jobs: []model.Job{{
			ID:   "publish",
			Uses: "./.github/workflows/middle.yml",
			With: []model.KV{{Key: "version", Value: "${{ inputs.version }}"}},
			Secrets: []model.KV{
				{Key: "GPG_KEY", Value: "${{ secrets.RELEASE_GPG_KEY }}"},
			},
		}},
	}
	middle := &model.Workflow{
		File: "middle.yml", Name: "Middle", On: []string{"workflow_call"},
		Jobs: []model.Job{{
			ID:   "build",
			Uses: "./.github/workflows/leaf.yml",
		}},
	}
	leaf := &model.Workflow{
		File: "leaf.yml", Name: "Leaf", On: []string{"workflow_call"},
		Tags: model.Tags{Secrets: []model.Param{{Name: "SIGNING_KEY"}}},
		Jobs: []model.Job{{
			ID: "compile", RunsOn: "ubuntu-latest",
			Permissions: &model.Permissions{Scopes: []model.Permission{
				{Scope: "contents", Level: "read"},
				{Scope: "id-token", Level: "write", OIDC: true},
			}},
			Steps: []model.Step{{
				Name: "Publish",
				Run:  `./publish.sh --region "${{ vars.DEPLOY_REGION }}" --token "${{ secrets.LEAF_TOKEN }}"`,
			}},
		}},
	}
	workflows := map[string]*model.Workflow{
		".github/workflows/release.yml": release,
		".github/workflows/middle.yml":  middle,
		".github/workflows/leaf.yml":    leaf,
	}
	var sources []callgraph.Source
	for path, w := range workflows {
		sources = append(sources, callgraph.Source{Path: path, Workflow: w})
	}
	return callgraph.Build(sources), workflows
}

// TestRenderCallerJobMatrixRow verifies a caller job's matrix axes render as a Matrix
// property row (a caller's matrix multiplies the reusable calls), including the
// include/exclude adjustment note.
func TestRenderCallerJobMatrixRow(t *testing.T) {
	g, workflows := chainGraph()
	id := ".github/workflows/release.yml"
	job := &workflows[id].Jobs[0]
	job.Matrix = []model.MatrixAxis{
		{Name: "os", Values: []string{"linux", "windows", "darwin"}},
		{Name: "arch", Values: []string{"amd64", "arm64"}},
	}
	job.MatrixAdjusted = true

	md := RenderMarkdownGraph(workflows[id], g, id)

	want := "| Matrix | `os`: linux, windows, darwin; `arch`: amd64, arm64 (combinations adjusted by include/exclude) |"
	if !strings.Contains(md, want) {
		t.Errorf("caller job missing Matrix row %q:\n%s", want, md)
	}
}

func TestRenderCallerJobForwarding(t *testing.T) {
	g, workflows := chainGraph()
	id := ".github/workflows/release.yml"

	md := RenderMarkdownGraph(workflows[id], g, id)

	checks := []string{
		"| Uses workflow | [Middle](#middle) |",
		"#### Inputs forwarded",
		"- `version`: `${{ inputs.version }}`",
		"#### Secrets forwarded",
		"- `GPG_KEY`: `${{ secrets.RELEASE_GPG_KEY }}`",
	}
	for _, want := range checks {
		if !strings.Contains(md, want) {
			t.Errorf("output missing %q\n\nFull output:\n%s", want, md)
		}
	}
	// Caller jobs have no steps and must not render an empty collapsed Steps block.
	if strings.Contains(md, "<summary>Steps") {
		t.Errorf("caller job rendered a Steps section:\n%s", md)
	}
}

// TestRenderCallerJobMultilineValue verifies a multi-line forwarded with:/secrets: value
// is collapsed to one line, so it can't inject newlines that break the Markdown list.
func TestRenderCallerJobMultilineValue(t *testing.T) {
	caller := &model.Workflow{
		File: "release.yml", Name: "Release", On: []string{"workflow_dispatch"},
		Jobs: []model.Job{{
			ID:   "publish",
			Uses: "./.github/workflows/build.yml",
			With: []model.KV{{Key: "config", Value: "line1\nline2\nline3"}},
		}},
	}
	build := &model.Workflow{File: "build.yml", Name: "Build", On: []string{"workflow_call"}}
	g := callgraph.Build([]callgraph.Source{
		{Path: ".github/workflows/release.yml", Workflow: caller},
		{Path: ".github/workflows/build.yml", Workflow: build},
	})

	md := RenderMarkdownGraph(caller, g, ".github/workflows/release.yml")
	if strings.Contains(md, "line1\nline2") {
		t.Errorf("multi-line forwarded value not collapsed (raw newline present):\n%s", md)
	}
	if !strings.Contains(md, "`config`: `line1 line2 line3`") {
		t.Errorf("expected collapsed forwarded value, got:\n%s", md)
	}
}

// TestRenderCallerJobOmitsEmptyForwardedInput verifies a forwarded `with:` entry whose value
// is empty/unset is omitted (no bare-dash row), while an entry with a real value is kept. When
// every forwarded input is empty, the "Inputs forwarded" header is dropped entirely.
func TestRenderCallerJobOmitsEmptyForwardedInput(t *testing.T) {
	caller := &model.Workflow{
		File: "release.yml", Name: "Release", On: []string{"workflow_dispatch"},
		Jobs: []model.Job{{
			ID:   "publish",
			Uses: "./.github/workflows/build.yml",
			With: []model.KV{
				{Key: "version", Value: "${{ inputs.version }}"},
				{Key: "test-name-separator", Value: ""},
				{Key: "blank", Value: "   "},
			},
		}},
	}
	build := &model.Workflow{File: "build.yml", Name: "Build", On: []string{"workflow_call"}}
	g := callgraph.Build([]callgraph.Source{
		{Path: ".github/workflows/release.yml", Workflow: caller},
		{Path: ".github/workflows/build.yml", Workflow: build},
	})

	md := RenderMarkdownGraph(caller, g, ".github/workflows/release.yml")

	if !strings.Contains(md, "- `version`: `${{ inputs.version }}`") {
		t.Errorf("real forwarded input dropped:\n%s", md)
	}
	if strings.Contains(md, "test-name-separator") {
		t.Errorf("empty forwarded input rendered:\n%s", md)
	}
	if strings.Contains(md, "`blank`") {
		t.Errorf("whitespace-only forwarded input rendered:\n%s", md)
	}

	// And when every forwarded input is empty, the header itself is suppressed.
	allEmpty := &model.Workflow{
		File: "release.yml", Name: "Release", On: []string{"workflow_dispatch"},
		Jobs: []model.Job{{
			ID:   "publish",
			Uses: "./.github/workflows/build.yml",
			With: []model.KV{{Key: "test-name-separator", Value: ""}},
		}},
	}
	g2 := callgraph.Build([]callgraph.Source{
		{Path: ".github/workflows/release.yml", Workflow: allEmpty},
		{Path: ".github/workflows/build.yml", Workflow: build},
	})
	md2 := RenderMarkdownGraph(allEmpty, g2, ".github/workflows/release.yml")
	if strings.Contains(md2, "#### Inputs forwarded") {
		t.Errorf("Inputs forwarded header rendered with no meaningful inputs:\n%s", md2)
	}
}

// TestRenderCallerJobTags verifies that ActionDoc tags declared on a reusable-workflow
// caller job (@secret/@env/@output/@example/@see) are rendered, not dropped on the early
// return.
func TestRenderCallerJobTags(t *testing.T) {
	caller := &model.Workflow{
		File: "release.yml", Name: "Release", On: []string{"workflow_dispatch"},
		Jobs: []model.Job{{
			ID:   "publish",
			Uses: "./.github/workflows/build.yml",
			Tags: model.Tags{
				Secrets: []model.Param{{Name: "DEPLOY_KEY", Description: "deploy key"}},
				Envs:    []model.Param{{Name: "REGION"}},
				Example: "gh workflow run release.yml",
				See:     []string{"https://example.com/runbook"},
			},
		}},
	}
	build := &model.Workflow{File: "build.yml", Name: "Build", On: []string{"workflow_call"}}
	g := callgraph.Build([]callgraph.Source{
		{Path: ".github/workflows/release.yml", Workflow: caller},
		{Path: ".github/workflows/build.yml", Workflow: build},
	})

	md := RenderMarkdownGraph(caller, g, ".github/workflows/release.yml")
	for _, want := range []string{
		"**Secrets:**", "`DEPLOY_KEY`",
		"**Environment Variables:**", "`REGION`",
		"**Example:**", "gh workflow run release.yml",
		"**See also:**", "https://example.com/runbook",
	} {
		if !strings.Contains(md, want) {
			t.Errorf("caller-job tags missing %q\n\nFull output:\n%s", want, md)
		}
	}
}

func TestRenderCallerJobSecretsInherit(t *testing.T) {
	w := &model.Workflow{
		File: "caller.yml", Name: "Caller", On: []string{"push"},
		Jobs: []model.Job{{
			ID:             "deploy",
			Uses:           "./.github/workflows/deploy.yml",
			SecretsInherit: true,
		}},
	}
	callee := &model.Workflow{
		File: "deploy.yml", Name: "Deploy", On: []string{"workflow_call"},
	}
	g := callgraph.Build([]callgraph.Source{
		{Path: ".github/workflows/caller.yml", Workflow: w},
		{Path: ".github/workflows/deploy.yml", Workflow: callee},
	})

	md := RenderMarkdownGraph(w, g, ".github/workflows/caller.yml")

	if !strings.Contains(md, "`secrets: inherit`") {
		t.Errorf("secrets: inherit not rendered verbatim:\n%s", md)
	}
}

func TestRenderCallerJobExternalCallee(t *testing.T) {
	w := &model.Workflow{
		File: "caller.yml", Name: "Caller", On: []string{"push"},
		Jobs: []model.Job{{
			ID:   "lint",
			Uses: "other-org/shared/.github/workflows/lint.yml@v3",
		}},
	}
	g := callgraph.Build([]callgraph.Source{
		{Path: ".github/workflows/caller.yml", Workflow: w},
	})

	md := RenderMarkdownGraph(w, g, ".github/workflows/caller.yml")

	// External callees render as inline code with the pin, marked external -- never as
	// an anchor link (there is no in-scope section to link to).
	if !strings.Contains(md, "`other-org/shared/.github/workflows/lint.yml@v3` (external)") {
		t.Errorf("external callee not rendered with pin and external marker:\n%s", md)
	}
	if strings.Contains(md, "[other-org") {
		t.Errorf("external callee must not be an anchor link:\n%s", md)
	}
}

func TestRenderCallerJobWithoutGraph(t *testing.T) {
	w := &model.Workflow{
		File: "caller.yml", Name: "Caller", On: []string{"push"},
		Jobs: []model.Job{{
			ID:   "publish",
			Uses: "./.github/workflows/publish.yml",
			With: []model.KV{{Key: "version", Value: "1.0"}},
		}},
	}

	// Single-file rendering (nil graph): the caller surface still renders, with the raw
	// uses: reference instead of a cross-link.
	md := RenderMarkdown(w)

	checks := []string{
		"| Uses workflow | `./.github/workflows/publish.yml` |",
		"- `version`: `1.0`",
	}
	for _, want := range checks {
		if !strings.Contains(md, want) {
			t.Errorf("output missing %q\n\nFull output:\n%s", want, md)
		}
	}
}

func TestCallGraphOnEntryPoint(t *testing.T) {
	g, workflows := chainGraph()
	id := ".github/workflows/release.yml"

	md := RenderMarkdownGraph(workflows[id], g, id)

	if !strings.Contains(md, "## Call graph (rooted at this workflow)") {
		t.Fatalf("missing call graph section:\n%s", md)
	}
	checks := []string{
		"- `publish` uses [middle.yml](#middle)",
		"  - `build` uses [leaf.yml](#leaf)",
	}
	for _, want := range checks {
		if !strings.Contains(md, want) {
			t.Errorf("call graph missing %q\n\nFull output:\n%s", want, md)
		}
	}
	// The tree starts directly at its children: it must not restate the root workflow's
	// file + triggers, which the section heading and property table already show.
	if strings.Contains(md, "`release.yml` [workflow_dispatch]") {
		t.Errorf("call graph should not restate the root file/triggers line:\n%s", md)
	}
	// The list rendering must not fall back to the old fenced ASCII tree.
	if strings.Contains(md, "+-- ") {
		t.Errorf("call graph still rendering ASCII tree:\n%s", md)
	}
}

func TestCallGraphSuppressedOnReusableWorkflow(t *testing.T) {
	g, workflows := chainGraph()
	// middle.yml calls leaf.yml but is workflow_call-only, so it is not an entry point
	// and gets a "Called by" section instead of a call graph.
	id := ".github/workflows/middle.yml"

	md := RenderMarkdownGraph(workflows[id], g, id)

	if strings.Contains(md, "## Call graph") {
		t.Errorf("call graph must not render on a reusable (non-entry-point) workflow:\n%s", md)
	}
	if !strings.Contains(md, "## Called by") {
		t.Errorf("missing Called by section on reusable workflow:\n%s", md)
	}
}

func TestCallGraphSuppressedOnFlatWorkflow(t *testing.T) {
	w := &model.Workflow{
		File: "ci.yml", Name: "CI", On: []string{"push"},
		Jobs: []model.Job{{ID: "test", RunsOn: "ubuntu-latest"}},
	}
	g := callgraph.Build([]callgraph.Source{
		{Path: ".github/workflows/ci.yml", Workflow: w},
	})

	md := RenderMarkdownGraph(w, g, ".github/workflows/ci.yml")

	for _, section := range []string{"## Call graph", "## Called by", "## Transitive requirements"} {
		if strings.Contains(md, section) {
			t.Errorf("flat workflow must not render %q:\n%s", section, md)
		}
	}
}

func TestCalledByTransitiveChain(t *testing.T) {
	g, workflows := chainGraph()
	id := ".github/workflows/leaf.yml"

	md := RenderMarkdownGraph(workflows[id], g, id)

	if !strings.Contains(md, "## Called by") {
		t.Fatalf("missing Called by section:\n%s", md)
	}
	checks := []string{
		"- [middle.yml](#middle) (job: `build`)",
		"  - [release.yml](#release) (job: `publish`) - entry point",
	}
	for _, want := range checks {
		if !strings.Contains(md, want) {
			t.Errorf("Called by chain missing %q\n\nFull output:\n%s", want, md)
		}
	}
	// The tree starts directly at its callers: it must not restate the root workflow's
	// file basename, which the section heading already names.
	if strings.Contains(md, "## Called by\n\n`leaf.yml`") {
		t.Errorf("Called by tree should not restate the root file basename:\n%s", md)
	}
}

func TestTransitiveRequirements(t *testing.T) {
	g, workflows := chainGraph()
	id := ".github/workflows/release.yml"

	md := RenderMarkdownGraph(workflows[id], g, id)

	const heading = "## Transitive requirements (from full call graph)"
	start := strings.Index(md, heading)
	if start < 0 {
		t.Fatalf("missing transitive requirements section:\n%s", md)
	}
	// Scope the assertions to just this section: a name dropped here may still legitimately
	// appear elsewhere (e.g. RELEASE_GPG_KEY in the per-workflow Used-by table, which stays).
	section := md[start:]
	if next := strings.Index(section[len(heading):], "\n## "); next >= 0 {
		section = section[:len(heading)+next]
	}

	// This section is the contract view: it aggregates only the DECLARED/forwarded secret
	// names. GPG_KEY comes from the entry point's forwarded secrets: mapping key; SIGNING_KEY
	// from the leaf's @secret tag two hops down. Names are sorted alphabetically.
	if !strings.Contains(section, "Secrets required (declared/forwarded names): `GPG_KEY`, `SIGNING_KEY`") {
		t.Errorf("declared secrets not aggregated across the chain:\n%s", section)
	}
	// Expression-scanned NAMES are intentionally NOT listed here anymore (the document-level
	// inventory and the per-workflow Used-by table cover those). LEAF_TOKEN and
	// RELEASE_GPG_KEY are only reachable by re-scanning expressions, so they must be absent
	// from this section, and the variables line (DEPLOY_REGION was scan-only) must not appear.
	for _, gone := range []string{"LEAF_TOKEN", "RELEASE_GPG_KEY", "Variables referenced", "DEPLOY_REGION"} {
		if strings.Contains(section, gone) {
			t.Errorf("transitive section must not contain expression-scanned %q:\n%s", gone, section)
		}
	}
	// Permission grants are no longer rolled up here: they are shown per-workflow, per-job, and
	// repo-wide, so the transitive line would only repeat them. The declared-secrets line (and
	// any external-workflows line) is what this section keeps.
	if strings.Contains(section, "Permissions declared across the chain") {
		t.Errorf("transitive section must not carry a permissions roll-up line:\n%s", section)
	}
}

func TestCallGraphCycleTerminates(t *testing.T) {
	a := &model.Workflow{
		File: "a.yml", Name: "A", On: []string{"push"},
		Jobs: []model.Job{{ID: "call-b", Uses: "./.github/workflows/b.yml"}},
	}
	bw := &model.Workflow{
		File: "b.yml", Name: "B", On: []string{"workflow_call"},
		Jobs: []model.Job{{ID: "call-a", Uses: "./.github/workflows/a.yml"}},
	}
	g := callgraph.Build([]callgraph.Source{
		{Path: ".github/workflows/a.yml", Workflow: a},
		{Path: ".github/workflows/b.yml", Workflow: bw},
	})

	// Both directions must terminate despite the a -> b -> a cycle.
	mdA := RenderMarkdownGraph(a, g, ".github/workflows/a.yml")
	mdB := RenderMarkdownGraph(bw, g, ".github/workflows/b.yml")

	if !strings.Contains(mdA, "## Call graph") {
		t.Errorf("entry point in cycle missing call graph:\n%s", mdA)
	}
	if !strings.Contains(mdB, "## Called by") {
		t.Errorf("reusable workflow in cycle missing Called by:\n%s", mdB)
	}
}

func TestRenderTreeShape(t *testing.T) {
	root := treeNode{
		label: "root",
		children: []treeNode{
			{label: "first", children: []treeNode{{label: "first-child"}}},
			{label: "last", children: []treeNode{{label: "last-child"}}},
		},
	}

	var b strings.Builder
	renderTreeList(&b, root)

	// The root is a plain lead line (a paragraph), then a nested `-` list with two spaces
	// of indentation per depth.
	want := strings.Join([]string{
		"root",
		"",
		"- first",
		"  - first-child",
		"- last",
		"  - last-child",
		"",
	}, "\n")
	if b.String() != want {
		t.Errorf("tree shape mismatch.\nGot:\n%s\nWant:\n%s", b.String(), want)
	}
}

// TestCallGraphExternalCalleeNoLink verifies a downstream external callee renders as plain
// inline code carrying its pin, never as an anchor link (there is no in-scope section).
func TestCallGraphExternalCalleeNoLink(t *testing.T) {
	w := &model.Workflow{
		File: "ci.yml", Name: "CI", On: []string{"push"},
		Jobs: []model.Job{{
			ID:   "lint",
			Uses: "other-org/shared/.github/workflows/lint.yml@v3",
		}},
	}
	g := callgraph.Build([]callgraph.Source{
		{Path: ".github/workflows/ci.yml", Workflow: w},
	})

	md := RenderMarkdownGraph(w, g, ".github/workflows/ci.yml")

	if !strings.Contains(md, "`lint` uses `other-org/shared/.github/workflows/lint.yml@v3`") {
		t.Errorf("external callee not rendered as plain code with pin:\n%s", md)
	}
	if strings.Contains(md, "[other-org") {
		t.Errorf("external callee must not be an anchor link:\n%s", md)
	}
}

// TestCallGraphOutsideScopeNoLink verifies a downstream callee whose target was not scanned
// renders as plain inline code annotated "(outside scan scope)", never as a link.
func TestCallGraphOutsideScopeNoLink(t *testing.T) {
	w := &model.Workflow{
		File: "ci.yml", Name: "CI", On: []string{"push"},
		Jobs: []model.Job{{
			ID:   "deploy",
			Uses: "./.github/workflows/missing.yml",
		}},
	}
	g := callgraph.Build([]callgraph.Source{
		{Path: ".github/workflows/ci.yml", Workflow: w},
	})

	md := RenderMarkdownGraph(w, g, ".github/workflows/ci.yml")

	if !strings.Contains(md, "`deploy` uses `./.github/workflows/missing.yml` (outside scan scope)") {
		t.Errorf("outside-scope callee not rendered as plain annotated code:\n%s", md)
	}
	if strings.Contains(md, "[./.github/workflows/missing.yml]") {
		t.Errorf("outside-scope callee must not be an anchor link:\n%s", md)
	}
}

// TestCalledByJobAnchorLink verifies the upstream caller links to the specific calling job
// heading via the document-wide JobAnchors, not the file's section anchor.
func TestCalledByJobAnchorLink(t *testing.T) {
	g, _ := chainGraph()
	// Simulate the assembler's document-wide job-anchor pass: middle.yml's single job
	// "build" lands on a disambiguated heading slug elsewhere in the document.
	g.Nodes[".github/workflows/middle.yml"].JobAnchors = []string{"build-7"}

	md := RenderMarkdownGraph(g.Nodes[".github/workflows/leaf.yml"].Workflow, g, ".github/workflows/leaf.yml")

	if !strings.Contains(md, "- [middle.yml](#build-7) (job: `build`)") {
		t.Errorf("upstream caller did not link to the job anchor:\n%s", md)
	}
}

// TestCallGraphSiblingCollapse verifies that repeated sibling subtrees (distinct caller jobs
// invoking the same callee with identical descendants) fold to one "(xN)" representative
// whose label drops the differing job ids.
func TestCallGraphSiblingCollapse(t *testing.T) {
	caller := &model.Workflow{
		File: "ci.yml", Name: "CI", On: []string{"push"},
		Jobs: []model.Job{
			{ID: "test-a", Uses: "./.github/workflows/shared.yml"},
			{ID: "test-b", Uses: "./.github/workflows/shared.yml"},
			{ID: "test-c", Uses: "./.github/workflows/shared.yml"},
		},
	}
	shared := &model.Workflow{File: "shared.yml", Name: "Shared", On: []string{"workflow_call"}}
	g := callgraph.Build([]callgraph.Source{
		{Path: ".github/workflows/ci.yml", Workflow: caller},
		{Path: ".github/workflows/shared.yml", Workflow: shared},
	})

	md := RenderMarkdownGraph(caller, g, ".github/workflows/ci.yml")

	if !strings.Contains(md, "- uses **[shared.yml](#shared)** (x3)") {
		t.Errorf("repeated callee siblings did not collapse to (x3):\n%s", md)
	}
	for _, jobID := range []string{"test-a", "test-b", "test-c"} {
		if strings.Contains(md, "`"+jobID+"` uses") {
			t.Errorf("collapsed representative should drop job id %q:\n%s", jobID, md)
		}
	}
}

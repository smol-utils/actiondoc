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
	// Caller jobs have no steps and must not render an empty Steps section.
	if strings.Contains(md, "#### Steps") {
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
		"release.yml [workflow_dispatch]",
		"+-- publish (uses middle.yml)",
		"    +-- build (uses leaf.yml)",
	}
	for _, want := range checks {
		if !strings.Contains(md, want) {
			t.Errorf("call graph missing %q\n\nFull output:\n%s", want, md)
		}
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
		"+-- middle.yml (job: build)",
		"    +-- release.yml (job: publish)  <- entry point",
	}
	for _, want := range checks {
		if !strings.Contains(md, want) {
			t.Errorf("Called by chain missing %q\n\nFull output:\n%s", want, md)
		}
	}
}

func TestTransitiveRequirements(t *testing.T) {
	g, workflows := chainGraph()
	id := ".github/workflows/release.yml"

	md := RenderMarkdownGraph(workflows[id], g, id)

	if !strings.Contains(md, "## Transitive requirements (from full call graph)") {
		t.Fatalf("missing transitive requirements section:\n%s", md)
	}
	// GPG_KEY comes from the entry point's forwarded secrets: mapping; SIGNING_KEY from
	// the leaf's @secret tag two hops down; LEAF_TOKEN from scanning the leaf's run:
	// expression; RELEASE_GPG_KEY from scanning the entry point's forwarded secret value.
	// Names are sorted alphabetically.
	if !strings.Contains(md, "Secrets referenced (literal names): `GPG_KEY`, `LEAF_TOKEN`, `RELEASE_GPG_KEY`, `SIGNING_KEY`") {
		t.Errorf("secrets not aggregated across the chain:\n%s", md)
	}
	// DEPLOY_REGION comes from scanning the leaf's run: expression two hops down.
	if !strings.Contains(md, "Variables referenced: `DEPLOY_REGION`") {
		t.Errorf("variables not aggregated across the chain:\n%s", md)
	}
	// The leaf job's permission grants surface on the entry point, with the OIDC marker.
	if !strings.Contains(md, "Permissions declared across the chain: `contents: read`, `id-token: write (OIDC)`") {
		t.Errorf("permissions not aggregated across the chain:\n%s", md)
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
	renderTree(&b, root)

	want := strings.Join([]string{
		"root",
		"+-- first",
		"|   +-- first-child",
		"+-- last",
		"    +-- last-child",
		"",
	}, "\n")
	if b.String() != want {
		t.Errorf("tree shape mismatch.\nGot:\n%s\nWant:\n%s", b.String(), want)
	}
}

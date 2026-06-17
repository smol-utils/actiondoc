package renderer

import (
	"strings"
	"testing"

	"github.com/smol-utils/actiondoc/internal/callgraph"
	"github.com/smol-utils/actiondoc/internal/model"
)

// triggerIndexGraph builds three in-memory workflows whose triggers overlap, so the index has
// to aggregate multiple workflows under one event, dedupe a workflow that lists the same event
// twice, and exclude the workflow_call event. It returns the built graph and the source slice
// in the order the document assembler would pass them.
func triggerIndexGraph() ([]callgraph.Source, *callgraph.Graph) {
	build := func(file, name string, on []string) *model.Workflow {
		return &model.Workflow{
			File: file, Name: name, On: on,
			Jobs: []model.Job{{ID: "run", RunsOn: "ubuntu-latest",
				Steps: []model.Step{{Name: "Run", Run: "echo hi"}}}},
		}
	}
	alpha := build("alpha.yml", "Alpha", []string{"push", "pull_request"})
	beta := build("beta.yml", "Beta", []string{"pull_request", "workflow_call"})
	// gamma lists push twice (e.g. the canonical flat list repeated a key) -- it must be
	// attributed to push exactly once.
	gamma := build("gamma.yml", "Gamma", []string{"push", "push"})
	sources := []callgraph.Source{
		{Path: ".github/workflows/alpha.yml", Workflow: alpha},
		{Path: ".github/workflows/beta.yml", Workflow: beta},
		{Path: ".github/workflows/gamma.yml", Workflow: gamma},
	}
	return sources, callgraph.Build(sources, "")
}

func TestRenderTriggerIndex(t *testing.T) {
	sources, g := triggerIndexGraph()
	out := RenderTriggerIndex(sources, g)

	if !strings.Contains(out, triggerIndexHeading) {
		t.Errorf("missing heading %q\n\n%s", triggerIndexHeading, out)
	}

	// push aggregates Alpha and Gamma; Gamma's duplicate push entry is deduped to one link.
	wantPush := "- **push**: [Alpha](#alpha), [Gamma](#gamma)"
	if !strings.Contains(out, wantPush) {
		t.Errorf("missing aggregated/deduped push row %q\n\n%s", wantPush, out)
	}

	// pull_request aggregates Alpha and Beta.
	wantPR := "- **pull_request**: [Alpha](#alpha), [Beta](#beta)"
	if !strings.Contains(out, wantPR) {
		t.Errorf("missing aggregated pull_request row %q\n\n%s", wantPR, out)
	}

	// workflow_call must never appear as an event row -- reusable workflows have their own group.
	if strings.Contains(out, "**workflow_call**") {
		t.Errorf("workflow_call should be excluded from the trigger index\n\n%s", out)
	}
}

// TestRenderTriggerIndexOrdering verifies events are ordered by descending workflow count, then
// alphabetically for ties. Here push (2) and pull_request (2) tie, so they sort alphabetically
// (pull_request before push), and the single-workflow schedule event comes last.
func TestRenderTriggerIndexOrdering(t *testing.T) {
	sources, g := triggerIndexGraph()
	// Add a fourth workflow on a less-common single event so the ordering has a clear tail.
	sched := &model.Workflow{
		File: "delta.yml", Name: "Delta", On: []string{"schedule"},
		Jobs: []model.Job{{ID: "run", RunsOn: "ubuntu-latest",
			Steps: []model.Step{{Name: "Run", Run: "echo hi"}}}},
	}
	sources = append(sources, callgraph.Source{Path: ".github/workflows/delta.yml", Workflow: sched})
	g = callgraph.Build(sources, "")
	out := RenderTriggerIndex(sources, g)

	posPR := strings.Index(out, "- **pull_request**")
	posPush := strings.Index(out, "- **push**")
	posSched := strings.Index(out, "- **schedule**")
	if posPR < 0 || posPush < 0 || posSched < 0 {
		t.Fatalf("missing expected event rows\n\n%s", out)
	}
	// pull_request (count 2) and push (count 2) tie -> alphabetical: pull_request first.
	if !(posPR < posPush) {
		t.Errorf("expected pull_request before push (alphabetical tie-break)\n\n%s", out)
	}
	// schedule (count 1) sorts after the two count-2 events.
	if !(posPush < posSched) {
		t.Errorf("expected schedule (lower count) last\n\n%s", out)
	}
}

// TestRenderTriggerIndexDuplicateNames verifies that when two workflows share a display name,
// their link labels are disambiguated with the filename so the otherwise identical link text
// (against distinct anchors) stays distinguishable.
func TestRenderTriggerIndexDuplicateNames(t *testing.T) {
	mk := func(file string) *model.Workflow {
		return &model.Workflow{
			File: file, Name: "Release", On: []string{"push"},
			Jobs: []model.Job{{ID: "build", RunsOn: "ubuntu-latest",
				Steps: []model.Step{{Name: "Run", Run: "echo hi"}}}},
		}
	}
	sources := []callgraph.Source{
		{Path: ".github/workflows/a.yml", Workflow: mk("a.yml")},
		{Path: ".github/workflows/b.yml", Workflow: mk("b.yml")},
	}
	out := RenderTriggerIndex(sources, callgraph.Build(sources, ""))

	for _, want := range []string{"[Release (a.yml)]", "[Release (b.yml)]"} {
		if !strings.Contains(out, want) {
			t.Errorf("missing disambiguated link label %q\n\n%s", want, out)
		}
	}
}

// TestRenderTriggerIndexWorkflowCallOnly verifies the section is suppressed when the only
// trigger present is workflow_call: after excluding it, no events remain.
func TestRenderTriggerIndexWorkflowCallOnly(t *testing.T) {
	w := &model.Workflow{
		File: "reusable.yml", Name: "Reusable", On: []string{"workflow_call"},
		Jobs: []model.Job{{ID: "run", RunsOn: "ubuntu-latest",
			Steps: []model.Step{{Name: "Run", Run: "echo hi"}}}},
	}
	sources := []callgraph.Source{{Path: ".github/workflows/reusable.yml", Workflow: w}}
	if out := RenderTriggerIndex(sources, callgraph.Build(sources, "")); out != "" {
		t.Errorf("expected empty index when only workflow_call is present, got:\n%s", out)
	}
}

// TestRenderTriggerIndexEmpty verifies the section is suppressed when there are no trigger
// events at all (e.g. an action-only or empty source set).
func TestRenderTriggerIndexEmpty(t *testing.T) {
	w := &model.Workflow{
		File: "plain.yml", Name: "Plain", On: nil,
		Jobs: []model.Job{{ID: "run", RunsOn: "ubuntu-latest",
			Steps: []model.Step{{Name: "Run", Run: "echo hi"}}}},
	}
	sources := []callgraph.Source{{Path: ".github/workflows/plain.yml", Workflow: w}}
	if out := RenderTriggerIndex(sources, callgraph.Build(sources, "")); out != "" {
		t.Errorf("expected empty index with no trigger events, got:\n%s", out)
	}
}

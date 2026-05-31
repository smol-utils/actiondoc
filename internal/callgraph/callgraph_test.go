package callgraph

import (
	"testing"

	"github.com/smol-utils/actiondoc/internal/model"
)

func TestBuildBasic(t *testing.T) {
	entry := &model.Workflow{
		File: "release.yml", Name: "Release", On: []string{"workflow_dispatch"},
		Jobs: []model.Job{
			{ID: "publish", Uses: "./.github/workflows/build.yml",
				With:    []model.KV{{Key: "version", Value: "${{ inputs.version }}"}},
				Secrets: []model.KV{{Key: "CONSUMER-KEY", Value: "${{ secrets.SDKMAN_KEY }}"}}},
			{ID: "external", Uses: "owner/repo/.github/workflows/linux_job.yml@main"},
		},
	}
	reusable := &model.Workflow{File: "build.yml", Name: "Build", On: []string{"workflow_call"}}

	g := Build([]Source{
		{Path: "/repo/.github/workflows/release.yml", Workflow: entry},
		{Path: "/repo/.github/workflows/build.yml", Workflow: reusable},
	})

	relID := "/repo/.github/workflows/release.yml"
	buildID := "/repo/.github/workflows/build.yml"

	if !g.IsEntryPoint(relID) {
		t.Error("release.yml should be an entry point")
	}
	if g.IsEntryPoint(buildID) {
		t.Error("build.yml is workflow_call-only; not an entry point")
	}
	if got := g.CallerCount(buildID); got != 1 {
		t.Errorf("CallerCount(build.yml) = %d, want 1", got)
	}
	// Local reusable edge resolves to the build.yml node.
	calls := g.Calls(relID)
	var foundLocal, foundExternal bool
	for _, e := range calls {
		if e.ToID == buildID && e.Kind == KindReusable {
			foundLocal = true
		}
		if e.Pin == "main" {
			foundExternal = true
			if n := g.Nodes[e.ToID]; n == nil || !n.External {
				t.Error("cross-repo target should be an external node")
			}
		}
	}
	if !foundLocal {
		t.Error("expected resolved local reusable edge release.yml -> build.yml")
	}
	if !foundExternal {
		t.Error("expected external edge with @main pin preserved")
	}
	if r := g.Reachable(relID); len(r) == 0 {
		t.Error("release.yml should reach at least build.yml")
	}
}

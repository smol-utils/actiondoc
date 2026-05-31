package callgraph

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/smol-utils/actiondoc/internal/model"
	"github.com/smol-utils/actiondoc/internal/parser"
)

func TestBuildBasic(t *testing.T) {
	entry := &model.Workflow{
		File: "release.yml", Name: "Release", On: []string{"workflow_dispatch"},
		Jobs: []model.Job{
			{ID: "publish", Uses: "./.github/workflows/build.yml",
				With:    []model.KV{{Key: "version", Value: "${{ inputs.version }}"}},
				Secrets: []model.KV{{Key: "CONSUMER-KEY", Value: "${{ secrets.SDKMAN_KEY }}"}}},
			{ID: "external", Uses: "pytorch/test-infra/.github/workflows/linux_job.yml@main"},
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

// TestBuildScala3Dogfood validates the graph against scala3 (item 7 driver: releases
// orchestrates several reusable-workflow caller jobs). Skipped if the dogfood corpus
// is not present.
func TestBuildScala3Dogfood(t *testing.T) {
	home, _ := os.UserHomeDir()
	wfDir := filepath.Join(home, "scratch", "actiondoc-dogfood", "scala3", ".github", "workflows")
	if _, err := os.Stat(wfDir); err != nil {
		t.Skipf("dogfood not present: %v", err)
	}
	entries, _ := os.ReadDir(wfDir)
	var sources []Source
	for _, e := range entries {
		ext := filepath.Ext(e.Name())
		if ext != ".yml" && ext != ".yaml" {
			continue
		}
		p := filepath.Join(wfDir, e.Name())
		w, err := parser.ParseFile(p)
		if err != nil {
			t.Logf("parse %s: %v", e.Name(), err)
			continue
		}
		sources = append(sources, Source{Path: p, Workflow: w})
	}
	g := Build(sources)

	reusableEdges := 0
	for _, e := range g.Edges {
		if e.Kind == KindReusable {
			reusableEdges++
		}
	}
	t.Logf("scala3: %d workflows, %d edges (%d reusable)", len(sources), len(g.Edges), reusableEdges)
	if reusableEdges == 0 {
		t.Error("expected scala3 to have reusable-workflow caller edges (releases.yml orchestration)")
	}
}

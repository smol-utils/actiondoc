package renderer

import (
	"flag"
	"os"
	"strings"
	"testing"

	"github.com/smol-utils/actiondoc/internal/callgraph"
	"github.com/smol-utils/actiondoc/internal/parser"
)

var updateGolden = flag.Bool("update", false, "rewrite golden files with current output")

// TestGoldenReusableWorkflows is an end-to-end test of the reusable-workflow and
// call-graph rendering: parse a multi-file workflow set (entry point -> middle reusable
// -> leaf reusable, plus a local composite action and a cross-repo external reference),
// build the call graph, render every file, and compare against the golden output.
//
// To regenerate the golden file after an intentional format change:
//
//	go test ./internal/renderer -run TestGoldenReusableWorkflows -update
func TestGoldenReusableWorkflows(t *testing.T) {
	workflowPaths := []string{
		testdataPath("s1/release.yml"),
		testdataPath("s1/build_and_publish.yml"),
		testdataPath("s1/build.yml"),
	}
	actionPath := testdataPath("s1/actions/setup/action.yml")
	goldenPath := testdataPath("s1/expected-output.md")

	var sources []callgraph.Source
	for _, p := range workflowPaths {
		w, err := parser.ParseFile(p)
		if err != nil {
			t.Fatalf("ParseFile(%s): %v", p, err)
		}
		sources = append(sources, callgraph.Source{Path: p, Workflow: w})
	}
	a, err := parser.ParseActionFile(actionPath)
	if err != nil {
		t.Fatalf("ParseActionFile(%s): %v", actionPath, err)
	}
	sources = append(sources, callgraph.Source{Path: actionPath, Action: a})

	g := callgraph.Build(sources)

	var out strings.Builder
	for _, s := range sources {
		if s.Workflow != nil {
			out.WriteString(RenderMarkdownGraph(s.Workflow, g, s.Path))
		} else {
			out.WriteString(RenderActionMarkdown(s.Action))
		}
	}
	got := out.String()

	if *updateGolden {
		if err := os.WriteFile(goldenPath, []byte(got), 0644); err != nil {
			t.Fatalf("writing golden file: %v", err)
		}
		return
	}

	want, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("reading golden file: %v", err)
	}
	if got != string(want) {
		t.Errorf("output does not match golden file.\n\nTo update, run:\n  go test ./internal/renderer -run TestGoldenReusableWorkflows -update\n\nGot:\n%s", got)
	}
}

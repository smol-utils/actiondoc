package renderer

import (
	"flag"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/smol-utils/actiondoc/internal/callgraph"
	"github.com/smol-utils/actiondoc/internal/parser"
)

var updateGolden = flag.Bool("update", false, "rewrite golden files with current output")

func testdataPath(name string) string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "testdata", name)
}

// goldenCase is one parse -> render -> compare fixture. Paths are relative to testdata/.
// When graph is set, all sources are parsed into a call graph and rendered with graph
// context (workflows first, then actions); otherwise the single source renders standalone.
type goldenCase struct {
	name      string
	workflows []string
	actions   []string
	golden    string
	graph     bool
}

var goldenCases = []goldenCase{
	{
		name:      "workflow",
		workflows: []string{"sample-workflow.yml"},
		golden:    "expected-output.md",
	},
	{
		name:    "action",
		actions: []string{"action.yml"},
		golden:  "expected-action-output.md",
	},
	{
		// Reusable-workflow and call-graph rendering: entry point -> middle reusable ->
		// leaf reusable, plus a local composite action and a cross-repo external reference.
		name:      "reusable-callgraph",
		graph:     true,
		workflows: []string{"s1/release.yml", "s1/build_and_publish.yml", "s1/build.yml"},
		actions:   []string{"s1/actions/setup/action.yml"},
		golden:    "s1/expected-output.md",
	},
	{
		// Declared surface: triggers, permissions, env, concurrency, defaults,
		// environment binding, implicit description.
		name:      "declared-surface",
		workflows: []string{"s2/surface.yml"},
		golden:    "s2/surface.expected.md",
	},
	{
		name:      "reusable-surface",
		workflows: []string{"s2/reusable.yml"},
		golden:    "s2/reusable.expected.md",
	},
	{
		name:      "license-header",
		workflows: []string{"s2/license-header.yml"},
		golden:    "s2/license-header.expected.md",
	},
	{
		// Step rendering, matrix job names, runs-on normalization, and secret/variable
		// aggregation.
		name:      "steps",
		workflows: []string{"s3/steps.yml"},
		golden:    "s3/steps.expected.md",
	},
}

// TestGolden runs every renderer golden fixture through parse -> render -> compare. To
// regenerate all golden files after an intentional format change:
//
//	make golden
func TestGolden(t *testing.T) {
	for _, tc := range goldenCases {
		t.Run(tc.name, func(t *testing.T) {
			got := renderCase(t, tc)
			goldenPath := testdataPath(tc.golden)

			if *updateGolden {
				if err := os.WriteFile(goldenPath, []byte(got), 0o644); err != nil {
					t.Fatalf("writing golden file: %v", err)
				}
				return
			}

			want, err := os.ReadFile(goldenPath)
			if err != nil {
				t.Fatalf("reading golden file: %v", err)
			}
			if got != string(want) {
				t.Errorf("output does not match %s.\n\nTo update all golden files, run:\n  make golden\n\nGot:\n%s", tc.golden, got)
			}
		})
	}
}

// renderCase parses a case's fixtures and renders them the way the fixture was designed
// to be rendered: standalone for single-source cases, via the call graph for multi-source
// cases (matching how directory mode assembles output).
func renderCase(t *testing.T, tc goldenCase) string {
	t.Helper()

	if !tc.graph {
		switch {
		case len(tc.workflows) == 1 && len(tc.actions) == 0:
			w, err := parser.ParseFile(testdataPath(tc.workflows[0]))
			if err != nil {
				t.Fatalf("ParseFile: %v", err)
			}
			return RenderMarkdown(w)
		case len(tc.actions) == 1 && len(tc.workflows) == 0:
			a, err := parser.ParseActionFile(testdataPath(tc.actions[0]))
			if err != nil {
				t.Fatalf("ParseActionFile: %v", err)
			}
			return RenderActionMarkdown(a)
		default:
			t.Fatalf("non-graph golden case needs exactly one source; got %d workflows, %d actions",
				len(tc.workflows), len(tc.actions))
			return ""
		}
	}

	var sources []callgraph.Source
	for _, p := range tc.workflows {
		w, err := parser.ParseFile(testdataPath(p))
		if err != nil {
			t.Fatalf("ParseFile(%s): %v", p, err)
		}
		sources = append(sources, callgraph.Source{Path: testdataPath(p), Workflow: w})
	}
	for _, p := range tc.actions {
		a, err := parser.ParseActionFile(testdataPath(p))
		if err != nil {
			t.Fatalf("ParseActionFile(%s): %v", p, err)
		}
		sources = append(sources, callgraph.Source{Path: testdataPath(p), Action: a})
	}
	g := callgraph.Build(sources)

	var out strings.Builder
	for _, s := range sources {
		if s.Workflow != nil {
			out.WriteString(RenderMarkdownGraph(s.Workflow, g, s.Path))
		} else {
			out.WriteString(RenderActionMarkdown(s.Action))
		}
	}
	return out.String()
}

package renderer

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/smol-utils/actiondoc/internal/model"
	"github.com/smol-utils/actiondoc/internal/parser"
)

func testdataPath(name string) string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "testdata", name)
}

// TestGoldenOutput is an end-to-end test: parse sample YAML -> render Markdown -> compare
// against the golden file. If the output format changes intentionally, regenerate
// the golden file with: go run . generate testdata/sample-workflow.yml > testdata/expected-output.md
func TestGoldenOutput(t *testing.T) {
	workflowPath := testdataPath("sample-workflow.yml")
	goldenPath := testdataPath("expected-output.md")

	w, err := parser.ParseFile(workflowPath)
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}

	got := RenderMarkdown(w)

	want, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("reading golden file: %v", err)
	}

	if got != string(want) {
		t.Errorf("output does not match golden file.\n\nTo update the golden file, run:\n  go run . generate testdata/sample-workflow.yml > testdata/expected-output.md\n\nGot:\n%s", got)
	}
}

func TestTableEscaping(t *testing.T) {
	w := &model.Workflow{
		File: "test.yml",
		Name: "Test",
		On:   []string{"push"},
		Tags: model.Tags{
			Secrets: []model.Param{
				{Name: "TOKEN", Description: "Use a|b for options"},
			},
		},
	}

	md := RenderMarkdown(w)

	// The pipe should be escaped so it doesn't break the table.
	if !strings.Contains(md, `a\|b`) {
		t.Errorf("pipe character not escaped in table cell.\n\nOutput:\n%s", md)
	}

	// The unescaped pipe should NOT appear in the table rows (header divider is ok).
	for _, line := range strings.Split(md, "\n") {
		if !strings.HasPrefix(line, "|") {
			continue
		}
		// Skip the header divider row.
		if strings.Contains(line, "---") {
			continue
		}
		// Count unescaped pipes (not preceded by backslash).
		content := strings.ReplaceAll(line, `\|`, "")
		pipes := strings.Count(content, "|")
		// A valid 3-column row has exactly 4 pipes: | col | col | col |
		if strings.Contains(line, "TOKEN") && pipes != 4 {
			t.Errorf("table row has %d unescaped pipes (expected 4): %s", pipes, line)
		}
	}
}

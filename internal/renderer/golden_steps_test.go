package renderer

import (
	"os"
	"testing"

	"github.com/smol-utils/actiondoc/internal/parser"
)

// TestGoldenSteps is the end-to-end test for step rendering, matrix job names, runs-on
// normalization, and secret/variable aggregation: parse the fixture -> render Markdown ->
// compare against the golden file. To regenerate after an intentional format change:
//
//	go run . generate testdata/s3/steps.yml > testdata/s3/steps.expected.md
func TestGoldenSteps(t *testing.T) {
	w, err := parser.ParseFile(testdataPath("s3/steps.yml"))
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}

	got := RenderMarkdown(w)

	want, err := os.ReadFile(testdataPath("s3/steps.expected.md"))
	if err != nil {
		t.Fatalf("reading golden file: %v", err)
	}

	if got != string(want) {
		t.Errorf("output does not match golden file.\n\nTo update:\n  go run . generate testdata/s3/steps.yml > testdata/s3/steps.expected.md\n\nGot:\n%s", got)
	}
}

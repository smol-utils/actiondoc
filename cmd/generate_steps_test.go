package cmd

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func testdataPath(name string) string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "testdata", name)
}

// TestGenerateDirectoryGolden is the end-to-end test for directory mode. The fixture
// deliberately exercises every render family in one pass: a reusable-workflow caller job
// (call graph + forwarded inputs/secrets), a leaf reusable workflow (called-by chain +
// workflow_call API), a discovered composite action (step `with:` docs), secret/var
// references, and a continue-on-error step. The substring guards below ensure no render
// family is silently dropped from the assembled Generate() output even if the byte golden
// is regenerated. To regenerate:
//
//	go run . generate testdata/s3/repo/.github/workflows > testdata/s3/repo.expected.md
func TestGenerateDirectoryGolden(t *testing.T) {
	out := filepath.Join(t.TempDir(), "out.md")

	err := Generate([]string{"-o", out, testdataPath("s3/repo/.github/workflows")})
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}

	got, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("reading output: %v", err)
	}

	// Guard against a render entry point dropping a whole section family (the exact
	// failure mode where graph sections vanished because the CLI called the non-graph
	// renderer). Each string belongs to a different session's feature.
	g := string(got)
	for _, must := range []string{
		"# Contents",                     // table of contents
		"## Call graph",                  // downstream call-graph tree
		"## Called by",                   // upstream caller chain on the reusable workflow
		"## Referenced secrets and variables", // auto-collected references
		"Uses workflow",                  // reusable-workflow caller job row
		"id-token",                       // caller job's own permissions must still render
		"- With:",                        // step with: block
		"`[continue-on-error]`",          // continue-on-error badge
	} {
		if !strings.Contains(g, must) {
			t.Errorf("Generate() output is missing %q -- a render entry point dropped a section family", must)
		}
	}

	want, err := os.ReadFile(testdataPath("s3/repo.expected.md"))
	if err != nil {
		t.Fatalf("reading golden file: %v", err)
	}

	if g != string(want) {
		t.Errorf("output does not match golden file.\n\nTo update:\n  go run . generate testdata/s3/repo/.github/workflows > testdata/s3/repo.expected.md\n\nGot:\n%s", got)
	}
}

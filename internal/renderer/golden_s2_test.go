package renderer

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/smol-utils/actiondoc/internal/parser"
)

// s2TestdataPath resolves a fixture under testdata/s2, the session-owned golden corpus
// kept separate from the shared master golden files.
func s2TestdataPath(name string) string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "testdata", "s2", name)
}

// TestGoldenS2 exercises the declared-surface rendering (triggers, permissions, env,
// concurrency, defaults, environment binding, implicit description) against per-fixture
// golden files. To regenerate after an intentional format change:
//
//	go run . generate testdata/s2/<name>.yml > testdata/s2/<name>.expected.md
func TestGoldenS2(t *testing.T) {
	fixtures := []string{"surface", "reusable", "license-header"}
	for _, name := range fixtures {
		t.Run(name, func(t *testing.T) {
			w, err := parser.ParseFile(s2TestdataPath(name + ".yml"))
			if err != nil {
				t.Fatalf("ParseFile: %v", err)
			}
			got := RenderMarkdown(w)

			want, err := os.ReadFile(s2TestdataPath(name + ".expected.md"))
			if err != nil {
				t.Fatalf("reading golden file: %v", err)
			}
			if got != string(want) {
				t.Errorf("output does not match golden file.\n\nTo update:\n  go run . generate testdata/s2/%s.yml > testdata/s2/%s.expected.md\n\nGot:\n%s", name, name, got)
			}
		})
	}
}

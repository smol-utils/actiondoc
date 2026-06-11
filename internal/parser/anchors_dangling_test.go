package parser

import (
	"os"
	"path/filepath"
	"testing"
)

// TestParseFileDanglingAlias locks the documented resolveAnchors fallback: an alias
// whose anchor was never defined is left as-is rather than guessed at, and parsing
// neither fails nor loses the rest of the document.
func TestParseFileDanglingAlias(t *testing.T) {
	path := filepath.Join(t.TempDir(), "dangling.yml")
	content := `name: Aliases
on: push
env: *missing
jobs:
  build:
    runs-on: *ghost
    steps:
      - run: echo hi
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	w, err := ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}
	if len(w.Jobs) != 1 || len(w.Jobs[0].Steps) != 1 {
		t.Errorf("rest of document must survive a dangling alias: %+v", w.Jobs)
	}
	// A dangling scalar-position alias is kept verbatim; a dangling mapping-position
	// alias contributes nothing.
	if got := w.Jobs[0].RunsOn; got != "*ghost" {
		t.Errorf("runs-on = %q, want the unresolved alias kept verbatim", got)
	}
	if len(w.Env) != 0 {
		t.Errorf("env from a dangling alias should be empty, got %+v", w.Env)
	}
}

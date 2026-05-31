package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

// TestResolveFilesCompositeDiscovery verifies item 6: pointing at a .github/workflows
// dir discovers both the workflow files and sibling composite actions under
// .github/actions/ (including a nested one and a composite placed in the workflows dir).
func TestResolveFilesCompositeDiscovery(t *testing.T) {
	root := t.TempDir()
	mkdir := func(p string) {
		if err := os.MkdirAll(p, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	write := func(p string) {
		mkdir(filepath.Dir(p))
		if err := os.WriteFile(p, []byte("name: x\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	gh := filepath.Join(root, ".github")
	write(filepath.Join(gh, "workflows", "ci.yml"))
	write(filepath.Join(gh, "workflows", "release.yaml"))
	write(filepath.Join(gh, "workflows", "action.yml"))             // composite in workflows dir (vets-website)
	write(filepath.Join(gh, "actions", "build", "action.yml"))      // sibling composite
	write(filepath.Join(gh, "actions", "deep", "x", "action.yaml")) // nested composite

	got, err := resolveFiles(filepath.Join(gh, "workflows"))
	if err != nil {
		t.Fatalf("resolveFiles: %v", err)
	}

	// Compare on paths relative to the temp root for stable assertions.
	found := map[string]bool{}
	for _, f := range got {
		rel, _ := filepath.Rel(root, f)
		found[filepath.ToSlash(rel)] = true
	}
	want := []string{
		".github/workflows/ci.yml",
		".github/workflows/release.yaml",
		".github/workflows/action.yml",
		".github/actions/build/action.yml",
		".github/actions/deep/x/action.yaml",
	}
	for _, w := range want {
		if !found[w] {
			t.Errorf("expected to discover %s; got %v", w, got)
		}
	}
	if len(got) != len(want) {
		t.Errorf("got %d files, want %d: %v", len(got), len(want), got)
	}
}

package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestResolveFilesCompositeDiscovery verifies that pointing at a .github/workflows dir
// discovers both the workflow files and sibling composite actions under
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
	write(filepath.Join(gh, "workflows", "action.yml"))             // composite placed in the workflows dir
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

	// Discovery order must be deterministic so directory-mode TOC/section order is stable.
	for i := 0; i < 3; i++ {
		again, err := resolveFiles(filepath.Join(gh, "workflows"))
		if err != nil {
			t.Fatalf("resolveFiles (repeat): %v", err)
		}
		if strings.Join(again, "\n") != strings.Join(got, "\n") {
			t.Errorf("resolveFiles order not deterministic:\n first: %v\n again: %v", got, again)
		}
	}
}

// TestGenerateNoWorkflowName verifies that a workflow without a `name:` still renders a
// non-empty title and TOC entry (the parser backfills the file name), so directory-mode
// output never produces an empty heading or anchor.
func TestGenerateNoWorkflowName(t *testing.T) {
	root := t.TempDir()
	wf := filepath.Join(root, ".github", "workflows")
	if err := os.MkdirAll(wf, 0o755); err != nil {
		t.Fatal(err)
	}
	// two unnamed workflows -> distinct file-name titles, no empty/duplicate anchors
	for _, name := range []string{"ci.yml", "release.yml"} {
		if err := os.WriteFile(filepath.Join(wf, name), []byte("on: push\njobs:\n  x:\n    runs-on: ubuntu-latest\n    steps:\n      - run: echo hi\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	out := filepath.Join(t.TempDir(), "out.md")
	if err := Generate([]string{"-o", out, wf}); err != nil {
		t.Fatalf("Generate: %v", err)
	}
	b, _ := os.ReadFile(out)
	md := string(b)
	for _, want := range []string{"# ci.yml", "# release.yml", "[ci.yml](#ciyml)", "[release.yml](#releaseyml)"} {
		if !strings.Contains(md, want) {
			t.Errorf("output missing %q (empty name not backfilled?)\n%s", want, md)
		}
	}
}

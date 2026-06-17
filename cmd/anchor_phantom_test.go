package cmd

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// generateToString writes the given workflow files into a temp .github/workflows directory,
// runs Generate over it, and returns the rendered document.
func generateToString(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	wf := filepath.Join(root, ".github", "workflows")
	if err := os.MkdirAll(wf, 0o755); err != nil {
		t.Fatal(err)
	}
	for name, body := range files {
		if err := os.WriteFile(filepath.Join(wf, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	out := filepath.Join(t.TempDir(), "out.md")
	if err := Generate([]string{"-o", out, wf}); err != nil {
		t.Fatalf("Generate: %v", err)
	}
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

// TestAnchorScanPhantomHeadingInDescription guards that a raw ATX heading emitted inside a
// description body (e.g. a "### Foo" line in an @desc, which GitHub renders and counts as a
// real heading) does not shift the mapping of link-target headings to their nodes. Such a
// heading targets no node, so every job mini-TOC link must still resolve to its own job
// heading rather than borrowing the phantom heading's slug.
func TestAnchorScanPhantomHeadingInDescription(t *testing.T) {
	doc := generateToString(t, map[string]string{
		"ci.yml": "# @desc ### Phantom Heading\n" +
			"name: CI\non: push\njobs:\n" +
			"  alpha:\n    name: Alpha\n    runs-on: ubuntu-latest\n    steps:\n      - run: echo a\n" +
			"  beta:\n    name: Beta\n    runs-on: ubuntu-latest\n    steps:\n      - run: echo b\n",
	})

	// The phantom heading must actually appear in the body (so the scan really sees the extra
	// level-3 heading), but it must not be used as any job's link target.
	if !strings.Contains(doc, "### Phantom Heading") {
		t.Fatalf("expected the raw description heading in the output:\n%s", doc)
	}
	for _, want := range []string{"[Alpha](#alpha-alpha)", "[Beta](#beta-beta)"} {
		if !strings.Contains(doc, want) {
			t.Errorf("expected mini-TOC link %q in output:\n%s", want, doc)
		}
	}
	if strings.Contains(doc, "(#phantom-heading)") {
		t.Errorf("a description's phantom heading became a link target:\n%s", doc)
	}

	// Every in-page link target must resolve to a real heading GitHub derives from the body.
	assertLinksResolve(t, doc)
}

// TestAnchorScanPhantomLevel1HeadingMultiDoc is the multi-document counterpart: a level-1 ATX
// heading inside one workflow's description must not shift the section-title-to-source mapping
// the table of contents and trigger index depend on.
func TestAnchorScanPhantomLevel1HeadingMultiDoc(t *testing.T) {
	doc := generateToString(t, map[string]string{
		"aaa.yml": "# @desc # Phantom Section\n" +
			"name: Aaa Workflow\non: push\njobs:\n" +
			"  build:\n    runs-on: ubuntu-latest\n    steps:\n      - run: echo a\n",
		"bbb.yml": "name: Bbb Workflow\non: push\njobs:\n" +
			"  test:\n    runs-on: ubuntu-latest\n    steps:\n      - run: echo b\n",
	})

	if !strings.Contains(doc, "# Phantom Section") {
		t.Fatalf("expected the raw description heading in the output:\n%s", doc)
	}
	// Both the table of contents and the by-trigger index must link the second workflow to its
	// own section, not to the phantom section the first workflow's description injected.
	for _, want := range []string{"(#aaa-workflow)", "(#bbb-workflow)"} {
		if !strings.Contains(doc, want) {
			t.Errorf("expected section link %q in output:\n%s", want, doc)
		}
	}
	if strings.Contains(doc, "(#phantom-section)") {
		t.Errorf("a description's phantom heading became a section link target:\n%s", doc)
	}
	assertLinksResolve(t, doc)
}

// assertLinksResolve checks every in-page "#anchor" link in the document points at a heading
// anchor GitHub actually derives from the rendered body, using the independent reimplementation
// of GitHub's slugger in githubAnchors.
func assertLinksResolve(t *testing.T, doc string) {
	t.Helper()
	anchors := githubAnchors(doc)
	linkRe := regexp.MustCompile(`\]\(#([a-z0-9_-]+)\)`)
	for _, m := range linkRe.FindAllStringSubmatch(doc, -1) {
		if m[1] == "contents" {
			continue
		}
		if !anchors[m[1]] {
			t.Errorf("link target #%s does not resolve to any rendered heading; anchors=%v", m[1], anchors)
		}
	}
}

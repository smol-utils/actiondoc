package cmd

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

// TestAnchorCollisionsAcrossNamespaces is the end-to-end guard that in-page links resolve to
// the right heading when a workflow title, a job name, and a structural section share a slug.
// GitHub numbers all same-slug headings of every level through one document-order counter, so:
//   - a workflow titled "Build" and a job id "build" in a later workflow must get "build" and
//     "build-1" (not two "build" anchors from separate counters), and
//   - a job id "permissions" must be numbered after the "## Permissions" structural section
//     GitHub also counts, yielding "permissions-1" rather than colliding with "#permissions".
//
// The assertion checks the link target against the anchor GitHub actually derives from each
// heading, computed independently here from the rendered document.
func TestAnchorCollisionsAcrossNamespaces(t *testing.T) {
	root := t.TempDir()
	wf := filepath.Join(root, ".github", "workflows")
	if err := os.MkdirAll(wf, 0o755); err != nil {
		t.Fatal(err)
	}
	write := func(name, body string) {
		if err := os.WriteFile(filepath.Join(wf, name), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	// "Build" sorts before "CI", so the workflow title "Build" is the first "build" slug and
	// CI's job id "build" is the second.
	write("build.yml", "name: Build\non: push\njobs:\n  compile:\n    runs-on: ubuntu-latest\n    steps:\n      - run: make\n  package:\n    runs-on: ubuntu-latest\n    steps:\n      - run: make package\n")
	write("ci.yml", "name: CI\non: push\npermissions:\n  contents: read\njobs:\n  build:\n    runs-on: ubuntu-latest\n    steps:\n      - run: echo build\n  permissions:\n    runs-on: ubuntu-latest\n    steps:\n      - run: echo perms\n")

	out := filepath.Join(t.TempDir(), "out.md")
	if err := Generate([]string{"-o", out, wf}); err != nil {
		t.Fatalf("Generate: %v", err)
	}
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	doc := string(data)

	// CI's job mini-TOC links must point at the job headings, not the colliding workflow title
	// or structural section.
	for _, want := range []string{
		"[`build`](#build-1)",             // CI's build job, not the "# Build" workflow at #build
		"[`permissions`](#permissions-1)", // CI's permissions job, not the "## Permissions" section
	} {
		if !strings.Contains(doc, want) {
			t.Errorf("expected mini-TOC link %q in output:\n%s", want, doc)
		}
	}

	// Every in-page link target must name a heading anchor that exists exactly once -- i.e. the
	// link points at a real, unambiguous heading. Build the set of anchors GitHub derives from
	// the rendered headings (document-order duplicate numbering) and confirm every "#..." link
	// resolves into it.
	anchors := githubAnchors(doc)
	linkRe := regexp.MustCompile(`\]\(#([a-z0-9_-]+)\)`)
	for _, m := range linkRe.FindAllStringSubmatch(doc, -1) {
		target := m[1]
		if target == "contents" {
			continue // the document-level "## Contents" heading (always the first occurrence)
		}
		if !anchors[target] {
			t.Errorf("link target #%s does not resolve to any rendered heading; anchors=%v", target, anchors)
		}
	}
}

// githubAnchors slugs every ATX heading of a rendered Markdown document the way GitHub does --
// one document-order counter across all heading levels, with "-N" suffixes for repeats -- and
// returns the resulting anchor set. It is an independent reimplementation of the production
// scan so the test does not merely assert the code agrees with itself.
func githubAnchors(doc string) map[string]bool {
	out := map[string]bool{}
	seen := map[string]int{}
	inFence := false
	headingRe := regexp.MustCompile(`^#{1,6} `)
	nonSlug := regexp.MustCompile(`[^a-z0-9 _-]`)
	for _, line := range strings.Split(doc, "\n") {
		if strings.HasPrefix(line, "```") {
			inFence = !inFence
			continue
		}
		if inFence || !headingRe.MatchString(line) {
			continue
		}
		text := strings.TrimLeft(line, "#")
		text = strings.TrimSpace(text)
		base := strings.ReplaceAll(nonSlug.ReplaceAllString(strings.ToLower(text), ""), " ", "-")
		slug := base
		if n := seen[base]; n > 0 {
			slug = base + "-" + strconv.Itoa(n)
		}
		seen[base]++
		out[slug] = true
	}
	return out
}

package renderer

import "testing"

// TestDocumentHeadingsCrossLevel is the core regression guard for anchor collisions across
// heading namespaces. GitHub slugs every heading of every level through one document-order
// counter, so a workflow title, a job heading, and a structural section that share a base
// slug must get distinct, correctly numbered anchors -- the property that broke when section
// titles and job headings were numbered by separate counters that ignored structural sections.
func TestDocumentHeadingsCrossLevel(t *testing.T) {
	doc := "# Build\n\n" + // workflow title -> build
		"## Jobs\n\n" +
		"### `compile`\n\n" + // compile
		"# CI\n\n" + // ci
		"## Permissions\n\n" + // structural section -> permissions (counted, not a link target)
		"## Jobs\n\n" + // jobs-1
		"### `build`\n\n" + // build-1 (second "build", even though the first was a workflow title)
		"### `permissions`\n\n" // permissions-1 (numbered after the structural "## Permissions")

	got := DocumentHeadings(doc)
	want := []DocHeading{
		{1, "build"},
		{2, "jobs"},
		{3, "compile"},
		{1, "ci"},
		{2, "permissions"},
		{2, "jobs-1"},
		{3, "build-1"},
		{3, "permissions-1"},
	}
	if len(got) != len(want) {
		t.Fatalf("DocumentHeadings returned %d headings, want %d: %+v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("heading %d = %+v, want %+v", i, got[i], want[i])
		}
	}
}

// TestDocumentHeadingsSkipsCodeFences verifies that "#" lines inside a fenced code block (an
// @example block can carry shell or YAML comments) are not mistaken for headings, matching
// GitHub, which does not slug headings inside code fences.
func TestDocumentHeadingsSkipsCodeFences(t *testing.T) {
	doc := "# Title\n\n" +
		"## Example\n\n" +
		"```\n" +
		"# this is a yaml comment, not a heading\n" +
		"### neither is this\n" +
		"```\n\n" +
		"## After\n"

	got := DocumentHeadings(doc)
	want := []DocHeading{
		{1, "title"},
		{2, "example"},
		{2, "after"},
	}
	if len(got) != len(want) {
		t.Fatalf("DocumentHeadings returned %d headings, want %d: %+v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("heading %d = %+v, want %+v", i, got[i], want[i])
		}
	}
}

// TestDocumentHeadingsNotHeadings guards the ATX-heading edge cases: a "#" with no following
// space is not a heading, and a run of seven or more "#" exceeds the maximum heading level.
func TestDocumentHeadingsNotHeadings(t *testing.T) {
	doc := "#nospace\n\n" +
		"####### too deep\n\n" +
		"# real\n"

	got := DocumentHeadings(doc)
	if len(got) != 1 || got[0] != (DocHeading{1, "real"}) {
		t.Fatalf("DocumentHeadings = %+v, want exactly [{1 real}]", got)
	}
}

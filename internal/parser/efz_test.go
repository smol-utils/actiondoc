package parser

import "testing"

// TestParseFileHeadCommentBeforeMarker covers the inverse of the license-header case: a
// non-license comment block (here an ActionDoc @desc) placed before the `---` marker must
// still be captured, not lost when the leading comment-only document is skipped.
func TestParseFileHeadCommentBeforeMarker(t *testing.T) {
	desc := "# @desc Runs the nightly cleanup against fresh mirrors.\n---\nname: Cleanup\non: schedule\njobs: {}\n"
	w, err := parseString(t, desc)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if w.Description != "Runs the nightly cleanup against fresh mirrors." {
		t.Errorf("Description = %q, want the @desc text from before the --- marker", w.Description)
	}

	// A plain (untagged) description comment before --- should also surface implicitly.
	plain := "# Builds and signs the release.\n---\nname: Release\non: push\njobs: {}\n"
	w2, err := parseString(t, plain)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if w2.Description != "Builds and signs the release." {
		t.Errorf("implicit Description = %q, want the pre-marker comment", w2.Description)
	}
}

// TestParseFileHeaderCommentDoc covers a file that opens with a '#' comment block
// followed by a '---' document-start marker, which makes goccy/go-yaml emit a leading
// comment-only document. The parser must skip it and parse the first mapping document
// rather than failing with "expected top-level mapping".
func TestParseFileHeaderCommentDoc(t *testing.T) {
	path := testdataPath("header-comment-doc.yml")

	w, err := ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile on header+--- file: %v", err)
	}
	if w.Name != "CI" {
		t.Errorf("Name = %q, want %q", w.Name, "CI")
	}
	if len(w.On) == 0 {
		t.Errorf("On is empty; triggers not parsed")
	}
	if len(w.Jobs) != 1 || w.Jobs[0].ID != "build" {
		t.Errorf("Jobs = %+v, want one job %q", w.Jobs, "build")
	}
	// The license header must NOT become the workflow description.
	if w.Description != "" {
		t.Errorf("Description = %q, want empty (license header must not leak in)", w.Description)
	}
}

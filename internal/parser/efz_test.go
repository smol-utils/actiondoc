package parser

import "testing"

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

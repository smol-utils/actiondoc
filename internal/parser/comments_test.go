package parser

import (
	"testing"

	"github.com/goccy/go-yaml/ast"
	yamlparser "github.com/goccy/go-yaml/parser"
)

// parseTopMapping is a test helper: parse a YAML snippet and return its root mapping.
func parseTopMapping(t *testing.T, src string) *ast.MappingNode {
	t.Helper()
	file, err := yamlparser.ParseBytes([]byte(src), yamlparser.ParseComments)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	_, root, _, ok := firstMappingDoc(file)
	if !ok {
		t.Fatalf("no mapping doc")
	}
	return root
}

// TestTrailingComment verifies the trailing-comment helper extracts an inline rationale
// comment from a mapping entry (e.g. a permission scope annotated with why it's needed).
func TestTrailingComment(t *testing.T) {
	root := parseTopMapping(t, "contents: read  # for actions/checkout to fetch code\n")
	if len(root.Values) == 0 {
		t.Fatal("no entries")
	}
	got := TrailingComment(root.Values[0])
	want := "for actions/checkout to fetch code"
	if got != want {
		t.Errorf("TrailingComment = %q, want %q", got, want)
	}
}

// TestLeadingComment verifies the leading-comment helper cleans a head comment block
// (used to derive an implicit description).
func TestLeadingComment(t *testing.T) {
	root := parseTopMapping(t, "# Runs the nightly cleanup pipeline.\nname: cleanup\n")
	if len(root.Values) == 0 {
		t.Fatal("no entries")
	}
	got := LeadingComment(root.Values[0], root.Values[0].Key)
	want := "Runs the nightly cleanup pipeline."
	if got != want {
		t.Errorf("LeadingComment = %q, want %q", got, want)
	}
}

// TestTrailingCommentIgnoresLeadingComment verifies that a comment on the line above an
// entry (a note about the whole block) is not attributed as that entry's rationale, while
// a same-line trailing comment still is.
func TestTrailingCommentIgnoresLeadingComment(t *testing.T) {
	root := parseTopMapping(t, `permissions:
  # All other permissions are set to none
  contents: read
  id-token: write  # for OIDC keyless signing
`)
	perms := root.Values[0].Value
	mapping, ok := perms.(*ast.MappingNode)
	if !ok {
		t.Fatalf("permissions value is %T, want mapping", perms)
	}
	if got := TrailingComment(mapping.Values[0]); got != "" {
		t.Errorf("leading block comment attributed as rationale: %q", got)
	}
	if got := TrailingComment(mapping.Values[1]); got != "for OIDC keyless signing" {
		t.Errorf("same-line trailing comment = %q, want it preserved", got)
	}
}

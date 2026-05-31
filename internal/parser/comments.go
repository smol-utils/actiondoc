package parser

import (
	"strings"

	"github.com/goccy/go-yaml/ast"
)

// cleanCommentText strips the leading "#"/"# " from each line of a raw comment group
// string and trims surrounding whitespace, joining multi-line comments with a space.
func cleanCommentText(raw string) string {
	if raw == "" {
		return ""
	}
	var parts []string
	for _, line := range strings.Split(raw, "\n") {
		parts = append(parts, stripCommentPrefix(line))
	}
	return strings.TrimSpace(strings.Join(parts, " "))
}

// TrailingComment returns the inline trailing comment attached to a mapping entry --
// the rationale in `contents: read  # for actions/checkout` -- with the leading "#"
// stripped. Empty if there is none. Foundation hook (M2) for the trailing-comment
// rationale features (roadmap items 4d, 15). The comment may be attached by goccy to
// either the value node or the entry node, so both are checked.
func TrailingComment(mv *ast.MappingValueNode) string {
	if mv == nil {
		return ""
	}
	for _, n := range []ast.Node{mv.Value, mv} {
		if n == nil {
			continue
		}
		if cg := n.GetComment(); cg != nil {
			if s := cleanCommentText(cg.String()); s != "" {
				return s
			}
		}
	}
	return ""
}

// LeadingComment returns the cleaned head-comment block attached to a node (the lines
// immediately above it), with "#" prefixes stripped. Foundation hook (M2) for implicit
// descriptions from leading comment blocks (roadmap item 19).
func LeadingComment(nodes ...ast.Node) string {
	return cleanCommentText(findComment(nodes...))
}

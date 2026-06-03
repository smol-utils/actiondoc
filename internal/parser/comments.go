package parser

import (
	"strings"

	"github.com/goccy/go-yaml/ast"
)

// cleanCommentText strips the leading "#"/"# " from each line of a raw comment group
// string and trims surrounding whitespace, joining multi-line comments with a space.
// Banner-style framing is decoration, not content, and is dropped: rule lines made
// entirely of '#' characters are skipped, and '#' runs decorating the start or end of a
// text line are trimmed, so a block like "###  TITLE  ###" reads back as just "TITLE".
func cleanCommentText(raw string) string {
	if raw == "" {
		return ""
	}
	var parts []string
	for _, line := range strings.Split(raw, "\n") {
		line = stripCommentPrefix(line)
		line = strings.TrimSpace(strings.Trim(line, "#"))
		if line == "" {
			continue
		}
		parts = append(parts, line)
	}
	return strings.TrimSpace(strings.Join(parts, " "))
}

// TrailingComment returns the inline trailing comment attached to a mapping entry --
// the rationale in `contents: read  # for actions/checkout` -- with the leading "#"
// stripped. Empty if there is none. The comment may be attached by goccy to either the
// value node or the entry node, so both are checked. Only a comment on the entry's own
// line counts: a comment on an earlier line is a leading note about what follows (e.g. a
// remark above a whole permissions block), not a description of this one entry.
func TrailingComment(mv *ast.MappingValueNode) string {
	if mv == nil {
		return ""
	}
	entryLine := 0
	if mv.Key != nil && mv.Key.GetToken() != nil {
		entryLine = mv.Key.GetToken().Position.Line
	}
	for _, n := range []ast.Node{mv.Value, mv} {
		if n == nil {
			continue
		}
		cg := n.GetComment()
		if cg == nil {
			continue
		}
		if tok := cg.GetToken(); tok != nil && entryLine != 0 && tok.Position.Line != entryLine {
			continue
		}
		if s := cleanCommentText(cg.String()); s != "" {
			return s
		}
	}
	return ""
}

// LeadingComment returns the cleaned head-comment block attached to a node (the lines
// immediately above it), with "#" prefixes stripped. Used to derive an implicit
// description from a leading comment block.
func LeadingComment(nodes ...ast.Node) string {
	return cleanCommentText(findComment(nodes...))
}

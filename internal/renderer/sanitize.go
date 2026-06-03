package renderer

import "strings"

// This file is the single seam for making arbitrary strings safe inside Markdown
// structures: table cells (escapeCell, cellOrDash, codeCellOrDash), inline code spans
// (codeSpan), bold/emphasis runs (escapeInline), and single-line contexts (oneLine).
// Renderers should reach for these instead of wrapping values in backticks, asterisks,
// or pipes directly, so a value containing Markdown-significant characters can never
// break the surrounding structure. The one deliberate exception: top-level section
// headings (workflow/action titles) render raw, because their anchors must match what
// GitHub's slugger derives from the unescaped text.

// escapeCell escapes characters that break Markdown table cells. Newlines become
// <br> (not a space) so multi-line values like multi-line `if:` conditions keep their
// visual line breaks instead of collapsing or, worse, being parsed as a new table row.
func escapeCell(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "|", "\\|")
	s = strings.ReplaceAll(s, "\n", "<br>")
	return s
}

// cellOrDash escapes a value for a table cell, substituting "-" when empty.
func cellOrDash(s string) string {
	if s == "" {
		return "-"
	}
	return escapeCell(s)
}

// codeCellOrDash escapes a value and wraps it in code formatting for a table cell,
// substituting "-" when empty. The code span is backtick-safe (see codeSpan).
func codeCellOrDash(s string) string {
	if s == "" {
		return "-"
	}
	return codeSpan(escapeCell(s))
}

// codeSpan wraps s in an inline Markdown code span. A plain single-backtick span breaks
// when the value itself contains backticks (e.g. a github-script step embedding JS
// template literals), so the delimiter is one backtick longer than the longest backtick
// run in s, with space padding when the content starts or ends with a backtick
// (CommonMark strips one such space). An empty value renders as "-", since a code span
// cannot be empty.
func codeSpan(s string) string {
	if s == "" {
		return "-"
	}
	longest, run := 0, 0
	for i := 0; i < len(s); i++ {
		if s[i] == '`' {
			run++
			if run > longest {
				longest = run
			}
		} else {
			run = 0
		}
	}
	if longest == 0 {
		return "`" + s + "`"
	}
	delim := strings.Repeat("`", longest+1)
	return delim + " " + s + " " + delim
}

// oneLine collapses newlines to spaces so a value renders safely inside a single-line
// context (inline code spans, list items).
func oneLine(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\n", " ")
	return strings.TrimSpace(s)
}

// escapeInline escapes the characters that would start or end Markdown inline markup
// (emphasis, code spans) inside running text, so an arbitrary value rendered inside
// **bold** -- a step title built from a name or a run: command line -- shows its
// backticks, asterisks, and underscores literally instead of changing the markup.
func escapeInline(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '*', '_', '`':
			b.WriteByte('\\')
		}
		b.WriteByte(s[i])
	}
	return b.String()
}

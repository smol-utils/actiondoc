package renderer

import "testing"

func TestEscapeCell(t *testing.T) {
	cases := []struct{ in, want string }{
		{"plain", "plain"},
		{"a | b", "a \\| b"},
		{"line1\nline2", "line1<br>line2"},
		{"crlf\r\nline", "crlf<br>line"},
		{"x &&\ny || z", "x &&<br>y \\|\\| z"},
		// A backslash before a pipe (a regex alternation `a\|b`, an authored `\|`) must
		// escape the backslash FIRST. Producing `\\|` would let cmark-gfm read `\\` as an
		// escaped backslash and the trailing `|` as a column delimiter, splitting the cell;
		// `\\\|` keeps both literal so the cell stays one column.
		{"a\\|b", "a\\\\\\|b"},
		{"x\\|y", "x\\\\\\|y"},
		{"path\\to", "path\\\\to"},
	}
	for _, c := range cases {
		if got := escapeCell(c.in); got != c.want {
			t.Errorf("escapeCell(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

// TestEscapeCellColumnCount is the behavioral guard for the backslash-before-pipe fix: the
// confirmed corruption was a cell splitting into extra table columns. A correctly escaped
// cell, joined into a `| a | b |` row, must still split back into exactly the original
// columns. cmark-gfm's table tokenizer treats `\\` as a literal backslash (consuming both)
// and a following bare `|` as a delimiter, so the count below mirrors how a renderer would
// see the row -- counting unescaped pipes after collapsing escaped backslash pairs.
func TestEscapeCellColumnCount(t *testing.T) {
	values := []string{"matches a\\|b alternation", "x\\|y"}
	for _, v := range values {
		cell := escapeCell(v)
		if delimiters := unescapedPipes(cell); delimiters != 0 {
			t.Errorf("escapeCell(%q) = %q leaves %d unescaped pipe delimiter(s); the cell would split into extra columns", v, cell, delimiters)
		}
	}
}

// unescapedPipes counts pipes that cmark-gfm's table row tokenizer would treat as column
// delimiters: a backslash escapes the next character (so `\\` consumes both and `\|` hides
// the pipe), and any pipe not so protected is a delimiter.
func unescapedPipes(s string) int {
	n := 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '\\':
			i++ // skip the escaped character
		case '|':
			n++
		}
	}
	return n
}

func TestCodeSpan(t *testing.T) {
	cases := []struct{ in, want string }{
		{"plain", "`plain`"},
		{"", "-"},
		{"a `quoted` word", "`` a `quoted` word ``"},
		{"console.log(`hi ${x}`)", "`` console.log(`hi ${x}`) ``"},
		{"``double run``", "``` ``double run`` ```"},
		{"`leading backtick", "`` `leading backtick ``"},
		{"trailing backtick`", "`` trailing backtick` ``"},
	}
	for _, c := range cases {
		if got := codeSpan(c.in); got != c.want {
			t.Errorf("codeSpan(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestCodeCellOrDash(t *testing.T) {
	cases := []struct{ in, want string }{
		{"", "-"},
		{"value", "`value`"},
		{"with `tick`", "`` with `tick` ``"},
		{"a|b", "`a\\|b`"},
	}
	for _, c := range cases {
		if got := codeCellOrDash(c.in); got != c.want {
			t.Errorf("codeCellOrDash(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestAnchor(t *testing.T) {
	cases := []struct{ in, want string }{
		{"build", "build"},
		{"Build and Test", "build-and-test"},
		{"publish-sdkman.yml", "publish-sdkmanyml"},
		{"Deploy (prod)", "deploy-prod"},
		// Non-ASCII letters are Unicode word characters; GitHub keeps them (lowercased)
		// rather than stripping to ASCII, so a cross-link to such a heading resolves.
		{"Café Déploy", "café-déploy"},
		{"Über Build", "über-build"},
		// Emoji and other non-letter/digit symbols are still dropped; a removed symbol
		// between two spaces leaves a double hyphen, matching GitHub's slugger.
		{"build 🚀 now", "build--now"},
	}
	for _, c := range cases {
		if got := anchor(c.in); got != c.want {
			t.Errorf("anchor(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestAssignAnchors(t *testing.T) {
	got := AssignAnchors([]string{"model jobs", "Build", "model jobs", "model jobs", "Build"})
	want := []string{"model-jobs", "build", "model-jobs-1", "model-jobs-2", "build-1"}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("AssignAnchors[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestEscapeInline(t *testing.T) {
	cases := []struct{ in, want string }{
		{"plain title", "plain title"},
		{"echo \"a `b` c\" | grep d", "echo \"a \\`b\\` c\" | grep d"},
		{"run **all** the tests", "run \\*\\*all\\*\\* the tests"},
		{"snake_case_name", "snake\\_case\\_name"},
		// A backslash is escaped first, so an authored `\*` shows literally instead of
		// rendering as a stray backslash followed by an active emphasis marker.
		{"a\\*b", "a\\\\\\*b"},
		{"back\\slash", "back\\\\slash"},
	}
	for _, c := range cases {
		if got := escapeInline(c.in); got != c.want {
			t.Errorf("escapeInline(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

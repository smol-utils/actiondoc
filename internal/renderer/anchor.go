package renderer

import "strings"

// anchor converts a heading string into a GitHub-style Markdown anchor slug: lowercase,
// spaces to hyphens, drop everything that isn't a letter, digit, hyphen, or underscore.
// Used for the table of contents and caller/callee cross-links.
func anchor(s string) string {
	s = strings.ToLower(s)
	var b strings.Builder
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9', r == '-', r == '_':
			b.WriteRune(r)
		case r == ' ':
			b.WriteByte('-')
		}
	}
	return b.String()
}

package renderer

import (
	"fmt"
	"strings"
)

// mdLinkLabel escapes the bracket characters that would otherwise terminate or corrupt a
// Markdown link label, so `[label](target)` stays valid for an arbitrary title (e.g. a
// workflow name containing `]`). Used by the TOC and caller/callee cross-links.
func mdLinkLabel(s string) string {
	s = strings.ReplaceAll(s, "[", "\\[")
	s = strings.ReplaceAll(s, "]", "\\]")
	return s
}

// AssignAnchors computes the document anchor slug for each title in render order,
// applying GitHub's "-N" suffix disambiguation for repeated titles (the second "Build"
// becomes "build-1"). It is the single owner of duplicate-name handling: the table of
// contents and every cross-link must use the same assignment, or links to a repeated
// title silently point at the wrong section.
func AssignAnchors(titles []string) []string {
	slugs := make([]string, len(titles))
	seen := map[string]int{}
	for i, t := range titles {
		base := anchor(t)
		slug := base
		if n := seen[base]; n > 0 {
			slug = fmt.Sprintf("%s-%d", base, n)
		}
		seen[base]++
		slugs[i] = slug
	}
	return slugs
}

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

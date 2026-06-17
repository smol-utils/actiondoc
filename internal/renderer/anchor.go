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

// DocHeading is one heading found in a rendered document: its level (1 for "#", 2 for
// "##", ...) and the GitHub anchor slug it resolves to, including the running "-N"
// duplicate suffix.
type DocHeading struct {
	Level int
	Slug  string
}

// DocumentHeadings scans a fully rendered Markdown document and returns every ATX heading,
// in document order, paired with the GitHub anchor slug it resolves to. It is the single
// source of truth for anchor assignment: GitHub's slugger numbers ALL same-slug headings of
// EVERY level together, in document order (the first "build" becomes "#build", the next
// "#build-1", ...), regardless of whether a heading is a section title, a job heading, or a
// structural section like "## Permissions". Assigning section and job anchors from separate
// counters -- or omitting the structural headings GitHub still counts -- shifts the suffixes
// and makes cross-links point at the wrong heading, so callers feed the whole rendered body
// through this one pass and read each link target's anchor back from it.
//
// Headings inside fenced code blocks (``` ... ```) are ignored, matching GitHub, which does
// not slug content inside code fences (an "@example" block can contain "# comment" lines).
func DocumentHeadings(doc string) []DocHeading {
	var out []DocHeading
	inFence := false
	var fenceChar byte
	fenceLen := 0
	for _, line := range strings.Split(doc, "\n") {
		// CommonMark allows up to three leading spaces before a fence or heading; four or
		// more makes it an indented code block, never a heading.
		body := strings.TrimLeft(line, " ")
		if len(line)-len(body) > 3 {
			continue
		}
		if c, n := fenceMarker(body); c != 0 {
			switch {
			case !inFence:
				inFence, fenceChar, fenceLen = true, c, n
			case c == fenceChar && n >= fenceLen:
				inFence, fenceChar, fenceLen = false, 0, 0
			}
			continue
		}
		if inFence {
			continue
		}
		if level, text, ok := atxHeading(body); ok {
			out = append(out, DocHeading{Level: level, Slug: anchor(text)})
		}
	}
	// One shared counter over every heading, in order: the Nth repeat of a base slug takes
	// the "-N" suffix GitHub assigns.
	seen := map[string]int{}
	for i := range out {
		base := out[i].Slug
		if n := seen[base]; n > 0 {
			out[i].Slug = fmt.Sprintf("%s-%d", base, n)
		}
		seen[base]++
	}
	return out
}

// fenceMarker reports the fence character ('`' or '~') and run length of a code-fence line
// (three or more of the same fence char), or (0, 0) when the line opens no fence.
func fenceMarker(s string) (byte, int) {
	if len(s) < 3 || (s[0] != '`' && s[0] != '~') {
		return 0, 0
	}
	c := s[0]
	n := 0
	for n < len(s) && s[n] == c {
		n++
	}
	if n < 3 {
		return 0, 0
	}
	return c, n
}

// atxHeading parses an ATX heading line ("# Title" through "###### Title"): one to six
// leading "#" characters followed by a space (or end of line). It returns the heading level
// and its trimmed text. A run of seven or more "#", or one not followed by a space, is not a
// heading. Any trailing "#" run (an optional ATX closing sequence) is dropped.
func atxHeading(s string) (level int, text string, ok bool) {
	i := 0
	for i < len(s) && s[i] == '#' {
		i++
	}
	if i == 0 || i > 6 {
		return 0, "", false
	}
	if i < len(s) && s[i] != ' ' {
		return 0, "", false
	}
	text = strings.TrimSpace(s[i:])
	text = strings.TrimRight(strings.TrimRight(text, "#"), " ")
	return i, text, true
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

package redact

import (
	"regexp"
	"strings"
)

// This file holds the low-level string scanners shared by the collect and rewrite
// passes. Two kinds of redaction live here:
//
//   - Exact, structural redaction of identifiers that appear inside ${{ }} expression
//     contexts (secrets.X, vars.Y, env.Z). These come from parsed fields, so the match
//     is exact: a context prefix on an identifier boundary.
//   - Fuzzy, regex redaction of hosts and URLs in free text (run: scripts, descriptions,
//     literal values). This is best-effort by nature; its limits are documented in the
//     spec.
//
// Both passes walk strings the same way, so collection (gather originals) and rewrite
// (substitute placeholders) never disagree about what counts as a redactable token.

// contexts are the expression namespaces whose identifiers carry sensitive names.
var contexts = []string{"secrets", "vars", "env"}

// urlRe matches a scheme-qualified URL (https://..., docker://..., git+ssh://...). The
// body stops at whitespace, quotes, backticks, and bracketing punctuation so a URL
// embedded in prose or a shell command does not swallow the surrounding text.
var urlRe = regexp.MustCompile("(?i)\\b[a-z][a-z0-9+.-]*://[^\\s\"'`<>()\\[\\]{}|]+")

// hostRe matches a bare dotted hostname (deploy.internal.corp, registry.example.com).
// It deliberately requires at least one dot and an alphabetic final label so that
// version strings (1.2.3) and most numeric tokens do not match. File-like names
// (deploy.yml) still match the shape, so the rewrite step filters them with a
// known-extension denylist; see looksLikeFilename.
var hostRe = regexp.MustCompile("(?i)\\b(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z][a-z0-9-]*[a-z]\\b")

// fileExts are trailing labels that are almost always file extensions rather than a
// domain's TLD. A dotted token ending in one of these is left untouched, so config.yaml
// or deploy.sh is not mistaken for a hostname.
var fileExts = map[string]bool{
	"yml": true, "yaml": true, "json": true, "md": true, "txt": true,
	"go": true, "js": true, "ts": true, "sh": true, "py": true, "rb": true,
	"lock": true, "mod": true, "sum": true, "toml": true, "cfg": true,
	"ini": true, "xml": true, "html": true, "css": true, "env": true,
}

// keepSecrets are built-in secret names that are universal GitHub knowledge, not
// sensitive-but-unclassified material. Leaving GITHUB_TOKEN readable keeps the output
// intelligible without leaking anything.
var keepSecrets = map[string]bool{"GITHUB_TOKEN": true}

// keepRunners are GitHub-hosted runner labels and generic selector keywords that reveal
// nothing about private infrastructure. Any runner label outside this set is treated as a
// self-hosted label and redacted.
var keepRunners = map[string]bool{
	"ubuntu-latest": true, "ubuntu-24.04": true, "ubuntu-22.04": true, "ubuntu-20.04": true,
	"windows-latest": true, "windows-2025": true, "windows-2022": true, "windows-2019": true,
	"macos-latest": true, "macos-15": true, "macos-14": true, "macos-13": true,
	"macos-12": true, "macos-11": true,
	"self-hosted": true, "linux": true, "windows": true, "macos": true,
	"x64": true, "x86": true, "arm": true, "arm64": true,
}

// looksLikeFilename reports whether a dotted token's final label is a common file
// extension, so the host scanner can skip it.
func looksLikeFilename(host string) bool {
	dot := strings.LastIndex(host, ".")
	if dot < 0 {
		return false
	}
	return fileExts[strings.ToLower(host[dot+1:])]
}

// isIdentChar reports whether c can appear in a secret/var/env identifier. Names may
// contain hyphens (workflow_call keys frequently do).
func isIdentChar(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '_' || c == '-'
}

// walkIdents calls fn(ctx, name) for every `ctx.IDENT` reference in s, for each ctx in
// contexts. The leading-boundary check avoids matching a longer identifier that merely
// ends in ctx (mysecrets.X) or property access on an unrelated value
// (steps.secrets.outputs.X). It is the single definition both passes use: collection
// records the names, rewrite replaces them.
func walkIdents(s string, fn func(ctx, name string)) {
	for _, ctx := range contexts {
		needle := ctx + "."
		body := s
		offset := 0
		for {
			i := strings.Index(body[offset:], needle)
			if i < 0 {
				break
			}
			i += offset
			// Require a non-identifier char (or start of string) before the context, and
			// that the char before is not a dot (property access).
			if i > 0 && (isIdentChar(body[i-1]) || body[i-1] == '.') {
				offset = i + len(needle)
				continue
			}
			start := i + len(needle)
			j := start
			for j < len(body) && isIdentChar(body[j]) {
				j++
			}
			if j > start {
				fn(ctx, body[start:j])
			}
			offset = j
		}
	}
}

// rewriteIdents replaces every `ctx.OLD` whose OLD is mapped in repl[ctx] with
// `ctx.NEW`, preserving everything else. It walks s once, honoring the same identifier
// boundaries as walkIdents so it can never replace a partial match.
func rewriteIdents(s string, repl map[string]map[string]string) string {
	if !strings.Contains(s, ".") {
		return s
	}
	var b strings.Builder
	b.Grow(len(s))
	i := 0
	for i < len(s) {
		matched := false
		// Only consider a context start at an identifier boundary.
		atBoundary := i == 0 || !(isIdentChar(s[i-1]) || s[i-1] == '.')
		if atBoundary {
			for _, ctx := range contexts {
				needle := ctx + "."
				if !strings.HasPrefix(s[i:], needle) {
					continue
				}
				start := i + len(needle)
				j := start
				for j < len(s) && isIdentChar(s[j]) {
					j++
				}
				name := s[start:j]
				if repl[ctx] != nil {
					if to, ok := repl[ctx][name]; ok {
						b.WriteString(needle)
						b.WriteString(to)
						i = j
						matched = true
						break
					}
				}
			}
		}
		if matched {
			continue
		}
		b.WriteByte(s[i])
		i++
	}
	return b.String()
}

// eachLiteralSpan calls fn for each maximal run of s that lies outside a ${{ ... }}
// expression. Host/URL detection runs only on these spans: expression-context
// identifiers (secrets.X, github.sha, matrix.os) are dotted and would otherwise be
// mistaken for hostnames. An unterminated ${{ is treated as expression text and not
// scanned.
func eachLiteralSpan(s string, fn func(string)) {
	for s != "" {
		i := strings.Index(s, "${{")
		if i < 0 {
			fn(s)
			return
		}
		if i > 0 {
			fn(s[:i])
		}
		rest := s[i+3:]
		j := strings.Index(rest, "}}")
		if j < 0 {
			return
		}
		s = rest[j+2:]
	}
}

// mapLiteralSpans rebuilds s by applying fn to each literal (non-expression) span and
// copying every ${{ ... }} span through verbatim. It is the rewrite-side counterpart to
// eachLiteralSpan, so host/URL substitution never touches expression bodies.
func mapLiteralSpans(s string, fn func(string) string) string {
	var b strings.Builder
	b.Grow(len(s))
	for s != "" {
		i := strings.Index(s, "${{")
		if i < 0 {
			b.WriteString(fn(s))
			break
		}
		b.WriteString(fn(s[:i]))
		rest := s[i+3:]
		j := strings.Index(rest, "}}")
		if j < 0 {
			b.WriteString(s[i:]) // unterminated expression: copy through verbatim
			break
		}
		b.WriteString(s[i : i+3+j+2])
		s = rest[j+2:]
	}
	return b.String()
}

// walkHostsURLs calls onURL for each scheme-qualified URL and onHost for each bare
// hostname in s. URLs are scanned first and removed from consideration so a URL's own
// host is not also counted as a bare host. Filenames are skipped.
func walkHostsURLs(s string, onURL, onHost func(string)) {
	rest := urlRe.ReplaceAllStringFunc(s, func(m string) string {
		clean := strings.TrimRight(m, ".,;:!?")
		onURL(clean)
		return " " // collapse so the embedded host is not re-scanned below
	})
	for _, h := range hostRe.FindAllString(rest, -1) {
		if looksLikeFilename(h) {
			continue
		}
		onHost(h)
	}
}

// rewriteHostsURLs substitutes URL and host placeholders into s. URLs are replaced first
// (longest, scheme-qualified matches), then bare hosts in what remains, mirroring the
// scan order in walkHostsURLs so the two passes agree.
func rewriteHostsURLs(s string, urls, hosts map[string]string) string {
	s = urlRe.ReplaceAllStringFunc(s, func(m string) string {
		trailing := ""
		clean := strings.TrimRight(m, ".,;:!?")
		if len(clean) < len(m) {
			trailing = m[len(clean):]
		}
		if to, ok := urls[clean]; ok {
			return to + trailing
		}
		return m
	})
	s = hostRe.ReplaceAllStringFunc(s, func(m string) string {
		if looksLikeFilename(m) {
			return m
		}
		if to, ok := hosts[m]; ok {
			return to
		}
		return m
	})
	return s
}

package renderer

import "strings"

// expressions returns the inner text of each ${{ ... }} occurrence in s, in order.
// Foundation hook (M4) for rendering matrix/run-name expressions (roadmap items 21, 10).
func expressions(s string) []string {
	var out []string
	for {
		i := strings.Index(s, "${{")
		if i < 0 {
			break
		}
		rest := s[i+3:]
		j := strings.Index(rest, "}}")
		if j < 0 {
			break
		}
		out = append(out, strings.TrimSpace(rest[:j]))
		s = rest[j+2:]
	}
	return out
}

// expressionKind classifies a GitHub Actions expression by its leading context token,
// e.g. "matrix.java.version" -> "matrix", "github.event.x" -> "github". Returns "other"
// if the context is unrecognized. Consumers decide render policy: matrix axes get
// resolved to value lists (item 21); everything else renders verbatim (item 10).
func expressionKind(expr string) string {
	expr = strings.TrimSpace(expr)
	for _, ctx := range []string{
		"matrix", "github", "secrets", "vars", "inputs",
		"needs", "steps", "env", "job", "jobs", "runner", "strategy",
	} {
		if expr == ctx || strings.HasPrefix(expr, ctx+".") {
			return ctx
		}
	}
	return "other"
}

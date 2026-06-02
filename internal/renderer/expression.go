package renderer

import "strings"

// expressions returns the inner text of each ${{ ... }} occurrence in s, in order.
// Used when rendering strings that embed expressions (matrix references, run-name).
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
// if the context is unrecognized. Consumers decide render policy: matrix axes can be
// resolved to value lists; everything else renders verbatim.
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

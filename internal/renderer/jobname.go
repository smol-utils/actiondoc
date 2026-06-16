package renderer

import (
	"regexp"
	"strings"

	"github.com/smol-utils/actiondoc/internal/model"
)

// templateExpr matches a single GitHub Actions expression placeholder `${{ ... }}`. It is
// non-greedy so two adjacent placeholders in one string are matched separately rather than
// swallowed into one.
var templateExpr = regexp.MustCompile(`\$\{\{(.*?)\}\}`)

// spaceRun collapses internal whitespace left behind after a placeholder is rewritten.
var spaceRun = regexp.MustCompile(`\s+`)

// jobDisplayName is the single origin of a job's rendered display name. A job NAME can embed
// GitHub Actions expressions (e.g. `CLI ${{ matrix.job.os }}`), and that raw expression
// would otherwise leak into the job heading, the job mini-TOC, and the call graph -- and,
// because the heading text drives the GitHub anchor, into ugly anchors that cross-links no
// longer match. De-templating the name here, where every consumer reads it, keeps the
// heading, its anchor, the mini-TOC label, and the call-graph label in agreement.
//
// Each `${{ expr }}` is replaced by the expression's trailing identifier in parentheses
// (`matrix.job.os` -> `(os)`, `inputs.build-type` -> `(build-type)`), keeping the literal
// text around it (`Build PROD ${{ inputs.build-type }} image ${{ matrix.python-version }}`
// -> `Build PROD (build-type) image (python-version)`). When the name carries no literal
// word of its own -- it is pure placeholder noise like `${{ inputs.workflow-name }}` or a
// bracketed conditional -- there is nothing readable to keep, so the function falls back to
// the job id. The transform is deterministic.
func jobDisplayName(job *model.Job) string {
	name := job.Name
	if !strings.Contains(name, "${{") {
		return name
	}
	// The text outside the expressions is the only human-written part of the name. If it
	// carries no actual word (letter or digit), the whole name is placeholder noise and the
	// job id is the only readable handle left.
	if !hasAlphanumeric(templateExpr.ReplaceAllString(name, "")) {
		return job.ID
	}
	out := templateExpr.ReplaceAllStringFunc(name, func(m string) string {
		inner := m[3 : len(m)-2] // strip the leading "${{" and trailing "}}"
		return "(" + exprTrailingIdent(inner) + ")"
	})
	return strings.TrimSpace(spaceRun.ReplaceAllString(out, " "))
}

// exprTrailingIdent reduces a GitHub Actions expression to the identifier a reader cares
// about: the segment after its last dot (`matrix.job.os` -> `os`, `inputs.build-type` ->
// `build-type`). Expressions can trail operators or string literals (`matrix.os || 'x'`),
// so only the leading run of identifier characters after the dot is kept. An expression
// with no usable identifier yields a generic "expr" placeholder.
func exprTrailingIdent(expr string) string {
	s := strings.TrimSpace(expr)
	if i := strings.LastIndex(s, "."); i >= 0 {
		s = s[i+1:]
	}
	s = strings.TrimSpace(s)
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '_' || c == '-' || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') {
			b.WriteByte(c)
		} else {
			break
		}
	}
	if b.Len() == 0 {
		return "expr"
	}
	return b.String()
}

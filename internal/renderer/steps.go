package renderer

import (
	"fmt"
	"strings"

	"github.com/smol-utils/actiondoc/internal/model"
)

// renderStep writes one step as a numbered list item with its metadata (id, uses, condition,
// with: inputs) and a continue-on-error badge when set.
func renderStep(b *strings.Builder, step *model.Step, num int) {
	fmt.Fprintf(b, "%d. **%s**", num, stepTitle(step, num))
	if step.ContinueOnError {
		b.WriteString(" `[continue-on-error]`")
	} else if step.ContinueOnErrorExpr != "" {
		fmt.Fprintf(b, " `[continue-on-error: %s]`", step.ContinueOnErrorExpr)
	}
	if step.Description != "" {
		fmt.Fprintf(b, " - %s", step.Description)
	}
	b.WriteString("\n")

	if step.ID != "" {
		fmt.Fprintf(b, "   - ID: `%s`\n", step.ID)
	}
	if step.Uses != "" {
		fmt.Fprintf(b, "   - Uses: `%s`%s\n", step.Uses, usesSuffix(step.Uses, step.UsesVersion))
	}
	if step.If != "" {
		fmt.Fprintf(b, "   - Condition: `%s`\n", oneLine(step.If))
	}
	if len(step.With) > 0 {
		b.WriteString("   - With:\n")
		for _, kv := range step.With {
			fmt.Fprintf(b, "     - `%s`: `%s`%s\n", kv.Key, oneLine(kv.Value), withDoc(step, kv.Key))
		}
	}

	// Step-level tags
	writeStepParams(b, "Output", step.Tags.Outputs)
	writeStepParams(b, "Secret", step.Tags.Secrets)
	writeStepParams(b, "Env", step.Tags.Envs)

	b.WriteString("\n")
}

// writeStepParams writes inline bullet points for step-level params.
func writeStepParams(b *strings.Builder, label string, params []model.Param) {
	for _, p := range params {
		typ := ""
		if p.Type != "" {
			typ = " {" + p.Type + "}"
		}
		desc := ""
		if p.Description != "" {
			desc = " - " + p.Description
		}
		fmt.Fprintf(b, "   - %s: `%s`%s%s\n", label, p.Name, typ, desc)
	}
}

// withDoc returns the declared documentation suffix for a `with:` key when the step's
// uses: target is a local composite action in the scanned set: the input's description,
// plus a "(required)" marker.
func withDoc(step *model.Step, key string) string {
	if step.UsesAction == nil {
		return ""
	}
	in := step.UsesAction.Input(key)
	if in == nil {
		return ""
	}
	var s string
	if in.Description != "" {
		s = " - " + oneLine(in.Description)
	}
	if in.Required {
		s += " (required)"
	}
	return s
}

// stepTitle picks the most readable heading for a step: name, then id, then a friendly
// uses: ref, then the first non-comment line of run:, then a positional fallback.
func stepTitle(step *model.Step, num int) string {
	switch {
	case step.Name != "":
		return step.Name
	case step.ID != "":
		return step.ID
	case step.Uses != "":
		return friendlyUses(step.Uses, step.UsesVersion)
	}
	if first := firstRunLine(step.Run); first != "" {
		return first
	}
	return fmt.Sprintf("Step %d", num)
}

// friendlyUses collapses a SHA-pinned action ref to its human form: `owner/repo@v4.1.1`
// when a trailing version comment is present, or `owner/repo` when only the bare SHA is.
// Non-SHA refs (tags, branches, local paths) pass through unchanged.
func friendlyUses(uses, version string) string {
	at := strings.LastIndex(uses, "@")
	if at < 0 {
		return uses
	}
	ref, pin := uses[:at], uses[at+1:]
	if !isSHA(pin) {
		return uses
	}
	if version != "" {
		return ref + "@" + version
	}
	return ref
}

// usesSuffix returns the parenthetical version annotation shown after a SHA-pinned uses:
// ref on its detail line, so the exact commit pin stays visible alongside the version.
func usesSuffix(uses, version string) string {
	if version == "" {
		return ""
	}
	if at := strings.LastIndex(uses, "@"); at >= 0 && isSHA(uses[at+1:]) {
		return " (" + version + ")"
	}
	return ""
}

// isSHA reports whether s is a 40-character hexadecimal commit SHA.
func isSHA(s string) bool {
	if len(s) != 40 {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if !(c >= '0' && c <= '9' || c >= 'a' && c <= 'f' || c >= 'A' && c <= 'F') {
			return false
		}
	}
	return true
}

// firstRunLine returns the first non-blank, non-comment line of a run: script, truncated,
// for use as a step title when nothing better is available.
func firstRunLine(run string) string {
	for _, line := range strings.Split(run, "\n") {
		t := strings.TrimSpace(line)
		if t == "" || strings.HasPrefix(t, "#") {
			continue
		}
		return truncate(t, 60)
	}
	return ""
}

// truncate shortens s to at most max characters (runes), appending "..." when it cuts.
// It counts and slices by rune so a multi-byte UTF-8 character is never split, which
// would otherwise emit invalid UTF-8 in the output. Behavior is identical for ASCII.
func truncate(s string, max int) string {
	r := []rune(s)
	if len(r) <= max {
		return s
	}
	if max < 3 {
		return string(r[:max])
	}
	return string(r[:max-3]) + "..."
}

// oneLine collapses newlines to spaces so a value renders safely inside an inline code span.
func oneLine(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\n", " ")
	return strings.TrimSpace(s)
}

// resolveJobName expands ${{ matrix.X }} references in a job name to their static value
// lists (e.g. "Java ${{ matrix.java.version }}" -> "Java 17, 21, 24"). References without a
// statically-resolvable axis -- and all non-matrix expressions -- are left verbatim.
func resolveJobName(job *model.Job) string {
	name := job.Name
	if !strings.Contains(name, "${{") {
		return name
	}
	var b strings.Builder
	for {
		i := strings.Index(name, "${{")
		if i < 0 {
			b.WriteString(name)
			break
		}
		b.WriteString(name[:i])
		rest := name[i+3:]
		j := strings.Index(rest, "}}")
		if j < 0 {
			b.WriteString(name[i:])
			break
		}
		token := name[i : i+3+j+2]
		inner := strings.TrimSpace(rest[:j])
		if expressionKind(inner) == "matrix" {
			if vals, ok := job.MatrixValues(strings.TrimPrefix(inner, "matrix.")); ok {
				b.WriteString(strings.Join(vals, ", "))
			} else {
				b.WriteString(token)
			}
		} else {
			b.WriteString(token)
		}
		name = rest[j+2:]
	}
	return b.String()
}

// renderReferences writes the auto-collected "Referenced secrets and variables" section.
func renderReferences(b *strings.Builder, refs model.References) {
	if refs.Empty() {
		return
	}
	b.WriteString("## Referenced secrets and variables\n\n")
	writeRefTable(b, "Secrets", refs.Secrets)
	writeRefTable(b, "Variables", refs.Vars)
}

func writeRefTable(b *strings.Builder, label string, refs []model.Reference) {
	if len(refs) == 0 {
		return
	}
	fmt.Fprintf(b, "**%s:**\n\n", label)
	b.WriteString("| Name | Used by |\n")
	b.WriteString("|------|---------|\n")
	for _, r := range refs {
		fmt.Fprintf(b, "| `%s` | %s |\n", escapeCell(r.Name), escapeCell(strings.Join(r.Sites, "; ")))
	}
	b.WriteString("\n")
}

// RenderTOC builds a table of contents linking each top-level heading title to its anchor,
// for navigating a single-file render of many workflows/actions. Returns "" for fewer than
// two entries. Duplicate titles get GitHub's `-N` anchor suffixes in document order.
func RenderTOC(titles []string) string {
	if len(titles) < 2 {
		return ""
	}
	var b strings.Builder
	b.WriteString("# Contents\n\n")
	seen := map[string]int{}
	for _, t := range titles {
		base := anchor(t)
		slug := base
		if n := seen[base]; n > 0 {
			slug = fmt.Sprintf("%s-%d", base, n)
		}
		seen[base]++
		fmt.Fprintf(&b, "- [%s](#%s)\n", mdLinkLabel(t), slug)
	}
	b.WriteString("\n")
	return b.String()
}

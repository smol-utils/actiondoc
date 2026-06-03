package renderer

import (
	"fmt"
	"strings"

	"github.com/smol-utils/actiondoc/internal/model"
)

// renderStep writes one step as a numbered list item with its metadata (id, uses, condition,
// with: inputs, env: variables) and a continue-on-error badge when set.
func renderStep(b *strings.Builder, step *model.Step, num int) {
	fmt.Fprintf(b, "%d. **%s**", num, escapeInline(stepTitle(step, num)))
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
		fmt.Fprintf(b, "   - Condition: %s\n", codeSpan(oneLine(step.If)))
	}
	if len(step.With) > 0 {
		b.WriteString("   - With:\n")
		for _, kv := range step.With {
			fmt.Fprintf(b, "     - `%s`: %s%s\n", kv.Key, codeSpan(oneLine(kv.Value)), withDoc(step, kv.Key))
		}
	}
	if len(step.Env) > 0 {
		b.WriteString("   - Env:\n")
		for _, kv := range step.Env {
			fmt.Fprintf(b, "     - `%s`: %s\n", kv.Key, codeSpan(oneLine(kv.Value)))
		}
	}

	// Step-level tags
	writeStepParams(b, "Input", step.Tags.Inputs)
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

// stepTitle picks the most readable heading for a step: the shared step label (name, id,
// or collapsed uses: ref), then the first meaningful run: line, then a positional
// fallback.
func stepTitle(step *model.Step, num int) string {
	if step.Name != "" || step.ID != "" || step.Uses != "" {
		return step.Label(num)
	}
	if first := firstRunLine(step.Run); first != "" {
		return first
	}
	return fmt.Sprintf("Step %d", num)
}

// usesSuffix returns the parenthetical version annotation shown after a SHA-pinned uses:
// ref on its detail line, so the exact commit pin stays visible alongside the version.
func usesSuffix(uses, version string) string {
	if version == "" {
		return ""
	}
	if at := strings.LastIndex(uses, "@"); at >= 0 && model.IsSHA(uses[at+1:]) {
		return " (" + version + ")"
	}
	return ""
}

// firstRunLine returns the first non-blank, non-comment line of a run: script that
// contains at least one letter or digit, truncated, for use as a step title when nothing
// better is available. Punctuation-only lines (a shell group's opening "{", a lone
// parenthesis) carry no meaning as a title.
func firstRunLine(run string) string {
	for _, line := range strings.Split(run, "\n") {
		t := strings.TrimSpace(line)
		if t == "" || strings.HasPrefix(t, "#") || !hasAlphanumeric(t) {
			continue
		}
		return truncate(t, 60)
	}
	return ""
}

// hasAlphanumeric reports whether s contains at least one letter or digit.
func hasAlphanumeric(s string) bool {
	for _, r := range s {
		if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' {
			return true
		}
	}
	return false
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
// two entries. Duplicate titles get GitHub's `-N` anchor suffixes in document order, via
// the same AssignAnchors pass that cross-links use.
func RenderTOC(titles []string) string {
	if len(titles) < 2 {
		return ""
	}
	var b strings.Builder
	b.WriteString("# Contents\n\n")
	slugs := AssignAnchors(titles)
	for i, t := range titles {
		fmt.Fprintf(&b, "- [%s](#%s)\n", mdLinkLabel(t), slugs[i])
	}
	b.WriteString("\n")
	return b.String()
}

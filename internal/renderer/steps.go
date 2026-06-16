package renderer

import (
	"fmt"
	"strings"
	"unicode"

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
		if disp := usesDisplay(step.Uses, step.UsesVersion); disp != stepTitle(step, num) {
			// Skip the Uses: line only when the bold title already shows exactly this ref
			// (an unnamed step titled by its collapsed pin). A named/id'd step keeps the
			// line so the action it runs is never hidden.
			fmt.Fprintf(b, "   - Uses: `%s`\n", disp)
		}
	}
	if step.If != "" {
		fmt.Fprintf(b, "   - Condition: %s\n", codeSpan(oneLine(step.If)))
	}
	writeStepPairs(b, "With", step.With, func(key string) string { return withDoc(step, key) })
	writeStepPairs(b, "Env", step.Env, nil)

	// Step-level tags
	writeStepParams(b, "Input", step.Tags.Inputs)
	writeStepParams(b, "Output", step.Tags.Outputs)
	writeStepParams(b, "Secret", step.Tags.Secrets)
	writeStepParams(b, "Env", step.Tags.Envs)

	b.WriteString("\n")
}

// writeStepPairs writes a step's `with:`/`env:` key/value pairs as an indented bullet list
// under a labelled header. Entries whose value is empty/unset are omitted (a bare key with no
// value carries no information and otherwise renders as a column of dashes); the header is
// written only when at least one entry has a real value. doc, when non-nil, supplies a
// trailing documentation suffix for a key (e.g. a composite-action input's description).
func writeStepPairs(b *strings.Builder, label string, pairs []model.KV, doc func(key string) string) {
	type kept struct {
		key, value, suffix string
	}
	var entries []kept
	for _, kv := range pairs {
		// Omit entries whose value is empty/unset (a bare key carries no information).
		// The raw value is retained so stepValue can summarize long inline values.
		if oneLine(kv.Value) == "" {
			continue
		}
		suffix := ""
		if doc != nil {
			suffix = doc(kv.Key)
		}
		entries = append(entries, kept{key: kv.Key, value: kv.Value, suffix: suffix})
	}
	if len(entries) == 0 {
		return
	}
	fmt.Fprintf(b, "   - %s:\n", label)
	for _, e := range entries {
		fmt.Fprintf(b, "     - `%s`: %s%s\n", e.key, stepValue(e.value), e.suffix)
	}
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

// usesDisplay is the ref shown on a step's Uses: line. A SHA pin with a known
// human-readable version collapses to `owner/repo@version` -- the 40-character commit SHA
// adds no signal a reader uses, and the version is already the title's form. A bare SHA pin
// (no version) keeps its full `owner/repo@sha` so the exact pin is never lost. Tags,
// branches, and local paths pass through unchanged.
func usesDisplay(uses, version string) string {
	if at := strings.LastIndex(uses, "@"); at >= 0 && model.IsSHA(uses[at+1:]) && version != "" {
		return uses[:at] + "@" + version
	}
	return uses
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

// hasAlphanumeric reports whether s contains at least one letter or digit in any script,
// so a run: line written in non-Latin text still counts as meaningful title material.
func hasAlphanumeric(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
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

// stepValueCharThreshold bounds how wide a step value (a with:/env: value) may be before it
// is summarized rather than dumped whole. A value longer than two lines, or wider than this
// many characters, is reduced to its first line plus a "... (+N more lines)" marker.
const stepValueCharThreshold = 200

// stepValue renders a step's with:/env: value for the step summary. A short value renders
// whole on a single line (newlines collapsed) as an inline code span. A long one -- more than
// two lines, or longer than stepValueCharThreshold characters -- is reduced to its first
// non-blank line (itself clamped to the width threshold) followed by a "... (+N more lines)"
// marker, so a whole inline script (e.g. a github-script `script:` value) is summarized
// instead of pasted into the doc. The full content lives in the source workflow.
func stepValue(s string) string {
	trimmed := strings.TrimRight(strings.ReplaceAll(s, "\r\n", "\n"), "\n")
	lines := strings.Split(trimmed, "\n")
	if len(lines) <= 2 && len(trimmed) <= stepValueCharThreshold {
		return codeSpan(oneLine(s))
	}
	// Lead with the first line that carries content, so a block scalar whose first physical
	// line is blank still shows a meaningful summary line.
	first := ""
	for _, ln := range lines {
		if t := strings.TrimSpace(ln); t != "" {
			first = t
			break
		}
	}
	if first == "" {
		first = strings.TrimSpace(lines[0])
	}
	// Clamp the displayed line without truncate's "..." -- the trailing marker already signals
	// that content was dropped.
	if r := []rune(first); len(r) > stepValueCharThreshold {
		first = string(r[:stepValueCharThreshold])
	}
	if more := len(lines) - 1; more > 0 {
		noun := "lines"
		if more == 1 {
			noun = "line"
		}
		return fmt.Sprintf("%s ... (+%d more %s)", codeSpan(first), more, noun)
	}
	return codeSpan(first) + " ..."
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
		fmt.Fprintf(b, "| `%s` | %s |\n", escapeCell(r.Name), usedByCell(r.Sites))
	}
	b.WriteString("\n")
}

// usedByCell renders a reference's usage sites for the "Used by" column, grouped one line per
// job so the cell reads as a short list instead of a semicolon-joined wall. Jobs appear in
// first-seen (source) order; within a job, sites keep source order. Each job line reads
// `<job-id>: <step> (<name>), ...`. A job-level site (no step: a job env/with/secrets/if)
// renders as just `(<name>)`; workflow-level env sites group under a leading `workflow env`
// line. Lines are joined with <br> so the GitHub table cell shows one job per visual line.
func usedByCell(sites []model.Site) string {
	if len(sites) == 0 {
		return ""
	}
	var order []string // group keys (job ids; "" for the workflow level) in first-seen order
	groups := map[string][]model.Site{}
	for _, s := range sites {
		if _, ok := groups[s.Job]; !ok {
			order = append(order, s.Job)
		}
		groups[s.Job] = append(groups[s.Job], s)
	}
	lines := make([]string, 0, len(order))
	for _, job := range order {
		prefix := "workflow env: "
		if job != "" {
			prefix = codeSpan(escapeCell(job)) + ": "
		}
		entries := make([]string, 0, len(groups[job]))
		for _, s := range groups[job] {
			entries = append(entries, refSiteEntry(s))
		}
		lines = append(lines, prefix+strings.Join(entries, ", "))
	}
	return strings.Join(lines, "<br>")
}

// refSiteEntry renders one usage site within its job line: `<step> (<name>)`, or just
// `(<name>)` when the site has no step (a job-level env/with/secrets/if, or workflow env).
func refSiteEntry(s model.Site) string {
	name := "(" + codeSpan(escapeCell(s.Name)) + ")"
	if s.Step == "" {
		return name
	}
	return escapeCell(s.Step) + " " + name
}

// TOCEntry is one link in the table of contents: a fully-formed visible label (already
// carrying any trigger annotation or duplicate-name disambiguation) and the anchor slug it
// links to.
type TOCEntry struct {
	Label  string
	Anchor string
}

// TOCGroup is a labelled run of TOC entries (e.g. "Workflows"). Empty groups are skipped
// when rendered.
type TOCGroup struct {
	Heading string
	Entries []TOCEntry
}

// RenderTOC builds a grouped table of contents for navigating a single-file render of many
// workflows/actions. Each non-empty group is rendered under its bold heading in the order
// given. Returns "" when fewer than two entries exist across all groups (a lone document
// needs no contents list). Entry labels are escaped for use as Markdown link text; anchors
// are taken as-is from the caller's AssignAnchors pass so they agree with cross-links.
func RenderTOC(groups []TOCGroup) string {
	total := 0
	for _, g := range groups {
		total += len(g.Entries)
	}
	if total < 2 {
		return ""
	}
	var b strings.Builder
	b.WriteString("## Contents\n\n")
	for _, g := range groups {
		if len(g.Entries) == 0 {
			continue
		}
		fmt.Fprintf(&b, "**%s**\n\n", g.Heading)
		for _, e := range g.Entries {
			fmt.Fprintf(&b, "- [%s](#%s)\n", mdLinkLabel(e.Label), e.Anchor)
		}
		b.WriteString("\n")
	}
	return b.String()
}

// RenderDocumentHeader emits the document's H1 title and a one-line inventory of what it
// documents. Counts that are zero are omitted; the rest are pluralized. It returns "" when
// title is empty so callers can suppress the header for single-document output.
func RenderDocumentHeader(title string, workflows, reusable, composites int) string {
	if title == "" {
		return ""
	}
	var b strings.Builder
	fmt.Fprintf(&b, "# %s\n\n", title)
	var parts []string
	if workflows > 0 {
		parts = append(parts, pluralize(workflows, "workflow", "workflows"))
	}
	if reusable > 0 {
		parts = append(parts, pluralize(reusable, "reusable workflow", "reusable workflows"))
	}
	if composites > 0 {
		parts = append(parts, pluralize(composites, "composite action", "composite actions"))
	}
	if len(parts) > 0 {
		b.WriteString(strings.Join(parts, ", "))
		b.WriteString("\n\n")
	}
	return b.String()
}

// pluralize formats a count with its singular or plural noun.
func pluralize(n int, singular, plural string) string {
	if n == 1 {
		return fmt.Sprintf("%d %s", n, singular)
	}
	return fmt.Sprintf("%d %s", n, plural)
}

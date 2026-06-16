package renderer

import (
	"fmt"
	"strings"

	"github.com/smol-utils/actiondoc/internal/callgraph"
	"github.com/smol-utils/actiondoc/internal/model"
)

// RenderMarkdown converts a Workflow IR into a Markdown document with no call-graph
// context (single-file rendering). Reusable-workflow cross-links and the call-graph
// sections are omitted; use RenderMarkdownGraph to include them.
func RenderMarkdown(w *model.Workflow) string {
	return RenderMarkdownGraph(w, nil, "")
}

// RenderMarkdownGraph converts a Workflow IR into a Markdown document, using the call
// graph g (built from the whole scan set) to resolve reusable-workflow cross-links and to
// render the call-graph / "called by" / transitive-requirements sections. id is this
// workflow's node id in g. g may be nil and id empty, in which case the graph-derived
// sections are skipped and the output matches single-file rendering.
func RenderMarkdownGraph(w *model.Workflow, g *callgraph.Graph, id string) string {
	var b strings.Builder

	// Title
	fmt.Fprintf(&b, "# %s\n\n", w.Name)

	// Deprecated banner
	if w.Tags.Deprecated != "" {
		fmt.Fprintf(&b, "> **Deprecated**: %s\n\n", w.Tags.Deprecated)
	}

	// Triggers: the fact a reader scans for first, promoted to a prominent line directly
	// under the heading instead of being buried as one row in the properties table.
	if len(w.On) > 0 {
		fmt.Fprintf(&b, "**Triggers:** %s\n\n", codelist(w.On))
	}

	// Description (rendered as written; never fabricated when absent). A long
	// description is folded so it does not bury the property table below.
	renderDescription(&b, w.Description)

	// Properties table (the Triggers row is promoted to the line above)
	b.WriteString("| Property | Value |\n")
	b.WriteString("|----------|-------|\n")
	fmt.Fprintf(&b, "| File | `%s` |\n", w.File)
	// The runs-on value shared by most jobs is stated once here; per-job tables then show
	// runs-on only where a job departs from it, instead of repeating the same (often long)
	// runner string on every job.
	defaultRunsOn := commonRunsOn(w.Jobs)
	if defaultRunsOn != "" {
		fmt.Fprintf(&b, "| Default runs-on | `%s` |\n", escapeCell(defaultRunsOn))
	}
	if w.Tags.Since != "" {
		fmt.Fprintf(&b, "| Since | %s |\n", w.Tags.Since)
	}
	b.WriteString("\n")

	writeSeeAlso(&b, w.Tags.See)

	// Job roster: a compact, scannable list of this workflow's jobs linking into the Jobs
	// section below, placed before the deep content so the reader sees the roster up front.
	// Job heading anchors are assigned document-wide by the assembler (stored on this node);
	// nil with no graph context falls back to local slugging in renderJobMiniTOC.
	var jobAnchors []string
	if g != nil {
		if n := g.Nodes[id]; n != nil {
			jobAnchors = n.JobAnchors
		}
	}
	renderJobMiniTOC(&b, w.Jobs, jobAnchors)

	renderWorkflowSurface(&b, w)

	// Call-graph sections (rendered only with graph context): downstream tree and
	// transitive requirements on entry points, upstream chain on reusable workflows.
	renderCallGraph(&b, g, id)
	renderTransitiveRequirements(&b, g, id)
	renderCalledBy(&b, g, id)

	writeParamSections(&b, styleHeading,
		paramSection{"Secrets", w.Tags.Secrets},
		paramSection{"Inputs", w.Tags.Inputs},
		paramSection{"Environment Variables", w.Tags.Envs},
		paramSection{"Outputs", w.Tags.Outputs},
	)

	writeExample(&b, styleHeading, w.Tags.Example)

	// Auto-collected secret and variable references found in expressions.
	renderReferences(&b, model.ScanReferences(w))

	// Jobs
	if len(w.Jobs) > 0 {
		b.WriteString("## Jobs\n\n")
		for i := range w.Jobs {
			renderJob(&b, &w.Jobs[i], g, id, defaultRunsOn)
		}
	}

	return b.String()
}

// renderJobMiniTOC writes a one-line roster of a workflow's jobs, each linking to its job
// heading, so a reader sees the job list without scrolling down to the Jobs section. A lone
// job needs no roster, so the line is emitted only for two or more jobs.
//
// anchors carries the job heading slugs assigned DOCUMENT-WIDE by the assembler (one per
// job, in job order): GitHub disambiguates repeated heading slugs across the whole rendered
// document, so a job heading text that recurs in a later workflow must keep the running "-N"
// suffix. Assigning per-workflow would restart the count and point a later workflow's link at
// the first occurrence in another workflow. When anchors is nil (single-file rendering, no
// assembly) the whole document is this one workflow, so per-workflow disambiguation equals
// document-wide and we slug the local jobs directly.
//
// Residual edge: GitHub numbers ALL same-slug headings together regardless of kind (job,
// section, H1, "#### Steps"). This unifies job-vs-job collisions only; a job slug that also
// collides with a non-job heading slug is not reconciled. That cross-kind case is rare and
// covering it would require threading every heading through one global pass.
func renderJobMiniTOC(b *strings.Builder, jobs []model.Job, anchors []string) {
	if len(jobs) < 2 {
		return
	}
	slugs := anchors
	if slugs == nil {
		texts := make([]string, len(jobs))
		for i := range jobs {
			texts[i] = JobHeadingText(&jobs[i])
		}
		slugs = AssignAnchors(texts)
	}
	parts := make([]string, len(jobs))
	for i := range jobs {
		parts[i] = fmt.Sprintf("[%s](#%s)", mdLinkLabel(jobMiniLabel(&jobs[i])), slugs[i])
	}
	// A compact comma-joined line stays scannable for a handful of jobs, but past this
	// threshold it becomes an unscannable wall of links right above the headings it mirrors,
	// so switch to a vertical bulleted list (one job per line). Same anchors, same labels.
	if len(jobs) > jobMiniTOCInlineMax {
		// On a giant workflow even the vertical list dominates the front matter, so past a
		// second threshold collapse it behind a <details>. GitHub renders the Markdown list
		// inside only when a blank line follows the summary.
		if len(jobs) > jobMiniTOCDetailsMax {
			fmt.Fprintf(b, "<details>\n<summary>Jobs (%d)</summary>\n\n", len(jobs))
			for _, p := range parts {
				fmt.Fprintf(b, "- %s\n", p)
			}
			b.WriteString("\n</details>\n\n")
			return
		}
		b.WriteString("**Jobs:**\n\n")
		for _, p := range parts {
			fmt.Fprintf(b, "- %s\n", p)
		}
		b.WriteString("\n")
		return
	}
	fmt.Fprintf(b, "**Jobs:** %s\n\n", strings.Join(parts, ", "))
}

// jobMiniTOCInlineMax is the largest job count rendered as a single inline comma-joined
// mini-TOC line; above it, the roster becomes a vertical bulleted list for scannability.
const jobMiniTOCInlineMax = 8

// jobMiniTOCDetailsMax is the largest job count whose vertical roster renders open inline;
// above it, the bulleted list is collapsed behind a <details> so it does not dominate the
// front matter of a giant workflow.
const jobMiniTOCDetailsMax = 15

// JobHeadingText returns the visible text of a job's heading: the basis for its GitHub anchor
// slug. It mirrors renderJob's heading construction so a mini-TOC link resolves to the
// heading GitHub actually emits (backticks and parentheses drop out of the slug either way).
// Exported so the assembler can feed the exact same text into its document-wide anchor pass.
func JobHeadingText(job *model.Job) string {
	if name := jobDisplayName(job); name != job.ID {
		return name + " (" + job.ID + ")"
	}
	return job.ID
}

// jobMiniLabel is a job's visible label in the mini-TOC: its name when distinct from the id,
// otherwise the id rendered as inline code (matching the job heading's own treatment).
func jobMiniLabel(job *model.Job) string {
	if name := jobDisplayName(job); name != job.ID {
		return name
	}
	return "`" + job.ID + "`"
}

// commonRunsOn returns the runs-on value shared by the most jobs (the workflow's de facto
// default), or "" when no value is shared by at least two jobs. Caller jobs (uses:) carry
// no runs-on and are skipped. Ties are broken by first appearance so the choice is
// deterministic for a given job order.
func commonRunsOn(jobs []model.Job) string {
	counts := map[string]int{}
	var order []string
	for i := range jobs {
		r := jobs[i].RunsOn
		if r == "" {
			continue
		}
		if counts[r] == 0 {
			order = append(order, r)
		}
		counts[r]++
	}
	best, bestN := "", 0
	for _, r := range order {
		if counts[r] > bestN {
			best, bestN = r, counts[r]
		}
	}
	if bestN < 2 {
		return ""
	}
	return best
}

func renderJob(b *strings.Builder, job *model.Job, g *callgraph.Graph, fromID, defaultRunsOn string) {
	// Job heading. A name embedding ${{ ... }} expressions is de-templated to a readable
	// label at jobDisplayName (the single origin every consumer shares), so the raw
	// expression never leaks into the heading -- and the same label flows into the anchor,
	// mini-TOC, and call graph. Placeholders are still never expanded into joined value
	// lists (GitHub creates one job per matrix combination; "Java 17, 21" is a job name
	// that never exists); the Matrix property row below shows the axis values they take.
	name := jobDisplayName(job)
	if name != job.ID {
		fmt.Fprintf(b, "### %s (`%s`)\n\n", escapeInline(name), job.ID)
	} else {
		fmt.Fprintf(b, "### `%s`\n\n", job.ID)
	}

	// Deprecated
	if job.Tags.Deprecated != "" {
		fmt.Fprintf(b, "> **Deprecated**: %s\n\n", job.Tags.Deprecated)
	}

	// Description
	renderDescription(b, job.Description)

	// A job that calls a reusable workflow uses `uses:` instead of `runs-on:`/`steps:`;
	// render its caller surface (callee link + forwarded inputs/secrets) and stop.
	if job.Uses != "" {
		renderCallerJob(b, job, g, fromID)
		return
	}

	// Properties table. The runs-on row is shown only when this job departs from the
	// workflow's stated default (or when there is no default to hoist), so the common runner
	// string is not repeated on every job.
	showRunsOn := job.RunsOn != "" && job.RunsOn != defaultRunsOn
	hasProps := showRunsOn || len(job.Needs) > 0 || job.If != "" || len(job.Matrix) > 0
	if hasProps {
		b.WriteString("| Property | Value |\n")
		b.WriteString("|----------|-------|\n")
		if showRunsOn {
			fmt.Fprintf(b, "| Runs on | `%s` |\n", escapeCell(job.RunsOn))
		}
		if len(job.Matrix) > 0 {
			fmt.Fprintf(b, "| Matrix | %s |\n", matrixCell(job.Matrix, job.MatrixAdjusted))
		}
		if len(job.Needs) > 0 {
			fmt.Fprintf(b, "| Depends on | %s |\n", codelist(job.Needs))
		}
		if job.If != "" {
			// Trim first: literal-block conditions carry a trailing newline that would
			// otherwise render as a dangling <br>.
			fmt.Fprintf(b, "| Condition | `%s` |\n", escapeCell(strings.TrimSpace(job.If)))
		}
		b.WriteString("\n")
	}

	renderJobSurface(b, job)
	renderJobTags(b, job)

	// Steps. The per-step detail is the bulkiest, least-scanned part of a job, so it is
	// folded behind a <details> whose summary states the step count. GitHub renders the
	// Markdown inside only when a blank line follows the summary, so that blank line is
	// load-bearing.
	if len(job.Steps) > 0 {
		fmt.Fprintf(b, "<details>\n<summary>Steps (%d)</summary>\n\n", len(job.Steps))
		for i, step := range job.Steps {
			renderStep(b, &step, i+1)
		}
		b.WriteString("</details>\n\n")
	}
}

// renderJobTags writes a job's ActionDoc tag sections (@secret/@env/@output/@example/
// @see). Shared by the normal and reusable-workflow-caller job renderers so caller jobs
// don't silently drop tags the spec allows on jobs.
func renderJobTags(b *strings.Builder, job *model.Job) {
	writeParamSections(b, styleBold,
		paramSection{"Secrets", job.Tags.Secrets},
		paramSection{"Inputs", job.Tags.Inputs},
		paramSection{"Environment Variables", job.Tags.Envs},
		paramSection{"Outputs", job.Tags.Outputs},
	)

	writeExample(b, styleBold, job.Tags.Example)
	writeSeeAlso(b, job.Tags.See)
}

// renderStep and writeStepParams live in steps.go.

// writeParamTable writes a Markdown table for a slice of Params.
func writeParamTable(b *strings.Builder, params []model.Param) {
	b.WriteString("| Name | Type | Description |\n")
	b.WriteString("|------|------|-------------|\n")
	for _, p := range params {
		typ := p.Type
		if typ == "" {
			typ = "-"
		}
		desc := p.Description
		if desc == "" {
			desc = "-"
		}
		fmt.Fprintf(b, "| `%s` | %s | %s |\n", escapeCell(p.Name), escapeCell(typ), escapeCell(desc))
	}
	b.WriteString("\n")
}

// sectionStyle selects how a param-table section heading is rendered.
type sectionStyle int

const (
	styleHeading sectionStyle = iota // "## Title"
	styleBold                        // "**Title:**"
)

// paramSection pairs a section title with its parameters.
type paramSection struct {
	title  string
	params []model.Param
}

// writeSeeAlso writes the @see links line, shared by the workflow, job, and action
// renderers.
func writeSeeAlso(b *strings.Builder, see []string) {
	if len(see) == 0 {
		return
	}
	b.WriteString("**See also:** ")
	b.WriteString(strings.Join(see, ", "))
	b.WriteString("\n\n")
}

// writeExample writes the fenced @example block with its heading in the given section
// style, shared by the workflow, job, and action renderers.
func writeExample(b *strings.Builder, style sectionStyle, example string) {
	if example == "" {
		return
	}
	if style == styleHeading {
		b.WriteString("## Example\n\n")
	} else {
		b.WriteString("**Example:**\n\n")
	}
	fmt.Fprintf(b, "```\n%s\n```\n\n", example)
}

// writeParamSections writes each non-empty param-table section in order, using the given
// heading style. Centralizes the "if present: heading + table" boilerplate shared by the
// workflow, job, and action renderers so adding a section is a one-line change per site.
func writeParamSections(b *strings.Builder, style sectionStyle, sections ...paramSection) {
	for _, s := range sections {
		if len(s.params) == 0 {
			continue
		}
		if style == styleBold {
			fmt.Fprintf(b, "**%s:**\n\n", s.title)
		} else {
			fmt.Fprintf(b, "## %s\n\n", s.title)
		}
		writeParamTable(b, s.params)
	}
}

// descriptionInlineMax is the longest single-line description rendered verbatim. Past it
// (or whenever the description spans multiple lines), only the lead sentence/line stays
// inline and the remainder folds behind a <details>, so a long boilerplate comment block
// does not push the property table far below the heading.
const descriptionInlineMax = 200

// renderDescription writes a heading's description paragraph (shared by the workflow, job,
// and action renderers so all three fold consistently). A short authored description renders
// verbatim. A long one -- multi-line, or longer than descriptionInlineMax characters -- keeps
// only its first sentence (or first line) inline and tucks the remainder inside a
// <details>/<summary>, so the scannable lead stays next to the heading while the bulk is one
// click away. An empty/whitespace description renders nothing (descriptions are never
// fabricated when absent).
func renderDescription(b *strings.Builder, desc string) {
	desc = strings.TrimSpace(desc)
	if desc == "" {
		return
	}
	lead, rest := splitDescriptionLead(desc)
	if rest == "" {
		fmt.Fprintf(b, "%s\n\n", desc)
		return
	}
	fmt.Fprintf(b, "%s\n\n", lead)
	// GitHub renders the Markdown inside <details> only when a blank line follows the
	// summary, so that blank line is load-bearing.
	fmt.Fprintf(b, "<details>\n<summary>more</summary>\n\n%s\n\n</details>\n\n", rest)
}

// splitDescriptionLead separates a description into an inline lead and a foldable remainder.
// It returns rest == "" (signalling render-whole) when the description is short: a single
// line within descriptionInlineMax characters, or a longer run with no usable split point.
// Otherwise the lead is the first sentence (text up to the first ". " or ".\n"), or, when
// there is no sentence boundary, the first line; rest is everything after it, trimmed.
func splitDescriptionLead(desc string) (lead, rest string) {
	if !strings.Contains(desc, "\n") && len(desc) <= descriptionInlineMax {
		return desc, ""
	}
	// Prefer a sentence boundary: the first period followed by a space or newline.
	if i := sentenceEnd(desc); i > 0 && i < len(desc) {
		return strings.TrimSpace(desc[:i]), strings.TrimSpace(desc[i:])
	}
	// No sentence boundary: fall back to the first line when the text is multi-line.
	if i := strings.IndexByte(desc, '\n'); i >= 0 {
		return strings.TrimSpace(desc[:i]), strings.TrimSpace(desc[i+1:])
	}
	// A single long line with no sentence boundary: render whole rather than chop mid-word.
	return desc, ""
}

// sentenceEnd returns the index just past the first sentence-ending period (a '.' followed
// by a space or newline), or -1 when the text has no such boundary. The returned index
// includes the period so the lead keeps its terminating punctuation.
func sentenceEnd(s string) int {
	for i := 0; i+1 < len(s); i++ {
		if s[i] == '.' && (s[i+1] == ' ' || s[i+1] == '\n') {
			return i + 1
		}
	}
	return -1
}

// codelist formats a slice of strings as inline code items.
func codelist(items []string) string {
	parts := make([]string, len(items))
	for i, s := range items {
		parts[i] = "`" + s + "`"
	}
	return strings.Join(parts, ", ")
}

// matrixCell formats a job's declared matrix axes for the properties table: each axis as
// `name`: v1, v2, joined with "; ". The job heading shows the name template as written;
// this row is what tells the reader which values its placeholders take. When the matrix
// also has include:/exclude: entries, the listed values are not the exact combination
// set, and the cell says so.
func matrixCell(axes []model.MatrixAxis, adjusted bool) string {
	parts := make([]string, len(axes))
	for i, a := range axes {
		parts[i] = codeSpan(escapeCell(a.Name)) + ": " + escapeCell(strings.Join(a.Values, ", "))
	}
	s := strings.Join(parts, "; ")
	if adjusted {
		s += " (combinations adjusted by include/exclude)"
	}
	return s
}

// RenderActionMarkdown converts an Action data model into a Markdown document.
func RenderActionMarkdown(a *model.Action) string {
	var b strings.Builder

	fmt.Fprintf(&b, "# %s\n\n", a.Name)

	if a.Tags.Deprecated != "" {
		fmt.Fprintf(&b, "> **Deprecated**: %s\n\n", a.Tags.Deprecated)
	}

	renderDescription(&b, a.Description)

	// Properties table
	b.WriteString("| Property | Value |\n")
	b.WriteString("|----------|-------|\n")
	fmt.Fprintf(&b, "| File | `%s` |\n", a.File)
	if a.Runs.Using != "" {
		fmt.Fprintf(&b, "| Runs with | `%s` |\n", a.Runs.Using)
	}
	if a.Tags.Since != "" {
		fmt.Fprintf(&b, "| Since | %s |\n", a.Tags.Since)
	}
	b.WriteString("\n")

	writeSeeAlso(&b, a.Tags.See)

	// Inputs
	if len(a.Inputs) > 0 {
		b.WriteString("## Inputs\n\n")
		b.WriteString("| Name | Description | Required | Default |\n")
		b.WriteString("|------|-------------|----------|--------|\n")
		for _, in := range a.Inputs {
			req := "No"
			if in.Required {
				req = "Yes"
			}
			def := "-"
			if in.Default != "" {
				def = "`" + escapeCell(in.Default) + "`"
			}
			fmt.Fprintf(&b, "| `%s` | %s | %s | %s |\n",
				escapeCell(in.Name), escapeCell(in.Description), req, def)
		}
		b.WriteString("\n")
	}

	// Outputs
	if len(a.Outputs) > 0 {
		b.WriteString("## Outputs\n\n")
		b.WriteString("| Name | Description |\n")
		b.WriteString("|------|-------------|\n")
		for _, out := range a.Outputs {
			desc := out.Description
			if desc == "" {
				desc = "-"
			}
			fmt.Fprintf(&b, "| `%s` | %s |\n", escapeCell(out.Name), escapeCell(desc))
		}
		b.WriteString("\n")
	}

	writeParamSections(&b, styleHeading,
		paramSection{"Secrets", a.Tags.Secrets},
		paramSection{"Environment Variables", a.Tags.Envs},
	)

	writeExample(&b, styleHeading, a.Tags.Example)

	return b.String()
}

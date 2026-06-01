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

	// Description
	if w.Description != "" {
		fmt.Fprintf(&b, "%s\n\n", w.Description)
	}

	// Properties table
	b.WriteString("| Property | Value |\n")
	b.WriteString("|----------|-------|\n")
	fmt.Fprintf(&b, "| File | `%s` |\n", w.File)
	if len(w.On) > 0 {
		fmt.Fprintf(&b, "| Triggers | %s |\n", codelist(w.On))
	}
	if w.Tags.Since != "" {
		fmt.Fprintf(&b, "| Since | %s |\n", w.Tags.Since)
	}
	b.WriteString("\n")

	// See also links
	if len(w.Tags.See) > 0 {
		b.WriteString("**See also:** ")
		b.WriteString(strings.Join(w.Tags.See, ", "))
		b.WriteString("\n\n")
	}

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

	// Auto-collected secret and variable references found in expressions.
	renderReferences(&b, model.ScanReferences(w))

	// Jobs
	if len(w.Jobs) > 0 {
		b.WriteString("## Jobs\n\n")
		for i := range w.Jobs {
			renderJob(&b, &w.Jobs[i], g, id)
		}
	}

	return b.String()
}

func renderJob(b *strings.Builder, job *model.Job, g *callgraph.Graph, fromID string) {
	// Job heading. The name may embed ${{ matrix.X }} references, expanded to value lists.
	name := resolveJobName(job)
	if name != job.ID {
		fmt.Fprintf(b, "### %s (`%s`)\n\n", name, job.ID)
	} else {
		fmt.Fprintf(b, "### `%s`\n\n", job.ID)
	}

	// Deprecated
	if job.Tags.Deprecated != "" {
		fmt.Fprintf(b, "> **Deprecated**: %s\n\n", job.Tags.Deprecated)
	}

	// Description
	if job.Description != "" {
		fmt.Fprintf(b, "%s\n\n", job.Description)
	}

	// A job that calls a reusable workflow uses `uses:` instead of `runs-on:`/`steps:`;
	// render its caller surface (callee link + forwarded inputs/secrets) and stop.
	if job.Uses != "" {
		renderCallerJob(b, job, g, fromID)
		return
	}

	// Properties table
	hasProps := job.RunsOn != "" || len(job.Needs) > 0 || job.If != ""
	if hasProps {
		b.WriteString("| Property | Value |\n")
		b.WriteString("|----------|-------|\n")
		if job.RunsOn != "" {
			fmt.Fprintf(b, "| Runs on | `%s` |\n", escapeCell(job.RunsOn))
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

	writeParamSections(b, styleBold,
		paramSection{"Secrets", job.Tags.Secrets},
		paramSection{"Environment Variables", job.Tags.Envs},
		paramSection{"Outputs", job.Tags.Outputs},
	)

	// Example
	if job.Tags.Example != "" {
		b.WriteString("**Example:**\n\n")
		fmt.Fprintf(b, "```\n%s\n```\n\n", job.Tags.Example)
	}

	// See also
	if len(job.Tags.See) > 0 {
		b.WriteString("**See also:** ")
		b.WriteString(strings.Join(job.Tags.See, ", "))
		b.WriteString("\n\n")
	}

	// Steps
	if len(job.Steps) > 0 {
		b.WriteString("#### Steps\n\n")
		for i, step := range job.Steps {
			renderStep(b, &step, i+1)
		}
	}
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

// codelist formats a slice of strings as inline code items.
func codelist(items []string) string {
	parts := make([]string, len(items))
	for i, s := range items {
		parts[i] = "`" + s + "`"
	}
	return strings.Join(parts, ", ")
}

// escapeCell escapes characters that break Markdown table cells. Newlines become
// <br> (not a space) so multi-line values like multi-line `if:` conditions keep their
// visual line breaks instead of collapsing or, worse, being parsed as a new table row.
func escapeCell(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "|", "\\|")
	s = strings.ReplaceAll(s, "\n", "<br>")
	return s
}

// RenderActionMarkdown converts an Action data model into a Markdown document.
func RenderActionMarkdown(a *model.Action) string {
	var b strings.Builder

	fmt.Fprintf(&b, "# %s\n\n", a.Name)

	if a.Tags.Deprecated != "" {
		fmt.Fprintf(&b, "> **Deprecated**: %s\n\n", a.Tags.Deprecated)
	}

	if a.Description != "" {
		fmt.Fprintf(&b, "%s\n\n", a.Description)
	}

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

	if len(a.Tags.See) > 0 {
		b.WriteString("**See also:** ")
		b.WriteString(strings.Join(a.Tags.See, ", "))
		b.WriteString("\n\n")
	}

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

	// Example
	if a.Tags.Example != "" {
		b.WriteString("## Example\n\n")
		fmt.Fprintf(&b, "```\n%s\n```\n\n", a.Tags.Example)
	}

	return b.String()
}

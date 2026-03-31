package renderer

import (
	"fmt"
	"strings"

	"github.com/smol-utils/actiondoc/internal/model"
)

// RenderMarkdown converts a Workflow IR into a Markdown document.
func RenderMarkdown(w *model.Workflow) string {
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

	// Workflow-level secrets
	if len(w.Tags.Secrets) > 0 {
		b.WriteString("## Secrets\n\n")
		writeParamTable(&b, w.Tags.Secrets)
	}

	// Workflow-level inputs
	if len(w.Tags.Inputs) > 0 {
		b.WriteString("## Inputs\n\n")
		writeParamTable(&b, w.Tags.Inputs)
	}

	// Workflow-level env
	if len(w.Tags.Envs) > 0 {
		b.WriteString("## Environment Variables\n\n")
		writeParamTable(&b, w.Tags.Envs)
	}

	// Workflow-level outputs
	if len(w.Tags.Outputs) > 0 {
		b.WriteString("## Outputs\n\n")
		writeParamTable(&b, w.Tags.Outputs)
	}

	// Jobs
	if len(w.Jobs) > 0 {
		b.WriteString("## Jobs\n\n")
		for _, job := range w.Jobs {
			renderJob(&b, &job)
		}
	}

	return b.String()
}

func renderJob(b *strings.Builder, job *model.Job) {
	// Job heading
	if job.Name != job.ID {
		fmt.Fprintf(b, "### %s (`%s`)\n\n", job.Name, job.ID)
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

	// Properties table
	hasProps := job.RunsOn != "" || len(job.Needs) > 0 || job.If != ""
	if hasProps {
		b.WriteString("| Property | Value |\n")
		b.WriteString("|----------|-------|\n")
		if job.RunsOn != "" {
			fmt.Fprintf(b, "| Runs on | `%s` |\n", job.RunsOn)
		}
		if len(job.Needs) > 0 {
			fmt.Fprintf(b, "| Depends on | %s |\n", codelist(job.Needs))
		}
		if job.If != "" {
			fmt.Fprintf(b, "| Condition | `%s` |\n", job.If)
		}
		b.WriteString("\n")
	}

	// Job-level secrets/envs/outputs
	if len(job.Tags.Secrets) > 0 {
		b.WriteString("**Secrets:**\n\n")
		writeParamTable(b, job.Tags.Secrets)
	}
	if len(job.Tags.Envs) > 0 {
		b.WriteString("**Environment Variables:**\n\n")
		writeParamTable(b, job.Tags.Envs)
	}
	if len(job.Tags.Outputs) > 0 {
		b.WriteString("**Outputs:**\n\n")
		writeParamTable(b, job.Tags.Outputs)
	}

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

func renderStep(b *strings.Builder, step *model.Step, num int) {
	name := step.Name
	if name == "" {
		name = step.ID
	}
	if name == "" && step.Uses != "" {
		name = step.Uses
	}
	if name == "" {
		name = fmt.Sprintf("Step %d", num)
	}

	fmt.Fprintf(b, "%d. **%s**", num, name)

	if step.Description != "" {
		fmt.Fprintf(b, " - %s", step.Description)
	}
	b.WriteString("\n")

	if step.ID != "" {
		fmt.Fprintf(b, "   - ID: `%s`\n", step.ID)
	}
	if step.Uses != "" {
		fmt.Fprintf(b, "   - Uses: `%s`\n", step.Uses)
	}
	if step.If != "" {
		fmt.Fprintf(b, "   - Condition: `%s`\n", step.If)
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

// codelist formats a slice of strings as inline code items.
func codelist(items []string) string {
	parts := make([]string, len(items))
	for i, s := range items {
		parts[i] = "`" + s + "`"
	}
	return strings.Join(parts, ", ")
}

// escapeCell escapes characters that break Markdown table cells.
func escapeCell(s string) string {
	s = strings.ReplaceAll(s, "|", "\\|")
	s = strings.ReplaceAll(s, "\n", " ")
	return s
}

package renderer

import (
	"fmt"
	"strings"

	"github.com/smol-utils/actiondoc/internal/model"
)

// renderWorkflowSurface writes the workflow-level declared surface (triggers,
// permissions, env, concurrency, defaults) using ## section headings. It is the single
// render call-site for this session's workflow-level additions.
func renderWorkflowSurface(b *strings.Builder, w *model.Workflow) {
	renderTriggers(b, w.Triggers)
	renderPermissions(b, w.Permissions, "## Permissions\n\n")
	renderEnv(b, w.Env, "## Environment (`env`)\n\n")
	renderConcurrency(b, w.Concurrency)
	renderDefaults(b, w.Defaults)
}

// renderJobSurface writes the job-level declared surface (environment binding,
// permissions, env, concurrency, defaults) using bold sub-headings, matching the
// job-level style. It is the single render call-site for this session's job-level
// additions.
func renderJobSurface(b *strings.Builder, job *model.Job) {
	renderEnvironment(b, job.Environment)
	renderPermissions(b, job.Permissions, "**Permissions:**\n\n")
	renderEnv(b, job.Env, "**Environment (`env`):**\n\n")
	renderConcurrency(b, job.Concurrency)
	renderDefaults(b, job.Defaults)
}

// renderEnv writes an env: block as a Variable/Value table under the given heading.
func renderEnv(b *strings.Builder, env []model.KV, heading string) {
	if len(env) == 0 {
		return
	}
	b.WriteString(heading)
	b.WriteString("| Variable | Value |\n")
	b.WriteString("|----------|-------|\n")
	for _, kv := range env {
		fmt.Fprintf(b, "| `%s` | %s |\n", escapeCell(kv.Key), codeCellOrDash(kv.Value))
	}
	b.WriteString("\n")
}

// renderConcurrency writes a concurrency: block as a one-line bold callout.
func renderConcurrency(b *strings.Builder, c *model.Concurrency) {
	if c == nil {
		return
	}
	fmt.Fprintf(b, "**Concurrency:** group `%s`", c.Group)
	if c.CancelInProgress != "" {
		fmt.Fprintf(b, ", cancel-in-progress: `%s`", c.CancelInProgress)
	}
	b.WriteString("\n\n")
}

// renderDefaults writes a defaults.run: block as a one-line bold callout.
func renderDefaults(b *strings.Builder, d *model.Defaults) {
	if d == nil {
		return
	}
	var parts []string
	if d.Shell != "" {
		parts = append(parts, fmt.Sprintf("shell `%s`", d.Shell))
	}
	if d.WorkingDirectory != "" {
		parts = append(parts, fmt.Sprintf("working-directory `%s`", d.WorkingDirectory))
	}
	fmt.Fprintf(b, "**Defaults:** %s\n\n", strings.Join(parts, ", "))
}

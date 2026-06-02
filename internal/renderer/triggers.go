package renderer

import (
	"fmt"
	"strings"

	"github.com/smol-utils/actiondoc/internal/model"
)

// renderTriggers writes the detailed trigger surface for a workflow: the manual-trigger
// inputs, the reusable-workflow public API, schedules, and per-event filters. Each
// subsection is emitted only when present.
func renderTriggers(b *strings.Builder, t *model.Triggers) {
	if t == nil {
		return
	}

	if t.Dispatch != nil && len(t.Dispatch.Inputs) > 0 {
		b.WriteString("## Manual trigger inputs\n\n")
		b.WriteString("Inputs for the `workflow_dispatch` event.\n\n")
		writeInputTable(b, t.Dispatch.Inputs)
	}

	if t.Call != nil {
		b.WriteString("## Workflow call API\n\n")
		b.WriteString("This workflow is reusable via `workflow_call`.\n\n")
		if len(t.Call.Inputs) > 0 {
			b.WriteString("**Inputs:**\n\n")
			writeInputTable(b, t.Call.Inputs)
		}
		if len(t.Call.Outputs) > 0 {
			b.WriteString("**Outputs:**\n\n")
			b.WriteString("| Name | Description | Value |\n")
			b.WriteString("|------|-------------|-------|\n")
			for _, o := range t.Call.Outputs {
				fmt.Fprintf(b, "| `%s` | %s | %s |\n",
					escapeCell(o.Name), cellOrDash(o.Description), codeCellOrDash(o.Value))
			}
			b.WriteString("\n")
		}
		if len(t.Call.Secrets) > 0 {
			b.WriteString("**Secrets:**\n\n")
			b.WriteString("| Name | Required | Description |\n")
			b.WriteString("|------|----------|-------------|\n")
			for _, s := range t.Call.Secrets {
				fmt.Fprintf(b, "| `%s` | %s | %s |\n",
					escapeCell(s.Name), yesNo(s.Required), cellOrDash(s.Description))
			}
			b.WriteString("\n")
		}
	}

	if len(t.Schedule) > 0 {
		b.WriteString("## Schedule\n\n")
		for _, c := range t.Schedule {
			fmt.Fprintf(b, "- `%s`", c.Cron)
			if c.Rationale != "" {
				fmt.Fprintf(b, " - %s", c.Rationale)
			}
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	if len(t.Events) > 0 {
		b.WriteString("## Event filters\n\n")
		for _, ev := range t.Events {
			fmt.Fprintf(b, "- **%s**\n", ev.Name)
			for _, f := range ev.Filters {
				fmt.Fprintf(b, "  - %s: %s\n", f.Key, codelist(f.Values))
			}
		}
		b.WriteString("\n")
	}
}

// writeInputTable writes a workflow_dispatch / workflow_call inputs table. For type:
// choice inputs, the options are appended inline to the description cell.
func writeInputTable(b *strings.Builder, inputs []model.WorkflowInput) {
	b.WriteString("| Name | Type | Required | Default | Description |\n")
	b.WriteString("|------|------|----------|---------|-------------|\n")
	for _, in := range inputs {
		desc := in.Description
		if len(in.Options) > 0 {
			opts := "Options: " + codelist(in.Options)
			if desc == "" {
				desc = opts
			} else {
				desc = desc + "<br>" + opts
			}
		}
		fmt.Fprintf(b, "| `%s` | %s | %s | %s | %s |\n",
			escapeCell(in.Name), cellOrDash(in.Type), yesNo(in.Required),
			codeCellOrDash(in.Default), cellOrDash(desc))
	}
	b.WriteString("\n")
}

// yesNo renders a boolean as Yes/No for a table cell.
func yesNo(v bool) string {
	if v {
		return "Yes"
	}
	return "No"
}

// cellOrDash escapes a value for a table cell, substituting "-" when empty.
func cellOrDash(s string) string {
	if s == "" {
		return "-"
	}
	return escapeCell(s)
}

// codeCellOrDash escapes a value and wraps it in code formatting, substituting "-" when
// empty.
func codeCellOrDash(s string) string {
	if s == "" {
		return "-"
	}
	return "`" + escapeCell(s) + "`"
}

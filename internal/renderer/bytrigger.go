package renderer

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/smol-utils/actiondoc/internal/callgraph"
)

// This file renders the by-trigger index: an inverted view that maps each trigger event to the
// workflows that fire on it, so a reader can answer "what runs on a pull_request?" at a glance
// without scanning every workflow's promoted Triggers line by hand. Like the document
// inventory, it is a read-only consumer of the parsed sources and the built call graph, and it
// lives outside RenderMarkdownGraph (which renders one source at a time and feeds the renderer
// goldens) so this cross-document aggregation stays off that single-source path.

// triggerIndexHeading deliberately differs from the per-workflow "**Triggers:**" promoted line
// so the two never collide on a shared GitHub anchor slug within the same rendered document.
const triggerIndexHeading = "## Workflows by trigger"

// triggerIndexEvent is one trigger event plus the workflows that fire on it, ordered
// deterministically.
type triggerIndexEvent struct {
	event  string
	usedBy []inventoryWorkflow
}

// RenderTriggerIndex builds the inverted trigger view: for every trigger event, the set of
// workflows that fire on it. The data source is model.Workflow.On -- the same canonical flat
// event list that drives the TOC trigger annotation and the per-workflow promoted Triggers
// line, so this index can never disagree with them. The workflow_call event is excluded:
// reusable workflows have their own TOC group, so a workflow_call row would only duplicate it
// (a workflow triggered by both push and workflow_call still appears under push). Each
// workflow is resolved to its display label and section anchor via the graph, with a filename
// suffix on names shared by more than one workflow so duplicate labels stay distinguishable.
// The whole section is suppressed (returns "") when no events remain after excluding
// workflow_call.
func RenderTriggerIndex(sources []callgraph.Source, g *callgraph.Graph) string {
	// event -> set of node ids (one per workflow), so a workflow that lists the same event
	// more than once is still attributed to that event exactly once.
	byEvent := map[string]map[string]bool{}
	for _, s := range sources {
		if s.Workflow == nil {
			continue
		}
		for _, event := range s.Workflow.On {
			if event == "workflow_call" {
				continue
			}
			if byEvent[event] == nil {
				byEvent[event] = map[string]bool{}
			}
			byEvent[event][s.Path] = true
		}
	}
	if len(byEvent) == 0 {
		return ""
	}

	events := buildTriggerIndexEvents(byEvent, g)

	var b strings.Builder
	b.WriteString(triggerIndexHeading + "\n\n")
	for _, e := range events {
		links := make([]string, len(e.usedBy))
		for i, wf := range e.usedBy {
			links[i] = fmt.Sprintf("[%s](#%s)", mdLinkLabel(wf.name), wf.anchor)
		}
		fmt.Fprintf(&b, "- **%s**: %s\n", e.event, strings.Join(links, ", "))
	}
	b.WriteString("\n")
	return b.String()
}

// buildTriggerIndexEvents turns an event -> set(node id) map into a sorted slice of
// triggerIndexEvent, resolving each node id to its display name and anchor via the graph.
// Events are ordered by descending workflow count, then alphabetically for ties (so the
// busiest triggers such as push and pull_request lead). Each event's workflow list is sorted
// by anchor (then display name) so rendering never depends on map iteration order.
func buildTriggerIndexEvents(byEvent map[string]map[string]bool, g *callgraph.Graph) []triggerIndexEvent {
	ambiguous := ambiguousWorkflowNames(g)
	out := make([]triggerIndexEvent, 0, len(byEvent))
	for event, ids := range byEvent {
		var wfs []inventoryWorkflow
		for id := range ids {
			wf := inventoryWorkflow{name: id, anchor: anchor(id)}
			if n := g.Nodes[id]; n != nil {
				// Disambiguate the visible label when another workflow shares this display
				// name: bare names would render identical link text against distinct anchors
				// (e.g. #model-jobs vs #model-jobs-1), so append the filename the way the
				// table of contents does.
				wf.name = n.Name
				if ambiguous[n.Name] {
					wf.name += " (" + filepath.Base(n.Path) + ")"
				}
				wf.anchor = nodeAnchor(n)
			}
			wfs = append(wfs, wf)
		}
		sort.Slice(wfs, func(i, j int) bool {
			if wfs[i].anchor != wfs[j].anchor {
				return wfs[i].anchor < wfs[j].anchor
			}
			return wfs[i].name < wfs[j].name
		})
		out = append(out, triggerIndexEvent{event: event, usedBy: wfs})
	}
	sort.Slice(out, func(i, j int) bool {
		if len(out[i].usedBy) != len(out[j].usedBy) {
			return len(out[i].usedBy) > len(out[j].usedBy)
		}
		return out[i].event < out[j].event
	})
	return out
}

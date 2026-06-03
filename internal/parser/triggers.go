package parser

import (
	"strings"

	"github.com/goccy/go-yaml/ast"
	"github.com/smol-utils/actiondoc/internal/model"
)

// parseTriggerSurface extracts the detailed trigger surface from an on: node. It returns
// nil when on: is a scalar or sequence (no per-event detail to surface) or when nothing
// of interest is declared. Workflow.On remains the canonical flat trigger list; this is
// purely additive.
func parseTriggerSurface(node ast.Node) *model.Triggers {
	mapping := toMapping(node)
	if mapping == nil {
		return nil
	}

	t := &model.Triggers{}
	for _, mv := range mapping.Values {
		switch mapKeyString(mv.Key) {
		case "workflow_dispatch":
			if inputs := parseWorkflowInputs(childMapping(mv.Value, "inputs")); len(inputs) > 0 {
				t.Dispatch = &model.DispatchTrigger{Inputs: inputs}
			}
		case "workflow_call":
			if t.Call = parseCallTrigger(mv.Value); t.Call == nil {
				// A bare `workflow_call:` with no inputs/outputs/secrets is still a
				// reusable workflow; record it so the "reusable via workflow_call" note
				// renders.
				t.Call = &model.CallTrigger{}
			}
		case "schedule":
			t.Schedule = parseSchedule(mv.Value)
		default:
			if ev := parseGenericEvent(mv); ev != nil {
				t.Events = append(t.Events, *ev)
			}
		}
	}

	if t.Dispatch == nil && t.Call == nil && len(t.Schedule) == 0 && len(t.Events) == 0 {
		return nil
	}
	return t
}

// childMapping returns the mapping under the named key within a parent node, or nil.
func childMapping(parent ast.Node, key string) *ast.MappingNode {
	mapping := toMapping(parent)
	if mapping == nil {
		return nil
	}
	for _, mv := range mapping.Values {
		if mapKeyString(mv.Key) == key {
			return toMapping(mv.Value)
		}
	}
	return nil
}

// parseCallTrigger parses a workflow_call block into its public API (inputs, outputs,
// secrets). Returns nil when none are declared.
func parseCallTrigger(node ast.Node) *model.CallTrigger {
	mapping := toMapping(node)
	if mapping == nil {
		return nil
	}
	c := &model.CallTrigger{}
	for _, mv := range mapping.Values {
		switch mapKeyString(mv.Key) {
		case "inputs":
			c.Inputs = parseWorkflowInputs(toMapping(mv.Value))
		case "outputs":
			c.Outputs = parseCallOutputs(toMapping(mv.Value))
		case "secrets":
			c.Secrets = parseCallSecrets(toMapping(mv.Value))
		}
	}
	if len(c.Inputs) == 0 && len(c.Outputs) == 0 && len(c.Secrets) == 0 {
		return nil
	}
	return c
}

// parseWorkflowInputs parses an inputs: mapping (shared by workflow_dispatch and
// workflow_call). Names preserve source casing; type: choice options are collected.
func parseWorkflowInputs(mapping *ast.MappingNode) []model.WorkflowInput {
	if mapping == nil {
		return nil
	}
	var inputs []model.WorkflowInput
	for _, mv := range mapping.Values {
		in := model.WorkflowInput{Name: mapKeyString(mv.Key)}
		spec := toMapping(mv.Value)
		if spec != nil {
			for _, f := range spec.Values {
				switch mapKeyString(f.Key) {
				case "type":
					in.Type = nodeString(f.Value)
				case "required":
					in.Required = isTrue(f.Value)
				case "default":
					in.Default = nodeString(f.Value)
				case "description":
					in.Description = strings.TrimSpace(nodeString(f.Value))
				case "options":
					in.Options = parseStringOrSequence(f.Value)
				}
			}
		}
		inputs = append(inputs, in)
	}
	return inputs
}

// parseCallOutputs parses a workflow_call outputs: mapping.
func parseCallOutputs(mapping *ast.MappingNode) []model.WorkflowOutput {
	if mapping == nil {
		return nil
	}
	var outputs []model.WorkflowOutput
	for _, mv := range mapping.Values {
		out := model.WorkflowOutput{Name: mapKeyString(mv.Key)}
		if spec := toMapping(mv.Value); spec != nil {
			for _, f := range spec.Values {
				switch mapKeyString(f.Key) {
				case "description":
					out.Description = strings.TrimSpace(nodeString(f.Value))
				case "value":
					out.Value = nodeString(f.Value)
				}
			}
		}
		outputs = append(outputs, out)
	}
	return outputs
}

// parseCallSecrets parses a workflow_call secrets: mapping. Names preserve source casing
// and hyphenation (e.g. CONSUMER-KEY).
func parseCallSecrets(mapping *ast.MappingNode) []model.WorkflowSecret {
	if mapping == nil {
		return nil
	}
	var secrets []model.WorkflowSecret
	for _, mv := range mapping.Values {
		s := model.WorkflowSecret{Name: mapKeyString(mv.Key)}
		if spec := toMapping(mv.Value); spec != nil {
			for _, f := range spec.Values {
				switch mapKeyString(f.Key) {
				case "required":
					s.Required = isTrue(f.Value)
				case "description":
					s.Description = strings.TrimSpace(nodeString(f.Value))
				}
			}
		}
		secrets = append(secrets, s)
	}
	return secrets
}

// parseSchedule parses an on.schedule sequence of `- cron: <expr>` entries, preserving
// each cron string verbatim and attaching any trailing-comment rationale.
func parseSchedule(node ast.Node) []model.CronEntry {
	seq, ok := node.(*ast.SequenceNode)
	if !ok {
		return nil
	}
	var entries []model.CronEntry
	for i, item := range seq.Values {
		mapping := toMapping(item)
		if mapping == nil {
			continue
		}
		for _, mv := range mapping.Values {
			if mapKeyString(mv.Key) != "cron" {
				continue
			}
			entry := model.CronEntry{Cron: nodeString(mv.Value)}
			entry.Rationale = TrailingComment(mv)
			if entry.Rationale == "" && i < len(seq.ValueHeadComments) {
				entry.Rationale = cleanCommentText(commentString(seq.ValueHeadComments[i]))
			}
			entries = append(entries, entry)
		}
	}
	return entries
}

// parseGenericEvent models any non-special event (push, pull_request, release,
// repository_dispatch, workflow_run, ...) as an ordered list of its filter sub-keys. It
// returns nil when the event has no sub-keys (those events are already in Workflow.On).
func parseGenericEvent(mv *ast.MappingValueNode) *model.TriggerEvent {
	mapping := toMapping(mv.Value)
	if mapping == nil {
		return nil
	}
	ev := &model.TriggerEvent{Name: mapKeyString(mv.Key)}
	for _, f := range mapping.Values {
		values := parseStringOrSequence(f.Value)
		if values == nil {
			if s := nodeString(f.Value); s != "" {
				values = []string{s}
			}
		}
		ev.Filters = append(ev.Filters, model.TriggerFilter{
			Key:    mapKeyString(f.Key),
			Values: values,
		})
	}
	if len(ev.Filters) == 0 {
		return nil
	}
	return ev
}

// isTrue reports whether a YAML node represents a boolean true.
func isTrue(node ast.Node) bool {
	return strings.EqualFold(strings.TrimSpace(nodeString(node)), "true")
}

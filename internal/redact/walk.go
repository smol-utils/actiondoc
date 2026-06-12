package redact

import (
	"strings"

	"github.com/smol-utils/actiondoc/internal/callgraph"
	"github.com/smol-utils/actiondoc/internal/model"
)

// This file is the single enumeration of the model's redactable fields. Both passes --
// collection and rewrite -- run the same walk with a different fieldOps, so a field
// wired in here is covered by both passes or by neither; they cannot drift apart. The
// completeness test plants a sentinel in every model string field and fails naming any
// field this walk misses.

// fieldOps is what a pass does to each kind of redactable string. Every method returns
// the value to store back: the collect pass records originals and returns its input
// unchanged, the rewrite pass returns the substituted value.
type fieldOps interface {
	// free handles free text and expression-bearing strings (descriptions, run scripts,
	// conditions): identifiers anywhere, hosts/URLs in literal spans.
	free(s string) string
	// value handles env:/with: values, which the aggressive profile replaces wholesale
	// when they are pure literals.
	value(s string) string
	// name handles a whole-string category name (secret/env/environment names).
	name(cat, s string) string
	// runnerLabel handles a single non-expression runs-on label.
	runnerLabel(s string) string
}

// walkSource applies ops to every redactable field of a parsed source.
func walkSource(s callgraph.Source, ops fieldOps) {
	if s.Workflow != nil {
		walkWorkflow(s.Workflow, ops)
	} else if s.Action != nil {
		walkAction(s.Action, ops)
	}
}

func walkWorkflow(w *model.Workflow, ops fieldOps) {
	walkTags(&w.Tags, ops)
	w.Description = ops.free(w.Description)
	walkKVs(w.Env, catEnv, ops)
	walkConcurrency(w.Concurrency, ops)
	walkDefaults(w.Defaults, ops)
	walkPermissions(w.Permissions, ops)
	walkTriggers(w.Triggers, ops)
	for ji := range w.Jobs {
		walkJob(&w.Jobs[ji], ops)
	}
}

func walkJob(j *model.Job, ops fieldOps) {
	walkTags(&j.Tags, ops)
	j.Description = ops.free(j.Description)
	j.RunsOn = walkRunsOn(j.RunsOn, ops)
	j.If = ops.free(j.If)
	walkKVs(j.Env, catEnv, ops)
	walkKVs(j.With, "", ops)
	for i := range j.Secrets {
		j.Secrets[i].Key = ops.name(catSecret, j.Secrets[i].Key)
		j.Secrets[i].Value = ops.free(j.Secrets[i].Value)
	}
	walkConcurrency(j.Concurrency, ops)
	walkDefaults(j.Defaults, ops)
	walkPermissions(j.Permissions, ops)
	if j.Environment != nil {
		// A literal deploy-environment name is redacted whole; an expression as free text.
		if strings.Contains(j.Environment.Name, "${{") {
			j.Environment.Name = ops.free(j.Environment.Name)
		} else {
			j.Environment.Name = ops.name(catEnvironment, j.Environment.Name)
		}
		j.Environment.URL = ops.free(j.Environment.URL)
	}
	for ai := range j.Matrix {
		for vi := range j.Matrix[ai].Values {
			j.Matrix[ai].Values[vi] = ops.free(j.Matrix[ai].Values[vi])
		}
	}
	for si := range j.Steps {
		walkStep(&j.Steps[si], ops)
	}
}

func walkStep(s *model.Step, ops fieldOps) {
	walkTags(&s.Tags, ops)
	s.Description = ops.free(s.Description)
	s.Run = ops.free(s.Run)
	s.If = ops.free(s.If)
	s.ContinueOnErrorExpr = ops.free(s.ContinueOnErrorExpr)
	walkKVs(s.With, "", ops)
	walkKVs(s.Env, catEnv, ops)
}

func walkAction(a *model.Action, ops fieldOps) {
	walkTags(&a.Tags, ops)
	a.Description = ops.free(a.Description)
	for i := range a.Inputs {
		a.Inputs[i].Description = ops.free(a.Inputs[i].Description)
		// An action's input default is part of its public contract, like a workflow_call
		// input default -- treat it as free text (redact secrets/hosts/URLs inside it) but
		// never blank it wholesale, even under the aggressive profile, whose value-blanking
		// is scoped to env:/with: values.
		a.Inputs[i].Default = ops.free(a.Inputs[i].Default)
	}
	for i := range a.Outputs {
		a.Outputs[i].Description = ops.free(a.Outputs[i].Description)
	}
	a.Runs.Image = ops.free(a.Runs.Image)
}

func walkTags(t *model.Tags, ops fieldOps) {
	t.Desc = ops.free(t.Desc)
	t.Deprecated = ops.free(t.Deprecated)
	t.Example = ops.free(t.Example)
	for i := range t.Secrets {
		t.Secrets[i].Name = ops.name(catSecret, t.Secrets[i].Name)
		t.Secrets[i].Description = ops.free(t.Secrets[i].Description)
	}
	for i := range t.Envs {
		t.Envs[i].Name = ops.name(catEnv, t.Envs[i].Name)
		t.Envs[i].Description = ops.free(t.Envs[i].Description)
	}
	for i := range t.Inputs {
		t.Inputs[i].Description = ops.free(t.Inputs[i].Description)
	}
	for i := range t.Outputs {
		t.Outputs[i].Description = ops.free(t.Outputs[i].Description)
	}
	for i := range t.See {
		t.See[i] = ops.free(t.See[i])
	}
}

// walkKVs sweeps an ordered key/value list. keyCat is the category the keys are
// redacted under, or "" when the keys are public contract (with: input names).
func walkKVs(kvs []model.KV, keyCat string, ops fieldOps) {
	for i := range kvs {
		if keyCat != "" {
			kvs[i].Key = ops.name(keyCat, kvs[i].Key)
		}
		kvs[i].Value = ops.value(kvs[i].Value)
	}
}

// walkRunsOn rebuilds a runs-on string, preserving the comma-joined list shape and the
// map form (group: ..., labels: ...): expressions go through free, labels (and map-form
// values) through runnerLabel.
func walkRunsOn(s string, ops fieldOps) string {
	if s == "" {
		return s
	}
	parts := strings.Split(s, ", ")
	for i, p := range parts {
		switch {
		case strings.Contains(p, "${{"):
			parts[i] = ops.free(p)
		case strings.Contains(p, ": "):
			kv := strings.SplitN(p, ": ", 2)
			parts[i] = kv[0] + ": " + ops.runnerLabel(kv[1])
		default:
			parts[i] = ops.runnerLabel(p)
		}
	}
	return strings.Join(parts, ", ")
}

// walkConcurrency sweeps a concurrency: block. Group is usually an expression and
// CancelInProgress is a raw value that may also be one; both are free text.
func walkConcurrency(c *model.Concurrency, ops fieldOps) {
	if c == nil {
		return
	}
	c.Group = ops.free(c.Group)
	c.CancelInProgress = ops.free(c.CancelInProgress)
}

// walkDefaults sweeps a defaults.run: block. The shell may be a custom command template
// and the working directory an arbitrary path; both are free text so a host or URL
// inside them is still caught.
func walkDefaults(d *model.Defaults, ops fieldOps) {
	if d == nil {
		return
	}
	d.Shell = ops.free(d.Shell)
	d.WorkingDirectory = ops.free(d.WorkingDirectory)
}

// walkPermissions sweeps the optional rationale comments on a permissions: block
// (e.g. `contents: read  # for the internal mirror`), which are free text that can
// carry a hostname. The scopes and levels themselves are not sensitive.
func walkPermissions(p *model.Permissions, ops fieldOps) {
	if p == nil {
		return
	}
	for i := range p.Scopes {
		p.Scopes[i].Rationale = ops.free(p.Scopes[i].Rationale)
	}
}

func walkTriggers(t *model.Triggers, ops fieldOps) {
	if t == nil {
		return
	}
	for i := range t.Schedule {
		t.Schedule[i].Rationale = ops.free(t.Schedule[i].Rationale)
	}
	if t.Call != nil {
		for i := range t.Call.Secrets {
			t.Call.Secrets[i].Name = ops.name(catSecret, t.Call.Secrets[i].Name)
			t.Call.Secrets[i].Description = ops.free(t.Call.Secrets[i].Description)
		}
		walkWorkflowInputs(t.Call.Inputs, ops)
		for i := range t.Call.Outputs {
			t.Call.Outputs[i].Value = ops.free(t.Call.Outputs[i].Value)
			t.Call.Outputs[i].Description = ops.free(t.Call.Outputs[i].Description)
		}
	}
	if t.Dispatch != nil {
		walkWorkflowInputs(t.Dispatch.Inputs, ops)
	}
	// Event filter values (branches, paths, tags, workflow names, dispatch types) are
	// free-form strings; sweep them so a host or URL in a filter is still caught.
	for ei := range t.Events {
		for fi := range t.Events[ei].Filters {
			f := &t.Events[ei].Filters[fi]
			for vi := range f.Values {
				f.Values[vi] = ops.free(f.Values[vi])
			}
		}
	}
}

// walkWorkflowInputs sweeps the free-text parts of declared workflow_dispatch /
// workflow_call inputs: defaults, descriptions, and choice options. Input names and
// types are public contract and stay readable.
func walkWorkflowInputs(ins []model.WorkflowInput, ops fieldOps) {
	for i := range ins {
		ins[i].Default = ops.free(ins[i].Default)
		ins[i].Description = ops.free(ins[i].Description)
		for oi := range ins[i].Options {
			ins[i].Options[oi] = ops.free(ins[i].Options[oi])
		}
	}
}

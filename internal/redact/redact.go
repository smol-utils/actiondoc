// Package redact rewrites a parsed workflow/action model in place so its generated
// documentation can be shared outside the owning team or handed to an external service
// (for example, an LLM asked to write prose docs) without leaking sensitive-but-
// unclassified material: secret and variable names, environment-variable names and
// values, hostnames and URLs, GitHub deploy-environment names, and self-hosted runner
// labels.
//
// Redaction is a single transform over the intermediate model, applied after parsing and
// call-graph construction but before rendering. Because every rendered surface -- the
// secret/variable inventories, the call-graph transitive requirements, and the per-step
// reference tables -- is derived from this model, redacting it once keeps every
// cross-reference internally consistent for free, in both Markdown and JSON output.
//
// Redaction is consistent pseudonymization, not blanking: the same secret maps to the
// same placeholder (SECRET_1) everywhere it appears. Placeholder numbering is
// deterministic -- originals are sorted before they are numbered -- so output diffs and
// snapshots stay reviewable across runs.
package redact

import (
	"encoding/json"
	"sort"
	"strconv"
	"strings"

	"github.com/smol-utils/actiondoc/internal/callgraph"
	"github.com/smol-utils/actiondoc/internal/model"
)

// Level selects how aggressively literal values in env:/with: blocks are redacted.
// Identifiers (secrets, vars, env names), hosts, URLs, runner labels, and environment
// names are redacted at every level.
type Level int

const (
	// Conservative redacts known-sensitive identifiers plus any host or URL detected in
	// values and free text, but leaves harmless literals (e.g. node-version: 20) readable.
	Conservative Level = iota
	// Aggressive additionally replaces every pure-literal value in an env:/with: block
	// with a VALUE_n placeholder, on the assumption any literal may be sensitive.
	Aggressive
)

// Options configures a redaction pass.
type Options struct {
	Level Level
}

// category names. The identifier categories ("secrets", "vars", "env") double as the
// expression-context keys consumed by rewriteIdents, so they must match contexts in
// scan.go.
const (
	catSecret      = "secrets"
	catVar         = "vars"
	catEnv         = "env"
	catRunner      = "runner"
	catEnvironment = "environment"
	catHost        = "host"
	catURL         = "url"
	catValue       = "value"
)

// prefixes maps each category to its placeholder prefix.
var prefixes = map[string]string{
	catSecret:      "SECRET",
	catVar:         "VAR",
	catEnv:         "ENV",
	catRunner:      "RUNNER",
	catEnvironment: "ENVIRONMENT",
	catHost:        "HOST",
	catURL:         "URL",
	catValue:       "VALUE",
}

// Mapping is the reverse map from placeholder to original. It exists only so a user can
// substitute real names back into an external service's prose output afterward; it is
// never part of the shareable documentation and is written only to an explicit local file.
type Mapping struct {
	entries map[string]string
}

// JSON renders the mapping as indented JSON keyed by placeholder. Returns nil when no
// substitutions were made, so the caller can skip writing an empty file.
func (m *Mapping) JSON() ([]byte, error) {
	if m == nil || len(m.entries) == 0 {
		return nil, nil
	}
	return json.MarshalIndent(m.entries, "", "  ")
}

// Empty reports whether the pass produced no substitutions.
func (m *Mapping) Empty() bool {
	return m == nil || len(m.entries) == 0
}

// redactor holds the placeholder maps built during collection and consumed during
// rewrite.
type redactor struct {
	opts Options
	// m[category][original] = placeholder. The identifier categories are passed directly
	// to rewriteIdents/rewriteHostsURLs.
	m map[string]map[string]string
	// seen[category][original] accumulates originals during collection, before numbering.
	seen map[string]map[string]bool
}

// Apply redacts every source in place and returns the reverse mapping. The call graph
// shares the same model pointers, so rendering it afterward sees the redacted values; the
// graph's denormalized names and `uses:` refs are intentionally left intact (titles and
// call-graph shape are not redacted).
func Apply(sources []callgraph.Source, opts Options) *Mapping {
	r := &redactor{
		opts: opts,
		m:    map[string]map[string]string{},
		seen: map[string]map[string]bool{},
	}
	for _, s := range sources {
		if s.Workflow != nil {
			r.collectWorkflow(s.Workflow)
		} else if s.Action != nil {
			r.collectAction(s.Action)
		}
	}
	r.assign()
	for _, s := range sources {
		if s.Workflow != nil {
			r.rewriteWorkflow(s.Workflow)
		} else if s.Action != nil {
			r.rewriteAction(s.Action)
		}
	}
	return r.mapping()
}

// note records an original under a category for later numbering.
func (r *redactor) note(cat, original string) {
	if original == "" {
		return
	}
	if cat == catSecret && keepSecrets[original] {
		return
	}
	if r.seen[cat] == nil {
		r.seen[cat] = map[string]bool{}
	}
	r.seen[cat][original] = true
}

// assign turns the collected originals into deterministic placeholders: each category's
// originals are sorted, then numbered from 1.
func (r *redactor) assign() {
	for cat, set := range r.seen {
		originals := make([]string, 0, len(set))
		for o := range set {
			originals = append(originals, o)
		}
		sort.Strings(originals)
		r.m[cat] = make(map[string]string, len(originals))
		for i, o := range originals {
			r.m[cat][o] = prefixes[cat] + "_" + strconv.Itoa(i+1)
		}
	}
}

// mapping builds the placeholder->original reverse map from every category.
func (r *redactor) mapping() *Mapping {
	entries := map[string]string{}
	for _, byOrig := range r.m {
		for orig, ph := range byOrig {
			entries[ph] = orig
		}
	}
	return &Mapping{entries: entries}
}

// lookup returns the placeholder for an original in a category, or the original unchanged.
func (r *redactor) lookup(cat, original string) string {
	if to, ok := r.m[cat][original]; ok {
		return to
	}
	return original
}

// --- shared field redactors -------------------------------------------------

// noteFree collects identifiers, hosts, and URLs from a free-text or expression-bearing
// string (run scripts, conditions, descriptions, generic values). Identifiers are scanned
// across the whole string; hosts and URLs only in literal spans, so expression-context
// refs (secrets.X, github.sha) are never mistaken for hostnames.
func (r *redactor) noteFree(s string) {
	walkIdents(s, func(ctx, name string) { r.note(ctx, name) })
	eachLiteralSpan(s, func(span string) {
		walkHostsURLs(span, func(u string) { r.note(catURL, u) }, func(h string) { r.note(catHost, h) })
	})
}

// redactFree applies identifier and host/URL substitution to a free-text string.
// Host/URL substitution runs first, on literal spans only: a literal hostname whose
// leading label looks like an expression context (secrets.example.com) would otherwise be
// corrupted by identifier rewriting (-> secrets.SECRET_1.com) before its host placeholder
// could land, leaving the host partly in the clear. Identifier rewriting then runs over
// the whole string, including expression bodies; the host/URL placeholders it now sees
// contain no `context.` pattern, so they pass through untouched.
func (r *redactor) redactFree(s string) string {
	s = mapLiteralSpans(s, func(span string) string {
		return rewriteHostsURLs(span, r.m[catURL], r.m[catHost])
	})
	s = rewriteIdents(s, r.m)
	return s
}

// noteValue collects from an env:/with: value. Under Aggressive, a pure literal (no
// expression) is recorded as a whole-value placeholder; otherwise it is treated as free
// text so its identifiers and hosts are still found.
func (r *redactor) noteValue(s string) {
	if r.opts.Level == Aggressive && s != "" && !strings.Contains(s, "${{") {
		r.note(catValue, s)
		return
	}
	r.noteFree(s)
}

// redactValue mirrors noteValue for the rewrite pass.
func (r *redactor) redactValue(s string) string {
	if r.opts.Level == Aggressive && s != "" && !strings.Contains(s, "${{") {
		return r.lookup(catValue, s)
	}
	return r.redactFree(s)
}

// noteRunsOn collects self-hosted runner labels (anything outside the well-known set).
func (r *redactor) noteRunsOn(s string) {
	r.eachRunnerLabel(s, func(label string) { r.note(catRunner, label) })
}

// redactRunsOn replaces self-hosted runner labels while preserving well-known labels,
// the comma-joined list shape, and any matrix expression.
func (r *redactor) redactRunsOn(s string) string {
	if s == "" {
		return s
	}
	parts := strings.Split(s, ", ")
	for i, p := range parts {
		switch {
		case strings.Contains(p, "${{"):
			parts[i] = r.redactFree(p)
		case strings.Contains(p, ": "):
			// runs-on map form (group: ..., labels: ...): redact only the value side.
			kv := strings.SplitN(p, ": ", 2)
			parts[i] = kv[0] + ": " + r.classifyRunner(kv[1])
		default:
			parts[i] = r.classifyRunner(p)
		}
	}
	return strings.Join(parts, ", ")
}

// eachRunnerLabel calls fn for every label in a runs-on string that should be redacted.
func (r *redactor) eachRunnerLabel(s string, fn func(string)) {
	if s == "" {
		return
	}
	for _, p := range strings.Split(s, ", ") {
		switch {
		case strings.Contains(p, "${{"):
			r.noteFree(p)
		case strings.Contains(p, ": "):
			kv := strings.SplitN(p, ": ", 2)
			if !keepRunners[strings.ToLower(kv[1])] {
				fn(kv[1])
			}
		default:
			if !keepRunners[strings.ToLower(p)] {
				fn(p)
			}
		}
	}
}

// classifyRunner returns the placeholder for a non-well-known runner label, or the label
// unchanged when it is a GitHub-hosted runner or generic keyword.
func (r *redactor) classifyRunner(label string) string {
	if keepRunners[strings.ToLower(label)] {
		return label
	}
	return r.lookup(catRunner, label)
}

// noteEnvironmentName collects a job's deploy-environment name (a literal name only; an
// expression is handled as free text).
func (r *redactor) noteEnvironmentName(s string) {
	if s == "" {
		return
	}
	if strings.Contains(s, "${{") {
		r.noteFree(s)
		return
	}
	r.note(catEnvironment, s)
}

// redactEnvironmentName mirrors noteEnvironmentName.
func (r *redactor) redactEnvironmentName(s string) string {
	if s == "" {
		return s
	}
	if strings.Contains(s, "${{") {
		return r.redactFree(s)
	}
	return r.lookup(catEnvironment, s)
}

// --- workflow walk ----------------------------------------------------------

func (r *redactor) collectWorkflow(w *model.Workflow) {
	r.collectTags(w.Tags)
	r.noteFree(w.Description)
	for _, kv := range w.Env {
		r.note(catEnv, kv.Key)
		r.noteValue(kv.Value)
	}
	if w.Concurrency != nil {
		r.noteFree(w.Concurrency.Group)
	}
	r.collectTriggers(w.Triggers)
	for ji := range w.Jobs {
		r.collectJob(&w.Jobs[ji])
	}
}

func (r *redactor) rewriteWorkflow(w *model.Workflow) {
	r.rewriteTags(&w.Tags)
	w.Description = r.redactFree(w.Description)
	for i := range w.Env {
		w.Env[i].Key = r.lookup(catEnv, w.Env[i].Key)
		w.Env[i].Value = r.redactValue(w.Env[i].Value)
	}
	if w.Concurrency != nil {
		w.Concurrency.Group = r.redactFree(w.Concurrency.Group)
	}
	r.rewriteTriggers(w.Triggers)
	for ji := range w.Jobs {
		r.rewriteJob(&w.Jobs[ji])
	}
}

func (r *redactor) collectJob(j *model.Job) {
	r.collectTags(j.Tags)
	r.noteFree(j.Description)
	r.noteRunsOn(j.RunsOn)
	r.noteFree(j.If)
	for _, kv := range j.Env {
		r.note(catEnv, kv.Key)
		r.noteValue(kv.Value)
	}
	for _, kv := range j.With {
		r.noteValue(kv.Value)
	}
	for _, kv := range j.Secrets {
		r.note(catSecret, kv.Key)
		r.noteFree(kv.Value)
	}
	if j.Concurrency != nil {
		r.noteFree(j.Concurrency.Group)
	}
	if j.Environment != nil {
		r.noteEnvironmentName(j.Environment.Name)
		r.noteFree(j.Environment.URL)
	}
	for _, ax := range j.Matrix {
		for _, v := range ax.Values {
			r.noteFree(v)
		}
	}
	for si := range j.Steps {
		r.collectStep(&j.Steps[si])
	}
}

func (r *redactor) rewriteJob(j *model.Job) {
	r.rewriteTags(&j.Tags)
	j.Description = r.redactFree(j.Description)
	j.RunsOn = r.redactRunsOn(j.RunsOn)
	j.If = r.redactFree(j.If)
	for i := range j.Env {
		j.Env[i].Key = r.lookup(catEnv, j.Env[i].Key)
		j.Env[i].Value = r.redactValue(j.Env[i].Value)
	}
	for i := range j.With {
		j.With[i].Value = r.redactValue(j.With[i].Value)
	}
	for i := range j.Secrets {
		j.Secrets[i].Key = r.lookup(catSecret, j.Secrets[i].Key)
		j.Secrets[i].Value = r.redactFree(j.Secrets[i].Value)
	}
	if j.Concurrency != nil {
		j.Concurrency.Group = r.redactFree(j.Concurrency.Group)
	}
	if j.Environment != nil {
		j.Environment.Name = r.redactEnvironmentName(j.Environment.Name)
		j.Environment.URL = r.redactFree(j.Environment.URL)
	}
	for ai := range j.Matrix {
		for vi := range j.Matrix[ai].Values {
			j.Matrix[ai].Values[vi] = r.redactFree(j.Matrix[ai].Values[vi])
		}
	}
	for si := range j.Steps {
		r.rewriteStep(&j.Steps[si])
	}
}

func (r *redactor) collectStep(s *model.Step) {
	r.collectTags(s.Tags)
	r.noteFree(s.Description)
	r.noteFree(s.Run)
	r.noteFree(s.If)
	for _, kv := range s.With {
		r.noteValue(kv.Value)
	}
	for _, kv := range s.Env {
		r.note(catEnv, kv.Key)
		r.noteValue(kv.Value)
	}
}

func (r *redactor) rewriteStep(s *model.Step) {
	r.rewriteTags(&s.Tags)
	s.Description = r.redactFree(s.Description)
	s.Run = r.redactFree(s.Run)
	s.If = r.redactFree(s.If)
	for i := range s.With {
		s.With[i].Value = r.redactValue(s.With[i].Value)
	}
	for i := range s.Env {
		s.Env[i].Key = r.lookup(catEnv, s.Env[i].Key)
		s.Env[i].Value = r.redactValue(s.Env[i].Value)
	}
}

func (r *redactor) collectTriggers(t *model.Triggers) {
	if t == nil {
		return
	}
	if t.Call != nil {
		for _, sec := range t.Call.Secrets {
			r.note(catSecret, sec.Name)
			r.noteFree(sec.Description)
		}
		for _, in := range t.Call.Inputs {
			r.noteFree(in.Default)
			r.noteFree(in.Description)
		}
		for _, out := range t.Call.Outputs {
			r.noteFree(out.Value)
			r.noteFree(out.Description)
		}
	}
	if t.Dispatch != nil {
		for _, in := range t.Dispatch.Inputs {
			r.noteFree(in.Default)
			r.noteFree(in.Description)
		}
	}
}

func (r *redactor) rewriteTriggers(t *model.Triggers) {
	if t == nil {
		return
	}
	if t.Call != nil {
		for i := range t.Call.Secrets {
			t.Call.Secrets[i].Name = r.lookup(catSecret, t.Call.Secrets[i].Name)
			t.Call.Secrets[i].Description = r.redactFree(t.Call.Secrets[i].Description)
		}
		for i := range t.Call.Inputs {
			t.Call.Inputs[i].Default = r.redactFree(t.Call.Inputs[i].Default)
			t.Call.Inputs[i].Description = r.redactFree(t.Call.Inputs[i].Description)
		}
		for i := range t.Call.Outputs {
			t.Call.Outputs[i].Value = r.redactFree(t.Call.Outputs[i].Value)
			t.Call.Outputs[i].Description = r.redactFree(t.Call.Outputs[i].Description)
		}
	}
	if t.Dispatch != nil {
		for i := range t.Dispatch.Inputs {
			t.Dispatch.Inputs[i].Default = r.redactFree(t.Dispatch.Inputs[i].Default)
			t.Dispatch.Inputs[i].Description = r.redactFree(t.Dispatch.Inputs[i].Description)
		}
	}
}

// --- action walk ------------------------------------------------------------

func (r *redactor) collectAction(a *model.Action) {
	r.collectTags(a.Tags)
	r.noteFree(a.Description)
	for _, in := range a.Inputs {
		r.noteFree(in.Description)
		// An action's input default is part of its public contract, like a workflow_call
		// input default -- treat it as free text (redact secrets/hosts/URLs inside it) but
		// never blank it wholesale, even under the aggressive profile, whose value-blanking
		// is scoped to env:/with: values.
		r.noteFree(in.Default)
	}
	for _, out := range a.Outputs {
		r.noteFree(out.Description)
	}
	r.noteFree(a.Runs.Image)
}

func (r *redactor) rewriteAction(a *model.Action) {
	r.rewriteTags(&a.Tags)
	a.Description = r.redactFree(a.Description)
	for i := range a.Inputs {
		a.Inputs[i].Description = r.redactFree(a.Inputs[i].Description)
		a.Inputs[i].Default = r.redactFree(a.Inputs[i].Default)
	}
	for i := range a.Outputs {
		a.Outputs[i].Description = r.redactFree(a.Outputs[i].Description)
	}
	a.Runs.Image = r.redactFree(a.Runs.Image)
}

// --- tags -------------------------------------------------------------------

func (r *redactor) collectTags(t model.Tags) {
	r.noteFree(t.Desc)
	r.noteFree(t.Deprecated)
	r.noteFree(t.Example)
	for _, p := range t.Secrets {
		r.note(catSecret, p.Name)
		r.noteFree(p.Description)
	}
	for _, p := range t.Envs {
		r.note(catEnv, p.Name)
		r.noteFree(p.Description)
	}
	for _, p := range t.Inputs {
		r.noteFree(p.Description)
	}
	for _, p := range t.Outputs {
		r.noteFree(p.Description)
	}
	for _, s := range t.See {
		r.noteFree(s)
	}
}

func (r *redactor) rewriteTags(t *model.Tags) {
	t.Desc = r.redactFree(t.Desc)
	t.Deprecated = r.redactFree(t.Deprecated)
	t.Example = r.redactFree(t.Example)
	for i := range t.Secrets {
		t.Secrets[i].Name = r.lookup(catSecret, t.Secrets[i].Name)
		t.Secrets[i].Description = r.redactFree(t.Secrets[i].Description)
	}
	for i := range t.Envs {
		t.Envs[i].Name = r.lookup(catEnv, t.Envs[i].Name)
		t.Envs[i].Description = r.redactFree(t.Envs[i].Description)
	}
	for i := range t.Inputs {
		t.Inputs[i].Description = r.redactFree(t.Inputs[i].Description)
	}
	for i := range t.Outputs {
		t.Outputs[i].Description = r.redactFree(t.Outputs[i].Description)
	}
	for i := range t.See {
		t.See[i] = r.redactFree(t.See[i])
	}
}

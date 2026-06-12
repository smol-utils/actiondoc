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
//
// Both passes run the same field walk (walk.go); only the string operations differ, so
// collection and rewrite cannot disagree about which fields are covered.
func Apply(sources []callgraph.Source, opts Options) *Mapping {
	r := &redactor{
		opts: opts,
		m:    map[string]map[string]string{},
		seen: map[string]map[string]bool{},
	}
	for _, s := range sources {
		walkSource(s, collector{r})
	}
	r.assign()
	for _, s := range sources {
		walkSource(s, rewriter{r})
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

// --- the two passes ----------------------------------------------------------

// collector implements fieldOps for the first pass: it records every redactable
// original and returns its inputs unchanged.
type collector struct{ r *redactor }

// free collects identifiers, hosts, and URLs from a free-text or expression-bearing
// string. Identifiers are scanned across the whole string; hosts and URLs only in
// literal spans, so expression-context refs (secrets.X, github.sha) are never mistaken
// for hostnames.
func (c collector) free(s string) string {
	walkIdents(s, func(ctx, name string) { c.r.note(ctx, name) })
	eachLiteralSpan(s, func(span string) {
		walkHostsURLs(span, func(u string) { c.r.note(catURL, u) }, func(h string) { c.r.note(catHost, h) })
	})
	return s
}

// value records a pure-literal env:/with: value as a whole-value placeholder under
// Aggressive; otherwise the value is free text.
func (c collector) value(s string) string {
	if c.r.opts.Level == Aggressive && s != "" && !strings.Contains(s, "${{") {
		c.r.note(catValue, s)
		return s
	}
	return c.free(s)
}

func (c collector) name(cat, s string) string {
	c.r.note(cat, s)
	return s
}

func (c collector) runnerLabel(s string) string {
	if !keepRunners[strings.ToLower(s)] {
		c.r.note(catRunner, s)
	}
	return s
}

// rewriter implements fieldOps for the second pass: it substitutes the placeholders
// assigned after collection.
type rewriter struct{ r *redactor }

// free applies identifier and host/URL substitution to a free-text string.
// Host/URL substitution runs first, on literal spans only: a literal hostname whose
// leading label looks like an expression context (secrets.example.com) would otherwise be
// corrupted by identifier rewriting (-> secrets.SECRET_1.com) before its host placeholder
// could land, leaving the host partly in the clear. Identifier rewriting then runs over
// the whole string, including expression bodies; the host/URL placeholders it now sees
// contain no `context.` pattern, so they pass through untouched.
func (w rewriter) free(s string) string {
	s = mapLiteralSpans(s, func(span string) string {
		return rewriteHostsURLs(span, w.r.m[catURL], w.r.m[catHost])
	})
	return rewriteIdents(s, w.r.m)
}

func (w rewriter) value(s string) string {
	if w.r.opts.Level == Aggressive && s != "" && !strings.Contains(s, "${{") {
		return w.r.lookup(catValue, s)
	}
	return w.free(s)
}

func (w rewriter) name(cat, s string) string {
	return w.r.lookup(cat, s)
}

// runnerLabel keeps GitHub-hosted runner labels and generic selector keywords readable;
// anything else is a self-hosted label and gets its placeholder.
func (w rewriter) runnerLabel(s string) string {
	if keepRunners[strings.ToLower(s)] {
		return s
	}
	return w.r.lookup(catRunner, s)
}

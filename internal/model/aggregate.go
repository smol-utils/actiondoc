package model

import "strings"

// Reference is a referenced secret or variable together with the human-readable sites
// where it is used.
type Reference struct {
	Name  string   `json:"name"`
	Sites []string `json:"sites,omitempty"`
}

// References holds the secrets and variables referenced by a workflow, deduplicated and
// in first-seen order.
type References struct {
	Secrets []Reference `json:"secrets,omitempty"`
	Vars    []Reference `json:"vars,omitempty"`
}

// Empty reports whether no references were found.
func (r References) Empty() bool {
	return len(r.Secrets) == 0 && len(r.Vars) == 0
}

// ScanReferences collects every ${{ secrets.X }} and ${{ vars.Y }} reference in a workflow
// across job/step `if:`, `run:`, `with:`, `env:`, and forwarded `secrets:` values,
// recording the site of each use. It is exported so call-graph consumers can union the
// references of every reachable workflow into a single transitive-requirements view.
func ScanReferences(w *Workflow) References {
	sc := newRefScan()
	for _, kv := range w.Env {
		sc.scan(kv.Value, "workflow env `"+kv.Key+"`")
	}
	for ji := range w.Jobs {
		job := &w.Jobs[ji]
		jobLabel := "job `" + job.ID + "`"
		sc.scanCondition(job.If, jobLabel+" (if)")
		for _, kv := range job.Env {
			sc.scan(kv.Value, jobLabel+" env `"+kv.Key+"`")
		}
		for _, kv := range job.With {
			sc.scan(kv.Value, jobLabel+" with `"+kv.Key+"`")
		}
		for _, kv := range job.Secrets {
			sc.scan(kv.Value, jobLabel+" secrets `"+kv.Key+"`")
		}
		for si := range job.Steps {
			step := &job.Steps[si]
			stepLabel := jobLabel + " step `" + stepRefLabel(step, si+1) + "`"
			sc.scanCondition(step.If, stepLabel+" (if)")
			sc.scan(step.Run, stepLabel+" (run)")
			for _, kv := range step.With {
				sc.scan(kv.Value, stepLabel+" with `"+kv.Key+"`")
			}
			for _, kv := range step.Env {
				sc.scan(kv.Value, stepLabel+" env `"+kv.Key+"`")
			}
		}
	}
	return References{Secrets: sc.secrets.list(), Vars: sc.vars.list()}
}

// stepRefLabel is a compact step identifier for reference site labels. It delegates to
// Step.Label so inventories and rendered step titles always name a step the same way.
func stepRefLabel(s *Step, num int) string {
	return s.Label(num)
}

// refScan accumulates secret and variable references during a workflow walk.
type refScan struct {
	secrets *refAcc
	vars    *refAcc
}

func newRefScan() *refScan {
	return &refScan{secrets: newRefAcc(), vars: newRefAcc()}
}

// scan extracts secrets.X / vars.Y references from every ${{ ... }} block in s and records
// them against the given site.
func (sc *refScan) scan(s, site string) {
	if s == "" || !strings.Contains(s, "${{") {
		return
	}
	for _, body := range exprBodies(s) {
		sc.collect(body, site)
	}
}

// scanCondition scans an if: value. A condition is itself a GitHub Actions expression, so
// references are commonly written WITHOUT the ${{ }} delimiters (for example an if: that
// compares secrets.DEPLOY_TOKEN to the empty string). Scan the whole string as an
// expression body; any ${{ }} the author did include is a substring and is still matched.
// This differs from scan(), used for run/with/env, where a bare secrets.X is plain text,
// not an expression, and must not be collected.
func (sc *refScan) scanCondition(s, site string) {
	if s == "" {
		return
	}
	sc.collect(s, site)
}

func (sc *refScan) collect(body, site string) {
	for _, name := range contextRefs(body, "secrets") {
		sc.secrets.add(name, site)
	}
	for _, name := range contextRefs(body, "vars") {
		sc.vars.add(name, site)
	}
}

// refAcc collects names in first-seen order, merging the (deduplicated) sites of each.
type refAcc struct {
	order []string
	sites map[string][]string
	seen  map[string]bool // name\x00site -> true, to avoid duplicate sites
}

func newRefAcc() *refAcc {
	return &refAcc{sites: map[string][]string{}, seen: map[string]bool{}}
}

func (a *refAcc) add(name, site string) {
	if _, ok := a.sites[name]; !ok {
		a.order = append(a.order, name)
	}
	key := name + "\x00" + site
	if a.seen[key] {
		return
	}
	a.seen[key] = true
	a.sites[name] = append(a.sites[name], site)
}

func (a *refAcc) list() []Reference {
	var out []Reference
	for _, n := range a.order {
		out = append(out, Reference{Name: n, Sites: a.sites[n]})
	}
	return out
}

// exprBodies returns the inner text of each ${{ ... }} occurrence in s, in order.
func exprBodies(s string) []string {
	var out []string
	for {
		i := strings.Index(s, "${{")
		if i < 0 {
			break
		}
		rest := s[i+3:]
		j := strings.Index(rest, "}}")
		if j < 0 {
			break
		}
		out = append(out, rest[:j])
		s = rest[j+2:]
	}
	return out
}

// contextRefs finds every `<ctx>.IDENT` reference in an expression body, e.g.
// contextRefs("vars.X && secrets.A || secrets.B", "secrets") -> ["A", "B"]. Names may
// contain hyphens (workflow_call secret/input keys frequently do). The leading boundary
// check avoids matching a longer identifier that merely ends in ctx (mysecrets.X,
// vault-secrets.X) and property access on some other value (steps.secrets.outputs.X,
// where a step happens to have the id "secrets").
func contextRefs(body, ctx string) []string {
	var out []string
	needle := ctx + "."
	for {
		i := strings.Index(body, needle)
		if i < 0 {
			break
		}
		if i > 0 && (isIdentChar(body[i-1]) || body[i-1] == '.') {
			body = body[i+len(needle):]
			continue
		}
		rest := body[i+len(needle):]
		j := 0
		for j < len(rest) && isIdentChar(rest[j]) {
			j++
		}
		if j > 0 {
			out = append(out, rest[:j])
		}
		body = rest[j:]
	}
	return out
}

func isIdentChar(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '_' || c == '-'
}

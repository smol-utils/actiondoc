// Package callgraph builds the directed graph of `uses:` relationships across a scanned
// set of workflows and composite actions: workflow -> reusable workflow (job-level
// `uses:`) and workflow -> composite action (step-level local `uses:`). It backs the
// reusable-workflow and call-graph rendering features and the caller-count index.
//
// Cross-repo references (e.g. owner/repo/.github/workflows/x.yml@ref) are recorded as
// external nodes and never fetched over the network.
package callgraph

import (
	"path/filepath"
	"strings"

	"github.com/smol-utils/actiondoc/internal/model"
)

// EdgeKind distinguishes the kind of `uses:` relationship.
type EdgeKind string

const (
	KindReusable  EdgeKind = "reusable-workflow"
	KindComposite EdgeKind = "composite-action"
)

// Node is a workflow or composite action in (or referenced by) the scan set.
type Node struct {
	ID       string          // canonical key (source path for in-scope, raw ref for external)
	Name     string          // display name
	Path     string          // source file path ("" for external)
	External bool            // true if outside the scan scope (cross-repo); never fetched
	IsAction bool            // true for composite actions
	Workflow *model.Workflow // non-nil for in-scope workflows
	Action   *model.Action   // non-nil for in-scope composite actions
}

// Edge is a `uses:` relationship from a caller to a target node.
type Edge struct {
	FromID   string
	JobID    string // calling job id (job-level reusable call)
	StepName string // calling step name (step-level composite call)
	ToID     string
	Ref      string // raw `uses:` string
	Pin      string // the @<ref> pin, if present (operational hazard worth surfacing)
	Kind     EdgeKind
}

// Source pairs a resolved file path with its parsed model (exactly one of Workflow or
// Action is non-nil).
type Source struct {
	Path     string
	Workflow *model.Workflow
	Action   *model.Action
}

// Graph is the built call graph.
type Graph struct {
	Nodes map[string]*Node
	Edges []Edge
}

// Build constructs the call graph from the parsed sources.
func Build(sources []Source) *Graph {
	g := &Graph{Nodes: map[string]*Node{}}

	// Register in-scope nodes and build resolution indexes.
	workflowByBase := map[string]string{}  // file base name -> node ID
	actionDirSuffix := map[string]string{} // normalized action dir -> node ID
	for _, s := range sources {
		switch {
		case s.Workflow != nil:
			n := &Node{ID: s.Path, Name: displayName(s.Workflow.Name, s.Path), Path: s.Path, Workflow: s.Workflow}
			g.Nodes[n.ID] = n
			workflowByBase[filepath.Base(s.Path)] = n.ID
		case s.Action != nil:
			n := &Node{ID: s.Path, Name: displayName(s.Action.Name, s.Path), Path: s.Path, IsAction: true, Action: s.Action}
			g.Nodes[n.ID] = n
			actionDirSuffix[filepath.ToSlash(filepath.Dir(s.Path))] = n.ID
		}
	}

	resolve := func(raw string) (toID string, kind EdgeKind, pin string) {
		ref, pin := splitPin(raw)
		if isLocal(ref) {
			clean := strings.TrimPrefix(strings.TrimPrefix(ref, "./"), "../")
			if strings.Contains(clean, ".github/workflows/") || strings.HasSuffix(clean, ".yml") || strings.HasSuffix(clean, ".yaml") {
				if id, ok := workflowByBase[filepath.Base(clean)]; ok {
					return id, KindReusable, pin
				}
			}
			// composite action: match by directory suffix
			clean = strings.TrimSuffix(strings.TrimSuffix(clean, "/action.yml"), "/action.yaml")
			for dir, id := range actionDirSuffix {
				if strings.HasSuffix(dir, clean) {
					return id, KindComposite, pin
				}
			}
			return "", KindReusable, pin // local but unresolved
		}
		// cross-repo external reference: record as an external node, do not fetch.
		ext := externalID(raw)
		if _, ok := g.Nodes[ext]; !ok {
			g.Nodes[ext] = &Node{ID: ext, Name: ext, External: true, IsAction: !strings.Contains(ref, "/.github/workflows/")}
		}
		kind = KindReusable
		if !strings.Contains(ref, "/.github/workflows/") {
			kind = KindComposite
		}
		return ext, kind, pin
	}

	for _, s := range sources {
		if s.Workflow == nil {
			continue
		}
		from := s.Path
		for _, job := range s.Workflow.Jobs {
			if job.Uses != "" { // job-level reusable workflow call
				to, kind, pin := resolve(job.Uses)
				g.Edges = append(g.Edges, Edge{FromID: from, JobID: job.ID, ToID: to, Ref: job.Uses, Pin: pin, Kind: kind})
			}
			for _, st := range job.Steps { // step-level local composite call
				if isLocal(st.Uses) {
					to, kind, pin := resolve(st.Uses)
					g.Edges = append(g.Edges, Edge{FromID: from, JobID: job.ID, StepName: st.Name, ToID: to, Ref: st.Uses, Pin: pin, Kind: kind})
				}
			}
		}
	}
	return g
}

// Calls returns the edges originating at node id.
func (g *Graph) Calls(id string) []Edge {
	var out []Edge
	for _, e := range g.Edges {
		if e.FromID == id {
			out = append(out, e)
		}
	}
	return out
}

// CalledBy returns the edges targeting node id.
func (g *Graph) CalledBy(id string) []Edge {
	var out []Edge
	for _, e := range g.Edges {
		if e.ToID == id {
			out = append(out, e)
		}
	}
	return out
}

// CallerCount returns the number of distinct caller nodes for id (used by the index and
// to collapse high fan-in cases where a reusable workflow has many callers).
func (g *Graph) CallerCount(id string) int {
	seen := map[string]bool{}
	for _, e := range g.Edges {
		if e.ToID == id {
			seen[e.FromID] = true
		}
	}
	return len(seen)
}

// Reachable returns the set of node IDs transitively reachable from id via `uses:`
// edges (excluding id itself), in deterministic discovery order.
func (g *Graph) Reachable(id string) []string {
	seen := map[string]bool{id: true}
	var order []string
	queue := []string{id}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for _, e := range g.Calls(cur) {
			if e.ToID == "" || seen[e.ToID] {
				continue
			}
			seen[e.ToID] = true
			order = append(order, e.ToID)
			queue = append(queue, e.ToID)
		}
	}
	return order
}

// IsEntryPoint reports whether id is an in-scope workflow with a trigger other than
// only `workflow_call` (i.e. a human/automation-facing entry point, not a pure
// reusable workflow).
func (g *Graph) IsEntryPoint(id string) bool {
	n := g.Nodes[id]
	if n == nil || n.Workflow == nil {
		return false
	}
	for _, t := range n.Workflow.On {
		if t != "workflow_call" {
			return true
		}
	}
	return false
}

// EntryPoints returns the IDs of all entry-point workflows.
func (g *Graph) EntryPoints() []string {
	var out []string
	for id := range g.Nodes {
		if g.IsEntryPoint(id) {
			out = append(out, id)
		}
	}
	return out
}

func isLocal(ref string) bool {
	return strings.HasPrefix(ref, "./") || strings.HasPrefix(ref, "../")
}

// splitPin separates a `uses:` value from its trailing @<ref> pin. Local paths have no
// pin; cross-repo refs do (owner/repo/...@v1).
func splitPin(raw string) (ref, pin string) {
	if isLocal(raw) {
		return raw, ""
	}
	if i := strings.LastIndex(raw, "@"); i >= 0 {
		return raw[:i], raw[i+1:]
	}
	return raw, ""
}

func externalID(raw string) string { return raw }

func displayName(name, path string) string {
	if name != "" {
		return name
	}
	return filepath.Base(path)
}

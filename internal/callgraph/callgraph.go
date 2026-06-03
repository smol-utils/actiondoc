// Package callgraph builds the directed graph of `uses:` relationships across a scanned
// set of workflows and composite actions: workflow -> reusable workflow (job-level
// `uses:`) and workflow -> composite action (step-level local `uses:`). It backs the
// reusable-workflow and call-graph rendering features and the caller-count index.
//
// Cross-repo references (e.g. owner/repo/.github/workflows/x.yml@ref) are recorded as
// external nodes and never fetched over the network.
package callgraph

import (
	"fmt"
	"path"
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

	// Anchor is the node's document anchor slug, assigned by the assembler that decides
	// section order (it owns duplicate-name disambiguation). Renderers building
	// cross-links use it when set, so links always agree with the table of contents;
	// when empty (single-file rendering, no assembly), links fall back to slugging Name.
	Anchor string
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

	// Register in-scope nodes and build resolution indexes. Workflows are indexed by base
	// filename (scan paths and `uses:` refs often live under different parent dirs, so
	// basename is the reliable join key); when several share a basename the full paths are
	// kept so resolution can disambiguate by longest shared path suffix. Composite actions
	// are matched by their normalized directory.
	workflowsByBase := map[string][]string{} // base filename -> node IDs (paths)
	actionDirSuffix := map[string]string{}   // normalized action dir -> node ID
	for _, s := range sources {
		switch {
		case s.Workflow != nil:
			n := &Node{ID: s.Path, Name: displayName(s.Workflow.Name, s.Path), Path: s.Path, Workflow: s.Workflow}
			g.Nodes[n.ID] = n
			// Key by base filename (normalize to slashes first so the key is consistent
			// whether the scan path is OS-native or already slash-separated); store the
			// ORIGINAL path (= node ID) so resolution returns a value that matches g.Nodes.
			base := path.Base(filepath.ToSlash(s.Path))
			workflowsByBase[base] = append(workflowsByBase[base], s.Path)
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
			// A path into .github/workflows or a bare .yml/.yaml file is a reusable-workflow
			// call; anything else (e.g. ./.github/actions/x) is a composite action.
			looksWorkflow := strings.Contains(clean, ".github/workflows/") || strings.HasSuffix(clean, ".yml") || strings.HasSuffix(clean, ".yaml")
			if looksWorkflow {
				if id := resolveWorkflow(workflowsByBase, clean); id != "" {
					return id, KindReusable, pin
				}
			}
			// composite action: match on a path-segment boundary, choosing the most
			// specific (longest) match so resolution is deterministic regardless of
			// map iteration order.
			clean = strings.TrimSuffix(strings.TrimSuffix(clean, "/action.yml"), "/action.yaml")
			if id := longestSuffixMatch(actionDirSuffix, clean); id != "" {
				return id, KindComposite, pin
			}
			// Local but unresolved: classify by what the ref looked like.
			if looksWorkflow {
				return "", KindReusable, pin
			}
			return "", KindComposite, pin
		}
		// cross-repo external reference: record as an external node (keyed without the
		// @pin so different pins of the same target collapse), do not fetch.
		isWorkflow := strings.Contains(ref, "/.github/workflows/")
		if _, ok := g.Nodes[ref]; !ok {
			g.Nodes[ref] = &Node{ID: ref, Name: ref, External: true, IsAction: !isWorkflow}
		}
		if isWorkflow {
			return ref, KindReusable, pin
		}
		return ref, KindComposite, pin
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
			for si, st := range job.Steps { // step-level local composite call
				if isLocal(st.Uses) {
					to, kind, pin := resolve(st.Uses)
					g.Edges = append(g.Edges, Edge{FromID: from, JobID: job.ID, StepName: stepLabel(st, si), ToID: to, Ref: st.Uses, Pin: pin, Kind: kind})
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

// longestSuffixMatch resolves a cleaned local ref against an index keyed by normalized
// dir, returning the node ID of the most specific match. The scan key and the ref may be
// qualified to different depths (a scan dir can be shallower than the ref, or vice
// versa), so a match holds when either is a path-segment suffix of the other. Among
// matches the longest key wins with a lexical tiebreak so the result is deterministic
// regardless of map iteration order. Returns "" when nothing matches. Used for
// composite-action resolution.
func longestSuffixMatch(index map[string]string, ref string) string {
	var bestKey, bestID string
	for key, id := range index {
		if key == ref || strings.HasSuffix(key, "/"+ref) || strings.HasSuffix(ref, "/"+key) {
			if bestID == "" || len(key) > len(bestKey) || (len(key) == len(bestKey) && key < bestKey) {
				bestKey, bestID = key, id
			}
		}
	}
	return bestID
}

// resolveWorkflow resolves a cleaned reusable-workflow `uses:` ref to a node ID. Scan
// paths and refs frequently live under different parent directories (the scan root is not
// necessarily .github/workflows), so the join key is the base filename. When a single
// workflow has that basename it wins outright; when several scanned workflows share the
// basename, pick the one whose path shares the longest trailing path segment with the
// ref, so colliding basenames in different directories can't silently resolve to the
// wrong file. Returns "" when nothing matches.
func resolveWorkflow(byBase map[string][]string, ref string) string {
	// ref is a slash-separated YAML path, so use path.Base (not filepath.Base, which
	// would not split on '/' on Windows).
	paths := byBase[path.Base(ref)]
	switch len(paths) {
	case 0:
		return ""
	case 1:
		return paths[0]
	}
	best, bestScore := "", -1
	for _, p := range paths {
		score := sharedSuffixSegments(p, ref)
		if score > bestScore || (score == bestScore && p < best) {
			best, bestScore = p, score
		}
	}
	return best
}

// sharedSuffixSegments counts the trailing path segments two paths share (e.g. "a/b/c"
// and "x/b/c" share 2). Either argument may use OS-native separators (scan paths) or
// slashes (YAML refs), so both are normalized to slashes before comparison.
func sharedSuffixSegments(a, b string) int {
	as := strings.Split(filepath.ToSlash(a), "/")
	bs := strings.Split(filepath.ToSlash(b), "/")
	n := 0
	for i, j := len(as)-1, len(bs)-1; i >= 0 && j >= 0 && as[i] == bs[j]; i, j = i-1, j-1 {
		n++
	}
	return n
}

func displayName(name, path string) string {
	if name != "" {
		return name
	}
	return filepath.Base(path)
}

// stepLabel identifies a step for the call-graph tree: its name, then id, then a
// positional fallback. Unnamed steps must still get a non-empty label so the renderer
// shows them in the "job / step" form rather than collapsing to the job-level form.
func stepLabel(s model.Step, idx int) string {
	switch {
	case s.Name != "":
		return s.Name
	case s.ID != "":
		return s.ID
	default:
		return fmt.Sprintf("step %d", idx+1)
	}
}

package renderer

import (
	"fmt"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/smol-utils/actiondoc/internal/callgraph"
	"github.com/smol-utils/actiondoc/internal/model"
)

// This file renders the reusable-workflow and call-graph surface: caller jobs that
// invoke another workflow via `uses:` (with forwarded inputs/secrets and a cross-link
// to the callee), the downstream call-graph tree on entry-point workflows, the upstream
// "Called by" chain on reusable workflows, and a flat aggregation of the requirements
// reachable across the whole chain. Everything here is a read-only consumer of the
// already-built callgraph: it walks Calls/CalledBy/Reachable and never mutates the graph.

// renderCallerJob renders the body of a job that calls a reusable workflow (it uses
// `uses:` in place of `runs-on:`/`steps:`). The shape is a small properties table whose
// first row links to the callee, followed by the forwarded `with:` inputs and `secrets:`
// (an explicit mapping, or `secrets: inherit` rendered verbatim). The caller passes the
// calling workflow's node id (fromID) so the callee edge can be located in the graph.
func renderCallerJob(b *strings.Builder, job *model.Job, g *callgraph.Graph, fromID string) {
	b.WriteString("| Property | Value |\n")
	b.WriteString("|----------|-------|\n")
	fmt.Fprintf(b, "| Uses workflow | %s |\n", callerUsesCell(g, fromID, job.ID, job.Uses))
	// A caller job's matrix multiplies the reusable calls; its axes are as much a part of
	// the job's surface as a regular job's.
	if len(job.Matrix) > 0 {
		fmt.Fprintf(b, "| Matrix | %s |\n", matrixCell(job.Matrix, job.MatrixAdjusted))
	}
	if len(job.Needs) > 0 {
		fmt.Fprintf(b, "| Depends on | %s |\n", codelist(job.Needs))
	}
	if job.If != "" {
		// Trim first: literal-block conditions carry a trailing newline that would
		// otherwise render as a dangling <br> (matches the normal job renderer).
		fmt.Fprintf(b, "| Condition | `%s` |\n", escapeCellCode(strings.TrimSpace(job.If)))
	}
	b.WriteString("\n")

	// A caller job can still declare its own permissions, environment binding, env,
	// concurrency, and defaults; render that declared surface (it would otherwise be
	// dropped, since caller jobs skip the normal job body).
	renderJobSurface(b, job)

	// Forwarded inputs whose value is empty/unset carry no information and would otherwise
	// render as a column of bare dashes; omit them (mirrors the step `with:`/`env:` path) and
	// drop the header entirely when nothing meaningful remains.
	type forwarded struct{ key, value string }
	var inputs []forwarded
	for _, kv := range job.With {
		value := oneLine(kv.Value)
		if value == "" {
			continue
		}
		inputs = append(inputs, forwarded{key: kv.Key, value: value})
	}
	if len(inputs) > 0 {
		b.WriteString("#### Inputs forwarded\n\n")
		for _, in := range inputs {
			fmt.Fprintf(b, "- `%s`: %s\n", in.key, codeSpan(in.value))
		}
		b.WriteString("\n")
	}

	switch {
	case job.SecretsInherit:
		b.WriteString("#### Secrets forwarded\n\n")
		b.WriteString("- `secrets: inherit` (all caller secrets are passed to the callee)\n\n")
	case len(job.Secrets) > 0:
		b.WriteString("#### Secrets forwarded\n\n")
		for _, kv := range job.Secrets {
			fmt.Fprintf(b, "- `%s`: %s\n", kv.Key, codeSpan(oneLine(kv.Value)))
		}
		b.WriteString("\n")
	}

	// A caller job may also carry ActionDoc tags (@secret/@env/@output/@example/@see);
	// render them like any other job rather than dropping them on the early return.
	renderJobTags(b, job)
}

// callerUsesCell renders the "Uses workflow" table cell: an anchor cross-link to the
// callee's rendered section when it is in scope, or the raw `uses:` string when there is
// no graph to resolve against.
func callerUsesCell(g *callgraph.Graph, fromID, jobID, rawUses string) string {
	if g != nil {
		for _, e := range g.Calls(fromID) {
			if e.JobID == jobID && e.StepName == "" && e.Kind == callgraph.KindReusable {
				return calleeLink(g, e)
			}
		}
	}
	return "`" + escapeCellCode(rawUses) + "`"
}

// calleeLink renders a reference to a reusable-workflow callee: an in-scope callee
// becomes an anchor cross-link to its rendered title; an out-of-scope (cross-repo)
// callee renders as inline code with its `@ref` pin surfaced (and is never fetched).
func calleeLink(g *callgraph.Graph, e callgraph.Edge) string {
	n := g.Nodes[e.ToID]
	if n == nil {
		// The ref points outside what was scanned: a path that exists in the repository
		// but was not discovered, or one that only exists at runtime (e.g. a checkout
		// into a subdirectory). Either way the tool did not see it -- which is different
		// from the ref being broken.
		return "`" + escapeCellCode(e.Ref) + "` (outside scan scope)"
	}
	if n.External {
		ref := n.Name
		if e.Pin != "" {
			ref += "@" + e.Pin
		}
		return "`" + escapeCellCode(ref) + "` (external)"
	}
	// An in-scope edge carrying a pin is a cross-repo self-reference (the repo calling
	// its own workflow at a branch/tag); keep the pin visible next to the link so the
	// reader knows the pinned version is what actually runs.
	link := fmt.Sprintf("[%s](#%s)", mdLinkLabel(n.Name), nodeAnchor(n))
	if e.Pin != "" {
		link += " (`@" + escapeCellCode(e.Pin) + "`)"
	}
	return link
}

// nodeAnchor is the anchor slug for an in-scope node's rendered section: the
// assembler-assigned anchor when present (which carries duplicate-name disambiguation),
// otherwise the slug of the node's name.
func nodeAnchor(n *callgraph.Node) string {
	if n.Anchor != "" {
		return n.Anchor
	}
	return anchor(n.Name)
}

// renderCallGraph renders the downstream `uses:` tree rooted at an entry-point workflow. It
// is suppressed for flat workflows (no outgoing `uses:`) and for pure reusable workflows
// (which get a "Called by" section instead). The tree is a nested Markdown list: each callee
// is a cross-link to its rendered section (or plain inline code when out of scope), and
// repeated sibling subtrees are folded to one "(xN)" representative.
func renderCallGraph(b *strings.Builder, g *callgraph.Graph, id string) {
	if g == nil || !g.IsEntryPoint(id) {
		return
	}
	edges := g.Calls(id)
	if len(edges) == 0 {
		return
	}
	root := treeNode{}
	path := []string{id}
	for _, e := range edges {
		root.children = append(root.children, callEdgeNode(g, e, path))
	}
	root.children = collapseSiblings(root.children)

	// The section heading already names the workflow, and its file/triggers are shown in
	// the property table directly above, so the tree renders straight from its children
	// rather than restating the root.
	b.WriteString("## Call graph (rooted at this workflow)\n\n")
	renderTreeListChildren(b, root.children, 0)
	b.WriteString("\n")
}

// callEdgeNode builds the subtree for a single outgoing `uses:` edge, recursing into the
// callee's own calls. The path slice records the node ids on the current branch so a
// cyclic `uses:` reference stops at a leaf marked "(cycle)" instead of recursing forever.
func callEdgeNode(g *callgraph.Graph, e callgraph.Edge, path []string) treeNode {
	label, rep := callEdgeLabels(g, e)
	node := treeNode{label: label, repLabel: rep}
	if e.ToID == "" {
		// Outside the scan scope: a leaf, but not a cycle.
		return node
	}
	if containsStr(path, e.ToID) {
		node.label += " (cycle)"
		node.repLabel += " (cycle)"
		return node
	}
	next := append(append([]string{}, path...), e.ToID)
	for _, ce := range g.Calls(e.ToID) {
		node.children = append(node.children, callEdgeNode(g, ce, next))
	}
	return node
}

// callEdgeLabels describes one call site in the downstream tree. The full label leads with
// the calling job (and step name, for a composite-action call) in an inline-code span, then
// links to the callee; the representative label drops the job/step context (which is what
// differs across collapsed siblings) and keeps the emphasized callee link.
func callEdgeLabels(g *callgraph.Graph, e callgraph.Edge) (label, rep string) {
	callee := calleeMarkdown(g, e)
	ctx := e.JobID
	if e.StepName != "" {
		ctx = e.JobID + " / " + e.StepName
	}
	label = fmt.Sprintf("%s uses %s", codeSpan(ctx), callee)
	rep = "uses **" + callee + "**"
	return label, rep
}

// calleeMarkdown renders a callee for the downstream list. An in-scope workflow or composite
// action becomes a cross-link to its rendered section (with the `@pin` surfaced for a
// cross-repo self-reference); an external callee renders as plain inline code carrying its
// `@pin` (never a link, there is no in-scope section); a callee outside the scan scope
// renders as plain inline code annotated as such.
func calleeMarkdown(g *callgraph.Graph, e callgraph.Edge) string {
	n := g.Nodes[e.ToID]
	if n == nil {
		return codeSpan(e.Ref) + " (outside scan scope)"
	}
	if n.External {
		ref := n.Name
		if e.Pin != "" {
			ref += "@" + e.Pin
		}
		return codeSpan(ref)
	}
	// In-scope: link to the callee's rendered section. Composite actions all live in an
	// action.yml, so the `uses:` ref is the only distinguishing display; workflows show the
	// file base name.
	display := filepath.Base(n.Path)
	if n.IsAction {
		display = e.Ref
	}
	link := fmt.Sprintf("[%s](#%s)", mdLinkLabel(display), nodeAnchor(n))
	if e.Pin != "" {
		link += " (" + codeSpan("@"+e.Pin) + ")"
	}
	return link
}

// calleeDisplay is the bare display string for a callee in an ASCII tree: the file base
// name for in-scope workflows, the `uses:` reference for in-scope composite actions
// (whose file base name is always just "action.yml"), or the raw cross-repo reference
// (with `@ref` pin) for external ones.
func calleeDisplay(g *callgraph.Graph, e callgraph.Edge) string {
	n := g.Nodes[e.ToID]
	if n == nil {
		return e.Ref + " (outside scan scope)"
	}
	if n.External {
		if e.Pin != "" {
			return n.Name + "@" + e.Pin
		}
		return n.Name
	}
	if n.IsAction {
		return e.Ref
	}
	// An in-scope edge carrying a pin is a cross-repo self-reference; keep the pinned
	// version visible in the tree label.
	if e.Pin != "" {
		return filepath.Base(n.Path) + "@" + e.Pin
	}
	return filepath.Base(n.Path)
}

// renderCalledBy renders the upstream caller chain on a workflow that is invoked by others:
// immediate callers at the top, each expanded to its own callers up to the entry points,
// which are marked. It reuses the same nested-list renderer as the downstream call graph;
// only the walk direction (CalledBy) differs. Each caller links to the specific calling job
// heading, and repeated sibling subtrees are folded to one "(xN)" representative.
func renderCalledBy(b *strings.Builder, g *callgraph.Graph, id string) {
	if g == nil {
		return
	}
	callers := g.CalledBy(id)
	if len(callers) == 0 {
		return
	}
	root := treeNode{}
	path := []string{id}
	for _, e := range callers {
		root.children = append(root.children, calledByNode(g, e, path))
	}
	root.children = collapseSiblings(root.children)

	// The root would be just this workflow's file basename, already named by the section
	// heading, so the tree renders straight from its children for consistency with the
	// downstream call graph.
	b.WriteString("## Called by\n\n")
	renderTreeListChildren(b, root.children, 0)
	b.WriteString("\n")
}

// calledByNode builds the subtree for a single caller edge, recursing upward into that
// caller's own callers. The path slice guards against cyclic call relationships, stopping at
// a leaf marked "(cycle)".
func calledByNode(g *callgraph.Graph, e callgraph.Edge, path []string) treeNode {
	label, rep := calledByLabels(g, e)
	node := treeNode{label: label, repLabel: rep}
	if containsStr(path, e.FromID) {
		node.label += " (cycle)"
		node.repLabel += " (cycle)"
		return node
	}
	next := append(append([]string{}, path...), e.FromID)
	for _, ce := range g.CalledBy(e.FromID) {
		node.children = append(node.children, calledByNode(g, ce, next))
	}
	return node
}

// calledByLabels describes one caller in the upstream list: the calling file and job, with
// an "entry point" marker when that caller is itself a human/automation-facing trigger. The
// full label links to the specific calling JOB heading (via the document-wide JobAnchors,
// falling back to the file's section anchor); the representative label drops the differing
// job id and links to the file's section.
func calledByLabels(g *callgraph.Graph, e callgraph.Edge) (label, rep string) {
	n := g.Nodes[e.FromID]
	base := e.FromID
	section, job := "", ""
	if n != nil {
		if n.Path != "" {
			base = filepath.Base(n.Path)
		}
		section = nodeAnchor(n)
		job = section
		if n.Workflow != nil {
			if i := jobIndex(n.Workflow, e.JobID); i >= 0 && i < len(n.JobAnchors) {
				job = n.JobAnchors[i]
			}
		}
	}
	entry := ""
	if g.IsEntryPoint(e.FromID) {
		entry = " - entry point"
	}
	label = fmt.Sprintf("[%s](#%s) (job: %s)%s", mdLinkLabel(base), job, codeSpan(e.JobID), entry)
	rep = fmt.Sprintf("**[%s](#%s)**%s", mdLinkLabel(base), section, entry)
	return label, rep
}

// renderTransitiveRequirements aggregates, across the entry point and everything
// reachable from it, what the whole pipeline needs by contract: the secret names each hop
// declares or forwards (workflow_call.secrets, forwarded `secrets:` keys, and `@secret`
// tags), and the external workflows it pulls in. It answers "what does this whole chain
// need?" without walking every hop.
//
// Permission grants are deliberately NOT aggregated here: they are already shown per-workflow
// (the `## Permissions` section), per-job (`**Permissions:**`), and repo-wide (the doc-level
// "Permissions across this repository" summary), so a transitive roll-up only repeats them.
//
// Expression-referenced secret/variable NAMES are deliberately NOT aggregated here: the
// document-level secrets/variables inventory already maps every referenced name to the
// workflows that use it, and the per-workflow "Referenced secrets and variables" table
// carries the within-workflow site detail. This section is the contract view; those are the
// usage views.
//
// Scope note: names are reported as written at each hop; values are not traced through
// per-hop `secrets:` renames.
func renderTransitiveRequirements(b *strings.Builder, g *callgraph.Graph, id string) {
	if g == nil || !g.IsEntryPoint(id) {
		return
	}
	reach := g.Reachable(id)
	if len(reach) == 0 {
		return
	}

	secrets := map[string]bool{}
	externals := map[string]bool{}
	for _, nid := range append([]string{id}, reach...) {
		n := g.Nodes[nid]
		if n == nil || n.External {
			continue
		}
		collectSecretNames(n, secrets)
		// External references are collected from this node's outgoing edges (not from the
		// external nodes themselves) so the `@ref` pin each call site uses is preserved.
		for _, e := range g.Calls(nid) {
			if to := g.Nodes[e.ToID]; to != nil && to.External {
				externals[calleeDisplay(g, e)] = true
			}
		}
	}
	if len(secrets) == 0 && len(externals) == 0 {
		return
	}

	b.WriteString("## Transitive requirements (from full call graph)\n\n")
	if len(secrets) > 0 {
		fmt.Fprintf(b, "Secrets required (declared/forwarded names): %s\n\n", codelist(sortedKeys(secrets)))
	}
	if len(externals) > 0 {
		fmt.Fprintf(b, "External workflows referenced: %s\n\n", codelist(sortedKeys(externals)))
	}
}

// collectSecretNames adds the literal secret names a node declares or forwards into set:
// workflow/job/step `@secret` tags and the keys of forwarded `secrets:` maps on caller
// jobs (for composite actions, the action-level `@secret` tags).
func collectSecretNames(n *callgraph.Node, set map[string]bool) {
	switch {
	case n.Workflow != nil:
		for _, p := range n.Workflow.Tags.Secrets {
			set[p.Name] = true
		}
		// A reusable workflow's declared workflow_call.secrets are part of its contract:
		// callers using `secrets: inherit` (or not forwarding an explicit key) still
		// require them, so they belong in the transitive requirements.
		if t := n.Workflow.Triggers; t != nil && t.Call != nil {
			for _, s := range t.Call.Secrets {
				set[s.Name] = true
			}
		}
		for _, job := range n.Workflow.Jobs {
			for _, p := range job.Tags.Secrets {
				set[p.Name] = true
			}
			for _, kv := range job.Secrets {
				set[kv.Key] = true
			}
			for _, st := range job.Steps {
				for _, p := range st.Tags.Secrets {
					set[p.Name] = true
				}
			}
		}
	case n.Action != nil:
		for _, p := range n.Action.Tags.Secrets {
			set[p.Name] = true
		}
	}
}

// treeNode is a single item in a dependency tree plus its children. Both the downstream
// call graph and the upstream "called by" chain build trees of these and hand them to
// renderTreeList, so there is exactly one tree-drawing implementation.
//
// label is the full Markdown the item renders as (a job/step context code span plus a
// cross-link to the callee or caller). repLabel is the same line with the differing
// job/step context dropped and the shared callee/caller emphasized; it doubles as the
// node's grouping identity (two siblings collapse only when their repLabel and their whole
// children block both match) and as the displayed label once a group is collapsed to a
// single "(xN)" representative. A node with an empty repLabel (the root, or a back-
// reference note) never participates in grouping.
type treeNode struct {
	label    string
	repLabel string
	children []treeNode
}

// collapseSiblings collapses repeated subtrees that sit side by side under the same parent.
// Both trees show the same shape of repetition: many distinct caller jobs invoking the SAME
// callee (downstream) -- or many distinct calling jobs in the same caller file reaching this
// workflow (upstream) -- each carrying an identical descendant block (e.g. airflow's dozen
// tests-* jobs, each `uses run-unit-tests.yml` with the same four-step subtree). The only
// thing that differs across them is the job id, which the repLabel deliberately omits.
//
// Children are grouped by repLabel + their (already collapsed) children block; each group is
// replaced by its first occurrence (so sibling order stays deterministic), and a group of
// more than one gets its differing job id dropped (repLabel becomes the label) with an
// "(xN)" count appended. Grandchildren are collapsed first, so a repeated subtree at any
// depth folds before its parent is keyed -- nested groups collapse too.
func collapseSiblings(children []treeNode) []treeNode {
	for i := range children {
		children[i].children = collapseSiblings(children[i].children)
	}
	var out []treeNode
	pos := map[string]int{}
	count := map[string]int{}
	for _, c := range children {
		// A node with no grouping identity (the root, or a note) is passed through verbatim.
		if c.repLabel == "" {
			out = append(out, c)
			continue
		}
		key := c.repLabel + "\x00" + childrenKey(c.children)
		if i, seen := pos[key]; seen {
			count[key]++
			_ = i
			continue
		}
		pos[key] = len(out)
		count[key] = 1
		out = append(out, c)
	}
	for key, i := range pos {
		if count[key] > 1 {
			out[i].label = fmt.Sprintf("%s (x%d)", out[i].repLabel, count[key])
		}
	}
	return out
}

// jobIndex returns the position of the job with the given id in a workflow's job list (the
// index into the workflow node's document-wide JobAnchors), or -1 when absent.
func jobIndex(w *model.Workflow, jobID string) int {
	for i := range w.Jobs {
		if w.Jobs[i].ID == jobID {
			return i
		}
	}
	return -1
}

// childrenKey serializes a node's full descendant structure (depth + label per line),
// ignoring the node's own label, so two caller nodes whose upstream chains are identical
// share a key regardless of their own differing job names.
func childrenKey(children []treeNode) string {
	if len(children) == 0 {
		return ""
	}
	var sb strings.Builder
	var walk func(ns []treeNode, depth int)
	walk = func(ns []treeNode, depth int) {
		for _, n := range ns {
			sb.WriteString(strconv.Itoa(depth))
			sb.WriteByte(':')
			sb.WriteString(n.label)
			sb.WriteByte('\n')
			walk(n.children, depth+1)
		}
	}
	walk(children, 0)
	return sb.String()
}

// renderTreeList writes a tree as a Markdown nested list: the root label as a plain lead
// line (a paragraph), then each child as a `-` list item, two spaces of indentation per
// depth. Labels are Markdown (cross-links plus inline-code context), so the listed
// callee/caller names -- which are headings elsewhere in the same document -- are clickable,
// unlike the old fenced-code-block rendering.
func renderTreeList(b *strings.Builder, root treeNode) {
	b.WriteString(root.label)
	b.WriteString("\n\n")
	renderTreeListChildren(b, root.children, 0)
}

func renderTreeListChildren(b *strings.Builder, children []treeNode, depth int) {
	indent := strings.Repeat("  ", depth)
	for _, c := range children {
		b.WriteString(indent + "- " + c.label + "\n")
		renderTreeListChildren(b, c.children, depth+1)
	}
}

// sortedKeys returns the keys of a string set in deterministic alphabetical order.
func sortedKeys(set map[string]bool) []string {
	out := make([]string, 0, len(set))
	for k := range set {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// containsStr reports whether s is present in xs.
func containsStr(xs []string, s string) bool {
	for _, x := range xs {
		if x == s {
			return true
		}
	}
	return false
}

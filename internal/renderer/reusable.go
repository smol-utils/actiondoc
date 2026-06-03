package renderer

import (
	"fmt"
	"path/filepath"
	"sort"
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
		fmt.Fprintf(b, "| Condition | `%s` |\n", escapeCell(strings.TrimSpace(job.If)))
	}
	b.WriteString("\n")

	// A caller job can still declare its own permissions, environment binding, env,
	// concurrency, and defaults; render that declared surface (it would otherwise be
	// dropped, since caller jobs skip the normal job body).
	renderJobSurface(b, job)

	if len(job.With) > 0 {
		b.WriteString("#### Inputs forwarded\n\n")
		for _, kv := range job.With {
			fmt.Fprintf(b, "- `%s`: %s\n", kv.Key, codeSpan(oneLine(kv.Value)))
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
	return "`" + escapeCell(rawUses) + "`"
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
		return "`" + escapeCell(e.Ref) + "` (outside scan scope)"
	}
	if n.External {
		ref := n.Name
		if e.Pin != "" {
			ref += "@" + e.Pin
		}
		return "`" + escapeCell(ref) + "` (external)"
	}
	// An in-scope edge carrying a pin is a cross-repo self-reference (the repo calling
	// its own workflow at a branch/tag); keep the pin visible next to the link so the
	// reader knows the pinned version is what actually runs.
	link := fmt.Sprintf("[%s](#%s)", mdLinkLabel(n.Name), nodeAnchor(n))
	if e.Pin != "" {
		link += " (`@" + escapeCell(e.Pin) + "`)"
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

// renderCallGraph renders the downstream `uses:` tree rooted at an entry-point workflow
// (item 12). It is suppressed for flat workflows (no outgoing `uses:`) and for pure
// reusable workflows (which get a "Called by" section instead). The tree is plain ASCII
// inside a code fence: full names, no middle-truncation.
func renderCallGraph(b *strings.Builder, g *callgraph.Graph, id string) {
	if g == nil || !g.IsEntryPoint(id) {
		return
	}
	edges := g.Calls(id)
	if len(edges) == 0 {
		return
	}
	root := treeNode{label: entryRootLabel(g, id)}
	path := []string{id}
	for _, e := range edges {
		root.children = append(root.children, callEdgeNode(g, e, path))
	}

	b.WriteString("## Call graph (rooted at this workflow)\n\n")
	b.WriteString("```\n")
	renderTree(b, root)
	b.WriteString("```\n\n")
}

// callEdgeNode builds the subtree for a single outgoing `uses:` edge, recursing into the
// callee's own calls. The path slice records the node ids on the current branch so a
// cyclic `uses:` reference stops instead of recursing forever.
func callEdgeNode(g *callgraph.Graph, e callgraph.Edge, path []string) treeNode {
	node := treeNode{label: callEdgeLabel(g, e)}
	if e.ToID == "" || containsStr(path, e.ToID) {
		return node
	}
	next := append(append([]string{}, path...), e.ToID)
	for _, ce := range g.Calls(e.ToID) {
		node.children = append(node.children, callEdgeNode(g, ce, next))
	}
	return node
}

// callEdgeLabel describes one call site in the downstream tree: the calling job (and step
// name, for a composite-action call) and the callee it targets.
func callEdgeLabel(g *callgraph.Graph, e callgraph.Edge) string {
	callee := calleeDisplay(g, e)
	if e.StepName != "" {
		return fmt.Sprintf("%s / %s (uses %s)", e.JobID, e.StepName, callee)
	}
	return fmt.Sprintf("%s (uses %s)", e.JobID, callee)
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

// entryRootLabel is the root line of the downstream call-graph tree: the entry-point file
// name annotated with its triggers, e.g. "release.yml [workflow_dispatch]".
func entryRootLabel(g *callgraph.Graph, id string) string {
	n := g.Nodes[id]
	base := id
	if n != nil && n.Path != "" {
		base = filepath.Base(n.Path)
	}
	if n != nil && n.Workflow != nil && len(n.Workflow.On) > 0 {
		return base + " [" + strings.Join(n.Workflow.On, ", ") + "]"
	}
	return base
}

// renderCalledBy renders the upstream caller chain on a workflow that is invoked by
// others (item 13): immediate callers at the top, each expanded to its own callers up to
// the entry points, which are marked. It reuses the same ASCII tree renderer as the
// downstream call graph; only the walk direction (CalledBy) differs.
func renderCalledBy(b *strings.Builder, g *callgraph.Graph, id string) {
	if g == nil {
		return
	}
	callers := g.CalledBy(id)
	if len(callers) == 0 {
		return
	}
	n := g.Nodes[id]
	base := id
	if n != nil && n.Path != "" {
		base = filepath.Base(n.Path)
	}
	root := treeNode{label: base}
	path := []string{id}
	for _, e := range callers {
		root.children = append(root.children, calledByNode(g, e, path))
	}

	b.WriteString("## Called by\n\n")
	b.WriteString("```\n")
	renderTree(b, root)
	b.WriteString("```\n\n")
}

// calledByNode builds the subtree for a single caller edge, recursing upward into that
// caller's own callers. The path slice guards against cyclic call relationships.
func calledByNode(g *callgraph.Graph, e callgraph.Edge, path []string) treeNode {
	node := treeNode{label: calledByLabel(g, e)}
	if containsStr(path, e.FromID) {
		return node
	}
	next := append(append([]string{}, path...), e.FromID)
	for _, ce := range g.CalledBy(e.FromID) {
		node.children = append(node.children, calledByNode(g, ce, next))
	}
	return node
}

// calledByLabel describes one caller in the upstream tree: the calling file and job, with
// an "entry point" marker when that caller is itself a human/automation-facing trigger.
func calledByLabel(g *callgraph.Graph, e callgraph.Edge) string {
	n := g.Nodes[e.FromID]
	base := e.FromID
	if n != nil && n.Path != "" {
		base = filepath.Base(n.Path)
	}
	label := fmt.Sprintf("%s (job: %s)", base, e.JobID)
	if g.IsEntryPoint(e.FromID) {
		label += "  <- entry point"
	}
	return label
}

// renderTransitiveRequirements aggregates, across the entry point and everything
// reachable from it, what the whole pipeline needs: secret and variable names (declared,
// forwarded, and referenced in expressions), the union of declared permission grants, and
// the external workflows it pulls in. It answers "what does this whole chain need?"
// without walking every hop.
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
	vars := map[string]bool{}
	perms := map[string]bool{}
	externals := map[string]bool{}
	for _, nid := range append([]string{id}, reach...) {
		n := g.Nodes[nid]
		if n == nil || n.External {
			continue
		}
		collectSecretNames(n, secrets)
		collectScannedRefs(n, secrets, vars)
		collectPermissions(n, perms)
		// External references are collected from this node's outgoing edges (not from the
		// external nodes themselves) so the `@ref` pin each call site uses is preserved.
		for _, e := range g.Calls(nid) {
			if to := g.Nodes[e.ToID]; to != nil && to.External {
				externals[calleeDisplay(g, e)] = true
			}
		}
	}
	if len(secrets) == 0 && len(vars) == 0 && len(perms) == 0 && len(externals) == 0 {
		return
	}

	b.WriteString("## Transitive requirements (from full call graph)\n\n")
	if len(secrets) > 0 {
		fmt.Fprintf(b, "Secrets referenced (literal names): %s\n\n", codelist(sortedKeys(secrets)))
	}
	if len(vars) > 0 {
		fmt.Fprintf(b, "Variables referenced: %s\n\n", codelist(sortedKeys(vars)))
	}
	if len(perms) > 0 {
		fmt.Fprintf(b, "Permissions declared across the chain: %s\n\n", codelist(sortedKeys(perms)))
	}
	if len(externals) > 0 {
		fmt.Fprintf(b, "External workflows referenced: %s\n\n", codelist(sortedKeys(externals)))
	}
}

// collectScannedRefs unions the secret/variable names referenced in a workflow node's
// expressions (run:, with:, env:, if:, forwarded secrets: values) into the given sets. It
// reuses the scanner that builds each workflow's own reference inventory, so the
// transitive view can never disagree with the per-workflow sections.
func collectScannedRefs(n *callgraph.Node, secrets, vars map[string]bool) {
	if n.Workflow == nil {
		return
	}
	refs := model.ScanReferences(n.Workflow)
	for _, r := range refs.Secrets {
		secrets[r.Name] = true
	}
	for _, r := range refs.Vars {
		vars[r.Name] = true
	}
}

// collectPermissions unions a workflow node's declared permission grants -- workflow-level
// and job-level -- as "scope: level" strings, with the (OIDC) marker on id-token: write.
// The scalar forms (read-all / write-all) are included as-is; an explicit default-deny
// (permissions: {}) grants nothing and contributes nothing.
func collectPermissions(n *callgraph.Node, set map[string]bool) {
	if n.Workflow == nil {
		return
	}
	add := func(p *model.Permissions) {
		if p == nil {
			return
		}
		if p.All != "" {
			set[p.All] = true
		}
		for _, s := range p.Scopes {
			grant := s.Scope + ": " + s.Level
			if s.OIDC {
				grant += " (OIDC)"
			}
			set[grant] = true
		}
	}
	add(n.Workflow.Permissions)
	for i := range n.Workflow.Jobs {
		add(n.Workflow.Jobs[i].Permissions)
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

// treeNode is a single line in an ASCII dependency tree plus its children. Both the
// downstream call graph and the upstream "called by" chain build trees of these and hand
// them to renderTree, so there is exactly one tree-drawing implementation.
type treeNode struct {
	label    string
	children []treeNode
}

// renderTree writes an ASCII tree: the root label on its own line, then each child under
// "+-- " with "|   " / "    " continuation guides, using only plain keyboard characters.
func renderTree(b *strings.Builder, root treeNode) {
	b.WriteString(root.label)
	b.WriteString("\n")
	renderTreeChildren(b, root.children, "")
}

func renderTreeChildren(b *strings.Builder, children []treeNode, prefix string) {
	for i, c := range children {
		last := i == len(children)-1
		b.WriteString(prefix + "+-- " + c.label + "\n")
		guide := prefix + "|   "
		if last {
			guide = prefix + "    "
		}
		renderTreeChildren(b, c.children, guide)
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

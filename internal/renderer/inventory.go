package renderer

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/smol-utils/actiondoc/internal/callgraph"
	"github.com/smol-utils/actiondoc/internal/model"
)

// This file renders the document-level configuration inventory: a single place that answers
// "what secrets and variables must I configure, and what permissions does CI need, to run
// this repository?" without unioning every per-workflow section by hand. It is a read-only
// consumer of the parsed sources and the built call graph; it never mutates either. It lives
// outside RenderMarkdownGraph (which renders one source at a time and feeds the renderer
// goldens) so the cross-document aggregation stays off that single-source path.

// Document-level inventory headings. They deliberately differ from the per-workflow
// "## Permissions" and "## Referenced secrets and variables" headings so the two never
// collide on a shared GitHub anchor slug within the same rendered document.
const (
	inventorySecretsVarsHeading = "## Secrets and variables used across this repository"
	inventoryPermissionsHeading = "## Permissions across this repository"
)

// RenderDocumentInventory builds the repo-wide secrets/variables and permissions inventory
// across every workflow source. Secrets and variables are the ACTUAL ${{ secrets.X }} /
// ${{ vars.Y }} references found by the per-workflow scanner (so this view can never disagree
// with the per-workflow "Referenced secrets and variables" sections), attributed to the
// workflows that use them. GITHUB_TOKEN is excluded (GitHub auto-provides it; a maintainer
// never configures it). Composite actions are excluded from the secrets/variables tables: the
// scanner is workflow-only, and an action's @secret declarations are a contract, not config.
// The whole section is suppressed when there are no secrets, no variables, and no permissions;
// each sub-table is suppressed independently when its own set is empty.
func RenderDocumentInventory(sources []callgraph.Source, g *callgraph.Graph) string {
	secrets, vars := aggregateReferences(sources, g)
	perms := aggregatePermissions(sources)
	if len(secrets) == 0 && len(vars) == 0 && len(perms) == 0 {
		return ""
	}

	var b strings.Builder
	if len(secrets) > 0 || len(vars) > 0 {
		b.WriteString(inventorySecretsVarsHeading + "\n\n")
		writeInventoryRefTable(&b, "Secrets", secrets)
		writeInventoryRefTable(&b, "Variables", vars)
	}
	if len(perms) > 0 {
		b.WriteString(inventoryPermissionsHeading + "\n\n")
		writeInventoryPermissionsTable(&b, perms)
	}
	return b.String()
}

// inventoryRef is one secret/variable name plus the workflows that reference it, deduplicated
// per workflow and ordered deterministically.
type inventoryRef struct {
	name   string
	usedBy []inventoryWorkflow
}

// inventoryWorkflow is a referencing workflow's display name and document anchor, used to
// render a cross-link from the inventory to that workflow's rendered section.
type inventoryWorkflow struct {
	name   string
	anchor string
}

// aggregateReferences unions the secret and variable references of every workflow source,
// attributing each name to the set of workflows that reference it. GITHUB_TOKEN is dropped
// from the secrets set. Names are sorted alphabetically; each name's workflow list is sorted
// by anchor (then display name) so the rendering never depends on map iteration order.
func aggregateReferences(sources []callgraph.Source, g *callgraph.Graph) (secrets, vars []inventoryRef) {
	// name -> set of node ids (one per referencing workflow), so a name used at several sites
	// in the same workflow is attributed to that workflow exactly once.
	secByName := map[string]map[string]bool{}
	varByName := map[string]map[string]bool{}
	add := func(m map[string]map[string]bool, name, id string) {
		if m[name] == nil {
			m[name] = map[string]bool{}
		}
		m[name][id] = true
	}
	for _, s := range sources {
		if s.Workflow == nil {
			continue
		}
		refs := model.ScanReferences(s.Workflow)
		for _, r := range refs.Secrets {
			if r.Name == "GITHUB_TOKEN" {
				continue
			}
			add(secByName, r.Name, s.Path)
		}
		for _, r := range refs.Vars {
			add(varByName, r.Name, s.Path)
		}
	}
	return buildInventoryRefs(secByName, g), buildInventoryRefs(varByName, g)
}

// buildInventoryRefs turns a name -> set(node id) map into a sorted slice of inventoryRef,
// resolving each node id to its display name and anchor via the graph.
func buildInventoryRefs(byName map[string]map[string]bool, g *callgraph.Graph) []inventoryRef {
	names := make([]string, 0, len(byName))
	for name := range byName {
		names = append(names, name)
	}
	sort.Strings(names)

	ambiguous := ambiguousWorkflowNames(g)
	out := make([]inventoryRef, 0, len(names))
	for _, name := range names {
		var wfs []inventoryWorkflow
		for id := range byName[name] {
			wf := inventoryWorkflow{name: id, anchor: anchor(id)}
			if n := g.Nodes[id]; n != nil {
				// Disambiguate the visible label when another workflow shares this
				// display name: bare names would render identical link text against
				// distinct anchors (e.g. #model-jobs vs #model-jobs-1), so append the
				// filename the way the table of contents does.
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
		out = append(out, inventoryRef{name: name, usedBy: wfs})
	}
	return out
}

// ambiguousWorkflowNames returns the set of workflow display names shared by more than one
// workflow node in the graph. Such names need a filename suffix in cross-links so their
// otherwise-identical link text stays distinguishable. Action nodes and external (cross-repo)
// nodes are excluded: neither appears in the secrets/variables "Used by" lists this set
// guards, so counting them could only invent a false collision that adds a needless filename
// suffix to a local workflow's label.
func ambiguousWorkflowNames(g *callgraph.Graph) map[string]bool {
	counts := map[string]int{}
	for _, n := range g.Nodes {
		if n.IsAction || n.External {
			continue
		}
		counts[n.Name]++
	}
	out := map[string]bool{}
	for name, c := range counts {
		if c > 1 {
			out[name] = true
		}
	}
	return out
}

// writeInventoryRefTable renders one secrets/variables table: each row is a name and the
// comma-separated cross-links to the workflows that reference it. Suppressed when empty.
func writeInventoryRefTable(b *strings.Builder, label string, refs []inventoryRef) {
	if len(refs) == 0 {
		return
	}
	fmt.Fprintf(b, "**%s:**\n\n", label)
	b.WriteString("| Name | Used by |\n")
	b.WriteString("|------|---------|\n")
	for _, r := range refs {
		links := make([]string, len(r.usedBy))
		for i, wf := range r.usedBy {
			links[i] = fmt.Sprintf("[%s](#%s)", mdLinkLabel(wf.name), wf.anchor)
		}
		fmt.Fprintf(b, "| `%s` | %s |\n", escapeCellCode(r.name), strings.Join(links, ", "))
	}
	b.WriteString("\n")
}

// inventoryPerm is the unioned grant for one permission scope across the repository: the
// effective (maximum) level, an OIDC marker, and any conflicting levels seen elsewhere.
type inventoryPerm struct {
	scope     string
	effective string
	oidc      bool
	conflicts []string // other levels seen for this scope, sorted, excluding the effective one
}

// permRank orders permission levels so the effective grant is the most permissive one seen.
// write outranks read outranks none; unknown levels sort below the known ones but are still
// surfaced rather than dropped.
func permRank(level string) int {
	switch level {
	case "write":
		return 3
	case "read":
		return 2
	case "none":
		return 1
	default:
		return 0
	}
}

// aggregatePermissions unions the declared permission grants across every workflow source --
// workflow-level and each job-level block -- into one scope -> set(levels) view. A scope seen
// at different levels (e.g. read in
// one workflow, write in another) keeps its effective (maximum) level plus a conflict note.
// The (OIDC) marker rides on id-token: write. Scalar read-all/write-all grants are surfaced as
// their own rows; an explicit permissions: {} (default-deny) grants nothing and contributes
// nothing.
func aggregatePermissions(sources []callgraph.Source) []inventoryPerm {
	type agg struct {
		levels map[string]bool
		oidc   bool
	}
	scopes := map[string]*agg{}
	scalars := map[string]bool{} // read-all / write-all
	addBlock := func(p *model.Permissions) {
		if p == nil || p.DefaultDeny {
			return
		}
		if p.All != "" {
			scalars[p.All] = true
			return
		}
		for _, s := range p.Scopes {
			a := scopes[s.Scope]
			if a == nil {
				a = &agg{levels: map[string]bool{}}
				scopes[s.Scope] = a
			}
			a.levels[s.Level] = true
			if s.OIDC {
				a.oidc = true
			}
		}
	}
	for _, s := range sources {
		if s.Workflow == nil {
			continue
		}
		addBlock(s.Workflow.Permissions)
		for i := range s.Workflow.Jobs {
			addBlock(s.Workflow.Jobs[i].Permissions)
		}
	}

	var out []inventoryPerm
	// Scalar repo-wide grants first, under a sentinel scope so they share the one table.
	for _, scalar := range sortedKeys(scalars) {
		out = append(out, inventoryPerm{scope: "(all scopes)", effective: scalar})
	}

	scopeNames := make([]string, 0, len(scopes))
	for name := range scopes {
		scopeNames = append(scopeNames, name)
	}
	sort.Strings(scopeNames)
	for _, name := range scopeNames {
		a := scopes[name]
		levels := sortedKeys(a.levels)
		effective := levels[0]
		for _, l := range levels {
			if permRank(l) > permRank(effective) {
				effective = l
			}
		}
		var conflicts []string
		for _, l := range levels {
			if l != effective {
				conflicts = append(conflicts, l)
			}
		}
		out = append(out, inventoryPerm{scope: name, effective: effective, oidc: a.oidc, conflicts: conflicts})
	}
	return out
}

// writeInventoryPermissionsTable renders the unioned permissions as a scope/level table. A
// scope granted different levels in different workflows shows its effective level with a short
// note naming the other levels seen.
func writeInventoryPermissionsTable(b *strings.Builder, perms []inventoryPerm) {
	b.WriteString("| Scope | Level |\n")
	b.WriteString("|-------|-------|\n")
	for _, p := range perms {
		level := "`" + p.effective + "`"
		if p.oidc {
			level += " (OIDC)"
		}
		if len(p.conflicts) > 0 {
			level += fmt.Sprintf(" (also granted as %s elsewhere)", codelist(p.conflicts))
		}
		fmt.Fprintf(b, "| `%s` | %s |\n", escapeCellCode(p.scope), level)
	}
	b.WriteString("\n")
}

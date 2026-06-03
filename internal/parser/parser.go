package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml/ast"
	yamlparser "github.com/goccy/go-yaml/parser"
	"github.com/smol-utils/actiondoc/internal/model"
)

// firstMappingDoc returns the first YAML document whose body is a mapping, wrapping a
// bare MappingValueNode if needed. goccy/go-yaml emits a leading comment-only document
// when a file opens with a '#' comment block followed by a '---' document-start marker
// (e.g. a license header). Taking Docs[0] unconditionally then fails with
// "expected top-level mapping"; iterating skips the comment-only document.
//
// It also returns any documents skipped before the mapping. Their comments must still be
// searched for the head comment -- otherwise a description or ActionDoc tags placed in a
// comment block before `---` would be silently lost.
func firstMappingDoc(file *ast.File) (doc *ast.DocumentNode, root *ast.MappingNode, skipped []*ast.DocumentNode, ok bool) {
	for _, d := range file.Docs {
		if d == nil || d.Body == nil {
			skipped = append(skipped, d)
			continue
		}
		switch body := d.Body.(type) {
		case *ast.MappingNode:
			return d, body, skipped, true
		case *ast.MappingValueNode:
			return d, &ast.MappingNode{Values: []*ast.MappingValueNode{body}}, skipped, true
		}
		skipped = append(skipped, d)
	}
	return nil, nil, skipped, false
}

// headCommentNodes builds the node list to search for the head comment: the leading
// comment-only documents skipped before the mapping (where a pre-`---` comment block
// lands) first, then the mapping document, root, first entry, and its key.
func headCommentNodes(doc *ast.DocumentNode, root *ast.MappingNode, skipped []*ast.DocumentNode) []ast.Node {
	var nodes []ast.Node
	for _, d := range skipped {
		if d != nil {
			nodes = append(nodes, d)
		}
	}
	nodes = append(nodes, doc, root)
	if len(root.Values) > 0 {
		nodes = append(nodes, root.Values[0], root.Values[0].Key)
	}
	return nodes
}

// ParseFile parses a single workflow YAML file into the IR.
func ParseFile(path string) (*model.Workflow, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}

	file, err := yamlparser.ParseBytes(data, yamlparser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("parsing %s: %w", path, err)
	}

	if len(file.Docs) == 0 {
		return nil, fmt.Errorf("%s: no YAML documents found", path)
	}
	doc, root, skipped, ok := firstMappingDoc(file)
	if !ok {
		return nil, fmt.Errorf("%s: expected top-level mapping", path)
	}
	resolveAnchors(root)

	w := &model.Workflow{
		File: filepath.Base(path),
	}

	// Workflow-level tags: the YAML library may attach the top-of-file comment to a
	// leading comment-only document (pre-`---` block), the mapping document, root
	// mapping, first MappingValueNode, or its key.
	headComment := findComment(headCommentNodes(doc, root, skipped)...)
	w.Tags = ParseTags(headComment)
	w.Description = w.Tags.Desc
	if w.Description == "" {
		w.Description = implicitDescription(headComment, w.Tags)
	}

	// Walk top-level keys.
	for _, mv := range root.Values {
		keyStr := mapKeyString(mv.Key)
		switch keyStr {
		case "name":
			w.Name = nameString(mv.Value)
		case "on", "true":
			w.On = parseTriggers(mv.Value)
			w.Triggers = parseTriggerSurface(mv.Value)
		case "jobs":
			w.Jobs = parseJobs(mv.Value)
		case "permissions":
			w.Permissions = parsePermissions(mv.Value)
		case "env":
			w.Env = parseKVMap(mv.Value)
		case "concurrency":
			w.Concurrency = parseConcurrency(mv.Value)
		case "defaults":
			w.Defaults = parseDefaults(mv.Value)
		}
	}

	if w.Name == "" {
		w.Name = w.File
	}

	return w, nil
}

// commentString extracts the text from a CommentGroupNode.
func commentString(cg *ast.CommentGroupNode) string {
	if cg == nil {
		return ""
	}
	return cg.String()
}

// findComment returns the first non-empty comment string from the given nodes. A node
// may carry its comment via GetComment(), or BE a comment group itself (goccy represents
// a pre-`---` comment block as a document whose body is a *ast.CommentGroupNode).
func findComment(nodes ...ast.Node) string {
	for _, n := range nodes {
		if n == nil {
			continue
		}
		if cg, ok := n.(*ast.CommentGroupNode); ok {
			if s := commentString(cg); s != "" {
				return s
			}
		}
		if doc, ok := n.(*ast.DocumentNode); ok {
			if cg, ok := doc.Body.(*ast.CommentGroupNode); ok {
				if s := commentString(cg); s != "" {
					return s
				}
			}
		}
		if s := commentString(n.GetComment()); s != "" {
			return s
		}
	}
	return ""
}

// mapKeyString extracts the string value from a mapping key node.
func mapKeyString(node ast.MapKeyNode) string {
	switch n := node.(type) {
	case *ast.StringNode:
		return n.Value
	case *ast.BoolNode:
		// YAML 1.1: "on" may be parsed as a boolean true.
		return n.GetToken().Value
	default:
		return node.String()
	}
}

// nodeString extracts a plain string value from a node.
func nodeString(node ast.Node) string {
	if node == nil {
		return ""
	}
	switch n := node.(type) {
	case *ast.StringNode:
		return n.Value
	case *ast.LiteralNode:
		return n.Value.Value
	case *ast.BoolNode:
		return n.GetToken().Value
	case *ast.IntegerNode:
		return n.GetToken().Value
	default:
		return node.String()
	}
}

// nameString reads a display-name value (workflow/job/step/action `name:`) as a single
// line: newlines from block scalars collapse to spaces and runs of whitespace collapse to
// one space. A display name is one line by definition; embedded newlines are YAML
// formatting accidents that would otherwise break the Markdown structures built around
// names (bold step titles, ASCII tree labels, heading anchors).
func nameString(node ast.Node) string {
	return strings.Join(strings.Fields(nodeString(node)), " ")
}

// parseTriggers extracts event names from the on: node.
// Handles: scalar ("push"), sequence ([push, pull_request]),
// and mapping (push: { branches: [...] }).
func parseTriggers(node ast.Node) []string {
	switch n := node.(type) {
	case *ast.StringNode:
		return []string{n.Value}
	case *ast.BoolNode:
		return []string{n.GetToken().Value}
	case *ast.SequenceNode:
		var triggers []string
		for _, item := range n.Values {
			triggers = append(triggers, nodeString(item))
		}
		return triggers
	case *ast.MappingNode:
		var triggers []string
		for _, mv := range n.Values {
			triggers = append(triggers, mapKeyString(mv.Key))
		}
		return triggers
	case *ast.MappingValueNode:
		return []string{mapKeyString(n.Key)}
	}
	return nil
}

// parseJobs extracts jobs from the jobs: mapping node.
func parseJobs(node ast.Node) []model.Job {
	mapping := toMapping(node)
	if mapping == nil {
		return nil
	}

	var jobs []model.Job
	for _, mv := range mapping.Values {
		job := model.Job{
			ID: mapKeyString(mv.Key),
		}

		job.Tags = ParseTags(findComment(mv, mv.Key))
		job.Description = job.Tags.Desc

		if jobMapping := toMapping(mv.Value); jobMapping != nil {
			parseJobFields(&job, jobMapping)
		}

		if job.Name == "" {
			job.Name = job.ID
		}

		jobs = append(jobs, job)
	}
	return jobs
}

// parseJobFields fills job fields from the job's mapping node.
func parseJobFields(job *model.Job, mapping *ast.MappingNode) {
	for _, mv := range mapping.Values {
		keyStr := mapKeyString(mv.Key)
		switch keyStr {
		case "name":
			job.Name = nameString(mv.Value)
		case "runs-on":
			job.RunsOn = parseRunsOn(mv.Value)
		case "needs":
			job.Needs = parseStringOrSequence(mv.Value)
		case "if":
			job.If = nodeString(mv.Value)
		case "strategy":
			job.Matrix = parseMatrix(mv.Value)
		case "steps":
			job.Steps = parseSteps(mv.Value)
		// Reusable-workflow caller jobs use uses:/with:/secrets: instead of steps.
		case "uses":
			job.Uses = nodeString(mv.Value)
		case "with":
			job.With = parseKVMap(mv.Value)
		case "secrets":
			// `secrets: inherit` (scalar) vs an explicit mapping of forwarded secrets.
			if strings.TrimSpace(nodeString(mv.Value)) == "inherit" {
				job.SecretsInherit = true
			} else {
				job.Secrets = parseKVMap(mv.Value)
			}
		case "permissions":
			job.Permissions = parsePermissions(mv.Value)
		case "env":
			job.Env = parseKVMap(mv.Value)
		case "concurrency":
			job.Concurrency = parseConcurrency(mv.Value)
		case "defaults":
			job.Defaults = parseDefaults(mv.Value)
		case "environment":
			job.Environment = parseEnvironment(mv.Value)
		}
	}
}

// parseKVMap parses a YAML mapping into an ordered slice of key/value pairs, preserving
// source order. Used for reusable-workflow `with:`/`secrets:` forwarding maps.
func parseKVMap(node ast.Node) []model.KV {
	mapping := toMapping(node)
	if mapping == nil {
		return nil
	}
	var out []model.KV
	for _, mv := range mapping.Values {
		out = append(out, model.KV{
			Key:   mapKeyString(mv.Key),
			Value: nodeString(mv.Value),
		})
	}
	return out
}

// parseSteps extracts steps from the steps: sequence node.
func parseSteps(node ast.Node) []model.Step {
	seq, ok := node.(*ast.SequenceNode)
	if !ok {
		return nil
	}

	var steps []model.Step
	for i, item := range seq.Values {
		step := model.Step{}

		itemMapping := toMapping(item)
		if itemMapping == nil {
			continue
		}

		step.Tags = ParseTags(stepComment(seq, i, item, itemMapping))
		step.Description = step.Tags.Desc

		for _, mv := range itemMapping.Values {
			keyStr := mapKeyString(mv.Key)
			switch keyStr {
			case "name":
				step.Name = nameString(mv.Value)
			case "id":
				step.ID = nodeString(mv.Value)
			case "uses":
				step.Uses = nodeString(mv.Value)
				step.UsesVersion = versionComment(TrailingComment(mv))
			case "run":
				step.Run = nodeString(mv.Value)
			case "if":
				step.If = nodeString(mv.Value)
			case "with":
				step.With = parseKVMap(mv.Value)
			case "env":
				step.Env = parseKVMap(mv.Value)
			case "continue-on-error":
				v := strings.TrimSpace(nodeString(mv.Value))
				if strings.EqualFold(v, "true") {
					step.ContinueOnError = true
				} else if strings.Contains(v, "${{") {
					// An expression-valued continue-on-error (e.g. matrix-driven) is still
					// failure-tolerant; keep the raw expression so the renderer can show it.
					step.ContinueOnErrorExpr = v
				}
			}
		}
		steps = append(steps, step)
	}
	return steps
}

// stepComment finds the comment for a step. The YAML library may place it in:
// the parallel ValueHeadComments array, the item node, or the first key inside.
func stepComment(seq *ast.SequenceNode, i int, item ast.Node, mapping *ast.MappingNode) string {
	if i < len(seq.ValueHeadComments) && seq.ValueHeadComments[i] != nil {
		if s := commentString(seq.ValueHeadComments[i]); s != "" {
			return s
		}
	}
	if s := commentString(item.GetComment()); s != "" {
		return s
	}
	if len(mapping.Values) > 0 {
		first := mapping.Values[0]
		return findComment(first, first.Key)
	}
	return ""
}

// parseStringOrSequence handles YAML values that can be a string or a sequence of strings.
func parseStringOrSequence(node ast.Node) []string {
	switch n := node.(type) {
	case *ast.StringNode:
		return []string{n.Value}
	case *ast.SequenceNode:
		var vals []string
		for _, item := range n.Values {
			vals = append(vals, nodeString(item))
		}
		return vals
	}
	return nil
}

// toMapping converts a node to a MappingNode, handling single MappingValueNode cases.
func toMapping(node ast.Node) *ast.MappingNode {
	switch n := node.(type) {
	case *ast.MappingNode:
		return n
	case *ast.MappingValueNode:
		return &ast.MappingNode{Values: []*ast.MappingValueNode{n}}
	}
	return nil
}

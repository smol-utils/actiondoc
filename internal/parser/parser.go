package parser

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml/ast"
	yamlparser "github.com/goccy/go-yaml/parser"
	"github.com/smol-utils/actiondoc/internal/model"
)

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
	doc := file.Docs[0]

	root, ok := doc.Body.(*ast.MappingNode)
	if !ok {
		// Could be a single MappingValueNode; wrap it.
		if mv, ok := doc.Body.(*ast.MappingValueNode); ok {
			root = &ast.MappingNode{Values: []*ast.MappingValueNode{mv}}
		} else {
			return nil, fmt.Errorf("%s: expected top-level mapping", path)
		}
	}

	w := &model.Workflow{
		File: filepath.Base(path),
	}

	// Workflow-level tags: the YAML library may attach the top-of-file comment
	// to the document, root mapping, first MappingValueNode, or its key.
	var firstMV, firstKey ast.Node
	if len(root.Values) > 0 {
		firstMV = root.Values[0]
		firstKey = root.Values[0].Key
	}
	headComment := findComment(doc, root, firstMV, firstKey)
	w.Tags = ParseTags(headComment)
	w.Description = w.Tags.Desc

	// Walk top-level keys.
	for _, mv := range root.Values {
		keyStr := mapKeyString(mv.Key)
		switch keyStr {
		case "name":
			w.Name = nodeString(mv.Value)
		case "on", "true":
			w.On = parseTriggers(mv.Value)
		case "jobs":
			w.Jobs = parseJobs(mv.Value)
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

// findComment returns the first non-empty comment string from the given nodes.
func findComment(nodes ...ast.Node) string {
	for _, n := range nodes {
		if n == nil {
			continue
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
			job.Name = nodeString(mv.Value)
		case "runs-on":
			job.RunsOn = nodeString(mv.Value)
		case "needs":
			job.Needs = parseStringOrSequence(mv.Value)
		case "if":
			job.If = nodeString(mv.Value)
		case "steps":
			job.Steps = parseSteps(mv.Value)
		}
	}
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
				step.Name = nodeString(mv.Value)
			case "id":
				step.ID = nodeString(mv.Value)
			case "uses":
				step.Uses = nodeString(mv.Value)
			case "run":
				step.Run = nodeString(mv.Value)
			case "if":
				step.If = nodeString(mv.Value)
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

// ParseActionFile parses a GitHub Action metadata file (action.yml) into the data model.
func ParseActionFile(path string) (*model.Action, error) {
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
	doc := file.Docs[0]

	root, ok := doc.Body.(*ast.MappingNode)
	if !ok {
		if mv, ok := doc.Body.(*ast.MappingValueNode); ok {
			root = &ast.MappingNode{Values: []*ast.MappingValueNode{mv}}
		} else {
			return nil, fmt.Errorf("%s: expected top-level mapping", path)
		}
	}

	a := &model.Action{
		File: filepath.Base(path),
	}

	// Action-level tags from top-of-file comments.
	var firstMV, firstKey ast.Node
	if len(root.Values) > 0 {
		firstMV = root.Values[0]
		firstKey = root.Values[0].Key
	}
	a.Tags = ParseTags(findComment(doc, root, firstMV, firstKey))

	for _, mv := range root.Values {
		keyStr := mapKeyString(mv.Key)
		switch keyStr {
		case "name":
			a.Name = nodeString(mv.Value)
		case "description":
			a.Description = nodeString(mv.Value)
		case "inputs":
			a.Inputs = parseActionInputs(mv.Value)
		case "outputs":
			a.Outputs = parseActionOutputs(mv.Value)
		case "runs":
			a.Runs = parseActionRuns(mv.Value)
		case "branding":
			a.Branding = parseBranding(mv.Value)
		}
	}

	// Fall back to @desc if no native description.
	if a.Description == "" {
		a.Description = a.Tags.Desc
	}

	if a.Name == "" {
		a.Name = a.File
	}

	return a, nil
}

func parseActionInputs(node ast.Node) []model.ActionInput {
	mapping := toMapping(node)
	if mapping == nil {
		return nil
	}

	var inputs []model.ActionInput
	for _, mv := range mapping.Values {
		input := model.ActionInput{
			Name: mapKeyString(mv.Key),
		}

		if inputMapping := toMapping(mv.Value); inputMapping != nil {
			for _, field := range inputMapping.Values {
				switch mapKeyString(field.Key) {
				case "description":
					input.Description = nodeString(field.Value)
				case "required":
					input.Required = nodeString(field.Value) == "true"
				case "default":
					input.Default = nodeString(field.Value)
				}
			}
		}

		inputs = append(inputs, input)
	}
	return inputs
}

func parseActionOutputs(node ast.Node) []model.ActionOutput {
	mapping := toMapping(node)
	if mapping == nil {
		return nil
	}

	var outputs []model.ActionOutput
	for _, mv := range mapping.Values {
		output := model.ActionOutput{
			Name: mapKeyString(mv.Key),
		}

		if outputMapping := toMapping(mv.Value); outputMapping != nil {
			for _, field := range outputMapping.Values {
				if mapKeyString(field.Key) == "description" {
					output.Description = nodeString(field.Value)
				}
			}
		}

		outputs = append(outputs, output)
	}
	return outputs
}

func parseActionRuns(node ast.Node) model.ActionRuns {
	var runs model.ActionRuns
	mapping := toMapping(node)
	if mapping == nil {
		return runs
	}

	for _, mv := range mapping.Values {
		switch mapKeyString(mv.Key) {
		case "using":
			runs.Using = nodeString(mv.Value)
		case "main":
			runs.Main = nodeString(mv.Value)
		case "image":
			runs.Image = nodeString(mv.Value)
		}
	}
	return runs
}

func parseBranding(node ast.Node) *model.Branding {
	mapping := toMapping(node)
	if mapping == nil {
		return nil
	}

	b := &model.Branding{}
	for _, mv := range mapping.Values {
		switch mapKeyString(mv.Key) {
		case "icon":
			b.Icon = nodeString(mv.Value)
		case "color":
			b.Color = nodeString(mv.Value)
		}
	}
	return b
}

package parser

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/smol-utils/actiondoc/internal/model"
	"gopkg.in/yaml.v3"
)

// ParseFile parses a single workflow YAML file into the IR.
func ParseFile(path string) (*model.Workflow, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}

	var doc yaml.Node
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("parsing %s: %w", path, err)
	}

	root := documentRoot(&doc)
	if root == nil || root.Kind != yaml.MappingNode {
		return nil, fmt.Errorf("%s: expected top-level mapping", path)
	}

	w := &model.Workflow{
		File: filepath.Base(path),
	}

	// Workflow-level tags from the top-of-file comment.
	headComment := firstComment(&doc, root)
	w.Tags = ParseTags(headComment)
	w.Description = w.Tags.Desc

	// Walk top-level keys (Content is [key, value, key, value, ...]).
	for i := 0; i+1 < len(root.Content); i += 2 {
		key := root.Content[i]
		val := root.Content[i+1]
		switch key.Value {
		case "name":
			w.Name = nodeString(val)
		case "on", "true":
			w.On = parseTriggers(val)
		case "jobs":
			w.Jobs = parseJobs(val)
		}
	}

	if w.Name == "" {
		w.Name = w.File
	}

	return w, nil
}

// documentRoot unwraps a DocumentNode to its single content node.
func documentRoot(doc *yaml.Node) *yaml.Node {
	if doc.Kind == yaml.DocumentNode {
		if len(doc.Content) > 0 {
			return doc.Content[0]
		}
		return nil
	}
	return doc
}

// firstComment returns the first non-empty head comment among the given nodes.
func firstComment(nodes ...*yaml.Node) string {
	for _, n := range nodes {
		if n == nil {
			continue
		}
		if n.HeadComment != "" {
			return n.HeadComment
		}
	}
	return ""
}

// nodeString extracts a plain scalar value from a node.
func nodeString(node *yaml.Node) string {
	if node == nil {
		return ""
	}
	return node.Value
}

// parseTriggers extracts event names from the on: node.
// Handles: scalar ("push"), sequence ([push, pull_request]),
// and mapping (push: { branches: [...] }).
func parseTriggers(node *yaml.Node) []string {
	switch node.Kind {
	case yaml.ScalarNode:
		return []string{node.Value}
	case yaml.SequenceNode:
		var triggers []string
		for _, item := range node.Content {
			triggers = append(triggers, item.Value)
		}
		return triggers
	case yaml.MappingNode:
		var triggers []string
		for i := 0; i+1 < len(node.Content); i += 2 {
			triggers = append(triggers, node.Content[i].Value)
		}
		return triggers
	}
	return nil
}

// parseJobs extracts jobs from the jobs: mapping node.
func parseJobs(node *yaml.Node) []model.Job {
	if node.Kind != yaml.MappingNode {
		return nil
	}

	var jobs []model.Job
	for i := 0; i+1 < len(node.Content); i += 2 {
		key := node.Content[i]
		val := node.Content[i+1]

		job := model.Job{ID: key.Value}
		job.Tags = ParseTags(firstComment(key))
		job.Description = job.Tags.Desc

		if val.Kind == yaml.MappingNode {
			parseJobFields(&job, val)
		}

		if job.Name == "" {
			job.Name = job.ID
		}

		jobs = append(jobs, job)
	}
	return jobs
}

// parseJobFields fills job fields from the job's mapping node.
func parseJobFields(job *model.Job, node *yaml.Node) {
	for i := 0; i+1 < len(node.Content); i += 2 {
		key := node.Content[i]
		val := node.Content[i+1]
		switch key.Value {
		case "name":
			job.Name = nodeString(val)
		case "runs-on":
			job.RunsOn = nodeString(val)
		case "needs":
			job.Needs = parseStringOrSequence(val)
		case "if":
			job.If = nodeString(val)
		case "steps":
			job.Steps = parseSteps(val)
		}
	}
}

// parseSteps extracts steps from the steps: sequence node.
func parseSteps(node *yaml.Node) []model.Step {
	if node.Kind != yaml.SequenceNode {
		return nil
	}

	var steps []model.Step
	for _, item := range node.Content {
		if item.Kind != yaml.MappingNode {
			continue
		}

		step := model.Step{}
		step.Tags = ParseTags(stepComment(item))
		step.Description = step.Tags.Desc

		for i := 0; i+1 < len(item.Content); i += 2 {
			key := item.Content[i]
			val := item.Content[i+1]
			switch key.Value {
			case "name":
				step.Name = nodeString(val)
			case "id":
				step.ID = nodeString(val)
			case "uses":
				step.Uses = nodeString(val)
			case "run":
				step.Run = nodeString(val)
			case "if":
				step.If = nodeString(val)
			}
		}
		steps = append(steps, step)
	}
	return steps
}

// stepComment finds the comment for a step. yaml.v3 may attach the head comment
// to the step's mapping node or to the first key inside it.
func stepComment(item *yaml.Node) string {
	if item.HeadComment != "" {
		return item.HeadComment
	}
	if len(item.Content) > 0 {
		return item.Content[0].HeadComment
	}
	return ""
}

// parseStringOrSequence handles YAML values that can be a string or a sequence of strings.
func parseStringOrSequence(node *yaml.Node) []string {
	switch node.Kind {
	case yaml.ScalarNode:
		return []string{node.Value}
	case yaml.SequenceNode:
		var vals []string
		for _, item := range node.Content {
			vals = append(vals, item.Value)
		}
		return vals
	}
	return nil
}

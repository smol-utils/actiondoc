package parser

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml/ast"
	yamlparser "github.com/goccy/go-yaml/parser"
	"github.com/smol-utils/actiondoc/internal/model"
)

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
	doc, root, skipped, ok := firstMappingDoc(file)
	if !ok {
		return nil, fmt.Errorf("%s: expected top-level mapping", path)
	}

	a := &model.Action{
		File: filepath.Base(path),
	}

	// Action-level tags from top-of-file comments, including a pre-`---` comment block
	// that goccy/go-yaml places in a leading comment-only document.
	a.Tags = ParseTags(findComment(headCommentNodes(doc, root, skipped)...))

	for _, mv := range root.Values {
		keyStr := mapKeyString(mv.Key)
		switch keyStr {
		case "name":
			a.Name = nameString(mv.Value)
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

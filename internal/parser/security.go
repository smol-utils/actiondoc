package parser

import (
	"strings"

	"github.com/goccy/go-yaml/ast"
	"github.com/smol-utils/actiondoc/internal/model"
)

// parsePermissions parses a permissions: block at the workflow or job level. It handles
// all three shapes: a scalar grant (read-all/write-all), an explicit empty block
// (permissions: {}, which grants nothing -- a default-deny posture), and a mapping of
// scope->level grants. Trailing-comment rationale and the id-token: write OIDC marker are
// captured per grant.
func parsePermissions(node ast.Node) *model.Permissions {
	if node == nil {
		return nil
	}

	// Scalar form: permissions: read-all / write-all.
	if scalar := scalarString(node); scalar != "" {
		return &model.Permissions{All: scalar}
	}

	mapping := toMapping(node)
	if mapping == nil {
		return nil
	}

	// Explicit empty block: permissions: {} -- grants nothing.
	if len(mapping.Values) == 0 {
		return &model.Permissions{DefaultDeny: true}
	}

	p := &model.Permissions{}
	for _, mv := range mapping.Values {
		scope := mapKeyString(mv.Key)
		level := nodeString(mv.Value)
		p.Scopes = append(p.Scopes, model.Permission{
			Scope:     scope,
			Level:     level,
			Rationale: TrailingComment(mv),
			OIDC:      scope == "id-token" && strings.EqualFold(strings.TrimSpace(level), "write"),
		})
	}
	return p
}

// parseEnvironment parses a job-level environment: binding. It accepts both the scalar
// form (environment: production) and the map form
// (environment: { name: production, url: ... }). Returns nil when no name resolves.
func parseEnvironment(node ast.Node) *model.Environment {
	if node == nil {
		return nil
	}
	if scalar := scalarString(node); scalar != "" {
		return &model.Environment{Name: scalar}
	}
	mapping := toMapping(node)
	if mapping == nil {
		return nil
	}
	env := &model.Environment{}
	for _, mv := range mapping.Values {
		switch mapKeyString(mv.Key) {
		case "name":
			env.Name = nodeString(mv.Value)
		case "url":
			env.URL = nodeString(mv.Value)
		}
	}
	if env.Name == "" {
		return nil
	}
	return env
}

// scalarString returns the string value of a scalar node (string/bool/int), or "" if the
// node is a mapping or sequence. Used to distinguish scalar grants from block forms.
func scalarString(node ast.Node) string {
	switch node.(type) {
	case *ast.StringNode, *ast.BoolNode, *ast.IntegerNode, *ast.LiteralNode:
		return strings.TrimSpace(nodeString(node))
	}
	return ""
}

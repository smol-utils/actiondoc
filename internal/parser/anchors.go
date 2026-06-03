package parser

import "github.com/goccy/go-yaml/ast"

// This file enforces one of the parser's field-reading rules: every field reader must
// handle the three forms a YAML Actions value can take -- a literal, a ${{ }} expression
// (kept verbatim as a string), and YAML indirection (anchors/aliases). The first two are
// each reader's responsibility; indirection is resolved document-wide here, before any
// reader runs, so no reader ever sees a raw &name or *name token.

// resolveAnchors replaces YAML anchor/alias indirection throughout a document: anchor
// definitions are unwrapped to their underlying values, and aliases are replaced by the
// nodes their anchors define. Anchors are a YAML authoring detail ("define once, reuse");
// the parsed model must hold the resolved values. YAML requires an anchor to be defined
// before its aliases, so a single document-order walk resolves references correctly; an
// alias whose anchor is unknown (including a self-reference inside its own definition) is
// left as-is rather than guessed at.
func resolveAnchors(root ast.Node) {
	resolveAnchorNode(root, map[string]ast.Node{})
}

// resolveAnchorNode walks node in document order, recording anchors and substituting
// aliases, and returns the node that should take node's place in the tree (the unwrapped
// value for anchors, the anchored node for known aliases, node itself otherwise).
func resolveAnchorNode(node ast.Node, anchors map[string]ast.Node) ast.Node {
	switch n := node.(type) {
	case *ast.AnchorNode:
		resolved := resolveAnchorNode(n.Value, anchors)
		if n.Name != nil {
			anchors[n.Name.String()] = resolved
		}
		return resolved
	case *ast.AliasNode:
		if n.Value != nil {
			if target, ok := anchors[n.Value.String()]; ok {
				return target
			}
		}
		return n
	case *ast.MappingNode:
		for _, mv := range n.Values {
			resolveAnchorNode(mv, anchors)
		}
	case *ast.MappingValueNode:
		n.Value = resolveAnchorNode(n.Value, anchors)
	case *ast.SequenceNode:
		for i, v := range n.Values {
			n.Values[i] = resolveAnchorNode(v, anchors)
		}
	}
	return node
}

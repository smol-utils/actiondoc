package parser

import (
	"strings"

	"github.com/goccy/go-yaml/ast"
	"github.com/smol-utils/actiondoc/internal/model"
)

// parseRunsOn normalizes a runs-on value to a readable string. It accepts a scalar
// (`ubuntu-latest`), a list (`[self-hosted, linux, x64]` -> "self-hosted, linux, x64"),
// or a map (matrix runner selection), rather than rendering raw JSON-ish text.
func parseRunsOn(node ast.Node) string {
	switch n := node.(type) {
	case *ast.SequenceNode:
		var vals []string
		for _, item := range n.Values {
			vals = append(vals, nodeString(item))
		}
		return strings.Join(vals, ", ")
	case *ast.MappingNode:
		return runsOnMap(n)
	case *ast.MappingValueNode:
		return runsOnMap(&ast.MappingNode{Values: []*ast.MappingValueNode{n}})
	default:
		return nodeString(node)
	}
}

func runsOnMap(m *ast.MappingNode) string {
	var parts []string
	for _, mv := range m.Values {
		parts = append(parts, mapKeyString(mv.Key)+": "+parseRunsOn(mv.Value))
	}
	return strings.Join(parts, ", ")
}

// versionComment returns the trailing comment of a SHA-pinned `uses:` ref when it looks
// like a version (e.g. "v4" or "4.1.0"), used to show a readable version next to the pin.
// Returns "" for non-version comments.
func versionComment(c string) string {
	if c == "" {
		return ""
	}
	tok := strings.Fields(c)[0]
	digits := tok
	if len(digits) > 0 && (digits[0] == 'v' || digits[0] == 'V') {
		digits = digits[1:]
	}
	if len(digits) > 0 && digits[0] >= '0' && digits[0] <= '9' {
		return tok
	}
	return ""
}

// parseMatrix extracts the statically-resolvable axes of a job's strategy.matrix. It
// returns nil (render names verbatim) for matrices that use include/exclude, a non-list
// axis, or a dynamic fromJSON source, since those can't be resolved without runtime data.
func parseMatrix(node ast.Node) []model.MatrixAxis {
	strat := toMapping(node)
	if strat == nil {
		return nil
	}
	var matrixNode ast.Node
	for _, mv := range strat.Values {
		if mapKeyString(mv.Key) == "matrix" {
			matrixNode = mv.Value
		}
	}
	m := toMapping(matrixNode)
	if m == nil {
		return nil // scalar/dynamic matrix (e.g. fromJSON) -> not statically resolvable
	}
	for _, mv := range m.Values {
		if k := mapKeyString(mv.Key); k == "include" || k == "exclude" {
			return nil // include/exclude alters the product set; render names verbatim
		}
	}

	acc := newAxisAccum()
	for _, mv := range m.Values {
		name := mapKeyString(mv.Key)
		seq, ok := mv.Value.(*ast.SequenceNode)
		if !ok {
			continue // non-list axis is not statically resolvable
		}
		for _, item := range seq.Values {
			if im := toMapping(item); im != nil {
				for _, sub := range im.Values {
					acc.add(name+"."+mapKeyString(sub.Key), nodeString(sub.Value))
				}
			} else {
				acc.add(name, nodeString(item))
			}
		}
	}
	return acc.axes()
}

// axisAccum collects matrix axis values in first-seen order so the rendered value lists
// are deterministic.
type axisAccum struct {
	order []string
	vals  map[string][]string
}

func newAxisAccum() *axisAccum {
	return &axisAccum{vals: map[string][]string{}}
}

func (a *axisAccum) add(name, val string) {
	if _, ok := a.vals[name]; !ok {
		a.order = append(a.order, name)
	}
	a.vals[name] = append(a.vals[name], val)
}

func (a *axisAccum) axes() []model.MatrixAxis {
	var out []model.MatrixAxis
	for _, n := range a.order {
		out = append(out, model.MatrixAxis{Name: n, Values: a.vals[n]})
	}
	return out
}

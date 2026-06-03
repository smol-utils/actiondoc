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

// parseMatrix extracts the declared axes of a job's strategy.matrix in every form the
// syntax allows:
//   - literal list axes (os: [a, b]) -> one axis per key
//   - list-of-objects axes (java: [{version: 17}, ...]) -> dotted axes (java.version)
//   - include: entries -> their keys and values merge into the axes (include adds
//     combinations, so its values are part of the declared surface)
//   - expression-valued axes (node: ${{ fromJSON(...) }}) -> the expression itself is the
//     axis's shown value, since it cannot be enumerated statically
//
// The second return reports whether an include:/exclude: key is present: the cartesian
// product of the listed values is then not the exact job set, and the renderer says so.
func parseMatrix(node ast.Node) ([]model.MatrixAxis, bool) {
	strat := toMapping(node)
	if strat == nil {
		return nil, false
	}
	var matrixNode ast.Node
	for _, mv := range strat.Values {
		if mapKeyString(mv.Key) == "matrix" {
			matrixNode = mv.Value
		}
	}
	m := toMapping(matrixNode)
	if m == nil {
		// The whole matrix is an expression (matrix: ${{ fromJSON(...) }}): there are no
		// axis names to list.
		return nil, false
	}

	acc := newAxisAccum()
	adjusted := false
	for _, mv := range m.Values {
		name := mapKeyString(mv.Key)
		switch name {
		case "exclude":
			adjusted = true
		case "include":
			adjusted = true
			// Each include entry is an object whose keys and values extend the matrix;
			// fold them into the axes so include-only matrices still document their values.
			if seq, ok := mv.Value.(*ast.SequenceNode); ok {
				for _, item := range seq.Values {
					if im := toMapping(item); im != nil {
						for _, sub := range im.Values {
							addAxisValue(acc, mapKeyString(sub.Key), sub.Value)
						}
					}
				}
			}
		default:
			seq, ok := mv.Value.(*ast.SequenceNode)
			if !ok {
				// Expression-valued axis: show the expression as the value.
				acc.add(name, nodeString(mv.Value))
				continue
			}
			for _, item := range seq.Values {
				addAxisValue(acc, name, item)
			}
		}
	}
	return acc.axes(), adjusted
}

// addAxisValue records one axis value, flattening an object value into dotted sub-axes
// (java: [{version: 17}] -> java.version: 17).
func addAxisValue(acc *axisAccum, name string, item ast.Node) {
	if im := toMapping(item); im != nil {
		for _, sub := range im.Values {
			acc.add(name+"."+mapKeyString(sub.Key), nodeString(sub.Value))
		}
		return
	}
	acc.add(name, nodeString(item))
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
	// Dedup: include entries commonly repeat values already declared on the axis.
	for _, v := range a.vals[name] {
		if v == val {
			return
		}
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

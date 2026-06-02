package model

// MatrixAxis is one statically-resolvable axis of a strategy.matrix. For a scalar list
// (`os: [ubuntu-latest, windows-latest]`) Name is the axis key. For a list of objects
// (`java: [{version: 17}, {version: 21}]`) each sub-field becomes its own dotted axis
// (`java.version`) so a `${{ matrix.java.version }}` reference resolves to a value list.
type MatrixAxis struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

// MatrixValues returns the value list for a matrix axis (matched on the full dotted name,
// e.g. "os" or "java.version") and whether the axis was statically resolved.
func (j *Job) MatrixValues(axis string) ([]string, bool) {
	for _, a := range j.Matrix {
		if a.Name == axis {
			return a.Values, true
		}
	}
	return nil, false
}

// Input returns the action's declared input with the given name, or nil.
func (a *Action) Input(name string) *ActionInput {
	for i := range a.Inputs {
		if a.Inputs[i].Name == name {
			return &a.Inputs[i]
		}
	}
	return nil
}

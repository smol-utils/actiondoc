package model

import (
	"fmt"
	"strings"
)

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

// Label is the human reference for a step: explicit name, then id, then a readable form
// of its uses: ref (SHA pins collapse to ref@version, or just ref), then a positional
// fallback. Reference inventories and rendered step titles share this rule, so the same
// step is always called the same thing everywhere it appears.
func (s *Step) Label(num int) string {
	switch {
	case s.Name != "":
		return s.Name
	case s.ID != "":
		return s.ID
	case s.Uses != "":
		return CollapsePin(s.Uses, s.UsesVersion)
	default:
		return fmt.Sprintf("step %d", num)
	}
}

// CollapsePin collapses a SHA-pinned action ref to its human form: `owner/repo@v4.1.1`
// when a version is known, or `owner/repo` when only the bare SHA is. Non-SHA refs
// (tags, branches, local paths) pass through unchanged.
func CollapsePin(uses, version string) string {
	at := strings.LastIndex(uses, "@")
	if at < 0 {
		return uses
	}
	ref, pin := uses[:at], uses[at+1:]
	if !IsSHA(pin) {
		return uses
	}
	if version != "" {
		return ref + "@" + version
	}
	return ref
}

// IsSHA reports whether s is a 40-character hexadecimal commit SHA.
func IsSHA(s string) bool {
	if len(s) != 40 {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		if !(c >= '0' && c <= '9' || c >= 'a' && c <= 'f' || c >= 'A' && c <= 'F') {
			return false
		}
	}
	return true
}

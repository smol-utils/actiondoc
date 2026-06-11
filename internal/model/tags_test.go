package model

import "testing"

func TestTagsEmpty(t *testing.T) {
	if !(Tags{}).Empty() {
		t.Error("zero Tags should be empty")
	}
	cases := map[string]Tags{
		"desc":    {Desc: "x"},
		"since":   {Since: "1.0"},
		"see":     {See: []string{"x"}},
		"secrets": {Secrets: []Param{{Name: "S"}}},
	}
	for name, tags := range cases {
		if tags.Empty() {
			t.Errorf("Tags with %s set should not be empty", name)
		}
	}
}

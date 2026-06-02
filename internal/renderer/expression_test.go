package renderer

import (
	"reflect"
	"testing"
)

func TestExpressions(t *testing.T) {
	cases := []struct {
		in   string
		want []string
	}{
		{"no expr here", nil},
		{"Java ${{ matrix.java.version }}", []string{"matrix.java.version"}},
		{"${{ inputs.version }}-${{ github.sha }}", []string{"inputs.version", "github.sha"}},
		{"unterminated ${{ matrix.x", nil},
	}
	for _, c := range cases {
		if got := expressions(c.in); !reflect.DeepEqual(got, c.want) {
			t.Errorf("expressions(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestExpressionKind(t *testing.T) {
	cases := []struct{ in, want string }{
		{"matrix.java.version", "matrix"},
		{"github.event.workflow_run.display_title", "github"},
		{"secrets.NPM_TOKEN", "secrets"},
		{"inputs.version", "inputs"},
		{"format('{0}', x)", "other"},
	}
	for _, c := range cases {
		if got := expressionKind(c.in); got != c.want {
			t.Errorf("expressionKind(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

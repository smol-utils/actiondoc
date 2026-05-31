package renderer

import "testing"

func TestEscapeCell(t *testing.T) {
	cases := []struct{ in, want string }{
		{"plain", "plain"},
		{"a | b", "a \\| b"},
		{"line1\nline2", "line1<br>line2"},
		{"crlf\r\nline", "crlf<br>line"},
		{"x &&\ny || z", "x &&<br>y \\|\\| z"},
	}
	for _, c := range cases {
		if got := escapeCell(c.in); got != c.want {
			t.Errorf("escapeCell(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestAnchor(t *testing.T) {
	cases := []struct{ in, want string }{
		{"build", "build"},
		{"Build and Test", "build-and-test"},
		{"publish-sdkman.yml", "publish-sdkmanyml"},
		{"Deploy (prod)", "deploy-prod"},
	}
	for _, c := range cases {
		if got := anchor(c.in); got != c.want {
			t.Errorf("anchor(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

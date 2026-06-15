package renderer

import (
	"fmt"
	"strings"
	"testing"

	"github.com/smol-utils/actiondoc/internal/model"
)

// TestJobDisplayName covers the de-templating rule at its single origin: a name embedding
// ${{ ... }} expressions becomes a readable label (each placeholder -> its trailing
// identifier in parentheses), a plain name is returned unchanged, and a name that is pure
// placeholder noise (no literal word) falls back to the job id.
func TestJobDisplayName(t *testing.T) {
	cases := []struct {
		name string
		id   string
		want string
	}{
		// Plain names pass through untouched.
		{"build", "build", "build"},
		{"Build and test", "build", "Build and test"},
		// Inline placeholders become the trailing identifier in parentheses.
		{"CLI ${{ matrix.job.os }}", "build-cli", "CLI (os)"},
		{"Build PROD ${{ inputs.build-type }} image ${{ matrix.python-version }}", "build-prod-images",
			"Build PROD (build-type) image (python-version)"},
		{"builder-${{matrix.os}}-${{matrix.arch}}", "builder", "builder-(os)-(arch)"},
		// Pure placeholder noise (no literal word) falls back to the id.
		{"${{ matrix.job.jdkOs }}", "jpackage", "jpackage"},
		{"${{ inputs.workflow-name || 'UI E2E Tests' }}", "test-ui-e2e-tests", "test-ui-e2e-tests"},
		{"[${{ inputs.target-branch }}]", "createupgrade-check", "createupgrade-check"},
	}
	for _, c := range cases {
		got := jobDisplayName(&model.Job{ID: c.id, Name: c.name})
		if got != c.want {
			t.Errorf("jobDisplayName(%q, id=%q) = %q, want %q", c.name, c.id, got, c.want)
		}
		if strings.Contains(got, "${{") {
			t.Errorf("jobDisplayName(%q) leaked a raw expression: %q", c.name, got)
		}
	}
}

// TestJobLabelAnchorConsistency verifies that the de-templated label flows into the heading,
// the mini-TOC label and link, and the anchor pass identically, so a mini-TOC link resolves
// to the heading GitHub actually emits. This is the property that breaks if de-templating is
// applied in one consumer but not another.
func TestJobLabelAnchorConsistency(t *testing.T) {
	w := &model.Workflow{
		File: "test.yml",
		Name: "Test",
		On:   []string{"push"},
		Jobs: []model.Job{
			{ID: "build-cli", Name: "CLI ${{ matrix.job.os }}"},
			{ID: "plain", Name: "plain"},
		},
	}
	md := RenderMarkdown(w)

	job := &w.Jobs[0]
	// The heading text and the mini-TOC label both come from jobDisplayName.
	if want := "### CLI (os) (`build-cli`)"; !strings.Contains(md, want) {
		t.Fatalf("missing de-templated heading %q in:\n%s", want, md)
	}
	// The mini-TOC link target must equal the anchor of the rendered heading text.
	slug := AssignAnchors([]string{JobHeadingText(job)})[0]
	link := fmt.Sprintf("[CLI (os)](#%s)", slug)
	if !strings.Contains(md, link) {
		t.Errorf("mini-TOC link %q does not match rendered heading anchor in:\n%s", link, md)
	}
	if slug != "cli-os-build-cli" {
		t.Errorf("unexpected anchor slug %q (heading text %q)", slug, JobHeadingText(job))
	}
}

package renderer

import (
	"strings"
	"testing"

	"github.com/smol-utils/actiondoc/internal/model"
)

func TestRenderMarkdownBasic(t *testing.T) {
	w := &model.Workflow{
		File:        "ci.yml",
		Name:        "CI Pipeline",
		Description: "Main CI pipeline.",
		On:          []string{"push", "pull_request"},
		Tags: model.Tags{
			Desc:  "Main CI pipeline.",
			Since: "v1.0.0",
			Secrets: []model.Param{
				{Name: "DEPLOY_KEY", Description: "SSH deploy key"},
			},
		},
		Jobs: []model.Job{
			{
				ID:          "build",
				Name:        "Build",
				Description: "Compile the application.",
				RunsOn:      "ubuntu-latest",
				Tags: model.Tags{
					Desc: "Compile the application.",
				},
				Steps: []model.Step{
					{Name: "Checkout", Uses: "actions/checkout@v4", Description: "Check out code."},
					{Name: "Build", Run: "make build"},
				},
			},
			{
				ID:          "test",
				Name:        "Run Tests",
				Description: "Run the test suite.",
				RunsOn:      "ubuntu-latest",
				Needs:       []string{"build"},
				If:          "github.event_name == 'push'",
				Tags: model.Tags{
					Desc: "Run the test suite.",
				},
				Steps: []model.Step{
					{Name: "Run tests", Run: "npm test"},
				},
			},
		},
	}

	md := RenderMarkdown(w)

	checks := []string{
		"# CI Pipeline",
		"Main CI pipeline.",
		"`ci.yml`",
		"`push`",
		"v1.0.0",
		"`DEPLOY_KEY`",
		"### Build (`build`)",
		"Compile the application.",
		"`ubuntu-latest`",
		"### Run Tests (`test`)",
		"`build`",
		"`github.event_name == 'push'`",
		"1. **Checkout** - Check out code.",
		"`actions/checkout@v4`",
	}

	for _, want := range checks {
		if !strings.Contains(md, want) {
			t.Errorf("output missing %q\n\nFull output:\n%s", want, md)
		}
	}
}

func TestRenderMarkdownDeprecated(t *testing.T) {
	w := &model.Workflow{
		File: "old.yml",
		Name: "Old Workflow",
		On:   []string{"push"},
		Tags: model.Tags{
			Deprecated: "Use new-workflow.yml instead.",
		},
	}

	md := RenderMarkdown(w)
	if !strings.Contains(md, "> **Deprecated**: Use new-workflow.yml instead.") {
		t.Errorf("missing deprecated banner in:\n%s", md)
	}
}

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

func TestRenderActionMarkdownBasic(t *testing.T) {
	a := &model.Action{
		File:        "action.yml",
		Name:        "My Action",
		Description: "Does a thing.",
		Inputs: []model.ActionInput{
			{Name: "token", Description: "GitHub token", Required: true},
			{Name: "env", Description: "Target environment", Required: false, Default: "staging"},
		},
		Outputs: []model.ActionOutput{
			{Name: "result", Description: "The result"},
		},
		Runs: model.ActionRuns{Using: "node20"},
	}

	md := RenderActionMarkdown(a)

	checks := []string{
		"# My Action",
		"Does a thing.",
		"`node20`",
		"## Inputs",
		"| `token`",
		"| Yes |",
		"| `env`",
		"| No |",
		"`staging`",
		"## Outputs",
		"| `result`",
	}

	for _, want := range checks {
		if !strings.Contains(md, want) {
			t.Errorf("output missing %q\n\nFull output:\n%s", want, md)
		}
	}
}

func TestTableEscaping(t *testing.T) {
	w := &model.Workflow{
		File: "test.yml",
		Name: "Test",
		On:   []string{"push"},
		Tags: model.Tags{
			Secrets: []model.Param{
				{Name: "TOKEN", Description: "Use a|b for options"},
			},
		},
	}

	md := RenderMarkdown(w)

	// The pipe should be escaped so it doesn't break the table.
	if !strings.Contains(md, `a\|b`) {
		t.Errorf("pipe character not escaped in table cell.\n\nOutput:\n%s", md)
	}

	// The unescaped pipe should NOT appear in the table rows (header divider is ok).
	for _, line := range strings.Split(md, "\n") {
		if !strings.HasPrefix(line, "|") {
			continue
		}
		// Skip the header divider row.
		if strings.Contains(line, "---") {
			continue
		}
		// Count unescaped pipes (not preceded by backslash).
		content := strings.ReplaceAll(line, `\|`, "")
		pipes := strings.Count(content, "|")
		// A valid 3-column row has exactly 4 pipes: | col | col | col |
		if strings.Contains(line, "TOKEN") && pipes != 4 {
			t.Errorf("table row has %d unescaped pipes (expected 4): %s", pipes, line)
		}
	}
}

// TestJobAndStepInputTagsRender locks the spec contract that @input applies at all three
// scopes: a job-level or step-level @input tag must render, not silently disappear.
func TestJobAndStepInputTagsRender(t *testing.T) {
	w := &model.Workflow{
		File: "test.yml",
		Name: "Test",
		On:   []string{"push"},
		Jobs: []model.Job{{
			ID:     "deploy",
			Name:   "deploy",
			RunsOn: "ubuntu-latest",
			Tags: model.Tags{
				Inputs: []model.Param{{Name: "environment", Type: "string", Description: "Target environment"}},
			},
			Steps: []model.Step{{
				Name: "Run",
				Uses: "./.github/actions/deploy",
				Tags: model.Tags{
					Inputs: []model.Param{{Name: "region", Description: "AWS region"}},
				},
			}},
		}},
	}

	md := RenderMarkdown(w)

	if !strings.Contains(md, "**Inputs:**") || !strings.Contains(md, "`environment`") {
		t.Errorf("job-level @input not rendered:\n%s", md)
	}
	if !strings.Contains(md, "Input: `region` - AWS region") {
		t.Errorf("step-level @input not rendered:\n%s", md)
	}
}

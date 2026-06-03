package parser

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/smol-utils/actiondoc/internal/model"
)

func testdataPath(name string) string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "testdata", name)
}

func TestParseFileSample(t *testing.T) {
	path := testdataPath("sample-workflow.yml")
	if _, err := os.Stat(path); err != nil {
		t.Skipf("testdata not found: %v", err)
	}

	w, err := ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}

	// Workflow-level checks
	if w.Name != "CI Pipeline" {
		t.Errorf("workflow name = %q", w.Name)
	}
	if w.Description == "" {
		t.Error("expected workflow description")
	}
	if w.Tags.Since != "v1.0.0" {
		t.Errorf("since = %q", w.Tags.Since)
	}
	if len(w.Tags.See) != 1 {
		t.Errorf("see = %+v", w.Tags.See)
	}
	if len(w.Tags.Secrets) != 1 || w.Tags.Secrets[0].Name != "DEPLOY_KEY" {
		t.Errorf("workflow secrets = %+v", w.Tags.Secrets)
	}
	if len(w.On) != 2 {
		t.Errorf("triggers = %+v", w.On)
	}

	// Jobs
	if len(w.Jobs) != 3 {
		t.Fatalf("expected 3 jobs, got %d", len(w.Jobs))
	}

	// Build job
	build := w.Jobs[0]
	if build.ID != "build" {
		t.Errorf("job 0 id = %q", build.ID)
	}
	if build.Name != "Build" {
		t.Errorf("job 0 name = %q", build.Name)
	}
	if build.Description == "" {
		t.Error("expected build job description")
	}
	if len(build.Steps) != 2 {
		t.Errorf("build steps = %d", len(build.Steps))
	}

	// Test job
	test := w.Jobs[1]
	if test.ID != "test" {
		t.Errorf("job 1 id = %q", test.ID)
	}
	if len(test.Needs) != 1 || test.Needs[0] != "build" {
		t.Errorf("test needs = %+v", test.Needs)
	}
	if test.If != "github.event_name == 'push'" {
		t.Errorf("test if = %q", test.If)
	}
	if len(test.Tags.Secrets) != 1 || test.Tags.Secrets[0].Name != "NPM_TOKEN" {
		t.Errorf("test secrets = %+v", test.Tags.Secrets)
	}

	// Deploy job
	deploy := w.Jobs[2]
	if deploy.Tags.Deprecated == "" {
		t.Error("expected deploy job to be deprecated")
	}
	if deploy.Tags.Example == "" {
		t.Error("expected deploy job to have an example")
	}
}

func TestParseActionFile(t *testing.T) {
	path := testdataPath("action.yml")
	if _, err := os.Stat(path); err != nil {
		t.Skipf("testdata not found: %v", err)
	}

	a, err := ParseActionFile(path)
	if err != nil {
		t.Fatalf("ParseActionFile: %v", err)
	}

	if a.Name != "Deploy Action" {
		t.Errorf("name = %q", a.Name)
	}
	if a.Description != "Deploy an application to the target environment." {
		t.Errorf("description = %q", a.Description)
	}
	if a.Runs.Using != "node20" {
		t.Errorf("runs.using = %q", a.Runs.Using)
	}
	if a.Tags.Since != "v1.0.0" {
		t.Errorf("since = %q", a.Tags.Since)
	}
	if len(a.Tags.See) != 1 {
		t.Errorf("see = %+v", a.Tags.See)
	}
	if len(a.Tags.Secrets) != 1 || a.Tags.Secrets[0].Name != "DEPLOY_TOKEN" {
		t.Errorf("secrets = %+v", a.Tags.Secrets)
	}
	if a.Tags.Example == "" {
		t.Error("expected example")
	}

	// Inputs
	if len(a.Inputs) != 3 {
		t.Fatalf("expected 3 inputs, got %d", len(a.Inputs))
	}
	env := a.Inputs[0]
	if env.Name != "environment" || !env.Required {
		t.Errorf("input 0 = %+v", env)
	}
	dryRun := a.Inputs[2]
	if dryRun.Name != "dry-run" || dryRun.Required || dryRun.Default != "false" {
		t.Errorf("input 2 = %+v", dryRun)
	}

	// Outputs
	if len(a.Outputs) != 2 {
		t.Fatalf("expected 2 outputs, got %d", len(a.Outputs))
	}
	if a.Outputs[0].Name != "deploy-url" {
		t.Errorf("output 0 name = %q", a.Outputs[0].Name)
	}

	// Branding
	if a.Branding == nil {
		t.Fatal("expected branding")
	}
	if a.Branding.Icon != "upload-cloud" || a.Branding.Color != "blue" {
		t.Errorf("branding = %+v", a.Branding)
	}
}

// TestResolveAnchors covers YAML anchor/alias indirection across field types: a scalar
// anchor on runs-on (the syft pattern), an alias reusing it, an anchored mapping reused
// as a job env: block, and an anchored sequence. Field readers must always see resolved
// values, never raw &name / *name tokens.
func TestResolveAnchors(t *testing.T) {
	src := `name: Anchors
on: push
jobs:
  define:
    runs-on: &test-runner "runs-on=large/cpu=8/ram=32"
    env: &common-env
      REGION: us-east-1
      TIER: production
    steps:
      - run: make build
  reuse:
    runs-on: *test-runner
    env: *common-env
    steps:
      - name: Reuse list
        uses: actions/checkout@v4
        with:
          paths: &path-list
            - src/**
            - docs/**
      - name: Alias list
        uses: actions/checkout@v4
        with:
          paths: *path-list
`
	w, err := parseString(t, src)
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}

	define, reuse := w.Jobs[0], w.Jobs[1]

	// Anchor definition site: the value, not "&test-runner ...".
	if define.RunsOn != "runs-on=large/cpu=8/ram=32" {
		t.Errorf("anchor definition RunsOn = %q, want the unwrapped value", define.RunsOn)
	}
	// Alias site: the resolved value, not "*test-runner".
	if reuse.RunsOn != "runs-on=large/cpu=8/ram=32" {
		t.Errorf("alias RunsOn = %q, want the resolved value", reuse.RunsOn)
	}

	// Anchored mapping reused as a whole env: block.
	for _, job := range []model.Job{define, reuse} {
		if len(job.Env) != 2 || job.Env[0].Key != "REGION" || job.Env[1].Value != "production" {
			t.Errorf("job %s Env = %+v, want resolved common-env mapping", job.ID, job.Env)
		}
	}

	// No raw anchor/alias tokens anywhere in the parsed model.
	for _, job := range w.Jobs {
		if strings.Contains(job.RunsOn, "&") || strings.Contains(job.RunsOn, "*test-runner") {
			t.Errorf("job %s RunsOn contains raw anchor/alias token: %q", job.ID, job.RunsOn)
		}
	}
}

// TestExpressionValuesKeptVerbatim locks the expression half of the field-reading rule:
// a ${{ }} expression value in any string field is kept verbatim, never interpreted,
// dropped, or coerced.
func TestExpressionValuesKeptVerbatim(t *testing.T) {
	src := `name: Expressions
on: push
concurrency:
  group: ci-${{ github.ref }}
jobs:
  build:
    runs-on: ${{ matrix.os }}
    environment: ${{ inputs.env }}
    strategy:
      matrix:
        os: [ubuntu-latest, macos-14]
    steps:
      - run: echo hi
        env:
          REF: ${{ github.ref }}
`
	w, err := parseString(t, src)
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}
	job := w.Jobs[0]
	for _, tc := range []struct{ got, want, what string }{
		{w.Concurrency.Group, "ci-${{ github.ref }}", "concurrency group"},
		{job.RunsOn, "${{ matrix.os }}", "runs-on"},
		{job.Environment.Name, "${{ inputs.env }}", "environment name"},
		{job.Steps[0].Env[0].Value, "${{ github.ref }}", "step env value"},
	} {
		if tc.got != tc.want {
			t.Errorf("%s = %q, want verbatim expression %q", tc.what, tc.got, tc.want)
		}
	}
}

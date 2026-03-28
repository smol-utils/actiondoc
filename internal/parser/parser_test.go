package parser

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
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

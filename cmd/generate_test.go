package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestResolveFilesCompositeDiscovery verifies that pointing at a .github/workflows dir
// discovers both the workflow files and sibling composite actions under
// .github/actions/ (including a nested one and a composite placed in the workflows dir).
func TestResolveFilesCompositeDiscovery(t *testing.T) {
	root := t.TempDir()
	mkdir := func(p string) {
		if err := os.MkdirAll(p, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	write := func(p string) {
		mkdir(filepath.Dir(p))
		if err := os.WriteFile(p, []byte("name: x\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	gh := filepath.Join(root, ".github")
	write(filepath.Join(gh, "workflows", "ci.yml"))
	write(filepath.Join(gh, "workflows", "release.yaml"))
	write(filepath.Join(gh, "workflows", "action.yml"))             // composite placed in the workflows dir
	write(filepath.Join(gh, "actions", "build", "action.yml"))      // sibling composite
	write(filepath.Join(gh, "actions", "deep", "x", "action.yaml")) // nested composite

	got, err := resolveFiles(filepath.Join(gh, "workflows"))
	if err != nil {
		t.Fatalf("resolveFiles: %v", err)
	}

	// Compare on paths relative to the temp root for stable assertions.
	found := map[string]bool{}
	for _, f := range got {
		rel, _ := filepath.Rel(root, f)
		found[filepath.ToSlash(rel)] = true
	}
	want := []string{
		".github/workflows/ci.yml",
		".github/workflows/release.yaml",
		".github/workflows/action.yml",
		".github/actions/build/action.yml",
		".github/actions/deep/x/action.yaml",
	}
	for _, w := range want {
		if !found[w] {
			t.Errorf("expected to discover %s; got %v", w, got)
		}
	}
	if len(got) != len(want) {
		t.Errorf("got %d files, want %d: %v", len(got), len(want), got)
	}

	// Discovery order must be deterministic so directory-mode TOC/section order is stable.
	for i := 0; i < 3; i++ {
		again, err := resolveFiles(filepath.Join(gh, "workflows"))
		if err != nil {
			t.Fatalf("resolveFiles (repeat): %v", err)
		}
		if strings.Join(again, "\n") != strings.Join(got, "\n") {
			t.Errorf("resolveFiles order not deterministic:\n first: %v\n again: %v", got, again)
		}
	}
}

// TestGenerateNoWorkflowName verifies that a workflow without a `name:` still renders a
// non-empty title and TOC entry (the parser backfills the file name), so directory-mode
// output never produces an empty heading or anchor.
func TestGenerateNoWorkflowName(t *testing.T) {
	root := t.TempDir()
	wf := filepath.Join(root, ".github", "workflows")
	if err := os.MkdirAll(wf, 0o755); err != nil {
		t.Fatal(err)
	}
	// two unnamed workflows -> distinct file-name titles, no empty/duplicate anchors
	for _, name := range []string{"ci.yml", "release.yml"} {
		if err := os.WriteFile(filepath.Join(wf, name), []byte("on: push\njobs:\n  x:\n    runs-on: ubuntu-latest\n    steps:\n      - run: echo hi\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	out := filepath.Join(t.TempDir(), "out.md")
	if err := Generate([]string{"-o", out, wf}); err != nil {
		t.Fatalf("Generate: %v", err)
	}
	b, _ := os.ReadFile(out)
	md := string(b)
	for _, want := range []string{"# ci.yml", "# release.yml", "[ci.yml](#ciyml)", "[release.yml](#releaseyml)"} {
		if !strings.Contains(md, want) {
			t.Errorf("output missing %q (empty name not backfilled?)\n%s", want, md)
		}
	}
}

// TestGenerateDuplicateNameCrossLinks verifies that when two workflows share a name, a
// caller's cross-link points at the section the TOC assigned to its actual callee (the
// "-N"-suffixed anchor), not at the first section with that name.
func TestGenerateDuplicateNameCrossLinks(t *testing.T) {
	root := t.TempDir()
	wf := filepath.Join(root, ".github", "workflows")
	if err := os.MkdirAll(wf, 0o755); err != nil {
		t.Fatal(err)
	}
	write := func(name, content string) {
		if err := os.WriteFile(filepath.Join(wf, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	// Both reusable workflows are named "model jobs"; the caller calls the SECOND one
	// (model_jobs_gaudi.yml), which the TOC will disambiguate as #model-jobs-1.
	write("caller.yml", `name: Scheduler
on: push
jobs:
  run:
    uses: ./.github/workflows/model_jobs_gaudi.yml
`)
	write("model_jobs.yml", `name: model jobs
on: workflow_call
jobs:
  x:
    runs-on: ubuntu-latest
    steps:
      - run: echo a
`)
	write("model_jobs_gaudi.yml", `name: model jobs
on: workflow_call
jobs:
  x:
    runs-on: ubuntu-latest
    steps:
      - run: echo b
`)

	out := filepath.Join(t.TempDir(), "out.md")
	if err := Generate([]string{"-o", out, wf}); err != nil {
		t.Fatalf("Generate: %v", err)
	}
	b, _ := os.ReadFile(out)
	md := string(b)

	// File order is alphabetical: caller.yml, model_jobs.yml, model_jobs_gaudi.yml.
	// So model_jobs.yml gets #model-jobs and model_jobs_gaudi.yml gets #model-jobs-1.
	if !strings.Contains(md, "[model jobs](#model-jobs-1)") {
		t.Errorf("caller link must point at the disambiguated anchor #model-jobs-1:\n%s", md)
	}
	// The TOC must contain both anchors.
	for _, want := range []string{"- [model jobs](#model-jobs)\n", "- [model jobs](#model-jobs-1)\n"} {
		if !strings.Contains(md, want) {
			t.Errorf("TOC missing %q", want)
		}
	}
}

// TestGenerateSkipsDisabledWorkflows verifies that a fully commented-out workflow file
// (the common way to disable a workflow) is skipped with a note rather than treated as a
// parse failure: the run still succeeds and the other workflows still render.
func TestGenerateSkipsDisabledWorkflows(t *testing.T) {
	root := t.TempDir()
	wf := filepath.Join(root, ".github", "workflows")
	if err := os.MkdirAll(wf, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(wf, "ci.yml"), []byte("name: CI\non: push\njobs:\n  x:\n    runs-on: ubuntu-latest\n    steps:\n      - run: echo hi\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	// Disabled workflow: every line commented out.
	disabled := "# name: Old Deploy\n# on: push\n# jobs:\n#   deploy:\n#     runs-on: ubuntu-latest\n#     steps:\n#       - run: ./deploy.sh\n"
	if err := os.WriteFile(filepath.Join(wf, "deploy.yml"), []byte(disabled), 0o644); err != nil {
		t.Fatal(err)
	}
	// A genuinely malformed file must still fail the run.
	if err := os.WriteFile(filepath.Join(wf, "broken.yml"), []byte("- this is a list\n- not a mapping\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	out := filepath.Join(t.TempDir(), "out.md")

	// With only the disabled file present (remove broken), Generate must succeed.
	if err := os.Remove(filepath.Join(wf, "broken.yml")); err != nil {
		t.Fatal(err)
	}
	if err := Generate([]string{"-o", out, wf}); err != nil {
		t.Fatalf("Generate with a disabled workflow must succeed, got: %v", err)
	}
	b, _ := os.ReadFile(out)
	if !strings.Contains(string(b), "# CI") {
		t.Errorf("remaining workflow did not render:\n%s", string(b))
	}

	// With a genuinely malformed file, Generate must still fail.
	if err := os.WriteFile(filepath.Join(wf, "broken.yml"), []byte("- this is a list\n- not a mapping\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := Generate([]string{"-o", out, wf}); err == nil {
		t.Error("Generate with a malformed (non-mapping) file must fail")
	}
}

// TestGenerateJSONOutput pins the -json output mode: the result must be valid JSON whose
// structure exposes the documented model fields. The model's JSON tags are a stable
// contract for downstream consumers; nothing else guards them.
func TestGenerateJSONOutput(t *testing.T) {
	out := filepath.Join(t.TempDir(), "out.json")

	err := Generate([]string{"-json", "-o", out, testdataPath("s3/repo/.github/workflows")})
	if err != nil {
		t.Fatalf("Generate -json: %v", err)
	}

	raw, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("reading output: %v", err)
	}

	var items []map[string]any
	if err := json.Unmarshal(raw, &items); err != nil {
		t.Fatalf("output is not valid JSON: %v\n%s", err, raw)
	}
	if len(items) == 0 {
		t.Fatal("JSON output is empty")
	}

	// The fixture's workflows and composite action must all be present, keyed by the
	// documented model fields.
	var names []string
	for _, item := range items {
		name, _ := item["name"].(string)
		names = append(names, name)
	}
	for _, want := range []string{"CI", "Reusable Build", "Deploy"} {
		found := false
		for _, n := range names {
			if n == want {
				found = true
			}
		}
		if !found {
			t.Errorf("JSON output missing item %q; got names %v", want, names)
		}
	}

	// Spot-check the documented field structure on the workflow items: jobs with ids,
	// steps, and trigger lists all surface under their JSON tags.
	for _, item := range items {
		if item["name"] == "CI" {
			jobs, ok := item["jobs"].([]any)
			if !ok || len(jobs) == 0 {
				t.Fatalf("CI workflow JSON has no jobs: %v", item)
			}
			job := jobs[0].(map[string]any)
			if job["id"] == "" {
				t.Errorf("job missing id field: %v", job)
			}
			if _, hasOn := item["on"]; !hasOn {
				t.Errorf("workflow missing 'on' field: %v", item)
			}
		}
	}
}

// TestGenerateRedact verifies that -redact rewrites the output (secret names become
// placeholders) and that -redact-map writes the reverse map to the named file only.
func TestGenerateRedact(t *testing.T) {
	root := t.TempDir()
	wf := filepath.Join(root, ".github", "workflows")
	if err := os.MkdirAll(wf, 0o755); err != nil {
		t.Fatal(err)
	}
	content := "name: Deploy\n" +
		"on: push\n" +
		"jobs:\n" +
		"  deploy:\n" +
		"    runs-on: [self-hosted, prod-gpu]\n" +
		"    steps:\n" +
		"      - name: Push\n" +
		"        env:\n" +
		"          API_URL: https://api.internal.corp/v1\n" +
		"        run: |\n" +
		"          curl -H \"x: ${{ secrets.DEPLOY_TOKEN }}\" \"$API_URL\"\n"
	if err := os.WriteFile(filepath.Join(wf, "deploy.yml"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	out := filepath.Join(t.TempDir(), "out.md")
	mapPath := filepath.Join(t.TempDir(), "map.json")
	if err := Generate([]string{"-redact", "-redact-map", mapPath, "-o", out, wf}); err != nil {
		t.Fatalf("Generate -redact: %v", err)
	}

	md, _ := os.ReadFile(out)
	got := string(md)
	if strings.Contains(got, "DEPLOY_TOKEN") || strings.Contains(got, "api.internal.corp") || strings.Contains(got, "prod-gpu") {
		t.Errorf("redacted output still contains sensitive material:\n%s", got)
	}
	for _, want := range []string{"SECRET_1", "URL_1", "RUNNER_1"} {
		if !strings.Contains(got, want) {
			t.Errorf("redacted output missing placeholder %q:\n%s", want, got)
		}
	}

	rawMap, err := os.ReadFile(mapPath)
	if err != nil {
		t.Fatalf("reading redaction map: %v", err)
	}
	var m map[string]string
	if err := json.Unmarshal(rawMap, &m); err != nil {
		t.Fatalf("redaction map is not valid JSON: %v\n%s", err, rawMap)
	}
	if m["SECRET_1"] != "DEPLOY_TOKEN" {
		t.Errorf("reverse map wrong: %v", m)
	}
}

// TestGenerateRedactMapRequiresRedact verifies that -redact-map without a redact flag is a
// usage error, so the map file is never produced for non-redacted output.
func TestGenerateRedactMapRequiresRedact(t *testing.T) {
	root := t.TempDir()
	wf := filepath.Join(root, ".github", "workflows")
	if err := os.MkdirAll(wf, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(wf, "ci.yml"), []byte("name: CI\non: push\njobs:\n  x:\n    runs-on: ubuntu-latest\n    steps:\n      - run: echo hi\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	mapPath := filepath.Join(t.TempDir(), "map.json")
	if err := Generate([]string{"-redact-map", mapPath, wf}); err == nil {
		t.Fatal("expected error when -redact-map is used without -redact")
	}
	if _, err := os.Stat(mapPath); err == nil {
		t.Error("redaction map should not be written on usage error")
	}
}

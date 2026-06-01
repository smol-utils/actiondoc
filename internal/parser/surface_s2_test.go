package parser

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/goccy/go-yaml/ast"
	"github.com/smol-utils/actiondoc/internal/model"
)

// firstValue is a test helper: parse a one-key YAML snippet and return the value node of
// its first entry (e.g. the block under "permissions:" or "on:").
func firstValue(t *testing.T, src string) ast.Node {
	t.Helper()
	root := parseTopMapping(t, src)
	if len(root.Values) == 0 {
		t.Fatal("no entries")
	}
	return root.Values[0].Value
}

func s2Path(name string) string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "testdata", "s2", name)
}

func TestParseDispatchInputs(t *testing.T) {
	node := firstValue(t, `on:
  workflow_dispatch:
    inputs:
      distribution:
        type: choice
        required: true
        description: Which build to release.
        options: [mandrel, graalvm]
      dry-run:
        type: boolean
        default: false
`)
	tr := parseTriggerSurface(node)
	if tr == nil || tr.Dispatch == nil {
		t.Fatalf("expected dispatch trigger, got %+v", tr)
	}
	inputs := tr.Dispatch.Inputs
	if len(inputs) != 2 {
		t.Fatalf("expected 2 inputs, got %d", len(inputs))
	}
	dist := inputs[0]
	if dist.Name != "distribution" || dist.Type != "choice" || !dist.Required {
		t.Errorf("distribution input = %+v", dist)
	}
	if len(dist.Options) != 2 || dist.Options[0] != "mandrel" {
		t.Errorf("options = %+v", dist.Options)
	}
	dry := inputs[1]
	if dry.Name != "dry-run" || dry.Type != "boolean" || dry.Required || dry.Default != "false" {
		t.Errorf("dry-run input = %+v", dry)
	}
}

func TestParseCallTrigger(t *testing.T) {
	node := firstValue(t, `on:
  workflow_call:
    inputs:
      version:
        type: string
        required: true
    outputs:
      artifact-url:
        description: Artifact URL.
        value: ${{ jobs.build.outputs.url }}
    secrets:
      CONSUMER-KEY:
        required: true
        description: Publish key.
`)
	tr := parseTriggerSurface(node)
	if tr == nil || tr.Call == nil {
		t.Fatalf("expected call trigger, got %+v", tr)
	}
	c := tr.Call
	if len(c.Inputs) != 1 || c.Inputs[0].Name != "version" || !c.Inputs[0].Required {
		t.Errorf("inputs = %+v", c.Inputs)
	}
	if len(c.Outputs) != 1 || c.Outputs[0].Name != "artifact-url" || c.Outputs[0].Value == "" {
		t.Errorf("outputs = %+v", c.Outputs)
	}
	// Casing and hyphenation must be preserved, never normalized.
	if len(c.Secrets) != 1 || c.Secrets[0].Name != "CONSUMER-KEY" || !c.Secrets[0].Required {
		t.Errorf("secrets = %+v", c.Secrets)
	}
}

func TestParseScheduleVerbatim(t *testing.T) {
	node := firstValue(t, `on:
  schedule:
    - cron: '0 16 * * THU'  # weekly security scan
    - cron: '30 2 * * *'
`)
	tr := parseTriggerSurface(node)
	if tr == nil || len(tr.Schedule) != 2 {
		t.Fatalf("expected 2 cron entries, got %+v", tr)
	}
	// Named day-of-week is preserved verbatim, not normalized to numeric form.
	if tr.Schedule[0].Cron != "0 16 * * THU" {
		t.Errorf("cron[0] = %q", tr.Schedule[0].Cron)
	}
	if tr.Schedule[0].Rationale != "weekly security scan" {
		t.Errorf("rationale[0] = %q", tr.Schedule[0].Rationale)
	}
	if tr.Schedule[1].Cron != "30 2 * * *" || tr.Schedule[1].Rationale != "" {
		t.Errorf("cron[1] = %+v", tr.Schedule[1])
	}
}

func TestParseGenericEventFilters(t *testing.T) {
	node := firstValue(t, `on:
  push:
    branches: [main, 'release/*']
    paths-ignore:
      - 'docs/**'
  repository_dispatch:
    types: [deploy-request]
  pull_request:
`)
	tr := parseTriggerSurface(node)
	if tr == nil {
		t.Fatal("expected trigger surface")
	}
	// pull_request has no sub-keys, so only push and repository_dispatch appear.
	if len(tr.Events) != 2 {
		t.Fatalf("expected 2 events, got %+v", tr.Events)
	}
	push := tr.Events[0]
	if push.Name != "push" || len(push.Filters) != 2 {
		t.Fatalf("push event = %+v", push)
	}
	if push.Filters[0].Key != "branches" || len(push.Filters[0].Values) != 2 {
		t.Errorf("branches filter = %+v", push.Filters[0])
	}
	if push.Filters[1].Key != "paths-ignore" || push.Filters[1].Values[0] != "docs/**" {
		t.Errorf("paths-ignore filter = %+v", push.Filters[1])
	}
	rd := tr.Events[1]
	if rd.Name != "repository_dispatch" || rd.Filters[0].Key != "types" || rd.Filters[0].Values[0] != "deploy-request" {
		t.Errorf("repository_dispatch event = %+v", rd)
	}
}

func TestParsePermissionsScopes(t *testing.T) {
	node := firstValue(t, `permissions:
  contents: read  # for actions/checkout to fetch code
  id-token: write
  packages: write
`)
	p := parsePermissions(node)
	if p == nil || len(p.Scopes) != 3 {
		t.Fatalf("permissions = %+v", p)
	}
	if p.Scopes[0].Scope != "contents" || p.Scopes[0].Level != "read" {
		t.Errorf("scope[0] = %+v", p.Scopes[0])
	}
	if p.Scopes[0].Rationale != "for actions/checkout to fetch code" {
		t.Errorf("rationale = %q", p.Scopes[0].Rationale)
	}
	// id-token: write is the OIDC enabler; other write grants are not flagged.
	if !p.Scopes[1].OIDC {
		t.Error("expected OIDC flag on id-token: write")
	}
	if p.Scopes[2].OIDC {
		t.Error("packages: write must not be flagged OIDC")
	}
}

func TestParsePermissionsDefaultDeny(t *testing.T) {
	p := parsePermissions(firstValue(t, "permissions: {}\n"))
	if p == nil || !p.DefaultDeny || len(p.Scopes) != 0 {
		t.Errorf("permissions = %+v", p)
	}
}

func TestParsePermissionsScalar(t *testing.T) {
	p := parsePermissions(firstValue(t, "permissions: read-all\n"))
	if p == nil || p.All != "read-all" {
		t.Errorf("permissions = %+v", p)
	}
}

func TestParseEnvironment(t *testing.T) {
	scalar := parseEnvironment(firstValue(t, "environment: production\n"))
	if scalar == nil || scalar.Name != "production" || scalar.URL != "" {
		t.Errorf("scalar environment = %+v", scalar)
	}

	mapped := parseEnvironment(firstValue(t, `environment:
  name: production
  url: https://app.example.com
`))
	if mapped == nil || mapped.Name != "production" || mapped.URL != "https://app.example.com" {
		t.Errorf("mapped environment = %+v", mapped)
	}
}

func TestParseConcurrency(t *testing.T) {
	scalar := parseConcurrency(firstValue(t, "concurrency: release-lock\n"))
	if scalar == nil || scalar.Group != "release-lock" {
		t.Errorf("scalar concurrency = %+v", scalar)
	}

	mapped := parseConcurrency(firstValue(t, `concurrency:
  group: ci-${{ github.ref }}
  cancel-in-progress: true
`))
	if mapped == nil || mapped.Group != "ci-${{ github.ref }}" || mapped.CancelInProgress != "true" {
		t.Errorf("mapped concurrency = %+v", mapped)
	}
}

func TestParseDefaults(t *testing.T) {
	d := parseDefaults(firstValue(t, `defaults:
  run:
    shell: bash
    working-directory: ./app
`))
	if d == nil || d.Shell != "bash" || d.WorkingDirectory != "./app" {
		t.Errorf("defaults = %+v", d)
	}
}

func TestImplicitDescription(t *testing.T) {
	prose := "# Builds and publishes the release artifacts."
	if got := implicitDescription(prose, model.Tags{}); got != "Builds and publishes the release artifacts." {
		t.Errorf("prose comment: got %q", got)
	}

	license := "# Licensed to the Apache Software Foundation (ASF) under one\n# or more contributor license agreements."
	if got := implicitDescription(license, model.Tags{}); got != "" {
		t.Errorf("license header must not become a description, got %q", got)
	}

	// A comment with recognized tags (but no @desc) must not leak into the description.
	tagged := "# @secret TOKEN - deploy token"
	if got := implicitDescription(tagged, model.Tags{Secrets: []model.Param{{Name: "TOKEN"}}}); got != "" {
		t.Errorf("tagged comment must not become a description, got %q", got)
	}

	// A block of only unknown @-markers (another tool's, e.g. Lula) must be ignored, not
	// rendered as prose -- tags is empty here because none are ActionDoc tags.
	markers := "# @lulaStart policy-block\n# @lulaEnd"
	if got := implicitDescription(markers, model.Tags{}); got != "" {
		t.Errorf("unknown @-marker block must not become a description, got %q", got)
	}

	// Prose mixed with an unknown marker keeps only the prose.
	mixed := "# Deploys the service.\n# @lulaStart policy"
	if got := implicitDescription(mixed, model.Tags{}); got != "Deploys the service." {
		t.Errorf("mixed prose+marker: got %q, want prose only", got)
	}
}

// TestParseFileS2Surface checks the end-to-end IR for the session fixture: implicit
// description plus the workflow- and job-level declared surface.
func TestParseFileS2Surface(t *testing.T) {
	w, err := ParseFile(s2Path("surface.yml"))
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}
	if w.Description == "" {
		t.Error("expected implicit description from leading comment block")
	}
	if w.Triggers == nil || w.Triggers.Dispatch == nil || len(w.Triggers.Schedule) != 2 || len(w.Triggers.Events) != 3 {
		t.Errorf("triggers = %+v", w.Triggers)
	}
	if w.Permissions == nil || len(w.Permissions.Scopes) != 3 {
		t.Errorf("workflow permissions = %+v", w.Permissions)
	}
	if len(w.Env) != 2 || w.Concurrency == nil || w.Defaults == nil {
		t.Errorf("workflow metadata: env=%+v concurrency=%+v defaults=%+v", w.Env, w.Concurrency, w.Defaults)
	}

	if len(w.Jobs) != 2 {
		t.Fatalf("expected 2 jobs, got %d", len(w.Jobs))
	}
	build, deploy := w.Jobs[0], w.Jobs[1]
	if build.Permissions == nil || !build.Permissions.DefaultDeny {
		t.Errorf("build permissions = %+v", build.Permissions)
	}
	if deploy.Environment == nil || deploy.Environment.Name != "production" || deploy.Environment.URL == "" {
		t.Errorf("deploy environment = %+v", deploy.Environment)
	}
	if deploy.Concurrency == nil || deploy.Defaults == nil || deploy.Defaults.WorkingDirectory != "./deploy" {
		t.Errorf("deploy metadata: concurrency=%+v defaults=%+v", deploy.Concurrency, deploy.Defaults)
	}
}

// TestParseFileLicenseHeaderSkipped checks that an ASF license header does not become an
// implicit description.
func TestParseFileLicenseHeaderSkipped(t *testing.T) {
	w, err := ParseFile(s2Path("license-header.yml"))
	if err != nil {
		t.Fatalf("ParseFile: %v", err)
	}
	if w.Description != "" {
		t.Errorf("license header leaked into description: %q", w.Description)
	}
}

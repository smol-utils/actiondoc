package model

import (
	"strings"
	"testing"
)

// TestScanReferences covers secret/var collection across job if:, forwarded secrets:,
// step if:/run:/with:, including dedup, site merging, and the no-references case.
func TestScanReferences(t *testing.T) {
	w := &Workflow{
		Name: "Release",
		Jobs: []Job{
			{
				ID: "publish",
				If: "${{ vars.RELEASE_ENABLED }}",
				Secrets: []KV{
					{Key: "CONSUMER-KEY", Value: "${{ secrets.SDKMAN_KEY }}"},
				},
			},
			{
				ID: "build",
				Steps: []Step{
					{
						Name: "Sign",
						Run:  "sign --key \"${{ secrets.GPG_KEY }}\" --pass \"${{ secrets.GPG_PASSPHRASE }}\"",
					},
					{
						ID: "login",
						If: "${{ vars.USE_PROD && secrets.PROD_TOKEN || secrets.DEV_TOKEN }}",
						With: []KV{
							{Key: "password", Value: "${{ secrets.GPG_KEY }}"}, // repeat -> merged
						},
					},
					{
						Run: "echo plain secrets.NOT_AN_EXPRESSION", // outside ${{ }} -> ignored
					},
				},
			},
		},
	}

	refs := ScanReferences(w)

	gotSecrets := map[string][]string{}
	var secretOrder []string
	for _, r := range refs.Secrets {
		gotSecrets[r.Name] = r.Sites
		secretOrder = append(secretOrder, r.Name)
	}

	wantOrder := []string{"SDKMAN_KEY", "GPG_KEY", "GPG_PASSPHRASE", "PROD_TOKEN", "DEV_TOKEN"}
	if strings.Join(secretOrder, ",") != strings.Join(wantOrder, ",") {
		t.Errorf("secret order = %v, want %v", secretOrder, wantOrder)
	}

	// GPG_KEY is referenced from two sites (run and with) -> both recorded.
	if sites := gotSecrets["GPG_KEY"]; len(sites) != 2 {
		t.Errorf("GPG_KEY sites = %v, want 2 sites", sites)
	}
	// Forwarded secret site labels the forwarding key.
	if sites := gotSecrets["SDKMAN_KEY"]; len(sites) != 1 || !strings.Contains(sites[0], "CONSUMER-KEY") {
		t.Errorf("SDKMAN_KEY sites = %v, want forwarding-key site", sites)
	}

	var varNames []string
	for _, r := range refs.Vars {
		varNames = append(varNames, r.Name)
	}
	if strings.Join(varNames, ",") != "RELEASE_ENABLED,USE_PROD" {
		t.Errorf("vars = %v, want [RELEASE_ENABLED USE_PROD]", varNames)
	}

	// Plain-text mention outside ${{ }} must not be collected.
	if _, ok := gotSecrets["NOT_AN_EXPRESSION"]; ok {
		t.Error("collected a secrets. mention outside an expression")
	}

	// A workflow with no expressions reports Empty.
	if !ScanReferences(&Workflow{Jobs: []Job{{ID: "x", Steps: []Step{{Run: "make build"}}}}}).Empty() {
		t.Error("expected Empty references for expression-free workflow")
	}
}

// TestScanReferencesEnv covers references in workflow-level and job-level env: blocks,
// a very common site for secret/var usage.
func TestScanReferencesEnv(t *testing.T) {
	w := &Workflow{
		Name: "CI",
		Env:  []KV{{Key: "GLOBAL_TOKEN", Value: "${{ secrets.NPM_TOKEN }}"}},
		Jobs: []Job{
			{
				ID:  "build",
				Env: []KV{{Key: "API_KEY", Value: "${{ secrets.SERVICE_KEY }}"}},
			},
		},
	}

	refs := ScanReferences(w)
	got := map[string][]string{}
	for _, r := range refs.Secrets {
		got[r.Name] = r.Sites
	}

	if sites, ok := got["NPM_TOKEN"]; !ok {
		t.Error("workflow-level env reference NPM_TOKEN was not collected")
	} else if len(sites) != 1 || !strings.Contains(sites[0], "workflow env") {
		t.Errorf("NPM_TOKEN sites = %v, want a workflow env site", sites)
	}

	if sites, ok := got["SERVICE_KEY"]; !ok {
		t.Error("job-level env reference SERVICE_KEY was not collected")
	} else if len(sites) != 1 || !strings.Contains(sites[0], "env `API_KEY`") {
		t.Errorf("SERVICE_KEY sites = %v, want a job env site", sites)
	}
}

// TestContextRefs covers boundary handling in the expression-body scanner.
func TestContextRefs(t *testing.T) {
	tests := []struct {
		body, ctx string
		want      []string
	}{
		{"secrets.A || secrets.B", "secrets", []string{"A", "B"}},
		{"mysecrets.X", "secrets", nil}, // boundary: not a secrets. ref
		{"vars.X && secrets.TOKEN", "vars", []string{"X"}},
		{"toJSON(secrets)", "secrets", nil}, // bare context, no member access
	}
	for _, tt := range tests {
		got := contextRefs(tt.body, tt.ctx)
		if strings.Join(got, ",") != strings.Join(tt.want, ",") {
			t.Errorf("contextRefs(%q, %q) = %v, want %v", tt.body, tt.ctx, got, tt.want)
		}
	}
}

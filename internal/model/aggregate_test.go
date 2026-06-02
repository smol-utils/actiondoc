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
						// A secret referenced ONLY in a step env: block must still be
						// inventoried (the most common omission in real workflows).
						Name: "Sign image",
						Env: []KV{
							{Key: "BUILDKIT", Value: "1"},
							{Key: "COSIGN_PASSWORD", Value: "${{ secrets.COSIGN_PASSWORD }}"},
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

	wantOrder := []string{"SDKMAN_KEY", "GPG_KEY", "GPG_PASSPHRASE", "PROD_TOKEN", "DEV_TOKEN", "COSIGN_PASSWORD"}
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
	// Step-env-only secret is collected and its site labels the env key.
	if sites := gotSecrets["COSIGN_PASSWORD"]; len(sites) != 1 || !strings.Contains(sites[0], "env `COSIGN_PASSWORD`") {
		t.Errorf("COSIGN_PASSWORD sites = %v, want a step env site", sites)
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

// TestScanReferencesBareCondition covers if: conditions that reference secrets/vars
// WITHOUT ${{ }} delimiters -- valid because a condition is itself an expression. Both
// job-level and step-level if: must be scanned this way.
func TestScanReferencesBareCondition(t *testing.T) {
	w := &Workflow{
		Name: "CI",
		Jobs: []Job{
			{
				ID: "deploy",
				If: "secrets.DEPLOY_TOKEN != '' && vars.ENABLED",
				Steps: []Step{
					{Name: "gate", If: "secrets.STEP_GATE == 'yes'"},
				},
			},
		},
	}

	refs := ScanReferences(w)
	secrets := map[string]bool{}
	for _, r := range refs.Secrets {
		secrets[r.Name] = true
	}
	vars := map[string]bool{}
	for _, r := range refs.Vars {
		vars[r.Name] = true
	}

	for _, name := range []string{"DEPLOY_TOKEN", "STEP_GATE"} {
		if !secrets[name] {
			t.Errorf("bare-condition secret %q not collected", name)
		}
	}
	if !vars["ENABLED"] {
		t.Error("bare-condition var ENABLED not collected")
	}
}

// TestScanReferencesRunIsNotExpression guards the inverse: a bare secrets.X mention in a
// run: shell script is plain text, not an expression, and must NOT be collected.
func TestScanReferencesRunIsNotExpression(t *testing.T) {
	w := &Workflow{
		Jobs: []Job{{ID: "x", Steps: []Step{
			{Run: "echo do not collect secrets.NOT_A_REF here"},
		}}},
	}
	if !ScanReferences(w).Empty() {
		t.Error("collected a bare secrets. mention from a run: script (should be plain text)")
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
		{"toJSON(secrets)", "secrets", nil},                                                                       // bare context, no member access
		{"secrets.registry-0-usr", "secrets", []string{"registry-0-usr"}},                                         // hyphenated name captured whole
		{"steps.vault-secrets.outputs.TOKEN", "secrets", nil},                                                     // hyphen boundary: step output, not the secrets context
		{"steps.secrets.outputs.github-token", "secrets", nil},                                                    // dot boundary: step with id "secrets", not the secrets context
		{"x != 'a' && secrets.GH_TOKEN || steps.vault-secrets.outputs.GH_TOKEN", "secrets", []string{"GH_TOKEN"}}, // real ref kept, phantom rejected
	}
	for _, tt := range tests {
		got := contextRefs(tt.body, tt.ctx)
		if strings.Join(got, ",") != strings.Join(tt.want, ",") {
			t.Errorf("contextRefs(%q, %q) = %v, want %v", tt.body, tt.ctx, got, tt.want)
		}
	}
}

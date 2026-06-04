package redact

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/smol-utils/actiondoc/internal/callgraph"
	"github.com/smol-utils/actiondoc/internal/model"
)

// wf wraps a workflow as a single-source slice for Apply.
func wf(w *model.Workflow) []callgraph.Source {
	return []callgraph.Source{{Path: "wf.yml", Workflow: w}}
}

func TestApply_IdentifiersAreConsistentAndDeterministic(t *testing.T) {
	// The same secret name in a condition, a run script, and a forwarded secrets map must
	// all collapse to one placeholder, and numbering must follow sorted order regardless of
	// first-seen order.
	w := &model.Workflow{
		Name: "W",
		Jobs: []model.Job{{
			ID: "deploy",
			If: "${{ secrets.ZEBRA != '' }}",
			Steps: []model.Step{{
				Name: "s",
				Run:  "echo ${{ secrets.ALPHA }} ${{ secrets.ZEBRA }} ${{ vars.BETA }}",
			}},
		}},
	}
	m := Apply(wf(w), Options{})

	if got := w.Jobs[0].If; got != "${{ secrets.SECRET_2 != '' }}" {
		t.Errorf("if: got %q", got)
	}
	// ALPHA sorts before ZEBRA -> SECRET_1, SECRET_2.
	if got := w.Jobs[0].Steps[0].Run; got != "echo ${{ secrets.SECRET_1 }} ${{ secrets.SECRET_2 }} ${{ vars.VAR_1 }}" {
		t.Errorf("run: got %q", got)
	}
	want := map[string]string{
		"SECRET_1": "ALPHA",
		"SECRET_2": "ZEBRA",
		"VAR_1":    "BETA",
	}
	if !reflect.DeepEqual(m.entries, want) {
		t.Errorf("mapping: got %v want %v", m.entries, want)
	}
}

func TestApply_StableAcrossRuns(t *testing.T) {
	build := func() *model.Workflow {
		return &model.Workflow{Jobs: []model.Job{{
			ID:    "j",
			Steps: []model.Step{{Run: "echo ${{ secrets.B }} ${{ secrets.A }}"}},
		}}}
	}
	w1, w2 := build(), build()
	m1 := Apply(wf(w1), Options{})
	m2 := Apply(wf(w2), Options{})
	if !reflect.DeepEqual(m1.entries, m2.entries) {
		t.Errorf("nondeterministic mapping: %v vs %v", m1.entries, m2.entries)
	}
	if w1.Jobs[0].Steps[0].Run != w2.Jobs[0].Steps[0].Run {
		t.Errorf("nondeterministic rewrite")
	}
}

func TestApply_EnvKeyAndRefShareMapping(t *testing.T) {
	// An env: key and a ${{ env.KEY }} reference to it must map to the same placeholder.
	w := &model.Workflow{
		Env: []model.KV{{Key: "DEPLOY_HOST", Value: "x"}},
		Jobs: []model.Job{{
			ID:    "j",
			Steps: []model.Step{{Run: "echo ${{ env.DEPLOY_HOST }}"}},
		}},
	}
	Apply(wf(w), Options{})
	if w.Env[0].Key != "ENV_1" {
		t.Errorf("env key: got %q", w.Env[0].Key)
	}
	if got := w.Jobs[0].Steps[0].Run; got != "echo ${{ env.ENV_1 }}" {
		t.Errorf("env ref: got %q", got)
	}
}

func TestApply_ExpressionStructurePreserved(t *testing.T) {
	// A value mixing a literal, a redactable secret, and a non-redactable context must keep
	// its ${{ }} structure and leave github.* untouched.
	w := &model.Workflow{Jobs: []model.Job{{
		ID: "j",
		Steps: []model.Step{{
			Env: []model.KV{{Key: "K", Value: "prefix-${{ secrets.TOK }}-${{ github.sha }}"}},
		}},
	}}}
	Apply(wf(w), Options{})
	if got := w.Jobs[0].Steps[0].Env[0].Value; got != "prefix-${{ secrets.SECRET_1 }}-${{ github.sha }}" {
		t.Errorf("got %q", got)
	}
}

func TestApply_BracketNotationIdentifiers(t *testing.T) {
	// GitHub requires bracket notation for names with characters outside the identifier
	// set (hyphens are common in workflow_call keys). Both quote styles must be redacted,
	// the bracket shape preserved, and a hyphenated name must not leak.
	w := &model.Workflow{Jobs: []model.Job{{
		ID: "j",
		Steps: []model.Step{{
			Run: "a ${{ secrets['my-secret'] }} b ${{ vars[\"REGION-WEST\"] }} c ${{ secrets.PLAIN }}",
		}},
	}}}
	m := Apply(wf(w), Options{})
	got := w.Jobs[0].Steps[0].Run
	for _, leak := range []string{"my-secret", "REGION-WEST"} {
		if strings.Contains(got, leak) {
			t.Errorf("bracket-form name leaked (%q) in: %q", leak, got)
		}
	}
	// Numbering follows sorted order: "PLAIN" sorts before "my-secret" (byte order).
	want := "a ${{ secrets['SECRET_2'] }} b ${{ vars[\"VAR_1\"] }} c ${{ secrets.SECRET_1 }}"
	if got != want {
		t.Errorf("got  %q\nwant %q", got, want)
	}
	if m.entries["SECRET_2"] != "my-secret" || m.entries["VAR_1"] != "REGION-WEST" {
		t.Errorf("mapping wrong: %v", m.entries)
	}
}

func TestApply_HostsAndURLs(t *testing.T) {
	w := &model.Workflow{
		Description: "see https://docs.corp.example/x and deploy.internal.corp",
		Jobs: []model.Job{{
			ID: "j",
			Steps: []model.Step{{
				Run: "curl https://api.internal.corp/v1; ping registry.internal.corp",
			}},
		}},
	}
	m := Apply(wf(w), Options{})
	// URLs and hosts each get their own category and numbering.
	if got := w.Description; got != "see URL_2 and HOST_1" {
		t.Errorf("desc: got %q", got)
	}
	if got := w.Jobs[0].Steps[0].Run; got != "curl URL_1; ping HOST_2" {
		t.Errorf("run: got %q", got)
	}
	if m.entries["URL_1"] != "https://api.internal.corp/v1" {
		t.Errorf("url map: %v", m.entries)
	}
}

func TestApply_HostWithContextLikeLabelFullyRedacted(t *testing.T) {
	// A literal hostname whose leading label matches an expression context (secrets./vars./
	// env.) must be redacted whole, not corrupted by identifier rewriting. Here "example" is
	// also a real secret, the worst case: the host substitution must still win on literal
	// text, while the genuine ${{ secrets.example }} reference is rewritten as an identifier.
	w := &model.Workflow{Jobs: []model.Job{{
		ID: "j",
		Steps: []model.Step{{
			Run: "echo ${{ secrets.example }}; curl secrets.example.com; curl https://secrets.example.com/x",
		}},
	}}}
	Apply(wf(w), Options{})
	got := w.Jobs[0].Steps[0].Run
	for _, leak := range []string{"example.com", ".com"} {
		if strings.Contains(got, leak) {
			t.Errorf("host left partly unredacted (%q) in: %q", leak, got)
		}
	}
	if !strings.Contains(got, "${{ secrets.SECRET_1 }}") {
		t.Errorf("genuine secret reference not redacted: %q", got)
	}
}

func TestApply_FilenameNotTreatedAsHost(t *testing.T) {
	w := &model.Workflow{Description: "edit config.yaml then run deploy.sh"}
	Apply(wf(w), Options{})
	if w.Description != "edit config.yaml then run deploy.sh" {
		t.Errorf("filenames should be untouched, got %q", w.Description)
	}
}

func TestApply_RunnerLabels(t *testing.T) {
	w := &model.Workflow{Jobs: []model.Job{
		{ID: "hosted", RunsOn: "ubuntu-latest"},
		{ID: "self", RunsOn: "self-hosted, linux, x64, prod-gpu"},
		{ID: "matrix", RunsOn: "${{ matrix.os }}"},
	}}
	Apply(wf(w), Options{})
	if w.Jobs[0].RunsOn != "ubuntu-latest" {
		t.Errorf("hosted runner should be kept, got %q", w.Jobs[0].RunsOn)
	}
	if got := w.Jobs[1].RunsOn; got != "self-hosted, linux, x64, RUNNER_1" {
		t.Errorf("self-hosted: got %q", got)
	}
	if w.Jobs[2].RunsOn != "${{ matrix.os }}" {
		t.Errorf("matrix expr should be preserved, got %q", w.Jobs[2].RunsOn)
	}
}

func TestApply_EnvironmentName(t *testing.T) {
	w := &model.Workflow{Jobs: []model.Job{{
		ID:          "j",
		Environment: &model.Environment{Name: "production-east", URL: "https://prod.corp.example"},
	}}}
	Apply(wf(w), Options{})
	if w.Jobs[0].Environment.Name != "ENVIRONMENT_1" {
		t.Errorf("env name: got %q", w.Jobs[0].Environment.Name)
	}
	if w.Jobs[0].Environment.URL != "URL_1" {
		t.Errorf("env url: got %q", w.Jobs[0].Environment.URL)
	}
}

func TestApply_RationaleFreeTextRedacted(t *testing.T) {
	// Free text in permission and cron rationale comments must be swept for hosts/URLs.
	w := &model.Workflow{
		Permissions: &model.Permissions{Scopes: []model.Permission{
			{Scope: "contents", Level: "read", Rationale: "mirror from registry.internal.corp"},
		}},
		Triggers: &model.Triggers{Schedule: []model.CronEntry{
			{Cron: "0 0 * * *", Rationale: "nightly sync to backup.internal.corp"},
		}},
	}
	Apply(wf(w), Options{})
	if got := w.Permissions.Scopes[0].Rationale; got != "mirror from HOST_2" {
		t.Errorf("permission rationale: got %q", got)
	}
	if got := w.Triggers.Schedule[0].Rationale; got != "nightly sync to HOST_1" {
		t.Errorf("cron rationale: got %q", got)
	}
	// Scope/level are not sensitive and stay intact.
	if w.Permissions.Scopes[0].Scope != "contents" || w.Permissions.Scopes[0].Level != "read" {
		t.Errorf("permission scope/level should be untouched: %+v", w.Permissions.Scopes[0])
	}
}

func TestApply_GitHubTokenKept(t *testing.T) {
	w := &model.Workflow{Jobs: []model.Job{{
		ID:    "j",
		Steps: []model.Step{{Run: "echo ${{ secrets.GITHUB_TOKEN }} ${{ secrets.OTHER }}"}},
	}}}
	Apply(wf(w), Options{})
	if got := w.Jobs[0].Steps[0].Run; got != "echo ${{ secrets.GITHUB_TOKEN }} ${{ secrets.SECRET_1 }}" {
		t.Errorf("GITHUB_TOKEN should be kept, got %q", got)
	}
}

func TestApply_ConservativeKeepsHarmlessLiterals(t *testing.T) {
	w := &model.Workflow{Jobs: []model.Job{{
		ID:    "j",
		Steps: []model.Step{{With: []model.KV{{Key: "node-version", Value: "20"}}}},
	}}}
	Apply(wf(w), Options{Level: Conservative})
	if w.Jobs[0].Steps[0].With[0].Value != "20" {
		t.Errorf("conservative should keep harmless literal, got %q", w.Jobs[0].Steps[0].With[0].Value)
	}
}

func TestApply_AggressiveRedactsLiterals(t *testing.T) {
	w := &model.Workflow{Jobs: []model.Job{{
		ID: "j",
		Steps: []model.Step{{With: []model.KV{
			{Key: "node-version", Value: "20"},
			{Key: "token", Value: "${{ secrets.TOK }}"},
		}}},
	}}}
	Apply(wf(w), Options{Level: Aggressive})
	if got := w.Jobs[0].Steps[0].With[0].Value; got != "VALUE_1" {
		t.Errorf("aggressive should redact literal, got %q", got)
	}
	// Expression-bearing values keep their structure even under Aggressive.
	if got := w.Jobs[0].Steps[0].With[1].Value; got != "${{ secrets.SECRET_1 }}" {
		t.Errorf("aggressive should preserve expression, got %q", got)
	}
}

func TestApply_TagSecretNameRedacted(t *testing.T) {
	w := &model.Workflow{Tags: model.Tags{
		Secrets: []model.Param{{Name: "SONATYPE_PW_X", Description: "for nexus.corp.example"}},
	}}
	m := Apply(wf(w), Options{})
	if w.Tags.Secrets[0].Name != "SECRET_1" {
		t.Errorf("tag secret name: got %q", w.Tags.Secrets[0].Name)
	}
	if w.Tags.Secrets[0].Description != "for HOST_1" {
		t.Errorf("tag secret desc: got %q", w.Tags.Secrets[0].Description)
	}
	if m.entries["SECRET_1"] != "SONATYPE_PW_X" {
		t.Errorf("mapping: %v", m.entries)
	}
}

func TestApply_ActionInputDefaultsAreFreeText(t *testing.T) {
	// An action.yml input default is public contract, like a workflow_call input default:
	// even under the aggressive profile it must not be blanked to VALUE_n, but a host
	// inside it is still redacted.
	a := &model.Action{
		Name: "A",
		Inputs: []model.ActionInput{
			{Name: "node-version", Default: "20"},
			{Name: "registry", Default: "registry.internal.corp"},
		},
	}
	sources := []callgraph.Source{{Path: "action.yml", Action: a}}
	Apply(sources, Options{Level: Aggressive})
	if a.Inputs[0].Default != "20" {
		t.Errorf("harmless default should be preserved, got %q", a.Inputs[0].Default)
	}
	if a.Inputs[1].Default != "HOST_1" {
		t.Errorf("host in default should be redacted as a host, got %q", a.Inputs[1].Default)
	}
}

func TestMappingJSON(t *testing.T) {
	w := &model.Workflow{Jobs: []model.Job{{ID: "j", Steps: []model.Step{{Run: "${{ secrets.A }}"}}}}}
	m := Apply(wf(w), Options{})
	data, err := m.JSON()
	if err != nil {
		t.Fatal(err)
	}
	var got map[string]string
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if got["SECRET_1"] != "A" {
		t.Errorf("json map: %v", got)
	}
}

func TestMappingJSON_EmptyWhenNothingMatched(t *testing.T) {
	w := &model.Workflow{Name: "Plain", Jobs: []model.Job{{ID: "j", RunsOn: "ubuntu-latest"}}}
	m := Apply(wf(w), Options{})
	if !m.Empty() {
		t.Errorf("expected empty mapping, got %v", m.entries)
	}
	data, err := m.JSON()
	if err != nil {
		t.Fatal(err)
	}
	if data != nil {
		t.Errorf("expected nil JSON for empty mapping, got %s", data)
	}
}

func TestApply_CrossSourceConsistency(t *testing.T) {
	// The same secret used in two different files must map to the same placeholder.
	a := &model.Workflow{Jobs: []model.Job{{ID: "a", Steps: []model.Step{{Run: "${{ secrets.SHARED }}"}}}}}
	b := &model.Workflow{Jobs: []model.Job{{ID: "b", Steps: []model.Step{{Run: "${{ secrets.SHARED }}"}}}}}
	sources := []callgraph.Source{{Path: "a.yml", Workflow: a}, {Path: "b.yml", Workflow: b}}
	Apply(sources, Options{})
	if a.Jobs[0].Steps[0].Run != b.Jobs[0].Steps[0].Run {
		t.Errorf("cross-source inconsistency: %q vs %q", a.Jobs[0].Steps[0].Run, b.Jobs[0].Steps[0].Run)
	}
}

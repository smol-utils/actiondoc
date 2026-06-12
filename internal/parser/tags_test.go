package parser

import (
	"testing"
)

func TestParseTagsEmpty(t *testing.T) {
	tags := ParseTags("")
	if tags.Desc != "" {
		t.Errorf("expected empty desc, got %q", tags.Desc)
	}
}

func TestParseTagsSimpleDesc(t *testing.T) {
	comment := "# @desc Build the application."
	tags := ParseTags(comment)
	if tags.Desc != "Build the application." {
		t.Errorf("desc = %q, want %q", tags.Desc, "Build the application.")
	}
}

func TestParseTagsMultilineDesc(t *testing.T) {
	comment := "# @desc Build the application\n# and run all tests."
	tags := ParseTags(comment)
	want := "Build the application\nand run all tests."
	if tags.Desc != want {
		t.Errorf("desc = %q, want %q", tags.Desc, want)
	}
}

func TestParseTagsSecret(t *testing.T) {
	comment := "# @secret NPM_TOKEN - Required for private packages"
	tags := ParseTags(comment)
	if len(tags.Secrets) != 1 {
		t.Fatalf("expected 1 secret, got %d", len(tags.Secrets))
	}
	s := tags.Secrets[0]
	if s.Name != "NPM_TOKEN" {
		t.Errorf("name = %q, want %q", s.Name, "NPM_TOKEN")
	}
	if s.Description != "Required for private packages" {
		t.Errorf("description = %q, want %q", s.Description, "Required for private packages")
	}
}

func TestParseTagsInputWithType(t *testing.T) {
	comment := "# @input deploy {boolean} - Whether to deploy after build"
	tags := ParseTags(comment)
	if len(tags.Inputs) != 1 {
		t.Fatalf("expected 1 input, got %d", len(tags.Inputs))
	}
	p := tags.Inputs[0]
	if p.Name != "deploy" {
		t.Errorf("name = %q", p.Name)
	}
	if p.Type != "boolean" {
		t.Errorf("type = %q", p.Type)
	}
	if p.Description != "Whether to deploy after build" {
		t.Errorf("description = %q", p.Description)
	}
}

func TestParseTagsMultipleTags(t *testing.T) {
	comment := "# @desc Deploy pipeline.\n# @secret DEPLOY_KEY - SSH key\n# @since v2.0.0\n# @see https://example.com"
	tags := ParseTags(comment)

	if tags.Desc != "Deploy pipeline." {
		t.Errorf("desc = %q", tags.Desc)
	}
	if len(tags.Secrets) != 1 || tags.Secrets[0].Name != "DEPLOY_KEY" {
		t.Errorf("secrets = %+v", tags.Secrets)
	}
	if tags.Since != "v2.0.0" {
		t.Errorf("since = %q", tags.Since)
	}
	if len(tags.See) != 1 || tags.See[0] != "https://example.com" {
		t.Errorf("see = %+v", tags.See)
	}
}

func TestParseTagsDeprecated(t *testing.T) {
	comment := "# @deprecated Use build-v2 instead."
	tags := ParseTags(comment)
	if tags.Deprecated != "Use build-v2 instead." {
		t.Errorf("deprecated = %q", tags.Deprecated)
	}
}

func TestParseTagsExample(t *testing.T) {
	comment := "# @example\n#   gh workflow run ci.yml\n#   gh workflow run ci.yml -f deploy=true"
	tags := ParseTags(comment)
	want := "  gh workflow run ci.yml\n  gh workflow run ci.yml -f deploy=true"
	if tags.Example != want {
		t.Errorf("example = %q, want %q", tags.Example, want)
	}
}

func TestParseTagsOutputWithType(t *testing.T) {
	comment := "# @output image-tag {semver} - The Docker image tag"
	tags := ParseTags(comment)
	if len(tags.Outputs) != 1 {
		t.Fatalf("expected 1 output, got %d", len(tags.Outputs))
	}
	o := tags.Outputs[0]
	if o.Name != "image-tag" || o.Type != "semver" || o.Description != "The Docker image tag" {
		t.Errorf("output = %+v", o)
	}
}

func TestParseTagsEnv(t *testing.T) {
	comment := "# @env DATABASE_URL - Connection string for test DB"
	tags := ParseTags(comment)
	if len(tags.Envs) != 1 {
		t.Fatalf("expected 1 env, got %d", len(tags.Envs))
	}
	if tags.Envs[0].Name != "DATABASE_URL" {
		t.Errorf("name = %q", tags.Envs[0].Name)
	}
}

func TestParseTagsContinuationIndentation(t *testing.T) {
	// Continuation lines are commonly indented under their tag for readability; that
	// source indentation must not leak into the value (it surfaces as '<br>   ' in
	// rendered tables). @example is the exception: its indentation is content.
	comment := "# @desc Deploys the application\n" +
		"#    to the staging environment\n" +
		"#      and verifies the rollout.\n" +
		"# @env DEPLOY_HOST - target host,\n" +
		"#    used by the deploy script"

	tags := ParseTags(comment)

	wantDesc := "Deploys the application\nto the staging environment\nand verifies the rollout."
	if tags.Desc != wantDesc {
		t.Errorf("Desc = %q, want %q", tags.Desc, wantDesc)
	}
	wantEnvDesc := "target host,\nused by the deploy script"
	if len(tags.Envs) != 1 || tags.Envs[0].Description != wantEnvDesc {
		t.Errorf("Envs = %+v, want description %q", tags.Envs, wantEnvDesc)
	}
}

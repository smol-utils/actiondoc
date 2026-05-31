package parser

import "testing"

// TestParseTagsAllowlist locks the fixed @-tag allowlist: unknown @-prefixed strings
// must be ignored entirely while interleaved allowlisted tags still parse. This lets
// ActionDoc coexist with other tooling that embeds its own @-markers in YAML comments.
func TestParseTagsAllowlist(t *testing.T) {
	comment := `@lulaStart compliance-block
@desc Deploys to production
@actiondoc-unknown should be ignored
@secret DEPLOY_KEY - production deploy key
@lulaEnd`

	tags := ParseTags(comment)

	if tags.Desc != "Deploys to production" {
		t.Errorf("Desc = %q, want %q", tags.Desc, "Deploys to production")
	}
	if len(tags.Secrets) != 1 || tags.Secrets[0].Name != "DEPLOY_KEY" {
		t.Errorf("Secrets = %+v, want one entry DEPLOY_KEY", tags.Secrets)
	}
	// Unknown tags must not leak into any field.
	if got := tags.Desc; got != "Deploys to production" {
		t.Errorf("unknown @-tags leaked into Desc: %q", got)
	}
}

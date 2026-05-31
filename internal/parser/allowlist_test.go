package parser

import "testing"

// TestParseTagsAllowlist locks the fixed @-tag allowlist (spec: "Tag Allowlist"):
// unknown @-prefixed strings must be ignored entirely while interleaved allowlisted
// tags still parse. Driver: defenseunicorns/uds-core embeds Lula markers
// (@lulaStart/@lulaEnd) in YAML comments alongside ActionDoc tags.
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

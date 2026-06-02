package model

// Permissions models a GitHub Actions permissions: block at the workflow or job level.
// A block can take three shapes: a scalar grant (All, e.g. read-all/write-all), an
// explicit empty block (DefaultDeny, written permissions: {}), or an ordered list of
// scope->level Scopes. Source order is preserved for deterministic rendering.
type Permissions struct {
	All         string       `json:"all,omitempty"`          // scalar form: read-all / write-all
	DefaultDeny bool         `json:"default_deny,omitempty"` // permissions: {} -- grants nothing
	Scopes      []Permission `json:"scopes,omitempty"`
}

// Permission is one scope->level grant within a permissions: block, with an optional
// rationale from a trailing comment (e.g. `contents: read  # for actions/checkout`).
// OIDC is set when the grant is id-token: write, the keyless-signing enabler worth
// calling out for security review.
type Permission struct {
	Scope     string `json:"scope"`
	Level     string `json:"level"`
	Rationale string `json:"rationale,omitempty"`
	OIDC      bool   `json:"oidc,omitempty"`
}

// Environment is a job's GitHub Environments binding -- the deploy gate. When a job
// declares environment: <name>, GitHub enforces the environment's protection rules
// (required reviewers, wait timers, branch policies) and exposes its scoped secrets.
// Those rules live in repository settings, not the workflow YAML, so they are not
// represented here; only the binding (and optional URL) is. Name may be an expression
// and is preserved verbatim.
type Environment struct {
	Name string `json:"name"`
	URL  string `json:"url,omitempty"`
}

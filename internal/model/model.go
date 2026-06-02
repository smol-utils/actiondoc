package model

// Workflow is the top-level IR for a single workflow file.
type Workflow struct {
	File        string   `json:"file"`
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	On          []string `json:"on"`
	Jobs        []Job    `json:"jobs"`
	Tags        Tags     `json:"tags,omitempty"`

	// Declared surface beyond the flat On list (see triggers.go, security.go, metadata.go).
	Triggers    *Triggers    `json:"triggers,omitempty"`
	Permissions *Permissions `json:"permissions,omitempty"`
	Env         []KV         `json:"env,omitempty"`
	Concurrency *Concurrency `json:"concurrency,omitempty"`
	Defaults    *Defaults    `json:"defaults,omitempty"`
}

// Job represents a single job in the workflow.
type Job struct {
	ID          string   `json:"id"`
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	RunsOn      string   `json:"runs_on"`
	Needs       []string `json:"needs,omitempty"`
	If          string   `json:"if,omitempty"`
	Steps       []Step   `json:"steps"`
	Tags        Tags     `json:"tags,omitempty"`

	// Reusable-workflow caller fields. A job that calls a reusable workflow uses
	// `uses:` instead of `runs-on:`/`steps:`.
	Uses           string `json:"uses,omitempty"`            // reusable workflow call target (raw `uses:` value)
	With           []KV   `json:"with,omitempty"`            // forwarded inputs (with:)
	Secrets        []KV   `json:"secrets,omitempty"`         // forwarded secrets (explicit mapping)
	SecretsInherit bool   `json:"secrets_inherit,omitempty"` // `secrets: inherit`

	// Declared job-level surface (see security.go, metadata.go).
	Permissions *Permissions `json:"permissions,omitempty"`
	Env         []KV         `json:"env,omitempty"`
	Concurrency *Concurrency `json:"concurrency,omitempty"`
	Defaults    *Defaults    `json:"defaults,omitempty"`
	Environment *Environment `json:"environment,omitempty"`

	// Matrix holds statically-resolvable strategy.matrix axes (axis name -> value list),
	// used to expand ${{ matrix.X }} references in the job name. Empty when the matrix
	// uses include/exclude or a dynamic fromJSON source, in which case names render verbatim.
	Matrix []MatrixAxis `json:"matrix,omitempty"`
}

// Step represents a single step within a job.
type Step struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Uses        string `json:"uses,omitempty"`
	Run         string `json:"run,omitempty"`
	If          string `json:"if,omitempty"`
	Tags        Tags   `json:"tags,omitempty"`

	// With holds the step's `with:` inputs and Env its `env:` variables, both in source
	// order. ContinueOnError records a literal `continue-on-error: true`;
	// ContinueOnErrorExpr holds the raw value instead when it is an expression (e.g.
	// `continue-on-error: ${{ matrix.experimental }}`), so the step is still flagged as
	// failure-tolerant. UsesVersion is the trailing version comment on a SHA-pinned
	// `uses:` ref (e.g. the "v4" in `uses: actions/checkout@<sha> # v4`).
	With                []KV   `json:"with,omitempty"`
	Env                 []KV   `json:"env,omitempty"`
	ContinueOnError     bool   `json:"continue_on_error,omitempty"`
	ContinueOnErrorExpr string `json:"continue_on_error_expr,omitempty"`
	UsesVersion         string `json:"uses_version,omitempty"`

	// UsesAction is the parsed local composite action that Uses points to, when that action
	// was discovered in the scan set; lets the renderer pair `with:` keys with declared inputs.
	UsesAction *Action `json:"-"`
}

// Tags holds ActionDoc comment annotations.
type Tags struct {
	Desc       string   `json:"desc,omitempty"`
	Secrets    []Param  `json:"secrets,omitempty"`
	Inputs     []Param  `json:"inputs,omitempty"`
	Envs       []Param  `json:"envs,omitempty"`
	Outputs    []Param  `json:"outputs,omitempty"`
	Deprecated string   `json:"deprecated,omitempty"`
	See        []string `json:"see,omitempty"`
	Since      string   `json:"since,omitempty"`
	Example    string   `json:"example,omitempty"`
}

// Param is a named parameter with optional type hint and description.
type Param struct {
	Name        string `json:"name"`
	Type        string `json:"type,omitempty"`
	Description string `json:"description,omitempty"`
}

// KV is an ordered key/value pair, used for reusable-workflow `with:`/`secrets:`
// forwarding maps where source order is preserved for deterministic rendering.
type KV struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Action is the top-level data model for a GitHub Action metadata file (action.yml).
type Action struct {
	File        string         `json:"file"`
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Inputs      []ActionInput  `json:"inputs,omitempty"`
	Outputs     []ActionOutput `json:"outputs,omitempty"`
	Runs        ActionRuns     `json:"runs"`
	Branding    *Branding      `json:"branding,omitempty"`
	Tags        Tags           `json:"tags,omitempty"`
}

// ActionInput is a single input defined in action.yml.
type ActionInput struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Required    bool   `json:"required"`
	Default     string `json:"default,omitempty"`
}

// ActionOutput is a single output defined in action.yml.
type ActionOutput struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// ActionRuns describes how the action executes.
type ActionRuns struct {
	Using string `json:"using"`
	Main  string `json:"main,omitempty"`
	Image string `json:"image,omitempty"`
}

// Branding holds optional branding metadata.
type Branding struct {
	Icon  string `json:"icon,omitempty"`
	Color string `json:"color,omitempty"`
}

package model

// Workflow is the top-level IR for a single workflow file.
type Workflow struct {
	File        string   `json:"file"`
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	On          []string `json:"on"`
	Jobs        []Job    `json:"jobs"`
	Tags        Tags     `json:"tags,omitempty"`
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

	// --- v0.2.0 M1: reusable-workflow caller fields ---
	// A job that calls a reusable workflow uses `uses:` instead of `runs-on:`/`steps:`.
	Uses           string `json:"uses,omitempty"`            // reusable workflow call target (raw `uses:` value)
	With           []KV   `json:"with,omitempty"`            // forwarded inputs (with:)
	Secrets        []KV   `json:"secrets,omitempty"`         // forwarded secrets (explicit mapping)
	SecretsInherit bool   `json:"secrets_inherit,omitempty"` // `secrets: inherit`
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

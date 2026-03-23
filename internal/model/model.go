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

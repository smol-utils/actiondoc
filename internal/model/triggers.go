package model

// Triggers captures a workflow's declared trigger surface beyond the flat event-name
// list in Workflow.On. It is intentionally a thin, generic container: it surfaces the
// pieces of the on: block that carry a documented contract (workflow_dispatch inputs,
// the workflow_call public API, schedules, and per-event filters) without mirroring
// GitHub's full on: grammar. Workflow.On remains the canonical flat trigger list.
type Triggers struct {
	Dispatch *DispatchTrigger `json:"workflow_dispatch,omitempty"`
	Call     *CallTrigger     `json:"workflow_call,omitempty"`
	Schedule []CronEntry      `json:"schedule,omitempty"`
	Events   []TriggerEvent   `json:"events,omitempty"`
}

// DispatchTrigger holds the inputs of a workflow_dispatch (manual) trigger.
type DispatchTrigger struct {
	Inputs []WorkflowInput `json:"inputs,omitempty"`
}

// CallTrigger holds the public API of a reusable workflow (workflow_call): the inputs,
// outputs, and secrets it declares as its caller contract.
type CallTrigger struct {
	Inputs  []WorkflowInput  `json:"inputs,omitempty"`
	Outputs []WorkflowOutput `json:"outputs,omitempty"`
	Secrets []WorkflowSecret `json:"secrets,omitempty"`
}

// WorkflowInput is a single declared input on a workflow_dispatch or workflow_call
// trigger. Names preserve their source casing and hyphenation (e.g. CONSUMER-KEY); they
// are never normalized. Options is populated only for type: choice inputs.
type WorkflowInput struct {
	Name        string   `json:"name"`
	Type        string   `json:"type,omitempty"`
	Required    bool     `json:"required,omitempty"`
	Default     string   `json:"default,omitempty"`
	Description string   `json:"description,omitempty"`
	Options     []string `json:"options,omitempty"`
}

// WorkflowOutput is a single declared workflow_call output.
type WorkflowOutput struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Value       string `json:"value,omitempty"`
}

// WorkflowSecret is a single declared workflow_call secret.
type WorkflowSecret struct {
	Name        string `json:"name"`
	Required    bool   `json:"required,omitempty"`
	Description string `json:"description,omitempty"`
}

// CronEntry is one on.schedule cron expression, preserved verbatim (named day-of-week
// such as THU is not normalized to numeric form), with an optional rationale taken from
// a trailing comment on the cron entry.
type CronEntry struct {
	Cron      string `json:"cron"`
	Rationale string `json:"rationale,omitempty"`
}

// TriggerEvent models any event under on: (push, pull_request, release,
// repository_dispatch, workflow_run, ...) generically as an ordered list of its filter
// sub-keys. It deliberately does not build a typed field per filter kind: a filter is
// just a key (types/branches/paths/branches-ignore/paths-ignore/...) and its value(s).
type TriggerEvent struct {
	Name    string          `json:"name"`
	Filters []TriggerFilter `json:"filters,omitempty"`
}

// TriggerFilter is one sub-key under an event, with its one-or-more values preserved
// in source order.
type TriggerFilter struct {
	Key    string   `json:"key"`
	Values []string `json:"values"`
}

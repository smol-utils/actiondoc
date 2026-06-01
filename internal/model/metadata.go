package model

// Concurrency models a concurrency: block (the release-serialization lock). Group may be
// an expression and is preserved verbatim; CancelInProgress holds the raw value
// ("true"/"false" or an expression) so it is never silently coerced.
type Concurrency struct {
	Group            string `json:"group"`
	CancelInProgress string `json:"cancel_in_progress,omitempty"`
}

// Defaults models a defaults.run: block. WorkingDirectory is load-bearing: it re-roots
// the relative paths in every run: step, so dropping it makes rendered script paths
// misleading.
type Defaults struct {
	Shell            string `json:"shell,omitempty"`
	WorkingDirectory string `json:"working_directory,omitempty"`
}

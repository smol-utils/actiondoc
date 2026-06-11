package redact

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/smol-utils/actiondoc/internal/callgraph"
	"github.com/smol-utils/actiondoc/internal/model"
)

// This test guards the redactor against the model growing a field it does not know
// about. The collect*/rewrite* walks enumerate model fields by hand, and a forgotten
// field fails open: its value ships in output the user asked to make shareable. Here
// reflection walks every string field reachable from the model schema, plants a unique
// redactable sentinel hostname in each, runs Apply, and fails naming any field whose
// sentinel survived -- unless that field is explicitly listed as intentionally readable.
//
// When this test fails after a model change, either wire the new field into the
// redactor's collect AND rewrite passes, or -- only if its content is structural
// (a name, an id, a ref, an enum) -- add it to intentionallyUnredacted below.

// intentionallyUnredacted lists the fields whose values are deliberately left readable,
// as dot-path suffixes (a suffix matches the same field at every nesting level).
// Structural names and ids keep the docs navigable; uses: refs and file paths preserve
// the call-graph shape; type/level/icon strings are enum-like metadata with no room for
// sensitive content.
var intentionallyUnredacted = []string{
	// Document identity and titles (titles are not redacted by design; see Apply).
	"Workflow.File", "Workflow.Name", "Workflow.On",
	"Action.File", "Action.Name",
	"Jobs.ID", "Jobs.Name", "Jobs.Needs",
	"Steps.ID", "Steps.Name",

	// Call-graph shape: uses: refs and the pinned-version comment.
	"Jobs.Uses", "Steps.Uses", "Steps.UsesVersion",

	// Declared parameter names are public contract (secret and env NAMES are redacted;
	// input/output names are not).
	"Inputs.Name", "Outputs.Name", "With.Key",
	"Tags.Inputs.Name", "Tags.Outputs.Name",

	// Enum-like metadata.
	"Inputs.Type",
	"Tags.Secrets.Type", "Tags.Inputs.Type", "Tags.Envs.Type", "Tags.Outputs.Type",
	"Tags.Since",
	"Permissions.All", "Scopes.Scope", "Scopes.Level",
	"Schedule.Cron",
	"Events.Name", "Filters.Key",
	"Matrix.Name",
	"Runs.Using", "Runs.Main",
	"Branding.Icon", "Branding.Color",
}

// sentinelMark is the recognizable core of every planted value. The full sentinel is a
// bare hostname (canary-N.redact.example) so any field swept as free text has its
// sentinel caught by host detection, and any field redacted as a whole-value lookup
// (env keys, secret names, environment names, runner labels) is replaced outright.
const sentinelMark = "canary-"

func sentinel(n int) string { return fmt.Sprintf("canary-%d.redact.example", n) }

// fillStrings walks v, allocating nil pointers and one element per slice, and sets every
// string field to a unique sentinel. Unexported fields and UsesAction are not entered:
// the linked action is a source of its own and is redacted as one.
func fillStrings(v reflect.Value, n *int) {
	switch v.Kind() {
	case reflect.Pointer:
		if v.IsNil() {
			if !v.CanSet() {
				return
			}
			v.Set(reflect.New(v.Type().Elem()))
		}
		fillStrings(v.Elem(), n)
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			f := v.Type().Field(i)
			if f.PkgPath != "" || f.Name == "UsesAction" {
				continue
			}
			fillStrings(v.Field(i), n)
		}
	case reflect.Slice:
		if v.IsNil() {
			v.Set(reflect.MakeSlice(v.Type(), 1, 1))
		}
		for i := 0; i < v.Len(); i++ {
			fillStrings(v.Index(i), n)
		}
	case reflect.String:
		(*n)++
		v.SetString(sentinel(*n))
	}
}

// collectSurvivors walks v the same way fillStrings does and records the dot path of
// every string field whose post-Apply value still carries a sentinel.
func collectSurvivors(v reflect.Value, path string, out map[string]bool) {
	switch v.Kind() {
	case reflect.Pointer:
		if !v.IsNil() {
			collectSurvivors(v.Elem(), path, out)
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			f := v.Type().Field(i)
			if f.PkgPath != "" || f.Name == "UsesAction" {
				continue
			}
			collectSurvivors(v.Field(i), path+"."+f.Name, out)
		}
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			collectSurvivors(v.Index(i), path, out)
		}
	case reflect.String:
		if strings.Contains(v.String(), sentinelMark) {
			out[path] = true
		}
	}
}

// matchesEntry reports whether an allowlist entry covers a field path: the entry is the
// whole path (root-level fields) or a trailing dot-path of it (nested fields).
func matchesEntry(path, entry string) bool {
	return path == entry || strings.HasSuffix(path, "."+entry)
}

func allowlisted(path string) bool {
	for _, suffix := range intentionallyUnredacted {
		if matchesEntry(path, suffix) {
			return true
		}
	}
	return false
}

func TestApply_CoversEveryModelStringField(t *testing.T) {
	w := &model.Workflow{}
	a := &model.Action{}
	n := 0
	fillStrings(reflect.ValueOf(w).Elem(), &n)
	fillStrings(reflect.ValueOf(a).Elem(), &n)

	Apply([]callgraph.Source{
		{Path: "wf.yml", Workflow: w},
		{Path: "action.yml", Action: a},
	}, Options{Level: Conservative})

	survivors := map[string]bool{}
	collectSurvivors(reflect.ValueOf(w).Elem(), "Workflow", survivors)
	collectSurvivors(reflect.ValueOf(a).Elem(), "Action", survivors)

	for path := range survivors {
		if !allowlisted(path) {
			t.Errorf("model field %s is not redacted and not listed as intentionally readable", path)
		}
	}
	// Keep the allowlist honest: every entry must still match a real field that survives.
	for _, suffix := range intentionallyUnredacted {
		used := false
		for path := range survivors {
			if matchesEntry(path, suffix) {
				used = true
				break
			}
		}
		if !used {
			t.Errorf("allowlist entry %q matches no surviving field; remove or update it", suffix)
		}
	}
}

package callgraph

import (
	"path/filepath"
	"testing"

	"github.com/smol-utils/actiondoc/internal/model"
)

func TestBuildBasic(t *testing.T) {
	entry := &model.Workflow{
		File: "release.yml", Name: "Release", On: []string{"workflow_dispatch"},
		Jobs: []model.Job{
			{ID: "publish", Uses: "./.github/workflows/build.yml",
				With:    []model.KV{{Key: "version", Value: "${{ inputs.version }}"}},
				Secrets: []model.KV{{Key: "CONSUMER-KEY", Value: "${{ secrets.SDKMAN_KEY }}"}}},
			{ID: "external", Uses: "owner/repo/.github/workflows/linux_job.yml@main"},
		},
	}
	reusable := &model.Workflow{File: "build.yml", Name: "Build", On: []string{"workflow_call"}}

	g := Build([]Source{
		{Path: "/repo/.github/workflows/release.yml", Workflow: entry},
		{Path: "/repo/.github/workflows/build.yml", Workflow: reusable},
	})

	relID := "/repo/.github/workflows/release.yml"
	buildID := "/repo/.github/workflows/build.yml"

	if !g.IsEntryPoint(relID) {
		t.Error("release.yml should be an entry point")
	}
	if g.IsEntryPoint(buildID) {
		t.Error("build.yml is workflow_call-only; not an entry point")
	}
	if got := g.CallerCount(buildID); got != 1 {
		t.Errorf("CallerCount(build.yml) = %d, want 1", got)
	}
	// Local reusable edge resolves to the build.yml node.
	calls := g.Calls(relID)
	var foundLocal, foundExternal bool
	for _, e := range calls {
		if e.ToID == buildID && e.Kind == KindReusable {
			foundLocal = true
		}
		if e.Pin == "main" {
			foundExternal = true
			if n := g.Nodes[e.ToID]; n == nil || !n.External {
				t.Error("cross-repo target should be an external node")
			}
		}
	}
	if !foundLocal {
		t.Error("expected resolved local reusable edge release.yml -> build.yml")
	}
	if !foundExternal {
		t.Error("expected external edge with @main pin preserved")
	}
	if r := g.Reachable(relID); len(r) == 0 {
		t.Error("release.yml should reach at least build.yml")
	}
}

// TestBuildCompositeResolution checks that a step-level local composite `uses:` resolves
// to the exact action directory on a path-segment boundary -- "build" must not match a
// sibling "my-build" -- and that the result is deterministic when dirs share a suffix.
func TestBuildCompositeResolution(t *testing.T) {
	caller := &model.Workflow{
		File: "ci.yml", Name: "CI", On: []string{"push"},
		Jobs: []model.Job{{ID: "build", RunsOn: "ubuntu-latest", Steps: []model.Step{
			{Uses: "./.github/actions/build"},
		}}},
	}
	g := Build([]Source{
		{Path: "/repo/.github/workflows/ci.yml", Workflow: caller},
		{Path: "/repo/.github/actions/my-build/action.yml", Action: &model.Action{Name: "My Build"}},
		{Path: "/repo/.github/actions/build/action.yml", Action: &model.Action{Name: "Build"}},
	})

	want := "/repo/.github/actions/build/action.yml"
	var got string
	var edge Edge
	for _, e := range g.Calls("/repo/.github/workflows/ci.yml") {
		if e.Kind == KindComposite {
			got = e.ToID
			edge = e
		}
	}
	if got != want {
		t.Errorf("composite ./.github/actions/build resolved to %q, want %q (must not match my-build)", got, want)
	}
	// The calling step is unnamed; it must still get a positional StepName so the call
	// graph renders it as "job / step N (uses ...)" rather than collapsing to job level.
	if edge.StepName != "step 1" {
		t.Errorf("unnamed composite step: StepName = %q, want %q", edge.StepName, "step 1")
	}
}

// TestResolveUnresolvedLocalKind checks that an unresolved local `uses:` is classified by
// its shape: a composite-action-looking path stays KindComposite (not KindReusable), and
// a workflow-looking path stays KindReusable.
func TestResolveUnresolvedLocalKind(t *testing.T) {
	caller := &model.Workflow{
		File: "ci.yml", Name: "CI", On: []string{"push"},
		Jobs: []model.Job{
			// references a local composite action that is NOT in the scan set
			{ID: "a", RunsOn: "ubuntu-latest", Steps: []model.Step{{Uses: "./.github/actions/missing"}}},
			// references a local reusable workflow that is NOT in the scan set
			{ID: "b", Uses: "./.github/workflows/missing.yml"},
		},
	}
	g := Build([]Source{{Path: "/r/.github/workflows/ci.yml", Workflow: caller}})

	var compKind, reusableKind EdgeKind
	for _, e := range g.Calls("/r/.github/workflows/ci.yml") {
		switch e.Ref {
		case "./.github/actions/missing":
			compKind = e.Kind
		case "./.github/workflows/missing.yml":
			reusableKind = e.Kind
		}
	}
	if compKind != KindComposite {
		t.Errorf("unresolved ./.github/actions/missing kind = %q, want %q", compKind, KindComposite)
	}
	if reusableKind != KindReusable {
		t.Errorf("unresolved ./.github/workflows/missing.yml kind = %q, want %q", reusableKind, KindReusable)
	}
}

// TestResolveWorkflowBasenameCollision verifies that two scanned workflows sharing a
// basename in different directories resolve to the correct one based on the caller's
// `uses:` path, rather than colliding on the bare filename.
func TestResolveWorkflowBasenameCollision(t *testing.T) {
	caller := &model.Workflow{
		File: "ci.yml", Name: "CI", On: []string{"push"},
		Jobs: []model.Job{
			{ID: "a", Uses: "./.github/workflows/sub/build.yml"},
			{ID: "b", Uses: "./.github/workflows/build.yml"},
		},
	}
	subBuild := &model.Workflow{File: "build.yml", Name: "Sub Build", On: []string{"workflow_call"}}
	topBuild := &model.Workflow{File: "build.yml", Name: "Top Build", On: []string{"workflow_call"}}

	g := Build([]Source{
		{Path: "repo/.github/workflows/ci.yml", Workflow: caller},
		{Path: "repo/.github/workflows/sub/build.yml", Workflow: subBuild},
		{Path: "repo/.github/workflows/build.yml", Workflow: topBuild},
	})

	got := map[string]string{} // jobID -> resolved ToID
	for _, e := range g.Calls("repo/.github/workflows/ci.yml") {
		got[e.JobID] = e.ToID
	}
	if got["a"] != "repo/.github/workflows/sub/build.yml" {
		t.Errorf("job a (uses sub/build.yml) resolved to %q, want the sub/ workflow", got["a"])
	}
	if got["b"] != "repo/.github/workflows/build.yml" {
		t.Errorf("job b (uses build.yml) resolved to %q, want the top-level workflow", got["b"])
	}
}

// TestResolveWorkflowOSNativePaths verifies reusable-workflow resolution when scan paths
// use OS-native separators (via filepath.FromSlash) while uses: refs stay slash-separated
// as in YAML. The resolved node ID must match the original node key, and basename-
// collision disambiguation must still work. On Windows FromSlash yields backslash paths,
// exercising the slash-normalization that keeps keys consistent; on Unix it is a no-op
// (still a valid regression guard against keying nodes and the index differently).
func TestResolveWorkflowOSNativePaths(t *testing.T) {
	caller := &model.Workflow{
		File: "ci.yml", Name: "CI", On: []string{"push"},
		Jobs: []model.Job{
			{ID: "a", Uses: "./.github/workflows/sub/build.yml"},
			{ID: "b", Uses: "./.github/workflows/build.yml"},
		},
	}
	subBuild := &model.Workflow{File: "build.yml", Name: "Sub", On: []string{"workflow_call"}}
	topBuild := &model.Workflow{File: "build.yml", Name: "Top", On: []string{"workflow_call"}}

	ciPath := filepath.FromSlash("repo/.github/workflows/ci.yml")
	subPath := filepath.FromSlash("repo/.github/workflows/sub/build.yml")
	topPath := filepath.FromSlash("repo/.github/workflows/build.yml")
	g := Build([]Source{
		{Path: ciPath, Workflow: caller},
		{Path: subPath, Workflow: subBuild},
		{Path: topPath, Workflow: topBuild},
	})

	got := map[string]string{}
	for _, e := range g.Calls(ciPath) {
		got[e.JobID] = e.ToID
	}
	for job, id := range got {
		if id == "" || g.Nodes[id] == nil {
			t.Errorf("job %s resolved to %q, which is not a node key", job, id)
		}
	}
	if got["a"] != subPath {
		t.Errorf("job a resolved to %q, want %q", got["a"], subPath)
	}
	if got["b"] != topPath {
		t.Errorf("job b resolved to %q, want %q", got["b"], topPath)
	}
}

// TestResolveCompositeRefMoreQualified verifies composite resolution when the uses: ref
// is MORE qualified than the scanned action dir (e.g. scan dir "actions/build", ref
// "./.github/actions/build") -- the asymmetric case that the one-directional suffix
// match missed.
func TestResolveCompositeRefMoreQualified(t *testing.T) {
	caller := &model.Workflow{
		File: "ci.yml", Name: "CI", On: []string{"push"},
		Jobs: []model.Job{{ID: "a", RunsOn: "ubuntu-latest", Steps: []model.Step{
			{Uses: "./.github/actions/build"},
		}}},
	}
	// Scan dir is shallower than the ref's path.
	action := &model.Action{Name: "Build"}
	g := Build([]Source{
		{Path: "ci.yml", Workflow: caller},
		{Path: "actions/build/action.yml", Action: action},
	})

	var got string
	var kind EdgeKind
	for _, e := range g.Calls("ci.yml") {
		if e.Ref == "./.github/actions/build" {
			got, kind = e.ToID, e.Kind
		}
	}
	if got != "actions/build/action.yml" || kind != KindComposite {
		t.Errorf("composite ref more-qualified than scan dir resolved to %q (kind %q), want the action node", got, kind)
	}
}

// TestCrossRepoSelfRefsResolve covers the inference that a cross-repo ref whose
// owner/repo prefix consistently points back into the scan set is the scanned repository
// calling itself (with a branch/tag pin), and must resolve to the in-scope node so the
// callee gets its Called-by chain. Prefixes that reference anything outside the scan set
// stay external, so a same-named workflow in a different repo is never mis-linked.
func TestCrossRepoSelfRefsResolve(t *testing.T) {
	release := &model.Workflow{
		File: "release.yml", Name: "Release", On: []string{"push"},
		Jobs: []model.Job{
			// Self-reference in cross-repo form, pinned to a branch.
			{ID: "precheck", Uses: "jreleaser/jreleaser/.github/workflows/step-precheck.yml@main"},
			// Genuinely external: basename matches an in-scope workflow, but the prefix
			// also references a workflow that is NOT in scope, so the whole prefix is
			// treated as a different repository.
			{ID: "shared-pre", Uses: "other-org/shared/.github/workflows/step-precheck.yml@v1"},
			{ID: "shared-scan", Uses: "other-org/shared/.github/workflows/security-scan.yml@v1"},
			// Genuinely external with no in-scope match at all.
			{ID: "slsa", Uses: "slsa-framework/slsa-github-generator/.github/workflows/generator_generic_slsa3.yml@v2.1.0"},
		},
	}
	precheck := &model.Workflow{
		File: "step-precheck.yml", Name: "Precheck", On: []string{"workflow_call"},
		Jobs: []model.Job{{ID: "check", RunsOn: "ubuntu-latest"}},
	}

	g := Build([]Source{
		{Path: ".github/workflows/release.yml", Workflow: release},
		{Path: ".github/workflows/step-precheck.yml", Workflow: precheck},
	})

	edges := g.Calls(".github/workflows/release.yml")
	byJob := map[string]Edge{}
	for _, e := range edges {
		byJob[e.JobID] = e
	}

	// The self-reference resolves in scope, keeping the pin.
	self := byJob["precheck"]
	if self.ToID != ".github/workflows/step-precheck.yml" {
		t.Errorf("self cross-repo ref ToID = %q, want in-scope node", self.ToID)
	}
	if self.Pin != "main" {
		t.Errorf("self cross-repo ref Pin = %q, want main", self.Pin)
	}
	// The callee's Called-by chain includes the caller.
	calledBy := g.CalledBy(".github/workflows/step-precheck.yml")
	if len(calledBy) != 1 || calledBy[0].FromID != ".github/workflows/release.yml" {
		t.Errorf("CalledBy = %+v, want one edge from release.yml", calledBy)
	}

	// The mixed prefix (one in-scope basename, one not) stays fully external.
	if to := byJob["shared-pre"].ToID; to != "other-org/shared/.github/workflows/step-precheck.yml" {
		t.Errorf("mixed-prefix ref ToID = %q, want external node", to)
	}
	if n := g.Nodes[byJob["shared-pre"].ToID]; n == nil || !n.External {
		t.Error("mixed-prefix ref must remain an external node")
	}

	// The no-match prefix stays external.
	if n := g.Nodes[byJob["slsa"].ToID]; n == nil || !n.External {
		t.Error("unrelated cross-repo ref must remain an external node")
	}
}

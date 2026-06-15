package renderer

import (
	"strings"
	"testing"

	"github.com/smol-utils/actiondoc/internal/callgraph"
	"github.com/smol-utils/actiondoc/internal/model"
)

// inventoryGraph builds two in-memory workflows that share a secret, reference a variable,
// touch GITHUB_TOKEN (which must be filtered), and declare conflicting permission levels for
// the same scope. It returns the built graph and the source slice in the order the document
// assembler would pass them.
func inventoryGraph() ([]callgraph.Source, *callgraph.Graph) {
	alpha := &model.Workflow{
		File: "alpha.yml", Name: "Alpha", On: []string{"push"},
		Permissions: &model.Permissions{Scopes: []model.Permission{
			{Scope: "contents", Level: "read"},
		}},
		Jobs: []model.Job{{
			ID: "build", RunsOn: "ubuntu-latest",
			Steps: []model.Step{{
				Name: "Run",
				Run:  `echo "${{ secrets.SHARED_TOKEN }} ${{ secrets.GITHUB_TOKEN }} ${{ vars.REGION }}"`,
			}},
		}},
	}
	beta := &model.Workflow{
		File: "beta.yml", Name: "Beta", On: []string{"workflow_dispatch"},
		Permissions: &model.Permissions{Scopes: []model.Permission{
			{Scope: "contents", Level: "write"},
			{Scope: "id-token", Level: "write", OIDC: true},
		}},
		Jobs: []model.Job{{
			ID: "deploy", RunsOn: "ubuntu-latest",
			Steps: []model.Step{{
				Name: "Deploy",
				Run:  `deploy --token "${{ secrets.SHARED_TOKEN }}"`,
			}},
		}},
	}
	sources := []callgraph.Source{
		{Path: ".github/workflows/alpha.yml", Workflow: alpha},
		{Path: ".github/workflows/beta.yml", Workflow: beta},
	}
	return sources, callgraph.Build(sources)
}

func TestRenderDocumentInventory(t *testing.T) {
	sources, g := inventoryGraph()
	out := RenderDocumentInventory(sources, g)

	// A shared secret is deduped to one row and links to BOTH workflows, sorted by anchor.
	wantSecretRow := "| `SHARED_TOKEN` | [Alpha](#alpha), [Beta](#beta) |"
	if !strings.Contains(out, wantSecretRow) {
		t.Errorf("missing deduped+linked shared secret row %q\n\n%s", wantSecretRow, out)
	}

	// GITHUB_TOKEN is auto-provided and must never appear in the secrets table.
	if strings.Contains(out, "GITHUB_TOKEN") {
		t.Errorf("GITHUB_TOKEN should be excluded from the inventory\n\n%s", out)
	}

	// Variables get their own table.
	if !strings.Contains(out, "**Variables:**") || !strings.Contains(out, "| `REGION` | [Alpha](#alpha) |") {
		t.Errorf("missing variables table / REGION row\n\n%s", out)
	}

	// A scope granted at different levels keeps its effective (max) level plus a conflict note.
	if !strings.Contains(out, "| `contents` | `write` (also granted as `read` elsewhere) |") {
		t.Errorf("missing contents read-vs-write conflict row\n\n%s", out)
	}

	// id-token: write carries the (OIDC) marker.
	if !strings.Contains(out, "| `id-token` | `write` (OIDC) |") {
		t.Errorf("missing id-token OIDC row\n\n%s", out)
	}

	// Secret names are sorted (no extra secrets here, but the heading must be present).
	if !strings.Contains(out, inventorySecretsVarsHeading) || !strings.Contains(out, inventoryPermissionsHeading) {
		t.Errorf("missing inventory headings\n\n%s", out)
	}
}

// TestRenderDocumentInventoryEmpty verifies the whole section is suppressed when there are no
// secrets, no variables, and no permissions to report.
func TestRenderDocumentInventoryEmpty(t *testing.T) {
	w := &model.Workflow{
		File: "plain.yml", Name: "Plain", On: []string{"push"},
		Jobs: []model.Job{{ID: "noop", RunsOn: "ubuntu-latest",
			Steps: []model.Step{{Name: "Echo", Run: "echo hi"}}}},
	}
	sources := []callgraph.Source{{Path: ".github/workflows/plain.yml", Workflow: w}}
	g := callgraph.Build(sources)

	if out := RenderDocumentInventory(sources, g); out != "" {
		t.Errorf("expected empty inventory, got:\n%s", out)
	}
}

// TestRenderDocumentInventorySubTableSuppression verifies each sub-table is suppressed
// independently: permissions-only input renders the permissions table but no secrets/vars one.
func TestRenderDocumentInventorySubTableSuppression(t *testing.T) {
	w := &model.Workflow{
		File: "perm.yml", Name: "Perm", On: []string{"push"},
		Permissions: &model.Permissions{Scopes: []model.Permission{
			{Scope: "contents", Level: "read"},
		}},
		Jobs: []model.Job{{ID: "noop", RunsOn: "ubuntu-latest",
			Steps: []model.Step{{Name: "Echo", Run: "echo hi"}}}},
	}
	sources := []callgraph.Source{{Path: ".github/workflows/perm.yml", Workflow: w}}
	out := RenderDocumentInventory(sources, callgraph.Build(sources))

	if strings.Contains(out, inventorySecretsVarsHeading) {
		t.Errorf("secrets/vars section should be suppressed when empty\n\n%s", out)
	}
	if !strings.Contains(out, inventoryPermissionsHeading) {
		t.Errorf("permissions section should render\n\n%s", out)
	}
}

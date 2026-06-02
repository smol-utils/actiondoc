package renderer

import (
	"fmt"
	"strings"

	"github.com/smol-utils/actiondoc/internal/model"
)

// renderPermissions writes a permissions: block under the given heading (a "## ...\n\n"
// for workflow level or a "**...:**\n\n" for job level). The three block shapes render
// distinctly: a scalar grant and an explicit empty block each get a one-line note; a
// scope list renders as bullets with OIDC and rationale annotations.
func renderPermissions(b *strings.Builder, p *model.Permissions, heading string) {
	if p == nil {
		return
	}
	b.WriteString(heading)

	switch {
	case p.DefaultDeny:
		b.WriteString("No permissions granted (`permissions: {}` -- default-deny).\n\n")
	case p.All != "":
		fmt.Fprintf(b, "All scopes: `%s`.\n\n", p.All)
	default:
		for _, s := range p.Scopes {
			fmt.Fprintf(b, "- `%s`: `%s`", s.Scope, s.Level)
			if s.OIDC {
				b.WriteString(" (OIDC)")
			}
			if s.Rationale != "" {
				fmt.Fprintf(b, " - %s", s.Rationale)
			}
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}
}

// renderEnvironment writes a job's GitHub Environments binding -- the deploy gate. The
// binding is marked [gated] and followed by a pointer to the protection rules, which live
// in repository settings rather than the workflow YAML.
func renderEnvironment(b *strings.Builder, e *model.Environment) {
	if e == nil {
		return
	}
	fmt.Fprintf(b, "**Deploys to environment:** `%s`", e.Name)
	if e.URL != "" {
		fmt.Fprintf(b, " (`%s`)", e.URL)
	}
	b.WriteString(" [gated]\n\n")
	b.WriteString("> Environment protection rules (required reviewers, wait timers, branch policies) " +
		"are configured in the repository's Settings -> Environments and are not represented here.\n\n")
}

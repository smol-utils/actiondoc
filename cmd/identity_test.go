package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestParseGitHubRemote(t *testing.T) {
	cases := []struct {
		url  string
		want string
	}{
		{"git@github.com:owner/repo.git", "owner/repo"},
		{"git@github.com:owner/repo", "owner/repo"},
		{"https://github.com/owner/repo.git", "owner/repo"},
		{"https://github.com/owner/repo", "owner/repo"},
		{"ssh://git@github.com/owner/repo.git", "owner/repo"},
		{"  https://github.com/owner/repo.git\n", "owner/repo"},   // trailing newline from git
		{"ssh://git@github.com:443/owner/repo.git", "owner/repo"}, // explicit port is stripped
		{"ssh://git@github.com:22/owner/repo.git", "owner/repo"},
		{"https://github.com:443/owner/repo.git", "owner/repo"},
		{"https://gitlab.com/owner/repo.git", ""},            // not GitHub
		{"git@example.com:owner/repo.git", ""},               // not GitHub
		{"https://notgithub.com/owner/repo.git", ""},         // host merely ends in github.com
		{"git@notgithub.com:owner/repo.git", ""},             // SSH host merely ends in github.com
		{"https://my-github.com/owner/repo.git", ""},         // lookalike host
		{"https://gitlab.com/github.com/owner/repo.git", ""}, // github.com is a path segment, not the host
		{"git@github.com.evil.com:owner/repo.git", ""},       // host extends past github.com
		{"https://github.com.au/owner/repo.git", ""},         // host extends past github.com
		{"github.company.com:owner/repo.git", ""},            // unrelated host containing "github"
		{"", ""},
		{"https://github.com/owner", ""}, // missing repo segment
	}
	for _, c := range cases {
		if got := parseGitHubRemote(c.url); got != c.want {
			t.Errorf("parseGitHubRemote(%q) = %q, want %q", c.url, got, c.want)
		}
	}
}

func TestNormalizeIdentity(t *testing.T) {
	cases := []struct{ in, want string }{
		{"owner/repo", "owner/repo"},
		{" owner/repo.git ", "owner/repo"},
		{"owner", ""},
		{"owner/repo/extra", ""},
		{"/repo", ""},
		{"owner/", ""},
		{"", ""},
	}
	for _, c := range cases {
		if got := normalizeIdentity(c.in); got != c.want {
			t.Errorf("normalizeIdentity(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

// TestResolveRepoIdentityPrecedence checks the flag > GITHUB_REPOSITORY > git-remote >
// unknown order, and that a non-GitHub / non-repo scan path resolves to "" (external by
// default).
func TestResolveRepoIdentityPrecedence(t *testing.T) {
	// A temp dir that is a git repo with a GitHub origin, used to exercise the git-remote
	// branch of resolution.
	repoDir := t.TempDir()
	for _, args := range [][]string{
		{"init"},
		{"remote", "add", "origin", "git@github.com:remote-org/remote-repo.git"},
	} {
		cmd := exec.Command("git", append([]string{"-C", repoDir}, args...)...)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Skipf("git unavailable or failed (%v): %s", err, out)
		}
	}

	// flag wins over everything.
	t.Setenv("GITHUB_REPOSITORY", "env-org/env-repo")
	if got := resolveRepoIdentity("flag-org/flag-repo", repoDir); got != "flag-org/flag-repo" {
		t.Errorf("flag precedence: got %q, want flag-org/flag-repo", got)
	}

	// env wins over git remote when no flag.
	if got := resolveRepoIdentity("", repoDir); got != "env-org/env-repo" {
		t.Errorf("env precedence: got %q, want env-org/env-repo", got)
	}

	// git remote used when neither flag nor env is set.
	os.Unsetenv("GITHUB_REPOSITORY")
	if got := resolveRepoIdentity("", repoDir); got != "remote-org/remote-repo" {
		t.Errorf("git-remote precedence: got %q, want remote-org/remote-repo", got)
	}

	// A path that is a file inside the repo resolves the same (git walks up).
	f := filepath.Join(repoDir, "workflow.yml")
	if err := os.WriteFile(f, []byte("name: x\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := resolveRepoIdentity("", f); got != "remote-org/remote-repo" {
		t.Errorf("git-remote from file path: got %q, want remote-org/remote-repo", got)
	}

	// Unknown: a non-repo directory with nothing set -> "" (external by default).
	if got := resolveRepoIdentity("", t.TempDir()); got != "" {
		t.Errorf("unknown identity: got %q, want \"\"", got)
	}
}

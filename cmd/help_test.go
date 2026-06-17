package cmd

import (
	"errors"
	"flag"
	"testing"
)

// TestGenerateHelpReturnsErrHelp confirms an explicit help request surfaces as flag.ErrHelp
// (and not some other error), so the top-level command can treat it as success rather than a
// failure and exit 0 with no "error:" line.
func TestGenerateHelpReturnsErrHelp(t *testing.T) {
	for _, arg := range []string{"-h", "--help"} {
		if err := Generate([]string{arg}); !errors.Is(err, flag.ErrHelp) {
			t.Errorf("Generate(%q) = %v, want flag.ErrHelp", arg, err)
		}
	}
}

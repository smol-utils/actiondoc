package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/smol-utils/actiondoc/internal/parser"
	"github.com/smol-utils/actiondoc/internal/renderer"
)

// Generate runs the "generate" subcommand.
func Generate(args []string) error {
	fs := flag.NewFlagSet("generate", flag.ContinueOnError)
	outFlag := fs.String("o", "", "output file (default: stdout)")
	jsonFlag := fs.Bool("json", false, "output JSON instead of Markdown")
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: actiondoc generate [flags] [path]\n\n")
		fmt.Fprintf(os.Stderr, "Generates documentation for GitHub Actions workflow files.\n\n")
		fmt.Fprintf(os.Stderr, "Arguments:\n")
		fmt.Fprintf(os.Stderr, "  path    Path to a YAML file or directory (default: .github/workflows)\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		return err
	}

	path := ".github/workflows"
	if fs.NArg() > 0 {
		path = fs.Arg(0)
	}

	files, err := resolveFiles(path)
	if err != nil {
		return fmt.Errorf("resolving %s: %w", path, err)
	}
	if len(files) == 0 {
		return fmt.Errorf("no YAML files found in %s", path)
	}

	// Process each file, collecting results for JSON or rendering Markdown inline.
	var jsonItems []any
	var md strings.Builder

	for _, f := range files {
		w, err := parser.ParseFile(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warning: %v\n", err)
			continue
		}
		if *jsonFlag {
			jsonItems = append(jsonItems, w)
		} else {
			md.WriteString(renderer.RenderMarkdown(w))
		}
	}

	var output string
	if *jsonFlag {
		data, err := json.MarshalIndent(jsonItems, "", "  ")
		if err != nil {
			return fmt.Errorf("marshaling JSON: %w", err)
		}
		output = string(data) + "\n"
	} else {
		output = md.String()
	}

	if *outFlag != "" {
		if err := os.WriteFile(*outFlag, []byte(output), 0644); err != nil {
			return fmt.Errorf("writing %s: %w", *outFlag, err)
		}
	} else {
		fmt.Print(output)
	}
	return nil
}

// resolveFiles returns a list of .yml/.yaml files from the given path.
func resolveFiles(path string) ([]string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		return []string{path}, nil
	}

	var files []string
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(e.Name()))
		if ext == ".yml" || ext == ".yaml" {
			files = append(files, filepath.Join(path, e.Name()))
		}
	}
	return files, nil
}

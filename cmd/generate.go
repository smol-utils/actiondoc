package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/smol-utils/actiondoc/internal/callgraph"
	"github.com/smol-utils/actiondoc/internal/model"
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
		fmt.Fprintf(os.Stderr, "Generates documentation for GitHub Actions workflow and action files.\n\n")
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

	// First pass: parse every file. The call graph needs the whole scan set before any
	// workflow can be rendered (reusable-workflow cross-links and the call-graph tree
	// resolve `uses:` targets across files), so parse all, then build the graph, then
	// render. A parsed file keeps its source path as its call-graph node id.
	type parsed struct {
		path     string
		workflow *model.Workflow
		action   *model.Action
	}
	var items []parsed
	var sources []callgraph.Source
	for _, f := range files {
		if isActionFile(f) {
			a, err := parser.ParseActionFile(f)
			if err != nil {
				fmt.Fprintf(os.Stderr, "warning: %v\n", err)
				continue
			}
			items = append(items, parsed{path: f, action: a})
			sources = append(sources, callgraph.Source{Path: f, Action: a})
		} else {
			w, err := parser.ParseFile(f)
			if err != nil {
				fmt.Fprintf(os.Stderr, "warning: %v\n", err)
				continue
			}
			items = append(items, parsed{path: f, workflow: w})
			sources = append(sources, callgraph.Source{Path: f, Workflow: w})
		}
	}
	graph := callgraph.Build(sources)

	// Second pass: emit JSON or rendered Markdown, now with graph context available.
	var jsonItems []any
	var md strings.Builder
	for _, it := range items {
		switch {
		case it.action != nil:
			if *jsonFlag {
				jsonItems = append(jsonItems, it.action)
			} else {
				md.WriteString(renderer.RenderActionMarkdown(it.action))
			}
		case it.workflow != nil:
			if *jsonFlag {
				jsonItems = append(jsonItems, it.workflow)
			} else {
				md.WriteString(renderer.RenderMarkdownGraph(it.workflow, graph, it.path))
			}
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

// isActionFile returns true if the file is a GitHub Action metadata file.
func isActionFile(path string) bool {
	base := strings.ToLower(filepath.Base(path))
	return base == "action.yml" || base == "action.yaml"
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

	// When pointed at a .github/workflows directory, also auto-discover sibling
	// composite actions under .github/actions/ so they render in the same output.
	// (Composite actions placed directly inside the workflows dir are already picked up
	// above, since action.yml/.yaml has a .yml/.yaml extension.)
	if strings.EqualFold(filepath.Base(path), "workflows") {
		discovered, err := discoverActionFiles(filepath.Join(filepath.Dir(path), "actions"))
		if err != nil {
			return nil, err
		}
		files = append(files, discovered...)
	}
	return files, nil
}

// discoverActionFiles walks dir (if present) for composite action metadata files named
// action.yml/action.yaml, at any depth. A missing dir is not an error.
func discoverActionFiles(dir string) ([]string, error) {
	if info, err := os.Stat(dir); err != nil || !info.IsDir() {
		return nil, nil
	}
	var out []string
	err := filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if isActionFile(p) {
			out = append(out, p)
		}
		return nil
	})
	return out, err
}

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

	// Parse the whole scan set first: the call graph needs every file before any workflow
	// can be rendered (reusable-workflow cross-links and the call-graph tree resolve
	// `uses:` targets across files). Build the graph once, then link local composite
	// actions into the steps that reference them so the renderer can pair `with:` keys
	// with declared inputs. A source's path is its call-graph node id.
	sources, parseFailures := parseSources(files)
	graph := callgraph.Build(sources)
	linkCompositeActions(sources, graph)

	var output string
	if *jsonFlag {
		var jsonItems []any
		for _, s := range sources {
			if s.Workflow != nil {
				jsonItems = append(jsonItems, s.Workflow)
			} else {
				jsonItems = append(jsonItems, s.Action)
			}
		}
		data, err := json.MarshalIndent(jsonItems, "", "  ")
		if err != nil {
			return fmt.Errorf("marshaling JSON: %w", err)
		}
		output = string(data) + "\n"
	} else {
		// Render each document as a section, with a table of contents linking them.
		// Workflows render with graph context so cross-links and call-graph sections
		// appear; actions render standalone.
		var sections, titles []string
		for _, s := range sources {
			if s.Workflow != nil {
				sections = append(sections, renderer.RenderMarkdownGraph(s.Workflow, graph, s.Path))
				titles = append(titles, s.Workflow.Name)
			} else {
				sections = append(sections, renderer.RenderActionMarkdown(s.Action))
				titles = append(titles, s.Action.Name)
			}
		}
		output = renderer.RenderTOC(titles) + strings.Join(sections, "")
	}

	if *outFlag != "" {
		if err := os.WriteFile(*outFlag, []byte(output), 0644); err != nil {
			return fmt.Errorf("writing %s: %w", *outFlag, err)
		}
	} else {
		fmt.Print(output)
	}

	// Output for the files that did parse is still written above, but a parse failure is a
	// real error: return non-zero so callers (e.g. the dogfood smoke test) don't treat a
	// partially-parsed run as success.
	if parseFailures > 0 {
		return fmt.Errorf("%d file(s) failed to parse", parseFailures)
	}
	return nil
}

// parseSources parses each file into a callgraph source, skipping (with a warning) files
// that fail to parse. Source order follows file order so rendering stays deterministic.
// The returned count is the number of files that failed to parse, so the caller can exit
// non-zero rather than silently reporting success on a partially-parsed scan.
func parseSources(files []string) ([]callgraph.Source, int) {
	var sources []callgraph.Source
	failed := 0
	for _, f := range files {
		if isActionFile(f) {
			a, err := parser.ParseActionFile(f)
			if err != nil {
				fmt.Fprintf(os.Stderr, "warning: %v\n", err)
				failed++
				continue
			}
			sources = append(sources, callgraph.Source{Path: f, Action: a})
		} else {
			w, err := parser.ParseFile(f)
			if err != nil {
				fmt.Fprintf(os.Stderr, "warning: %v\n", err)
				failed++
				continue
			}
			sources = append(sources, callgraph.Source{Path: f, Workflow: w})
		}
	}
	return sources, failed
}

// linkCompositeActions attaches each parsed local composite action to the step whose
// `uses:` references it, using the prebuilt call graph, so the renderer can pair `with:`
// keys with the action's declared input descriptions.
func linkCompositeActions(sources []callgraph.Source, g *callgraph.Graph) {
	for _, e := range g.Edges {
		if e.Kind != callgraph.KindComposite {
			continue
		}
		target := g.Nodes[e.ToID]
		caller := g.Nodes[e.FromID]
		if target == nil || target.Action == nil || caller == nil || caller.Workflow == nil {
			continue
		}
		for ji := range caller.Workflow.Jobs {
			job := &caller.Workflow.Jobs[ji]
			if job.ID != e.JobID {
				continue
			}
			// Match on the raw uses: ref -- the exact string the edge was built from.
			for si := range job.Steps {
				if job.Steps[si].Uses == e.Ref {
					job.Steps[si].UsesAction = target.Action
				}
			}
		}
	}
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
// action.yml/action.yaml, at any depth. A missing dir is not an error; other stat errors
// (permissions, I/O) are surfaced rather than silently swallowed.
func discoverActionFiles(dir string) ([]string, error) {
	info, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	if !info.IsDir() {
		return nil, nil
	}
	var out []string
	err = filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
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

package cmd

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
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
		output, err = renderJSONOutput(sources)
		if err != nil {
			return err
		}
	} else {
		output = renderMarkdownOutput(sources, graph, path)
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

// renderJSONOutput marshals the parsed models as a JSON array -- the machine-readable
// form of everything the Markdown renderer documents. The model's JSON field tags are a
// stable contract for downstream consumers.
func renderJSONOutput(sources []callgraph.Source) (string, error) {
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
		return "", fmt.Errorf("marshaling JSON: %w", err)
	}
	return string(data) + "\n", nil
}

// renderMarkdownOutput renders each document as a section, prefaced by a document header
// (repo title + inventory) and a grouped table of contents. Anchors are assigned before
// any section renders: cross-links built during rendering must use the same duplicate-name
// disambiguation the TOC will use, so the assignment is computed once here and stored on
// the graph nodes. inputPath is the path the scan was launched from; it supplies the repo
// title. The header and TOC are emitted only for multi-document output.
func renderMarkdownOutput(sources []callgraph.Source, graph *callgraph.Graph, inputPath string) string {
	// Order the sources to match the table of contents: entry-point workflows, then reusable
	// workflows, then composite actions, sorted within each group. Every downstream pass --
	// the section-anchor pass, the document-wide job-anchor pass, the rendered body sections,
	// and the TOC -- iterates this one ordered slice, so the body and the TOC are guaranteed to
	// agree and the order-dependent anchor "-N" numbering is computed over the final order.
	sources = orderSourcesForRender(sources, graph)

	var titles []string
	for _, s := range sources {
		if s.Workflow != nil {
			titles = append(titles, s.Workflow.Name)
		} else {
			titles = append(titles, s.Action.Name)
		}
	}
	slugs := renderer.AssignAnchors(titles)
	for i, slug := range slugs {
		if n := graph.Nodes[sources[i].Path]; n != nil {
			n.Anchor = slug
		}
	}

	// Job heading anchors are assigned the same way section anchors are: document-wide.
	// GitHub disambiguates repeated heading slugs across the whole rendered document, so a
	// job heading text that recurs in a later workflow must carry the running "-N" suffix.
	// Collect every job heading in document order (source order, then job order within each
	// workflow) -- using the renderer's own JobHeadingText so the slug input matches the
	// rendered "### ..." heading exactly -- run one AssignAnchors pass, then store each
	// workflow's slice on its graph node for renderJobMiniTOC to use.
	var jobTexts []string
	type jobSpan struct {
		path  string
		start int
		count int
	}
	var spans []jobSpan
	for _, s := range sources {
		if s.Workflow == nil {
			continue
		}
		spans = append(spans, jobSpan{path: s.Path, start: len(jobTexts), count: len(s.Workflow.Jobs)})
		for i := range s.Workflow.Jobs {
			jobTexts = append(jobTexts, renderer.JobHeadingText(&s.Workflow.Jobs[i]))
		}
	}
	jobSlugs := renderer.AssignAnchors(jobTexts)
	for _, sp := range spans {
		if n := graph.Nodes[sp.path]; n != nil {
			n.JobAnchors = jobSlugs[sp.start : sp.start+sp.count]
		}
	}

	// Workflows render with graph context so cross-links and call-graph sections
	// appear; actions render standalone.
	var sections []string
	for _, s := range sources {
		if s.Workflow != nil {
			sections = append(sections, renderer.RenderMarkdownGraph(s.Workflow, graph, s.Path))
		} else {
			sections = append(sections, renderer.RenderActionMarkdown(s.Action))
		}
	}
	// A single document is self-describing (its own H1 + properties); the orientation
	// header, contents list, and per-section back-to-contents links only earn their space
	// once there are several sections to navigate between.
	if len(sources) < 2 {
		return strings.Join(sections, "")
	}

	// Each top-level section ends with a link back to the Contents list, so a reader deep in
	// one section can return to navigation without scrolling.
	for i := range sections {
		sections[i] += "[Back to contents](#contents)\n\n"
	}

	// The document-level inventory (repo-wide secrets/variables + permissions) sits between
	// the table of contents and the per-section bodies, after anchors are assigned so its
	// "used by" cross-links resolve to the same section anchors the TOC uses. It is
	// multi-document only -- the same condition that gates the header and TOC above -- and
	// returns "" when there is nothing to inventory.
	header, toc := renderDocumentNav(sources, graph, slugs, inputPath)
	triggerIndex := renderer.RenderTriggerIndex(sources, graph)
	inventory := renderer.RenderDocumentInventory(sources, graph)
	return header + toc + triggerIndex + inventory + strings.Join(sections, "")
}

// tocGroup indexes the three TOC families in render order.
const (
	groupWorkflow  = iota // entry-point workflows
	groupReusable         // workflow_call-only workflows
	groupComposite        // composite actions
)

// groupOf classifies a source into its TOC family: composite action, entry-point workflow
// (a trigger other than workflow_call), or reusable (workflow_call-only) workflow. The body
// section order and the TOC grouping both derive from this single classifier, so they cannot
// disagree about which family a source belongs to.
func groupOf(s callgraph.Source, graph *callgraph.Graph) int {
	switch {
	case s.Action != nil:
		return groupComposite
	case graph.IsEntryPoint(s.Path):
		return groupWorkflow
	default:
		return groupReusable
	}
}

// orderSourcesForRender returns the sources reordered to match the table of contents:
// entry-point workflows first, then reusable workflows, then composite actions. Within each
// group entries are sorted by display title (case-insensitively, for a legible reading order),
// with the exact title and then the file path breaking ties so the result is fully
// deterministic. The body and the TOC both iterate the returned slice, so scrolling follows
// the same mental model the TOC sets up. The input slice is not mutated.
func orderSourcesForRender(sources []callgraph.Source, graph *callgraph.Graph) []callgraph.Source {
	ordered := make([]callgraph.Source, len(sources))
	copy(ordered, sources)
	sort.SliceStable(ordered, func(i, j int) bool {
		gi, gj := groupOf(ordered[i], graph), groupOf(ordered[j], graph)
		if gi != gj {
			return gi < gj
		}
		ti, tj := titleOf(ordered[i]), titleOf(ordered[j])
		if li, lj := strings.ToLower(ti), strings.ToLower(tj); li != lj {
			return li < lj
		}
		if ti != tj {
			return ti < tj
		}
		return ordered[i].Path < ordered[j].Path
	})
	return ordered
}

// renderDocumentNav builds the document header (title + inventory) and the grouped table of
// contents. Each source is classified via the call graph: composite actions, entry-point
// workflows (a trigger other than workflow_call), and reusable (workflow_call-only)
// workflows. Entry-point labels carry their trigger list; labels that would otherwise read
// identically are disambiguated with their source filename (for actions, the action's
// directory, since every action file is named action.yml).
func renderDocumentNav(sources []callgraph.Source, graph *callgraph.Graph, slugs []string, inputPath string) (header, toc string) {
	type item struct {
		label  string
		anchor string
		group  int
	}
	items := make([]item, len(sources))
	var nWorkflow, nReusable, nComposite int
	for i, s := range sources {
		label := titleOf(s)
		group := groupOf(s, graph)
		switch group {
		case groupComposite:
			nComposite++
		case groupWorkflow:
			nWorkflow++
			if len(s.Workflow.On) > 0 {
				label += " - " + strings.Join(s.Workflow.On, ", ")
			}
		default:
			nReusable++
		}
		items[i] = item{label: label, anchor: slugs[i], group: group}
	}

	// Disambiguate entries whose visible label collides with another's.
	counts := map[string]int{}
	for _, it := range items {
		counts[it.label]++
	}
	for i := range items {
		if counts[items[i].label] > 1 {
			items[i].label += " (" + sourceDisambiguator(sources[i]) + ")"
		}
	}

	groups := []renderer.TOCGroup{
		groupWorkflow:  {Heading: "Workflows"},
		groupReusable:  {Heading: "Reusable workflows"},
		groupComposite: {Heading: "Composite actions"},
	}
	for _, it := range items {
		groups[it.group].Entries = append(groups[it.group].Entries,
			renderer.TOCEntry{Label: it.label, Anchor: it.anchor})
	}

	header = renderer.RenderDocumentHeader(documentTitle(inputPath), nWorkflow, nReusable, nComposite)
	toc = renderer.RenderTOC(groups)
	return header, toc
}

// titleOf returns a source's display name (workflow or action).
func titleOf(s callgraph.Source) string {
	if s.Workflow != nil {
		return s.Workflow.Name
	}
	return s.Action.Name
}

// sourceDisambiguator returns the path fragment that distinguishes a source from a
// same-named sibling in the TOC. Workflows use their filename; composite actions use the
// last two segments of their containing directory, since every action metadata file is
// named action.yml/action.yaml and actions are discovered at any depth -- a single leaf
// directory name can repeat across depths, so two segments keep the label unambiguous while
// staying compact.
func sourceDisambiguator(s callgraph.Source) string {
	if s.Action != nil {
		return lastTwoPathSegments(filepath.Dir(s.Path))
	}
	return filepath.Base(s.Path)
}

// lastTwoPathSegments returns the final two segments of a slash- or OS-separated path joined
// with "/", or the whole path when it has fewer than two segments.
func lastTwoPathSegments(p string) string {
	parts := strings.Split(filepath.ToSlash(filepath.Clean(p)), "/")
	if len(parts) >= 2 {
		return parts[len(parts)-2] + "/" + parts[len(parts)-1]
	}
	return parts[len(parts)-1]
}

// documentTitle derives the document/repo title from the scanned path. A path of the form
// .../<repo>/.github/workflows yields <repo>; anything else falls back to a generic title.
func documentTitle(inputPath string) string {
	p := filepath.ToSlash(filepath.Clean(inputPath))
	const suffix = ".github/workflows"
	if p == suffix || !strings.HasSuffix(p, "/"+suffix) {
		return "GitHub Actions"
	}
	repoPath := strings.TrimSuffix(p, "/"+suffix)
	if i := strings.LastIndex(repoPath, "/"); i >= 0 {
		repoPath = repoPath[i+1:]
	}
	if repoPath == "" || repoPath == "." {
		return "GitHub Actions"
	}
	return repoPath
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
				// A fully commented-out file is a disabled workflow, not a broken one:
				// note it and move on without failing the run.
				if errors.Is(err, parser.ErrOnlyComments) {
					fmt.Fprintf(os.Stderr, "note: skipping %v\n", err)
					continue
				}
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
	// WalkDir already visits in lexical order, but sort explicitly so deterministic
	// output (TOC + section ordering) doesn't silently depend on that implementation
	// detail.
	sort.Strings(out)
	return out, err
}

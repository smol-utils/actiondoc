package parser

import (
	"strings"

	"github.com/smol-utils/actiondoc/internal/model"
)

// ParseTags extracts ActionDoc tags from a comment string.
// Each line is prefixed with "# " by the YAML library.
func ParseTags(comment string) model.Tags {
	if comment == "" {
		return model.Tags{}
	}

	var tags model.Tags
	var currentTag string
	var currentValue strings.Builder

	lines := strings.Split(comment, "\n")
	for _, line := range lines {
		line = stripCommentPrefix(line)

		if tag, rest, ok := parseTagLine(line); ok {
			flushTag(&tags, currentTag, &currentValue)
			currentTag = tag
			currentValue.WriteString(rest)
		} else if strings.HasPrefix(line, "@") {
			// An @-prefixed line that is not in the allowlist (e.g. another tool's
			// markers like Lula's @lulaStart/@lulaEnd). Ignore it entirely -- do NOT
			// treat it as a continuation, so it never leaks into a preceding tag value.
			continue
		} else if currentTag != "" {
			// Continuation line -- append to current tag.
			if currentValue.Len() > 0 {
				currentValue.WriteString("\n")
			}
			currentValue.WriteString(line)
		}
	}
	flushTag(&tags, currentTag, &currentValue)

	return tags
}

// stripCommentPrefix removes the leading "# " or "#" from a comment line.
func stripCommentPrefix(line string) string {
	line = strings.TrimSpace(line)
	if strings.HasPrefix(line, "# ") {
		return line[2:]
	}
	if strings.HasPrefix(line, "#") {
		return line[1:]
	}
	return line
}

// tagAppliers is the single source of truth for the recognized @-tag allowlist: it maps
// each tag name to the function that stores its value into the Tags struct. parseTagLine
// uses it for membership; flushTag dispatches through it. Adding a tag is a one-line
// change here.
var tagAppliers = map[string]func(t *model.Tags, v string){
	"desc":       func(t *model.Tags, v string) { t.Desc = v },
	"secret":     func(t *model.Tags, v string) { t.Secrets = append(t.Secrets, parseParam(v)) },
	"input":      func(t *model.Tags, v string) { t.Inputs = append(t.Inputs, parseParam(v)) },
	"env":        func(t *model.Tags, v string) { t.Envs = append(t.Envs, parseParam(v)) },
	"output":     func(t *model.Tags, v string) { t.Outputs = append(t.Outputs, parseParam(v)) },
	"deprecated": func(t *model.Tags, v string) { t.Deprecated = v },
	"see":        func(t *model.Tags, v string) { t.See = append(t.See, v) },
	"since":      func(t *model.Tags, v string) { t.Since = v },
	"example":    func(t *model.Tags, v string) { t.Example = v },
}

// parseTagLine checks if a line starts with an allowlisted @tagname and returns the tag
// and remainder.
func parseTagLine(line string) (tag, rest string, ok bool) {
	if !strings.HasPrefix(line, "@") {
		return "", "", false
	}
	// Split "@tagname rest of line"
	parts := strings.SplitN(line, " ", 2)
	tag = parts[0][1:] // strip leading "@"
	if len(parts) == 2 {
		rest = parts[1]
	}
	if _, ok := tagAppliers[tag]; ok {
		return tag, rest, true
	}
	return "", "", false
}

// flushTag stores the accumulated value for the current tag into the Tags struct.
func flushTag(tags *model.Tags, tag string, value *strings.Builder) {
	if tag == "" {
		return
	}
	v := value.String()
	if tag == "example" {
		// Preserve indentation for examples; only trim trailing whitespace.
		v = strings.TrimRight(v, " \t\n")
	} else {
		v = strings.TrimSpace(v)
	}
	if apply, ok := tagAppliers[tag]; ok {
		apply(tags, v)
	}
	value.Reset()
}

// parseParam parses "name {type} - description" or "name - description" or just "name".
func parseParam(s string) model.Param {
	var p model.Param

	rest := s
	idx := strings.IndexAny(rest, " \t")
	if idx == -1 {
		p.Name = rest
		return p
	}
	p.Name = rest[:idx]
	rest = strings.TrimSpace(rest[idx:])

	// Optional {type}.
	if strings.HasPrefix(rest, "{") {
		end := strings.Index(rest, "}")
		if end != -1 {
			p.Type = rest[1:end]
			rest = strings.TrimSpace(rest[end+1:])
		}
	}

	// Optional "- description".
	if strings.HasPrefix(rest, "-") {
		rest = strings.TrimSpace(rest[1:])
	}

	p.Description = rest
	return p
}

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

// parseTagLine checks if a line starts with @tagname and returns the tag and remainder.
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
	switch tag {
	case "desc", "secret", "input", "env", "output", "deprecated", "see", "since", "example":
		return tag, rest, true
	default:
		return "", "", false
	}
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
	switch tag {
	case "desc":
		tags.Desc = v
	case "secret":
		tags.Secrets = append(tags.Secrets, parseParam(v))
	case "input":
		tags.Inputs = append(tags.Inputs, parseParam(v))
	case "env":
		tags.Envs = append(tags.Envs, parseParam(v))
	case "output":
		tags.Outputs = append(tags.Outputs, parseParam(v))
	case "deprecated":
		tags.Deprecated = v
	case "see":
		tags.See = append(tags.See, v)
	case "since":
		tags.Since = v
	case "example":
		tags.Example = v
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

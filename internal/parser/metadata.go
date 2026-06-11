package parser

import (
	"strings"

	"github.com/goccy/go-yaml/ast"
	"github.com/smol-utils/actiondoc/internal/model"
)

// parseConcurrency parses a concurrency: block. It accepts both the scalar form
// (concurrency: my-group) and the map form
// (concurrency: { group: ..., cancel-in-progress: true }). Returns nil when no group
// resolves.
func parseConcurrency(node ast.Node) *model.Concurrency {
	if node == nil {
		return nil
	}
	if scalar := scalarString(node); scalar != "" {
		return &model.Concurrency{Group: scalar}
	}
	mapping := toMapping(node)
	if mapping == nil {
		return nil
	}
	c := &model.Concurrency{}
	for _, mv := range mapping.Values {
		switch mapKeyString(mv.Key) {
		case "group":
			c.Group = nodeString(mv.Value)
		case "cancel-in-progress":
			c.CancelInProgress = strings.TrimSpace(nodeString(mv.Value))
		}
	}
	if c.Group == "" {
		return nil
	}
	return c
}

// parseDefaults parses a defaults: block, reading the run: sub-map's shell and
// working-directory. Returns nil when neither is present.
func parseDefaults(node ast.Node) *model.Defaults {
	run := childMapping(node, "run")
	if run == nil {
		return nil
	}
	d := &model.Defaults{}
	for _, mv := range run.Values {
		switch mapKeyString(mv.Key) {
		case "shell":
			d.Shell = nodeString(mv.Value)
		case "working-directory":
			d.WorkingDirectory = nodeString(mv.Value)
		}
	}
	if d.Shell == "" && d.WorkingDirectory == "" {
		return nil
	}
	return d
}

// implicitDescription derives a description from a leading comment block when no explicit
// ActionDoc tags were found. It returns "" when the comment is a license header (which
// should not masquerade as documentation) or when tags are already present.
func implicitDescription(headComment string, tags model.Tags) string {
	if !tags.Empty() {
		return ""
	}
	// If the block contains ANY @-marker, it belongs to another tool (e.g. a Lula
	// @lulaStart...@lulaEnd block); the marker's surrounding prose is that tool's content,
	// not an ActionDoc description. Suppress the implicit description entirely rather than
	// guess which interior lines are "ours" -- per the allowlist rule, unknown markers and
	// their blocks are ignored, not rendered as prose.
	for _, line := range strings.Split(headComment, "\n") {
		if strings.HasPrefix(strings.TrimSpace(stripCommentPrefix(line)), "@") {
			return ""
		}
	}
	text := cleanCommentText(headComment)
	if text == "" || isLicenseHeader(text) {
		return ""
	}
	return text
}

// licenseMarkers are substrings that identify a leading comment block as a license or
// copyright header rather than a human-facing description. The Apache Software Foundation
// header is the common driver; SPDX and copyright lines cover the rest.
var licenseMarkers = []string{
	"Licensed to the Apache Software Foundation",
	"SPDX-License-Identifier",
	"Licensed under the",
	"http://www.apache.org/licenses/",
}

// isLicenseHeader reports whether cleaned comment text is a license/copyright header.
func isLicenseHeader(text string) bool {
	for _, m := range licenseMarkers {
		if strings.Contains(text, m) {
			return true
		}
	}
	return strings.HasPrefix(text, "Copyright ")
}

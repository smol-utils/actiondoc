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
	if !tagsEmpty(tags) {
		return ""
	}
	// Drop any @-marker lines before forming the description. Allowlisted tags are already
	// captured in `tags`; unknown markers (e.g. another tool's Lula blocks) must be ignored
	// per the allowlist rule, not rendered as prose. If the block is nothing but markers,
	// there is no implicit description.
	var prose []string
	for _, line := range strings.Split(headComment, "\n") {
		if strings.HasPrefix(strings.TrimSpace(stripCommentPrefix(line)), "@") {
			continue
		}
		prose = append(prose, line)
	}
	text := cleanCommentText(strings.Join(prose, "\n"))
	if text == "" || isLicenseHeader(text) {
		return ""
	}
	return text
}

// tagsEmpty reports whether no ActionDoc tag was recognized in a comment block.
func tagsEmpty(t model.Tags) bool {
	return t.Desc == "" && t.Deprecated == "" && t.Since == "" && t.Example == "" &&
		len(t.Secrets) == 0 && len(t.Inputs) == 0 && len(t.Envs) == 0 &&
		len(t.Outputs) == 0 && len(t.See) == 0
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

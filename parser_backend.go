package main

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclsyntax"
)

func parseBackendBlock(block *hclsyntax.Block) (*Workspace, error) {
	// On backend block, label must be "remote"
	if !isRemoteBackend(block.Labels) {
		return nil, fmt.Errorf(`Backend label must be "remote"`)
	}

	hostname, err := getAttribute(block.Body.Attributes, "hostname")
	if err != nil {
		return nil, nil
	}
	// hostname must be "app.terraform.io"
	if hostname != "app.terraform.io" {
		return nil, nil
	}

	org, err := getAttribute(block.Body.Attributes, "organization")
	if err != nil {
		return nil, fmt.Errorf("Failed to get organization attribute: %w", err)
	}

	for _, b := range block.Body.Blocks {
		if b.Type != "workspaces" {
			continue
		}
		workspace, err := getAttribute(b.Body.Attributes, "name")
		if err != nil {
			return nil, fmt.Errorf("name attribute not found in backend.workspaces")
		}
		return &Workspace{
			Organization: org,
			Workspace:    workspace,
		}, nil
	}

	return nil, fmt.Errorf("No enough information got in terraform block")
}

func isRemoteBackend(labels []string) bool {
	for i := range labels {
		if labels[i] == "remote" {
			return true
		}
	}
	return false
}

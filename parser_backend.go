package main

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/pkg/errors"
)

const HCPTerraformHost = "app.terraform.io"

func parseBackendBlock(block *hclsyntax.Block) (*Workspace, error) {
	// On backend block, label must be "remote"
	if !isRemoteBackend(block.Labels) {
		return nil, errors.WithStack(fmt.Errorf(`backend labels must have "remote"`))
	}

	hostname, err := getAttribute(block.Body.Attributes, "hostname")
	if err != nil {
		return nil, nil
	}
	// hostname must be "app.terraform.io"
	if hostname != HCPTerraformHost {
		return nil, nil
	}

	org, err := getAttribute(block.Body.Attributes, "organization")
	if err != nil {
		return nil, errors.WithStack(fmt.Errorf("failed to get organization attribute: %w", err))
	}

	for _, b := range block.Body.Blocks {
		if b.Type != "workspaces" {
			continue
		}
		workspace, err := getAttribute(b.Body.Attributes, "name")
		if err != nil {
			return nil, errors.WithStack(fmt.Errorf("name attribute not found in backend.workspaces"))
		}
		return &Workspace{
			Organization: org,
			Workspace:    workspace,
		}, nil
	}

	return nil, errors.WithStack(fmt.Errorf("not enough informations got in terraform block"))
}

func isRemoteBackend(labels []string) bool {
	for i := range labels {
		if labels[i] == "remote" {
			return true
		}
	}
	return false
}

package main

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/pkg/errors"
)

func parseCloudBlock(block *hclsyntax.Block) (*Workspace, error) {
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
		name, err := getAttribute(b.Body.Attributes, "name")
		if err != nil {
			return nil, errors.WithStack(fmt.Errorf("failed to get name attribute in cloud.workspaces: %w", err))
		}
		return &Workspace{
			Organization: org,
			Workspace:    name,
		}, nil
	}

	return nil, errors.WithStack(fmt.Errorf("not enough information got in terraform block"))
}

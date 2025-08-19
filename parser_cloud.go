package main

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclsyntax"
)

func parseCloudBlock(block *hclsyntax.Block) (*Workspace, error) {
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
		attr, ok := b.Body.Attributes["project"]
		if !ok {
			return nil, fmt.Errorf("project attribute not found in cloud.workspaces")
		}
		ws, err := evalExpression(attr.Expr)
		if err != nil {
			return nil, err
		}
		return &Workspace{
			Organization: org,
			Workspace:    ws,
		}, nil
	}

	return nil, fmt.Errorf("No enough information got in terraform block")
}

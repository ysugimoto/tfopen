package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

const (
	blockTypeTerraform = "terraform"
	blockTypeCloud     = "cloud"
	blockTypeBackend   = "backend"
)

type Workspace struct {
	Organization string
	Workspace    string
}

func (w *Workspace) toURL() string {
	return fmt.Sprintf(
		"https://app.terraform.io/app/%s/workspaces/%s",
		w.Organization,
		w.Workspace,
	)
}

func parseHCLFile(p *hclparse.Parser, filename string) (*Workspace, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	file, diag := p.ParseHCL(content, filename)
	if diag.HasErrors() {
		return nil, fmt.Errorf("HCL file parse error: %s", diag.Error())
	}

	hcl, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		return nil, fmt.Errorf("Failed to assert to hclsyntax.Body")
	}

	for _, block := range hcl.Blocks {
		// We will process for "terraform" block only
		if block.Type != blockTypeTerraform {
			continue
		}

		for _, tf := range block.Body.Blocks {
			switch tf.Type {
			case blockTypeCloud:
				ws, err := parseCloudBlock(tf)
				if err != nil {
					return nil, fmt.Errorf("Failed to parse cloud block: %w", err)
				}
				return ws, nil
			case blockTypeBackend:
				ws, err := parseBackendBlock(tf)
				if err != nil {
					return nil, fmt.Errorf("Failed to parse backend block: %w", err)
				}
				return ws, nil
			}
		}
	}

	return nil, nil
}

package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

func TestParseBackendBlock(t *testing.T) {
	testCases := []struct {
		name        string
		hcl         string
		expectedW   *Workspace
		expectedErr bool
	}{
		{
			name: "valid remote backend",
			hcl: `
backend "remote" {
  hostname = "app.terraform.io"
  organization = "dummy_org"

  workspaces {
    name = "dummy_workspace"
  }
}
`,
			expectedW: &Workspace{
				Organization: "dummy_org",
				Workspace:    "dummy_workspace",
			},
			expectedErr: false,
		},
		{
			name:        "not remote backend",
			hcl:         `backend "local" {}`,
			expectedW:   nil,
			expectedErr: true,
		},
		{
			name: "hostname is not app.terraform.io",
			hcl: `
backend "remote" {
  hostname = "example.com"
  organization = "dummy_org"

  workspaces {
    name = "dummy_workspace"
  }
}
`,
			expectedW:   nil,
			expectedErr: false,
		},
		{
			name: "missing organization",
			hcl: `
backend "remote" {
  hostname = "app.terraform.io"

  workspaces {
    name = "dummy_workspace"
  }
}
`,
			expectedW:   nil,
			expectedErr: true,
		},
		{
			name: "missing workspaces block",
			hcl: `
backend "remote" {
  hostname = "app.terraform.io"
  organization = "dummy_org"
}
`,
			expectedW:   nil,
			expectedErr: true,
		},
		{
			name: "missing name in workspaces",
			hcl: `
backend "remote" {
  hostname = "app.terraform.io"
  organization = "dummy_org"

  workspaces {}
}
`,
			expectedW:   nil,
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := hclparse.NewParser()
			file, diags := p.ParseHCL([]byte(tc.hcl), "test.tf")
			if diags.HasErrors() {
				t.Fatalf("failed to parse hcl: %s", diags.Error())
			}
			body, ok := file.Body.(*hclsyntax.Body)
			if !ok {
				t.Fatalf("failed to assert to hclsyntax.Body")
			}

			// Handle empty file case
			if len(body.Blocks) == 0 {
				return
			}

			ws, err := parseBackendBlock(body.Blocks[0])
			if tc.expectedErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %s", err)
				}
			}
			if diff := cmp.Diff(tc.expectedW, ws); diff != "" {
				t.Errorf("workspace mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

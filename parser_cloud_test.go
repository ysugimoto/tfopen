package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

func TestParseCloudBlock(t *testing.T) {
	testCases := []struct {
		name        string
		hcl         string
		expectedW   *Workspace
		expectedErr bool
	}{
		{
			name: "valid cloud block",
			hcl: `
cloud {
  hostname     = "app.terraform.io"
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
			name: "hostname is not app.terraform.io",
			hcl: `
cloud {
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
cloud {
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
cloud {
  hostname = "app.terraform.io"
  organization = "dummy_org"
}
`,
			expectedW:   nil,
			expectedErr: true,
		},
		{
			name: "missing project in workspaces",
			hcl: `
cloud {
  hostname = "app.terraform.io"
  organization = "dummy_org"

  workspaces {}
}
`,
			expectedW:   nil,
			expectedErr: true,
		},
		{
			name: "workspaces has tags instead of project",
			hcl: `
cloud {
  hostname = "app.terraform.io"
  organization = "dummy_org"

  workspaces {
    tags = ["a", "b"]
  }
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

			ws, err := parseCloudBlock(body.Blocks[0])
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

package main

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/hcl/v2/hclparse"
)

func TestParseHCLFile(t *testing.T) {
	testCases := []struct {
		name          string
		filename      string
		expectedW     *Workspace
		expectedErr   bool
		expectedErrContains string
	}{
		{
			name:     "cloud block",
			filename: filepath.Join("__test__", "cloud.tf"),
			expectedW: &Workspace{
				Organization: "dummy_org",
				Workspace:    "dummy_workspace",
			},
			expectedErr: false,
		},
		{
			name:     "backend block",
			filename: filepath.Join("__test__", "backend.tf"),
			expectedW: &Workspace{
				Organization: "dummy_org",
				Workspace:    "dummy_workspace",
			},
			expectedErr: false,
		},
		{
			name:        "no cloud or backend block",
			filename:    filepath.Join("__test__", "empty.tf"),
			expectedW:   nil,
			expectedErr: false,
		},
		{
			name:        "file not found",
			filename:    "nonexistent.tf",
			expectedW:   nil,
			expectedErr: true,
			expectedErrContains: "no such file or directory",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := hclparse.NewParser()
			ws, err := parseHCLFile(p, tc.filename)

			if tc.expectedErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				if tc.expectedErrContains != "" && !strings.Contains(err.Error(), tc.expectedErrContains) {
					t.Errorf("expected error to contain '%s' but got '%s'", tc.expectedErrContains, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %s", err)
				}
				if diff := cmp.Diff(tc.expectedW, ws); diff != "" {
					t.Errorf("workspace mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}
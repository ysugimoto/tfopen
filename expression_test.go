package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
)

func TestEvalExpression(t *testing.T) {
	testCases := []struct {
		name          string
		expr          hclsyntax.Expression
		expectedVal   string
		expectedErr   bool
	}{
		{
			name: "string literal",
			expr: &hclsyntax.LiteralValueExpr{
				Val: cty.StringVal("hello"),
			},
			expectedVal: "hello",
			expectedErr: false,
		},
		{
			name: "template expression",
			expr: &hclsyntax.TemplateExpr{
				Parts: []hclsyntax.Expression{
					&hclsyntax.LiteralValueExpr{
						Val: cty.StringVal("hello"),
					},
				},
			},
			expectedVal: "hello",
			expectedErr: false,
		},
		{
			name: "template expression with interpolation",
			expr: &hclsyntax.TemplateExpr{
				Parts: []hclsyntax.Expression{
					&hclsyntax.LiteralValueExpr{
						Val: cty.StringVal("hello "),
					},
					&hclsyntax.TemplateWrapExpr{
						Wrapped: &hclsyntax.ScopeTraversalExpr{},
					},
				},
			},
			expectedVal: "",
			expectedErr: true,
		},
		{
			name: "number literal",
			expr: &hclsyntax.LiteralValueExpr{
				Val: cty.NumberIntVal(123),
			},
			expectedVal: "",
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			val, err := evalExpression(tc.expr)

			if tc.expectedErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %s", err)
				}
				if diff := cmp.Diff(tc.expectedVal, val); diff != "" {
					t.Errorf("value mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

func TestGetAttribute(t *testing.T) {
	attributes := hclsyntax.Attributes{
		"name": &hclsyntax.Attribute{
			Name: "name",
			Expr: &hclsyntax.LiteralValueExpr{
				Val: cty.StringVal("test-workspace"),
			},
		},
		"invalid": &hclsyntax.Attribute{
			Name: "invalid",
			Expr: &hclsyntax.LiteralValueExpr{
				Val: cty.NumberIntVal(123),
			},
		},
	}

	testCases := []struct {
		name          string
		attrName      string
		expectedVal   string
		expectedErr   bool
	}{
		{
			name:        "valid attribute",
			attrName:    "name",
			expectedVal: "test-workspace",
			expectedErr: false,
		},
		{
			name:        "missing attribute",
			attrName:    "nonexistent",
			expectedVal: "",
			expectedErr: true,
		},
		{
			name:        "invalid expression",
			attrName:    "invalid",
			expectedVal: "",
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			val, err := getAttribute(attributes, tc.attrName)

			if tc.expectedErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %s", err)
				}
				if diff := cmp.Diff(tc.expectedVal, val); diff != "" {
					t.Errorf("value mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}

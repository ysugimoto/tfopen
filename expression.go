package main

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
)

func getAttribute(attributes hclsyntax.Attributes, name string) (string, error) {
	attr, ok := attributes[name]
	if !ok {
		return "", fmt.Errorf("attribute %s is not found in attributes", name)
	}
	val, err := evalExpression(attr.Expr)
	if err != nil {
		return "", fmt.Errorf("Failed to evaluate expression: %w", err)
	}
	return val, nil
}

func evalExpression(expr hcl.Expression) (string, error) {
	switch t := expr.(type) {
	case *hclsyntax.LiteralValueExpr:
		if t.Val.Type() != cty.String {
			return "", fmt.Errorf("Expression must be a string literal")
		}
		return t.Val.AsString(), nil
	case *hclsyntax.TemplateExpr:
		if len(t.Parts) == 1 {
			l, ok := t.Parts[0].(*hclsyntax.LiteralValueExpr)
			if !ok {
				return "", fmt.Errorf("Expression must be a string literal")
			}
			if l.Val.Type() != cty.String {
				return "", fmt.Errorf("Expression must be a string literal")
			}
			return l.Val.AsString(), nil
		}
	}
	return "", fmt.Errorf("Expression must be a string literal")
}

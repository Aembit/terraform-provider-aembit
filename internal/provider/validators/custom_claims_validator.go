package validators

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// CustomClaimsValidator validates custom claims sets on credential provider resources at plan time.
type CustomClaimsValidator struct {
	// dynamicVarRegex matches canonical dynamic template variable syntax ${...}.
	dynamicVarRegex *regexp.Regexp

	// literalPatternRegex matches any template expression syntax ${...} forbidden in literal values.
	literalPatternRegex *regexp.Regexp
}

// customClaimModel represents a single custom claim element within the custom_claims set attribute.
type customClaimModel struct {
	// Key represents the custom claim key name.
	Key types.String `tfsdk:"key"`

	// Value represents the custom claim value expression or literal.
	Value types.String `tfsdk:"value"`

	// ValueType represents whether the custom claim value is literal or dynamic.
	ValueType types.String `tfsdk:"value_type"`
}

// Description returns a plain text description of the validator.
func (v CustomClaimsValidator) Description(_ context.Context) string {
	return "Validates that custom claim values match their specified value_type (literal or dynamic)."
}

// MarkdownDescription returns a markdown-formatted description of the validator.
func (v CustomClaimsValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// formatFieldName constructs the diagnostic field name string for a custom claim key.
func formatFieldName(key string) string {
	if key == "" {
		return "CustomClaim"
	}
	return fmt.Sprintf("CustomClaim '%s'", key)
}

// ValidateSet executes plan-time validation on the custom_claims set attribute.
func (v CustomClaimsValidator) ValidateSet(
	ctx context.Context,
	req validator.SetRequest,
	resp *validator.SetResponse,
) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	dynamicRegex := v.dynamicVarRegex
	if dynamicRegex == nil {
		dynamicRegex = regexp.MustCompile(`\$\{[\\\w\-_.\[\]'":/]+\}`)
	}

	literalRegex := v.literalPatternRegex
	if literalRegex == nil {
		literalRegex = regexp.MustCompile(`\$\{[^}]*\}`)
	}

	for _, elem := range req.ConfigValue.Elements() {
		if elem.IsNull() || elem.IsUnknown() {
			continue
		}

		obj, ok := elem.(types.Object)
		if !ok {
			continue
		}

		var claim customClaimModel
		diags := obj.As(ctx, &claim, basetypes.ObjectAsOptions{})
		resp.Diagnostics.Append(diags...)
		if diags.HasError() {
			continue
		}

		if claim.Value.IsNull() || claim.Value.IsUnknown() ||
			claim.ValueType.IsNull() || claim.ValueType.IsUnknown() {
			continue
		}

		key := ""
		if !claim.Key.IsNull() && !claim.Key.IsUnknown() {
			key = claim.Key.ValueString()
		}

		val := claim.Value.ValueString()
		valType := strings.TrimSpace(claim.ValueType.ValueString())

		if strings.EqualFold(valType, "literal") {
			if val != "" && literalRegex.MatchString(val) {
				resp.Diagnostics.AddAttributeError(
					req.Path,
					"Invalid Custom Claim Value",
					fmt.Sprintf("Field '%s' cannot contain dynamic template syntax '${...}' when configured as literal.", formatFieldName(key)),
				)
			}
		} else if strings.EqualFold(valType, "dynamic") {
			if val == "" {
				resp.Diagnostics.AddAttributeError(
					req.Path,
					"Invalid Custom Claim Value",
					fmt.Sprintf("Field '%s' must contain a valid template expression '${...}' when configured as dynamic.", formatFieldName(key)),
				)
			} else if strings.Contains(val, "${") {
				if !dynamicRegex.MatchString(val) || strings.Contains(dynamicRegex.ReplaceAllString(val, ""), "${") {
					resp.Diagnostics.AddAttributeError(
						req.Path,
						"Invalid Custom Claim Value",
						fmt.Sprintf("Field '%s' has invalid template expression syntax.", formatFieldName(key)),
					)
				}
			} else {
				resp.Diagnostics.AddAttributeError(
					req.Path,
					"Invalid Custom Claim Value",
					fmt.Sprintf("Field '%s' must contain a valid template expression '${...}' when configured as dynamic.", formatFieldName(key)),
				)
			}
		}
	}
}

// NewCustomClaimsValidator returns a new instance of the custom claims validator.
func NewCustomClaimsValidator() validator.Set {
	return CustomClaimsValidator{
		dynamicVarRegex:     regexp.MustCompile(`\$\{[\\\w\-_.\[\]'":/]+\}`),
		literalPatternRegex: regexp.MustCompile(`\$\{[^}]*\}`),
	}
}

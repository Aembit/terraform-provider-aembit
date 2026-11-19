package validators

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// invalidCustomClaimValueSummary is the diagnostic error summary for custom claim value errors.
const invalidCustomClaimValueSummary = "Invalid Custom Claim Value"

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

// keyString returns the claim key as a string or empty if null or unknown.
func (c *customClaimModel) keyString() string {
	if c.Key.IsNull() || c.Key.IsUnknown() {
		return ""
	}
	return c.Key.ValueString()
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

// withDefaults ensures compiled regular expressions are initialized.
func (v CustomClaimsValidator) withDefaults() CustomClaimsValidator {
	if v.dynamicVarRegex == nil {
		v.dynamicVarRegex = regexp.MustCompile(`\$\{[\\\w\-_.\[\]'":/]+\}`)
	}
	if v.literalPatternRegex == nil {
		v.literalPatternRegex = regexp.MustCompile(`\$\{[^}]*\}`)
	}
	return v
}

// extractClaim handles object casting, null/unknown checks, and model mapping for a set element.
func (v CustomClaimsValidator) extractClaim(
	ctx context.Context,
	elem attr.Value,
	resp *validator.SetResponse,
) (*customClaimModel, bool) {
	if elem.IsNull() || elem.IsUnknown() {
		return nil, false
	}

	obj, ok := elem.(types.Object)
	if !ok {
		return nil, false
	}

	var claim customClaimModel
	diags := obj.As(ctx, &claim, basetypes.ObjectAsOptions{})
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return nil, false
	}

	if claim.Value.IsNull() || claim.Value.IsUnknown() ||
		claim.ValueType.IsNull() || claim.ValueType.IsUnknown() {
		return nil, false
	}

	return &claim, true
}

// claimValidationError represents a diagnostic validation error for a custom claim.
type claimValidationError struct {
	// message describes the specific claim validation failure.
	message string
}

// Error returns the validation error message.
func (e claimValidationError) Error() string {
	return e.message
}

// validateLiteralClaim validates that a literal claim value does not contain template syntax.
func (v CustomClaimsValidator) validateLiteralClaim(key, val string) error {
	if val != "" && v.literalPatternRegex.MatchString(val) {
		return claimValidationError{
			message: fmt.Sprintf("Field '%s' cannot contain dynamic template syntax '${...}' when configured as literal.", formatFieldName(key)),
		}
	}
	return nil
}

// validateDynamicClaim validates that a dynamic claim value contains valid template expression syntax.
func (v CustomClaimsValidator) validateDynamicClaim(key, val string) error {
	if val == "" || !strings.Contains(val, "${") {
		return claimValidationError{
			message: fmt.Sprintf("Field '%s' must contain a valid template expression '${...}' when configured as dynamic.", formatFieldName(key)),
		}
	}

	if !v.dynamicVarRegex.MatchString(val) || strings.Contains(v.dynamicVarRegex.ReplaceAllString(val, ""), "${") {
		return claimValidationError{
			message: fmt.Sprintf("Field '%s' has invalid template expression syntax.", formatFieldName(key)),
		}
	}

	return nil
}

// validateClaimValue delegates claim value validation according to its declared value type.
func (v CustomClaimsValidator) validateClaimValue(key, val, valType string) error {
	if strings.EqualFold(valType, "literal") {
		return v.validateLiteralClaim(key, val)
	}
	if strings.EqualFold(valType, "dynamic") {
		return v.validateDynamicClaim(key, val)
	}
	return nil
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

	v = v.withDefaults()

	for _, elem := range req.ConfigValue.Elements() {
		claim, ok := v.extractClaim(ctx, elem, resp)
		if !ok {
			continue
		}

		key := claim.keyString()
		val := claim.Value.ValueString()
		valType := strings.TrimSpace(claim.ValueType.ValueString())

		if err := v.validateClaimValue(key, val, valType); err != nil {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				invalidCustomClaimValueSummary,
				err.Error(),
			)
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

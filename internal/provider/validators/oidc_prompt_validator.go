package validators

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// oidcPromptValidator validates that oidc.prompt adheres to OIDC Core 1.0 Section 3.1.2.1 rules.
type oidcPromptValidator struct{}

var _ validator.String = oidcPromptValidator{}

// OidcPromptValidation returns a validator.String that validates OIDC prompt configuration.
// It verifies that the value is a single-space-delimited list of allowed OIDC Core 1.0 tokens
// (none, login, consent, select_account), contains no duplicate tokens, and never combines none
// with any other value.
func OidcPromptValidation() validator.String {
	return oidcPromptValidator{}
}

// Description returns a plain text description of the validator.
func (v oidcPromptValidator) Description(_ context.Context) string {
	return "Validates that the OIDC prompt parameter consists of valid space-delimited OIDC Core 1.0 values (none, login, consent, select_account), contains no duplicates, and does not combine 'none' with other values."
}

// MarkdownDescription returns a markdown description of the validator.
func (v oidcPromptValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

// isAllowedOidcPrompt reports whether token is a valid OIDC Core 1.0 prompt value.
func isAllowedOidcPrompt(token string) bool {
	switch token {
	case "none", "login", "consent", "select_account":
		return true
	default:
		return false
	}
}

// ValidateString executes plan-time validation on the prompt string attribute.
func (v oidcPromptValidator) ValidateString(
	_ context.Context,
	req validator.StringRequest,
	resp *validator.StringResponse,
) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	val := req.ConfigValue.ValueString()
	tokens := strings.Fields(val)

	// Validate non-empty and canonically single-space-delimited to prevent post-apply state mismatch.
	if len(tokens) == 0 || val != strings.Join(tokens, " ") {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid OIDC Prompt Value",
			fmt.Sprintf("OIDC prompt must be a non-empty, single-space-delimited list of prompt values without leading, trailing, or consecutive whitespace. Got: %q", val),
		)
		return
	}

	seen := make(map[string]struct{}, len(tokens))
	hasNone := false

	for _, token := range tokens {
		if !isAllowedOidcPrompt(token) {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid OIDC Prompt Value",
				fmt.Sprintf("Invalid OIDC prompt value %q. Allowed values are: 'none', 'login', 'consent', 'select_account'.", token),
			)
			return
		}

		if _, exists := seen[token]; exists {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid OIDC Prompt Value",
				fmt.Sprintf("Duplicate OIDC prompt value %q found. Each prompt value must be unique.", token),
			)
			return
		}
		seen[token] = struct{}{}

		if token == "none" {
			hasNone = true
		}
	}

	if hasNone && len(tokens) > 1 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid OIDC Prompt Value",
			"OIDC prompt 'none' cannot be combined with any other prompt values.",
		)
		return
	}
}

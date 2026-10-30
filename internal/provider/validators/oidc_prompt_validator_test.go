package validators

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestOidcPromptValidation_Descriptions(t *testing.T) {
	t.Parallel()

	v := OidcPromptValidation()
	require.NotNil(t, v)

	ctx := context.Background()
	require.NotEmpty(t, v.Description(ctx))
	require.Equal(t, v.Description(ctx), v.MarkdownDescription(ctx))
}

func TestOidcPromptValidation_NullAndUnknown(t *testing.T) {
	t.Parallel()

	v := OidcPromptValidation()
	ctx := context.Background()

	testCases := []struct {
		name  string
		value types.String
	}{
		{
			name:  "null string",
			value: types.StringNull(),
		},
		{
			name:  "unknown string",
			value: types.StringUnknown(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			req := validator.StringRequest{
				Path:        path.Root("oidc").AtName("prompt"),
				ConfigValue: tc.value,
			}
			resp := &validator.StringResponse{
				Diagnostics: diag.Diagnostics{},
			}
			v.ValidateString(ctx, req, resp)
			require.False(t, resp.Diagnostics.HasError(), "Expected null/unknown value to pass validation")
		})
	}
}

func TestOidcPromptValidation_Valid(t *testing.T) {
	t.Parallel()

	v := OidcPromptValidation()
	ctx := context.Background()

	validCases := []struct {
		name  string
		input string
	}{
		{
			name:  "single none",
			input: "none",
		},
		{
			name:  "single login",
			input: "login",
		},
		{
			name:  "single consent",
			input: "consent",
		},
		{
			name:  "single select_account",
			input: "select_account",
		},
		{
			name:  "login and consent",
			input: "login consent",
		},
		{
			name:  "consent and login",
			input: "consent login",
		},
		{
			name:  "login, consent, and select_account",
			input: "login consent select_account",
		},
		{
			name:  "all valid interactive tokens",
			input: "select_account consent login",
		},
	}

	for _, tc := range validCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			req := validator.StringRequest{
				Path:        path.Root("oidc").AtName("prompt"),
				ConfigValue: types.StringValue(tc.input),
			}
			resp := &validator.StringResponse{
				Diagnostics: diag.Diagnostics{},
			}
			v.ValidateString(ctx, req, resp)
			require.False(t, resp.Diagnostics.HasError(), "Expected %q to pass validation, got: %s", tc.input, resp.Diagnostics)
		})
	}
}

func TestOidcPromptValidation_Invalid(t *testing.T) {
	t.Parallel()

	v := OidcPromptValidation()
	ctx := context.Background()

	invalidCases := []struct {
		name          string
		input         string
		expectedError string
	}{
		{
			name:          "empty string",
			input:         "",
			expectedError: "non-empty, single-space-delimited list of prompt values",
		},
		{
			name:          "whitespace only",
			input:         "   ",
			expectedError: "non-empty, single-space-delimited list of prompt values",
		},
		{
			name:          "leading whitespace",
			input:         " login",
			expectedError: "non-empty, single-space-delimited list of prompt values",
		},
		{
			name:          "trailing whitespace",
			input:         "login ",
			expectedError: "non-empty, single-space-delimited list of prompt values",
		},
		{
			name:          "multiple consecutive spaces",
			input:         "login  consent",
			expectedError: "non-empty, single-space-delimited list of prompt values",
		},
		{
			name:          "tab delimiter",
			input:         "login\tconsent",
			expectedError: "non-empty, single-space-delimited list of prompt values",
		},
		{
			name:          "unknown single token",
			input:         "invalid_prompt",
			expectedError: "Invalid OIDC prompt value \"invalid_prompt\". Allowed values are: 'none', 'login', 'consent', 'select_account'.",
		},
		{
			name:          "uppercase token",
			input:         "LOGIN",
			expectedError: "Invalid OIDC prompt value \"LOGIN\". Allowed values are: 'none', 'login', 'consent', 'select_account'.",
		},
		{
			name:          "mixed valid and invalid token",
			input:         "login invalid",
			expectedError: "Invalid OIDC prompt value \"invalid\". Allowed values are: 'none', 'login', 'consent', 'select_account'.",
		},
		{
			name:          "duplicate token login",
			input:         "login login",
			expectedError: "Duplicate OIDC prompt value \"login\" found. Each prompt value must be unique.",
		},
		{
			name:          "duplicate token consent in multi-value",
			input:         "consent login consent",
			expectedError: "Duplicate OIDC prompt value \"consent\" found. Each prompt value must be unique.",
		},
		{
			name:          "none combined with login",
			input:         "none login",
			expectedError: "OIDC prompt 'none' cannot be combined with any other prompt values.",
		},
		{
			name:          "login combined with none",
			input:         "login none",
			expectedError: "OIDC prompt 'none' cannot be combined with any other prompt values.",
		},
		{
			name:          "none combined with multiple tokens",
			input:         "login consent none",
			expectedError: "OIDC prompt 'none' cannot be combined with any other prompt values.",
		},
		{
			name:          "duplicate none",
			input:         "none none",
			expectedError: "Duplicate OIDC prompt value \"none\" found. Each prompt value must be unique.",
		},
	}

	for _, tc := range invalidCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			req := validator.StringRequest{
				Path:        path.Root("oidc").AtName("prompt"),
				ConfigValue: types.StringValue(tc.input),
			}
			resp := &validator.StringResponse{
				Diagnostics: diag.Diagnostics{},
			}
			v.ValidateString(ctx, req, resp)
			require.True(t, resp.Diagnostics.HasError(), "Expected %q to fail validation", tc.input)

			hasExpectedMessage := false
			for _, d := range resp.Diagnostics.Errors() {
				require.Equal(t, "Invalid OIDC Prompt Value", d.Summary())
				if d.Detail() == tc.expectedError || (len(tc.expectedError) > 0 && d.Detail() != "" && (d.Detail() == tc.expectedError || (len(d.Detail()) >= len(tc.expectedError) && d.Detail()[:len(tc.expectedError)] == tc.expectedError))) {
					hasExpectedMessage = true
					break
				}
				// Also check contains if full match not exact
				if d.Detail() != "" && strings.Contains(d.Detail(), tc.expectedError) {
					hasExpectedMessage = true
					break
				}
			}
			require.True(t, hasExpectedMessage, "Expected diagnostic error detail containing %q, got %v", tc.expectedError, resp.Diagnostics.Errors())
		})
	}
}

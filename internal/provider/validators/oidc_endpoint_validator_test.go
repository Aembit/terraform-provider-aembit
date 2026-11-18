package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestOidcEndpointValidation(t *testing.T) {
	t.Parallel()

	v := OidcEndpointValidation()
	require.NotNil(t, v)

	ctx := context.Background()
	require.Equal(t, "OIDC Endpoint must not include well-known endpoint.", v.Description(ctx))
	require.Equal(t, v.Description(ctx), v.MarkdownDescription(ctx))

	testCases := []struct {
		name        string
		value       types.String
		expectError bool
	}{
		{
			name:        "null value",
			value:       types.StringNull(),
			expectError: false,
		},
		{
			name:        "unknown value",
			value:       types.StringUnknown(),
			expectError: false,
		},
		{
			name:        "valid issuer url",
			value:       types.StringValue("https://accounts.google.com"),
			expectError: false,
		},
		{
			name:        "valid issuer url with path",
			value:       types.StringValue("https://auth.example.com/oauth2/default"),
			expectError: false,
		},
		{
			name:        "valid empty string",
			value:       types.StringValue(""),
			expectError: false,
		},
		{
			name:        "invalid exact well-known suffix",
			value:       types.StringValue("https://auth.example.com/.well-known/openid-configuration"),
			expectError: true,
		},
		{
			name:        "invalid nested path containing well-known endpoint",
			value:       types.StringValue("https://auth.example.com/tenant/.well-known/openid-configuration?foo=bar"),
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req := validator.StringRequest{
				Path:        path.Root("oidc").AtName("endpoint"),
				ConfigValue: tc.value,
			}
			resp := &validator.StringResponse{
				Diagnostics: diag.Diagnostics{},
			}

			v.ValidateString(ctx, req, resp)

			if tc.expectError {
				require.True(t, resp.Diagnostics.HasError(), "expected validation error for %v", tc.value)
				require.Equal(t, "Invalid OIDC Endpoint", resp.Diagnostics.Errors()[0].Summary())
				require.Equal(t, "OIDC Endpoint must not include well-known endpoint.", resp.Diagnostics.Errors()[0].Detail())
			} else {
				require.False(t, resp.Diagnostics.HasError(), "unexpected validation error: %v", resp.Diagnostics)
			}
		})
	}
}

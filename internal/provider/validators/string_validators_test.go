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

func TestCimdUrlValidation(t *testing.T) {
	t.Parallel()

	v := CimdUrlValidation()
	require.NotNil(t, v)

	ctx := context.Background()

	// Verify Description and MarkdownDescription
	require.Contains(t, v.Description(ctx), "CIMD client ID URL")
	require.Contains(t, v.MarkdownDescription(ctx), "CIMD client ID URL")

	validInputs := []string{
		"https://mcpjam.com/.well-known/oauth/client-metadata.json",
		"https://example.com/some/path",
		"https://sub.example.com/path/to/resource",
		"https://a.b.co/d",
	}

	for _, input := range validInputs {
		t.Run("Valid_"+input, func(t *testing.T) {
			req := validator.StringRequest{
				ConfigValue: types.StringValue(input),
				Path:        path.Root("cimd_url"),
			}
			resp := &validator.StringResponse{
				Diagnostics: diag.Diagnostics{},
			}
			v.ValidateString(ctx, req, resp)
			require.False(t, resp.Diagnostics.HasError(), "Expected valid CIMD URL, but got error: %s", resp.Diagnostics)
		})
	}

	invalidInputs := []string{
		"http://example.com/path",             // Not https
		"https://example.com",                 // No path
		"https://example.com/",                // Path is just /
		"https://example.com/path?query=1",    // Query string not allowed
		"https://example.com/path#fragment",   // Fragment not allowed
		"https://user@example.com/path",       // User info not allowed
		"https://example.com/path with space", // Spaces not allowed
	}

	for _, input := range invalidInputs {
		t.Run("Invalid_"+input, func(t *testing.T) {
			req := validator.StringRequest{
				ConfigValue: types.StringValue(input),
				Path:        path.Root("cimd_url"),
			}
			resp := &validator.StringResponse{
				Diagnostics: diag.Diagnostics{},
			}
			v.ValidateString(ctx, req, resp)
			require.True(t, resp.Diagnostics.HasError(), "Expected invalid CIMD URL to fail validation, but it succeeded: %s", input)
		})
	}
}

func TestSpiffeSubjectValidation(t *testing.T) {
	t.Parallel()

	v := SpiffeSubjectValidation()
	require.NotNil(t, v)

	ctx := context.Background()

	// Verify Description and MarkdownDescription
	require.NotEmpty(t, v.Description(ctx))
	require.NotEmpty(t, v.MarkdownDescription(ctx))

	validInputs := []string{
		"spiffe://trust-domain-name/path",
		"spiffe://example.org/ns/prod/sa/app",
		"spiffe://foo/bar",
	}

	for _, input := range validInputs {
		t.Run("Valid_"+input, func(t *testing.T) {
			req := validator.StringRequest{
				ConfigValue: types.StringValue(input),
				Path:        path.Root("spiffe_subject"),
			}
			resp := &validator.StringResponse{
				Diagnostics: diag.Diagnostics{},
			}
			v.ValidateString(ctx, req, resp)
			require.False(t, resp.Diagnostics.HasError(), "Expected valid SPIFFE subject, but got error: %s", resp.Diagnostics)
		})
	}

	invalidInputs := []string{
		"http://example.com",
		"https://example.com",
		"not-spiffe://domain/path",
		"spiffe:",
		"spiffe:/",
		"spiffe",
	}

	for _, input := range invalidInputs {
		t.Run("Invalid_"+input, func(t *testing.T) {
			req := validator.StringRequest{
				ConfigValue: types.StringValue(input),
				Path:        path.Root("spiffe_subject"),
			}
			resp := &validator.StringResponse{
				Diagnostics: diag.Diagnostics{},
			}
			v.ValidateString(ctx, req, resp)
			require.True(t, resp.Diagnostics.HasError(), "Expected invalid SPIFFE subject to fail validation, but it succeeded: %s", input)
		})
	}
}

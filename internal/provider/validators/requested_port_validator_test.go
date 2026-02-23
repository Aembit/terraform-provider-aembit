package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/require"
)

func TestRequestedPortEqualsPortForMCPValidator(t *testing.T) {
	testCases := []struct {
		name          string
		appProtocol   string
		port          int64
		requestedPort int64
		expectError   bool
	}{
		{
			name:          "MCP protocol, ports equal",
			appProtocol:   "MCP",
			port:          8080,
			requestedPort: 8080,
			expectError:   false,
		},
		{
			name:          "MCP protocol, ports not equal",
			appProtocol:   "MCP",
			port:          8080,
			requestedPort: 9090,
			expectError:   true,
		},
		{
			name:          "Non-MCP protocol, ports not equal",
			appProtocol:   "HTTP",
			port:          8080,
			requestedPort: 9090,
			expectError:   false,
		},
		{
			name:          "Non-MCP protocol, ports equal",
			appProtocol:   "HTTP",
			port:          8080,
			requestedPort: 8080,
			expectError:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v := requestedPortEqualsPortForMCPValidator{}
			ctx := context.Background()

			schema := schema.Schema{
				Attributes: map[string]schema.Attribute{
					"app_protocol":   schema.StringAttribute{},
					"port":           schema.Int64Attribute{},
					"requested_port": schema.Int64Attribute{},
				},
			}

			// Create a mock config with the schema structure
			configValue := tftypes.NewValue(tftypes.Object{
				AttributeTypes: map[string]tftypes.Type{
					"app_protocol":   tftypes.String,
					"port":           tftypes.Number,
					"requested_port": tftypes.Number,
				},
			}, map[string]tftypes.Value{
				"app_protocol":   tftypes.NewValue(tftypes.String, tc.appProtocol),
				"port":           tftypes.NewValue(tftypes.Number, tc.port),
				"requested_port": tftypes.NewValue(tftypes.Number, tc.requestedPort),
			})

			config := tfsdk.Config{
				Raw:    configValue,
				Schema: schema,
			}

			req := validator.Int64Request{
				ConfigValue: types.Int64Value(tc.requestedPort),
				Config:      config,
				Path:        path.Root("requested_port"),
			}

			resp := &validator.Int64Response{
				Diagnostics: diag.Diagnostics{},
			}

			v.ValidateInt64(ctx, req, resp)

			hasError := resp.Diagnostics.HasError()
			if tc.expectError {
				require.True(t, hasError, "expected error but got none")
			} else {
				require.False(t, hasError, "expected no error but got one: %v", resp.Diagnostics)
			}
		})
	}
}

func TestRequestedPortEqualsPortForMCPValidation(t *testing.T) {
	v := RequestedPortEqualsPortForMCPValidation()
	require.NotNil(t, v, "validator should not be nil")
}

func TestRequestedPortEqualsPortForMCPValidator_Description(t *testing.T) {
	v := requestedPortEqualsPortForMCPValidator{}
	ctx := context.Background()

	desc := v.Description(ctx)
	require.NotEmpty(t, desc, "description should not be empty")
	require.Equal(t, "Ensures requested_port equals port when app_protocol is MCP", desc)
}

func TestRequestedPortEqualsPortForMCPValidator_MarkdownDescription(t *testing.T) {
	v := requestedPortEqualsPortForMCPValidator{}
	ctx := context.Background()

	mdDesc := v.MarkdownDescription(ctx)
	require.NotEmpty(t, mdDesc, "markdown description should not be empty")
	require.Equal(t, v.Description(ctx), mdDesc, "markdown description should match description")
}

package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestRequestedPortEqualsPortForMCPValidator_Description(t *testing.T) {
	v := RequestedPortEqualsPortForMCPValidation()
	ctx := context.Background()

	expected := "Ensures requested_port equals port when app_protocol is MCP"
	if got := v.Description(ctx); got != expected {
		t.Errorf("Description() = %v, want %v", got, expected)
	}
}

func TestRequestedPortEqualsPortForMCPValidator_MarkdownDescription(t *testing.T) {
	v := RequestedPortEqualsPortForMCPValidation()
	ctx := context.Background()

	expected := "Ensures requested_port equals port when app_protocol is MCP"
	if got := v.MarkdownDescription(ctx); got != expected {
		t.Errorf("MarkdownDescription() = %v, want %v", got, expected)
	}
}

func TestRequestedPortEqualsPortForMCPValidator_ValidateInt64(t *testing.T) {
	tests := []struct {
		name            string
		requestedPort   types.Int64
		port            types.Int64
		appProtocol     types.String
		expectedError   bool
		expectedSummary string
		expectedDetail  string
	}{
		{
			name:          "valid - MCP with matching ports",
			requestedPort: types.Int64Value(443),
			port:          types.Int64Value(443),
			appProtocol:   types.StringValue("MCP"),
			expectedError: false,
		},
		{
			name:            "invalid - MCP with different ports",
			requestedPort:   types.Int64Value(8080),
			port:            types.Int64Value(443),
			appProtocol:     types.StringValue("MCP"),
			expectedError:   true,
			expectedSummary: "Invalid requested_port for MCP protocol",
			expectedDetail:  "When app_protocol is \"MCP\", requested_port must be equal to port.",
		},
		{
			name:          "valid - non-MCP protocol with different ports",
			requestedPort: types.Int64Value(8080),
			port:          types.Int64Value(443),
			appProtocol:   types.StringValue("HTTP"),
			expectedError: false,
		},
		{
			name:          "valid - non-MCP protocol with matching ports",
			requestedPort: types.Int64Value(443),
			port:          types.Int64Value(443),
			appProtocol:   types.StringValue("HTTPS"),
			expectedError: false,
		},
		{
			name:          "valid - requested_port is null",
			requestedPort: types.Int64Null(),
			port:          types.Int64Value(443),
			appProtocol:   types.StringValue("MCP"),
			expectedError: false,
		},
		{
			name:          "valid - requested_port is unknown",
			requestedPort: types.Int64Unknown(),
			port:          types.Int64Value(443),
			appProtocol:   types.StringValue("MCP"),
			expectedError: false,
		},
		{
			name:          "valid - app_protocol is null",
			requestedPort: types.Int64Value(8080),
			port:          types.Int64Value(443),
			appProtocol:   types.StringNull(),
			expectedError: false,
		},
		{
			name:          "valid - app_protocol is unknown",
			requestedPort: types.Int64Value(8080),
			port:          types.Int64Value(443),
			appProtocol:   types.StringUnknown(),
			expectedError: false,
		},
		{
			name:          "valid - port is null",
			requestedPort: types.Int64Value(443),
			port:          types.Int64Null(),
			appProtocol:   types.StringValue("MCP"),
			expectedError: false,
		},
		{
			name:          "valid - port is unknown",
			requestedPort: types.Int64Value(443),
			port:          types.Int64Unknown(),
			appProtocol:   types.StringValue("MCP"),
			expectedError: false,
		},
		{
			name:          "valid - case sensitive, mcp lowercase",
			requestedPort: types.Int64Value(8080),
			port:          types.Int64Value(443),
			appProtocol:   types.StringValue("mcp"),
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			// Create a schema with the service_endpoint block structure
			schema := tfsdk.Schema{
				Attributes: map[string]tfsdk.Attribute{
					"service_endpoint": {
						Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
							"requested_port": {
								Type:     types.Int64Type,
								Required: true,
							},
							"port": {
								Type:     types.Int64Type,
								Required: true,
							},
							"app_protocol": {
								Type:     types.StringType,
								Required: true,
							},
						}),
						Required: true,
					},
				},
			}

			// Create the config value
			configValue := tftypes.NewValue(
				tftypes.Object{
					AttributeTypes: map[string]tftypes.Type{
						"service_endpoint": tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"requested_port": tftypes.Number,
								"port":           tftypes.Number,
								"app_protocol":   tftypes.String,
							},
						},
					},
				},
				map[string]tftypes.Value{
					"service_endpoint": tftypes.NewValue(
						tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"requested_port": tftypes.Number,
								"port":           tftypes.Number,
								"app_protocol":   tftypes.String,
							},
						},
						map[string]tftypes.Value{
							"requested_port": getTfTypeValue(tt.requestedPort),
							"port":           getTfTypeValue(tt.port),
							"app_protocol":   getTfTypeValue(tt.appProtocol),
						},
					),
				},
			)

			config := tfsdk.Config{
				Schema: schema,
				Raw:    configValue,
			}

			req := validator.Int64Request{
				Path:           path.Root("service_endpoint").AtName("requested_port"),
				PathExpression: path.MatchRoot("service_endpoint").AtName("requested_port"),
				ConfigValue:    tt.requestedPort,
				Config:         config,
			}

			resp := &validator.Int64Response{
				Diagnostics: diag.Diagnostics{},
			}

			v := RequestedPortEqualsPortForMCPValidation()
			v.ValidateInt64(ctx, req, resp)

			if tt.expectedError {
				if !resp.Diagnostics.HasError() {
					t.Fatal("expected error, but got none")
				}

				if len(resp.Diagnostics.Errors()) != 1 {
					t.Fatalf("expected 1 error, got %d", len(resp.Diagnostics.Errors()))
				}

				err := resp.Diagnostics.Errors()[0]
				if err.Summary() != tt.expectedSummary {
					t.Errorf("expected error summary %q, got %q", tt.expectedSummary, err.Summary())
				}
				if err.Detail() != tt.expectedDetail {
					t.Errorf("expected error detail %q, got %q", tt.expectedDetail, err.Detail())
				}
			} else {
				if resp.Diagnostics.HasError() {
					t.Fatalf("unexpected error: %v", resp.Diagnostics.Errors())
				}
			}
		})
	}
}

// Helper function to convert framework types to tftypes.Value
func getTfTypeValue(v attr.Value) tftypes.Value {
	val, err := v.ToTerraformValue(context.Background())
	if err != nil {
		panic(err)
	}
	return val
}
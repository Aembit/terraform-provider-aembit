package validators

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type requestedPortEqualsPortForMCPValidator struct{}

var _ validator.Int64 = requestedPortEqualsPortForMCPValidator{}

func RequestedPortEqualsPortForMCPValidation() validator.Int64 {
	return requestedPortEqualsPortForMCPValidator{}
}

func (v requestedPortEqualsPortForMCPValidator) Description(_ context.Context) string {
	return "Ensures requested_port equals port when app_protocol is MCP"
}

func (v requestedPortEqualsPortForMCPValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v requestedPortEqualsPortForMCPValidator) ValidateInt64(
	ctx context.Context,
	req validator.Int64Request,
	resp *validator.Int64Response,
) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	parentPath := req.Path.ParentPath()

	var appProtocol basetypes.StringValue
	resp.Diagnostics.Append(
		req.Config.GetAttribute(ctx, parentPath.AtName("app_protocol"), &appProtocol)...,
	)
	if resp.Diagnostics.HasError() || appProtocol.IsNull() || appProtocol.IsUnknown() {
		return
	}

	var port basetypes.Int64Value
	resp.Diagnostics.Append(
		req.Config.GetAttribute(ctx, parentPath.AtName("port"), &port)...,
	)
	if resp.Diagnostics.HasError() || port.IsNull() || port.IsUnknown() {
		return
	}

	if appProtocol.ValueString() == "MCP" &&
		req.ConfigValue.ValueInt64() != port.ValueInt64() {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid requested_port for MCP protocol",
			"When app_protocol is \"MCP\", requested_port must be equal to port.",
		)
	}
}
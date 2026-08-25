package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// models.McpToolAccessRuleResourceModel maps the resource schema.
type McpToolAccessRuleResourceModel struct {
	ID                types.String `tfsdk:"id"`
	ContentSecurityID types.String `tfsdk:"content_security_id"`
	ToolName          types.String `tfsdk:"tool_name"`
	IsVisible         types.Bool   `tfsdk:"is_visible"`
	IsInvocable       types.Bool   `tfsdk:"is_invocable"`
}

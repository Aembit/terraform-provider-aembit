package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// models.AgentControllerDeviceCodeDataSourceModel maps the resource schema.
type AgentControllerDeviceCodeDataSourceModel struct {
	ID         types.String `tfsdk:"agent_controller_id"`
	DeviceCode types.String `tfsdk:"device_code"`
}

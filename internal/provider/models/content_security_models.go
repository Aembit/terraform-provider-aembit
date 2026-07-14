package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// models.ContentSecurityResourceModel maps the resource schema.
type ContentSecurityResourceModel struct {
	ID                    types.String                               `tfsdk:"id"`
	ResourceSetID         types.String                               `tfsdk:"resource_set_id"`
	Name                  types.String                               `tfsdk:"name"`
	Description           types.String                               `tfsdk:"description"`
	IsActive              types.Bool                                 `tfsdk:"is_active"`
	Tags                  types.Map                                  `tfsdk:"tags"`
	TagsAll               types.Map                                  `tfsdk:"tags_all"`
	Type                  types.String                               `tfsdk:"type"`
	CrowdStrikeFalconAIDR *CrowdStrikeFalconAIDRContentSecurityModel `tfsdk:"crowdstrike_falcon_aidr"`
}

type CrowdStrikeFalconAIDRContentSecurityModel struct {
	EncryptedToken types.String `tfsdk:"encrypted_token"`
	BaseUrl        types.String `tfsdk:"base_url"`
	FailOpen       types.Bool   `tfsdk:"fail_open"`
	TimeoutMs      types.Int64  `tfsdk:"timeout_ms"`
	MaxRetries     types.Int64  `tfsdk:"max_retries"`
}

// ContentSecuritiesDataSourceModel maps the datasource schema.
type ContentSecuritiesDataSourceModel struct {
	ResourceSetID     types.String                   `tfsdk:"resource_set_id"`
	ContentSecurities []ContentSecurityResourceModel `tfsdk:"content_securities"`
}

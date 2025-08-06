package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CallerIdentityDataSourceModel struct {
	TenantId types.String `tfsdk:"tenant_id"`
}

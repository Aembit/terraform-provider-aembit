package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type timezoneResourceModel struct {
	Timezone types.String `tfsdk:"timezone"`
	Group    types.String `tfsdk:"group"`
	Label    types.String `tfsdk:"label"`
}

// integrationDataSourceModel maps the datasource schema.
type timezonesDataSourceModel struct {
	Timezones []*timezoneResourceModel `tfsdk:"timezones"`
}

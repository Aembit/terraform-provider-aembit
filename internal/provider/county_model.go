package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type countrySubdivisionResourceModel struct {
	CountryCode     types.String `tfsdk:"country_code"`
	Name            types.String `tfsdk:"name"`
	SubdivisionCode types.String `tfsdk:"subdivision_code"`
}

type countryResourceModel struct {
	CountryCode  types.String                       `tfsdk:"country_code"`
	ShortName    types.String                       `tfsdk:"short_name"`
	Subdivisions []*countrySubdivisionResourceModel `tfsdk:"subdivisions"`
}

// integrationDataSourceModel maps the datasource schema.
type countriesDataSourceModel struct {
	Countries []*countryResourceModel `tfsdk:"countries"`
}

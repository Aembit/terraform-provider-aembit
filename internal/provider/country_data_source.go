package provider

import (
	"context"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &countriesDataSource{}
	_ datasource.DataSourceWithConfigure = &countriesDataSource{}
)

// NewcountriesDataSource is a helper function to simplify the provider implementation.
func NewCountriesDataSource() datasource.DataSource {
	return &countriesDataSource{}
}

// countriesDataSource is the data source implementation.
type countriesDataSource struct {
	client *aembit.CloudClient
}

// Configure adds the provider configured client to the data source.
func (d *countriesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = datasourceConfigure(req, resp)
}

// Metadata returns the data source type name.
func (d *countriesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_countries"
}

// Schema defines the schema for the resource.
func (d *countriesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Available countries",
		Attributes: map[string]schema.Attribute{
			"countries": schema.ListNestedAttribute{
				Description: "List of countries.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"alpha2_code": schema.StringAttribute{
							Description: "Country alpha2Code",
							Computed:    true,
						},
						"short_name": schema.StringAttribute{
							Description: "Country shortName",
							Computed:    true,
						},
						"subdivisions": schema.ListNestedAttribute{
							Description: "Country sub-divisions",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"alpha2_code": schema.StringAttribute{
										Description: "alpha2Code of the parent Country",
										Computed:    true,
									},
									"name": schema.StringAttribute{
										Description: "name",
										Computed:    true,
									},
									"subdivision_code": schema.StringAttribute{
										Description: "subdivisionCode",
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *countriesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	state := GetCountries(d.client)
	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func GetCountries(client *aembit.CloudClient) countriesDataSourceModel {
	var state countriesDataSourceModel

	// include staticcountries if a type is filtered
	countries, _ := client.GetCountries(nil)

	// Map response body to model
	for _, country := range countries {
		model := countryResourceModel{}
		model.Alpha2Code = types.StringValue(country.Alpha2Code)
		model.ShortName = types.StringValue(country.ShortName)

		for _, sd := range country.Subdivisions {
			model.Subdivisions = append(model.Subdivisions, &countrySubdivisionResourceModel{
				Name:            types.StringValue(sd.Name),
				Alpha2Code:      types.StringValue(sd.Alpha2Code),
				SubdivisionCode: types.StringValue(sd.SubdivisionCode),
			})
		}

		state.Countries = append(state.Countries, &model)
	}

	return state
}

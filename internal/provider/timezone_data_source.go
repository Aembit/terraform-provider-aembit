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
	_ datasource.DataSource              = &timezonesDataSource{}
	_ datasource.DataSourceWithConfigure = &timezonesDataSource{}
)

// NewtimezonesDataSource is a helper function to simplify the provider implementation.
func NewtimezonesDataSource() datasource.DataSource {
	return &timezonesDataSource{}
}

// timezonesDataSource is the data source implementation.
type timezonesDataSource struct {
	client *aembit.CloudClient
}

// Configure adds the provider configured client to the data source.
func (d *timezonesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = datasourceConfigure(req, resp)
}

// Metadata returns the data source type name.
func (d *timezonesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_timezones"
}

// Schema defines the schema for the resource.
func (d *timezonesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Available timeZones",
		Attributes: map[string]schema.Attribute{
			"timezones": schema.ListNestedAttribute{
				Description: "List of timeZones.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"timezone": schema.StringAttribute{
							Description: "Timezone value such as America/Chicago, Europe/Istanbul",
							Computed:    true,
						},
						"group": schema.StringAttribute{
							Description: "Group",
							Computed:    true,
						},
						"label": schema.StringAttribute{
							Description: "Label",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *timezonesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	state := GetTimezones(d.client)
	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func GetTimezones(client *aembit.CloudClient) timezonesDataSourceModel {
	var state timezonesDataSourceModel

	timeZones, _ := client.GetTimezones(nil)
	// Map response body to model
	for _, tz := range timeZones {
		model := timezoneResourceModel{}
		model.Timezone = types.StringValue(tz.Timezone)
		model.Group = types.StringValue(tz.Group)
		model.Label = types.StringValue(tz.Label)

		state.Timezones = append(state.Timezones, &model)
	}

	return state
}

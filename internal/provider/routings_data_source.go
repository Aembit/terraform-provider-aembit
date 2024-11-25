package provider

import (
	"context"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &routingsDataSource{}
	_ datasource.DataSourceWithConfigure = &routingsDataSource{}
)

// NewroutingsDataSource is a helper function to simplify the provider implementation.
func NewRoutingsDataSource() datasource.DataSource {
	return &routingsDataSource{}
}

// routingsDataSource is the data source implementation.
type routingsDataSource struct {
	client *aembit.CloudClient
}

// Configure adds the provider configured client to the data source.
func (d *routingsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = datasourceConfigure(req, resp)
}

// Metadata returns the data source type name.
func (d *routingsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_routings"
}

// Schema defines the schema for the resource.
func (d *routingsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a routing.",
		Attributes: map[string]schema.Attribute{
			"routings": schema.ListNestedAttribute{
				Description: "List of routings.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						// ID field is required for Terraform Framework acceptance testing.
						"id": schema.StringAttribute{
							Description: "Unique identifier of the routing.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "User-provided name of the routing.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "User-provided description of the routing.",
							Computed:    true,
						},
						"is_active": schema.BoolAttribute{
							Description: "Active/Inactive status of the routing.",
							Computed:    true,
						},
						"proxy_url": schema.StringAttribute{
							Description: "User-provided proxyUrl of the routing.",
							Computed:    true,
						},
						"resource_set_name": schema.StringAttribute{
							Description: "User-provided resourceSetName of the routing.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *routingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state routingsDataSourceModel

	routings, err := d.client.GetRoutings(nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Aembit routings",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, routing := range routings {
		routingState := convertListRoutingDTOToModel(routing, routingListResourceModel{})
		state.Routings = append(state.Routings, routingState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

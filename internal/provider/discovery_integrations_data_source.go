package provider

import (
	"context"
	"terraform-provider-aembit/internal/provider/models"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &discoveryIntegrationsDataSource{}
	_ datasource.DataSourceWithConfigure = &discoveryIntegrationsDataSource{}
)

// NewIntegrationsDataSource is a helper function to simplify the provider implementation.
func NewDiscoveryIntegrationsDataSource() datasource.DataSource {
	return &discoveryIntegrationsDataSource{}
}

// integrationsDataSource is the data source implementation.
type discoveryIntegrationsDataSource struct {
	client *aembit.CloudClient
}

// Configure adds the provider configured client to the data source.
func (d *discoveryIntegrationsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = datasourceConfigure(req, resp)
}

// Metadata returns the data source type name.
func (d *discoveryIntegrationsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_discovery_integrations"
}

// Schema defines the schema for the resource.
func (d *discoveryIntegrationsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a discovery integration.",
		Attributes: map[string]schema.Attribute{
			"discovery_integrations": schema.ListNestedAttribute{
				Description: "List of discovery integrations.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						// ID field is required for Terraform Framework acceptance testing.
						"id": schema.StringAttribute{
							Description: "Unique identifier of the discovery integration.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "User-defined name of the discovery integration.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "User-defined description of the discovery integration.",
							Computed:    true,
						},
						"is_active": schema.BoolAttribute{
							Description: "Active status of the discovery integration.",
							Optional:    true,
							Computed:    true,
						},
						"tags": schema.MapAttribute{
							Description: "Tags are key-value pairs.",
							ElementType: types.StringType,
							Optional:    true,
							Computed:    true,
						},
						"type": schema.StringAttribute{
							Description: "Type of discovery integration. Possible value is: `WizIntegrationApi`.",
							Required:    true,
						},
						"sync_frequency_seconds": schema.Int64Attribute{
							Description: "Frequency (in seconds) for synchronizing the discovery integration. Accepted range: 300-3600 seconds",
							Computed:    true,
						},
						"last_sync": schema.StringAttribute{
							Description: "ISO 8601-formatted last sync date of the discovery integration.",
							Computed:    true,
						},
						"last_sync_status": schema.StringAttribute{
							Description: "Status of the last sync of the discovery integration.",
							Computed:    true,
						},
						"endpoint": schema.StringAttribute{
							Description: "Endpoint that performs the discovery integration.",
							Required:    true,
						},
						"wiz_integration": schema.SingleNestedAttribute{
							Description: "Wiz-specific properties for the discovery integration.",
							Optional:    true,
							Attributes: map[string]schema.Attribute{
								"token_url": schema.StringAttribute{
									Description: "Token URL for the Wiz Endpoint of the discovery integration.",
									Required:    true,
								},
								"client_id": schema.StringAttribute{
									Description: "Client ID for the Wiz Endpoint of the discovery integration.",
									Required:    true,
								},
								"client_secret": schema.StringAttribute{
									Description: "Client Secret for the Wiz Endpoint of the discovery integration.",
									Required:    true,
									Sensitive:   true,
								},
								"audience": schema.StringAttribute{
									Description: "Audience for the Wiz Endpoint of the discovery integration.",
									Required:    true,
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
func (d *discoveryIntegrationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state models.DiscoveryIntegrationsDataSourceModel

	discoveryIntegrations, err := d.client.GetDiscoveryIntegrations(nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Discovery Integrations",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, discoveryIntegration := range discoveryIntegrations {
		discoveryIntegtationState := convertDiscoveryIntegrationDTOToModel(ctx, discoveryIntegration, models.DiscoveryIntegrationResourceModel{})
		state.DiscoveryIntegrations = append(state.DiscoveryIntegrations, discoveryIntegtationState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

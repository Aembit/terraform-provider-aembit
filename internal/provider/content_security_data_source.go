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
	_ datasource.DataSource              = &contentSecuritiesDataSource{}
	_ datasource.DataSourceWithConfigure = &contentSecuritiesDataSource{}
)

// NewContentSecuritiesDataSource is a helper function to simplify the provider implementation.
func NewContentSecuritiesDataSource() datasource.DataSource {
	return &contentSecuritiesDataSource{}
}

// contentSecuritiesDataSource is the data source implementation.
type contentSecuritiesDataSource struct {
	client *aembit.CloudClient
}

// Configure adds the provider configured client to the data source.
func (d *contentSecuritiesDataSource) Configure(
	_ context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
	d.client = datasourceConfigure(req, resp)
}

// Metadata returns the data source type name.
func (d *contentSecuritiesDataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_content_securities"
}

// Schema defines the schema for the data source.
func (d *contentSecuritiesDataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Queries the list of Content Security resources.",
		Attributes: map[string]schema.Attribute{
			"resource_set_id": schema.StringAttribute{
				Description: "ResourceSet unique identifier to filter the Content Security resources, or provider-level override.",
				Optional:    true,
				Computed:    true,
			},
			"content_securities": schema.ListNestedAttribute{
				Description: "List of Content Security resources.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Unique identifier of the Content Security resource.",
							Computed:    true,
						},
						"resource_set_id": schema.StringAttribute{
							Description: "ResourceSet unique identifier of the Content Security resource.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name for the Content Security resource.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "Description for the Content Security resource.",
							Computed:    true,
						},
						"is_active": schema.BoolAttribute{
							Description: "Status of the Content Security resource (`active` or `inactive`).",
							Computed:    true,
						},
						"tags":     TagsComputedMapAttribute(),
						"tags_all": TagsAllMapAttribute(),
						"type": schema.StringAttribute{
							Description: "Type of the Content Security resource.",
							Computed:    true,
						},
						"crowdstrike_falcon_aidr": schema.SingleNestedAttribute{
							Description: "CrowdStrike Falcon AIDR configuration settings.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"encrypted_token": schema.StringAttribute{
									Description: "The encrypted API token or client secret used to authenticate with the CrowdStrike Falcon AIDR service.",
									Computed:    true,
									Sensitive:   true,
								},
								"base_url": schema.StringAttribute{
									Description: "The base URL of the CrowdStrike Falcon AIDR service endpoint or API gateway.",
									Computed:    true,
								},
								"fail_open": schema.BoolAttribute{
									Description: "Indicates whether requests should be allowed (fail-open) or blocked (fail-closed) if the CrowdStrike Falcon AIDR service is unreachable.",
									Computed:    true,
								},
								"timeout_ms": schema.Int64Attribute{
									Description: "The connection timeout in milliseconds for requests sent to the CrowdStrike Falcon AIDR service.",
									Computed:    true,
								},
								"max_retries": schema.Int64Attribute{
									Description: "The maximum number of retry attempts for requests to the CrowdStrike Falcon AIDR service.",
									Computed:    true,
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
func (d *contentSecuritiesDataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var state models.ContentSecuritiesDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceSetId := getResourceSetId(state.ResourceSetID, d.client)

	contentSecurities, err := d.client.GetContentSecurities(nil, &resourceSetId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Aembit Content Securities",
			err.Error(),
		)
		return
	}

	state.ResourceSetID = types.StringValue(resourceSetId)

	// Map response body to model
	for _, contentSecurity := range contentSecurities {
		contentSecurityState := convertContentSecurityDTOToModel(
			ctx,
			contentSecurity,
			&models.ContentSecurityResourceModel{},
		)
		contentSecurityState.Tags = newTagsModel(ctx, contentSecurity.Tags)
		state.ContentSecurities = append(state.ContentSecurities, contentSecurityState)
	}

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

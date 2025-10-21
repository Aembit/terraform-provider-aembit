package provider

import (
	"context"

	"terraform-provider-aembit/internal/provider/models"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &standaloneCertificateAuthoritiesDataSource{}
	_ datasource.DataSourceWithConfigure = &standaloneCertificateAuthoritiesDataSource{}
)

// NewIntegrationsDataSource is a helper function to simplify the provider implementation.
func NewStandaloneCertificateAuthoritiesDataSource() datasource.DataSource {
	return &standaloneCertificateAuthoritiesDataSource{}
}

// integrationsDataSource is the data source implementation.
type standaloneCertificateAuthoritiesDataSource struct {
	client *aembit.CloudClient
}

// Configure adds the provider configured client to the data source.
func (d *standaloneCertificateAuthoritiesDataSource) Configure(
	_ context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
	d.client = datasourceConfigure(req, resp)
}

// Metadata returns the data source type name.
func (d *standaloneCertificateAuthoritiesDataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_standalone_certificate_authorities"
}

// Schema defines the schema for the resource.
func (d *standaloneCertificateAuthoritiesDataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Manages a standalone certificate authority.",
		Attributes: map[string]schema.Attribute{
			"standalone_certificate_authorities": schema.ListNestedAttribute{
				Description: "List of standalone certificate authorities.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						// ID field is required for Terraform Framework acceptance testing.
						"id": schema.StringAttribute{
							Description: "Unique identifier of the standalone certificate authority.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "User-provided name of the standalone certificate authority.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "User-provided description of the standalone certificate authority.",
							Computed:    true,
						},
						"tags": TagsComputedMapAttribute(),
						"leaf_lifetime": schema.Int32Attribute{
							Description: "Leaf certificate lifetime(in minutes) of the standalone certificate authority. Valid options; 60, 1440, 10080.",
							Computed:    true,
						},
						"not_before": schema.StringAttribute{
							Description: "ISO 8601 formatted not before date of the standalone certificate authority.",
							Computed:    true,
						},
						"not_after": schema.StringAttribute{
							Description: "ISO 8601 formatted not after date of the standalone certificate authority.",
							Computed:    true,
						},
						"client_workload_count": schema.Int32Attribute{
							Description: "Client workloads associated with the standalone certificate authority.",
							Optional:    true,
							Computed:    true,
						},
						"resource_set_count": schema.Int32Attribute{
							Description: "Resource sets associated with the standalone certificate authority.",
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *standaloneCertificateAuthoritiesDataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var state models.StandaloneCertificateAuthoritiesDataSourceModel

	standaloneCertificates, err := d.client.GetStandaloneCertificates(nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Standalone Certificate Authorities",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, standaloneCertificate := range standaloneCertificates {
		standaloneCertificateState := convertStandaloneCertificateDTOToModel(
			ctx,
			standaloneCertificate,
		)
		state.StandaloneCertificates = append(
			state.StandaloneCertificates,
			standaloneCertificateState,
		)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

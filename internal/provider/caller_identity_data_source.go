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
	_ datasource.DataSource              = &callerIdentityDataSource{}
	_ datasource.DataSourceWithConfigure = &callerIdentityDataSource{}
)

// NewCallerIdentityDataSource returns a new instance of the caller identity data source.
func NewCallerIdentityDataSource() datasource.DataSource {
	return &callerIdentityDataSource{}
}

// callerIdentityDataSource implements the data source logic.
type callerIdentityDataSource struct {
	client *aembit.CloudClient
}

// Configure adds the provider client to the data source.
func (d *callerIdentityDataSource) Configure(
	_ context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
	d.client = datasourceConfigure(req, resp)
}

// Metadata returns the data source type name.
func (d *callerIdentityDataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_caller_identity"
}

// Schema defines the schema for the data source.
func (d *callerIdentityDataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Retrieves the Tenant ID associated with the current session.",
		Attributes: map[string]schema.Attribute{
			"tenant_id": schema.StringAttribute{
				Computed:    true,
				Description: "The Tenant ID associated with the current session.",
			},
		},
	}
}

// Read retrieves the tenant ID from the client and sets it in state.
func (d *callerIdentityDataSource) Read(
	ctx context.Context,
	_ datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	if d.client == nil {
		resp.Diagnostics.AddError(
			"Missing client",
			"Client was not configured for the data source.",
		)
		return
	}

	state := models.CallerIdentityDataSourceModel{
		TenantId: types.StringValue(d.client.Tenant),
	}

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

package provider

import (
	"context"
	"fmt"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &credentialProvidersDataSource{}
	_ datasource.DataSourceWithConfigure = &credentialProvidersDataSource{}
)

// NewCredentialProvidersDataSource is a helper function to simplify the provider implementation.
func NewCredentialProvidersDataSource() datasource.DataSource {
	return &credentialProvidersDataSource{}
}

// credentialProvidersDataSource is the data source implementation.
type credentialProvidersDataSource struct {
	client *aembit.Client
}

// Configure adds the provider configured client to the data source.
func (d *credentialProvidersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*aembit.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *aembit.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

// Metadata returns the data source type name.
func (d *credentialProvidersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_credential_providers"
}

// Schema defines the schema for the resource.
func (r *credentialProvidersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an credential provider.",
		Attributes: map[string]schema.Attribute{
			"credential_providers": schema.ListNestedAttribute{
				Description: "List of credential providers.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						// ID field is required for Terraform Framework acceptance testing.
						"id": schema.StringAttribute{
							Description: "Alphanumeric identifier of the credential provider.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "User-provided name of the credential provider.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "User-provided description of the credential provider.",
							Computed:    true,
						},
						"is_active": schema.BoolAttribute{
							Description: "Active/Inactive status of the credential provider.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *credentialProvidersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state credentialProvidersDataSourceModel

	credential_providers, err := d.client.GetCredentialProviders(nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Aembit Credential Providers",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, credential_provider := range credential_providers {
		credentialProviderState := credentialProviderResourceModel{
			ID:          types.StringValue(credential_provider.EntityDTO.ExternalId),
			Name:        types.StringValue(credential_provider.EntityDTO.Name),
			Description: types.StringValue(credential_provider.EntityDTO.Description),
			IsActive:    types.BoolValue(credential_provider.EntityDTO.IsActive),
		}
		state.CredentialProviders = append(state.CredentialProviders, credentialProviderState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

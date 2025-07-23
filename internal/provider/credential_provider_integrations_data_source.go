package provider

import (
	"context"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"terraform-provider-aembit/internal/provider/models"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &credentialProviderIntegrationsDataSource{}
	_ datasource.DataSourceWithConfigure = &credentialProviderIntegrationsDataSource{}
)

// NewCredentialProviderIntegrationsDataSource is a helper function to simplify the provider implementation.
func NewCredentialProviderIntegrationsDataSource() datasource.DataSource {
	return &credentialProviderIntegrationsDataSource{}
}

// credentialProviderIntegrationsDataSource is the data source implementation.
type credentialProviderIntegrationsDataSource struct {
	client *aembit.CloudClient
}

// Configure adds the provider configured client to the data source.
func (d *credentialProviderIntegrationsDataSource) Configure(
	_ context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
	d.client = datasourceConfigure(req, resp)
}

// Metadata returns the data source type name.
func (d *credentialProviderIntegrationsDataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_credential_provider_integrations"
}

// Schema defines the schema for the resource.
func (d *credentialProviderIntegrationsDataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Manages a credential provider integration.",
		Attributes: map[string]schema.Attribute{
			"credential_provider_integrations": schema.ListNestedAttribute{
				Description: "List of credential provider integrations.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						// ID field is required for Terraform Framework acceptance testing.
						"id": schema.StringAttribute{
							Description: "Unique identifier of the credential provider integration.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "User-provided name of the credential provider integration.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "User-provided description of the credential provider integration.",
							Computed:    true,
						},
						"gitlab": schema.SingleNestedAttribute{
							Description: "GitLab Managed Account type Credential Provider Integration configuration.",
							Optional:    true,
							Attributes: map[string]schema.Attribute{
								"url": schema.StringAttribute{
									Description: "GitLab URL.",
									Computed:    true,
								},
								"personal_access_token": schema.StringAttribute{
									Description: "GitLab personal access token value.",
									Computed:    true,
									Sensitive:   true,
								},
							},
						},
						"aws_iam_role": schema.SingleNestedAttribute{
							Description: "AWS IAM Role type Credential Provider Integration configuration.",
							Optional:    true,
							Attributes: map[string]schema.Attribute{
								"role_arn": schema.StringAttribute{
									Description: "AWS IAM Role ARN.",
									Computed:    true,
								},
								"lifetime_in_seconds": schema.Int32Attribute{
									Description: "Lifetime in seconds for the requested AWS temporary credentials.",
									Computed:    true,
								},
								"oidc_issuer_url": schema.StringAttribute{
									Description: "OIDC Issuer URL for AWS IAM Identity Provider configuration",
									Computed:    true,
								},
								"token_audience": schema.StringAttribute{
									Description: "Token Audience for AWS IAM Identity Provider configuration",
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
func (d *credentialProviderIntegrationsDataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var state models.CredentialProviderIntegrationsDataSourceModel

	req.Config.Get(ctx, &state)

	integrations, err := d.client.GetCredentialProviderIntegrations(nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Aembit Credential Provider Integrations",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, integration := range integrations {
		integrationState := convertCredentialProviderIntegrationDTOToModel(
			integration,
			models.CredentialProviderIntegrationResourceModel{},
			d.client.Tenant,
			d.client.StackDomain,
		)

		state.CredentialProviderIntegrations = append(
			state.CredentialProviderIntegrations,
			integrationState,
		)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

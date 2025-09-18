package provider

import (
	"context"
	"fmt"

	"terraform-provider-aembit/internal/provider/models"
	"terraform-provider-aembit/internal/provider/validators"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int32default"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &credentialProviderIntegrationResource{}
	_ resource.ResourceWithConfigure   = &credentialProviderIntegrationResource{}
	_ resource.ResourceWithImportState = &credentialProviderIntegrationResource{}
)

// NewIntegrationResource is a helper function to simplify the provider implementation.
func NewCredentialProviderIntegrationResource() resource.Resource {
	return &credentialProviderIntegrationResource{}
}

// credentialProviderIntegrationResource is the resource implementation.
type credentialProviderIntegrationResource struct {
	client *aembit.CloudClient
}

// Metadata returns the resource type name.
func (r *credentialProviderIntegrationResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_credential_provider_integration"
}

// Configure adds the provider configured client to the resource.
func (r *credentialProviderIntegrationResource) Configure(
	_ context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	r.client = resourceConfigure(req, resp)
}

// Schema defines the schema for the resource.
func (r *credentialProviderIntegrationResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			// ID field is required for Terraform Framework acceptance testing.
			"id": schema.StringAttribute{
				Description: "Unique identifier of the Credential Provider Integration.",
				Computed:    true,
				Validators: []validator.String{
					validators.UUIDRegexValidation(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name for the Credential Provider Integration.",
				Required:    true,
				Validators: []validator.String{
					validators.NameLengthValidation(),
				},
			},
			"description": schema.StringAttribute{
				Description: "Description for the Credential Provider Integration.",
				Optional:    true,
				Computed:    true,
			},
			"gitlab": schema.SingleNestedAttribute{
				Description: "GitLab Managed Account type Credential Provider Integration configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"url": schema.StringAttribute{
						Description: "GitLab URL.",
						Required:    true,
						Validators: []validator.String{
							validators.SecureURLValidation(),
						},
					},
					"personal_access_token": schema.StringAttribute{
						Description: "GitLab personal access token value.",
						Required:    true,
						Sensitive:   true,
					},
				},
			},
			"aws_iam_role": schema.SingleNestedAttribute{
				Description: "Configuration of Credential Provider Integration of type AWS IAM Role",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"role_arn": schema.StringAttribute{
						Description: "AWS IAM Role ARN.",
						Required:    true,
						Validators: []validator.String{
							validators.AwsIamRoleArnValidation(),
						},
					},
					"lifetime_in_seconds": schema.Int32Attribute{
						Description: "Lifetime in seconds for the AWS IAM Role credentials.",
						Optional:    true,
						Computed:    true,
						Default:     int32default.StaticInt32(3600),
						Validators: []validator.Int32{
							int32validator.Between(900, 43200), // 15 minutes to 12 hours
						},
					},
					"fetch_secret_arns": schema.BoolAttribute{
						Description: "Whether to fetch secret ARNs from AWS Secret Manager in the Aembit UI.",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
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
			"azure_entra_federation": schema.SingleNestedAttribute{
				Description: "Configuration of Credential Provider Integration of type Azure Entra Federation",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"audience": schema.StringAttribute{
						Description: "The audience claim (aud) for the federated JWT token.",
						Required:    true,
					},
					"subject": schema.StringAttribute{
						Description: "The subject claim (sub) for the federated JWT token.",
						Required:    true,
					},
					"azure_tenant": schema.StringAttribute{
						Description: "The Azure AD tenant ID that the application belongs to.",
						Required:    true,
						Validators: []validator.String{
							validators.UUIDRegexValidation(),
						},
					},
					"client_id": schema.StringAttribute{
						Description: "The client ID of the Azure AD application.",
						Required:    true,
						Validators: []validator.String{
							validators.UUIDRegexValidation(),
						},
					},
					"key_vault_name": schema.StringAttribute{
						Description: "The name of the Azure Key Vault.",
						Required:    true,
					},
					"fetch_secret_names": schema.BoolAttribute{
						Description: "Whether to fetch secret names from Azure Key Vault in the Aembit UI.",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
					},
					"oidc_issuer_url": schema.StringAttribute{
						Description: "OIDC Issuer URL for Azure Entra Federation configuration",
						Computed:    true,
					},
				},
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *credentialProviderIntegrationResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	// Retrieve values from plan
	var plan models.CredentialProviderIntegrationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	dto := convertCredentialProviderIntegrationModelToDTO(plan, nil)

	// Create new Integration
	credentialIntegration, err := r.client.CreateCredentialProviderIntegration(dto, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Integration",
			"Could not create Credential Provider Integration, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan = convertCredentialProviderIntegrationDTOToModel(
		*credentialIntegration,
		plan,
		r.client.Tenant,
		r.client.StackDomain,
	)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *credentialProviderIntegrationResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	// Get current state
	var state models.CredentialProviderIntegrationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed trust value from Aembit
	credentialIntegration, err, notFound := r.client.GetCredentialProviderIntegration(
		state.ID.ValueString(),
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Error reading Aembit Credential Provider Integration",
			"Could not read Aembit External ID from Terraform state "+state.ID.ValueString()+": "+err.Error(),
		)

		// If the resource is not found on Aembit Cloud, delete it locally
		if notFound {
			resp.State.RemoveResource(ctx)
		}
		return
	}

	state = convertCredentialProviderIntegrationDTOToModel(
		credentialIntegration,
		state,
		r.client.Tenant,
		r.client.StackDomain,
	)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *credentialProviderIntegrationResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	// Get current state
	var state models.CredentialProviderIntegrationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract external ID from state
	externalID := state.ID.ValueString()

	// Retrieve values from plan
	var plan models.CredentialProviderIntegrationResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	dto := convertCredentialProviderIntegrationModelToDTO(plan, &externalID)

	// Update Credential Provider Integration
	credentialIntegration, err := r.client.UpdateCredentialProviderIntegration(dto, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Credential Provider Integration",
			"Could not update Credential Provider Integration, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	state = convertCredentialProviderIntegrationDTOToModel(
		*credentialIntegration,
		plan,
		r.client.Tenant,
		r.client.StackDomain,
	)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *credentialProviderIntegrationResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	// Retrieve values from state
	var state models.CredentialProviderIntegrationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing Credential Provider Integration
	_, err := r.client.DeleteCredentialProviderIntegration(ctx, state.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Credential Provider Integration",
			"Could not delete Credential Provider Integration, unexpected error: "+err.Error(),
		)
		return
	}
}

// Imports an existing resource by passing externalId.
func (r *credentialProviderIntegrationResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	// Retrieve import externalId and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func convertCredentialProviderIntegrationModelToDTO(
	model models.CredentialProviderIntegrationResourceModel,
	externalID *string,
) aembit.CredentialProviderIntegrationDTO {
	var credentialIntegration aembit.CredentialProviderIntegrationDTO
	credentialIntegration.EntityDTO = aembit.EntityDTO{
		Name:        model.Name.ValueString(),
		Description: model.Description.ValueString(),
	}

	if externalID != nil {
		credentialIntegration.ExternalID = *externalID
	}

	if model.GitLab != nil {
		credentialIntegration.Type = "GitLab"
		credentialIntegration.Url = model.GitLab.Url.ValueString()
		credentialIntegration.PersonalAccessToken = model.GitLab.PersonalAccessToken.ValueString()
	}

	if model.AwsIamRole != nil {
		credentialIntegration.Type = "AwsIamRole"
		credentialIntegration.RoleArn = model.AwsIamRole.RoleArn.ValueString()
		credentialIntegration.LifetimeInSeconds = model.AwsIamRole.LifetimeInSeconds.ValueInt32()
		credentialIntegration.FetchSecretArns = model.AwsIamRole.FetchSecretArns.ValueBool()
	}

	if model.AzureEntraFederation != nil {
		credentialIntegration.Type = "AzureEntraFederation"
		credentialIntegration.Audience = model.AzureEntraFederation.Audience.ValueString()
		credentialIntegration.Subject = model.AzureEntraFederation.Subject.ValueString()
		credentialIntegration.AzureTenant = model.AzureEntraFederation.AzureTenant.ValueString()
		credentialIntegration.ClientID = model.AzureEntraFederation.ClientID.ValueString()
		credentialIntegration.KeyVaultName = model.AzureEntraFederation.KeyVaultName.ValueString()
		credentialIntegration.FetchSecretNames = model.AzureEntraFederation.FetchSecretNames.ValueBool()
	}
	return credentialIntegration
}

func convertCredentialProviderIntegrationDTOToModel(
	dto aembit.CredentialProviderIntegrationDTO,
	state models.CredentialProviderIntegrationResourceModel,
	tenant string,
	stackDomain string,
) models.CredentialProviderIntegrationResourceModel {
	var model models.CredentialProviderIntegrationResourceModel
	model.ID = types.StringValue(dto.ExternalID)
	model.Name = types.StringValue(dto.Name)
	model.Description = types.StringValue(dto.Description)

	switch dto.Type {
	case "GitLab":
		model.GitLab = &models.CredentialProviderIntegrationGitlabModel{
			Url:                 types.StringValue(dto.Url),
			PersonalAccessToken: types.StringValue(dto.PersonalAccessToken),
		}
		if len(dto.PersonalAccessToken) == 0 && state.GitLab != nil {
			model.GitLab.PersonalAccessToken = state.GitLab.PersonalAccessToken
		}
	case "AwsIamRole":
		model.AwsIamRole = &models.CredentialProviderIntegrationAwsIamRoleModel{
			RoleArn:           types.StringValue(dto.RoleArn),
			LifetimeInSeconds: types.Int32Value(dto.LifetimeInSeconds),
			FetchSecretArns:   types.BoolValue(dto.FetchSecretArns),
			OIDCIssuerUrl: types.StringValue(
				fmt.Sprintf(oidcIssuerTemplate, tenant, stackDomain),
			),
			TokenAudience: types.StringValue("sts.amazonaws.com"),
		}
	case "AzureEntraFederation":
		model.AzureEntraFederation = &models.CredentialProviderIntegrationAzureEntraFederationModel{
			Audience:         types.StringValue(dto.Audience),
			Subject:          types.StringValue(dto.Subject),
			AzureTenant:      types.StringValue(dto.AzureTenant),
			ClientID:         types.StringValue(dto.ClientID),
			KeyVaultName:     types.StringValue(dto.KeyVaultName),
			FetchSecretNames: types.BoolValue(dto.FetchSecretNames),
			OIDCIssuerUrl: types.StringValue(
				fmt.Sprintf(oidcIssuerTemplate, tenant, stackDomain),
			),
		}
	default:
		// This should never happen as the API restricts the type field to known values
	}
	return model
}

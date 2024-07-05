package provider

import (
	"context"
	"fmt"
	"strings"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &credentialProviderV2Resource{}
	_ resource.ResourceWithConfigure   = &credentialProviderV2Resource{}
	_ resource.ResourceWithImportState = &credentialProviderV2Resource{}
)

// NewCredentialProviderV2Resource is a helper function to simplify the provider implementation.
func NewCredentialProviderV2Resource() resource.Resource {
	return &credentialProviderV2Resource{}
}

// credentialProviderV2Resource is the resource implementation.
type credentialProviderV2Resource struct {
	client *aembit.CloudClient
}

// Metadata returns the resource type name.
func (r *credentialProviderV2Resource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_credential_provider_v2"
}

// Configure adds the provider configured client to the resource.
func (r *credentialProviderV2Resource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = resourceConfigure(req, resp)
}

// Schema defines the schema for the resource.
func (r *credentialProviderV2Resource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			// ID field is required for Terraform Framework acceptance testing.
			"id": schema.StringAttribute{
				Description: "Unique identifier of the Credential Provider.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name for the Credential Provider.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Description: "Description for the Credential Provider.",
				Optional:    true,
				Computed:    true,
			},
			"is_active": schema.BoolAttribute{
				Description: "Active status of the Credential Provider.",
				Optional:    true,
				Computed:    true,
			},
			"tags": schema.MapAttribute{
				Description: "Tags are key-value pairs.",
				ElementType: types.StringType,
				Optional:    true,
			},
			"aembit_access_token": schema.SingleNestedAttribute{
				Description: "Aembit Access Token type Credential Provider configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"audience": schema.StringAttribute{
						Description: "Audience of the Credential Provider.",
						Computed:    true,
					},
					"role_id": schema.StringAttribute{
						Description: "Aembit Role ID of the Credential Provider.",
						Required:    true,
					},
					"lifetime": schema.Int64Attribute{
						Description: "Lifetime of the Credential Provider.",
						Required:    true,
					},
				},
			},
			"api_key": schema.SingleNestedAttribute{
				Description: "API Key type Credential Provider configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"api_key": schema.StringAttribute{
						Description: "API Key secret of the Credential Provider.",
						Required:    true,
						Sensitive:   true,
					},
				},
			},
			"aws_sts": schema.SingleNestedAttribute{
				Description: "AWS Security Token Service Federation type Credential Provider configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"oidc_issuer": schema.StringAttribute{
						Description: "OIDC Issuer for AWS IAM Identity Provider configuration of the Credential Provider.",
						Computed:    true,
					},
					"role_arn": schema.StringAttribute{
						Description: "AWS Role Arn to be used for the AWS Session credentials requested by the Credential Provider.",
						Required:    true,
					},
					"token_audience": schema.StringAttribute{
						Description: "Token Audience for AWS IAM Identity Provider configuration of the Credential Provider.",
						Computed:    true,
					},
					"lifetime": schema.Int64Attribute{
						Description: "Lifetime (seconds) of the AWS Session credentials requested by the Credential Provider.",
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(3600),
					},
				},
			},
			"google_workload_identity": schema.SingleNestedAttribute{
				Description: "Google Workload Identity Federation type Credential Provider configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"oidc_issuer": schema.StringAttribute{
						Description: "OIDC Issuer for AWS IAM Identity Provider configuration of the Credential Provider.",
						Computed:    true,
					},
					"audience": schema.StringAttribute{
						Description: "Audience for GCP Workload Identity Federation configuration of the Credential Provider.",
						Required:    true,
					},
					"service_account": schema.StringAttribute{
						Description: "Service Account email of the GCP Session credentials requested by the Credential Provider.",
						Required:    true,
					},
					"lifetime": schema.Int64Attribute{
						Description: "Lifetime (seconds) of the GCP Session credentials requested by the Credential Provider.",
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(3600),
					},
				},
			},
			"snowflake_jwt": schema.SingleNestedAttribute{
				Description: "JSON Web Token type Credential Provider configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"account_id": schema.StringAttribute{
						Description: "Snowflake Account ID of the Credential Provider.",
						Required:    true,
					},
					"username": schema.StringAttribute{
						Description: "Snowflake Username of the Credential Provider.",
						Required:    true,
					},
					"alter_user_command": schema.StringAttribute{
						Description: "Snowflake Alter User Command generated for configuration of Snowflake by the Credential Provider.",
						Computed:    true,
					},
				},
			},
			"oauth_client_credentials": schema.SingleNestedAttribute{
				Description: "OAuth Client Credentials Flow type Credential Provider configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"token_url": schema.StringAttribute{
						Description: "Token URL for the OAuth Credential Provider.",
						Required:    true,
					},
					"client_id": schema.StringAttribute{
						Description: "Client ID for the OAuth Credential Provider.",
						Required:    true,
					},
					"client_secret": schema.StringAttribute{
						Description: "Client Secret for the OAuth Credential Provider.",
						Required:    true,
						Sensitive:   true,
					},
					"scopes": schema.StringAttribute{
						Description: "Scopes for the OAuth Credential Provider.",
						Optional:    true,
					},
					"credential_style": schema.StringAttribute{
						Description: "Credential Style for the OAuth Credential Provider.",
						Required:    true,
					},
					"custom_parameters": schema.SetNestedAttribute{
						Description: "Set Custom Parameters for the OAuth Credential Provider.",
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"key": schema.StringAttribute{
									Description: "Key for Custom Parameter of the OAuth Credential Provider.",
									Required:    true,
								},
								"value": schema.StringAttribute{
									Description: "Value for Custom Parameter of the OAuth Credential Provider.",
									Required:    true,
								},
								"value_type": schema.StringAttribute{
									Description: "Type of value for Custom Parameter of the OAuth Credential Provider. Possible values are `literal` or `dynamic`.",
									Optional:    true,
									Computed:    true,
									Default:     stringdefault.StaticString("literal"),
									Validators: []validator.String{
										stringvalidator.OneOf([]string{
											"literal",
											"dynamic",
										}...),
									},
								},
							},
						},
					},
				},
			},
			"username_password": schema.SingleNestedAttribute{
				Description: "Username/Password type Credential Provider configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"username": schema.StringAttribute{
						Description: "Username of the Credential Provider.",
						Required:    true,
					},
					"password": schema.StringAttribute{
						Description: "Password of the Credential Provider.",
						Required:    true,
						Sensitive:   true,
					},
				},
			},
			"vault_client_token": schema.SingleNestedAttribute{
				Description: "Vault Client Token type Credential Provider configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"subject": schema.StringAttribute{
						Description: "Subject of the JWT Token used to authenticate to the Vault Cluster.",
						Required:    true,
					},
					"subject_type": schema.StringAttribute{
						Description: "Type of value for the JWT Token Subject. Possible values are `literal` or `dynamic`.",
						Required:    true,
						Validators: []validator.String{
							stringvalidator.OneOf([]string{
								"literal",
								"dynamic",
							}...),
						},
					},
					"custom_claims": schema.SetNestedAttribute{
						Description: "Set of Custom Claims for the JWT Token used to authenticate to the Vault Cluster.",
						Optional:    true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"key": schema.StringAttribute{
									Description: "Key for the JWT Token Custom Claim.",
									Required:    true,
								},
								"value": schema.StringAttribute{
									Description: "Value for the JWT Token Custom Claim.",
									Required:    true,
								},
								"value_type": schema.StringAttribute{
									Description: "Type of value for the JWT Token Custom Claim. Possible values are `literal` or `dynamic`.",
									Required:    true,
									Validators: []validator.String{
										stringvalidator.OneOf([]string{
											"literal",
											"dynamic",
										}...),
									},
								},
							},
						},
					},
					"lifetime": schema.Int64Attribute{
						Description: "Lifetime of the JWT Token used to authenticate to the Vault Cluster. Note: The lifetime of the retrieved Vault Client Token is managed within Vault configuration.",
						Required:    true,
					},
					"vault_host": schema.StringAttribute{
						Description: "Hostname of the Vault Cluster to be used for executing the login API.",
						Required:    true,
					},
					"vault_port": schema.Int64Attribute{
						Description: "Port of the Vault Cluster to be used for executing the login API.",
						Required:    true,
					},
					"vault_tls": schema.BoolAttribute{
						Description: "Configuration to utilize TLS for connectivity to the Vault Cluster.",
						Required:    true,
					},
					"vault_namespace": schema.StringAttribute{
						Description: "Namespace to utilize when executing the login API on the Vault Cluster.",
						Optional:    true,
					},
					"vault_role": schema.StringAttribute{
						Description: "Role to utilize when executing the login API on the Vault Cluster.",
						Optional:    true,
					},
					"vault_path": schema.StringAttribute{
						Description: "Path to utilize when executing the login API on the Vault Cluster.",
						Required:    true,
					},
					"vault_forwarding": schema.StringAttribute{
						Description: "If Vault Forwarding is required, this configuration can be set to `unconditional` or `conditional`.",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
						Validators: []validator.String{
							stringvalidator.OneOf([]string{
								"",
								"unconditional",
								"conditional",
							}...),
						},
					},
				},
			},
		},
	}
}

// Configure validators to ensure that only one Credential Provider type is specified.
func (r *credentialProviderV2Resource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.ExactlyOneOf(
			path.MatchRoot("aembit_access_token"),
			path.MatchRoot("api_key"),
			path.MatchRoot("aws_sts"),
			path.MatchRoot("google_workload_identity"),
			path.MatchRoot("snowflake_jwt"),
			path.MatchRoot("oauth_client_credentials"),
			path.MatchRoot("username_password"),
			path.MatchRoot("vault_client_token"),
		),
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *credentialProviderV2Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan credentialProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var credential aembit.CredentialProviderV2DTO = convertCredentialProviderModelToV2DTO(ctx, plan, nil, r.client.Tenant, r.client.StackDomain)

	// Create new Credential Provider
	credentialProvider, err := r.client.CreateCredentialProviderV2(credential, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Credential Provider",
			"Could not create Credential Provider, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan = convertCredentialProviderV2DTOToModel(ctx, *credentialProvider, plan, r.client.Tenant, r.client.StackDomain)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *credentialProviderV2Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state credentialProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed credential value from Aembit
	credentialProvider, err := r.client.GetCredentialProviderV2(state.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Error reading Aembit Credential Provider",
			"Could not read Aembit External ID from Terraform state "+state.ID.ValueString()+": "+err.Error(),
		)
		resp.State.RemoveResource(ctx)
		return
	}

	state = convertCredentialProviderV2DTOToModel(ctx, credentialProvider, state, r.client.Tenant, r.client.StackDomain)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *credentialProviderV2Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get current state
	var state credentialProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract external ID from state
	var externalID string = state.ID.ValueString()

	// Retrieve values from plan
	var plan credentialProviderResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var credential aembit.CredentialProviderV2DTO = convertCredentialProviderModelToV2DTO(ctx, plan, &externalID, r.client.Tenant, r.client.StackDomain)

	// Update Credential Provider
	credentialProvider, err := r.client.UpdateCredentialProviderV2(credential, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Credential Provider",
			"Could not update Credential Provider, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan = convertCredentialProviderV2DTOToModel(ctx, *credentialProvider, plan, r.client.Tenant, r.client.StackDomain)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *credentialProviderV2Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state credentialProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check if Credential Provider is Active - if it is, disable it first
	if state.IsActive == types.BoolValue(true) {
		_, err := r.client.DisableCredentialProvider(state.ID.ValueString(), nil)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error disabling Credential Provider",
				"Could not disable Credential Provider, unexpected error: "+err.Error(),
			)
			return
		}
	}

	// Delete existing Credential Provider
	_, err := r.client.DeleteCredentialProvider(state.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Credential Provider",
			"Could not delete Credential Provider, unexpected error: "+err.Error(),
		)
		return
	}
}

// Imports an existing resource by passing externalId.
func (r *credentialProviderV2Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import externalId and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func convertCredentialProviderModelToV2DTO(ctx context.Context, model credentialProviderResourceModel, externalID *string, tenantID string, stackDomain string) aembit.CredentialProviderV2DTO {
	var credential aembit.CredentialProviderV2DTO
	credential.EntityDTO = aembit.EntityDTO{
		Name:        model.Name.ValueString(),
		Description: model.Description.ValueString(),
		IsActive:    model.IsActive.ValueBool(),
	}
	if len(model.Tags.Elements()) > 0 {
		tagsMap := make(map[string]string)
		_ = model.Tags.ElementsAs(ctx, &tagsMap, true)

		for key, value := range tagsMap {
			credential.Tags = append(credential.Tags, aembit.TagDTO{
				Key:   key,
				Value: value,
			})
		}
	}
	if externalID != nil {
		credential.EntityDTO.ExternalID = *externalID
	}

	// Handle the Aembit Token use case
	if model.AembitToken != nil {
		credential.Type = "aembit-access-token"
		credential.Audience = fmt.Sprintf("%s.api.%s", tenantID, stackDomain)
		credential.CredentialAembitTokenV2DTO = aembit.CredentialAembitTokenV2DTO{
			RoleID:            model.AembitToken.Role.ValueString(),
			LifetimeInSeconds: model.AembitToken.Lifetime,
		}
	}

	// Handle the API Key use case
	if model.APIKey != nil {
		credential.Type = "apikey"
		credential.CredentialAPIKeyDTO = aembit.CredentialAPIKeyDTO{APIKey: model.APIKey.APIKey.ValueString()}
	}

	// Handle the AWS STS use case
	if model.AwsSTS != nil {
		credential.Type = "aws-sts-oidc"
		credential.Lifetime = model.AwsSTS.Lifetime
		credential.CredentialAwsSTSV2DTO = aembit.CredentialAwsSTSV2DTO{
			RoleArn: model.AwsSTS.RoleARN.ValueString(),
		}
	}

	// Handle the GCP Workload Identity Federation use case
	if model.GoogleWorkload != nil {
		credential.Type = "gcp-identity-federation"
		credential.Audience = model.GoogleWorkload.Audience.ValueString()
		credential.Lifetime = model.GoogleWorkload.Lifetime
		credential.CredentialGoogleWorkloadV2DTO = aembit.CredentialGoogleWorkloadV2DTO{
			ServiceAccount: model.GoogleWorkload.ServiceAccount.ValueString(),
		}
	}

	// Handle the Snowflake JWT use case
	if model.SnowflakeToken != nil {
		credential.Type = "signed-jwt"
		credential.Issuer = fmt.Sprintf("%s.%s.SHA256:{sha256(publicKey)}", model.SnowflakeToken.AccountID.ValueString(), model.SnowflakeToken.Username.ValueString())
		credential.Subject = fmt.Sprintf("%s.%s", model.SnowflakeToken.AccountID.ValueString(), model.SnowflakeToken.Username.ValueString())
		credential.Lifetime = 1
		credential.CredentialSnowflakeTokenV2DTO = aembit.CredentialSnowflakeTokenV2DTO{
			TokenConfiguration: "snowflake",
			AlgorithmType:      "RS256",
		}
	}

	// Handle the OAuth Client Credentials use case
	if model.OAuthClientCredentials != nil {
		credential.Type = "oauth-client-credential"
		credential.CredentialOAuthClientCredentialDTO = aembit.CredentialOAuthClientCredentialDTO{
			TokenURL:         model.OAuthClientCredentials.TokenURL.ValueString(),
			ClientID:         model.OAuthClientCredentials.ClientID.ValueString(),
			ClientSecret:     model.OAuthClientCredentials.ClientSecret.ValueString(),
			Scope:            model.OAuthClientCredentials.Scopes.ValueString(),
			CredentialStyle:  model.OAuthClientCredentials.CredentialStyle.ValueString(),
			CustomParameters: convertCredentialOAuthCustomParameters(model),
		}
	}

	// Handle the Username Password use case
	if model.UsernamePassword != nil {
		credential.Type = "username-password"
		credential.CredentialUsernamePasswordDTO = aembit.CredentialUsernamePasswordDTO{
			Username: model.UsernamePassword.Username.ValueString(),
			Password: model.UsernamePassword.Password.ValueString(),
		}
	}

	// Handle the Vault Client Token use case
	if model.VaultClientToken != nil {
		credential.Type = "vaultClientToken"
		credential.Issuer = fmt.Sprintf("https://%s.id.%s/", tenantID, stackDomain)
		credential.Subject = model.VaultClientToken.Subject
		credential.Lifetime = model.VaultClientToken.Lifetime
		credential.CredentialVaultClientTokenV2DTO = aembit.CredentialVaultClientTokenV2DTO{
			SubjectType:        model.VaultClientToken.SubjectType,
			CustomClaims:       make([]aembit.CredentialVaultClientTokenClaimsDTO, len(model.VaultClientToken.CustomClaims)),
			VaultHost:          model.VaultClientToken.VaultHost,
			Port:               model.VaultClientToken.VaultPort,
			TLS:                model.VaultClientToken.VaultTLS,
			Namespace:          model.VaultClientToken.VaultNamespace,
			Role:               model.VaultClientToken.VaultRole,
			AuthenticationPath: model.VaultClientToken.VaultPath,
			ForwardingConfig:   model.VaultClientToken.VaultForwarding,
		}
		for i, claim := range model.VaultClientToken.CustomClaims {
			credential.CredentialVaultClientTokenV2DTO.CustomClaims[i] = aembit.CredentialVaultClientTokenClaimsDTO{
				Key:       claim.Key,
				Value:     claim.Value,
				ValueType: claim.ValueType,
			}
		}
	}

	return credential
}

func convertCredentialProviderV2DTOToModel(ctx context.Context, dto aembit.CredentialProviderV2DTO, state credentialProviderResourceModel, tenant, stackDomain string) credentialProviderResourceModel {
	var model credentialProviderResourceModel
	model.ID = types.StringValue(dto.EntityDTO.ExternalID)
	model.Name = types.StringValue(dto.EntityDTO.Name)
	model.Description = types.StringValue(dto.EntityDTO.Description)
	model.IsActive = types.BoolValue(dto.EntityDTO.IsActive)
	model.Tags = newTagsModel(ctx, dto.EntityDTO.Tags)

	// Set the objects to null to begin with
	model.AembitToken = nil
	model.APIKey = nil
	model.AwsSTS = nil
	model.GoogleWorkload = nil
	model.OAuthClientCredentials = nil
	model.UsernamePassword = nil
	model.VaultClientToken = nil

	// Now fill in the objects based on the Credential Provider type
	switch dto.Type {
	case "aembit-access-token":
		model.AembitToken = convertAembitTokenV2DTOToModel(dto)
	case "apikey":
		model.APIKey = convertAPIKeyV2DTOToModel(dto, state)
	case "aws-sts-oidc":
		model.AwsSTS = convertAwsSTSV2DTOToModel(dto, tenant, stackDomain)
	case "gcp-identity-federation":
		model.GoogleWorkload = convertGoogleWorkloadV2DTOToModel(dto, tenant, stackDomain)
	case "signed-jwt":
		model.SnowflakeToken = convertSnowflakeTokenV2DTOToModel(dto)
	case "oauth-client-credential":
		model.OAuthClientCredentials = convertOAuthClientCredentialV2DTOToModel(dto, state)
	case "username-password":
		model.UsernamePassword = convertUserPassV2DTOToModel(dto, state)
	case "vaultClientToken":
		model.VaultClientToken = convertVaultClientTokenV2DTOToModel(dto, state)
	}
	return model
}

// convertAembitTokenV2DTOToModel converts the Aembit Token state object into a model ready for terraform processing.
func convertAembitTokenV2DTOToModel(dto aembit.CredentialProviderV2DTO) *credentialProviderAembitTokenModel {
	// First, parse the credentialProvider.ProviderDetail JSON returned from Aembit Cloud
	value := credentialProviderAembitTokenModel{
		Audience: types.StringValue(dto.Audience),
		Role:     types.StringValue(dto.RoleID),
		Lifetime: dto.LifetimeInSeconds,
	}
	return &value
}

// convertAPIKeyV2DTOToModel converts the API Key state object into a model ready for terraform processing.
// Note: Since Aembit vaults the API Key and does not return it in the API, the DTO will never contain the stored value.
func convertAPIKeyV2DTOToModel(_ aembit.CredentialProviderV2DTO, state credentialProviderResourceModel) *credentialProviderAPIKeyModel {
	value := credentialProviderAPIKeyModel{APIKey: types.StringNull()}
	if state.APIKey != nil {
		value.APIKey = state.APIKey.APIKey
	}
	return &value
}

// convertAwsSTSV2DTOToModel converts the AWS STS state object into a model ready for terraform processing.
func convertAwsSTSV2DTOToModel(dto aembit.CredentialProviderV2DTO, tenant, stackDomain string) *credentialProviderAwsSTSModel {
	value := credentialProviderAwsSTSModel{
		OIDCIssuer:    types.StringValue(fmt.Sprintf("https://%s.id.%s", tenant, stackDomain)),
		TokenAudience: types.StringValue("sts.amazonaws.com"),
		RoleARN:       types.StringValue(dto.RoleArn),
		Lifetime:      dto.Lifetime,
	}
	return &value
}

// convertGoogleWorkloadV2DTOToModel converts the Google Workload state object into a model ready for terraform processing.
func convertGoogleWorkloadV2DTOToModel(dto aembit.CredentialProviderV2DTO, tenant, stackDomain string) *credentialProviderGoogleWorkloadModel {
	value := credentialProviderGoogleWorkloadModel{
		OIDCIssuer:     types.StringValue(fmt.Sprintf("https://%s.id.%s", tenant, stackDomain)),
		Audience:       types.StringValue(dto.Audience),
		ServiceAccount: types.StringValue(dto.ServiceAccount),
		Lifetime:       dto.Lifetime,
	}
	return &value
}

// convertSnowflakeTokenV2DTOToModel converts the Snowflake JWT Token state object into a model ready for terraform processing.
func convertSnowflakeTokenV2DTOToModel(dto aembit.CredentialProviderV2DTO) *credentialProviderSnowflakeTokenModel {
	acctData := strings.Split(dto.Subject, ".")
	keyData := strings.ReplaceAll(dto.KeyContent, "\n", "")
	keyData = strings.Replace(keyData, "-----BEGIN PUBLIC KEY-----", "", 1)
	keyData = strings.Replace(keyData, "-----END PUBLIC KEY-----", "", 1)
	value := credentialProviderSnowflakeTokenModel{
		AccountID:        types.StringValue(acctData[0]),
		Username:         types.StringValue(acctData[1]),
		AlertUserCommand: types.StringValue(fmt.Sprintf("ALTER USER %s SET RSA_PUBLIC_KEY='%s'", acctData[1], keyData)),
	}
	return &value
}

// convertOAuthClientCredentialV2DTOToModel converts the OAuth Client Credential state object into a model ready for terraform processing.
// Note: Since Aembit vaults the Client Secret and does not return it in the API, the DTO will never contain the stored value.
func convertOAuthClientCredentialV2DTOToModel(dto aembit.CredentialProviderV2DTO, state credentialProviderResourceModel) *credentialProviderOAuthClientCredentialsModel {
	value := credentialProviderOAuthClientCredentialsModel{ClientSecret: types.StringNull()}
	value.TokenURL = types.StringValue(dto.TokenURL)
	value.ClientID = types.StringValue(dto.ClientID)
	value.Scopes = types.StringValue(dto.Scope)
	value.CredentialStyle = types.StringValue(dto.CredentialStyle)
	if state.OAuthClientCredentials != nil {
		value.ClientSecret = state.OAuthClientCredentials.ClientSecret
	}

	// Get the custom parameters to be injected into the model
	parameters := make([]*credentialProviderOAuthClientCustomParametersModel, len(dto.CustomParameters))
	for i, parameter := range dto.CustomParameters {
		parameters[i] = &credentialProviderOAuthClientCustomParametersModel{
			Key:       parameter.Key,
			Value:     parameter.Value,
			ValueType: parameter.ValueType,
		}
	}
	value.CustomParameters = parameters

	return &value
}

// convertUserPassV2DTOToModel converts the API Key state object into a model ready for terraform processing.
// Note: Since Aembit vaults the Password and does not return it in the API, the DTO will never contain the stored value.
func convertUserPassV2DTOToModel(dto aembit.CredentialProviderV2DTO, state credentialProviderResourceModel) *credentialProviderUserPassModel {
	value := credentialProviderUserPassModel{
		Username: types.StringValue(dto.Username),
		Password: types.StringNull(),
	}
	if state.UsernamePassword != nil {
		value.Password = state.UsernamePassword.Password
	}
	return &value
}

// convertVaultClientTokenV2DTOToModel converts the VaultClientToken state object into a model ready for terraform processing.
func convertVaultClientTokenV2DTOToModel(dto aembit.CredentialProviderV2DTO, _ credentialProviderResourceModel) *credentialProviderVaultClientTokenModel {
	value := credentialProviderVaultClientTokenModel{
		Subject:     dto.Subject,
		SubjectType: dto.SubjectType,
		Lifetime:    dto.Lifetime,

		VaultHost:       dto.VaultHost,
		VaultPort:       dto.Port,
		VaultTLS:        dto.TLS,
		VaultNamespace:  dto.Namespace,
		VaultRole:       dto.Role,
		VaultPath:       dto.AuthenticationPath,
		VaultForwarding: dto.ForwardingConfig,
	}

	// Get the custom claims to be injected into the model
	claims := make([]*credentialProviderVaultClientTokenCustomClaimsModel, len(dto.CustomClaims))
	//types.ObjectValue(credentialProviderVaultClientTokenCustomClaimsModel.AttrTypes),
	//claims := getSetObjectAttr(ctx, model.VaultClientToken, "custom_claims")
	for i, claim := range dto.CustomClaims {
		claims[i] = &credentialProviderVaultClientTokenCustomClaimsModel{
			Key:       claim.Key,
			Value:     claim.Value,
			ValueType: claim.ValueType,
		}
	}
	value.CustomClaims = claims
	return &value
}

package provider

import (
	"context"
	"fmt"
	"strings"
	"time"

	"terraform-provider-aembit/internal/provider/models"
	"terraform-provider-aembit/internal/provider/validators"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &credentialProviderResource{}
	_ resource.ResourceWithConfigure   = &credentialProviderResource{}
	_ resource.ResourceWithImportState = &credentialProviderResource{}
	_ resource.ResourceWithModifyPlan  = &credentialProviderResource{}
)

const oidcIssuerTemplate string = "https://%s.id.%s"

// NewCredentialProviderResource is a helper function to simplify the provider implementation.
func NewCredentialProviderResource() resource.Resource {
	return &credentialProviderResource{}
}

// credentialProviderResource is the resource implementation.
type credentialProviderResource struct {
	client *aembit.CloudClient
}

// Metadata returns the resource type name.
func (r *credentialProviderResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_credential_provider"
}

// Configure adds the provider configured client to the resource.
func (r *credentialProviderResource) Configure(
	_ context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	r.client = resourceConfigure(req, resp)
}

// Schema defines the schema for the resource.
func (r *credentialProviderResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			// ID field is required for Terraform Framework acceptance testing.
			"id": schema.StringAttribute{
				Description: "Unique identifier of the Credential Provider.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					validators.UUIDRegexValidation(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name for the Credential Provider.",
				Required:    true,
				Validators: []validator.String{
					validators.NameLengthValidation(),
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
			"tags":     TagsMapAttribute(),
			"tags_all": TagsAllMapAttribute(),
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
						Validators: []validator.Int64{
							int64validator.Between(900, 43200),
						},
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
						Validators: []validator.Int64{
							int64validator.Between(900, 43200),
						},
					},
				},
			},
			"google_workload_identity": schema.SingleNestedAttribute{
				Description: "Google Workload Identity Federation type Credential Provider configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"oidc_issuer": schema.StringAttribute{
						Description: "OIDC Issuer for GCP Workload Identity Federation configuration of the Credential Provider.",
						Computed:    true,
					},
					"audience": schema.StringAttribute{
						Description: "Audience for GCP Workload Identity Federation configuration of the Credential Provider.",
						Required:    true,
					},
					"service_account": schema.StringAttribute{
						Description: "Service Account email of the GCP Session credentials requested by the Credential Provider.",
						Required:    true,
						Validators: []validator.String{
							validators.EmailValidation(),
						},
					},
					"lifetime": schema.Int64Attribute{
						Description: "Lifetime (seconds) of the GCP Session credentials requested by the Credential Provider.",
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(3600),
						Validators: []validator.Int64{
							int64validator.Between(1, 43200),
						},
					},
				},
			},
			"azure_entra_workload_identity": schema.SingleNestedAttribute{
				Description: "Azure Entra Workload Identity Federation type Credential Provider configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"oidc_issuer": schema.StringAttribute{
						Description: "OIDC Issuer for Azure Entra Workload Identity Federation configuration of the Credential Provider.",
						Computed:    true,
					},
					"audience": schema.StringAttribute{
						Description: "Audience for Azure Entra Workload Identity Federation configuration of the Credential Provider.",
						Required:    true,
					},
					"subject": schema.StringAttribute{
						Description: "Subject for JWT Token for Azure Entra Workload Identity Federation configuration of the Credential Provider.",
						Required:    true,
					},
					"scope": schema.StringAttribute{
						Description: "Scope for Azure Entra Workload Identity Federation configuration of the Credential Provider.",
						Required:    true,
					},
					"azure_tenant": schema.StringAttribute{
						Description: "Azure Tenant ID for Azure Entra Workload Identity Federation configuration of the Credential Provider.",
						Required:    true,
						Validators: []validator.String{
							validators.UUIDRegexValidation(),
						},
					},
					"client_id": schema.StringAttribute{
						Description: "Azure Client ID for Azure Entra Workload Identity Federation configuration of the Credential Provider.",
						Required:    true,
						Validators: []validator.String{
							validators.UUIDRegexValidation(),
						},
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
						Validators: []validator.String{
							stringvalidator.LengthBetween(1, 24),
							validators.SnowflakeAccountValidation(),
						},
					},
					"username": schema.StringAttribute{
						Description: "Snowflake Username of the Credential Provider.",
						Required:    true,
						Validators: []validator.String{
							validators.SnowflakeUserNameValidation(),
						},
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
						Validators: []validator.String{
							validators.UrlSchemeValidation(),
						},
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
						Required:    true,
					},
					"credential_style": schema.StringAttribute{
						Description: "Credential Style for the OAuth Credential Provider.",
						Required:    true,
					},
					"custom_parameters": schema.SetNestedAttribute{
						Description: "Set Custom Parameters for the OAuth Credential Provider.",
						Optional:    true,
						Computed:    true,
						Default: setdefault.StaticValue(
							types.SetValueMust(types.ObjectType{AttrTypes: map[string]attr.Type{
								"key":        types.StringType,
								"value":      types.StringType,
								"value_type": types.StringType,
							}}, []attr.Value{}),
						),
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
									Description: "Type of value for Custom Parameter of the OAuth Credential Provider. Only `literal` is currently supported.",
									Optional:    true,
									Computed:    true,
									Default:     stringdefault.StaticString("literal"),
									Validators: []validator.String{
										stringvalidator.OneOf([]string{
											"literal",
										}...),
									},
								},
							},
						},
					},
				},
			},
			"oauth_authorization_code": schema.SingleNestedAttribute{
				Description: "OAuth Authorization Code Flow type Credential Provider configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"oauth_discovery_url": schema.StringAttribute{
						Description: "OAuth URL for the OAuth Credential Provider.",
						Required:    true,
						Validators: []validator.String{
							validators.SecureURLValidation(),
						},
					},
					"oauth_authorization_url": schema.StringAttribute{
						Description: "Authorization URL for the OAuth Credential Provider.",
						Required:    true,
					},
					"oauth_token_url": schema.StringAttribute{
						Description: "Token URL for the OAuth Credential Provider.",
						Required:    true,
					},
					"oauth_introspection_url": schema.StringAttribute{
						Description: "Introspection Url of the OAuth 2.0 introspection endpoint, used to validate and obtain metadata about access tokens",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
					},
					"user_authorization_url": schema.StringAttribute{
						Description: "3rd Party Authorization URL for User Consent for the OAuth Credential Provider.",
						Computed:    true,
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
						Required:    true,
					},
					"custom_parameters": schema.SetNestedAttribute{
						Description: "Set Custom Parameters for the OAuth Credential Provider.",
						Optional:    true,
						Computed:    true,
						Default: setdefault.StaticValue(
							types.SetValueMust(types.ObjectType{AttrTypes: map[string]attr.Type{
								"key":        types.StringType,
								"value":      types.StringType,
								"value_type": types.StringType,
							}}, []attr.Value{}),
						),
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
									Description: "Type of value for Custom Parameter of the OAuth Credential Provider. Only `literal` is currently supported.",
									Optional:    true,
									Computed:    true,
									Default:     stringdefault.StaticString("literal"),
									Validators: []validator.String{
										stringvalidator.OneOf([]string{
											"literal",
										}...),
									},
								},
							},
						},
					},
					"is_pkce_required": schema.BoolAttribute{
						Description: "PKCE required flag for the OAuth Credential Provider.",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
					},
					"callback_url": schema.StringAttribute{
						Description: "Authorization Callback URL for the OAuth Credential Provider. Callback URL can be pre-generated by matching the following template; \n" +
							"\t* https://AEMBIT_TENANT_ID.AEMBIT_STACK_DOMAIN/api/v1/credential-providers/RESOURCE_ID/callback \n" +
							"\t* RESOURCE_ID should be a valid uuid pre-generated for the Credential Provider resource. \n",
						Computed: true,
					},
					"state": schema.StringAttribute{
						Description: "State for the OAuth Credential Provider.",
						Computed:    true,
					},
					"lifetime": schema.Int64Attribute{
						Description: "Lifetime (in seconds) of the OAuth Authorization Code credentials requested by the Credential Provider.",
						Optional:    true,
						Computed:    true,
						Default:     int64default.StaticInt64(31536000),
						Validators: []validator.Int64{
							int64validator.AtLeast(86400),
						},
					},
					"lifetime_expiration": schema.StringAttribute{
						Description: "ISO 8601 formatted Lifetime Expiration of the OAuth Authorization Code credentials requested by the Credential Provider. This expiration timer begins when the user successfully completes an authorization of the Credential Provider and will be set to the authorization time plus the Credential Provider Lifetime value at that moment.",
						Computed:    true,
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
						Computed:    true,
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
						Default: setdefault.StaticValue(
							types.SetValueMust(types.ObjectType{AttrTypes: map[string]attr.Type{
								"key":        types.StringType,
								"value":      types.StringType,
								"value_type": types.StringType,
							}}, []attr.Value{}),
						),
					},
					"lifetime": schema.Int64Attribute{
						Description: "Lifetime (in seconds) of the JWT Token used to authenticate to the Vault Cluster. Note: The lifetime of the retrieved Vault Client Token is managed within Vault configuration.",
						Required:    true,
						Validators: []validator.Int64{
							int64validator.Between(1, 3600),
						},
					},
					"vault_host": schema.StringAttribute{
						Description: "Hostname of the Vault Cluster to be used for executing the login API.",
						Required:    true,
						Validators: []validator.String{
							validators.HostValidation(),
						},
					},
					"vault_port": schema.Int64Attribute{
						Description: "Port of the Vault Cluster to be used for executing the login API.",
						Required:    true,
						Validators: []validator.Int64{
							int64validator.Between(1, 65535),
						},
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
					"vault_private_network_access": schema.BoolAttribute{
						Description: "Indicates if the Vault instance operates within a private network.",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
					},
				},
			},
			"managed_gitlab_account": schema.SingleNestedAttribute{
				Description: "Managed GitLab Account type Credential Provider configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"service_account_username": schema.StringAttribute{
						Description: "Service Account Username value to be used for the service account created for this Credential Provider." +
							" If not specified, the service account username will be generated by Aembit using the following format:" +
							" ``Aembit_<credential_provider_name>_managed_service_account``.<br>**Note on Service Account Username Updates for Self-Managed GitLab or Dedicated GitLab Instances**<br>" +
							"The service_account_username attribute for this resource cannot be updated on Self-Managed or Dedicated GitLab instances after creation." +
							" This is due to a limitation in the GitLab API, which does not currently provide a method for modifying instance-level service accounts.<br>" +
							"Consequently, in order to change this attribute in your Terraform configuration you will need to destroy and reecreating the Credential Provider." +
							" This behavior does not affect GitLab.com SaaS instances.",
						Optional: true,
					},
					"group_ids": schema.SetAttribute{
						Description: "The set of GitLab group IDs.",
						ElementType: types.StringType,
						Required:    true,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(1),
							setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
						},
					},
					"project_ids": schema.SetAttribute{
						Description: "The set of GitLab project IDs.",
						ElementType: types.StringType,
						Required:    true,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(1),
							setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
						},
					},
					"access_level": schema.Int32Attribute{
						Description: "The access level of authorization. Valid values: 0 (No Access), 5 (Minimal Access), 10 (Guest), 15 (Planner), 20 (Reporter), 30 (Developer), 40 (Maintainer), 50 (Owner).",
						Required:    true,
						Validators: []validator.Int32{
							int32validator.OneOf([]int32{
								0,
								5,
								10,
								15,
								20,
								30,
								40,
								50,
							}...),
						},
					},
					"lifetime_in_hours": schema.Int32Attribute{
						Description: "Lifetime of the managed GitLab token in hours.",
						Required:    true,
						Validators: []validator.Int32{
							int32validator.Between(1, 8760),
						},
					},
					"scope": schema.StringAttribute{
						Description: "Scope for Managed Gitlab Account configuration of the Credential Provider.",
						Required:    true,
					},
					"credential_provider_integration_id": schema.StringAttribute{
						Description: "The unique identifier of the credential provider integration.",
						Required:    true,
					},
				},
			},
			"oidc_id_token": schema.SingleNestedAttribute{
				Description: "OIDC ID Token type Credential Provider configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"subject": schema.StringAttribute{
						Description: "Subject for JWT Token for OIDC ID Token configuration of the Credential Provider.",
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
					"issuer": schema.StringAttribute{
						Description: "OIDC Issuer for OIDC ID Token configuration of the Credential Provider.",
						Computed:    true,
					},
					"lifetime_in_minutes": schema.Int32Attribute{
						Description: "Lifetime of the Credential Provider in minutes.",
						Required:    true,
						Validators: []validator.Int32{
							int32validator.Between(1, 5256000), // max ten years
						},
					},
					"algorithm_type": schema.StringAttribute{
						Description: "JWT Signing algorithm type (RS256 or ES256)",
						Required:    true,
						Validators: []validator.String{
							stringvalidator.OneOf([]string{"RS256", "ES256"}...),
						},
					},
					"audience": schema.StringAttribute{
						Description: "Audience for OIDC ID Token configuration of the Credential Provider.",
						Required:    true,
					},
					"custom_claims": schema.SetNestedAttribute{
						Description: "Set of Custom Claims for the JWT Token.",
						Optional:    true,
						Computed:    true,
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
						Default: setdefault.StaticValue(
							types.SetValueMust(types.ObjectType{AttrTypes: map[string]attr.Type{
								"key":        types.StringType,
								"value":      types.StringType,
								"value_type": types.StringType,
							}}, []attr.Value{}),
						),
					},
				},
			},
			"aws_secrets_manager_value": schema.SingleNestedAttribute{
				Description: "AWS Secrets Manager Value type Credential Provider configuration. This type of credential provider" +
					" supports secret values in plaintext or JSON formats. Do not provide values in `secret_key_1` and `secret_key_2`" +
					" fields for plaintext secrets. These fields are used to specify property names when a secret contains a JSON.",
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"secret_arn": schema.StringAttribute{
						Description: "ARN of the AWS Secrets Manager secret to be used by the Credential Provider.",
						Required:    true,
						Validators: []validator.String{
							validators.AwsSecretArnValidation(),
						},
					},
					"secret_key_1": schema.StringAttribute{
						Description: "Used when an AWS Secrets Manager secret object is in JSON format. Specifies a key of an element with the secret value.",
						Optional:    true,
					},
					"secret_key_2": schema.StringAttribute{
						Description: "Similar to `secret_key_1` but used when you need a credential provider to work with 2 secret values." +
							" For example, a username / password pair.",
						Optional: true,
					},
					"private_network_access": schema.BoolAttribute{
						Description: "Indicates that the AWS Secrets Manager is accessible via a private network only.",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
					},
					"credential_provider_integration_id": schema.StringAttribute{
						Description: "The unique identifier of the Credential Provider Integration of type AWS IAM Role.",
						Required:    true,
					},
				},
			},
			"azure_key_vault_value": schema.SingleNestedAttribute{
				Description: "Azure Key Vault Value type Credential Provider configuration. This type of credential provider" +
					" supports secret values in plaintext formats.",
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"secret_name_1": schema.StringAttribute{
						Description: "Name of the Azure Key Vault secret to be used by the Credential Provider.",
						Required:    true,
					},
					"secret_name_2": schema.StringAttribute{
						Description: "Similar to `secret_name_1` but used when you need a credential provider to work with 2 secret values." +
							" For example, a username / password pair.",
						Optional: true,
					},
					"private_network_access": schema.BoolAttribute{
						Description: "Indicates that the Azure Key Vault is accessible via a private network only.",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
					},
					"credential_provider_integration_id": schema.StringAttribute{
						Description: "The unique identifier of the Credential Provider Integration of type Azure Entra Federation.",
						Required:    true,
					},
				},
			},
			"jwt_svid_token": schema.SingleNestedAttribute{
				Description: "JWT-SVID Token type Credential Provider configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"subject": schema.StringAttribute{
						Description: "Subject for JWT Token for JWT-SVID Token configuration of the Credential Provider.",
						Required:    true,
						Validators: []validator.String{
							validators.SpiffeSubjectValidation(),
						},
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
					"issuer": schema.StringAttribute{
						Description: "Issuer claim for JWT-SVID Token configuration of the Credential Provider.",
						Computed:    true,
					},
					"lifetime_in_minutes": schema.Int32Attribute{
						Description: "Lifetime of the Credential Provider in minutes.",
						Required:    true,
						Validators: []validator.Int32{
							int32validator.Between(1, 5256000), // max ten years
						},
					},
					"algorithm_type": schema.StringAttribute{
						Description: "JWT Signing algorithm type (RS256 or ES256)",
						Required:    true,
						Validators: []validator.String{
							stringvalidator.OneOf([]string{"RS256", "ES256"}...),
						},
					},
					"audience": schema.StringAttribute{
						Description: "Audience for JWT-SVID Token configuration of the Credential Provider.",
						Required:    true,
					},
					"custom_claims": schema.SetNestedAttribute{
						Description: "Set of Custom Claims for the JWT Token.",
						Optional:    true,
						Computed:    true,
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
						Default: setdefault.StaticValue(
							types.SetValueMust(types.ObjectType{AttrTypes: map[string]attr.Type{
								"key":        types.StringType,
								"value":      types.StringType,
								"value_type": types.StringType,
							}}, []attr.Value{}),
						),
					},
				},
			},
		},
	}
}

// Configure validators to ensure that only one Credential Provider type is specified.
func (r *credentialProviderResource) ConfigValidators(
	_ context.Context,
) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.ExactlyOneOf(
			path.MatchRoot("aembit_access_token"),
			path.MatchRoot("api_key"),
			path.MatchRoot("aws_sts"),
			path.MatchRoot("google_workload_identity"),
			path.MatchRoot("azure_entra_workload_identity"),
			path.MatchRoot("snowflake_jwt"),
			path.MatchRoot("oauth_client_credentials"),
			path.MatchRoot("oauth_authorization_code"),
			path.MatchRoot("username_password"),
			path.MatchRoot("vault_client_token"),
			path.MatchRoot("managed_gitlab_account"),
			path.MatchRoot("oidc_id_token"),
			path.MatchRoot("aws_secrets_manager_value"),
			path.MatchRoot("azure_key_vault_value"),
			path.MatchRoot("jwt_svid_token"),
		),
	}
}

func (r *credentialProviderResource) ModifyPlan(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
) {
	modifyPlanForTagsAll(ctx, req, resp, r.client.DefaultTags)
}

// Create creates the resource and sets the initial Terraform state.
func (r *credentialProviderResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	// Retrieve values from plan
	var plan models.CredentialProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	credential := convertCredentialProviderModelToV2DTO(
		ctx,
		plan,
		nil,
		r.client.Tenant,
		r.client.StackDomain,
		r.client.DefaultTags,
	)

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
	plan = convertCredentialProviderV2DTOToModel(
		ctx,
		*credentialProvider,
		&plan,
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
func (r *credentialProviderResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	// Get current state
	var state models.CredentialProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed credential value from Aembit
	credentialProvider, err, notFound := r.client.GetCredentialProviderV2(
		state.ID.ValueString(),
		nil,
	)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Error reading Aembit Credential Provider",
			"Could not read Aembit External ID from Terraform state "+state.ID.ValueString()+": "+err.Error(),
		)

		// If the resource is not found on Aembit Cloud, delete it locally
		if notFound {
			resp.State.RemoveResource(ctx)
		}
		return
	}

	state = convertCredentialProviderV2DTOToModel(
		ctx,
		credentialProvider,
		&state,
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
func (r *credentialProviderResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	// Get current state
	var state models.CredentialProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract external ID from state
	externalID := state.ID.ValueString()

	// Retrieve values from plan
	var plan models.CredentialProviderResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	credential := convertCredentialProviderModelToV2DTO(
		ctx,
		plan,
		&externalID,
		r.client.Tenant,
		r.client.StackDomain,
		r.client.DefaultTags,
	)

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
	plan = convertCredentialProviderV2DTOToModel(
		ctx,
		*credentialProvider,
		&plan,
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

// Delete deletes the resource and removes the Terraform state on success.
func (r *credentialProviderResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	// Retrieve values from state
	var state models.CredentialProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check if Credential Provider is Active - if it is, disable it first
	if state.IsActive == types.BoolValue(true) {
		_, err := r.client.DisableCredentialProviderV2(state.ID.ValueString(), nil)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error disabling Credential Provider",
				"Could not disable Credential Provider, unexpected error: "+err.Error(),
			)
			return
		}
	}

	// Delete existing Credential Provider
	_, err := r.client.DeleteCredentialProviderV2(ctx, state.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Credential Provider",
			"Could not delete Credential Provider, unexpected error: "+err.Error(),
		)
		return
	}
}

// Imports an existing resource by passing externalId.
func (r *credentialProviderResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	// Retrieve import externalId and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func convertCredentialProviderModelToV2DTO(
	ctx context.Context,
	model models.CredentialProviderResourceModel,
	externalID *string,
	tenantID string,
	stackDomain string,
	defaultTags map[string]string,
) aembit.CredentialProviderV2DTO {
	var credential aembit.CredentialProviderV2DTO
	credential.EntityDTO = aembit.EntityDTO{
		Name:        model.Name.ValueString(),
		Description: model.Description.ValueString(),
		IsActive:    model.IsActive.ValueBool(),
	}

	if externalID != nil {
		credential.ExternalID = *externalID
	}

	// Handle the Aembit Token use case
	if model.AembitToken != nil {
		convertToAembitTokenDTO(&credential, model, fmt.Sprintf("%s.api.%s", tenantID, stackDomain))
	}

	// Handle the API Key use case
	if model.APIKey != nil {
		convertToApiKeyDTO(&credential, model)
	}

	// Handle the AWS STS use case
	if model.AwsSTS != nil {
		convertToAwsSTSDTO(&credential, model)
	}

	// Handle the GCP Workload Identity Federation use case
	if model.GoogleWorkload != nil {
		convertToGoogleWorkloadDTO(&credential, model)
	}

	// Handle the Azure Entra Workload Identity Federation use case
	if model.AzureEntraWorkload != nil {
		convertToAzureEntraDTO(&credential, model)
	}

	// Handle the Snowflake JWT use case
	if model.SnowflakeToken != nil {
		convertToSnowflakeTokenDTO(&credential, model)
	}

	// Handle the OAuth Client Credentials use case
	if model.OAuthClientCredentials != nil {
		convertToOAuthClientCredentialsDTO(&credential, model)
	}

	// Handle the OAuth Authorization Code use case
	if model.OAuthAuthorizationCode != nil {
		convertToOAuthAuthorizationCodeDTO(&credential, model)
	}

	// Handle the Username Password use case
	if model.UsernamePassword != nil {
		convertToUsernamePasswordDTO(&credential, model)
	}

	// Handle the Vault Client Token use case
	if model.VaultClientToken != nil {
		convertToVaultClientTokenDTO(
			&credential,
			model,
			fmt.Sprintf(oidcIssuerTemplate, tenantID, stackDomain),
		)
	}

	// Handle the Managed Gitlab Account use case
	if model.ManagedGitlabAccount != nil {
		convertToManagedGitlabAccountDTO(&credential, model)
	}

	// Handle the OidcIdToken use case
	if model.OidcIdToken != nil {
		convertToOidcIdTokenDTO(
			&credential,
			*model.OidcIdToken,
			fmt.Sprintf(oidcIssuerTemplate, tenantID, stackDomain),
			"oidc-id-token",
		)
	}

	// Handle the JWT-SVID use case
	if model.JwtSvidToken != nil {
		convertToOidcIdTokenDTO(
			&credential,
			*model.JwtSvidToken,
			fmt.Sprintf(oidcIssuerTemplate, tenantID, stackDomain),
			"jwt-svid-token",
		)
	}

	// Handle the AWS Secret use case
	if model.AwsSecretsManagerValue != nil {
		convertToAwsSecretsManagerValueDTO(&credential, model)
	}

	// Handle the Azure Key Vault use case
	if model.AzureKeyVaultValue != nil {
		convertToAzureKeyVaultValueDTO(&credential, model)
	}

	credential.Tags = collectAllTagsDto(ctx, defaultTags, model.Tags)
	return credential
}

func convertCredentialProviderV2DTOToModel(
	ctx context.Context,
	dto aembit.CredentialProviderV2DTO,
	planModel *models.CredentialProviderResourceModel,
	tenant, stackDomain string,
) models.CredentialProviderResourceModel {
	var model models.CredentialProviderResourceModel
	model.ID = types.StringValue(dto.ExternalID)
	model.Name = types.StringValue(dto.Name)
	model.Description = types.StringValue(dto.Description)
	model.IsActive = types.BoolValue(dto.IsActive)
	// handle tags
	model.Tags = newTagsModelFromPlan(ctx, planModel.Tags)
	model.TagsAll = newTagsModel(ctx, dto.Tags)

	// Set the objects to null to begin with
	model.AembitToken = nil
	model.APIKey = nil
	model.AwsSTS = nil
	model.GoogleWorkload = nil
	model.AzureEntraWorkload = nil
	model.OAuthClientCredentials = nil
	model.UsernamePassword = nil
	model.VaultClientToken = nil
	model.ManagedGitlabAccount = nil
	model.OidcIdToken = nil
	model.AwsSecretsManagerValue = nil
	model.AzureKeyVaultValue = nil
	model.JwtSvidToken = nil

	// Now fill in the objects based on the Credential Provider type
	switch dto.Type {
	case "aembit-access-token":
		model.AembitToken = convertAembitTokenV2DTOToModel(dto)
	case "apikey":
		model.APIKey = convertAPIKeyV2DTOToModel(dto, *planModel)
	case "aws-sts-oidc":
		model.AwsSTS = convertAwsSTSV2DTOToModel(dto, tenant, stackDomain)
	case "gcp-identity-federation":
		model.GoogleWorkload = convertGoogleWorkloadV2DTOToModel(dto, tenant, stackDomain)
	case "azure-entra-federation":
		model.AzureEntraWorkload = convertAzureEntraWorkloadV2DTOToModel(dto, tenant, stackDomain)
	case "signed-jwt":
		model.SnowflakeToken = convertSnowflakeTokenV2DTOToModel(dto)
	case "oauth-client-credential":
		model.OAuthClientCredentials = convertOAuthClientCredentialV2DTOToModel(dto, *planModel)
	case "oauth-authorization-code":
		model.OAuthAuthorizationCode = convertOAuthAuthorizationCodeV2DTOToModel(dto, *planModel)
	case "username-password":
		model.UsernamePassword = convertUserPassV2DTOToModel(dto, *planModel)
	case "vaultClientToken":
		model.VaultClientToken = convertVaultClientTokenV2DTOToModel(dto, *planModel)
	case "gitlab-managed-account":
		model.ManagedGitlabAccount = convertManagedGitlabAccountDTOToModel(dto, *planModel)
	case "oidc-id-token":
		model.OidcIdToken = convertOidcIdTokenDTOToModel(dto, *planModel)
	case "aws-secret-manager-value":
		model.AwsSecretsManagerValue = &models.CredentialProviderAwsSecretsManagerValueModel{
			SecretArn:            types.StringValue(dto.SecretArn),
			SecretKey1:           types.StringValue(dto.SecretKey1),
			SecretKey2:           types.StringValue(dto.SecretKey2),
			PrivateNetworkAccess: types.BoolValue(dto.PrivateNetworkAccess),
			CredentialProviderIntegrationExternalId: types.StringValue(
				dto.CredentialProviderIntegrationExternalId,
			),
		}
	case "azure-key-vault-value":
		model.AzureKeyVaultValue = &models.CredentialProviderAzureKeyVaultValueModel{
			SecretName1:          types.StringValue(dto.SecretName1),
			SecretName2:          types.StringValue(dto.SecretName2),
			PrivateNetworkAccess: types.BoolValue(dto.PrivateNetworkAccess),
			CredentialProviderIntegrationExternalId: types.StringValue(
				dto.CredentialProviderIntegrationExternalId,
			),
		}
	case "jwt-svid-token":
		model.JwtSvidToken = convertOidcIdTokenDTOToModel(dto, *planModel)
	}
	return model
}

// convertAembitTokenV2DTOToModel converts the Aembit Token state object into a model ready for terraform processing.
func convertAembitTokenV2DTOToModel(
	dto aembit.CredentialProviderV2DTO,
) *models.CredentialProviderAembitTokenModel {
	// First, parse the credentialProvider.ProviderDetail JSON returned from Aembit Cloud
	value := models.CredentialProviderAembitTokenModel{
		Audience: types.StringValue(dto.Audience),
		Role:     types.StringValue(dto.RoleID),
		Lifetime: dto.LifetimeInSeconds,
	}
	return &value
}

// convertAPIKeyV2DTOToModel converts the API Key state object into a model ready for terraform processing.
// Note: Since Aembit vaults the API Key and does not return it in the API, the DTO will never contain the stored value.
func convertAPIKeyV2DTOToModel(
	_ aembit.CredentialProviderV2DTO,
	state models.CredentialProviderResourceModel,
) *models.CredentialProviderAPIKeyModel {
	value := models.CredentialProviderAPIKeyModel{APIKey: types.StringNull()}
	if state.APIKey != nil {
		value.APIKey = state.APIKey.APIKey
	}
	return &value
}

// convertAwsSTSV2DTOToModel converts the AWS STS state object into a model ready for terraform processing.
func convertAwsSTSV2DTOToModel(
	dto aembit.CredentialProviderV2DTO,
	tenant, stackDomain string,
) *models.CredentialProviderAwsSTSModel {
	value := models.CredentialProviderAwsSTSModel{
		OIDCIssuer:    types.StringValue(fmt.Sprintf(oidcIssuerTemplate, tenant, stackDomain)),
		TokenAudience: types.StringValue("sts.amazonaws.com"),
		RoleARN:       types.StringValue(dto.RoleArn),
		Lifetime:      dto.Lifetime,
	}
	return &value
}

// convertGoogleWorkloadV2DTOToModel converts the Google Workload state object into a model ready for terraform processing.
func convertGoogleWorkloadV2DTOToModel(
	dto aembit.CredentialProviderV2DTO,
	tenant, stackDomain string,
) *models.CredentialProviderGoogleWorkloadModel {
	value := models.CredentialProviderGoogleWorkloadModel{
		OIDCIssuer:     types.StringValue(fmt.Sprintf(oidcIssuerTemplate, tenant, stackDomain)),
		Audience:       types.StringValue(dto.Audience),
		ServiceAccount: types.StringValue(dto.ServiceAccount),
		Lifetime:       dto.Lifetime,
	}
	return &value
}

// convertAzureEntraWorkloadV2DTOTOModel converts the Azure Entra Workload state object into a model ready for terraform processing.
func convertAzureEntraWorkloadV2DTOToModel(
	dto aembit.CredentialProviderV2DTO,
	tenant, stackDomain string,
) *models.CredentialProviderAzureEntraWorkloadModel {
	value := models.CredentialProviderAzureEntraWorkloadModel{
		OIDCIssuer:  types.StringValue(fmt.Sprintf(oidcIssuerTemplate, tenant, stackDomain)),
		Audience:    types.StringValue(dto.Audience),
		Subject:     types.StringValue(dto.Subject),
		Scope:       types.StringValue(dto.Scope),
		AzureTenant: types.StringValue(dto.AzureTenant),
		ClientID:    types.StringValue(dto.ClientID),
	}
	return &value
}

// convertSnowflakeTokenV2DTOToModel converts the Snowflake JWT Token state object into a model ready for terraform processing.
func convertSnowflakeTokenV2DTOToModel(
	dto aembit.CredentialProviderV2DTO,
) *models.CredentialProviderSnowflakeTokenModel {
	acctData := strings.Split(dto.Subject, ".")
	keyData := strings.ReplaceAll(dto.KeyContent, "\n", "")
	keyData = strings.Replace(keyData, "-----BEGIN PUBLIC KEY-----", "", 1)
	keyData = strings.Replace(keyData, "-----END PUBLIC KEY-----", "", 1)
	value := models.CredentialProviderSnowflakeTokenModel{
		AccountID: types.StringValue(acctData[0]),
		Username:  types.StringValue(acctData[1]),
		AlertUserCommand: types.StringValue(
			fmt.Sprintf("ALTER USER %s SET RSA_PUBLIC_KEY='%s'", acctData[1], keyData),
		),
	}
	return &value
}

// convertOAuthClientCredentialV2DTOToModel converts the OAuth Client Credential state object into a model ready for terraform processing.
// Note: Since Aembit vaults the Client Secret and does not return it in the API, the DTO will never contain the stored value.
func convertOAuthClientCredentialV2DTOToModel(
	dto aembit.CredentialProviderV2DTO,
	state models.CredentialProviderResourceModel,
) *models.CredentialProviderOAuthClientCredentialsModel {
	value := models.CredentialProviderOAuthClientCredentialsModel{ClientSecret: types.StringNull()}
	value.TokenURL = types.StringValue(dto.TokenURL)
	value.ClientID = types.StringValue(dto.ClientID)
	value.Scopes = types.StringValue(dto.Scope)
	value.CredentialStyle = types.StringValue(dto.CredentialStyle)
	if state.OAuthClientCredentials != nil {
		value.ClientSecret = state.OAuthClientCredentials.ClientSecret
	}

	// Get the custom parameters to be injected into the model
	parameters := make(
		[]*models.CredentialProviderOAuthClientCustomParametersModel,
		len(dto.CustomParameters),
	)
	for i, parameter := range dto.CustomParameters {
		parameters[i] = &models.CredentialProviderOAuthClientCustomParametersModel{
			Key:       parameter.Key,
			Value:     parameter.Value,
			ValueType: parameter.ValueType,
		}
	}
	value.CustomParameters = parameters

	return &value
}

// convertOAuthAuthorizationCodeV2DTOToModel converts the OAuth Authorization Code state object into a model ready for terraform processing.
// Note: Since Aembit vaults the Client Secret and does not return it in the API, the DTO will never contain the stored value.
func convertOAuthAuthorizationCodeV2DTOToModel(
	dto aembit.CredentialProviderV2DTO,
	state models.CredentialProviderResourceModel,
) *models.CredentialProviderOAuthAuthorizationCodeModel {
	value := models.CredentialProviderOAuthAuthorizationCodeModel{ClientSecret: types.StringNull()}
	value.OAuthDiscoveryUrl = types.StringValue(dto.OAuthUrl)
	value.OAuthAuthorizationUrl = types.StringValue(dto.AuthorizationUrl)
	value.OAuthTokenUrl = types.StringValue(dto.TokenUrl)
	value.UserAuthorizationUrl = types.StringValue(dto.UserAuthorizationUrl)
	value.OAuthIntrospectionUrl = types.StringValue(dto.IntrospectionUrl)
	value.ClientID = types.StringValue(dto.ClientID)
	value.Scopes = types.StringValue(dto.Scope)
	value.IsPkceRequired = types.BoolValue(dto.IsPkceRequired)
	value.CallBackUrl = types.StringValue(dto.CallBackUrl)
	value.State = types.StringValue(dto.State)
	if state.OAuthAuthorizationCode != nil {
		value.ClientSecret = state.OAuthAuthorizationCode.ClientSecret
	}

	// Get the custom parameters to be injected into the model
	parameters := make(
		[]*models.CredentialProviderOAuthClientCustomParametersModel,
		len(dto.CustomParameters),
	)
	for i, parameter := range dto.CustomParameters {
		parameters[i] = &models.CredentialProviderOAuthClientCustomParametersModel{
			Key:       parameter.Key,
			Value:     parameter.Value,
			ValueType: parameter.ValueType,
		}
	}
	value.CustomParameters = parameters

	value.Lifetime = dto.LifetimeTimeSpanSeconds

	if dto.LifetimeExpiration != nil {
		// Add Z to indicate that Date string is in UTC format, API returns it without region info
		timeParsed, err := time.Parse(time.RFC3339, *dto.LifetimeExpiration+"Z")
		if err != nil {
			fmt.Println("Error parsing LifetimeExpiration:", err)
		}

		value.LifetimeExpiration = types.StringValue(timeParsed.Local().Format(time.RFC3339))
	}

	return &value
}

// convertUserPassV2DTOToModel converts the API Key state object into a model ready for terraform processing.
// Note: Since Aembit vaults the Password and does not return it in the API, the DTO will never contain the stored value.
func convertUserPassV2DTOToModel(
	dto aembit.CredentialProviderV2DTO,
	state models.CredentialProviderResourceModel,
) *models.CredentialProviderUserPassModel {
	value := models.CredentialProviderUserPassModel{
		Username: types.StringValue(dto.Username),
		Password: types.StringNull(),
	}
	if state.UsernamePassword != nil {
		value.Password = state.UsernamePassword.Password
	}
	return &value
}

// convertVaultClientTokenV2DTOToModel converts the VaultClientToken state object into a model ready for terraform processing.
func convertVaultClientTokenV2DTOToModel(
	dto aembit.CredentialProviderV2DTO,
	_ models.CredentialProviderResourceModel,
) *models.CredentialProviderVaultClientTokenModel {
	value := models.CredentialProviderVaultClientTokenModel{
		Subject:                   dto.Subject,
		SubjectType:               dto.SubjectType,
		Lifetime:                  dto.Lifetime,
		VaultHost:                 dto.VaultHost,
		VaultPort:                 dto.Port,
		VaultTLS:                  dto.TLS,
		VaultNamespace:            dto.Namespace,
		VaultRole:                 dto.Role,
		VaultPath:                 dto.AuthenticationPath,
		VaultForwarding:           dto.ForwardingConfig,
		VaultPrivateNetworkAccess: types.BoolValue(dto.PrivateNetworkAccess),
	}

	// Get the custom claims to be injected into the model
	claims := make([]*models.CredentialProviderCustomClaimsModel, len(dto.CustomClaims))
	// types.ObjectValue(models.CredentialProviderCustomClaimsModel.AttrTypes),
	// claims := getSetObjectAttr(ctx, model.VaultClientToken, "custom_claims")
	for i, claim := range dto.CustomClaims {
		claims[i] = &models.CredentialProviderCustomClaimsModel{
			Key:       claim.Key,
			Value:     claim.Value,
			ValueType: claim.ValueType,
		}
	}
	value.CustomClaims = claims
	return &value
}

// convertManagedGitlabAccountDTOToModel converts the GitlabManagedAccount state object into a model ready for terraform processing.
func convertManagedGitlabAccountDTOToModel(
	dto aembit.CredentialProviderV2DTO,
	_ models.CredentialProviderResourceModel,
) *models.CredentialProviderManagedGitlabAccountModel {
	value := models.CredentialProviderManagedGitlabAccountModel{
		ServiceAccountUsername: types.StringValue(dto.Username),
		GroupIds: convertSliceToSet(
			strings.Split(dto.GroupIds, ","),
		),
		ProjectIds: convertSliceToSet(
			strings.Split(dto.ProjectIds, ","),
		),
		LifetimeInHours:                         types.Int32Value(dto.LifetimeInSeconds / 3600),
		Scope:                                   dto.Scope,
		AccessLevel:                             dto.AccessLevel,
		CredentialProviderIntegrationExternalId: dto.CredentialProviderIntegrationExternalId,
	}

	return &value
}

// convertOidcIdTokenDTOToModel converts the OidcIdToken state object into a model ready for terraform processing.
func convertOidcIdTokenDTOToModel(
	dto aembit.CredentialProviderV2DTO,
	_ models.CredentialProviderResourceModel,
) *models.CredentialProviderManagedOidcIdToken {
	value := models.CredentialProviderManagedOidcIdToken{
		Subject:           dto.Subject,
		SubjectType:       dto.SubjectType,
		LifetimeInMinutes: dto.LifetimeTimeSpanSeconds / 60,
		Audience:          dto.Audience,
		AlgorithmType:     dto.AlgorithmType,
		Issuer:            types.StringValue(dto.Issuer),
	}

	// Get the custom claims to be injected into the model
	claims := make([]*models.CredentialProviderCustomClaimsModel, len(dto.CustomClaims))
	// types.ObjectValue(models.CredentialProviderCustomClaimsModel.AttrTypes),
	// claims := getSetObjectAttr(ctx, model.VaultClientToken, "custom_claims")
	for i, claim := range dto.CustomClaims {
		claims[i] = &models.CredentialProviderCustomClaimsModel{
			Key:       claim.Key,
			Value:     claim.Value,
			ValueType: claim.ValueType,
		}
	}
	value.CustomClaims = claims
	return &value
}

// Get the custom parameters to be injected into the model.
func convertCredentialOAuthClientCredentialsCustomParameters(
	model models.CredentialProviderResourceModel,
) []aembit.CustomClaimsDTO {
	parameters := make([]aembit.CustomClaimsDTO, len(model.OAuthClientCredentials.CustomParameters))
	for i, param := range model.OAuthClientCredentials.CustomParameters {
		parameters[i] = aembit.CustomClaimsDTO{
			Key:       param.Key,
			Value:     param.Value,
			ValueType: param.ValueType,
		}
	}
	return parameters
}

// Get the custom parameters to be injected into the model.
func convertCredentialOAuthAuthorizationCodeCustomParameters(
	model models.CredentialProviderResourceModel,
) []aembit.CustomClaimsDTO {
	parameters := make([]aembit.CustomClaimsDTO, len(model.OAuthAuthorizationCode.CustomParameters))
	for i, param := range model.OAuthAuthorizationCode.CustomParameters {
		parameters[i] = aembit.CustomClaimsDTO{
			Key:       param.Key,
			Value:     param.Value,
			ValueType: param.ValueType,
		}
	}
	return parameters
}

func convertToAembitTokenDTO(
	credential *aembit.CredentialProviderV2DTO,
	model models.CredentialProviderResourceModel,
	audience string,
) {
	credential.Type = "aembit-access-token"
	credential.LifetimeInSeconds = model.AembitToken.Lifetime
	credential.Audience = audience
	credential.CredentialAembitTokenV2DTO = aembit.CredentialAembitTokenV2DTO{
		RoleID: model.AembitToken.Role.ValueString(),
	}
}

func convertToApiKeyDTO(
	credential *aembit.CredentialProviderV2DTO,
	model models.CredentialProviderResourceModel,
) {
	credential.Type = "apikey"
	credential.CredentialAPIKeyDTO = aembit.CredentialAPIKeyDTO{
		APIKey: model.APIKey.APIKey.ValueString(),
	}
}

func convertToAwsSTSDTO(
	credential *aembit.CredentialProviderV2DTO,
	model models.CredentialProviderResourceModel,
) {
	credential.Type = "aws-sts-oidc"
	credential.Lifetime = model.AwsSTS.Lifetime
	credential.CredentialAwsSTSV2DTO = aembit.CredentialAwsSTSV2DTO{
		RoleArn:  model.AwsSTS.RoleARN.ValueString(),
		Lifetime: model.AwsSTS.Lifetime,
	}
}

func convertToGoogleWorkloadDTO(
	credential *aembit.CredentialProviderV2DTO,
	model models.CredentialProviderResourceModel,
) {
	credential.Type = "gcp-identity-federation"
	credential.Audience = model.GoogleWorkload.Audience.ValueString()
	credential.Lifetime = model.GoogleWorkload.Lifetime
	credential.CredentialGoogleWorkloadV2DTO = aembit.CredentialGoogleWorkloadV2DTO{
		ServiceAccount: model.GoogleWorkload.ServiceAccount.ValueString(),
	}
}

func convertToAzureEntraDTO(
	credential *aembit.CredentialProviderV2DTO,
	model models.CredentialProviderResourceModel,
) {
	credential.Type = "azure-entra-federation"
	credential.Audience = model.AzureEntraWorkload.Audience.ValueString()
	credential.Subject = model.AzureEntraWorkload.Subject.ValueString()
	credential.Scope = model.AzureEntraWorkload.Scope.ValueString()
	credential.ClientID = model.AzureEntraWorkload.ClientID.ValueString()
	credential.CredentialAzureEntraWorkloadV2DTO = aembit.CredentialAzureEntraWorkloadV2DTO{
		AzureTenant: model.AzureEntraWorkload.AzureTenant.ValueString(),
	}
}

func convertToSnowflakeTokenDTO(
	credential *aembit.CredentialProviderV2DTO,
	model models.CredentialProviderResourceModel,
) {
	credential.Type = "signed-jwt"
	credential.Issuer = fmt.Sprintf(
		"%s.%s.SHA256:{sha256(publicKey)}",
		model.SnowflakeToken.AccountID.ValueString(),
		model.SnowflakeToken.Username.ValueString(),
	)
	credential.Subject = fmt.Sprintf(
		"%s.%s",
		model.SnowflakeToken.AccountID.ValueString(),
		model.SnowflakeToken.Username.ValueString(),
	)
	credential.Lifetime = 1
	credential.AlgorithmType = "RS256"
	credential.CredentialSnowflakeTokenV2DTO = aembit.CredentialSnowflakeTokenV2DTO{
		TokenConfiguration: "snowflake",
	}
}

func convertToOAuthClientCredentialsDTO(
	credential *aembit.CredentialProviderV2DTO,
	model models.CredentialProviderResourceModel,
) {
	credential.Type = "oauth-client-credential"
	credential.ClientID = model.OAuthClientCredentials.ClientID.ValueString()
	credential.ClientSecret = model.OAuthClientCredentials.ClientSecret.ValueString()
	credential.Scope = model.OAuthClientCredentials.Scopes.ValueString()
	credential.CustomParameters = convertCredentialOAuthClientCredentialsCustomParameters(model)
	credential.CredentialOAuthClientCredentialV2DTO = aembit.CredentialOAuthClientCredentialV2DTO{
		TokenURL:        model.OAuthClientCredentials.TokenURL.ValueString(),
		CredentialStyle: model.OAuthClientCredentials.CredentialStyle.ValueString(),
	}
}

func convertToOAuthAuthorizationCodeDTO(
	credential *aembit.CredentialProviderV2DTO,
	model models.CredentialProviderResourceModel,
) {
	credential.Type = "oauth-authorization-code"
	credential.ClientID = model.OAuthAuthorizationCode.ClientID.ValueString()
	credential.ClientSecret = model.OAuthAuthorizationCode.ClientSecret.ValueString()
	credential.Scope = model.OAuthAuthorizationCode.Scopes.ValueString()
	credential.CustomParameters = convertCredentialOAuthAuthorizationCodeCustomParameters(model)
	credential.CredentialOAuthAuthorizationCodeV2DTO = aembit.CredentialOAuthAuthorizationCodeV2DTO{
		OAuthUrl:             model.OAuthAuthorizationCode.OAuthDiscoveryUrl.ValueString(),
		AuthorizationUrl:     model.OAuthAuthorizationCode.OAuthAuthorizationUrl.ValueString(),
		IntrospectionUrl:     model.OAuthAuthorizationCode.OAuthIntrospectionUrl.ValueString(),
		TokenUrl:             model.OAuthAuthorizationCode.OAuthTokenUrl.ValueString(),
		UserAuthorizationUrl: model.OAuthAuthorizationCode.UserAuthorizationUrl.ValueString(),
		IsPkceRequired:       model.OAuthAuthorizationCode.IsPkceRequired.ValueBool(),
		CallBackUrl:          model.OAuthAuthorizationCode.CallBackUrl.ValueString(),
		State:                model.OAuthAuthorizationCode.State.ValueString(),
	}
	if len(model.ID.ValueString()) > 0 {
		credential.ExternalID = model.ID.ValueString()
	}

	credential.LifetimeTimeSpanSeconds = model.OAuthAuthorizationCode.Lifetime
}

func convertToUsernamePasswordDTO(
	credential *aembit.CredentialProviderV2DTO,
	model models.CredentialProviderResourceModel,
) {
	credential.Type = "username-password"
	credential.CredentialUsernamePasswordDTO = aembit.CredentialUsernamePasswordDTO{
		Username: model.UsernamePassword.Username.ValueString(),
		Password: model.UsernamePassword.Password.ValueString(),
	}
}

func convertToVaultClientTokenDTO(
	credential *aembit.CredentialProviderV2DTO,
	model models.CredentialProviderResourceModel,
	issuer string,
) {
	credential.Type = "vaultClientToken"
	credential.Issuer = issuer
	credential.Subject = model.VaultClientToken.Subject
	credential.Lifetime = model.VaultClientToken.Lifetime
	credential.SubjectType = model.VaultClientToken.SubjectType
	credential.CredentialVaultClientTokenV2DTO = aembit.CredentialVaultClientTokenV2DTO{
		VaultHost:            model.VaultClientToken.VaultHost,
		Port:                 model.VaultClientToken.VaultPort,
		TLS:                  model.VaultClientToken.VaultTLS,
		Namespace:            model.VaultClientToken.VaultNamespace,
		Role:                 model.VaultClientToken.VaultRole,
		AuthenticationPath:   model.VaultClientToken.VaultPath,
		ForwardingConfig:     model.VaultClientToken.VaultForwarding,
		PrivateNetworkAccess: model.VaultClientToken.VaultPrivateNetworkAccess.ValueBool(),
	}
	credential.CustomClaims = make(
		[]aembit.CustomClaimsDTO,
		len(model.VaultClientToken.CustomClaims),
	)
	for i, claim := range model.VaultClientToken.CustomClaims {
		credential.CustomClaims[i] = aembit.CustomClaimsDTO{
			Key:       claim.Key,
			Value:     claim.Value,
			ValueType: claim.ValueType,
		}
	}
}

func convertToManagedGitlabAccountDTO(
	credential *aembit.CredentialProviderV2DTO,
	model models.CredentialProviderResourceModel,
) {
	credential.Type = "gitlab-managed-account"
	credential.Username = model.ManagedGitlabAccount.ServiceAccountUsername.ValueString()
	credential.GroupIds = strings.Join(
		convertSetToSlice(model.ManagedGitlabAccount.GroupIds),
		",",
	)
	credential.ProjectIds = strings.Join(
		convertSetToSlice(model.ManagedGitlabAccount.ProjectIds),
		",",
	)
	credential.LifetimeInSeconds = model.ManagedGitlabAccount.LifetimeInHours.ValueInt32() * 3600
	credential.AccessLevel = model.ManagedGitlabAccount.AccessLevel
	credential.Scope = model.ManagedGitlabAccount.Scope
	credential.CredentialProviderIntegrationExternalId = model.ManagedGitlabAccount.CredentialProviderIntegrationExternalId
}

func convertToOidcIdTokenDTO(
	credential *aembit.CredentialProviderV2DTO,
	oidcToken models.CredentialProviderManagedOidcIdToken,
	issuer string,
	credentialType string,
) {
	credential.Type = credentialType
	credential.LifetimeTimeSpanSeconds = oidcToken.LifetimeInMinutes * 60
	credential.Subject = oidcToken.Subject
	credential.SubjectType = oidcToken.SubjectType
	credential.Issuer = issuer
	credential.Audience = oidcToken.Audience
	credential.AlgorithmType = oidcToken.AlgorithmType

	credential.CustomClaims = make(
		[]aembit.CustomClaimsDTO,
		len(oidcToken.CustomClaims),
	)
	for i, claim := range oidcToken.CustomClaims {
		credential.CustomClaims[i] = aembit.CustomClaimsDTO{
			Key:       claim.Key,
			Value:     claim.Value,
			ValueType: claim.ValueType,
		}
	}
}

func convertToAwsSecretsManagerValueDTO(
	credential *aembit.CredentialProviderV2DTO,
	model models.CredentialProviderResourceModel,
) {
	credential.Type = "aws-secret-manager-value"
	credential.SecretArn = model.AwsSecretsManagerValue.SecretArn.ValueString()
	credential.SecretKey1 = model.AwsSecretsManagerValue.SecretKey1.ValueString()
	credential.SecretKey2 = model.AwsSecretsManagerValue.SecretKey2.ValueString()
	credential.PrivateNetworkAccess = model.AwsSecretsManagerValue.PrivateNetworkAccess.ValueBool()
	credential.CredentialProviderIntegrationExternalId = model.AwsSecretsManagerValue.CredentialProviderIntegrationExternalId.ValueString()
}

func convertToAzureKeyVaultValueDTO(
	credential *aembit.CredentialProviderV2DTO,
	model models.CredentialProviderResourceModel,
) {
	credential.Type = "azure-key-vault-value"
	credential.SecretName1 = model.AzureKeyVaultValue.SecretName1.ValueString()
	credential.SecretName2 = model.AzureKeyVaultValue.SecretName2.ValueString()
	credential.PrivateNetworkAccess = model.AzureKeyVaultValue.PrivateNetworkAccess.ValueBool()
	credential.CredentialProviderIntegrationExternalId = model.AzureKeyVaultValue.CredentialProviderIntegrationExternalId.ValueString()
}

func convertSetToSlice(set []types.String) []string {
	result := []string{}

	for _, val := range set {
		result = append(result, val.ValueString())
	}

	return result
}

func convertSliceToSet(slice []string) []types.String {
	result := []types.String{}

	for _, val := range slice {
		result = append(result, types.StringValue(val))
	}

	return result
}

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// CredentialProviderResourceModel maps the resource schema.
type CredentialProviderResourceModel struct {
	// ID is required for Framework acceptance testing
	ID                     types.String                                   `tfsdk:"id"`
	Name                   types.String                                   `tfsdk:"name"`
	Description            types.String                                   `tfsdk:"description"`
	IsActive               types.Bool                                     `tfsdk:"is_active"`
	Tags                   types.Map                                      `tfsdk:"tags"`
	AembitToken            *CredentialProviderAembitTokenModel            `tfsdk:"aembit_access_token"`
	APIKey                 *CredentialProviderAPIKeyModel                 `tfsdk:"api_key"`
	AwsSTS                 *CredentialProviderAwsSTSModel                 `tfsdk:"aws_sts"`
	GoogleWorkload         *CredentialProviderGoogleWorkloadModel         `tfsdk:"google_workload_identity"`
	AzureEntraWorkload     *CredentialProviderAzureEntraWorkloadModel     `tfsdk:"azure_entra_workload_identity"`
	SnowflakeToken         *CredentialProviderSnowflakeTokenModel         `tfsdk:"snowflake_jwt"`
	OAuthClientCredentials *CredentialProviderOAuthClientCredentialsModel `tfsdk:"oauth_client_credentials"`
	OAuthAuthorizationCode *CredentialProviderOAuthAuthorizationCodeModel `tfsdk:"oauth_authorization_code"`
	UsernamePassword       *CredentialProviderUserPassModel               `tfsdk:"username_password"`
	VaultClientToken       *CredentialProviderVaultClientTokenModel       `tfsdk:"vault_client_token"`
	ManagedGitlabAccount   *CredentialProviderManagedGitlabAccountModel   `tfsdk:"managed_gitlab_account"`
	OidcIdToken            *CredentialProviderManagedOidcIdToken          `tfsdk:"oidc_id_token"`
	AwsSecretsManagerValue *CredentialProviderAwsSecretsManagerValueModel `tfsdk:"aws_secrets_manager_value"`
	JwtSvidToken           *CredentialProviderManagedOidcIdToken          `tfsdk:"jwt_svid_token"`
}

// credentialProviderDataSourceModel maps the datasource schema.
type CredentialProvidersDataSourceModel struct {
	CredentialProviders []CredentialProviderResourceModel `tfsdk:"credential_providers"`
}

type CredentialProviderAembitTokenModel struct {
	Audience types.String `tfsdk:"audience"`
	Role     types.String `tfsdk:"role_id"`
	Lifetime int32        `tfsdk:"lifetime"`
}

type CredentialProviderAPIKeyModel struct {
	APIKey types.String `tfsdk:"api_key"`
}

type CredentialProviderAwsSTSModel struct {
	OIDCIssuer    types.String `tfsdk:"oidc_issuer"`
	RoleARN       types.String `tfsdk:"role_arn"`
	TokenAudience types.String `tfsdk:"token_audience"`
	Lifetime      int32        `tfsdk:"lifetime"`
}

type CredentialProviderGoogleWorkloadModel struct {
	OIDCIssuer     types.String `tfsdk:"oidc_issuer"`
	Audience       types.String `tfsdk:"audience"`
	ServiceAccount types.String `tfsdk:"service_account"`
	Lifetime       int32        `tfsdk:"lifetime"`
}

type CredentialProviderAzureEntraWorkloadModel struct {
	OIDCIssuer  types.String `tfsdk:"oidc_issuer"`
	Audience    types.String `tfsdk:"audience"`
	Subject     types.String `tfsdk:"subject"`
	Scope       types.String `tfsdk:"scope"`
	AzureTenant types.String `tfsdk:"azure_tenant"`
	ClientID    types.String `tfsdk:"client_id"`
}

type CredentialProviderSnowflakeTokenModel struct {
	AccountID        types.String `tfsdk:"account_id"`
	Username         types.String `tfsdk:"username"`
	AlertUserCommand types.String `tfsdk:"alter_user_command"`
}

// CredentialProviderOAuthClientCredentialsModel maps OAuth Client Credentials Flow configuration.
type CredentialProviderOAuthClientCredentialsModel struct {
	TokenURL         types.String                                          `tfsdk:"token_url"`
	ClientID         types.String                                          `tfsdk:"client_id"`
	ClientSecret     types.String                                          `tfsdk:"client_secret"`
	Scopes           types.String                                          `tfsdk:"scopes"`
	CredentialStyle  types.String                                          `tfsdk:"credential_style"`
	CustomParameters []*CredentialProviderOAuthClientCustomParametersModel `tfsdk:"custom_parameters"`
}

// CredentialProviderOAuthAuthorizationCodeModel maps OAuth Authorization Code Flow configuration.
type CredentialProviderOAuthAuthorizationCodeModel struct {
	OAuthDiscoveryUrl     types.String                                          `tfsdk:"oauth_discovery_url"`
	OAuthAuthorizationUrl types.String                                          `tfsdk:"oauth_authorization_url"`
	OAuthTokenUrl         types.String                                          `tfsdk:"oauth_token_url"`
	OAuthIntrospectionUrl types.String                                          `tfsdk:"oauth_introspection_url"`
	UserAuthorizationUrl  types.String                                          `tfsdk:"user_authorization_url"`
	ClientID              types.String                                          `tfsdk:"client_id"`
	ClientSecret          types.String                                          `tfsdk:"client_secret"`
	Scopes                types.String                                          `tfsdk:"scopes"`
	CustomParameters      []*CredentialProviderOAuthClientCustomParametersModel `tfsdk:"custom_parameters"`
	IsPkceRequired        types.Bool                                            `tfsdk:"is_pkce_required"`
	CallBackUrl           types.String                                          `tfsdk:"callback_url"`
	State                 types.String                                          `tfsdk:"state"`
	Lifetime              int64                                                 `tfsdk:"lifetime"`
	LifetimeExpiration    types.String                                          `tfsdk:"lifetime_expiration"`
}

type CredentialProviderOAuthClientCustomParametersModel struct {
	Key       string `tfsdk:"key"`
	Value     string `tfsdk:"value"`
	ValueType string `tfsdk:"value_type"`
}

type CredentialProviderUserPassModel struct {
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

// CredentialProviderVaultClientTokenModel maps Vault Client configuration.
type CredentialProviderVaultClientTokenModel struct {
	Subject                   string                                 `tfsdk:"subject"`
	SubjectType               string                                 `tfsdk:"subject_type"`
	CustomClaims              []*CredentialProviderCustomClaimsModel `tfsdk:"custom_claims"`
	Lifetime                  int32                                  `tfsdk:"lifetime"`
	VaultHost                 string                                 `tfsdk:"vault_host"`
	VaultTLS                  bool                                   `tfsdk:"vault_tls"`
	VaultPort                 int32                                  `tfsdk:"vault_port"`
	VaultNamespace            string                                 `tfsdk:"vault_namespace"`
	VaultRole                 string                                 `tfsdk:"vault_role"`
	VaultPath                 string                                 `tfsdk:"vault_path"`
	VaultForwarding           string                                 `tfsdk:"vault_forwarding"`
	VaultPrivateNetworkAccess types.Bool                             `tfsdk:"vault_private_network_access"`
}

// CredentialProviderManagedGitlabAccountModel maps Managed Gitlab Account configuration.
type CredentialProviderManagedGitlabAccountModel struct {
	ServiceAccountUsername                  types.String           `tfsdk:"service_account_username"`
	GroupIds                                []types.String         `tfsdk:"group_ids"`
	ProjectIds                              []types.String         `tfsdk:"project_ids"`
	AccessLevel                             int32                  `tfsdk:"access_level"`
	LifetimeInDays                          basetypes.Float32Value `tfsdk:"lifetime_in_days"`
	LifetimeInHours                         basetypes.Int32Value   `tfsdk:"lifetime_in_hours"`
	Scope                                   string                 `tfsdk:"scope"`
	CredentialProviderIntegrationExternalId string                 `tfsdk:"credential_provider_integration_id"`
}

// CredentialProviderManagedOidcIdToken maps OIDC ID Token configuration.
type CredentialProviderManagedOidcIdToken struct {
	Subject           string                                 `tfsdk:"subject"`
	SubjectType       string                                 `tfsdk:"subject_type"`
	LifetimeInMinutes int64                                  `tfsdk:"lifetime_in_minutes"`
	Audience          string                                 `tfsdk:"audience"`
	AlgorithmType     string                                 `tfsdk:"algorithm_type"`
	Issuer            types.String                           `tfsdk:"issuer"`
	CustomClaims      []*CredentialProviderCustomClaimsModel `tfsdk:"custom_claims"`
}

type CredentialProviderCustomClaimsModel struct {
	Key       string `tfsdk:"key"`
	Value     string `tfsdk:"value"`
	ValueType string `tfsdk:"value_type"`
}

type CredentialProviderAwsSecretsManagerValueModel struct {
	SecretArn                               types.String `tfsdk:"secret_arn"`
	SecretKey1                              types.String `tfsdk:"secret_key_1"`
	SecretKey2                              types.String `tfsdk:"secret_key_2"`
	PrivateNetworkAccess                    types.Bool   `tfsdk:"private_network_access"`
	CredentialProviderIntegrationExternalId types.String `tfsdk:"credential_provider_integration_id"`
}

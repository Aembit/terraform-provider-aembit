package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TrustProviderResourceModel maps the resource schema.
type TrustProviderResourceModel struct {
	// ID is required for Framework acceptance testing
	ID                 types.String                     `tfsdk:"id"`
	Name               types.String                     `tfsdk:"name"`
	Description        types.String                     `tfsdk:"description"`
	IsActive           types.Bool                       `tfsdk:"is_active"`
	Tags               types.Map                        `tfsdk:"tags"`
	AzureMetadata      *TrustProviderAzureMetadataModel `tfsdk:"azure_metadata"`
	AwsRole            *TrustProviderAwsRoleModel       `tfsdk:"aws_role"`
	AwsMetadata        *TrustProviderAwsMetadataModel   `tfsdk:"aws_metadata"`
	GcpIdentity        *TrustProviderGcpIdentityModel   `tfsdk:"gcp_identity"`
	GitHubAction       *TrustProviderGitHubActionModel  `tfsdk:"github_action"`
	GitLabJob          *TrustProviderGitLabJobModel     `tfsdk:"gitlab_job"`
	Kerberos           *TrustProviderKerberosModel      `tfsdk:"kerberos"`
	KubernetesService  *TrustProviderKubernetesModel    `tfsdk:"kubernetes_service_account"`
	OidcIdToken        *TrustProviderOidcIdTokenModel   `tfsdk:"oidc_id_token"`
	TerraformWorkspace *TrustProviderTerraformModel     `tfsdk:"terraform_workspace"`
}

// trustProviderDataSourceModel maps the datasource schema.
type TrustProvidersDataSourceModel struct {
	TrustProviders []TrustProviderResourceModel `tfsdk:"trust_providers"`
}

type TrustProviderAzureMetadataModel struct {
	Sku             types.String   `tfsdk:"sku"`
	Skus            []types.String `tfsdk:"skus"`
	VMID            types.String   `tfsdk:"vm_id"`
	VMIDs           []types.String `tfsdk:"vm_ids"`
	SubscriptionID  types.String   `tfsdk:"subscription_id"`
	SubscriptionIDs []types.String `tfsdk:"subscription_ids"`
}

type TrustProviderAwsRoleModel struct {
	AccountID    types.String   `tfsdk:"account_id"`
	AccountIDs   []types.String `tfsdk:"account_ids"`
	AssumedRole  types.String   `tfsdk:"assumed_role"`
	AssumedRoles []types.String `tfsdk:"assumed_roles"`
	RoleARN      types.String   `tfsdk:"role_arn"`
	RoleARNs     []types.String `tfsdk:"role_arns"`
	Username     types.String   `tfsdk:"username"`
	Usernames    []types.String `tfsdk:"usernames"`
}

type TrustProviderAwsMetadataModel struct {
	Certificate             types.String   `tfsdk:"certificate"`
	AccountID               types.String   `tfsdk:"account_id"`
	AccountIDs              []types.String `tfsdk:"account_ids"`
	Architecture            types.String   `tfsdk:"architecture"`
	AvailabilityZone        types.String   `tfsdk:"availability_zone"`
	AvailabilityZones       []types.String `tfsdk:"availability_zones"`
	BillingProducts         types.String   `tfsdk:"billing_products"`
	ImageID                 types.String   `tfsdk:"image_id"`
	InstanceID              types.String   `tfsdk:"instance_id"`
	InstanceIDs             []types.String `tfsdk:"instance_ids"`
	InstanceType            types.String   `tfsdk:"instance_type"`
	InstanceTypes           []types.String `tfsdk:"instance_types"`
	KernelID                types.String   `tfsdk:"kernel_id"`
	MarketplaceProductCodes types.String   `tfsdk:"marketplace_product_codes"`
	PendingTime             types.String   `tfsdk:"pending_time"`
	PrivateIP               types.String   `tfsdk:"private_ip"`
	RamdiskID               types.String   `tfsdk:"ramdisk_id"`
	Region                  types.String   `tfsdk:"region"`
	Regions                 []types.String `tfsdk:"regions"`
	Version                 types.String   `tfsdk:"version"`
}

type TrustProviderKerberosModel struct {
	AgentControllerIDs []types.String `tfsdk:"agent_controller_ids"`
	Principal          types.String   `tfsdk:"principal"`
	Principals         []types.String `tfsdk:"principals"`
	RealmOrDomain      types.String   `tfsdk:"realm_domain"`
	RealmsOrDomains    []types.String `tfsdk:"realms_domains"`
	SourceIP           types.String   `tfsdk:"source_ip"`
	SourceIPs          []types.String `tfsdk:"source_ips"`
}

type TrustProviderKubernetesModel struct {
	Issuer              types.String     `tfsdk:"issuer"`
	Issuers             []types.String   `tfsdk:"issuers"`
	Namespace           types.String     `tfsdk:"namespace"`
	Namespaces          []types.String   `tfsdk:"namespaces"`
	PodName             types.String     `tfsdk:"pod_name"`
	PodNames            []types.String   `tfsdk:"pod_names"`
	ServiceAccountName  types.String     `tfsdk:"service_account_name"`
	ServiceAccountNames []types.String   `tfsdk:"service_account_names"`
	Subject             types.String     `tfsdk:"subject"`
	Subjects            []types.String   `tfsdk:"subjects"`
	OIDCEndpoint        types.String     `tfsdk:"oidc_endpoint"`
	PublicKey           types.String     `tfsdk:"public_key"`
	Jwks                JsonWebKeysModel `tfsdk:"jwks"`
}

type TrustProviderOidcIdTokenModel struct {
	Issuer       types.String     `tfsdk:"issuer"`
	Issuers      []types.String   `tfsdk:"issuers"`
	Subject      types.String     `tfsdk:"subject"`
	Subjects     []types.String   `tfsdk:"subjects"`
	Audience     types.String     `tfsdk:"audience"`
	Audiences    []types.String   `tfsdk:"audiences"`
	OIDCEndpoint types.String     `tfsdk:"oidc_endpoint"`
	PublicKey    types.String     `tfsdk:"public_key"`
	Jwks         JsonWebKeysModel `tfsdk:"jwks"`
}

type JsonWebKeysModel struct {
	Keys []JsonWebKeyModel `tfsdk:"keys"`
}

type JsonWebKeyModel struct {
	Kid types.String `tfsdk:"kid"`
	Kty types.String `tfsdk:"kty"`
	Use types.String `tfsdk:"use"`
	Alg types.String `tfsdk:"alg"`

	// RSA
	N types.String `tfsdk:"n"`
	E types.String `tfsdk:"e"`

	// EC
	X   types.String `tfsdk:"x"`
	Y   types.String `tfsdk:"y"`
	Crv types.String `tfsdk:"crv"`
}

type TrustProviderGcpIdentityModel struct {
	EMail  types.String   `tfsdk:"email"`
	EMails []types.String `tfsdk:"emails"`
}

type TrustProviderGitHubActionModel struct {
	Actor        types.String   `tfsdk:"actor"`
	Actors       []types.String `tfsdk:"actors"`
	Repository   types.String   `tfsdk:"repository"`
	Repositories []types.String `tfsdk:"repositories"`
	Workflow     types.String   `tfsdk:"workflow"`
	Workflows    []types.String `tfsdk:"workflows"`
}

type TrustProviderGitLabJobModel struct {
	OIDCEndpoint   types.String   `tfsdk:"oidc_endpoint"`
	OIDCClientID   types.String   `tfsdk:"oidc_client_id"`
	OIDCAudience   types.String   `tfsdk:"oidc_audience"`
	NamespacePath  types.String   `tfsdk:"namespace_path"`
	NamespacePaths []types.String `tfsdk:"namespace_paths"`
	ProjectPath    types.String   `tfsdk:"project_path"`
	ProjectPaths   []types.String `tfsdk:"project_paths"`
	RefPath        types.String   `tfsdk:"ref_path"`
	RefPaths       []types.String `tfsdk:"ref_paths"`
	Subject        types.String   `tfsdk:"subject"`
	Subjects       []types.String `tfsdk:"subjects"`
}

type TrustProviderTerraformModel struct {
	OrganizationID  types.String   `tfsdk:"organization_id"`
	OrganizationIDs []types.String `tfsdk:"organization_ids"`
	ProjectID       types.String   `tfsdk:"project_id"`
	ProjectIDs      []types.String `tfsdk:"project_ids"`
	WorkspaceID     types.String   `tfsdk:"workspace_id"`
	WorkspaceIDs    []types.String `tfsdk:"workspace_ids"`
}

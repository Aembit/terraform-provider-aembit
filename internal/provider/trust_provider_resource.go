package provider

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"terraform-provider-aembit/internal/provider/models"
	"terraform-provider-aembit/internal/provider/validators"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &trustProviderResource{}
	_ resource.ResourceWithConfigure   = &trustProviderResource{}
	_ resource.ResourceWithImportState = &trustProviderResource{}
)

// NewTrustProviderResource is a helper function to simplify the provider implementation.
func NewTrustProviderResource() resource.Resource {
	return &trustProviderResource{}
}

// trustProviderResource is the resource implementation.
type trustProviderResource struct {
	client *aembit.CloudClient
}

// Metadata returns the resource type name.
func (r *trustProviderResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_trust_provider"
}

// Configure adds the provider configured client to the resource.
func (r *trustProviderResource) Configure(
	_ context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	r.client = resourceConfigure(req, resp)
}

// Schema defines the schema for the resource.
func (r *trustProviderResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "**Note:** One and only one nested schema (e.g. `aws_metadata`) must be provided for the Trust Provider to be configured.",
		Attributes: map[string]schema.Attribute{
			// ID field is required for Terraform Framework acceptance testing.
			"id": schema.StringAttribute{
				Description: "Unique identifier of the Trust Provider.",
				Computed:    true,
				Validators: []validator.String{
					validators.UUIDRegexValidation(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name for the Trust Provider.",
				Required:    true,
				Validators: []validator.String{
					validators.NameLengthValidation(),
				},
			},
			"description": schema.StringAttribute{
				Description: "Description for the Trust Provider.",
				Optional:    true,
				Computed:    true,
			},
			"is_active": schema.BoolAttribute{
				Description: "Active status of the Trust Provider.",
				Optional:    true,
				Computed:    true,
			},
			"tags": schema.MapAttribute{
				Description: "Tags are key-value pairs.",
				ElementType: types.StringType,
				Optional:    true,
			},
			"azure_metadata": schema.SingleNestedAttribute{
				Description: "Azure Metadata type Trust Provider configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"sku": schema.StringAttribute{
						Description: "Specific SKU for the Virtual Machine image.",
						Optional:    true,
					},
					"skus": schema.SetAttribute{
						Description: "The set of accepted Azure SKUs which are hosting the Client Workloads. Used only for cases where multiple Azure SKUs can be matched.",
						ElementType: types.StringType,
						Optional:    true,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(2),
							setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
						},
					},
					"vm_id": schema.StringAttribute{
						Description: "Unique identifier for the Virtual Machine.",
						Optional:    true,
					},
					"vm_ids": schema.SetAttribute{
						Description: "The set of accepted Azure VM IDs which are hosting the Client Workloads. Used only for cases where multiple Azure VM IDs can be matched.",
						ElementType: types.StringType,
						Optional:    true,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(2),
							setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
						},
					},
					"subscription_id": schema.StringAttribute{
						Description: "Azure subscription for the Virtual Machine.",
						Optional:    true,
					},
					"subscription_ids": schema.SetAttribute{
						Description: "The set of accepted Azure Subscription IDs which are hosting the Client Workloads. Used only for cases where multiple Azure Subscription IDs can be matched.",
						ElementType: types.StringType,
						Optional:    true,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(2),
							setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
						},
					},
				},
			},
			"aws_role": schema.SingleNestedAttribute{
				Description: "AWS Role type Trust Provider configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"account_id": schema.StringAttribute{
						Description: "The ID of the AWS account that is hosting the Client Workload.",
						Optional:    true,
					},
					"account_ids": schema.SetAttribute{
						Description: "The set of accepted AWS Account IDs which are hosting the Client Workloads. Used only for cases where multiple AWS Account IDs can be matched.",
						ElementType: types.StringType,
						Optional:    true,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(2),
							setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
						},
					},
					"assumed_role": schema.StringAttribute{
						Description: "The Name of the AWS IAM Role which is running the Client Workload.",
						Optional:    true,
					},
					"assumed_roles": schema.SetAttribute{
						Description: "The set of accepted AWS Assumed Roles which are hosting the Client Workloads. Used only for cases where multiple AWS Assumed Roles can be matched.",
						ElementType: types.StringType,
						Optional:    true,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(2),
							setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
						},
					},
					"role_arn": schema.StringAttribute{
						Description: "The ARN of the AWS IAM Role which is running the Client Workload.",
						Optional:    true,
					},
					"role_arns": schema.SetAttribute{
						Description: "The set of accepted Role ARNs which are hosting the Client Workloads. Used only for cases where multiple Role ARNs can be matched.",
						ElementType: types.StringType,
						Optional:    true,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(2),
							setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
						},
					},
					"username": schema.StringAttribute{
						Description: "The UserID of the AWS IAM Account which is running the Client Workload (not commonly used).",
						Optional:    true,
					},
					"usernames": schema.SetAttribute{
						Description: "The set of accepted AWS UserIDs which are hosting the Client Workloads. Used only for cases where multiple AWS UserIDs can be matched.",
						ElementType: types.StringType,
						Optional:    true,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(2),
							setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
						},
					},
				},
			},
			"aws_metadata": schema.SingleNestedAttribute{
				Description: "AWS Metadata type Trust Provider configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"certificate": schema.StringAttribute{
						Description: "PEM Certificate to be used for Signature verification.",
						Optional:    true,
					},
					"account_id": schema.StringAttribute{
						Description: "The ID of the AWS account that launched the instance.",
						Optional:    true,
					},
					"account_ids": schema.SetAttribute{
						Description: "The set of accepted AWS account IDs which are hosting the Client Workloads. Used only for cases where multiple AWS account IDs can be matched.",
						ElementType: types.StringType,
						Optional:    true,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(2),
							setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
						},
					},
					"architecture": schema.StringAttribute{
						Description: "The architecture of the AMI used to launch the instance (i386 | x86_64 | arm64).",
						Optional:    true,
					},
					"availability_zone": schema.StringAttribute{
						Description: "The Availability Zone in which the instance is running.",
						Optional:    true,
					},
					"availability_zones": schema.SetAttribute{
						Description: "The set of accepted AWS Availability Zones which are hosting the Client Workloads. Used only for cases where multiple AWS Availability Zones can be matched.",
						ElementType: types.StringType,
						Optional:    true,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(2),
							setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
						},
					},
					"billing_products": schema.StringAttribute{
						Description: "The billing products of the instance.",
						Optional:    true,
					},
					"image_id": schema.StringAttribute{
						Description: "The ID of the AMI used to launch the instance.",
						Optional:    true,
					},
					"instance_id": schema.StringAttribute{
						Description: "The ID of the instance.",
						Optional:    true,
					},
					"instance_ids": schema.SetAttribute{
						Description: "The set of accepted AWS Instance IDs which are hosting the Client Workloads. Used only for cases where multiple AWS Instance IDs can be matched.",
						ElementType: types.StringType,
						Optional:    true,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(2),
							setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
						},
					},
					"instance_type": schema.StringAttribute{
						Description: "The instance type of the instance.",
						Optional:    true,
					},
					"instance_types": schema.SetAttribute{
						Description: "The set of accepted AWS Instance types which are hosting the Client Workloads. Used only for cases where multiple AWS Instance types can be matched.",
						ElementType: types.StringType,
						Optional:    true,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(2),
							setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
						},
					},
					"kernel_id": schema.StringAttribute{
						Description: "The ID of the kernel associated with the instance, if applicable.",
						Optional:    true,
					},
					"marketplace_product_codes": schema.StringAttribute{
						Description: "The AWS Marketplace product code of the AMI used to launch the instance.",
						Optional:    true,
					},
					"pending_time": schema.StringAttribute{
						Description: "The date and time that the instance was launched.",
						Optional:    true,
					},
					"private_ip": schema.StringAttribute{
						Description: "The private IPv4 address of the instance.",
						Optional:    true,
					},
					"ramdisk_id": schema.StringAttribute{
						Description: "The ID of the RAM disk associated with the instance, if applicable.",
						Optional:    true,
					},
					"region": schema.StringAttribute{
						Description: "The Region in which the instance is running.",
						Optional:    true,
					},
					"regions": schema.SetAttribute{
						Description: "The set of accepted AWS Regions which are hosting the Client Workloads. Used only for cases where multiple AWS Regions can be matched.",
						ElementType: types.StringType,
						Optional:    true,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(2),
							setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
						},
					},
					"version": schema.StringAttribute{
						Description: "The version of the instance identity document format.",
						Optional:    true,
					},
				},
			},
			"gcp_identity": schema.SingleNestedAttribute{
				Description: "GCP Identity type Trust Provider configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"email": schema.StringAttribute{
						Description: "The GCP Service Account email address associated with the resource.",
						Optional:    true,
						Validators: []validator.String{
							validators.EmailValidation(),
						},
					},
					"emails": schema.SetAttribute{
						Description: "A set of GCP Service Account email addresses that are associated with the resource(s).",
						ElementType: types.StringType,
						Optional:    true,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(2),
							setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
							setvalidator.ValueStringsAre(validators.EmailValidation()),
						},
					},
				},
			},
			"github_action": schema.SingleNestedAttribute{
				Description: "GitHub Action type Trust Provider configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"actor": schema.StringAttribute{
						Description: "The GitHub Actor which initiated the GitHub Action.",
						Optional:    true,
					},
					"actors": schema.SetAttribute{
						Description: "The set of accepted GitHub ID Token Actors which initiated the GitHub Action. Used only for cases where multiple GitHub ID Token Actors can be matched.",
						ElementType: types.StringType,
						Optional:    true,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(2),
							setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
						},
					},
					"repository": schema.StringAttribute{
						Description: "The GitHub Repository associated with the GitHub Action ID Token.",
						Optional:    true,
					},
					"repositories": schema.SetAttribute{
						Description: "The set of accepted GitHub ID Token Repositories which initiated the GitHub Action. Used only for cases where multiple GitHub ID Token Repositories can be matched.",
						ElementType: types.StringType,
						Optional:    true,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(2),
							setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
						},
					},
					"workflow": schema.StringAttribute{
						Description: "The GitHub Workflow execution associated with the GitHub Action ID Token.",
						Optional:    true,
					},
					"workflows": schema.SetAttribute{
						Description: "The set of accepted GitHub ID Token Workflows which initiated the GitHub Action. Used only for cases where multiple GitHub ID Token Workflows can be matched.",
						ElementType: types.StringType,
						Optional:    true,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(2),
							setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
						},
					},
					"oidc_client_id": schema.StringAttribute{
						Description: "The OAuth Client ID value required for authenticating a GitHub Action.",
						Computed:    true,
					},
					"oidc_audience": schema.StringAttribute{
						Description: "The audience value required for the GitHub Action ID Token.",
						Computed:    true,
					},
				},
			},
			"gitlab_job": schema.SingleNestedAttribute{
				Description: "GitLab Job type Trust Provider configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"oidc_endpoint": schema.StringAttribute{
						Description: "The GitLab OIDC Endpoint used for validating GitLab Job generated ID Tokens. Default: `https://gitlab.com`.",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString("https://gitlab.com"),
					},
					"oidc_client_id": schema.StringAttribute{
						Description: "The OAuth Client ID value required for authenticating a GitLab Job.",
						Computed:    true,
					},
					"oidc_audience": schema.StringAttribute{
						Description: "The audience value required for the GitLab Job ID Token.",
						Computed:    true,
					},
					"namespace_path": schema.StringAttribute{
						Description: "The GitLab ID Token Namespace Path which initiated the GitLab Job.",
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"namespace_paths": schema.SetAttribute{
						Description: "The set of accepted GitLab ID Token Namespace Paths which initiated the GitLab Job. Used only for cases where multiple GitLab ID Token Namespace Paths can be matched.",
						ElementType: types.StringType,
						Optional:    true,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(2),
							setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
						},
					},
					"project_path": schema.StringAttribute{
						Description: "The GitLab ID Token Project Path which initiated the GitLab Job.",
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"project_paths": schema.SetAttribute{
						Description: "The set of accepted GitLab ID Token Project Paths which initiated the GitLab Job. Used only for cases where multiple GitLab ID Token Project Paths can be matched.",
						ElementType: types.StringType,
						Optional:    true,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(2),
							setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
						},
					},
					"ref_path": schema.StringAttribute{
						Description: "The GitLab ID Token Ref Path which initiated the GitLab Job.",
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"ref_paths": schema.SetAttribute{
						Description: "The set of accepted GitLab ID Token Ref Paths which initiated the GitLab Job. Used only for cases where multiple GitLab ID Token Ref Paths can be matched.",
						ElementType: types.StringType,
						Optional:    true,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(2),
							setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
						},
					},
					"subject": schema.StringAttribute{
						Description: "The GitLab ID Token Subject which initiated the GitLab Job.",
						Optional:    true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"subjects": schema.SetAttribute{
						Description: "The set of accepted GitLab ID Token Subjects which initiated the GitLab Job. Used only for cases where multiple GitLab ID Token Subjects can be matched.",
						ElementType: types.StringType,
						Optional:    true,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(2),
							setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
						},
					},
				},
			},
			"kerberos": schema.SingleNestedAttribute{
				Description: "Kerberos type Trust Provider configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"agent_controller_ids": schema.SetAttribute{
						Description: "Unique identifier for the Aembit Agent Controller to use for Signature verification.",
						Required:    true,
						ElementType: types.StringType,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(1),
						},
					},
					"principal": schema.StringAttribute{
						Description: "The Kerberos Principal of the authenticated Agent Proxy.",
						Optional:    true,
					},
					"principals": schema.SetAttribute{
						Description: "The set of accepted Kerberos Principals which initiated the authenticated Agent Proxy. Used only for cases where multiple Kerberos Principals can be matched.",
						ElementType: types.StringType,
						Optional:    true,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(2),
							setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
						},
					},
					"realm_domain": schema.StringAttribute{
						Description: "The Kerberos Realm or ActiveDirectory Domain of the authenticated Agent Proxy.",
						Optional:    true,
					},
					"realms_domains": schema.SetAttribute{
						Description: "The set of accepted Kerberos Realms or ActiveDirectory Domains which initiated the authenticated Agent Proxy. Used for cases where multiple Kerberos Realms or ActiveDirectory Domains can be matched.",
						ElementType: types.StringType,
						Optional:    true,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(2),
							setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
						},
					},
					"source_ip": schema.StringAttribute{
						Description: "The Source IP Address of the authenticated Agent Proxy.",
						Optional:    true,
					},
					"source_ips": schema.SetAttribute{
						Description: "The set of accepted Source IPs which initiated the authenticated Agent Proxy. Used only for cases where multiple Source IPs can be matched.",
						ElementType: types.StringType,
						Optional:    true,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(2),
							setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
						},
					},
				},
			},
			"kubernetes_service_account": schema.SingleNestedAttribute{
				Description: "Kubernetes Service Account type Trust Provider configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"issuer": schema.StringAttribute{
						Description: "The Issuer (`iss` claim) of the Kubernetes Service Account Token.",
						Optional:    true,
					},
					"issuers": schema.SetAttribute{
						Description: "The set of accepted Issuer values of the associated Kubernetes Service Account Token. Used only for cases where multiple Issuers can be matched.",
						ElementType: types.StringType,
						Optional:    true,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(2),
							setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
						},
					},
					"namespace": schema.StringAttribute{
						Description: "The Namespace of the Kubernetes Service Account Token.",
						Optional:    true,
					},
					"namespaces": schema.SetAttribute{
						Description: "The set of accepted Namespace values of the associated Kubernetes Service Account Token. Used only for cases where multiple Namespaces can be matched.",
						ElementType: types.StringType,
						Optional:    true,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(2),
							setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
						},
					},
					"pod_name": schema.StringAttribute{
						Description: "The Pod Name of the Kubernetes Service Account Token.",
						Optional:    true,
					},
					"pod_names": schema.SetAttribute{
						Description: "The set of accepted Pod Name values of the associated Kubernetes Service Account Token. Used only for cases where multiple Pod Names can be matched.",
						ElementType: types.StringType,
						Optional:    true,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(2),
							setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
						},
					},
					"service_account_name": schema.StringAttribute{
						Description: "The Service Account Name of the Kubernetes Service Account Token.",
						Optional:    true,
					},
					"service_account_names": schema.SetAttribute{
						Description: "The set of accepted Service Account Name values of the associated Kubernetes Service Account Token. Used only for cases where multiple Service Account Names can be matched.",
						ElementType: types.StringType,
						Optional:    true,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(2),
							setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
						},
					},
					"subject": schema.StringAttribute{
						Description: "The Subject (`sub` claim) of the Kubernetes Service Account Token.",
						Optional:    true,
					},
					"subjects": schema.SetAttribute{
						Description: "The set of accepted Subject values of the associated Kubernetes Service Account Token. Used only for cases where multiple Subjects can be matched.",
						ElementType: types.StringType,
						Optional:    true,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(2),
							setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
						},
					},
					"oidc_endpoint": schema.StringAttribute{
						Description: "The OIDC Endpoint from which Public Keys can be retrieved for verifying the signature of the Kubernetes Service Account Token.",
						Optional:    true,
						Validators: []validator.String{
							validators.OidcEndpointValidation(),
						},
					},
					"public_key": schema.StringAttribute{
						Description: "The Public Key that can be used to verify the signature of the Kubernetes Service Account Token.",
						Optional:    true,
					},
					"jwks": schema.StringAttribute{
						CustomType:  jsontypes.NormalizedType{},
						Description: "The JSON Web Key Set (JWKS) containing public keys used for signature verification.<br>**Note:** Only strictly valid JSON, with no trailing commas, will pass validation for this field.",
						Optional:    true,
						Computed:    true,
					},
					"symmetric_key": schema.StringAttribute{
						Description: "The Symmetric Key that can be used to verify the signature of the Kubernetes Service Account Token.",
						Optional:    true,
						Sensitive:   true,
						Validators: []validator.String{
							validators.Base64Validation(),
						},
					},
				},
			},
			"terraform_workspace": schema.SingleNestedAttribute{
				Description: "Terraform Workspace type Trust Provider configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"organization_id": schema.StringAttribute{
						Description: "The Organization ID of the calling Terraform Workspace.",
						Optional:    true,
					},
					"organization_ids": schema.SetAttribute{
						Description: "The set of accepted Organization ID values of the calling Terraform Workspace. Used only for cases where multiple Terraform ID Token Organization IDs can be matched.",
						ElementType: types.StringType,
						Optional:    true,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(2),
							setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
						},
					},
					"project_id": schema.StringAttribute{
						Description: "The Project ID of the calling Terraform Workspace.",
						Optional:    true,
					},
					"project_ids": schema.SetAttribute{
						Description: "The set of accepted Project ID values of the calling Terraform Workspace. Used only for cases where multiple Terraform ID Token Project IDs can be matched.",
						ElementType: types.StringType,
						Optional:    true,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(2),
							setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
						},
					},
					"workspace_id": schema.StringAttribute{
						Description: "The Workspace ID of the calling Terraform Workspace.",
						Optional:    true,
					},
					"workspace_ids": schema.SetAttribute{
						Description: "The set of accepted Workspace ID values of the calling Terraform Workspace. Used only for cases where multiple Terraform ID Token Workspace IDs can be matched.",
						ElementType: types.StringType,
						Optional:    true,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(2),
							setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
						},
					},
				},
			},
			"oidc_id_token": schema.SingleNestedAttribute{
				Description: "OIDC ID Token type Trust Provider configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"issuer": schema.StringAttribute{
						Description: "The Issuer (`iss` claim) of the OIDC ID Token.",
						Optional:    true,
					},
					"issuers": schema.SetAttribute{
						Description: "The set of accepted Issuer values of the associated OIDC ID Token. Used only for cases where multiple Issuers can be matched.",
						ElementType: types.StringType,
						Optional:    true,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(2),
							setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
						},
					},
					"subject": schema.StringAttribute{
						Description: "The Subject (`sub` claim) of the OIDC ID Token.",
						Optional:    true,
					},
					"subjects": schema.SetAttribute{
						Description: "The set of accepted Subject values of the associated OIDC ID Token. Used only for cases where multiple Subjects can be matched.",
						ElementType: types.StringType,
						Optional:    true,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(2),
							setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
						},
					},
					"audience": schema.StringAttribute{
						Description: "The Audience (`aud` claim) of the OIDC ID Token.",
						Optional:    true,
					},
					"audiences": schema.SetAttribute{
						Description: "The set of accepted Audience values of the associated OIDC ID Token. Used only for cases where multiple Audiences can be matched.",
						ElementType: types.StringType,
						Optional:    true,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(2),
							setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
						},
					},
					"oidc_endpoint": schema.StringAttribute{
						Description: "The OIDC Endpoint from which Public Keys can be retrieved for verifying the signature of the OIDC ID Token.",
						Optional:    true,
						Validators: []validator.String{
							validators.OidcEndpointValidation(),
						},
					},
					"public_key": schema.StringAttribute{
						Description: "The Public Key that can be used to verify the signature of the OIDC ID Token.",
						Optional:    true,
					},
					"jwks": schema.StringAttribute{
						CustomType:  jsontypes.NormalizedType{},
						Description: "The JSON Web Key Set (JWKS) containing public keys used for signature verification.<br>**Note:** Only strictly valid JSON, with no trailing commas, will pass validation for this field.",
						Optional:    true,
						Computed:    true,
					},
					"symmetric_key": schema.StringAttribute{
						Description: "The Symmetric Key that can be used to verify the signature of the OIDC ID Token.",
						Optional:    true,
						Sensitive:   true,
						Validators: []validator.String{
							validators.Base64Validation(),
						},
					},
				},
			},
		},
	}
}

// Configure validators to ensure that only one Trust Provider type is specified.
func (r *trustProviderResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.ExactlyOneOf(
			path.MatchRoot("aws_role"),
			path.MatchRoot("aws_metadata"),
			path.MatchRoot("azure_metadata"),
			path.MatchRoot("gcp_identity"),
			path.MatchRoot("github_action"),
			path.MatchRoot("gitlab_job"),
			path.MatchRoot("kerberos"),
			path.MatchRoot("kubernetes_service_account"),
			path.MatchRoot("terraform_workspace"),
			path.MatchRoot("oidc_id_token"),
		),
		// Ensure we don't have conflicting single and multiple match rule configurations (Azure Metadata)
		resourcevalidator.Conflicting(
			path.MatchRoot("azure_metadata").AtName("sku"),
			path.MatchRoot("azure_metadata").AtName("skus"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("azure_metadata").AtName("vm_id"),
			path.MatchRoot("azure_metadata").AtName("vm_ids"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("azure_metadata").AtName("subscription_id"),
			path.MatchRoot("azure_metadata").AtName("subscription_ids"),
		),
		// Ensure we don't have conflicting single and multiple match rule configurations (AWS Role)
		resourcevalidator.Conflicting(
			path.MatchRoot("aws_role").AtName("account_id"),
			path.MatchRoot("aws_role").AtName("account_ids"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("aws_role").AtName("assumed_role"),
			path.MatchRoot("aws_role").AtName("assumed_roles"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("aws_role").AtName("role_arn"),
			path.MatchRoot("aws_role").AtName("role_arns"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("aws_role").AtName("username"),
			path.MatchRoot("aws_role").AtName("usernames"),
		),
		// Ensure we don't have conflicting single and multiple match rule configurations (AWS Metadata)
		resourcevalidator.Conflicting(
			path.MatchRoot("aws_metadata").AtName("account_id"),
			path.MatchRoot("aws_metadata").AtName("account_ids"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("aws_metadata").AtName("availability_zone"),
			path.MatchRoot("aws_metadata").AtName("availability_zones"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("aws_metadata").AtName("instance_id"),
			path.MatchRoot("aws_metadata").AtName("instance_ids"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("aws_metadata").AtName("instance_type"),
			path.MatchRoot("aws_metadata").AtName("instance_types"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("aws_metadata").AtName("region"),
			path.MatchRoot("aws_metadata").AtName("regions"),
		),
		// Ensure we don't have conflicting single and multiple match rule configurations (Kerberos)
		resourcevalidator.Conflicting(
			path.MatchRoot("kerberos").AtName("principal"),
			path.MatchRoot("kerberos").AtName("principals"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("kerberos").AtName("realm_domain"),
			path.MatchRoot("kerberos").AtName("realms_domains"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("kerberos").AtName("source_ip"),
			path.MatchRoot("kerberos").AtName("source_ips"),
		),
		// Ensure we don't have conflicting single and multiple match rule configurations (Kubernetes Service Account)
		resourcevalidator.Conflicting(
			path.MatchRoot("kubernetes_service_account").AtName("issuer"),
			path.MatchRoot("kubernetes_service_account").AtName("issuers"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("kubernetes_service_account").AtName("namespace"),
			path.MatchRoot("kubernetes_service_account").AtName("namespaces"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("kubernetes_service_account").AtName("pod_name"),
			path.MatchRoot("kubernetes_service_account").AtName("pod_names"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("kubernetes_service_account").AtName("service_account_name"),
			path.MatchRoot("kubernetes_service_account").AtName("service_account_names"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("kubernetes_service_account").AtName("subject"),
			path.MatchRoot("kubernetes_service_account").AtName("subjects"),
		),
		// Ensure we don't have conflicting single and multiple match rule configurations (GCP Identity)
		resourcevalidator.Conflicting(
			path.MatchRoot("gcp_identity").AtName("email"),
			path.MatchRoot("gcp_identity").AtName("emails"),
		),
		// Ensure we don't have conflicting single and multiple match rule configurations (GitHub Action)
		resourcevalidator.Conflicting(
			path.MatchRoot("github_action").AtName("actor"),
			path.MatchRoot("github_action").AtName("actors"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("github_action").AtName("repository"),
			path.MatchRoot("github_action").AtName("repositories"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("github_action").AtName("workflow"),
			path.MatchRoot("github_action").AtName("workflows"),
		),
		// Ensure we don't have conflicting single and multiple match rule configurations (GitLab Job)
		resourcevalidator.Conflicting(
			path.MatchRoot("gitlab_job").AtName("namespace_path"),
			path.MatchRoot("gitlab_job").AtName("namespace_paths"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("gitlab_job").AtName("ref_path"),
			path.MatchRoot("gitlab_job").AtName("ref_paths"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("gitlab_job").AtName("project_path"),
			path.MatchRoot("gitlab_job").AtName("project_paths"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("gitlab_job").AtName("subject"),
			path.MatchRoot("gitlab_job").AtName("subjects"),
		),
		// Ensure we don't have conflicting single and multiple match rule configurations (Terraform Cloud)
		resourcevalidator.Conflicting(
			path.MatchRoot("terraform_workspace").AtName("organization_id"),
			path.MatchRoot("terraform_workspace").AtName("organization_ids"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("terraform_workspace").AtName("project_id"),
			path.MatchRoot("terraform_workspace").AtName("project_ids"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("terraform_workspace").AtName("workspace_id"),
			path.MatchRoot("terraform_workspace").AtName("workspace_ids"),
		),
		// Ensure we don't have conflicting single and multiple match rule configurations (OIDC ID Token)
		resourcevalidator.Conflicting(
			path.MatchRoot("oidc_id_token").AtName("issuer"),
			path.MatchRoot("oidc_id_token").AtName("issuers"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("oidc_id_token").AtName("subject"),
			path.MatchRoot("oidc_id_token").AtName("subjects"),
		),
		resourcevalidator.Conflicting(
			path.MatchRoot("oidc_id_token").AtName("audience"),
			path.MatchRoot("oidc_id_token").AtName("audiences"),
		),
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *trustProviderResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	// Retrieve values from plan
	var plan models.TrustProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	trust, err := convertTrustProviderModelToDTO(ctx, plan, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Trust Provider",
			err.Error(),
		)
		return
	}

	// Create new Trust Provider
	trustProvider, err := r.client.CreateTrustProvider(trust, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Trust Provider",
			"Could not create Trust Provider, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan = convertTrustProviderDTOToModel(
		ctx,
		*trustProvider,
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
func (r *trustProviderResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	// Get current state
	var state models.TrustProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed trust value from Aembit
	trustProvider, err, notFound := r.client.GetTrustProvider(state.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Error reading Aembit Trust Provider",
			"Could not read Aembit External ID from Terraform state "+state.ID.ValueString()+": "+err.Error(),
		)

		// If the resource is not found on Aembit Cloud, delete it locally
		if notFound {
			resp.State.RemoveResource(ctx)
		}
		return
	}

	state = convertTrustProviderDTOToModel(
		ctx,
		trustProvider,
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
func (r *trustProviderResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	// Get current state
	var state models.TrustProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract external ID from state
	externalID := state.ID.ValueString()

	// Retrieve values from plan
	var plan models.TrustProviderResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	trust, err := convertTrustProviderModelToDTO(ctx, plan, &externalID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Trust Provider",
			err.Error(),
		)
		return
	}

	// Update Trust Provider
	trustProvider, err := r.client.UpdateTrustProvider(trust, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Trust Provider",
			"Could not update Trust Provider, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	state = convertTrustProviderDTOToModel(
		ctx,
		*trustProvider,
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
func (r *trustProviderResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	// Retrieve values from state
	var state models.TrustProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check if Trust Provider is Active - if it is, disable it first
	if state.IsActive == types.BoolValue(true) {
		_, err := r.client.DisableTrustProvider(state.ID.ValueString(), nil)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error disabling Trust Provider",
				"Could not disable Trust Provider, unexpected error: "+err.Error(),
			)
			return
		}
	}

	// Delete existing Trust Provider
	_, err := r.client.DeleteTrustProvider(ctx, state.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Trust Provider",
			"Could not delete Trust Provider, unexpected error: "+err.Error(),
		)
		return
	}
}

// Imports an existing resource by passing externalId.
func (r *trustProviderResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	// Retrieve import externalId and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Model to DTO conversion methods.
func convertTrustProviderModelToDTO(
	ctx context.Context,
	model models.TrustProviderResourceModel,
	externalID *string,
) (aembit.TrustProviderDTO, error) {
	var trust aembit.TrustProviderDTO
	trust.EntityDTO = aembit.EntityDTO{
		Name:        model.Name.ValueString(),
		Description: model.Description.ValueString(),
		IsActive:    model.IsActive.ValueBool(),
	}
	if len(model.Tags.Elements()) > 0 {
		tagsMap := make(map[string]string)
		_ = model.Tags.ElementsAs(ctx, &tagsMap, true)

		for key, value := range tagsMap {
			trust.Tags = append(trust.Tags, aembit.TagDTO{
				Key:   key,
				Value: value,
			})
		}
	}
	if externalID != nil {
		trust.ExternalID = *externalID
	}

	var err error = nil
	// Transform the various Trust Provider types
	if model.AwsMetadata != nil {
		convertAwsMetadataModelToDTO(model, &trust)
	}
	if model.AwsRole != nil {
		convertAwsRoleModelToDTO(model, &trust)
	}
	if model.AzureMetadata != nil {
		convertAzureMetadataModelToDTO(model, &trust)
	}
	if model.GcpIdentity != nil {
		convertGcpIdentityModelToDTO(model, &trust)
	}
	if model.GitHubAction != nil {
		convertGitHubActionModelToDTO(model, &trust)
	}
	if model.GitLabJob != nil {
		convertGitLabJobModelToDTO(model, &trust)
	}
	if model.Kerberos != nil {
		convertKerberosModelToDTO(model, &trust)
	}
	if model.KubernetesService != nil {
		err = convertKubernetesModelToDTO(model, &trust)
	}
	if model.TerraformWorkspace != nil {
		convertTerraformModelToDTO(model, &trust)
	}
	if model.OidcIdToken != nil {
		err = convertOidcIdTokenTpModelToDTO(model, &trust)
	}

	return trust, err
}

func appendMatchRuleIfExists(
	matchRules []aembit.TrustProviderMatchRuleDTO,
	value basetypes.StringValue,
	attrName string,
) []aembit.TrustProviderMatchRuleDTO {
	if len(value.ValueString()) > 0 {
		return append(matchRules, aembit.TrustProviderMatchRuleDTO{
			Attribute: attrName, Value: value.ValueString(),
		})
	}
	return matchRules
}

func appendMatchRulesIfExists(
	matchRules []aembit.TrustProviderMatchRuleDTO,
	values []basetypes.StringValue,
	attrName string,
) []aembit.TrustProviderMatchRuleDTO {
	if len(values) > 0 {
		for _, value := range values {
			matchRules = appendMatchRuleIfExists(matchRules, value, attrName)
		}
	}
	return matchRules
}

func convertAzureMetadataModelToDTO(
	model models.TrustProviderResourceModel,
	dto *aembit.TrustProviderDTO,
) {
	dto.Provider = "AzureMetadataService"

	dto.MatchRules = make([]aembit.TrustProviderMatchRuleDTO, 0)
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.AzureMetadata.Sku, "AzureSku")
	dto.MatchRules = appendMatchRulesIfExists(dto.MatchRules, model.AzureMetadata.Skus, "AzureSku")
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.AzureMetadata.VMID, "AzureVmId")
	dto.MatchRules = appendMatchRulesIfExists(
		dto.MatchRules,
		model.AzureMetadata.VMIDs,
		"AzureVmId",
	)
	dto.MatchRules = appendMatchRuleIfExists(
		dto.MatchRules,
		model.AzureMetadata.SubscriptionID,
		"AzureSubscriptionId",
	)
	dto.MatchRules = appendMatchRulesIfExists(
		dto.MatchRules,
		model.AzureMetadata.SubscriptionIDs,
		"AzureSubscriptionId",
	)
}

func convertAwsRoleModelToDTO(
	model models.TrustProviderResourceModel,
	dto *aembit.TrustProviderDTO,
) {
	dto.Provider = "AWSRole"

	dto.MatchRules = make([]aembit.TrustProviderMatchRuleDTO, 0)
	dto.MatchRules = appendMatchRuleIfExists(
		dto.MatchRules,
		model.AwsRole.AccountID,
		"AwsAccountId",
	)
	dto.MatchRules = appendMatchRulesIfExists(
		dto.MatchRules,
		model.AwsRole.AccountIDs,
		"AwsAccountId",
	)
	dto.MatchRules = appendMatchRuleIfExists(
		dto.MatchRules,
		model.AwsRole.AssumedRole,
		"AwsAssumedRole",
	)
	dto.MatchRules = appendMatchRulesIfExists(
		dto.MatchRules,
		model.AwsRole.AssumedRoles,
		"AwsAssumedRole",
	)
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.AwsRole.RoleARN, "AwsRoleARN")
	dto.MatchRules = appendMatchRulesIfExists(dto.MatchRules, model.AwsRole.RoleARNs, "AwsRoleARN")
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.AwsRole.Username, "AwsUsername")
	dto.MatchRules = appendMatchRulesIfExists(
		dto.MatchRules,
		model.AwsRole.Usernames,
		"AwsUsername",
	)
}

func convertAwsMetadataModelToDTO(
	model models.TrustProviderResourceModel,
	dto *aembit.TrustProviderDTO,
) {
	dto.Provider = "AWSMetadataService"
	dto.Certificate = base64.StdEncoding.EncodeToString(
		[]byte(model.AwsMetadata.Certificate.ValueString()),
	)
	dto.PemType = "Certificate"

	dto.MatchRules = make([]aembit.TrustProviderMatchRuleDTO, 0)
	dto.MatchRules = appendMatchRuleIfExists(
		dto.MatchRules,
		model.AwsMetadata.AccountID,
		"AwsAccountId",
	)
	dto.MatchRules = appendMatchRulesIfExists(
		dto.MatchRules,
		model.AwsMetadata.AccountIDs,
		"AwsAccountId",
	)
	dto.MatchRules = appendMatchRuleIfExists(
		dto.MatchRules,
		model.AwsMetadata.Architecture,
		"AwsArchitecture",
	)
	dto.MatchRules = appendMatchRuleIfExists(
		dto.MatchRules,
		model.AwsMetadata.AvailabilityZone,
		"AwsAvailabilityZone",
	)
	dto.MatchRules = appendMatchRulesIfExists(
		dto.MatchRules,
		model.AwsMetadata.AvailabilityZones,
		"AwsAvailabilityZone",
	)
	dto.MatchRules = appendMatchRuleIfExists(
		dto.MatchRules,
		model.AwsMetadata.BillingProducts,
		"AwsBillingProducts",
	)
	dto.MatchRules = appendMatchRuleIfExists(
		dto.MatchRules,
		model.AwsMetadata.ImageID,
		"AwsImageId",
	)
	dto.MatchRules = appendMatchRuleIfExists(
		dto.MatchRules,
		model.AwsMetadata.InstanceID,
		"AwsInstanceId",
	)
	dto.MatchRules = appendMatchRulesIfExists(
		dto.MatchRules,
		model.AwsMetadata.InstanceIDs,
		"AwsInstanceId",
	)
	dto.MatchRules = appendMatchRuleIfExists(
		dto.MatchRules,
		model.AwsMetadata.InstanceType,
		"AwsInstanceType",
	)
	dto.MatchRules = appendMatchRulesIfExists(
		dto.MatchRules,
		model.AwsMetadata.InstanceTypes,
		"AwsInstanceType",
	)
	dto.MatchRules = appendMatchRuleIfExists(
		dto.MatchRules,
		model.AwsMetadata.KernelID,
		"AwsKernelId",
	)
	dto.MatchRules = appendMatchRuleIfExists(
		dto.MatchRules,
		model.AwsMetadata.MarketplaceProductCodes,
		"AwsMarketplaceProductCodes",
	)
	dto.MatchRules = appendMatchRuleIfExists(
		dto.MatchRules,
		model.AwsMetadata.PendingTime,
		"AwsPendingTime",
	)
	dto.MatchRules = appendMatchRuleIfExists(
		dto.MatchRules,
		model.AwsMetadata.PrivateIP,
		"AwsPrivateIp",
	)
	dto.MatchRules = appendMatchRuleIfExists(
		dto.MatchRules,
		model.AwsMetadata.RamdiskID,
		"AwsRamdiskId",
	)
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.AwsMetadata.Region, "AwsRegion")
	dto.MatchRules = appendMatchRulesIfExists(
		dto.MatchRules,
		model.AwsMetadata.Regions,
		"AwsRegion",
	)
	dto.MatchRules = appendMatchRuleIfExists(
		dto.MatchRules,
		model.AwsMetadata.Version,
		"AwsVersion",
	)
}

func convertGcpIdentityModelToDTO(
	model models.TrustProviderResourceModel,
	dto *aembit.TrustProviderDTO,
) {
	dto.Provider = "GcpIdentityToken"

	dto.MatchRules = make([]aembit.TrustProviderMatchRuleDTO, 0)
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.GcpIdentity.EMail, "Email")
	dto.MatchRules = appendMatchRulesIfExists(dto.MatchRules, model.GcpIdentity.EMails, "Email")
}

func convertGitHubActionModelToDTO(
	model models.TrustProviderResourceModel,
	dto *aembit.TrustProviderDTO,
) {
	dto.Provider = "GitHubIdentityToken"

	dto.MatchRules = make([]aembit.TrustProviderMatchRuleDTO, 0)
	dto.MatchRules = appendMatchRuleIfExists(
		dto.MatchRules,
		model.GitHubAction.Actor,
		"GithubActor",
	)
	dto.MatchRules = appendMatchRulesIfExists(
		dto.MatchRules,
		model.GitHubAction.Actors,
		"GithubActor",
	)
	dto.MatchRules = appendMatchRuleIfExists(
		dto.MatchRules,
		model.GitHubAction.Repository,
		"GithubRepository",
	)
	dto.MatchRules = appendMatchRulesIfExists(
		dto.MatchRules,
		model.GitHubAction.Repositories,
		"GithubRepository",
	)
	dto.MatchRules = appendMatchRuleIfExists(
		dto.MatchRules,
		model.GitHubAction.Workflow,
		"GithubWorkflow",
	)
	dto.MatchRules = appendMatchRulesIfExists(
		dto.MatchRules,
		model.GitHubAction.Workflows,
		"GithubWorkflow",
	)
}

func convertGitLabJobModelToDTO(
	model models.TrustProviderResourceModel,
	dto *aembit.TrustProviderDTO,
) {
	dto.Provider = "GitLabIdentityToken"

	dto.OidcUrl = model.GitLabJob.OIDCEndpoint.ValueString()
	dto.MatchRules = make([]aembit.TrustProviderMatchRuleDTO, 0)
	dto.MatchRules = appendMatchRuleIfExists(
		dto.MatchRules,
		model.GitLabJob.Subject,
		"GitLabSubject",
	)
	dto.MatchRules = appendMatchRulesIfExists(
		dto.MatchRules,
		model.GitLabJob.Subjects,
		"GitLabSubject",
	)
	dto.MatchRules = appendMatchRuleIfExists(
		dto.MatchRules,
		model.GitLabJob.ProjectPath,
		"GitLabProjectPath",
	)
	dto.MatchRules = appendMatchRulesIfExists(
		dto.MatchRules,
		model.GitLabJob.ProjectPaths,
		"GitLabProjectPath",
	)
	dto.MatchRules = appendMatchRuleIfExists(
		dto.MatchRules,
		model.GitLabJob.NamespacePath,
		"GitLabNamespacePath",
	)
	dto.MatchRules = appendMatchRulesIfExists(
		dto.MatchRules,
		model.GitLabJob.NamespacePaths,
		"GitLabNamespacePath",
	)
	dto.MatchRules = appendMatchRuleIfExists(
		dto.MatchRules,
		model.GitLabJob.RefPath,
		"GitLabRefPath",
	)
	dto.MatchRules = appendMatchRulesIfExists(
		dto.MatchRules,
		model.GitLabJob.RefPaths,
		"GitLabRefPath",
	)
}

func convertKerberosModelToDTO(
	model models.TrustProviderResourceModel,
	dto *aembit.TrustProviderDTO,
) {
	dto.Provider = "Kerberos"
	dto.AgentControllerIDs = make([]string, len(model.Kerberos.AgentControllerIDs))
	for i, controllerID := range model.Kerberos.AgentControllerIDs {
		dto.AgentControllerIDs[i] = controllerID.ValueString()
	}

	dto.MatchRules = make([]aembit.TrustProviderMatchRuleDTO, 0)
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.Kerberos.Principal, "Principal")
	dto.MatchRules = appendMatchRulesIfExists(
		dto.MatchRules,
		model.Kerberos.Principals,
		"Principal",
	)
	dto.MatchRules = appendMatchRuleIfExists(
		dto.MatchRules,
		model.Kerberos.RealmOrDomain,
		"RealmOrDomain",
	)
	dto.MatchRules = appendMatchRulesIfExists(
		dto.MatchRules,
		model.Kerberos.RealmsOrDomains,
		"RealmOrDomain",
	)
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.Kerberos.SourceIP, "SourceIp")
	dto.MatchRules = appendMatchRulesIfExists(dto.MatchRules, model.Kerberos.SourceIPs, "SourceIp")
}

func convertKubernetesModelToDTO(
	model models.TrustProviderResourceModel,
	dto *aembit.TrustProviderDTO,
) error {
	dto.Provider = "KubernetesServiceAccount"
	dto.Certificate = base64.StdEncoding.EncodeToString(
		[]byte(model.KubernetesService.PublicKey.ValueString()),
	)
	if len(dto.Certificate) > 0 {
		dto.PemType = "PublicKey"
	}
	dto.OidcUrl = model.KubernetesService.OIDCEndpoint.ValueString()

	err := convertJWKSModelToDto(model.KubernetesService.Jwks.ValueString(), dto)
	if err != nil {
		return err
	}

	dto.MatchRules = make([]aembit.TrustProviderMatchRuleDTO, 0)
	dto.MatchRules = appendMatchRuleIfExists(
		dto.MatchRules,
		model.KubernetesService.Issuer,
		"KubernetesIss",
	)
	dto.MatchRules = appendMatchRulesIfExists(
		dto.MatchRules,
		model.KubernetesService.Issuers,
		"KubernetesIss",
	)
	dto.MatchRules = appendMatchRuleIfExists(
		dto.MatchRules,
		model.KubernetesService.Namespace,
		"KubernetesIoNamespace",
	)
	dto.MatchRules = appendMatchRulesIfExists(
		dto.MatchRules,
		model.KubernetesService.Namespaces,
		"KubernetesIoNamespace",
	)
	dto.MatchRules = appendMatchRuleIfExists(
		dto.MatchRules,
		model.KubernetesService.PodName,
		"KubernetesIoPodName",
	)
	dto.MatchRules = appendMatchRulesIfExists(
		dto.MatchRules,
		model.KubernetesService.PodNames,
		"KubernetesIoPodName",
	)
	dto.MatchRules = appendMatchRuleIfExists(
		dto.MatchRules,
		model.KubernetesService.ServiceAccountName,
		"KubernetesIoServiceAccountName",
	)
	dto.MatchRules = appendMatchRulesIfExists(
		dto.MatchRules,
		model.KubernetesService.ServiceAccountNames,
		"KubernetesIoServiceAccountName",
	)
	dto.MatchRules = appendMatchRuleIfExists(
		dto.MatchRules,
		model.KubernetesService.Subject,
		"KubernetesSub",
	)
	dto.MatchRules = appendMatchRulesIfExists(
		dto.MatchRules,
		model.KubernetesService.Subjects,
		"KubernetesSub",
	)

	return nil
}

func convertOidcIdTokenTpModelToDTO(
	model models.TrustProviderResourceModel,
	dto *aembit.TrustProviderDTO,
) error {
	dto.Provider = "OidcIdToken"
	dto.Certificate = base64.StdEncoding.EncodeToString(
		[]byte(model.OidcIdToken.PublicKey.ValueString()),
	)
	if len(dto.Certificate) > 0 {
		dto.PemType = "PublicKey"
	}
	dto.OidcUrl = model.OidcIdToken.OIDCEndpoint.ValueString()
	err := convertJWKSModelToDto(model.OidcIdToken.Jwks.ValueString(), dto)
	if err != nil {
		return err
	}

	dto.MatchRules = make([]aembit.TrustProviderMatchRuleDTO, 0)
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.OidcIdToken.Issuer, "OidcIssuer")
	dto.MatchRules = appendMatchRulesIfExists(
		dto.MatchRules,
		model.OidcIdToken.Issuers,
		"OidcIssuer",
	)
	dto.MatchRules = appendMatchRuleIfExists(
		dto.MatchRules,
		model.OidcIdToken.Subject,
		"OidcSubject",
	)
	dto.MatchRules = appendMatchRulesIfExists(
		dto.MatchRules,
		model.OidcIdToken.Subjects,
		"OidcSubject",
	)
	dto.MatchRules = appendMatchRuleIfExists(
		dto.MatchRules,
		model.OidcIdToken.Audience,
		"OidcAudience",
	)
	dto.MatchRules = appendMatchRulesIfExists(
		dto.MatchRules,
		model.OidcIdToken.Audiences,
		"OidcAudience",
	)

	return nil
}

func convertJWKSModelToDto(jwksJson string, dto *aembit.TrustProviderDTO) error {
	if jwksJson == "" {
		return nil
	}

	// convert jwksJson to model
	var jwks aembit.JsonWebKeysDTO
	err := json.Unmarshal([]byte(jwksJson), &jwks)
	if err != nil {
		return fmt.Errorf("JWKS content is not valid")
	}

	if len(jwks.Keys) == 0 {
		return fmt.Errorf("JWKS does not have any keys")
	}

	for _, key := range jwks.Keys {
		if key.Kty == "" {
			return fmt.Errorf("kty (Key Type) must be present in a JWK")
		}

		if key.Kty == "RSA" {
			if key.E == "" || key.N == "" {
				return fmt.Errorf("JWKS key does not have RSA required fields: e, n")
			}
		}

		if key.Kty == "EC" {
			if key.Crv == "" || key.X == "" || key.Y == "" {
				return fmt.Errorf("JWKS key does not have ECDSA required fields: x, y, crv")
			}
		}
	}

	// normalizing JSON
	// var buf bytes.Buffer
	// json.Compact(&buf, []byte(jwksJson))
	// dto.Jwks = buf.String()
	dto.Jwks = jwksJson
	return nil
}

func convertTerraformModelToDTO(
	model models.TrustProviderResourceModel,
	dto *aembit.TrustProviderDTO,
) {
	dto.Provider = "TerraformIdentityToken"

	dto.MatchRules = make([]aembit.TrustProviderMatchRuleDTO, 0)
	dto.MatchRules = appendMatchRuleIfExists(
		dto.MatchRules,
		model.TerraformWorkspace.OrganizationID,
		"TerraformOrganizationId",
	)
	dto.MatchRules = appendMatchRulesIfExists(
		dto.MatchRules,
		model.TerraformWorkspace.OrganizationIDs,
		"TerraformOrganizationId",
	)
	dto.MatchRules = appendMatchRuleIfExists(
		dto.MatchRules,
		model.TerraformWorkspace.ProjectID,
		"TerraformProjectId",
	)
	dto.MatchRules = appendMatchRulesIfExists(
		dto.MatchRules,
		model.TerraformWorkspace.ProjectIDs,
		"TerraformProjectId",
	)
	dto.MatchRules = appendMatchRuleIfExists(
		dto.MatchRules,
		model.TerraformWorkspace.WorkspaceID,
		"TerraformWorkspaceId",
	)
	dto.MatchRules = appendMatchRulesIfExists(
		dto.MatchRules,
		model.TerraformWorkspace.WorkspaceIDs,
		"TerraformWorkspaceId",
	)
}

// DTO to Model conversion methods.
func convertTrustProviderDTOToModel(
	ctx context.Context,
	dto aembit.TrustProviderDTO,
	state models.TrustProviderResourceModel,
	tenant, stackDomain string,
) models.TrustProviderResourceModel {
	var model models.TrustProviderResourceModel
	model.ID = types.StringValue(dto.ExternalID)
	model.Name = types.StringValue(dto.Name)
	model.Description = types.StringValue(dto.Description)
	model.IsActive = types.BoolValue(dto.IsActive)
	model.Tags = newTagsModel(ctx, dto.Tags)

	switch dto.Provider {
	case "AWSRole":
		model.AwsRole = convertAwsRoleDTOToModel(dto)
	case "AWSMetadataService":
		model.AwsMetadata = convertAwsMetadataDTOToModel(dto)
	case "AzureMetadataService":
		model.AzureMetadata = convertAzureMetadataDTOToModel(dto)
	case "GcpIdentityToken":
		model.GcpIdentity = convertGcpIdentityDTOToModel(dto)
	case "GitHubIdentityToken":
		model.GitHubAction = convertGitHubActionDTOToModel(dto)
	case "GitLabIdentityToken":
		model.GitLabJob = convertGitLabJobDTOToModel(dto, tenant, stackDomain)
	case "Kerberos":
		model.Kerberos = convertKerberosDTOToModel(dto)
	case "KubernetesServiceAccount":
		model.KubernetesService = convertKubernetesDTOToModel(dto, state)
	case "OidcIdToken":
		model.OidcIdToken = convertOidcIdTokenTpDTOToModel(dto, state)
	case "TerraformIdentityToken":
		model.TerraformWorkspace = convertTerraformDTOToModel(dto)
	}

	return model
}

func convertAzureMetadataDTOToModel(
	dto aembit.TrustProviderDTO,
) *models.TrustProviderAzureMetadataModel {
	model := &models.TrustProviderAzureMetadataModel{
		Sku:            types.StringNull(),
		VMID:           types.StringNull(),
		SubscriptionID: types.StringNull(),
	}

	if slices.ContainsFunc(dto.MatchRules, matchRuleAttributeFunc("AzureSku")) {
		model.Sku, model.Skus = extractMatchRules(dto.MatchRules, "AzureSku")
	}
	if slices.ContainsFunc(dto.MatchRules, matchRuleAttributeFunc("AzureVmId")) {
		model.VMID, model.VMIDs = extractMatchRules(dto.MatchRules, "AzureVmId")
	}
	if slices.ContainsFunc(dto.MatchRules, matchRuleAttributeFunc("AzureSubscriptionId")) {
		model.SubscriptionID, model.SubscriptionIDs = extractMatchRules(
			dto.MatchRules,
			"AzureSubscriptionId",
		)
	}
	return model
}

func convertAwsRoleDTOToModel(dto aembit.TrustProviderDTO) *models.TrustProviderAwsRoleModel {
	model := &models.TrustProviderAwsRoleModel{
		AccountID:   types.StringNull(),
		AssumedRole: types.StringNull(),
		RoleARN:     types.StringNull(),
		Username:    types.StringNull(),
	}

	if slices.ContainsFunc(dto.MatchRules, matchRuleAttributeFunc("AwsAccountId")) {
		model.AccountID, model.AccountIDs = extractMatchRules(dto.MatchRules, "AwsAccountId")
	}
	if slices.ContainsFunc(dto.MatchRules, matchRuleAttributeFunc("AwsAssumedRole")) {
		model.AssumedRole, model.AssumedRoles = extractMatchRules(dto.MatchRules, "AwsAssumedRole")
	}
	if slices.ContainsFunc(dto.MatchRules, matchRuleAttributeFunc("AwsRoleARN")) {
		model.RoleARN, model.RoleARNs = extractMatchRules(dto.MatchRules, "AwsRoleARN")
	}
	if slices.ContainsFunc(dto.MatchRules, matchRuleAttributeFunc("AwsUsername")) {
		model.Username, model.Usernames = extractMatchRules(dto.MatchRules, "AwsUsername")
	}
	return model
}

func convertAwsMetadataDTOToModel(
	dto aembit.TrustProviderDTO,
) *models.TrustProviderAwsMetadataModel {
	decodedCert, _ := base64.StdEncoding.DecodeString(dto.Certificate)

	model := &models.TrustProviderAwsMetadataModel{
		Certificate:             types.StringValue(string(decodedCert)),
		AccountID:               types.StringNull(),
		Architecture:            types.StringNull(),
		AvailabilityZone:        types.StringNull(),
		BillingProducts:         types.StringNull(),
		ImageID:                 types.StringNull(),
		InstanceID:              types.StringNull(),
		InstanceType:            types.StringNull(),
		KernelID:                types.StringNull(),
		MarketplaceProductCodes: types.StringNull(),
		PendingTime:             types.StringNull(),
		PrivateIP:               types.StringNull(),
		RamdiskID:               types.StringNull(),
		Region:                  types.StringNull(),
		Version:                 types.StringNull(),
	}

	if slices.ContainsFunc(dto.MatchRules, matchRuleAttributeFunc("AwsAccountId")) {
		model.AccountID, model.AccountIDs = extractMatchRules(dto.MatchRules, "AwsAccountId")
	}
	if slices.ContainsFunc(dto.MatchRules, matchRuleAttributeFunc("AwsArchitecture")) {
		model.Architecture, _ = extractMatchRules(dto.MatchRules, "AwsArchitecture")
	}
	if slices.ContainsFunc(dto.MatchRules, matchRuleAttributeFunc("AwsAvailabilityZone")) {
		model.AvailabilityZone, model.AvailabilityZones = extractMatchRules(
			dto.MatchRules,
			"AwsAvailabilityZone",
		)
	}
	if slices.ContainsFunc(dto.MatchRules, matchRuleAttributeFunc("AwsBillingProducts")) {
		model.BillingProducts, _ = extractMatchRules(dto.MatchRules, "AwsBillingProducts")
	}
	if slices.ContainsFunc(dto.MatchRules, matchRuleAttributeFunc("AwsImageId")) {
		model.ImageID, _ = extractMatchRules(dto.MatchRules, "AwsImageId")
	}
	if slices.ContainsFunc(dto.MatchRules, matchRuleAttributeFunc("AwsInstanceId")) {
		model.InstanceID, model.InstanceIDs = extractMatchRules(dto.MatchRules, "AwsInstanceId")
	}
	if slices.ContainsFunc(dto.MatchRules, matchRuleAttributeFunc("AwsInstanceType")) {
		model.InstanceType, model.InstanceTypes = extractMatchRules(
			dto.MatchRules,
			"AwsInstanceType",
		)
	}
	if slices.ContainsFunc(dto.MatchRules, matchRuleAttributeFunc("AwsKernelId")) {
		model.KernelID, _ = extractMatchRules(dto.MatchRules, "AwsKernelId")
	}
	if slices.ContainsFunc(dto.MatchRules, matchRuleAttributeFunc("AwsMarketplaceProductCodes")) {
		model.MarketplaceProductCodes, _ = extractMatchRules(
			dto.MatchRules,
			"AwsMarketplaceProductCodes",
		)
	}
	if slices.ContainsFunc(dto.MatchRules, matchRuleAttributeFunc("AwsPendingTime")) {
		model.PendingTime, _ = extractMatchRules(dto.MatchRules, "AwsPendingTime")
	}
	if slices.ContainsFunc(dto.MatchRules, matchRuleAttributeFunc("AwsPrivateIp")) {
		model.PrivateIP, _ = extractMatchRules(dto.MatchRules, "AwsPrivateIp")
	}
	if slices.ContainsFunc(dto.MatchRules, matchRuleAttributeFunc("AwsRamdiskId")) {
		model.RamdiskID, _ = extractMatchRules(dto.MatchRules, "AwsRamdiskId")
	}
	if slices.ContainsFunc(dto.MatchRules, matchRuleAttributeFunc("AwsRegion")) {
		model.Region, model.Regions = extractMatchRules(dto.MatchRules, "AwsRegion")
	}
	if slices.ContainsFunc(dto.MatchRules, matchRuleAttributeFunc("AwsVersion")) {
		model.Version, _ = extractMatchRules(dto.MatchRules, "AwsVersion")
	}
	return model
}

func convertKerberosDTOToModel(dto aembit.TrustProviderDTO) *models.TrustProviderKerberosModel {
	model := &models.TrustProviderKerberosModel{
		Principal:     types.StringNull(),
		RealmOrDomain: types.StringNull(),
		SourceIP:      types.StringNull(),
	}
	model.AgentControllerIDs = make([]types.String, len(dto.AgentControllerIDs))
	for i, controllerID := range dto.AgentControllerIDs {
		model.AgentControllerIDs[i] = types.StringValue(controllerID)
	}

	if slices.ContainsFunc(dto.MatchRules, matchRuleAttributeFunc("Principal")) {
		model.Principal, model.Principals = extractMatchRules(dto.MatchRules, "Principal")
	}
	if slices.ContainsFunc(dto.MatchRules, matchRuleAttributeFunc("Realm")) {
		model.RealmOrDomain, model.RealmsOrDomains = extractMatchRules(dto.MatchRules, "Realm")
	}
	if slices.ContainsFunc(dto.MatchRules, matchRuleAttributeFunc("SourceIp")) {
		model.SourceIP, model.SourceIPs = extractMatchRules(dto.MatchRules, "SourceIp")
	}
	return model
}

func convertKubernetesDTOToModel(
	dto aembit.TrustProviderDTO,
	state models.TrustProviderResourceModel,
) *models.TrustProviderKubernetesModel {
	decodedKey, _ := base64.StdEncoding.DecodeString(dto.Certificate)

	model := &models.TrustProviderKubernetesModel{
		Issuer:             types.StringNull(),
		Namespace:          types.StringNull(),
		PodName:            types.StringNull(),
		ServiceAccountName: types.StringNull(),
		Subject:            types.StringNull(),
		PublicKey:          types.StringNull(),
		OIDCEndpoint:       types.StringNull(),
		Jwks:               jsontypes.NewNormalizedNull(),
		SymmetricKey:       types.StringNull(),
	}

	if len(dto.Certificate) > 0 {
		model.PublicKey = types.StringValue(string(decodedKey))
	} else if len(dto.OidcUrl) > 0 {
		model.OIDCEndpoint = types.StringValue(dto.OidcUrl)
	} else if len(dto.Jwks) > 0 {
		model.Jwks = jsontypes.NewNormalizedValue(dto.Jwks)
	} else {
		model.SymmetricKey = state.KubernetesService.SymmetricKey
	}

	if slices.ContainsFunc(dto.MatchRules, matchRuleAttributeFunc("KubernetesIss")) {
		model.Issuer, model.Issuers = extractMatchRules(dto.MatchRules, "KubernetesIss")
	}
	if slices.ContainsFunc(dto.MatchRules, matchRuleAttributeFunc("KubernetesIoNamespace")) {
		model.Namespace, model.Namespaces = extractMatchRules(
			dto.MatchRules,
			"KubernetesIoNamespace",
		)
	}
	if slices.ContainsFunc(dto.MatchRules, matchRuleAttributeFunc("KubernetesIoPodName")) {
		model.PodName, model.PodNames = extractMatchRules(dto.MatchRules, "KubernetesIoPodName")
	}
	if slices.ContainsFunc(
		dto.MatchRules,
		matchRuleAttributeFunc("KubernetesIoServiceAccountName"),
	) {
		model.ServiceAccountName, model.ServiceAccountNames = extractMatchRules(
			dto.MatchRules,
			"KubernetesIoServiceAccountName",
		)
	}
	if slices.ContainsFunc(dto.MatchRules, matchRuleAttributeFunc("KubernetesSub")) {
		model.Subject, model.Subjects = extractMatchRules(dto.MatchRules, "KubernetesSub")
	}
	return model
}

func convertOidcIdTokenTpDTOToModel(
	dto aembit.TrustProviderDTO,
	state models.TrustProviderResourceModel,
) *models.TrustProviderOidcIdTokenModel {
	decodedKey, _ := base64.StdEncoding.DecodeString(dto.Certificate)

	model := &models.TrustProviderOidcIdTokenModel{
		Issuer:       types.StringNull(),
		Subject:      types.StringNull(),
		Audience:     types.StringNull(),
		PublicKey:    types.StringNull(),
		OIDCEndpoint: types.StringNull(),
		Jwks:         jsontypes.NewNormalizedNull(),
		SymmetricKey: types.StringNull(),
	}
	if len(dto.Certificate) > 0 {
		model.PublicKey = types.StringValue(string(decodedKey))
	} else if len(dto.OidcUrl) > 0 {
		model.OIDCEndpoint = types.StringValue(dto.OidcUrl)
	} else if len(dto.Jwks) > 0 {
		model.Jwks = jsontypes.NewNormalizedValue(dto.Jwks)
	} else {
		model.SymmetricKey = state.OidcIdToken.SymmetricKey
	}

	if slices.ContainsFunc(dto.MatchRules, matchRuleAttributeFunc("OidcIssuer")) {
		model.Issuer, model.Issuers = extractMatchRules(dto.MatchRules, "OidcIssuer")
	}
	if slices.ContainsFunc(dto.MatchRules, matchRuleAttributeFunc("OidcSubject")) {
		model.Subject, model.Subjects = extractMatchRules(dto.MatchRules, "OidcSubject")
	}
	if slices.ContainsFunc(dto.MatchRules, matchRuleAttributeFunc("OidcAudience")) {
		model.Audience, model.Audiences = extractMatchRules(dto.MatchRules, "OidcAudience")
	}
	return model
}

func convertGcpIdentityDTOToModel(
	dto aembit.TrustProviderDTO,
) *models.TrustProviderGcpIdentityModel {
	model := &models.TrustProviderGcpIdentityModel{
		EMail: types.StringNull(),
	}

	if slices.ContainsFunc(dto.MatchRules, matchRuleAttributeFunc("Email")) {
		model.EMail, model.EMails = extractMatchRules(dto.MatchRules, "Email")
	}
	return model
}

func convertGitHubActionDTOToModel(
	dto aembit.TrustProviderDTO,
) *models.TrustProviderGitHubActionModel {
	model := &models.TrustProviderGitHubActionModel{
		Actor:      types.StringNull(),
		Repository: types.StringNull(),
		Workflow:   types.StringNull(),
	}

	if slices.ContainsFunc(dto.MatchRules, matchRuleAttributeFunc("GithubActor")) {
		model.Actor, model.Actors = extractMatchRules(dto.MatchRules, "GithubActor")
	}
	if slices.ContainsFunc(dto.MatchRules, matchRuleAttributeFunc("GithubRepository")) {
		model.Repository, model.Repositories = extractMatchRules(dto.MatchRules, "GithubRepository")
	}
	if slices.ContainsFunc(dto.MatchRules, matchRuleAttributeFunc("GithubWorkflow")) {
		model.Workflow, model.Workflows = extractMatchRules(dto.MatchRules, "GithubWorkflow")
	}
	return model
}

func convertGitLabJobDTOToModel(
	dto aembit.TrustProviderDTO,
	tenant, stackDomain string,
) *models.TrustProviderGitLabJobModel {
	stackDomain = strings.ToLower(stackDomain) // Force the stack/domain to be lowercase
	stack := strings.Split(stackDomain, ".")[0]
	model := &models.TrustProviderGitLabJobModel{
		OIDCEndpoint:  types.StringValue(dto.OidcUrl),
		Subject:       types.StringNull(),
		ProjectPath:   types.StringNull(),
		NamespacePath: types.StringNull(),
		RefPath:       types.StringNull(),
		OIDCClientID: types.StringValue(
			fmt.Sprintf("aembit:%s:%s:identity:gitlab_idtoken:%s", stack, tenant, dto.ExternalID),
		),
		OIDCAudience: types.StringValue(fmt.Sprintf("https://%s.id.%s", tenant, stackDomain)),
	}

	if slices.ContainsFunc(dto.MatchRules, matchRuleAttributeFunc("GitLabSubject")) {
		model.Subject, model.Subjects = extractMatchRules(dto.MatchRules, "GitLabSubject")
	}
	if slices.ContainsFunc(dto.MatchRules, matchRuleAttributeFunc("GitLabProjectPath")) {
		model.ProjectPath, model.ProjectPaths = extractMatchRules(
			dto.MatchRules,
			"GitLabProjectPath",
		)
	}
	if slices.ContainsFunc(dto.MatchRules, matchRuleAttributeFunc("GitLabNamespacePath")) {
		model.NamespacePath, model.NamespacePaths = extractMatchRules(
			dto.MatchRules,
			"GitLabNamespacePath",
		)
	}
	if slices.ContainsFunc(dto.MatchRules, matchRuleAttributeFunc("GitLabRefPath")) {
		model.RefPath, model.RefPaths = extractMatchRules(dto.MatchRules, "GitLabRefPath")
	}
	return model
}

func convertTerraformDTOToModel(dto aembit.TrustProviderDTO) *models.TrustProviderTerraformModel {
	model := &models.TrustProviderTerraformModel{
		OrganizationID: types.StringNull(),
		ProjectID:      types.StringNull(),
		WorkspaceID:    types.StringNull(),
	}

	if slices.ContainsFunc(dto.MatchRules, matchRuleAttributeFunc("TerraformOrganizationId")) {
		model.OrganizationID, model.OrganizationIDs = extractMatchRules(
			dto.MatchRules,
			"TerraformOrganizationId",
		)
	}
	if slices.ContainsFunc(dto.MatchRules, matchRuleAttributeFunc("TerraformProjectId")) {
		model.ProjectID, model.ProjectIDs = extractMatchRules(dto.MatchRules, "TerraformProjectId")
	}
	if slices.ContainsFunc(dto.MatchRules, matchRuleAttributeFunc("TerraformWorkspaceId")) {
		model.WorkspaceID, model.WorkspaceIDs = extractMatchRules(
			dto.MatchRules,
			"TerraformWorkspaceId",
		)
	}
	return model
}

// matchRuleOccurrences returns the number of match rules with a specific attributeName.
func matchRuleOccurrences(matchRules []aembit.TrustProviderMatchRuleDTO, attributeName string) int {
	count := 0
	for _, rule := range matchRules {
		if rule.Attribute == attributeName {
			count++
		}
	}
	return count
}

// matchRuleAttributeFunc provides a function that can be used by the slices.ContainsFunc utility to determine if a slice of TrustProviderMatchRuleDTO contains a specific attribute.
func matchRuleAttributeFunc(attributeName string) func(aembit.TrustProviderMatchRuleDTO) bool {
	return func(matchRule aembit.TrustProviderMatchRuleDTO) bool {
		return matchRule.Attribute == attributeName
	}
}

// extractMatchRules pulls out the match rules with the specified attribute name.
func extractMatchRules(
	matchRules []aembit.TrustProviderMatchRuleDTO,
	attributeName string,
) (types.String, []types.String) {
	singleValue := types.StringNull()
	var multiValue []types.String = nil
	multipleMatchRules := matchRuleOccurrences(matchRules, attributeName) > 1

	for _, rule := range matchRules {
		if rule.Attribute == attributeName {
			if multipleMatchRules {
				multiValue = append(multiValue, types.StringValue(rule.Value))
			} else {
				singleValue = types.StringValue(rule.Value)
			}
		}
	}
	return singleValue, multiValue
}

// normalizeJSON takes raw JSON (string) and returns canonicalized version (minified)
// func normalizeJSON(input string) (string, error) {
// 	var buf bytes.Buffer
// 	err := json.Compact(&buf, []byte(input))
// 	if err != nil {
// 		return "", fmt.Errorf("failed to normalize JSON: %w", err)
// 	}
// 	return buf.String(), nil
// }

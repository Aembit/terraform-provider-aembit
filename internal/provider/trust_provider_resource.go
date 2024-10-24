package provider

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"aembit.io/aembit"
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
func (r *trustProviderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_trust_provider"
}

// Configure adds the provider configured client to the resource.
func (r *trustProviderResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = resourceConfigure(req, resp)
}

// Schema defines the schema for the resource.
func (r *trustProviderResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "**Note:** One and only one nested schema (e.g. `aws_metadata`) must be provided for the Trust Provider to be configured.",
		Attributes: map[string]schema.Attribute{
			// ID field is required for Terraform Framework acceptance testing.
			"id": schema.StringAttribute{
				Description: "Unique identifier of the Trust Provider.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name for the Trust Provider.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
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
					"vm_id": schema.StringAttribute{
						Description: "Unique identifier for the Virtual Machine.",
						Optional:    true,
					},
					"subscription_id": schema.StringAttribute{
						Description: "Azure subscription for the Virtual Machine.",
						Optional:    true,
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
					"assumed_role": schema.StringAttribute{
						Description: "The Name of the AWS IAM Role which is running the Client Workload.",
						Optional:    true,
					},
					"role_arn": schema.StringAttribute{
						Description: "The ARN of the AWS IAM Role which is running the Client Workload.",
						Optional:    true,
					},
					"username": schema.StringAttribute{
						Description: "The UsernID of the AWS IAM Account which is running the Client Workload (not commonly used).",
						Optional:    true,
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
					"architecture": schema.StringAttribute{
						Description: "The architecture of the AMI used to launch the instance (i386 | x86_64 | arm64).",
						Optional:    true,
					},
					"availability_zone": schema.StringAttribute{
						Description: "The Availability Zone in which the instance is running.",
						Optional:    true,
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
					"instance_type": schema.StringAttribute{
						Description: "The instance type of the instance.",
						Optional:    true,
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
						Description: "The Email of the GCP Service Account used by the associated GCP resource.",
						Optional:    true,
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
					"repository": schema.StringAttribute{
						Description: "The GitHub Repository associated with the GitHub Action ID Token.",
						Optional:    true,
					},
					"workflow": schema.StringAttribute{
						Description: "The GitHub Workflow execution associated with the GitHub Action ID Token.",
						Optional:    true,
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
					"realm": schema.StringAttribute{
						Description: "The Kerberos Realm of the authenticated Agent Proxy.",
						Optional:    true,
					},
					"source_ip": schema.StringAttribute{
						Description: "The Source IP Address of the authenticated Agent Proxy.",
						Optional:    true,
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
					"namespace": schema.StringAttribute{
						Description: "The Namespace of the Kubernetes Service Account Token.",
						Optional:    true,
					},
					"pod_name": schema.StringAttribute{
						Description: "The Pod Name of the Kubernetes Service Account Token.",
						Optional:    true,
					},
					"service_account_name": schema.StringAttribute{
						Description: "The Service Account Name of the Kubernetes Service Account Token.",
						Optional:    true,
					},
					"subject": schema.StringAttribute{
						Description: "The Subject (`sub` claim) of the Kubernetes Service Account Token.",
						Optional:    true,
					},
					"oidc_endpoint": schema.StringAttribute{
						Description: "The OIDC Endpoint from which Public Keys can be retrieved for verifying the signature of the Kubernetes Service Account Token.",
						Optional:    true,
					},
					"public_key": schema.StringAttribute{
						Description: "The Public Key that can be used to verify the signature of the Kubernetes Service Account Token.",
						Optional:    true,
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
					"project_id": schema.StringAttribute{
						Description: "The Project ID of the calling Terraform Workspace.",
						Optional:    true,
					},
					"workspace_id": schema.StringAttribute{
						Description: "The Workspace ID of the calling Terraform Workspace.",
						Optional:    true,
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
		),
		// For the GitLab Job type, ensure we don't have conflicting single and multiple match rule configurations
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
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *trustProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan trustProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var trust aembit.TrustProviderDTO = convertTrustProviderModelToDTO(ctx, plan, nil)

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
	plan = convertTrustProviderDTOToModel(ctx, *trustProvider, r.client.Tenant, r.client.StackDomain)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *trustProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state trustProviderResourceModel
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

	state = convertTrustProviderDTOToModel(ctx, trustProvider, r.client.Tenant, r.client.StackDomain)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *trustProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get current state
	var state trustProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract external ID from state
	externalID := state.ID.ValueString()

	// Retrieve values from plan
	var plan trustProviderResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var trust aembit.TrustProviderDTO = convertTrustProviderModelToDTO(ctx, plan, &externalID)

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
	state = convertTrustProviderDTOToModel(ctx, *trustProvider, r.client.Tenant, r.client.StackDomain)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *trustProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state trustProviderResourceModel
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
func (r *trustProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import externalId and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Model to DTO conversion methods.
func convertTrustProviderModelToDTO(ctx context.Context, model trustProviderResourceModel, externalID *string) aembit.TrustProviderDTO {
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
		trust.EntityDTO.ExternalID = *externalID
	}

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
		convertKubernetesModelToDTO(model, &trust)
	}
	if model.TerraformWorkspace != nil {
		convertTerraformModelToDTO(model, &trust)
	}

	return trust
}

func appendMatchRuleIfExists(matchRules []aembit.TrustProviderMatchRuleDTO, value basetypes.StringValue, attrName string) []aembit.TrustProviderMatchRuleDTO {
	if len(value.ValueString()) > 0 {
		return append(matchRules, aembit.TrustProviderMatchRuleDTO{
			Attribute: attrName, Value: value.ValueString(),
		})
	}
	return matchRules
}

func appendMatchRulesIfExists(matchRules []aembit.TrustProviderMatchRuleDTO, values []basetypes.StringValue, attrName string) []aembit.TrustProviderMatchRuleDTO {
	if len(values) > 0 {
		for _, value := range values {
			matchRules = appendMatchRuleIfExists(matchRules, value, attrName)
		}
	}
	return matchRules
}

func convertAzureMetadataModelToDTO(model trustProviderResourceModel, dto *aembit.TrustProviderDTO) {
	dto.Provider = "AzureMetadataService"

	dto.MatchRules = make([]aembit.TrustProviderMatchRuleDTO, 0)
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.AzureMetadata.Sku, "AzureSku")
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.AzureMetadata.VMID, "AzureVmId")
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.AzureMetadata.SubscriptionID, "AzureSubscriptionId")
}

func convertAwsRoleModelToDTO(model trustProviderResourceModel, dto *aembit.TrustProviderDTO) {
	dto.Provider = "AWSRole"

	dto.MatchRules = make([]aembit.TrustProviderMatchRuleDTO, 0)
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.AwsRole.AccountID, "AwsAccountId")
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.AwsRole.AssumedRole, "AwsAssumedRole")
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.AwsRole.RoleARN, "AwsRoleARN")
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.AwsRole.Username, "AwsUsername")
}

func convertAwsMetadataModelToDTO(model trustProviderResourceModel, dto *aembit.TrustProviderDTO) {
	dto.Provider = "AWSMetadataService"
	dto.Certificate = base64.StdEncoding.EncodeToString([]byte(model.AwsMetadata.Certificate.ValueString()))
	dto.PemType = "Certificate"

	dto.MatchRules = make([]aembit.TrustProviderMatchRuleDTO, 0)
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.AwsMetadata.AccountID, "AwsAccountId")
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.AwsMetadata.Architecture, "AwsArchitecture")
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.AwsMetadata.AvailabilityZone, "AwsAvailabilityZone")
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.AwsMetadata.BillingProducts, "AwsBillingProducts")
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.AwsMetadata.ImageID, "AwsImageId")
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.AwsMetadata.InstanceID, "AwsInstanceId")
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.AwsMetadata.InstanceType, "AwsInstanceType")
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.AwsMetadata.KernelID, "AwsKernelId")
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.AwsMetadata.MarketplaceProductCodes, "AwsMarketplaceProductCodes")
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.AwsMetadata.PendingTime, "AwsPendingTime")
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.AwsMetadata.PrivateIP, "AwsPrivateIp")
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.AwsMetadata.RamdiskID, "AwsRamdiskId")
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.AwsMetadata.Region, "AwsRegion")
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.AwsMetadata.Version, "AwsVersion")
}

func convertGcpIdentityModelToDTO(model trustProviderResourceModel, dto *aembit.TrustProviderDTO) {
	dto.Provider = "GcpIdentityToken"

	dto.MatchRules = make([]aembit.TrustProviderMatchRuleDTO, 0)
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.GcpIdentity.EMail, "Email")
}

func convertGitHubActionModelToDTO(model trustProviderResourceModel, dto *aembit.TrustProviderDTO) {
	dto.Provider = "GitHubIdentityToken"

	dto.MatchRules = make([]aembit.TrustProviderMatchRuleDTO, 0)
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.GitHubAction.Actor, "GithubActor")
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.GitHubAction.Repository, "GithubRepository")
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.GitHubAction.Workflow, "GithubWorkflow")
}

func convertGitLabJobModelToDTO(model trustProviderResourceModel, dto *aembit.TrustProviderDTO) {
	dto.Provider = "GitLabIdentityToken"

	dto.OidcUrl = model.GitLabJob.OIDCEndpoint.ValueString()
	dto.MatchRules = make([]aembit.TrustProviderMatchRuleDTO, 0)
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.GitLabJob.Subject, "GitLabSubject")
	dto.MatchRules = appendMatchRulesIfExists(dto.MatchRules, model.GitLabJob.Subjects, "GitLabSubject")
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.GitLabJob.ProjectPath, "GitLabProjectPath")
	dto.MatchRules = appendMatchRulesIfExists(dto.MatchRules, model.GitLabJob.ProjectPaths, "GitLabProjectPath")
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.GitLabJob.NamespacePath, "GitLabNamespacePath")
	dto.MatchRules = appendMatchRulesIfExists(dto.MatchRules, model.GitLabJob.NamespacePaths, "GitLabNamespacePath")
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.GitLabJob.RefPath, "GitLabRefPath")
	dto.MatchRules = appendMatchRulesIfExists(dto.MatchRules, model.GitLabJob.RefPaths, "GitLabRefPath")
}

func convertKerberosModelToDTO(model trustProviderResourceModel, dto *aembit.TrustProviderDTO) {
	dto.Provider = "Kerberos"
	dto.AgentControllerIDs = make([]string, len(model.Kerberos.AgentControllerIDs))
	for i, controllerID := range model.Kerberos.AgentControllerIDs {
		dto.AgentControllerIDs[i] = controllerID.ValueString()
	}

	dto.MatchRules = make([]aembit.TrustProviderMatchRuleDTO, 0)
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.Kerberos.Principal, "Principal")
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.Kerberos.Realm, "Realm")
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.Kerberos.SourceIP, "SourceIp")
}

func convertKubernetesModelToDTO(model trustProviderResourceModel, dto *aembit.TrustProviderDTO) {
	dto.Provider = "KubernetesServiceAccount"
	dto.Certificate = base64.StdEncoding.EncodeToString([]byte(model.KubernetesService.PublicKey.ValueString()))
	if len(dto.Certificate) > 0 {
		dto.PemType = "PublicKey"
	}
	dto.OidcUrl = model.KubernetesService.OIDCEndpoint.ValueString()

	dto.MatchRules = make([]aembit.TrustProviderMatchRuleDTO, 0)
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.KubernetesService.Issuer, "KubernetesIss")
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.KubernetesService.Namespace, "KubernetesIoNamespace")
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.KubernetesService.PodName, "KubernetesIoPodName")
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.KubernetesService.ServiceAccountName, "KubernetesIoServiceAccountName")
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.KubernetesService.Subject, "KubernetesSub")
}

func convertTerraformModelToDTO(model trustProviderResourceModel, dto *aembit.TrustProviderDTO) {
	dto.Provider = "TerraformIdentityToken"

	dto.MatchRules = make([]aembit.TrustProviderMatchRuleDTO, 0)
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.TerraformWorkspace.OrganizationID, "TerraformOrganizationId")
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.TerraformWorkspace.ProjectID, "TerraformProjectId")
	dto.MatchRules = appendMatchRuleIfExists(dto.MatchRules, model.TerraformWorkspace.WorkspaceID, "TerraformWorkspaceId")
}

// DTO to Model conversion methods.
func convertTrustProviderDTOToModel(ctx context.Context, dto aembit.TrustProviderDTO, tenant, stackDomain string) trustProviderResourceModel {
	var model trustProviderResourceModel
	model.ID = types.StringValue(dto.EntityDTO.ExternalID)
	model.Name = types.StringValue(dto.EntityDTO.Name)
	model.Description = types.StringValue(dto.EntityDTO.Description)
	model.IsActive = types.BoolValue(dto.EntityDTO.IsActive)
	model.Tags = newTagsModel(ctx, dto.EntityDTO.Tags)

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
		model.KubernetesService = convertKubernetesDTOToModel(dto)
	case "TerraformIdentityToken":
		model.TerraformWorkspace = convertTerraformDTOToModel(dto)
	}

	return model
}

func convertAzureMetadataDTOToModel(dto aembit.TrustProviderDTO) *trustProviderAzureMetadataModel {
	model := &trustProviderAzureMetadataModel{
		Sku:            types.StringNull(),
		VMID:           types.StringNull(),
		SubscriptionID: types.StringNull(),
	}

	for _, rule := range dto.MatchRules {
		switch rule.Attribute {
		case "AzureSku":
			model.Sku = types.StringValue(rule.Value)
		case "AzureVmId":
			model.VMID = types.StringValue(rule.Value)
		case "AzureSubscriptionId":
			model.SubscriptionID = types.StringValue(rule.Value)
		}
	}
	return model
}

func convertAwsRoleDTOToModel(dto aembit.TrustProviderDTO) *trustProviderAwsRoleModel {
	model := &trustProviderAwsRoleModel{
		AccountID:   types.StringNull(),
		AssumedRole: types.StringNull(),
		RoleARN:     types.StringNull(),
		Username:    types.StringNull(),
	}

	for _, rule := range dto.MatchRules {
		switch rule.Attribute {
		case "AwsAccountId":
			model.AccountID = types.StringValue(rule.Value)
		case "AwsAssumedRole":
			model.AssumedRole = types.StringValue(rule.Value)
		case "AwsRoleARN":
			model.RoleARN = types.StringValue(rule.Value)
		case "AwsUsername":
			model.Username = types.StringValue(rule.Value)
		}
	}
	return model
}

func convertAwsMetadataDTOToModel(dto aembit.TrustProviderDTO) *trustProviderAwsMetadataModel {
	decodedCert, _ := base64.StdEncoding.DecodeString(dto.Certificate)

	model := &trustProviderAwsMetadataModel{
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

	for _, rule := range dto.MatchRules {
		switch rule.Attribute {
		case "AwsAccountId":
			model.AccountID = types.StringValue(rule.Value)
		case "AwsArchitecture":
			model.Architecture = types.StringValue(rule.Value)
		case "AwsAvailabilityZone":
			model.AvailabilityZone = types.StringValue(rule.Value)
		case "AwsBillingProducts":
			model.BillingProducts = types.StringValue(rule.Value)
		case "AwsImageId":
			model.ImageID = types.StringValue(rule.Value)
		case "AwsInstanceId":
			model.InstanceID = types.StringValue(rule.Value)
		case "AwsInstanceType":
			model.InstanceType = types.StringValue(rule.Value)
		case "AwsKernelId":
			model.KernelID = types.StringValue(rule.Value)
		case "AwsMarketplaceProductCodes":
			model.MarketplaceProductCodes = types.StringValue(rule.Value)
		case "AwsPendingTime":
			model.PendingTime = types.StringValue(rule.Value)
		case "AwsPrivateIp":
			model.PrivateIP = types.StringValue(rule.Value)
		case "AwsRamdiskId":
			model.RamdiskID = types.StringValue(rule.Value)
		case "AwsRegion":
			model.Region = types.StringValue(rule.Value)
		case "AwsVersion":
			model.Version = types.StringValue(rule.Value)
		}
	}
	return model
}

func convertKerberosDTOToModel(dto aembit.TrustProviderDTO) *trustProviderKerberosModel {
	model := &trustProviderKerberosModel{
		Principal: types.StringNull(),
		Realm:     types.StringNull(),
		SourceIP:  types.StringNull(),
	}
	model.AgentControllerIDs = make([]types.String, len(dto.AgentControllerIDs))
	for i, controllerID := range dto.AgentControllerIDs {
		model.AgentControllerIDs[i] = types.StringValue(controllerID)
	}

	for _, rule := range dto.MatchRules {
		switch rule.Attribute {
		case "Principal":
			model.Principal = types.StringValue(rule.Value)
		case "Realm":
			model.Realm = types.StringValue(rule.Value)
		case "SourceIp":
			model.SourceIP = types.StringValue(rule.Value)
		}
	}
	return model
}

func convertKubernetesDTOToModel(dto aembit.TrustProviderDTO) *trustProviderKubernetesModel {
	decodedKey, _ := base64.StdEncoding.DecodeString(dto.Certificate)

	model := &trustProviderKubernetesModel{
		Issuer:             types.StringNull(),
		Namespace:          types.StringNull(),
		PodName:            types.StringNull(),
		ServiceAccountName: types.StringNull(),
		Subject:            types.StringNull(),
		PublicKey:          types.StringNull(),
		OIDCEndpoint:       types.StringNull(),
	}
	if len(dto.Certificate) > 0 {
		model.PublicKey = types.StringValue(string(decodedKey))
	} else {
		model.OIDCEndpoint = types.StringValue(dto.OidcUrl)
	}

	for _, rule := range dto.MatchRules {
		switch rule.Attribute {
		case "KubernetesIss":
			model.Issuer = types.StringValue(rule.Value)
		case "KubernetesIoNamespace":
			model.Namespace = types.StringValue(rule.Value)
		case "KubernetesIoPodName":
			model.PodName = types.StringValue(rule.Value)
		case "KubernetesIoServiceAccountName":
			model.ServiceAccountName = types.StringValue(rule.Value)
		case "KubernetesSub":
			model.Subject = types.StringValue(rule.Value)
		}
	}
	return model
}

func convertGcpIdentityDTOToModel(dto aembit.TrustProviderDTO) *trustProviderGcpIdentityModel {
	model := &trustProviderGcpIdentityModel{
		EMail: types.StringNull(),
	}

	for _, rule := range dto.MatchRules {
		switch rule.Attribute {
		case "Email":
			model.EMail = types.StringValue(rule.Value)
		}
	}
	return model
}

func convertGitHubActionDTOToModel(dto aembit.TrustProviderDTO) *trustProviderGitHubActionModel {
	model := &trustProviderGitHubActionModel{
		Actor:      types.StringNull(),
		Repository: types.StringNull(),
		Workflow:   types.StringNull(),
	}

	for _, rule := range dto.MatchRules {
		switch rule.Attribute {
		case "GithubActor":
			model.Actor = types.StringValue(rule.Value)
		case "GithubRepository":
			model.Repository = types.StringValue(rule.Value)
		case "GithubWorkflow":
			model.Workflow = types.StringValue(rule.Value)
		}
	}
	return model
}

func convertGitLabJobDTOToModel(dto aembit.TrustProviderDTO, tenant, stackDomain string) *trustProviderGitLabJobModel {
	stackDomain = strings.ToLower(stackDomain) // Force the stack/domain to be lowercase
	stack := strings.Split(stackDomain, ".")[0]
	model := &trustProviderGitLabJobModel{
		OIDCEndpoint:  types.StringValue(dto.OidcUrl),
		Subject:       types.StringNull(),
		ProjectPath:   types.StringNull(),
		NamespacePath: types.StringNull(),
		RefPath:       types.StringNull(),
		OIDCClientID:  types.StringValue(fmt.Sprintf("aembit:%s:%s:identity:gitlab_idtoken:%s", stack, tenant, dto.ExternalID)),
		OIDCAudience:  types.StringValue(fmt.Sprintf("https://%s.id.%s", tenant, stackDomain)),
	}

	for _, rule := range dto.MatchRules {
		switch rule.Attribute {
		case "GitLabSubject":
			if matchRuleOccurrences(dto.MatchRules, rule.Attribute) > 1 {
				model.Subjects = append(model.Subjects, types.StringValue(rule.Value))
			} else {
				model.Subject = types.StringValue(rule.Value)
			}
		case "GitLabProjectPath":
			if matchRuleOccurrences(dto.MatchRules, rule.Attribute) > 1 {
				model.ProjectPaths = append(model.ProjectPaths, types.StringValue(rule.Value))
			} else {
				model.ProjectPath = types.StringValue(rule.Value)
			}
		case "GitLabNamespacePath":
			if matchRuleOccurrences(dto.MatchRules, rule.Attribute) > 1 {
				model.NamespacePaths = append(model.NamespacePaths, types.StringValue(rule.Value))
			} else {
				model.NamespacePath = types.StringValue(rule.Value)
			}
		case "GitLabRefPath":
			if matchRuleOccurrences(dto.MatchRules, rule.Attribute) > 1 {
				model.RefPaths = append(model.RefPaths, types.StringValue(rule.Value))
			} else {
				model.RefPath = types.StringValue(rule.Value)
			}
		}
	}
	return model
}

func convertTerraformDTOToModel(dto aembit.TrustProviderDTO) *trustProviderTerraformModel {
	model := &trustProviderTerraformModel{
		OrganizationID: types.StringNull(),
		ProjectID:      types.StringNull(),
		WorkspaceID:    types.StringNull(),
	}

	for _, rule := range dto.MatchRules {
		switch rule.Attribute {
		case "TerraformOrganizationId":
			model.OrganizationID = types.StringValue(rule.Value)
		case "TerraformProjectId":
			model.ProjectID = types.StringValue(rule.Value)
		case "TerraformWorkspaceId":
			model.WorkspaceID = types.StringValue(rule.Value)
		}
	}
	return model
}

func matchRuleOccurrences(matchRules []aembit.TrustProviderMatchRuleDTO, attribute string) int {
	count := 0
	for _, rule := range matchRules {
		if rule.Attribute == attribute {
			count++
		}
	}
	return count
}

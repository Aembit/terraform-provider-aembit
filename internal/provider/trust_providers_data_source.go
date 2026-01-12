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
	_ datasource.DataSource              = &trustProvidersDataSource{}
	_ datasource.DataSourceWithConfigure = &trustProvidersDataSource{}
)

// NewTrustProvidersDataSource is a helper function to simplify the provider implementation.
func NewTrustProvidersDataSource() datasource.DataSource {
	return &trustProvidersDataSource{}
}

// trustProvidersDataSource is the data source implementation.
type trustProvidersDataSource struct {
	client *aembit.CloudClient
}

// Configure adds the provider configured client to the data source.
func (d *trustProvidersDataSource) Configure(
	_ context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
	d.client = datasourceConfigure(req, resp)
}

// Metadata returns the data source type name.
func (d *trustProvidersDataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_trust_providers"
}

// Schema defines the schema for the resource.
func (d *trustProvidersDataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Manages an trust provider.",
		Attributes: map[string]schema.Attribute{
			"trust_providers": schema.ListNestedAttribute{
				Description: "List of trust providers.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						// ID field is required for Terraform Framework acceptance testing.
						"id": schema.StringAttribute{
							Description: "Unique identifier of the trust provider.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "User-provided name of the trust provider.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "User-provided description of the trust provider.",
							Computed:    true,
						},
						"is_active": schema.BoolAttribute{
							Description: "Active/Inactive status of the trust provider.",
							Computed:    true,
						},
						"tags":     TagsComputedMapAttribute(),
						"tags_all": TagsAllMapAttribute(),
						"azure_metadata": schema.SingleNestedAttribute{
							Description: "Azure Metadata type Trust Provider configuration.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"sku": schema.StringAttribute{Computed: true},
								"skus": schema.SetAttribute{
									Description: "The set of accepted Azure SKUs that are hosting the Client Workloads.",
									ElementType: types.StringType,
									Computed:    true,
								},
								"vm_id": schema.StringAttribute{Computed: true},
								"vm_ids": schema.SetAttribute{
									Description: "The set of accepted Azure VM IDs that are hosting the Client Workloads.",
									ElementType: types.StringType,
									Computed:    true,
								},
								"subscription_id": schema.StringAttribute{Computed: true},
								"subscription_ids": schema.SetAttribute{
									Description: "The set of accepted Azure Subscription IDs that are hosting the Client Workloads.",
									ElementType: types.StringType,
									Computed:    true,
								},
							},
						},
						"aws_role": schema.SingleNestedAttribute{
							Description: "AWS Role type Trust Provider configuration.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"account_id": schema.StringAttribute{
									Description: "The ID of the AWS account that is hosting the Client Workload.",
									Computed:    true,
								},
								"account_ids": schema.SetAttribute{
									Description: "The set of accepted AWS account IDs that are hosting the Client Workloads.",
									ElementType: types.StringType,
									Computed:    true,
								},
								"assumed_role": schema.StringAttribute{
									Description: "The Name of the AWS IAM Role which is running the Client Workload.",
									Computed:    true,
								},
								"assumed_roles": schema.SetAttribute{
									Description: "The set of accepted AWS IAM Roles that are hosting the Client Workloads.",
									ElementType: types.StringType,
									Computed:    true,
								},
								"role_arn": schema.StringAttribute{
									Description: "The ARN of the AWS IAM Role which is running the Client Workload.",
									Computed:    true,
								},
								"role_arns": schema.SetAttribute{
									Description: "The set of accepted AWS IAM Role ARNs that are hosting the Client Workloads.",
									ElementType: types.StringType,
									Computed:    true,
								},
								"username": schema.StringAttribute{
									Description: "The UserID of the AWS IAM Account which is running the Client Workload (not commonly used).",
									Computed:    true,
								},
								"usernames": schema.SetAttribute{
									Description: "The set of accepted AWS IAM Account UserIDs that are hosting the Client Workloads.",
									ElementType: types.StringType,
									Computed:    true,
								},
							},
						},
						"aws_metadata": schema.SingleNestedAttribute{
							Description: "AWS Metadata type Trust Provider configuration.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"certificate": schema.StringAttribute{
									Computed:    true,
									Description: "PEM Certificate to be used for Signature verification",
								},
								"account_id": schema.StringAttribute{Computed: true},
								"account_ids": schema.SetAttribute{
									Description: "The set of accepted AWS Account IDs which are hosting the Client Workloads. Used only for cases where multiple AWS Account IDs can be matched.",
									ElementType: types.StringType,
									Computed:    true,
								},
								"architecture":      schema.StringAttribute{Computed: true},
								"availability_zone": schema.StringAttribute{Computed: true},
								"availability_zones": schema.SetAttribute{
									Description: "The set of accepted AWS Availability Zones which are hosting the Client Workloads. Used only for cases where multiple AWS Availability Zones can be matched.",
									ElementType: types.StringType,
									Computed:    true,
								},
								"billing_products": schema.StringAttribute{Computed: true},
								"image_id":         schema.StringAttribute{Computed: true},
								"instance_id":      schema.StringAttribute{Computed: true},
								"instance_ids": schema.SetAttribute{
									Description: "The set of accepted AWS Instance IDs which are hosting the Client Workloads. Used only for cases where multiple AWS Instance IDs can be matched.",
									ElementType: types.StringType,
									Computed:    true,
								},
								"instance_type": schema.StringAttribute{Computed: true},
								"instance_types": schema.SetAttribute{
									Description: "The set of accepted AWS Instance Types which are hosting the Client Workloads. Used only for cases where multiple AWS Instance Types can be matched.",
									ElementType: types.StringType,
									Computed:    true,
								},
								"kernel_id":                 schema.StringAttribute{Computed: true},
								"marketplace_product_codes": schema.StringAttribute{Computed: true},
								"pending_time":              schema.StringAttribute{Computed: true},
								"private_ip":                schema.StringAttribute{Computed: true},
								"ramdisk_id":                schema.StringAttribute{Computed: true},
								"region":                    schema.StringAttribute{Computed: true},
								"regions": schema.SetAttribute{
									Description: "The set of accepted AWS Regions which are hosting the Client Workloads. Used only for cases where multiple AWS Regions can be matched.",
									ElementType: types.StringType,
									Computed:    true,
								},
								"version": schema.StringAttribute{Computed: true},
							},
						},
						"gcp_identity": schema.SingleNestedAttribute{
							Description: "GCP Identity type Trust Provider configuration.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"email": schema.StringAttribute{
									Description: "The GCP Service Account email address associated with the resource.",
									Computed:    true,
								},
								"emails": schema.SetAttribute{
									Description: "A set of GCP Service Account email addresses that are associated with the resource(s).",
									ElementType: types.StringType,
									Computed:    true,
								},
							},
						},
						"github_action": schema.SingleNestedAttribute{
							Description: "GitHub Action type Trust Provider configuration.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"actor": schema.StringAttribute{
									Description: "The GitHub Actor which initiated the GitHub Action.",
									Computed:    true,
								},
								"actors": schema.SetAttribute{
									Description: "The set of accepted GitHub ID Token Actors which initiated the GitHub Action.",
									ElementType: types.StringType,
									Computed:    true,
								},
								"repository": schema.StringAttribute{
									Description: "The GitHub Repository associated with the GitHub Action ID Token.",
									Computed:    true,
								},
								"repositories": schema.SetAttribute{
									Description: "The set of accepted GitHub ID Token Repositories which initiated the GitHub Action.",
									ElementType: types.StringType,
									Computed:    true,
								},
								"workflow": schema.StringAttribute{
									Description: "The GitHub Workflow execution associated with the GitHub Action ID Token.",
									Computed:    true,
								},
								"workflows": schema.SetAttribute{
									Description: "The set of accepted GitHub ID Token Workflows which initiated the GitHub Action.",
									ElementType: types.StringType,
									Computed:    true,
								},
								"oidc_endpoint": schema.StringAttribute{
									Description: "The Github OIDC Endpoint used for validating Github Action generated ID Tokens.",
									Computed:    true,
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
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"oidc_endpoint": schema.StringAttribute{
									Description: "The GitLab OIDC Endpoint used for validating GitLab Job generated ID Tokens.",
									Computed:    true,
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
									Computed:    true,
								},
								"namespace_paths": schema.SetAttribute{
									Description: "The set of accepted GitLab ID Token Namespace Paths which initiated the GitLab Job.",
									ElementType: types.StringType,
									Computed:    true,
								},
								"project_path": schema.StringAttribute{
									Description: "The GitLab ID Token Project Path which initiated the GitLab Job.",
									Computed:    true,
								},
								"project_paths": schema.SetAttribute{
									Description: "The set of accepted GitLab ID Token Project Paths which initiated the GitLab Job.",
									ElementType: types.StringType,
									Computed:    true,
								},
								"ref_path": schema.StringAttribute{
									Description: "The GitLab ID Token Ref Path which initiated the GitLab Job.",
									Computed:    true,
								},
								"ref_paths": schema.SetAttribute{
									Description: "The set of accepted GitLab ID Token Ref Paths which initiated the GitLab Job.",
									ElementType: types.StringType,
									Computed:    true,
								},
								"subject": schema.StringAttribute{
									Description: "The GitLab ID Token Subject which initiated the GitLab Job.",
									Computed:    true,
								},
								"subjects": schema.SetAttribute{
									Description: "The set of accepted GitLab ID Token Subjects which initiated the GitLab Job.",
									ElementType: types.StringType,
									Computed:    true,
								},
							},
						},
						"kerberos": schema.SingleNestedAttribute{
							Description: "Kerberos type Trust Provider configuration.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"agent_controller_ids": schema.SetAttribute{
									Computed:    true,
									ElementType: types.StringType,
								},
								"principal": schema.StringAttribute{Computed: true},
								"principals": schema.SetAttribute{
									Description: "The set of accepted Kerberos Principals of the authenticated Agent Proxy.",
									ElementType: types.StringType,
									Computed:    true,
								},
								"realm_domain": schema.StringAttribute{
									Description: "The Kerberos Realm or ActiveDirectory Domain of the authenticated Agent Proxy.",
									Computed:    true,
								},
								"realms_domains": schema.SetAttribute{
									Description: "The set of accepted Kerberos Realms or ActiveDirectory Domains which initiated the authenticated Agent Proxy.",
									ElementType: types.StringType,
									Computed:    true,
								},
								"source_ip": schema.StringAttribute{Computed: true},
								"source_ips": schema.SetAttribute{
									Description: "The set of accepted Source IPs of the authenticated Agent Proxy.",
									ElementType: types.StringType,
									Computed:    true,
								},
							},
						},
						"kubernetes_service_account": schema.SingleNestedAttribute{
							Description: "Kubernetes Service Account type Trust Provider configuration.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"issuer": schema.StringAttribute{
									Description: "The Issuer (`iss` claim) of the Kubernetes Service Account Token.",
									Computed:    true,
								},
								"issuers": schema.SetAttribute{
									Description: "The set of accepted Issuer values of the Kubernetes Service Account Token.",
									ElementType: types.StringType,
									Computed:    true,
								},
								"namespace": schema.StringAttribute{
									Description: "The Namespace of the Kubernetes Service Account Token.",
									Computed:    true,
								},
								"namespaces": schema.SetAttribute{
									Description: "The set of accepted Namespace values of the Kubernetes Service Account Token.",
									ElementType: types.StringType,
									Computed:    true,
								},
								"pod_name": schema.StringAttribute{
									Description: "The Pod Name of the Kubernetes Service Account Token.",
									Computed:    true,
								},
								"pod_names": schema.SetAttribute{
									Description: "The set of accepted Pod Name values of the Kubernetes Service Account Token.",
									ElementType: types.StringType,
									Computed:    true,
								},
								"service_account_name": schema.StringAttribute{
									Description: "The Service Account Name of the Kubernetes Service Account Token.",
									Computed:    true,
								},
								"service_account_names": schema.SetAttribute{
									Description: "The set of accepted Service Account Name values of the Kubernetes Service Account Token.",
									ElementType: types.StringType,
									Computed:    true,
								},
								"subject": schema.StringAttribute{
									Description: "The Subject (`sub` claim) of the Kubernetes Service Account Token.",
									Computed:    true,
								},
								"subjects": schema.SetAttribute{
									Description: "The set of accepted Subject values of the Kubernetes Service Account Token.",
									ElementType: types.StringType,
									Computed:    true,
								},
								"oidc_endpoint": schema.StringAttribute{
									Description: "The OIDC Endpoint from which Public Keys can be retrieved for verifying the signature of the Kubernetes Service Account Token.",
									Computed:    true,
								},
								"public_key": schema.StringAttribute{
									Description: "The Public Key that can be used to verify the signature of the Kubernetes Service Account Token.",
									Computed:    true,
								},
								"jwks": schema.StringAttribute{
									Description: "The JSON Web Key Set (JWKS) containing public keys used for signature verification.",
									Computed:    true,
								},
							},
						},
						"terraform_workspace": schema.SingleNestedAttribute{
							Description: "Terraform Workspace type Trust Provider configuration.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"organization_id": schema.StringAttribute{
									Description: "The Organization ID of the calling Terraform Workspace.",
									Computed:    true,
								},
								"organization_ids": schema.SetAttribute{
									Description: "The set of accepted Organization ID values of the calling Terraform Workspace.",
									ElementType: types.StringType,
									Computed:    true,
								},
								"project_id": schema.StringAttribute{
									Description: "The Project ID of the calling Terraform Workspace.",
									Computed:    true,
								},
								"project_ids": schema.SetAttribute{
									Description: "The set of accepted Project ID values of the calling Terraform Workspace.",
									ElementType: types.StringType,
									Computed:    true,
								},
								"workspace_id": schema.StringAttribute{
									Description: "The Workspace ID of the calling Terraform Workspace.",
									Computed:    true,
								},
								"workspace_ids": schema.SetAttribute{
									Description: "The set of accepted Workspace ID values of the calling Terraform Workspace.",
									ElementType: types.StringType,
									Computed:    true,
								},
							},
						},
						"oidc_id_token": schema.SingleNestedAttribute{
							Description: "OIDC ID Token type Trust Provider configuration.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"issuer": schema.StringAttribute{
									Description: "The Issuer (`iss` claim) of the OIDC ID Token Token.",
									Computed:    true,
								},
								"issuers": schema.SetAttribute{
									Description: "The set of accepted Issuer values of the OIDC ID Token Token.",
									ElementType: types.StringType,
									Computed:    true,
								},
								"subject": schema.StringAttribute{
									Description: "The Subject (`sub` claim) of the OIDC ID Token Token.",
									Computed:    true,
								},
								"subjects": schema.SetAttribute{
									Description: "The set of accepted Subject values of the OIDC ID Token Token.",
									ElementType: types.StringType,
									Computed:    true,
								},
								"audience": schema.StringAttribute{
									Description: "The Audience (`aud` claim) of the OIDC ID Token Token.",
									Computed:    true,
								},
								"audiences": schema.SetAttribute{
									Description: "The set of accepted Audience values of the OIDC ID Token Token.",
									ElementType: types.StringType,
									Computed:    true,
								},
								"oidc_endpoint": schema.StringAttribute{
									Description: "The OIDC Endpoint from which Public Keys can be retrieved for verifying the signature of the OIDC ID Token Token.",
									Computed:    true,
								},
								"public_key": schema.StringAttribute{
									Description: "The Public Key that can be used to verify the signature of the OIDC ID Token Token.",
									Computed:    true,
								},
								"jwks": schema.StringAttribute{
									Description: "The JSON Web Key Set (JWKS) containing public keys used for signature verification.",
									Computed:    true,
								},
							},
						},
						"certificate_signed_attestation": schema.SingleNestedAttribute{
							Description: "Certificate Signed Attestation type Trust Provider configuration.",
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *trustProvidersDataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var state models.TrustProvidersDataSourceModel

	trustProviders, err := d.client.GetTrustProviders(nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Aembit Trust Providers",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, trustProvider := range trustProviders {
		trustProviderState := convertTrustProviderDTOToModel(
			ctx,
			trustProvider,
			&models.TrustProviderResourceModel{},
			d.client.Tenant,
			d.client.StackDomain,
		)
		trustProviderState.Tags = newTagsModel(ctx, trustProvider.Tags)
		state.TrustProviders = append(state.TrustProviders, trustProviderState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

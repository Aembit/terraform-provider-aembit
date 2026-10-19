package provider

import (
	"context"
	"strings"

	"terraform-provider-aembit/internal/provider/models"
	"terraform-provider-aembit/internal/provider/validators"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                   = &clientWorkloadResource{}
	_ resource.ResourceWithConfigure      = &clientWorkloadResource{}
	_ resource.ResourceWithImportState    = &clientWorkloadResource{}
	_ resource.ResourceWithModifyPlan     = &clientWorkloadResource{}
	_ resource.ResourceWithValidateConfig = &clientWorkloadResource{}
)

// NewClientWorkloadResource is a helper function to simplify the provider implementation.
func NewClientWorkloadResource() resource.Resource {
	return &clientWorkloadResource{}
}

// clientWorkloadResource is the resource implementation.
type clientWorkloadResource struct {
	client *aembit.CloudClient
}

// Metadata returns the resource type name.
func (r *clientWorkloadResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_client_workload"
}

// Configure adds the provider configured client to the resource.
func (r *clientWorkloadResource) Configure(
	_ context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	r.client = resourceConfigure(req, resp)
}

// Schema defines the schema for the resource.
func (r *clientWorkloadResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			// ID field is required for Terraform Framework acceptance testing.
			"id": schema.StringAttribute{
				Description: "Unique identifier of the Client Workload.",
				Computed:    true,
				Validators: []validator.String{
					validators.UUIDRegexValidation(),
				},
				// Prevent ID from becoming "known after apply" on updates
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"resource_set_id": schema.StringAttribute{
				Description: "ResourceSet unique identifier of the Client Workload.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					validators.UUIDRegexValidation(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name for the Client Workload.",
				Required:    true,
				Validators: []validator.String{
					validators.NameLengthValidation(),
				},
			},
			"description": schema.StringAttribute{
				Description: "Description for the Client Workload.",
				Optional:    true,
				Computed:    true,
			},
			"is_active": schema.BoolAttribute{
				Description: "Active status of the Client Workload.",
				Optional:    true,
				Computed:    true,
			},
			"identities": schema.SetNestedAttribute{
				Description: "Set of Client Workload identities.",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Description: "Client identity type. Possible values are: \n" +
								"\t* `aembitClientId`\n" +
								"\t* `cimdClientId`\n\t\t" +
								"When configured, it cannot be mixed with non-MCP client workload identity types.\n" +
								"\t* `awsAccountId`\n" +
								"\t* `awsEc2InstanceId`\n" +
								//"\t* `awsEcsServiceName`\n" +	// Hiding for now
								"\t* `awsEcsTaskFamily`\n" +
								"\t* `awsLambdaArn`\n" +
								"\t* `awsRegion`\n" +
								"\t* `azureSubscriptionId`\n" +
								"\t* `azureVmId`\n" +
								"\t* `gcpIdentityToken`\n" +
								"\t* `githubIdTokenRepository`\n" +
								"\t* `githubIdTokenSubject`\n" +
								"\t* `gitlabIdTokenSubject`\n" +
								"\t* `gitlabIdTokenProjectPath`\n" +
								"\t* `gitlabIdTokenNamespacePath`\n" +
								"\t* `gitlabIdTokenRefPath`\n" +
								"\t* `hostname`\n" +
								"\t* `k8sNamespace`\n" +
								"\t* `k8sPodName`\n" +
								"\t* `k8sPodNamePrefix`\n" +
								"\t* `k8sServiceAccountName`\n" +
								"\t* `k8sServiceAccountUID`\n" +
								"\t* `oauthRedirectUri`\n\t\t" +
								"When configured, it cannot be mixed with non-MCP client workload identity types.\n" +
								"\t* `oauthScope`\n" +
								"\t* `oidcIdToken`\n" +
								"\t* `oidcIdTokenAudience`\n" +
								"\t* `oidcIdTokenIssuer`\n" +
								"\t* `oidcIdTokenSubject`\n" +
								"\t* `processName`\n" +
								"\t* `processUserName`\n" +
								"\t* `processPath`\n" +
								"\t* `processCommandLine`\n" +
								"\t* `sourceIPAddress`\n" +
								"\t* `terraformIdTokenOrganizationId`\n" +
								"\t* `terraformIdTokenProjectId`\n" +
								"\t* `terraformIdTokenWorkspaceId`\n" +
								"\t* `userGroup`\n",
							Required: true,
							Validators: []validator.String{
								stringvalidator.OneOf([]string{
									"aembitClientId",
									"cimdClientId",
									"awsAccountId",
									"awsEc2InstanceId",
									//"awsEcsServiceName",	// Hiding for now
									"awsEcsTaskFamily",
									"awsLambdaArn",
									"awsRegion",
									"azureSubscriptionId",
									"azureVmId",
									"gcpIdentityToken",
									"githubIdTokenRepository",
									"githubIdTokenSubject",
									"gitlabIdTokenSubject",
									"gitlabIdTokenProjectPath",
									"gitlabIdTokenNamespacePath",
									"gitlabIdTokenRefPath",
									"hostname",
									"k8sNamespace",
									"k8sPodName",
									"k8sPodNamePrefix",
									"k8sServiceAccountName",
									"k8sServiceAccountUID",
									"oauthRedirectUri",
									"oauthScope",
									"oidcIdToken",
									"oidcIdTokenAudience",
									"oidcIdTokenIssuer",
									"oidcIdTokenSubject",
									"processName",
									"processUserName",
									"processPath",
									"processCommandLine",
									"sourceIPAddress",
									"terraformIdTokenOrganizationId",
									"terraformIdTokenProjectId",
									"terraformIdTokenWorkspaceId",
									"userGroup",
								}...),
							},
						},
						"value": schema.StringAttribute{
							Description: "Client identity value.",
							Required:    true,
						},
						"claim_name": schema.StringAttribute{
							Description: "Client identity claim name. Applicable for oidcIdToken Client identity type.",
							Optional:    true,
						},
						"identity_provider_id": schema.StringAttribute{
							Description: "Identity Provider unique identifier. Applicable and required for userGroup Client identity type.",
							Optional:    true,
							Validators: []validator.String{
								validators.UUIDRegexValidation(),
							},
						},
					},
				},
			},
			"enforce_sso": schema.BoolAttribute{
				Description: "Whether SSO authentication is enforced for MCP authorization. This is only applicable when the client workload identities use an MCP identity (`oauthRedirectUri` or `cimdClientId`), which cannot be mixed with other identity types.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"sso_identity_providers": schema.SetAttribute{
				Description: "Set of SSO Identity Provider IDs used for MCP authorization. This is only applicable to MCP client workloads (`oauthRedirectUri` or `cimdClientId`) when `enforce_sso` is true.",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(validators.UUIDRegexValidation()),
				},
			},
			"tags":     TagsMapAttribute(),
			"tags_all": TagsAllMapAttribute(),
			"standalone_certificate_authority": schema.StringAttribute{
				Description: "Standalone Certificate Authority ID configured for this Client Workload.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validators.UUIDRegexValidation(),
				},
			},
		},
	}
}

func (r *clientWorkloadResource) ModifyPlan(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
) {
	modifyPlanForTagsAll(ctx, req, resp, r.client.DefaultTags)
	modifyPlanForResourceSetId(ctx, req, resp, r.client)
}

func (r *clientWorkloadResource) ValidateConfig(
	ctx context.Context,
	req resource.ValidateConfigRequest,
	resp *resource.ValidateConfigResponse,
) {
	var config models.ClientWorkloadResourceModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var identities []models.IdentitiesModel
	if config.Identities.IsUnknown() {
		return
	}
	if len(config.Identities.Elements()) > 0 {
		diags = config.Identities.ElementsAs(ctx, &identities, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	if len(identities) == 0 {
		return
	}

	// MCP client identities cannot be combined with other identity types.
	identityDiags, hasMcpIdentity := validateMcpIdentityTypeForConfig(identities)
	resp.Diagnostics.Append(identityDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// UserGroup client identities cross-field validation.
	resp.Diagnostics.Append(validateUserGroupIdentityTypeForConfig(identities)...)
	if resp.Diagnostics.HasError() {
		return
	}

	effectiveEnforceSso := true
	if config.EnforceSso.IsUnknown() {
		return
	}
	if !config.EnforceSso.IsNull() {
		effectiveEnforceSso = config.EnforceSso.ValueBool()
	}

	if config.SsoIdentityProviders.IsUnknown() {
		return
	}

	resp.Diagnostics.Append(validateMcpConfigurationForConfig(
		hasMcpIdentity,
		effectiveEnforceSso,
		config.EnforceSso,
		config.SsoIdentityProviders,
	)...)

	// CIMD Client ID values must be valid URLs.
	resp.Diagnostics.Append(validateCimdClientIdIdentityTypeForConfig(identities)...)
}

// Create creates the resource and sets the initial Terraform state.
func (r *clientWorkloadResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	// Retrieve values from plan
	var plan models.ClientWorkloadResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	workload, diags := convertClientWorkloadModelToDTO(ctx, plan, nil, r.client.DefaultTags)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceSetId := getResourceSetId(plan.ResourceSetId, r.client)

	// Create new Client Workload
	clientWorkload, err := r.client.CreateClientWorkload(workload, nil, &resourceSetId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Client Workload",
			"Could not create Client Workload, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan = convertClientWorkloadDTOToModel(ctx, *clientWorkload, &plan)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *clientWorkloadResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	// Get current state
	var state models.ClientWorkloadResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceSetId := getResourceSetId(state.ResourceSetId, r.client)

	// Get refreshed trust value from Aembit
	clientWorkload, err, notFound := r.client.GetClientWorkload(ctx, state.ID.ValueString(), nil, &resourceSetId)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Error reading Aembit Client Workload",
			"Could not read Aembit External ID from Terraform state "+state.ID.ValueString()+": "+err.Error(),
		)

		// If the resource is not found on Aembit Cloud, delete it locally
		if notFound {
			resp.State.RemoveResource(ctx)
		}
		return
	}

	// Overwrite items with refreshed state
	state = convertClientWorkloadDTOToModel(ctx, clientWorkload, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *clientWorkloadResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	// Get current state
	var state models.ClientWorkloadResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract external ID from state
	externalID := state.ID.ValueString()

	// Retrieve values from plan
	var plan models.ClientWorkloadResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !state.ResourceSetId.Equal(plan.ResourceSetId) {
		resp.Diagnostics.AddError(
			"Error updating Client Workload",
			"Changing the ResourceSet of the resource is not supported.",
		)
		return
	}

	// Generate API request body from plan
	workload, diags := convertClientWorkloadModelToDTO(ctx, plan, &externalID, r.client.DefaultTags)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update Client Workload
	clientWorkload, err := r.client.UpdateClientWorkload(workload, nil, &workload.ResourceSet)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Client Workload",
			"Could not update Client Workload, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	state = convertClientWorkloadDTOToModel(ctx, *clientWorkload, &plan)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *clientWorkloadResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	// Retrieve values from state
	var state models.ClientWorkloadResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceSetId := state.ResourceSetId.ValueString()

	// Check if Client Workload is Active - if it is, disable it first
	if state.IsActive == types.BoolValue(true) {
		_, err := r.client.DisableClientWorkload(state.ID.ValueString(), nil, &resourceSetId)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error disabling Client Workload",
				"Could not disable Client Workload, unexpected error: "+err.Error(),
			)
			return
		}
	}

	// Delete existing Client Workload
	_, err := r.client.DeleteClientWorkload(ctx, state.ID.ValueString(), nil, &resourceSetId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting Client Workload",
			"Could not delete client workload, unexpected error: "+err.Error(),
		)
		return
	}
}

// Imports an existing resource by passing externalId.
func (r *clientWorkloadResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	idParts := strings.Split(req.ID, ",")

	if len(idParts) == 2 {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("resource_set_id"), idParts[0])...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[1])...)
		return
	}

	// Retrieve import externalId and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func convertClientWorkloadModelToDTO(
	ctx context.Context,
	model models.ClientWorkloadResourceModel,
	externalID *string,
	defaultTags map[string]string,
) (aembit.ClientWorkloadExternalDTO, diag.Diagnostics) {
	var workload aembit.ClientWorkloadExternalDTO
	var diags diag.Diagnostics

	var identities []models.IdentitiesModel
	if len(model.Identities.Elements()) > 0 {
		diags.Append(model.Identities.ElementsAs(ctx, &identities, false)...)
	}
	if diags.HasError() {
		diags.AddAttributeError(
			path.Root("identities"),
			"Invalid identities",
			"Client Workload must contain valid client type identities.",
		)
		return workload, diags
	}
	if len(identities) == 0 {
		diags.AddAttributeError(
			path.Root("identities"),
			"Missing identities",
			"Client Workload must contain  at least one client workload identity.",
		)
		return workload, diags
	}

	diags.Append(validateMcpIdentityType(identities, model.EnforceSso, model.SsoIdentityProviders)...)
	if diags.HasError() {
		return workload, diags
	}

	workload.EntityDTO = aembit.EntityDTO{
		Name:        model.Name.ValueString(),
		Description: model.Description.ValueString(),
		IsActive:    model.IsActive.ValueBool(),
	}

	for _, identity := range identities {
		key := identity.ClaimName.ValueString()
		if identity.Type.ValueString() == "userGroup" {
			key = identity.IdentityProviderId.ValueString()
		}
		workload.Identities = append(workload.Identities, aembit.ClientWorkloadIdentityDTO{
			Type:  identity.Type.ValueString(),
			Value: identity.Value.ValueString(),
			Key:   key,
		})
	}

	if externalID != nil {
		workload.ExternalID = *externalID
	}

	hasMcp := hasMcpIdentity(identities)
	workload.EnforceSso = true
	if hasMcp {
		if !model.EnforceSso.IsNull() {
			workload.EnforceSso = model.EnforceSso.ValueBool()
		}
	}

	if hasMcp {
		if !model.SsoIdentityProviders.IsNull() {
			_ = model.SsoIdentityProviders.ElementsAs(ctx, &workload.SsoIdentityProviders, false)
		}
	}

	workload.StandaloneCertificateAuthority = model.StandaloneCertificateAuthority.ValueString()

	workload.Tags = collectAllTagsDto(ctx, defaultTags, model.Tags)

	workload.ResourceSet = model.ResourceSetId.ValueString()
	return workload, diags
}

func convertClientWorkloadDTOToModel(
	ctx context.Context,
	dto aembit.ClientWorkloadExternalDTO,
	planModel *models.ClientWorkloadResourceModel,
) models.ClientWorkloadResourceModel {
	var model models.ClientWorkloadResourceModel
	hasMcp := hasMcpIdentityDTO(dto.Identities)
	model.ID = types.StringValue(dto.ExternalID)
	model.Name = types.StringValue(dto.Name)
	model.Description = types.StringValue(dto.Description)
	model.IsActive = types.BoolValue(dto.IsActive)
	model.Identities = newClientWorkloadIdentityModel(ctx, dto.Identities)
	model.EnforceSso = types.BoolValue(true)
	if hasMcp {
		model.EnforceSso = types.BoolValue(dto.EnforceSso)
	}
	model.SsoIdentityProviders = newStringSetModel(ctx, dto.SsoIdentityProviders)
	// handle tags
	model.Tags = newTagsModelFromPlan(ctx, planModel.Tags)
	model.TagsAll = newTagsModel(ctx, dto.Tags)

	if dto.StandaloneCertificateAuthority == "" {
		model.StandaloneCertificateAuthority = types.StringNull()
	} else {
		model.StandaloneCertificateAuthority = types.StringValue(dto.StandaloneCertificateAuthority)
	}

	model.ResourceSetId = types.StringValue(dto.ResourceSet)
	return model
}

func newClientWorkloadIdentityModel(
	ctx context.Context,
	clientWorkloadIdentities []aembit.ClientWorkloadIdentityDTO,
) types.Set {
	identities := make([]models.IdentitiesModel, len(clientWorkloadIdentities))

	for i, identity := range clientWorkloadIdentities {
		var claimName types.String
		var identityProviderId types.String

		if identity.Type == "userGroup" {
			identityProviderId = types.StringValue(identity.Key)
			claimName = types.StringNull()
		} else {
			identityProviderId = types.StringNull()
			if identity.Key == "" {
				claimName = types.StringNull()
			} else {
				claimName = types.StringValue(identity.Key)
			}
		}

		identities[i] = models.IdentitiesModel{
			Type:               types.StringValue(identity.Type),
			Value:              types.StringValue(identity.Value),
			ClaimName:          claimName,
			IdentityProviderId: identityProviderId,
		}
	}

	s, _ := types.SetValueFrom(ctx, models.TfIdentityObjectType, identities)
	return s
}

// hasMcpIdentity returns true if any configured identity is an MCP client workload type (oauthRedirectUri or cimdClientId).
func hasMcpIdentity(identities []models.IdentitiesModel) bool {
	for _, identity := range identities {
		if identity.Type.IsNull() {
			continue
		}

		val := identity.Type.ValueString()
		if val == "oauthRedirectUri" || val == "cimdClientId" {
			return true
		}
	}

	return false
}

// validateMcpIdentityType checks identity types and MCP configuration for client workloads.
func validateMcpIdentityType(
	identities []models.IdentitiesModel,
	enforceSso types.Bool,
	ssoIdentityProviders types.Set,
) diag.Diagnostics {
	var diags diag.Diagnostics
	identityDiags, hasMcp := validateMcpIdentityTypeForConfig(identities)
	diags.Append(identityDiags...)
	if diags.HasError() {
		return diags
	}

	effectiveEnforceSso := true
	if !enforceSso.IsNull() {
		effectiveEnforceSso = enforceSso.ValueBool()
	}

	diags.Append(validateMcpConfigurationForConfig(
		hasMcp,
		effectiveEnforceSso,
		enforceSso,
		ssoIdentityProviders,
	)...)
	if diags.HasError() {
		return diags
	}

	diags.Append(validateCimdClientIdIdentityTypeForConfig(identities)...)
	if diags.HasError() {
		return diags
	}

	diags.Append(validateUserGroupIdentityTypeForConfig(identities)...)

	return diags
}

func validateUserGroupIdentityTypeForConfig(
	identities []models.IdentitiesModel,
) diag.Diagnostics {
	var diags diag.Diagnostics

	for _, identity := range identities {
		if identity.Type.IsNull() || identity.Type.IsUnknown() {
			continue
		}

		if identity.Type.ValueString() == "userGroup" {
			if !identity.IdentityProviderId.IsUnknown() && (identity.IdentityProviderId.IsNull() || identity.IdentityProviderId.ValueString() == "") {
				diags.AddAttributeError(
					path.Root("identities"),
					"Invalid client workload identities configuration",
					"Identity Provider ID is required when client identity type is userGroup.",
				)
			}
			if !identity.ClaimName.IsUnknown() && !identity.ClaimName.IsNull() && identity.ClaimName.ValueString() != "" {
				diags.AddAttributeError(
					path.Root("identities"),
					"Invalid client workload identities configuration",
					"Claim name is only applicable for oidcIdToken client identity type and cannot be specified for userGroup.",
				)
			}
		} else {
			if !identity.IdentityProviderId.IsNull() && (identity.IdentityProviderId.IsUnknown() || identity.IdentityProviderId.ValueString() != "") {
				diags.AddAttributeError(
					path.Root("identities"),
					"Invalid client workload identities configuration",
					"Identity Provider ID is only applicable for userGroup client identity type.",
				)
			}
		}
	}

	return diags
}

func validateCimdClientIdIdentityTypeForConfig(
	identities []models.IdentitiesModel,
) diag.Diagnostics {
	var diags diag.Diagnostics

	for _, identity := range identities {
		if identity.Type.IsNull() || identity.Value.IsNull() {
			continue
		}

		if identity.Type.ValueString() == "cimdClientId" {
			val := identity.Value.ValueString()
			if !validators.CimdUrlRegex.MatchString(val) {
				diags.AddAttributeError(
					path.Root("identities"),
					"Invalid CIMD Client ID Value",
					"CIMD Client ID must be a valid https URL (starts with https://, contains a path, and has no query string, fragment, or user information).",
				)
			}
		}
	}

	return diags
}

// validateMcpIdentityTypeForConfig ensures MCP client identities (oauthRedirectUri or cimdClientId) cannot be mixed with non-MCP identities.
func validateMcpIdentityTypeForConfig(
	identities []models.IdentitiesModel,
) (diag.Diagnostics, bool) {
	var diags diag.Diagnostics
	hasMcp := false
	hasNonMcp := false

	for _, identity := range identities {
		if identity.Type.IsNull() {
			continue
		}

		val := identity.Type.ValueString()
		if val == "oauthRedirectUri" || val == "cimdClientId" {
			hasMcp = true
			continue
		}

		hasNonMcp = true
	}

	if hasMcp && hasNonMcp {
		diags.AddAttributeError(
			path.Root("identities"),
			"Invalid client workload identities configuration",
			"MCP client workload identities (`oauthRedirectUri` or `cimdClientId`) cannot be combined with other identity types.",
		)
	}

	return diags, hasMcp
}

// validateMcpConfigurationForConfig validates enforce_sso and sso_identity_providers consistency for MCP and non-MCP workloads.
func validateMcpConfigurationForConfig(
	hasMcp bool,
	effectiveEnforceSso bool,
	enforceSso types.Bool,
	ssoIdentityProviders types.Set,
) diag.Diagnostics {
	var diags diag.Diagnostics

	if hasMcp {
		if effectiveEnforceSso && (ssoIdentityProviders.IsNull() || len(ssoIdentityProviders.Elements()) == 0) {
			diags.AddAttributeError(
				path.Root("sso_identity_providers"),
				"Missing SSO Identity Providers configuration",
				"`sso_identity_providers` must contain at least one identity provider when an MCP client identity (`oauthRedirectUri` or `cimdClientId`) is configured and `enforce_sso` is `true`.",
			)
		}

		if !effectiveEnforceSso && !ssoIdentityProviders.IsNull() && len(ssoIdentityProviders.Elements()) > 0 {
			diags.AddAttributeError(
				path.Root("sso_identity_providers"),
				"Invalid SSO Identity Providers configuration",
				"`sso_identity_providers` cannot be configured when `enforce_sso` is `false`.",
			)
		}

		return diags
	}

	if !enforceSso.IsNull() && !enforceSso.ValueBool() {
		diags.AddAttributeError(
			path.Root("enforce_sso"),
			"Invalid enforce_sso configuration",
			"`enforce_sso` can only be set to `false` when one of the client workload identities has an MCP type (`oauthRedirectUri` or `cimdClientId`).",
		)
	}

	if !ssoIdentityProviders.IsNull() && len(ssoIdentityProviders.Elements()) > 0 {
		diags.AddAttributeError(
			path.Root("sso_identity_providers"),
			"Invalid SSO Identity Providers configuration",
			"`sso_identity_providers` can only be configured when one of the client workload identities has an MCP type (`oauthRedirectUri` or `cimdClientId`).",
		)
	}

	return diags
}

// hasMcpIdentityDTO returns true if any identity in the DTO is an MCP client workload type (oauthRedirectUri or cimdClientId).
func hasMcpIdentityDTO(identities []aembit.ClientWorkloadIdentityDTO) bool {
	for _, identity := range identities {
		if identity.Type == "oauthRedirectUri" || identity.Type == "cimdClientId" {
			return true
		}
	}

	return false
}

func newStringSetModel(ctx context.Context, values []string) types.Set {
	if len(values) == 0 {
		return types.SetNull(types.StringType)
	}

	s, _ := types.SetValueFrom(ctx, types.StringType, values)
	return s
}

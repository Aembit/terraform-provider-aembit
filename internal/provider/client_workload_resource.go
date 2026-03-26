package provider

import (
	"context"

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
								"\t* `oauthRedirectUri`\n\t \t" +
								"When configured, it must be the only client workload identity type in the set.\n" +
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
								"\t* `terraformIdTokenWorkspaceId`\n",
							Required: true,
							Validators: []validator.String{
								stringvalidator.OneOf([]string{
									"aembitClientId",
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
					},
				},
			},
			"enforce_sso": schema.BoolAttribute{
				Description: "Whether SSO authentication is enforced for MCP authorization. This is only applicable when the client workload identities use `oauthRedirectUri`, which must be the only identity type in the set.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"sso_identity_providers": schema.SetAttribute{
				Description: "Set of SSO Identity Provider IDs used for MCP authorization. This is only applicable when 'enable_sso' is true.",
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

	// Redirect URI client identities cannot be combined with other identity types.
	identityDiags, hasRedirectURI := validateRedirectURIIdentityTypeForConfig(identities)
	resp.Diagnostics.Append(identityDiags...)
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

	resp.Diagnostics.Append(validateRedirectURIConfigurationForConfig(
		hasRedirectURI,
		effectiveEnforceSso,
		config.EnforceSso,
		config.SsoIdentityProviders,
	)...)
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

	// Create new Client Workload
	clientWorkload, err := r.client.CreateClientWorkload(workload, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating client workload",
			"Could not create client workload, unexpected error: "+err.Error(),
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

	// Get refreshed workload value from Aembit
	clientWorkload, err, notFound := r.client.GetClientWorkload(ctx, state.ID.ValueString(), nil)
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

	// Generate API request body from plan
	workload, diags := convertClientWorkloadModelToDTO(ctx, plan, &externalID, r.client.DefaultTags)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update Client Workload
	clientWorkload, err := r.client.UpdateClientWorkload(workload, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating client workload",
			"Could not update client workload, unexpected error: "+err.Error(),
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

	// Check if Client Workload is Active - if it is, disable it first
	if state.IsActive == types.BoolValue(true) {
		_, err := r.client.DisableClientWorkload(state.ID.ValueString(), nil)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error disabling Client Workload",
				"Could not disable Client Workload, unexpected error: "+err.Error(),
			)
			return
		}
	}

	// Delete existing Client Workload
	_, err := r.client.DeleteClientWorkload(ctx, state.ID.ValueString(), nil)
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

	diags.Append(validateRedirectURIIdentityType(identities, model.EnforceSso, model.SsoIdentityProviders)...)
	if diags.HasError() {
		return workload, diags
	}

	workload.EntityDTO = aembit.EntityDTO{
		Name:        model.Name.ValueString(),
		Description: model.Description.ValueString(),
		IsActive:    model.IsActive.ValueBool(),
	}

	for _, identity := range identities {
		workload.Identities = append(workload.Identities, aembit.ClientWorkloadIdentityDTO{
			Type:  identity.Type.ValueString(),
			Value: identity.Value.ValueString(),
			Key:   identity.ClaimName.ValueString(),
		})
	}

	if externalID != nil {
		workload.ExternalID = *externalID
	}

	hasRedirectURI := hasRedirectURIIdentity(identities)
	workload.EnforceSso = true
	if hasRedirectURI {
		if !model.EnforceSso.IsNull() {
			workload.EnforceSso = model.EnforceSso.ValueBool()
		}
	}

	if hasRedirectURI {
		if !model.SsoIdentityProviders.IsNull() {
			_ = model.SsoIdentityProviders.ElementsAs(ctx, &workload.SsoIdentityProviders, false)
		}
	}

	workload.StandaloneCertificateAuthority = model.StandaloneCertificateAuthority.ValueString()

	workload.Tags = collectAllTagsDto(ctx, defaultTags, model.Tags)

	return workload, diags
}

func convertClientWorkloadDTOToModel(
	ctx context.Context,
	dto aembit.ClientWorkloadExternalDTO,
	planModel *models.ClientWorkloadResourceModel,
) models.ClientWorkloadResourceModel {
	var model models.ClientWorkloadResourceModel
	hasRedirectURI := hasOAuthRedirectURIIdentityDTO(dto.Identities)
	model.ID = types.StringValue(dto.ExternalID)
	model.Name = types.StringValue(dto.Name)
	model.Description = types.StringValue(dto.Description)
	model.IsActive = types.BoolValue(dto.IsActive)
	model.Identities = newClientWorkloadIdentityModel(ctx, dto.Identities)
	model.EnforceSso = types.BoolValue(true)
	if hasRedirectURI {
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

	return model
}

func newClientWorkloadIdentityModel(
	ctx context.Context,
	clientWorkloadIdentities []aembit.ClientWorkloadIdentityDTO,
) types.Set {
	identities := make([]models.IdentitiesModel, len(clientWorkloadIdentities))

	for i, identity := range clientWorkloadIdentities {
		identities[i] = models.IdentitiesModel{
			Type:  types.StringValue(identity.Type),
			Value: types.StringValue(identity.Value),
			ClaimName: func() types.String {
				if identity.Key == "" {
					return types.StringNull()
				}
				return types.StringValue(identity.Key)
			}(),
		}
	}

	s, _ := types.SetValueFrom(ctx, models.TfIdentityObjectType, identities)
	return s
}

func hasRedirectURIIdentity(identities []models.IdentitiesModel) bool {
	for _, identity := range identities {
		if identity.Type.IsNull() {
			continue
		}

		if identity.Type.ValueString() == "oauthRedirectUri" {
			return true
		}
	}

	return false
}

func validateRedirectURIIdentityType(
	identities []models.IdentitiesModel,
	enforceSso types.Bool,
	ssoIdentityProviders types.Set,
) diag.Diagnostics {
	var diags diag.Diagnostics
	identityDiags, hasRedirectURI := validateRedirectURIIdentityTypeForConfig(identities)
	diags.Append(identityDiags...)
	if diags.HasError() {
		return diags
	}

	effectiveEnforceSso := true
	if !enforceSso.IsNull() {
		effectiveEnforceSso = enforceSso.ValueBool()
	}

	diags.Append(validateRedirectURIConfigurationForConfig(
		hasRedirectURI,
		effectiveEnforceSso,
		enforceSso,
		ssoIdentityProviders,
	)...)

	return diags
}

func validateRedirectURIIdentityTypeForConfig(
	identities []models.IdentitiesModel,
) (diag.Diagnostics, bool) {
	var diags diag.Diagnostics
	hasRedirectURI := false
	hasNonRedirectURI := false

	for _, identity := range identities {
		if identity.Type.IsNull() {
			continue
		}

		if identity.Type.ValueString() == "oauthRedirectUri" {
			hasRedirectURI = true
			continue
		}

		hasNonRedirectURI = true
	}

	if hasRedirectURI && hasNonRedirectURI {
		diags.AddAttributeError(
			path.Root("identities"),
			"Invalid client workload identities configuration",
			"`oauthRedirectUri` must be the only client workload identity type when it is configured.",
		)
	}

	return diags, hasRedirectURI
}

func validateRedirectURIConfigurationForConfig(
	hasRedirectURI bool,
	effectiveEnforceSso bool,
	enforceSso types.Bool,
	ssoIdentityProviders types.Set,
) diag.Diagnostics {
	var diags diag.Diagnostics

	if hasRedirectURI {
		if effectiveEnforceSso && (ssoIdentityProviders.IsNull() || len(ssoIdentityProviders.Elements()) == 0) {
			diags.AddAttributeError(
				path.Root("sso_identity_providers"),
				"Missing SSO Identity Providers configuration",
				"`sso_identity_providers` must contain at least one identity provider when `oauthRedirectUri` is configured and `enforce_sso` is `true`.",
			)
		}

		if !effectiveEnforceSso && !ssoIdentityProviders.IsNull() && len(ssoIdentityProviders.Elements()) > 0 {
			diags.AddAttributeError(
				path.Root("sso_identity_providers"),
				"Invalid SSO Identity Providers configuration",
				"`sso_identity_providers` cannot be configured when `oauthRedirectUri` is used and `enforce_sso` is `false`.",
			)
		}

		return diags
	}

	if !enforceSso.IsNull() && !enforceSso.ValueBool() {
		diags.AddAttributeError(
			path.Root("enforce_sso"),
			"Invalid enforce_sso configuration",
			"`enforce_sso` can only be set to `false` when one of the client workload identities has type `oauthRedirectUri`.",
		)
	}

	if !ssoIdentityProviders.IsNull() && len(ssoIdentityProviders.Elements()) > 0 {
		diags.AddAttributeError(
			path.Root("sso_identity_providers"),
			"Invalid SSO Identity Providers configuration",
			"`sso_identity_providers` can only be configured when one of the client workload identities has type `oauthRedirectUri`.",
		)
	}

	return diags
}

func hasOAuthRedirectURIIdentityDTO(identities []aembit.ClientWorkloadIdentityDTO) bool {
	for _, identity := range identities {
		if identity.Type == "oauthRedirectUri" {
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

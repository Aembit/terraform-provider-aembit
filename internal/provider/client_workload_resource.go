package provider

import (
	"context"

	"terraform-provider-aembit/internal/provider/models"
	"terraform-provider-aembit/internal/provider/validators"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &clientWorkloadResource{}
	_ resource.ResourceWithConfigure   = &clientWorkloadResource{}
	_ resource.ResourceWithImportState = &clientWorkloadResource{}
	_ resource.ResourceWithModifyPlan  = &clientWorkloadResource{}
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
								"\t* `oidcIdTokenAudience`\n" +
								"\t* `oidcIdTokenIssuer`\n" +
								"\t* `oidcIdTokenSubject`\n" +
								"\t* `processName`\n" +
								"\t* `processUserName`\n" +
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
									"oidcIdTokenAudience",
									"oidcIdTokenIssuer",
									"oidcIdTokenSubject",
									"processName",
									"processUserName",
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
					},
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
	workload := convertClientWorkloadModelToDTO(ctx, plan, nil, r.client.DefaultTags)

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
	workload := convertClientWorkloadModelToDTO(ctx, plan, &externalID, r.client.DefaultTags)

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
) aembit.ClientWorkloadExternalDTO {
	var workload aembit.ClientWorkloadExternalDTO
	workload.EntityDTO = aembit.EntityDTO{
		Name:        model.Name.ValueString(),
		Description: model.Description.ValueString(),
		IsActive:    model.IsActive.ValueBool(),
	}

	var identities []models.IdentitiesModel
	if len(model.Identities.Elements()) > 0 {
		_ = model.Identities.ElementsAs(ctx, &identities, false)

		for _, identity := range identities {
			workload.Identities = append(workload.Identities, aembit.ClientWorkloadIdentityDTO{
				Type:  identity.Type.ValueString(),
				Value: identity.Value.ValueString(),
			})
		}

	}

	if externalID != nil {
		workload.ExternalID = *externalID
	}

	workload.StandaloneCertificateAuthority = model.StandaloneCertificateAuthority.ValueString()

	workload.Tags = collectAllTagsDto(ctx, defaultTags, model.Tags)

	return workload
}

func convertClientWorkloadDTOToModel(
	ctx context.Context,
	dto aembit.ClientWorkloadExternalDTO,
	planModel *models.ClientWorkloadResourceModel,
) models.ClientWorkloadResourceModel {
	var model models.ClientWorkloadResourceModel
	model.ID = types.StringValue(dto.ExternalID)
	model.Name = types.StringValue(dto.Name)
	model.Description = types.StringValue(dto.Description)
	model.IsActive = types.BoolValue(dto.IsActive)
	model.Identities = newClientWorkloadIdentityModel(ctx, dto.Identities)
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
		}
	}

	s, _ := types.SetValueFrom(ctx, models.TfIdentityObjectType, identities)
	return s
}

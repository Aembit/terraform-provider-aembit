package provider

import (
	"context"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-aembit/internal/provider/models"
	"terraform-provider-aembit/internal/provider/validators"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &accessPolicyResource{}
	_ resource.ResourceWithConfigure   = &accessPolicyResource{}
	_ resource.ResourceWithImportState = &accessPolicyResource{}
)

// NewAccessPolicyResource is a helper function to simplify the provider implementation.
func NewAccessPolicyResource() resource.Resource {
	return &accessPolicyResource{}
}

// accessPolicyResource is the resource implementation.
type accessPolicyResource struct {
	client *aembit.CloudClient
}

// Metadata returns the resource type name.
func (r *accessPolicyResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_access_policy"
}

// Configure adds the provider configured client to the resource.
func (r *accessPolicyResource) Configure(
	_ context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	r.client = resourceConfigure(req, resp)
}

// Schema defines the schema for the resource.
func (r *accessPolicyResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			// ID field is required for Terraform Framework acceptance testing.
			"id": schema.StringAttribute{
				Description: "Unique identifier of the Access Policy.",
				Computed:    true,
				Validators: []validator.String{
					validators.UUIDRegexValidation(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name for the Access Policy.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("Placeholder"),
			},
			"is_active": schema.BoolAttribute{
				Description: "Active/Inactive status of the Access Policy.",
				Optional:    true,
				Computed:    true,
			},
			"client_workload": schema.StringAttribute{
				Description: "Client workload ID configured in the Access Policy.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validators.UUIDRegexValidation(),
				},
			},
			"trust_providers": schema.SetAttribute{
				Description: "Set of Trust Providers to enforce on the Access Policy.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(validators.UUIDRegexValidation()),
				},
				Default: setdefault.StaticValue(
					types.SetValueMust(types.StringType, []attr.Value{}),
				),
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},
			},
			"access_conditions": schema.SetAttribute{
				Description: "Set of Access Conditions to enforce on the Access Policy.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(validators.UUIDRegexValidation()),
				},
				Default: setdefault.StaticValue(
					types.SetValueMust(types.StringType, []attr.Value{}),
				),
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},
			},
			"credential_provider": schema.StringAttribute{
				Description: "Credential Provider ID configured in the Access Policy.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validators.UUIDRegexValidation(),
				},
			},
			"credential_providers": schema.SetNestedAttribute{
				Description: "Set of Credential Providers to enforce on the Access Policy.",
				Optional:    true,
				Computed:    true,
				Default: setdefault.StaticValue(
					types.SetNull(types.ObjectType{AttrTypes: map[string]attr.Type{
						"credential_provider_id": types.StringType,
						"mapping_type":           types.StringType,
						"header_name":            types.StringType,
						"header_value":           types.StringType,
						"httpbody_field_path":    types.StringType,
						"httpbody_field_value":   types.StringType,
						"account_name":           types.StringType,
					}}),
				),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"credential_provider_id": schema.StringAttribute{
							Description: "ID of associated Credential Provider.",
							Required:    true,
							Validators: []validator.String{
								validators.UUIDRegexValidation(),
							},
						},
						"mapping_type": schema.StringAttribute{
							Description: "Mapping type for the credential provider.",
							Required:    true,
						},
						"header_name": schema.StringAttribute{
							Description: "Name of the header for the credential provider.",
							Optional:    true,
							Computed:    true,
						},
						"header_value": schema.StringAttribute{
							Description: "Value of the header for the credential provider.",
							Optional:    true,
							Computed:    true,
						},
						"httpbody_field_path": schema.StringAttribute{
							Description: "Field path in the HTTP body for the credential provider.",
							Optional:    true,
							Computed:    true,
						},
						"httpbody_field_value": schema.StringAttribute{
							Description: "Field value in the HTTP body for the credential provider.",
							Optional:    true,
							Computed:    true,
						},
						"account_name": schema.StringAttribute{
							Description: "Name of the Snowflake account for the credential provider.",
							Optional:    true,
							Computed:    true,
						},
					},
				},
				Validators: []validator.Set{
					validators.NewCredentialProviderMappingValidator(),
				},
			},
			"server_workload": schema.StringAttribute{
				Description: "Server workload ID configured in the Access Policy.",
				Required:    true,
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

// Create creates the resource and sets the initial Terraform state.
func (r *accessPolicyResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	// Retrieve values from plan
	var plan models.AccessPolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	credentialProvider := !plan.CredentialProvider.IsNull() && !plan.CredentialProvider.IsUnknown() && plan.CredentialProvider.ValueString() != ""
	credentialProviders := len(plan.CredentialProviders) > 0

	if (credentialProvider && credentialProviders) || (!credentialProvider && !credentialProviders) {
		resp.Diagnostics.AddError(
			"Error creating Access Policy",
			"Only one of 'credential_provider' or 'credential_providers' should be set.",
		)
		return
	}

	initialOrderOfCredentialProviders := make([]string, len(plan.CredentialProviders))
	for i, cp := range plan.CredentialProviders {
		initialOrderOfCredentialProviders[i] = cp.CredentialProviderId.ValueString()
	}
	// Generate API request body from plan
	policy := convertAccessPolicyModelToPolicyDTO(plan, nil)

	// Create new Access Policy
	accessPolicy, err := r.client.CreateAccessPolicyV2(policy, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating access policy",
			"Could not create access policy, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan = convertAccessPolicyDTOToModel(plan, *accessPolicy)
	plan.CredentialProviders = sortCredentialProviders(
		plan.CredentialProviders,
		initialOrderOfCredentialProviders,
	)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *accessPolicyResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	// Get current state
	var state models.AccessPolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	initialOrderOfCredentialProviders := make([]string, len(state.CredentialProviders))

	for i, cp := range state.CredentialProviders {
		initialOrderOfCredentialProviders[i] = cp.CredentialProviderId.ValueString()
	}

	// Get refreshed policy value from Aembit
	accessPolicy, err, notFound := r.client.GetAccessPolicyV2(state.ID.ValueString(), nil)

	if err != nil {
		resp.Diagnostics.AddWarning(
			"Error reading Aembit Access Policy",
			"Could not read Aembit External ID from Terraform state "+state.ID.ValueString()+": "+err.Error(),
		)

		// If the resource is not found on Aembit Cloud, delete it locally
		if notFound {
			resp.State.RemoveResource(ctx)
		}
		return
	}

	// fetch mappings values individually
	credentialMappings, err, _ := r.client.GetAccessPolicyV2CredentialMappings(
		state.ID.ValueString(),
		nil,
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error retrieving credential mappings",
			"Could not get credential mappings, unexpected error: "+err.Error(),
		)
		return
	}

	state = convertAccessPolicyExternalDTOToModel(accessPolicy, credentialMappings)
	state.CredentialProviders = sortCredentialProviders(
		state.CredentialProviders,
		initialOrderOfCredentialProviders,
	)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *accessPolicyResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	// Get current state
	var state models.AccessPolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	credentialProvider := !state.CredentialProvider.IsNull() && !state.CredentialProvider.IsUnknown() && state.CredentialProvider.ValueString() != ""
	credentialProviders := len(state.CredentialProviders) > 0

	if (credentialProvider && credentialProviders) || (!credentialProvider && !credentialProviders) {
		resp.Diagnostics.AddError(
			"Error creating Access Policy",
			"Only one of 'credential_provider' or 'credential_providers' should be set.",
		)
		return
	}

	// Extract external ID from state
	externalID := state.ID.ValueString()

	// Retrieve values from plan
	var plan models.AccessPolicyResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	initialOrderOfCredentialProviders := make([]string, len(state.CredentialProviders))

	for i, cp := range state.CredentialProviders {
		initialOrderOfCredentialProviders[i] = cp.CredentialProviderId.ValueString()
	}

	// Generate API request body from plan
	policy := convertAccessPolicyModelToPolicyDTO(plan, &externalID)

	// Update Access Policy
	accessPolicy, err := r.client.UpdateAccessPolicyV2(policy, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating access policy",
			"Could not update access policy, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	state = convertAccessPolicyDTOToModel(plan, *accessPolicy)

	state.CredentialProviders = sortCredentialProviders(
		state.CredentialProviders,
		initialOrderOfCredentialProviders,
	)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *accessPolicyResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	// Retrieve values from state
	var state models.AccessPolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check if Access Policy is Active - if it is, disable it first
	if state.IsActive == types.BoolValue(true) {
		_, err := r.client.DisableAccessPolicyV2(state.ID.ValueString(), nil)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error disabling Access Policy",
				"Could not disable Access Policy, unexpected error: "+err.Error(),
			)
			return
		}
	}

	// Delete existing Access Policy
	_, err := r.client.DeleteAccessPolicyV2(ctx, state.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Access Policy",
			"Could not delete access policy, unexpected error: "+err.Error(),
		)
		return
	}
}

// Imports an existing resource by passing externalId.
func (r *accessPolicyResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	// Retrieve import externalId and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func convertAccessPolicyModelToPolicyDTO(
	model models.AccessPolicyResourceModel,
	externalID *string,
) aembit.CreatePolicyDTO {
	var policy aembit.CreatePolicyDTO
	policy.EntityDTO = aembit.EntityDTO{
		Name:     model.Name.ValueString(),
		IsActive: model.IsActive.ValueBool(),
	}
	if externalID != nil {
		policy.ExternalID = *externalID
	}

	policy.ClientWorkload = model.ClientWorkload.ValueString()
	policy.ServerWorkload = model.ServerWorkload.ValueString()
	policy.TrustProviders = make([]string, len(model.TrustProviders))
	policy.AccessConditions = make([]string, len(model.AccessConditions))

	if len(model.TrustProviders) > 0 {
		for i, trustProvider := range model.TrustProviders {
			policy.TrustProviders[i] = trustProvider.ValueString()
		}
	}
	if len(model.AccessConditions) > 0 {
		for i, accessConditions := range model.AccessConditions {
			policy.AccessConditions[i] = accessConditions.ValueString()
		}
	}

	// populate CredentialProviders statically if only credentialProvider provided
	if len(model.CredentialProvider.ValueString()) > 0 {
		policy.CredentialProviders = make([]aembit.PolicyCredentialMappingDTO, 1)
		policy.CredentialProviders[0] = aembit.PolicyCredentialMappingDTO{
			CredentialProviderId: model.CredentialProvider.ValueString(),
			MappingType:          "None",
			AccountName:          "",
			HeaderName:           "",
			HeaderValue:          "",
			HttpbodyFieldPath:    "",
			HttpbodyFieldValue:   "",
		}
	} else {
		policy.CredentialProviders = make([]aembit.PolicyCredentialMappingDTO, len(model.CredentialProviders))

		if len(model.CredentialProviders) > 0 {
			for i, credentialProvider := range model.CredentialProviders {
				policy.CredentialProviders[i] = aembit.PolicyCredentialMappingDTO{
					CredentialProviderId: credentialProvider.CredentialProviderId.ValueString(),
					MappingType:          credentialProvider.MappingType.ValueString(),
					AccountName:          credentialProvider.AccountName.ValueString(),
					HeaderName:           credentialProvider.HeaderName.ValueString(),
					HeaderValue:          credentialProvider.HeaderValue.ValueString(),
					HttpbodyFieldPath:    credentialProvider.HttpbodyFieldPath.ValueString(),
					HttpbodyFieldValue:   credentialProvider.HttpbodyFieldValue.ValueString(),
				}
			}
		}
	}

	return policy
}

func convertAccessPolicyDTOToModel(
	plan models.AccessPolicyResourceModel,
	dto aembit.CreatePolicyDTO,
) models.AccessPolicyResourceModel {
	var model models.AccessPolicyResourceModel
	model.ID = types.StringValue(dto.ExternalID)
	model.Name = types.StringValue(dto.Name)
	model.IsActive = types.BoolValue(dto.IsActive)
	model.ClientWorkload = types.StringValue(dto.ClientWorkload)
	model.ServerWorkload = types.StringValue(dto.ServerWorkload)

	if len(dto.CredentialProviders) > 0 {
		model.CredentialProvider = types.StringValue(
			dto.CredentialProviders[0].CredentialProviderId,
		)

		// discard credentialProvider mappings if there is a single credential provider
		if len(dto.CredentialProviders) > 1 && dto.CredentialProviders[0].MappingType != "None" {
			model.CredentialProviders = make(
				[]*models.PolicyCredentialMappingModel,
				len(dto.CredentialProviders),
			)

			for i, credentialProvider := range dto.CredentialProviders {
				model.CredentialProviders[i] = &models.PolicyCredentialMappingModel{
					CredentialProviderId: types.StringValue(
						credentialProvider.CredentialProviderId,
					),
					MappingType:        types.StringValue(credentialProvider.MappingType),
					AccountName:        types.StringValue(""),
					HeaderName:         types.StringValue(""),
					HeaderValue:        types.StringValue(""),
					HttpbodyFieldPath:  types.StringValue(""),
					HttpbodyFieldValue: types.StringValue(""),
				}

				model.CredentialProviders[i].AccountName = types.StringValue(
					credentialProvider.AccountName,
				)
				model.CredentialProviders[i].HeaderName = types.StringValue(
					credentialProvider.HeaderName,
				)
				model.CredentialProviders[i].HeaderValue = types.StringValue(
					credentialProvider.HeaderValue,
				)
				model.CredentialProviders[i].HttpbodyFieldPath = types.StringValue(
					credentialProvider.HttpbodyFieldPath,
				)
				model.CredentialProviders[i].HttpbodyFieldValue = types.StringValue(
					credentialProvider.HttpbodyFieldValue,
				)
			}
		}
	}

	model.TrustProviders = make([]types.String, len(dto.TrustProviders))
	if len(dto.TrustProviders) > 0 {
		for i, trustProvider := range dto.TrustProviders {
			model.TrustProviders[i] = types.StringValue(trustProvider)
		}
	}
	model.AccessConditions = make([]types.String, len(dto.AccessConditions))
	if len(dto.AccessConditions) > 0 {
		for i, accessConditions := range dto.AccessConditions {
			model.AccessConditions[i] = types.StringValue(accessConditions)
		}
	}

	// Handle cases where the plan does not specify TrustProviders or AccessConditions at all
	if plan.TrustProviders == nil {
		model.TrustProviders = nil
	}
	if plan.AccessConditions == nil {
		model.AccessConditions = nil
	}

	return model
}

func convertAccessPolicyExternalDTOToModel(
	dto aembit.GetPolicyDTO,
	credentialMappings []aembit.PolicyCredentialMappingDTO,
) models.AccessPolicyResourceModel {
	var model models.AccessPolicyResourceModel
	model.ID = types.StringValue(dto.ExternalID)
	model.Name = types.StringValue(dto.Name)
	model.IsActive = types.BoolValue(dto.IsActive)
	model.ClientWorkload = types.StringValue(dto.ClientWorkload.ExternalID)
	model.ServerWorkload = types.StringValue(dto.ServerWorkload.ExternalID)

	if len(dto.CredentialProviders) > 0 {
		model.CredentialProvider = types.StringValue(dto.CredentialProviders[0].ExternalID)

		// discard credentialProvider mappings if there is a single credential provider
		if len(credentialMappings) > 1 && credentialMappings[0].MappingType != "None" {
			model.CredentialProviders = make(
				[]*models.PolicyCredentialMappingModel,
				len(dto.CredentialProviders),
			)

			for i, credentialProvider := range dto.CredentialProviders {
				// find related mapping
				relatedMapping := aembit.PolicyCredentialMappingDTO{}

				for _, mapping := range credentialMappings {
					if mapping.CredentialProviderId == credentialProvider.ExternalID {
						relatedMapping = mapping
						break
					}
				}

				model.CredentialProviders[i] = &models.PolicyCredentialMappingModel{
					CredentialProviderId: types.StringValue(credentialProvider.ExternalID),
					MappingType:          types.StringValue(relatedMapping.MappingType),
					AccountName:          types.StringValue(relatedMapping.AccountName),
					HeaderName:           types.StringValue(relatedMapping.HeaderName),
					HeaderValue:          types.StringValue(relatedMapping.HeaderValue),
					HttpbodyFieldPath:    types.StringValue(relatedMapping.HttpbodyFieldPath),
					HttpbodyFieldValue:   types.StringValue(relatedMapping.HttpbodyFieldValue),
				}
			}
		}
	}

	model.TrustProviders = make([]types.String, len(dto.TrustProviders))
	if len(dto.TrustProviders) > 0 {
		for i, trustProvider := range dto.TrustProviders {
			model.TrustProviders[i] = types.StringValue(trustProvider.ExternalID)
		}
	}
	model.AccessConditions = make([]types.String, len(dto.AccessConditions))
	if len(dto.AccessConditions) > 0 {
		for i, accessConditions := range dto.AccessConditions {
			model.AccessConditions[i] = types.StringValue(accessConditions.ExternalID)
		}
	}

	return model
}

func sortCredentialProviders(
	credentialProviders []*models.PolicyCredentialMappingModel,
	initialOrderOfCredentialProviders []string,
) []*models.PolicyCredentialMappingModel {
	finalOrderOfCredentialProviders := make(
		[]*models.PolicyCredentialMappingModel,
		len(credentialProviders),
	)

	// make sure the order is correct if there are more than 1 credential providers
	if len(credentialProviders) == len(initialOrderOfCredentialProviders) &&
		len(finalOrderOfCredentialProviders) > 1 {
		for i := 0; i < len(credentialProviders); i++ {
			for _, cp := range credentialProviders {
				if cp.CredentialProviderId.ValueString() == initialOrderOfCredentialProviders[i] {
					finalOrderOfCredentialProviders[i] = cp
					break
				}
			}
		}
		return finalOrderOfCredentialProviders
	}

	// preserve the incoming order
	return credentialProviders
}

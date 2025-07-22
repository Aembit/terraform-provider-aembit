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
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &identityProviderResource{}
	_ resource.ResourceWithConfigure   = &identityProviderResource{}
	_ resource.ResourceWithImportState = &identityProviderResource{}
)

func NewIdentityProviderResource() resource.Resource {
	return &identityProviderResource{}
}

// identityProviderResource is the resource implementation.
type identityProviderResource struct {
	client *aembit.CloudClient
}

// Metadata returns the resource type name.
func (r *identityProviderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_identity_provider"
}

// Configure adds the provider configured client to the resource.
func (r *identityProviderResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = resourceConfigure(req, resp)
}

// Schema defines the schema for the resource.
func (r *identityProviderResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			// ID field is required for Terraform Framework acceptance testing.
			"id": schema.StringAttribute{
				Description: "Unique identifier of the Identity Provider.",
				Computed:    true,
				Validators: []validator.String{
					validators.UUIDRegexValidation(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name for the Identity Provider.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Description: "Description for the Identity Provider.",
				Optional:    true,
				Computed:    true,
			},
			"is_active": schema.BoolAttribute{
				Description: "Active status of the Identity Provider.",
				Optional:    true,
				Computed:    true,
			},
			"tags": schema.MapAttribute{
				Description: "Tags are key-value pairs.",
				ElementType: types.StringType,
				Optional:    true,
			},
			"entity_id": schema.StringAttribute{
				Description: "Entity ID of the Identity Provider.",
				Required:    true,
			},
			"metadata_url": schema.StringAttribute{
				Description: "URL pointing to the metadata for the Identity Provider.",
				Optional:    true,
			},
			"metadata_xml": schema.StringAttribute{
				Description: "XML containing the metadata for the Identity Provider.",
				Optional:    true,
			},
			"saml_statement_role_mappings": schema.SingleNestedAttribute{
				Description: "Mapping between SAML attributes for the Identity Provider and Aembit user roles. This set of attributes is used to assign Aembit Roles to users during automatic user creation during the SSO flow.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"attribute_name": schema.StringAttribute{
						Description: "SAML attribute name.",
						Required:    true,
					},
					"attribute_value": schema.StringAttribute{
						Description: "SAML attribute value.",
						Required:    true,
					},
					"roles_ids": schema.SetAttribute{
						Description: "List of Aembit Role Identifiers to be assigned to a user.",
						Required:    true,
					},
				},
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *identityProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan models.IdentityProviderDataSourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var dto aembit.IdentityProviderDTO = convertIdentityProviderModelToDTO(ctx, plan, nil)

	// Create new identity provider
	idp, err := r.client.CreateIdentityProvider(dto, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Identity Provider",
			"Could not create Identity Provider, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan = convertIdentityProviderDTOToModel(ctx, *idp, plan)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *identityProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state models.IdentityProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	idpDto, err, notFound := r.client.GetIdentityProvider(state.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Error reading Identity Provider",
			"Could not read Aembit Identity Provider with External ID from Terraform state "+state.ID.ValueString()+": "+err.Error(),
		)

		// If the resource is not found on Aembit Cloud, delete it locally
		if notFound {
			resp.State.RemoveResource(ctx)
		}
		return
	}

	state = convertIdentityProviderDTOToModel(ctx, idpDto, state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *identityProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get current state
	var state identityProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract external ID from state
	externalID := state.ID.ValueString()

	// Retrieve values from plan
	var plan identityProviderResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var dto aembit.IdentityProviderDTO = convertIdentityProviderModelToDTO(ctx, plan, &externalID)

	// Update Identity Provider
	idpDto, err := r.client.UpdateIdentityProvider(dto, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Identity Provider",
			"Could not update Identity Provider, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	state = convertIdentityProviderDTOToModel(ctx, *idpDto, plan)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *identityProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state identityProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check if Identity Provider is Active - if it is, disable it first
	if state.IsActive == types.BoolValue(true) {
		_, err := r.client.DisableIdentityProvider(state.ID.ValueString(), nil)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error disabling Identity Provider",
				"Could not disable Identity Provider, unexpected error: "+err.Error(),
			)
			return
		}
	}

	// Delete existing Identity Provider
	_, err := r.client.DeleteIdentityProvider(ctx, state.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Identity Provider",
			"Could not delete Identity Provider, unexpected error: "+err.Error(),
		)
		return
	}
}

// Imports an existing resource by passing externalId.
func (r *identityProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import externalId and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func convertIdentityProviderModelToDTO(ctx context.Context, model models.IdentityProviderResourceModel, externalID *string) aembit.IdentityProviderDTO {
	var identityProvider aembit.IdentityProviderDTO
	identityProvider.EntityDTO = aembit.EntityDTO{
		Name:        model.Name.ValueString(),
		Description: model.Description.ValueString(),
		IsActive:    model.IsActive.ValueBool(),
	}
	if len(model.Tags.Elements()) > 0 {
		tagsMap := make(map[string]string)
		_ = model.Tags.ElementsAs(ctx, &tagsMap, true)

		for key, value := range tagsMap {
			identityProvider.Tags = append(identityProvider.Tags, aembit.TagDTO{
				Key:   key,
				Value: value,
			})
		}
	}

	if externalID != nil {
		identityProvider.EntityDTO.ExternalID = *externalID
	}

	identityProvider.EntityId = model.EntityId.ValueString()
	identityProvider.MetadataUrl = model.MetadataUrl.ValueString()
	identityProvider.MetadataXml = model.MetadataXml.ValueString()

	for _, mapping := range model.SamlStatementRoleMappings {
		mappingDto := aembit.SamlStatementRoleMappingDTO{
			AttributeName:  mapping.AttributeName.ValueString(),
			AttributeValue: mapping.AttributeValue.ValueString(),
		}
		for _, roleId := range mapping.Roles {
			mappingDto.RoleExternalId = roleId.ValueString()
		}
		identityProvider.SamlStatementRoleMappings = append(identityProvider.SamlStatementRoleMappings, mappingDto)
	}

	return identityProvider
}

func convertIdentityProviderDTOToModel(ctx context.Context, dto aembit.IdentityProviderDTO, state identityProviderResourceModel) identityProviderResourceModel {
	var model identityProviderResourceModel
	model.ID = types.StringValue(dto.EntityDTO.ExternalID)
	model.Name = types.StringValue(dto.EntityDTO.Name)
	model.Description = types.StringValue(dto.EntityDTO.Description)
	model.IsActive = types.BoolValue(dto.EntityDTO.IsActive)
	model.Tags = newTagsModel(ctx, dto.EntityDTO.Tags)

	model.EntityId = types.StringValue(dto.EntityId)
	model.MetadataUrl = types.StringValue(dto.MetadataUrl)
	model.MetadataXml = types.StringValue(dto.MetadataXml)

	//convert the mapping array from flat to unflatten form
	tempMap := make(map[string]samlStatementRoleMappings)
	for _, mapping := range dto.SamlStatementRoleMappings {
		key := mapping.AttributeName + mapping.AttributeValue
		if newItem, exists := tempMap[key]; exists {
			newItem.Roles = append(newItem.Roles, types.StringValue(mapping.RoleExternalId))
			tempMap[key] = newItem
		} else {
			tempMap[key] = samlStatementRoleMappings{
				AttributeName:  types.StringValue(mapping.AttributeName),
				AttributeValue: types.StringValue(mapping.AttributeValue),
				Roles:          []types.String{types.StringValue(mapping.RoleExternalId)},
			}
		}
	}
	model.SamlStatementRoleMappings = make([]samlStatementRoleMappings, 0, len(tempMap))
	for _, item := range tempMap {
		model.SamlStatementRoleMappings = append(model.SamlStatementRoleMappings, item)
	}

	return model
}

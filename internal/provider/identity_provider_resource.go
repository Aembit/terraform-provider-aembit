package provider

import (
	"context"
	"encoding/base64"
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
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &identityProviderResource{}
	_ resource.ResourceWithConfigure   = &identityProviderResource{}
	_ resource.ResourceWithImportState = &identityProviderResource{}
	_ resource.ResourceWithModifyPlan  = &identityProviderResource{}
)

func NewIdentityProviderResource() resource.Resource {
	return &identityProviderResource{}
}

// identityProviderResource is the resource implementation.
type identityProviderResource struct {
	client *aembit.CloudClient
}

// Metadata returns the resource type name.
func (r *identityProviderResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_identity_provider"
}

// Configure adds the provider configured client to the resource.
func (r *identityProviderResource) Configure(
	_ context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	r.client = resourceConfigure(req, resp)
}

// Schema defines the schema for the resource.
func (r *identityProviderResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
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
			"tags":     TagsMapAttribute(),
			"tags_all": TagsAllMapAttribute(),
			"sso_statement_role_mappings": schema.SetNestedAttribute{
				Description: "Mapping between SAML attributes for the Identity Provider and Aembit user roles. This set of attributes is used to assign Aembit Roles to users during automatic user creation during the SSO flow.",
				Optional:    true,
				Validators: []validator.Set{
					setvalidator.SizeAtLeast(1),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"attribute_name": schema.StringAttribute{
							Description: "SAML attribute name.",
							Required:    true,
						},
						"attribute_value": schema.StringAttribute{
							Description: "SAML attribute value.",
							Required:    true,
						},
						"roles": schema.SetAttribute{
							Description: "List of Aembit Role Identifiers to be assigned to a user.",
							ElementType: types.StringType,
							Required:    true,
							Validators: []validator.Set{
								setvalidator.SizeAtLeast(1),
								setvalidator.ValueStringsAre(validators.UUIDRegexValidation()),
							},
						},
					},
				},
			},
			"saml": schema.SingleNestedAttribute{
				Description: "SAML type Identity Provider configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"metadata_url": schema.StringAttribute{
						Description: "URL pointing to the metadata for the Identity Provider.",
						Optional:    true,
						Computed:    true,
					},
					"metadata_xml": schema.StringAttribute{
						Description: "XML containing the metadata for the Identity Provider.",
						Optional:    true,
						Computed:    true,
					},
					"service_provider_entity_id": schema.StringAttribute{
						Description: "The unique identifier (Entity ID) for the SAML Service Provider.",
						Computed:    true,
					},
					"service_provider_sso_url": schema.StringAttribute{
						Description: "The Single Sign-On (SSO) endpoint URL of the Service Provider.",
						Computed:    true,
					},
				},
			},
			"oidc": schema.SingleNestedAttribute{
				Description: "OIDC type Identity Provider configuration.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"oidc_base_url": schema.StringAttribute{
						Description: "The base URL of the OIDC Identity Provider.",
						Required:    true,
					},
					"client_id": schema.StringAttribute{
						Description: "The client identifier registered with the OIDC provider.",
						Required:    true,
					},
					"scopes": schema.StringAttribute{
						Description: "A space-separated list of OIDC scopes to request during authentication (e.g., 'openid profile email').",
						Required:    true,
					},
					"client_secret": schema.StringAttribute{
						Description: "The client secret associated with the OIDC client.",
						Optional:    true,
						Sensitive:   true,
					},
					"auth_type": schema.StringAttribute{
						Description: "Authentication method. Possible values are: \n" +
							"\t* `ClientSecret`\n" +
							"\t* `KeyPair`\n",
						Required: true,
						Validators: []validator.String{
							stringvalidator.OneOf([]string{
								"ClientSecret",
								"KeyPair",
							}...),
						},
					},
					"pcke_required": schema.BoolAttribute{
						Description: "Indicates whether Proof Key for Code Exchange (PKCE) is required during the OIDC authorization flow.",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(true),
					},
					"aembit_redirect_url": schema.StringAttribute{
						Description: "The redirect URI registered with the OIDC provider.",
						Computed:    true,
					},
					"aembit_jwks_url": schema.StringAttribute{
						Description: "The URL where the OIDC provider's JSON Web Key Set (JWKS) can be retrieved.",
						Computed:    true,
					},
				},
			},
		},
	}
}

func (r *identityProviderResource) ModifyPlan(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
) {
	modifyPlanForTagsAll(ctx, req, resp, r.client.DefaultTags)
}

func ValidatePlan(
	diagnostics *diag.Diagnostics,
	plan *models.IdentityProviderResourceModel,
) {
	if plan.Oidc != nil && plan.Oidc.AuthType.ValueString() == "ClientSecret" &&
		plan.Oidc.ClientSecret.IsNull() {

		diagnostics.AddError(
			"Error creating Identity Provider",
			"When auth_type is 'ClientSecret', the 'ClientSecret' attribute must be set.",
		)
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *identityProviderResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	// Retrieve values from plan
	var plan models.IdentityProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)

	ValidatePlan(&diags, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var dto = convertIdentityProviderModelToDTO(ctx, plan, nil, r.client.DefaultTags)

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
	plan = convertIdentityProviderDTOToModel(ctx, &plan, *idp)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *identityProviderResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
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

	state = convertIdentityProviderDTOToModel(ctx, &state, idpDto)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *identityProviderResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	// Get current state
	var state models.IdentityProviderResourceModel
	diags := req.State.Get(ctx, &state)

	ValidatePlan(&diags, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Extract external ID from state
	externalID := state.ID.ValueString()

	// Retrieve values from plan
	var plan models.IdentityProviderResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	var dto = convertIdentityProviderModelToDTO(ctx, plan, &externalID, r.client.DefaultTags)

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
	state = convertIdentityProviderDTOToModel(ctx, &plan, *idpDto)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *identityProviderResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	// Retrieve values from state
	var state models.IdentityProviderResourceModel
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
func (r *identityProviderResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	// Retrieve import externalId and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func convertIdentityProviderModelToDTO(
	ctx context.Context,
	model models.IdentityProviderResourceModel,
	externalID *string,
	defaultTags map[string]string,
) aembit.IdentityProviderDTO {
	var identityProvider aembit.IdentityProviderDTO
	identityProvider.EntityDTO = aembit.EntityDTO{
		Name:        model.Name.ValueString(),
		Description: model.Description.ValueString(),
		IsActive:    model.IsActive.ValueBool(),
	}

	if externalID != nil {
		identityProvider.ExternalID = *externalID
	}

	if model.Saml != nil {
		identityProvider.Type = "SAMLv2"
		identityProvider.MetadataUrl = model.Saml.MetadataUrl.ValueString()
		identityProvider.MetadataXml = base64.StdEncoding.EncodeToString(
			[]byte(strings.TrimRight(model.Saml.MetadataXml.ValueString(), "\n\r")),
		)
	}

	if model.Oidc != nil {
		identityProvider.Type = "OIDCv1"
		identityProvider.OidcBaseUrl = model.Oidc.OidcBaseUrl.ValueString()
		identityProvider.ClientId = model.Oidc.ClientId.ValueString()
		identityProvider.Scopes = model.Oidc.Scopes.ValueString()
		identityProvider.ClientSecret = model.Oidc.ClientSecret.ValueString()
		identityProvider.PkceRequired = model.Oidc.PkceRequired.ValueBool()
		identityProvider.AuthType = model.Oidc.AuthType.ValueString()
	}

	for _, mapping := range model.SsoStatementRoleMappings {
		mappingDto := aembit.SsoStatementRoleMappingDTO{
			AttributeName:  mapping.AttributeName.ValueString(),
			AttributeValue: mapping.AttributeValue.ValueString(),
		}
		for _, roleId := range mapping.Roles {
			mappingDto.RoleExternalId = roleId.ValueString()
			identityProvider.SsoStatementRoleMappings = append(
				identityProvider.SsoStatementRoleMappings,
				mappingDto,
			)
		}
	}

	identityProvider.Tags = collectAllTagsDto(ctx, defaultTags, model.Tags)
	return identityProvider
}

func convertIdentityProviderDTOToModel(
	ctx context.Context,
	planModel *models.IdentityProviderResourceModel,
	dto aembit.IdentityProviderDTO,
) models.IdentityProviderResourceModel {
	var model models.IdentityProviderResourceModel
	model.ID = types.StringValue(dto.ExternalID)
	model.Name = types.StringValue(dto.Name)
	model.Description = types.StringValue(dto.Description)
	model.IsActive = types.BoolValue(dto.IsActive)

	// handle tags
	model.Tags = newTagsModelFromPlan(ctx, planModel.Tags)
	model.TagsAll = newTagsModel(ctx, dto.Tags)

	// Set the objects to null to begin with
	model.Saml = nil
	model.Oidc = nil

	// Now fill in the objects based on the Identity Provider type
	switch dto.Type {
	case "SAMLv2":
		model.Saml = &models.IdentityProviderSamlModel{}

		if dto.MetadataUrl == "" {
			model.Saml.MetadataUrl = types.StringNull()
		} else {
			model.Saml.MetadataUrl = types.StringValue(dto.MetadataUrl)
		}

		// Add a carriage return if the original plan had one, to avoid diffs on re-read
		if planModel != nil && planModel.Saml != nil &&
			strings.HasSuffix(planModel.Saml.MetadataXml.ValueString(), "\n") &&
			!strings.HasSuffix(dto.MetadataXml, "\n") {
			dto.MetadataXml += "\n"
		}
		model.Saml.MetadataXml = types.StringValue(dto.MetadataXml)
		model.Saml.ServiceProviderEntityId = types.StringValue(dto.ServiceProviderEntityId)
		model.Saml.ServiceProviderSsoUrl = types.StringValue(dto.ServiceProviderSsoUrl)
	case "OIDCv1":
		model.Oidc = &models.IdentityProviderOidcModel{}

		model.Oidc.ClientId = types.StringValue(dto.ClientId)
		model.Oidc.OidcBaseUrl = types.StringValue(dto.OidcBaseUrl)
		model.Oidc.Scopes = types.StringValue(dto.Scopes)
		model.Oidc.AuthType = types.StringValue(dto.AuthType)
		model.Oidc.PkceRequired = types.BoolValue(dto.PkceRequired)

		model.Oidc.AembitRedirectUrl = types.StringValue(dto.AembitRedirectUrl)
		model.Oidc.AembitJwksUrl = types.StringValue(dto.AembitJwksUrl)

		model.Oidc.ClientSecret = types.StringNull()
		if planModel.Oidc != nil {
			model.Oidc.ClientSecret = planModel.Oidc.ClientSecret
		}
	}

	//convert the mapping array from flat to unflatten form
	if len(dto.SsoStatementRoleMappings) == 0 {
		model.SsoStatementRoleMappings = nil
		return model
	}

	tempMap := make(map[string]models.SsoStatementRoleMappings)
	for _, mapping := range dto.SsoStatementRoleMappings {
		key := mapping.AttributeName + mapping.AttributeValue
		if newItem, exists := tempMap[key]; exists {
			newItem.Roles = append(newItem.Roles, types.StringValue(mapping.RoleExternalId))
			tempMap[key] = newItem
		} else {
			tempMap[key] = models.SsoStatementRoleMappings{
				AttributeName:  types.StringValue(mapping.AttributeName),
				AttributeValue: types.StringValue(mapping.AttributeValue),
				Roles:          []types.String{types.StringValue(mapping.RoleExternalId)},
			}
		}
	}
	model.SsoStatementRoleMappings = make([]models.SsoStatementRoleMappings, 0, len(tempMap))
	for _, item := range tempMap {
		model.SsoStatementRoleMappings = append(model.SsoStatementRoleMappings, item)
	}

	return model
}

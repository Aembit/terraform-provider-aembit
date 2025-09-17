package provider

import (
	"context"
	"terraform-provider-aembit/internal/provider/models"
	"terraform-provider-aembit/internal/provider/validators"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &identityProviderDataSource{}
	_ datasource.DataSourceWithConfigure = &identityProviderDataSource{}
)

func NewIdentityProviderDataSource() datasource.DataSource {
	return &identityProviderDataSource{}
}

type identityProviderDataSource struct {
	client *aembit.CloudClient
}

func (r *identityProviderDataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_identity_providers"
}

func (r *identityProviderDataSource) Configure(
	_ context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
	r.client = datasourceConfigure(req, resp)
}

func (r *identityProviderDataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Data source for Aembit Identity Providers.",
		Attributes: map[string]schema.Attribute{
			"identity_providers": schema.ListNestedAttribute{
				Description: "List of Identity Providers.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
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
						"saml_statement_role_mappings": schema.SetNestedAttribute{
							Description: "Mapping between SAML attributes for the Identity Provider and Aembit user roles. This set of attributes is used to assign Aembit Roles to users during automatic user creation during the SSO flow.",
							Optional:    true,
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
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *identityProviderDataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	// Get current state
	var state models.IdentityProviderDataSourceModel
	req.Config.Get(ctx, &state)

	idps, err := r.client.GetIdentityProviders(nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading Aembit Identity Providers",
			"Could not read Aembit Identity Providers: "+err.Error(),
		)
		return
	}

	for _, idp := range idps {
		idpState := convertIdentityProviderDTOToModel(ctx, nil, idp)

		state.IdentityProviders = append(state.IdentityProviders, idpState)
	}

	// Set refreshed state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

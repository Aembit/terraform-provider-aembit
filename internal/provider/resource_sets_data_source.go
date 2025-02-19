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
	_ datasource.DataSource              = &resourceSetsDataSource{}
	_ datasource.DataSourceWithConfigure = &resourceSetsDataSource{}
)

// NewResourceSetsDataSource is a helper function to simplify the provider implementation.
func NewResourceSetsDataSource() datasource.DataSource {
	return &resourceSetsDataSource{}
}

// resourceSetsDataSource is the data source implementation.
type resourceSetsDataSource struct {
	client *aembit.CloudClient
}

// Configure adds the provider configured client to the data source.
func (d *resourceSetsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = datasourceConfigure(req, resp)
}

// Metadata returns the data source type name.
func (d *resourceSetsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resource_sets"
}

// Schema defines the schema for the resource.
func (d *resourceSetsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to get information about all Aembit ResourceSets.",
		Attributes: map[string]schema.Attribute{
			"resource_sets": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						// ID field is required for Terraform Framework acceptance testing.
						"id": schema.StringAttribute{
							Description: "Unique identifier of the ResourceSet.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "User-provided name of the ResourceSet.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "User-provided description of the ResourceSet.",
							Computed:    true,
						},
						"roles": schema.SetAttribute{
							Description: "IDs of the Roles to associate with this ResourceSet.",
							Computed:    true,
							ElementType: types.StringType,
						},
					},
				},
			},
		},
	}
}

// Read checks for the datasource ID and retrieves the associated ResourceSet.
func (d *resourceSetsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state models.ResourceSetsDataModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceSets, err := d.client.GetResourceSets(nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Aembit Resource Sets",
			err.Error(),
		)
		return
	}

	// Map response body to model
	state = convertResourceSetsDTOToModel(ctx, resourceSets)

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// DTO to Model conversion methods.
func convertResourceSetsDTOToModel(_ context.Context, dto []aembit.ResourceSetDTO) models.ResourceSetsDataModel {
	var model models.ResourceSetsDataModel = models.ResourceSetsDataModel{
		ResourceSets: make([]models.ResourceSetResourceModel, len(dto)),
	}
	for index, resourceSet := range dto {
		model.ResourceSets[index] = models.ResourceSetResourceModel{
			ID:          types.StringValue(resourceSet.EntityDTO.ExternalID),
			Name:        types.StringValue(resourceSet.EntityDTO.Name),
			Description: types.StringValue(resourceSet.EntityDTO.Description),
			Roles:       make([]types.String, len(resourceSet.Roles)),
		}
		model.ResourceSets[index].Roles = make([]types.String, len(resourceSet.Roles))
		if len(resourceSet.Roles) > 0 {
			for i, role := range resourceSet.Roles {
				model.ResourceSets[index].Roles[i] = types.StringValue(role)
			}
		}
	}

	return model
}

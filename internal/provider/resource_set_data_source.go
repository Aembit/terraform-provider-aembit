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
	_ datasource.DataSource              = &resourceSetDataSource{}
	_ datasource.DataSourceWithConfigure = &resourceSetDataSource{}
)

// NewResourceSetDataSource is a helper function to simplify the provider implementation.
func NewResourceSetDataSource() datasource.DataSource {
	return &resourceSetDataSource{}
}

// resourceSetDataSource is the data source implementation.
type resourceSetDataSource struct {
	client *aembit.CloudClient
}

// Implementations for resourceSetDataSource (singular)

// Configure adds the provider configured client to the data source.
func (d *resourceSetDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = datasourceConfigure(req, resp)
}

// Metadata returns the data source type name.
func (d *resourceSetDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resource_set"
}

// Schema defines the schema for the resource.
func (d *resourceSetDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to get information about an Aembit ResourceSet.",
		Attributes: map[string]schema.Attribute{
			// ID field is required for Terraform Framework acceptance testing.
			"id": schema.StringAttribute{
				Description: "Unique identifier of the ResourceSet.",
				Required:    true,
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
				Description: "IDs of the Roles associated with this ResourceSet.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"standalone_certificate_authority": schema.StringAttribute{
				Description: "Standalone Certificate Authority ID configured for this ResourceSet.",
				Computed:    true,
				Optional:    true,
			},
		},
	}
}

// Read checks for the datasource ID and retrieves the associated ResourceSet.
func (d *resourceSetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state models.ResourceSetResourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceSet, err, _ := d.client.GetResourceSet(state.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Aembit Resource Sets",
			err.Error(),
		)
		return
	}

	// Map response body to model
	state = convertResourceSetDTOToModel(ctx, resourceSet)

	// Set state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

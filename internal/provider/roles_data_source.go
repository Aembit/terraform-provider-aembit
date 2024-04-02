package provider

import (
	"context"
	"fmt"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &rolesDataSource{}
	_ datasource.DataSourceWithConfigure = &rolesDataSource{}
)

// NewRolesDataSource is a helper function to simplify the provider implementation.
func NewRolesDataSource() datasource.DataSource {
	return &rolesDataSource{}
}

// rolesDataSource is the data source implementation.
type rolesDataSource struct {
	client *aembit.CloudClient
}

// Configure adds the provider configured client to the data source.
func (d *rolesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*aembit.CloudClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *aembit.CloudClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

// Metadata returns the data source type name.
func (d *rolesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_roles"
}

// Schema defines the schema for the resource.
func (d *rolesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an role.",
		Attributes: map[string]schema.Attribute{
			"roles": schema.ListNestedAttribute{
				Description: "List of roles.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						// ID field is required for Terraform Framework acceptance testing.
						"id": schema.StringAttribute{
							Description: "Unique identifier of the role.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "User-provided name of the role.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "User-provided description of the role.",
							Computed:    true,
						},
						"is_active": schema.BoolAttribute{
							Description: "Active/Inactive status of the role.",
							Computed:    true,
						},
						"tags": schema.MapAttribute{
							ElementType: types.StringType,
							Computed:    true,
						},
						"access_policies":             definePermissionAttribute("Access Policy", true),
						"client_workloads":            definePermissionAttribute("Client Workload", true),
						"trust_providers":             definePermissionAttribute("Trust Provider", true),
						"access_conditions":           definePermissionAttribute("Access Condition", true),
						"integrations":                definePermissionAttribute("Integration", true),
						"credential_providers":        definePermissionAttribute("Credential Provider", true),
						"server_workloads":            definePermissionAttribute("Server Workload", true),
						"agent_controllers":           definePermissionAttribute("Agent Controller", true),
						"access_authorization_events": definePermissionReadOnlyAttribute("Access Authorization Event", true),
						"audit_logs":                  definePermissionReadOnlyAttribute("Audit Log", true),
						"workload_events":             definePermissionReadOnlyAttribute("Workload Event", true),
						"users":                       definePermissionAttribute("User", true),
						"roles":                       definePermissionAttribute("Role", true),
						"log_streams":                 definePermissionAttribute("Log Stream", true),
						"identity_providers":          definePermissionAttribute("Identity Provider", true),
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *rolesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state rolesDataSourceModel

	roles, err := d.client.GetRoles(nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Aembit Trust Providers",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, role := range roles {
		roleState := convertRoleDTOToModel(ctx, role)
		state.Roles = append(state.Roles, roleState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

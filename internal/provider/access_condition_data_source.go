package provider

import (
	"context"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-aembit/internal/provider/models"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &accessConditionsDataSource{}
	_ datasource.DataSourceWithConfigure = &accessConditionsDataSource{}
)

// NewAccessConditionsDataSource is a helper function to simplify the provider implementation.
func NewAccessConditionsDataSource() datasource.DataSource {
	return &accessConditionsDataSource{}
}

// accessConditionsDataSource is the data source implementation.
type accessConditionsDataSource struct {
	client *aembit.CloudClient
}

// Configure adds the provider configured client to the data source.
func (d *accessConditionsDataSource) Configure(
	_ context.Context,
	req datasource.ConfigureRequest,
	resp *datasource.ConfigureResponse,
) {
	d.client = datasourceConfigure(req, resp)
}

// Metadata returns the data source type name.
func (d *accessConditionsDataSource) Metadata(
	_ context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_access_conditions"
}

// Schema defines the schema for the resource.
func (d *accessConditionsDataSource) Schema(
	_ context.Context,
	_ datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Description: "Manages an accessCondition.",
		Attributes: map[string]schema.Attribute{
			"access_conditions": schema.ListNestedAttribute{
				Description: "List of accessConditions.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						// ID field is required for Terraform Framework acceptance testing.
						"id": schema.StringAttribute{
							Description: "Unique identifier of the accessCondition.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "User-provided name of the accessCondition.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "User-provided description of the accessCondition.",
							Computed:    true,
						},
						"is_active": schema.BoolAttribute{
							Description: "Active/Inactive status of the accessCondition.",
							Computed:    true,
						},
						"tags": schema.MapAttribute{
							Description: "Tags are key-value pairs.",
							ElementType: types.StringType,
							Computed:    true,
						},
						"integration_id": schema.StringAttribute{
							Description: "ID of the Integration used by the Access Condition.",
							Computed:    true,
						},
						"wiz_conditions": schema.SingleNestedAttribute{
							Description: "Wiz Specific rules for the Access Condition.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"max_last_seen": schema.Int64Attribute{
									Required: true,
								},
								"container_cluster_connected": schema.BoolAttribute{Required: true},
							},
						},
						"crowdstrike_conditions": schema.SingleNestedAttribute{
							Description: "CrowdStrike Specific rules for the Access Condition.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"max_last_seen":       schema.Int64Attribute{Required: true},
								"match_hostname":      schema.BoolAttribute{Required: true},
								"match_serial_number": schema.BoolAttribute{Required: true},
								"prevent_rfm":         schema.BoolAttribute{Required: true},
							},
						},
						"geoip_conditions": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "Defines geographical conditions for filtering based on country and administrative subdivisions.",
							Attributes: map[string]schema.Attribute{
								"locations": schema.ListNestedAttribute{
									Required:    true,
									Description: "A list of geographical locations, each containing a country code and optional subdivisions.",
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"country_code": schema.StringAttribute{
												Description: "A list of two-letter country code identifiers (as defined by ISO 3166-1) to allow as part of the validation for this access condition.",
												Required:    true,
											},
											"subdivisions": schema.ListNestedAttribute{
												Description: "A list of subdivision identifiers (as defined by ISO 3166) to allow as part of the validation for this access condition.",
												Optional:    true,
												NestedObject: schema.NestedAttributeObject{
													Attributes: map[string]schema.Attribute{
														"subdivision_code": schema.StringAttribute{
															Description: "The subdivision identifier as defined by ISO 3166.",
															Required:    true,
														},
													},
												},
											},
										},
									},
								},
							},
						},
						"time_conditions": schema.SingleNestedAttribute{
							Computed:    true,
							Description: "Defines the conditions for scheduling based on time, including specific time slots and timezone settings for the Access Condition.",
							Attributes: map[string]schema.Attribute{
								"schedule": schema.ListNestedAttribute{
									Required: true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"start_time": schema.StringAttribute{
												Required:    true,
												Description: "The start time of the schedule in 24-hour format (HH:mm), e.g., '07:00' for 7:00 AM.",
											},
											"end_time": schema.StringAttribute{
												Required:    true,
												Description: "The end time of the schedule in 24-hour format (HH:mm), e.g., '18:00' for 6:00 PM.",
											},
											"day": schema.StringAttribute{
												Required:    true,
												Description: "Day of Week, for example: Tuesday",
											},
										},
									},
								},
								"timezone": schema.StringAttribute{
									Description: "Timezone value such as America/Chicago, Europe/Istanbul",
									Required:    true,
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
func (d *accessConditionsDataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var state models.AccessConditionsDataSourceModel

	accessConditions, err := d.client.GetAccessConditions(nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Aembit AccessConditions",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, accessCondition := range accessConditions {
		accessConditionState := convertAccessConditionDTOToModel(
			ctx,
			accessCondition,
			models.AccessConditionResourceModel{},
		)
		state.AccessConditions = append(state.AccessConditions, accessConditionState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

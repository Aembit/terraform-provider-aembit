package provider

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-aembit/internal/provider/models"
	"terraform-provider-aembit/internal/provider/validators"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &accessConditionResource{}
	_ resource.ResourceWithConfigure   = &accessConditionResource{}
	_ resource.ResourceWithImportState = &accessConditionResource{}
)

// NewAccessConditionResource is a helper function to simplify the provider implementation.
func NewAccessConditionResource() resource.Resource {
	return &accessConditionResource{}
}

// accessConditionResource is the resource implementation.
type accessConditionResource struct {
	client *aembit.CloudClient
}

// Metadata returns the resource type name.
func (r *accessConditionResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_access_condition"
}

// Configure adds the provider configured client to the resource.
func (r *accessConditionResource) Configure(
	_ context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	r.client = resourceConfigure(req, resp)
}

// Schema defines the schema for the resource.
func (r *accessConditionResource) Schema(
	_ context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			// ID field is required for Terraform Framework acceptance testing.
			"id": schema.StringAttribute{
				Description: "Unique identifier of the Access Condition.",
				Computed:    true,
				Validators: []validator.String{
					validators.UUIDRegexValidation(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name for the Access Condition.",
				Required:    true,
				Validators: []validator.String{
					validators.NameLengthValidation(),
				},
			},
			"description": schema.StringAttribute{
				Description: "Description for the Access Condition.",
				Optional:    true,
				Computed:    true,
			},
			"is_active": schema.BoolAttribute{
				Description: "Active status of the Access Condition.",
				Optional:    true,
				Computed:    true,
			},
			"tags": schema.MapAttribute{
				Description: "Tags are key-value pairs.",
				ElementType: types.StringType,
				Optional:    true,
			},
			"integration_id": schema.StringAttribute{
				Description: "Reference to the Integration used for this Access Condition.",
				Required:    true,
			},
			"wiz_conditions": schema.SingleNestedAttribute{
				Description: "Wiz Specific rules for the Access Condition.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"max_last_seen": schema.Int64Attribute{
						Description: "The maximum number of seconds since the managed Cluster was last seen by Wiz. Accepted range: 1-31449600 seconds",
						Required:    true,
						Validators: []validator.Int64{
							int64validator.Between(1, 31449600),
						},
					},
					"container_cluster_connected": schema.BoolAttribute{
						Description: "The condition requires that managed Clusters be defined as Container Cluster Connected by Wiz.",
						Required:    true,
					},
				},
			},
			"crowdstrike_conditions": schema.SingleNestedAttribute{
				Description: "CrowdStrike Specific rules for the Access Condition.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"max_last_seen": schema.Int64Attribute{
						Description: "The maximum number of seconds since the managed Cluster was last seen by CrowdStrike. Accepted range: 1-31449600 seconds",
						Required:    true,
						Validators: []validator.Int64{
							int64validator.Between(1, 31449600),
						},
					},
					"match_hostname": schema.BoolAttribute{
						Description: "The condition requires that managed hosts have a hostname which matches the CrowdStrike identified hostname.",
						Required:    true,
					},
					"match_serial_number": schema.BoolAttribute{
						Description: "The condition requires that managed hosts have a system serial number which matches the CrowdStrike identified serial number.",
						Required:    true,
					},
					"prevent_rfm": schema.BoolAttribute{
						Description: "The condition requires that managed hosts not be in CrowdStrike Reduced Functionality Mode.",
						Required:    true,
					},
					"match_mac_address": schema.BoolAttribute{
						Description: "The condition requires that managed hosts have a MAC address which matches the CrowdStrike identified MAC address.",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
					},
					"match_local_ip": schema.BoolAttribute{
						Description: "The condition requires that managed hosts have a local IP that matches the CrowdStrike-identified local or connection IP.",
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
					},
				},
			},
			"geoip_conditions": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"locations": schema.ListNestedAttribute{
						Required: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"country_code": schema.StringAttribute{
									Required:    true,
									Description: "A list of two-letter country code identifiers (as defined by ISO 3166-1) to allow as part of the validation for this access condition.",
								},
								"subdivisions": schema.ListNestedAttribute{
									Optional: true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"subdivision_code": schema.StringAttribute{
												Required:    true,
												Description: "A list of subdivision identifiers (as defined by ISO 3166) to allow as part of the validation for this access condition.",
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
				Optional:    true,
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
						Required:    true,
						Description: "Timezone value such as America/Chicago, Europe/Istanbul",
					},
				},
			},
		},
	}
}

// Configure validators to ensure that only one Access Condition type is specified.
func (r *accessConditionResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.ExactlyOneOf(
			path.MatchRoot("wiz_conditions"),
			path.MatchRoot("crowdstrike_conditions"),
			path.MatchRoot("geoip_conditions"),
			path.MatchRoot("time_conditions"),
		),
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *accessConditionResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	// Retrieve values from plan
	var plan models.AccessConditionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	dto, err := convertAccessConditionModelToDTO(ctx, plan, nil, r.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Access Condition",
			err.Error(),
		)
		return
	}

	// Create new AccessCondition
	accessCondition, err := r.client.CreateAccessCondition(dto, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Access Condition",
			"Could not create Access Condition, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan = convertAccessConditionDTOToModel(ctx, *accessCondition, plan)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *accessConditionResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	// Get current state
	var state models.AccessConditionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed trust value from Aembit
	accessCondition, err, notFound := r.client.GetAccessCondition(state.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Error reading Aembit Access Condition",
			"Could not read Aembit External ID from Terraform state "+state.ID.ValueString()+": "+err.Error(),
		)

		// If the resource is not found on Aembit Cloud, delete it locally
		if notFound {
			resp.State.RemoveResource(ctx)
		}
		return
	}

	state = convertAccessConditionDTOToModel(ctx, accessCondition, state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *accessConditionResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	// Get current state
	var state models.AccessConditionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract external ID from state
	externalID := state.ID.ValueString()

	// Retrieve values from plan
	var plan models.AccessConditionResourceModel
	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	dto, err := convertAccessConditionModelToDTO(ctx, plan, &externalID, r.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Access Condition",
			err.Error(),
		)
		return
	}

	// Update AccessCondition
	accessCondition, err := r.client.UpdateAccessCondition(dto, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Access Condition",
			"Could not update Access Condition, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	state = convertAccessConditionDTOToModel(ctx, *accessCondition, state)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *accessConditionResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	// Retrieve values from state
	var state models.AccessConditionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check if Access Condition is Active - if it is, disable it first
	if state.IsActive == types.BoolValue(true) {
		_, err := r.client.DisableAccessCondition(state.ID.ValueString(), nil)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error disabling Access Condition",
				"Could not disable Access Condition, unexpected error: "+err.Error(),
			)
			return
		}
	}

	// Delete existing AccessCondition
	_, err := r.client.DeleteAccessCondition(ctx, state.ID.ValueString(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting AccessCondition",
			"Could not delete Access Condition, unexpected error: "+err.Error(),
		)
		return
	}
}

// Imports an existing resource by passing externalId.
func (r *accessConditionResource) ImportState(
	ctx context.Context,
	req resource.ImportStateRequest,
	resp *resource.ImportStateResponse,
) {
	// Retrieve import externalId and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func convertAccessConditionModelToDTO(
	ctx context.Context,
	model models.AccessConditionResourceModel,
	externalID *string,
	client *aembit.CloudClient,
) (aembit.AccessConditionDTO, error) {
	var accessCondition aembit.AccessConditionDTO
	accessCondition.EntityDTO = aembit.EntityDTO{
		Name:        model.Name.ValueString(),
		Description: model.Description.ValueString(),
		IsActive:    model.IsActive.ValueBool(),
	}
	if externalID != nil {
		accessCondition.ExternalID = *externalID
	}

	if len(model.Tags.Elements()) > 0 {
		tagsMap := make(map[string]string)
		_ = model.Tags.ElementsAs(ctx, &tagsMap, true)

		for key, value := range tagsMap {
			accessCondition.Tags = append(accessCondition.Tags, aembit.TagDTO{
				Key:   key,
				Value: value,
			})
		}
	}

	accessCondition.IntegrationID = model.IntegrationID.ValueString()
	if model.Wiz != nil {
		accessCondition.Conditions.MaxLastSeenSeconds = model.Wiz.MaxLastSeen.ValueInt64()
		accessCondition.Conditions.ContainerClusterConnected = model.Wiz.ContainerClusterConnected.ValueBool()
	}
	if model.CrowdStrike != nil {
		accessCondition.Conditions.MaxLastSeenSeconds = model.CrowdStrike.MaxLastSeen.ValueInt64()
		accessCondition.Conditions.MatchHostname = model.CrowdStrike.MatchHostname.ValueBool()
		accessCondition.Conditions.MatchSerialNumber = model.CrowdStrike.MatchSerialNumber.ValueBool()
		accessCondition.Conditions.PreventRestrictedFunctionalityMode = model.CrowdStrike.PreventRestrictedFunctionalityMode.ValueBool()
		accessCondition.Conditions.MatchMacAddress = model.CrowdStrike.MatchMacAddress.ValueBool()
		accessCondition.Conditions.MatchLocalIP = model.CrowdStrike.MatchLocalIP.ValueBool()
	}
	if model.GeoIp != nil {
		// retrieve countries datasource for validation
		countriesResource := GetCountries(client)

		for _, location := range model.GeoIp.Locations {
			countryCodeInput := location.CountryCode.ValueString()

			countryIndex := slices.IndexFunc(
				countriesResource.Countries,
				func(c *countryResourceModel) bool {
					return c.CountryCode.ValueString() == countryCodeInput
				},
			)

			if countryIndex == -1 {
				return accessCondition, fmt.Errorf(
					"%v is not a valid CountryCode",
					countryCodeInput,
				)
			}

			countryFound := countriesResource.Countries[countryIndex]

			loc := aembit.CountryDTO{
				Alpha2Code: countryFound.CountryCode.ValueString(),
				ShortName:  countryFound.ShortName.ValueString(),
			}

			err := FillSubdivisions(&loc, location.Subdivisions, countryFound)
			if err != nil {
				return accessCondition, err
			}
			accessCondition.Conditions.Locations = append(accessCondition.Conditions.Locations, loc)
		}
	}
	if model.Time != nil {
		// retrieve timezones datasource for validation
		timezoneResource := GetTimezones(client)

		timezoneInput := model.Time.Timezone.ValueString()

		tsIndex := slices.IndexFunc(
			timezoneResource.Timezones,
			func(ts *timezoneResourceModel) bool {
				return ts.Timezone.ValueString() == timezoneInput
			},
		)

		if tsIndex == -1 {
			return accessCondition, fmt.Errorf("%v is not a valid timezone", timezoneInput)
		}

		timeZoneFound := timezoneResource.Timezones[tsIndex]

		accessCondition.Conditions.Timezone = &aembit.TimezoneDTO{
			Timezone: timeZoneFound.Timezone.ValueString(),
			Group:    timeZoneFound.Group.ValueString(),
			Label:    timeZoneFound.Label.ValueString(),
		}

		for _, schedule := range model.Time.Schedule {
			ordinal, err := findOrdinal(schedule.Day.ValueString())
			if err != nil {
				return accessCondition, err
			}

			accessCondition.Conditions.Schedule = append(
				accessCondition.Conditions.Schedule,
				aembit.ScheduleDTO{
					StartTime: schedule.StartTime.ValueString(),
					EndTime:   schedule.EndTime.ValueString(),
					WeekDay: &aembit.WeekDayDTO{
						Name:    schedule.Day.ValueString(),
						Ordinal: ordinal,
					},
				},
			)
		}
	}

	return accessCondition, nil
}

func FillSubdivisions(
	loc *aembit.CountryDTO,
	subDivisions []*models.GeoIpSubdivisionModel,
	countryFound *countryResourceModel,
) error {
	for _, subDivision := range subDivisions {
		subDivisionInput := subDivision.SubdivisionCode.ValueString()

		subDivisionIndex := slices.IndexFunc(
			countryFound.Subdivisions,
			func(s *countrySubdivisionResourceModel) bool {
				return s.SubdivisionCode.ValueString() == subDivisionInput
			},
		)

		if subDivisionIndex == -1 {
			return fmt.Errorf("%v is not a valid SubdivisionCode", subDivisionInput)
		}

		subdivisionFound := countryFound.Subdivisions[subDivisionIndex]

		loc.Subdivisions = append(loc.Subdivisions, aembit.SubdivisionDTO{
			SubdivisionCode: subdivisionFound.SubdivisionCode.ValueString(),
			Alpha2Code:      subdivisionFound.CountryCode.ValueString(),
			Name:            subdivisionFound.Name.ValueString(),
		})
	}

	return nil
}

func convertAccessConditionDTOToModel(
	ctx context.Context,
	dto aembit.AccessConditionDTO,
	_ models.AccessConditionResourceModel,
) models.AccessConditionResourceModel {
	var model models.AccessConditionResourceModel
	model.ID = types.StringValue(dto.ExternalID)
	model.Name = types.StringValue(dto.Name)
	model.Description = types.StringValue(dto.Description)
	model.IsActive = types.BoolValue(dto.IsActive)
	model.Tags = newTagsModel(ctx, dto.Tags)

	if len(dto.IntegrationID) == 0 {
		model.IntegrationID = types.StringValue(dto.Integration.ExternalID)
	} else {
		model.IntegrationID = types.StringValue(dto.IntegrationID)
	}
	switch dto.Integration.Type {
	case "WizIntegrationApi":
		model.Wiz = &models.AccessConditionWizModel{
			MaxLastSeen:               types.Int64Value(dto.Conditions.MaxLastSeenSeconds),
			ContainerClusterConnected: types.BoolValue(dto.Conditions.ContainerClusterConnected),
		}
	case "CrowdStrike":
		model.CrowdStrike = &models.AccessConditionCrowdstrikeModel{
			MaxLastSeen:       types.Int64Value(dto.Conditions.MaxLastSeenSeconds),
			MatchHostname:     types.BoolValue(dto.Conditions.MatchHostname),
			MatchSerialNumber: types.BoolValue(dto.Conditions.MatchSerialNumber),
			PreventRestrictedFunctionalityMode: types.BoolValue(
				dto.Conditions.PreventRestrictedFunctionalityMode,
			),
			MatchMacAddress: types.BoolValue(dto.Conditions.MatchMacAddress),
			MatchLocalIP:    types.BoolValue(dto.Conditions.MatchLocalIP),
		}
	case "AembitGeoIPCondition":
		geoIpModel := models.AccessConditionGeoIpModel{}

		for _, location := range dto.Conditions.Locations {
			loc := models.GeoIpLocationModel{
				CountryCode:  types.StringValue(location.Alpha2Code),
				Subdivisions: []*models.GeoIpSubdivisionModel{},
			}

			for _, subDivision := range location.Subdivisions {
				loc.Subdivisions = append(loc.Subdivisions, &models.GeoIpSubdivisionModel{
					SubdivisionCode: types.StringValue(subDivision.SubdivisionCode),
				})
			}

			geoIpModel.Locations = append(geoIpModel.Locations, &loc)
		}

		model.GeoIp = &geoIpModel
	case "AembitTimeCondition":
		acTimeZone := models.AccessConditionTimeZoneModel{}

		for _, schedule := range dto.Conditions.Schedule {
			acTimeZone.Schedule = append(acTimeZone.Schedule, &models.ScheduleModel{
				StartTime: types.StringValue(schedule.StartTime),
				EndTime:   types.StringValue(schedule.EndTime),
				Day:       types.StringValue(schedule.WeekDay.Name),
			})
		}

		acTimeZone.Timezone = types.StringValue(dto.Conditions.Timezone.Timezone)

		model.Time = &acTimeZone
	}

	return model
}

func findOrdinal(weekDay string) (int, error) {
	weekdays := map[string]int{
		"sunday":    0,
		"monday":    1,
		"tuesday":   2,
		"wednesday": 3,
		"thursday":  4,
		"friday":    5,
		"saturday":  6,
	}

	if val, ok := weekdays[strings.ToLower(weekDay)]; ok {
		return val, nil
	}

	return -1, fmt.Errorf("%v is not a valid weekday name", weekDay)
}

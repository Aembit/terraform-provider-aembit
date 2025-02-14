package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// accessConditionResourceModel maps the resource schema.
type accessConditionResourceModel struct {
	// ID is required for Framework acceptance testing
	ID            types.String                     `tfsdk:"id"`
	Name          types.String                     `tfsdk:"name"`
	Description   types.String                     `tfsdk:"description"`
	IsActive      types.Bool                       `tfsdk:"is_active"`
	Tags          types.Map                        `tfsdk:"tags"`
	IntegrationID types.String                     `tfsdk:"integration_id"`
	Wiz           *accessConditionWizModel         `tfsdk:"wiz_conditions"`
	CrowdStrike   *accessConditionCrowdstrikeModel `tfsdk:"crowdstrike_conditions"`
	GeoIp         *accessConditionGeoIpModel       `tfsdk:"geoip_conditions"`
	Timezone      *accessConditionTimeZoneModel    `tfsdk:"timezone_conditions"`
}

type accessConditionWizModel struct {
	MaxLastSeen               types.Int64 `tfsdk:"max_last_seen"`
	ContainerClusterConnected types.Bool  `tfsdk:"container_cluster_connected"`
}

type accessConditionCrowdstrikeModel struct {
	MaxLastSeen                        types.Int64 `tfsdk:"max_last_seen"`
	MatchHostname                      types.Bool  `tfsdk:"match_hostname"`
	MatchSerialNumber                  types.Bool  `tfsdk:"match_serial_number"`
	PreventRestrictedFunctionalityMode types.Bool  `tfsdk:"prevent_rfm"`
}

type accessConditionGeoIpModel struct {
	Locations []*geoIpLocationModel `tfsdk:"locations"`
}

type geoIpSubdivisionModel struct {
	SubdivisionCode types.String `tfsdk:"subdivision_code"`
}

type geoIpLocationModel struct {
	Alpha2Code   types.String             `tfsdk:"alpha2_code"`
	Subdivisions []*geoIpSubdivisionModel `tfsdk:"subdivisions"`
}

type accessConditionTimeZoneModel struct {
	Schedule []*scheduleModel `tfsdk:"schedule"`
	Timezone types.String     `tfsdk:"timezone"`
}

type scheduleModel struct {
	EndTime   types.String `tfsdk:"end_time"`
	StartTime types.String `tfsdk:"start_time"`
	Day       types.String `tfsdk:"day"`
}

// accessConditionDataSourceModel maps the datasource schema.
type accessConditionsDataSourceModel struct {
	AccessConditions []accessConditionResourceModel `tfsdk:"access_conditions"`
}

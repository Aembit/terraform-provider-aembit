package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// models.AccessConditionResourceModel maps the resource schema.
type AccessConditionResourceModel struct {
	// ID is required for Framework acceptance testing
	ID            types.String                     `tfsdk:"id"`
	Name          types.String                     `tfsdk:"name"`
	Description   types.String                     `tfsdk:"description"`
	IsActive      types.Bool                       `tfsdk:"is_active"`
	Tags          types.Map                        `tfsdk:"tags"`
	IntegrationID types.String                     `tfsdk:"integration_id"`
	Wiz           *AccessConditionWizModel         `tfsdk:"wiz_conditions"`
	CrowdStrike   *AccessConditionCrowdstrikeModel `tfsdk:"crowdstrike_conditions"`
	GeoIp         *AccessConditionGeoIpModel       `tfsdk:"geoip_conditions"`
	Time          *AccessConditionTimeZoneModel    `tfsdk:"time_conditions"`
}

type AccessConditionWizModel struct {
	MaxLastSeen               types.Int64 `tfsdk:"max_last_seen"`
	ContainerClusterConnected types.Bool  `tfsdk:"container_cluster_connected"`
}

type AccessConditionCrowdstrikeModel struct {
	MaxLastSeen                        types.Int64 `tfsdk:"max_last_seen"`
	MatchHostname                      types.Bool  `tfsdk:"match_hostname"`
	MatchSerialNumber                  types.Bool  `tfsdk:"match_serial_number"`
	PreventRestrictedFunctionalityMode types.Bool  `tfsdk:"prevent_rfm"`
	MatchMacAddress                    types.Bool  `tfsdk:"match_mac_address"`
	MatchLocalIP                       types.Bool  `tfsdk:"match_local_ip"`
}

type AccessConditionGeoIpModel struct {
	Locations []*GeoIpLocationModel `tfsdk:"locations"`
}

type GeoIpSubdivisionModel struct {
	SubdivisionCode types.String `tfsdk:"subdivision_code"`
}

type GeoIpLocationModel struct {
	CountryCode  types.String             `tfsdk:"country_code"`
	Subdivisions []*GeoIpSubdivisionModel `tfsdk:"subdivisions"`
}

type AccessConditionTimeZoneModel struct {
	Schedule []*ScheduleModel `tfsdk:"schedule"`
	Timezone types.String     `tfsdk:"timezone"`
}

type ScheduleModel struct {
	EndTime   types.String `tfsdk:"end_time"`
	StartTime types.String `tfsdk:"start_time"`
	Day       types.String `tfsdk:"day"`
}

// accessConditionDataSourceModel maps the datasource schema.
type AccessConditionsDataSourceModel struct {
	AccessConditions []AccessConditionResourceModel `tfsdk:"access_conditions"`
}

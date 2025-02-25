provider "aembit" {
}

data "aembit_integrations" "filtered" {
	type = "AembitGeoIPCondition"
}

resource "aembit_access_condition" "geoip" {
	name = "TF Acceptance GeoIp"
	is_active = true
	integration_id = data.aembit_integrations.filtered.integrations[0].id
	geoip_conditions = {
		locations = [
			{
				country_code = "TR",
				subdivisions = []
			}
		]
	}
    tags = {
        color = "blue"
        day   = "Sunday"
    }
}
provider "aembit" {
}

resource "aembit_access_condition" "geoip" {
	name = "TF Acceptance GeoIp"
	is_active = true
	integration_id = "426efdca-4626-4730-8264-44d6fc6f8553"
	geoip_conditions = {
		locations = [
			{
				alpha2_code = "Türkiye",
				short_name = "TR",
				subdivisions = []
			}
		]
	}
    tags = {
        color = "blue"
        day   = "Sunday"
    }
}
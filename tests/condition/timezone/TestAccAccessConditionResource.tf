provider "aembit" {
}

data "aembit_integrations" "filtered" {
	type = "AembitTimeCondition"
}

resource "aembit_access_condition" "timezone" {
	name = "TF Acceptance Timezone"
	is_active = true
	integration_id = data.aembit_integrations.filtered.integrations[0].id
	time_conditions = {
		schedule = [
			{
                end_time: "17:00",
                day: "Wednesday"
                start_time: "08:00"
            },
            {
                end_time: "17:00",
                day: "Friday"
                start_time: "08:00"
            },
		]
		timezone = "America/Metlakatla"
	}
    tags = {
        color = "blue"
        day   = "Sunday"
    }
}
provider "aembit" {
}

resource "aembit_access_condition" "timezone" {
	name = "TF Acceptance TimeZone"
	is_active = true
	integration_id = "aa7c2571-4a39-4eff-a4b2-92b0df0a9540"
		timezone_conditions = {
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
		timezone = {
            timezone: "America/Metlakatla"
        }
	}
    tags = {
        color = "blue"
        day   = "Sunday"
    }
}
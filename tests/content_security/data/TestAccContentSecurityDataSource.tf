provider "aembit" {
}

resource "aembit_content_security" "crowdstrike" {
	name = "TF Acceptance CrowdStrike Content Security"
	is_active = true
	type = "CrowdStrikeAIDR"
	crowdstrike_falcon_aidr = {
		base_url = "https://api.crowdstrike.com/auth"
		encrypted_token = "test_token"
		fail_open = true
		max_retries = 10
		timeout_ms = 3000
	}
}

data "aembit_content_securities" "crowdstrike_datasource" {
    depends_on = [ aembit_content_security.crowdstrike ]
}

provider "aembit" {
}

resource "aembit_routing" "default" {
	name = "TF Acceptance Routing"
	is_active = false
	proxy_url = "http://test.com:9876"
	description = "TF Acceptance Routing desc"
	resource_set_id = "ffffffff-ffff-ffff-ffff-ffffffffffff"
}

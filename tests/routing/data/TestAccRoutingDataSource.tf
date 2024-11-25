provider "aembit" {
}

resource "aembit_routing" "routing" {
	name = "TF Acceptance Routing"
	is_active = true
	proxy_url = "http://test.com:9876"
	description = "TF Acceptance Routing desc"
	resource_set_id = "ffffffff-ffff-ffff-ffff-ffffffffffff"
}

data "aembit_routings" "test" {
    depends_on = [ aembit_routing.routing ]
}

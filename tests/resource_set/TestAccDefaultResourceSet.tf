provider "aembit" {
}

data "aembit_resource_set" "default" {
	id = "ffffffff-ffff-ffff-ffff-ffffffffffff"
}

data "aembit_resource_sets" "all" {
}
package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_GPC_DataSource_Read(t *testing.T) {
	const gpc_datasource_name string = "data.aembit_global_policy_compliance_data.test"
	const gpc_datasource_def string = `provider "aembit" {}
		data "aembit_global_policy_compliance_data" "test" {}`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create
			{
				ResourceName: gpc_datasource_name,
				Config:       gpc_datasource_def,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(gpc_datasource_name, "access_policy_trust_provider_compliance", "Recommended"),
					resource.TestCheckResourceAttr(gpc_datasource_name, "access_policy_access_condition_compliance", "Recommended"),
					resource.TestCheckResourceAttr(gpc_datasource_name, "agent_controller_trust_provider_compliance", "Recommended"),
					resource.TestCheckResourceAttr(gpc_datasource_name, "agent_controller_allowed_tls_hostname_compliance", "Recommended"),
				),
			},
		},
	})
}

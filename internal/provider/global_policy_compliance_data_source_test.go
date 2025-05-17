package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_GPC_DataSource_Read(t *testing.T) {
	const gpcDatasourceName string = "data.aembit_global_policy_compliance_data.test"
	const gpcDatasourceDef string = `provider "aembit" {}
		
		resource "aembit_global_policy_compliance" "test" {
			access_policy_trust_provider_compliance = "Optional"
			access_policy_access_condition_compliance = "Optional"
			agent_controller_trust_provider_compliance = "Optional"
			agent_controller_allowed_tls_hostname_compliance = "Optional"
		}
		
		data "aembit_global_policy_compliance_data" "test" {
			depends_on = [ aembit_global_policy_compliance.test ]
		}`

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create
			{
				ResourceName: gpcDatasourceName,
				Config:       gpcDatasourceDef,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(gpcDatasourceName, "access_policy_trust_provider_compliance", "Optional"),
					resource.TestCheckResourceAttr(gpcDatasourceName, "access_policy_access_condition_compliance", "Optional"),
					resource.TestCheckResourceAttr(gpcDatasourceName, "agent_controller_trust_provider_compliance", "Optional"),
					resource.TestCheckResourceAttr(gpcDatasourceName, "agent_controller_allowed_tls_hostname_compliance", "Optional"),
				),
			},
		},
	})
}

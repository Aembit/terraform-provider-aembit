package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const (
	testAccessConditionsDataSource string = "data.aembit_access_conditions.test"
	testAccessConditionResource    string = "aembit_access_condition.crowdstrike"
)

func testFindAccessCondition(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var rs *terraform.ResourceState
		var err error
		var ok bool
		var notFound bool
		if rs, ok = s.RootModule().Resources[resourceName]; !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if _, err, notFound = testClient.GetAccessCondition(rs.Primary.ID, nil); notFound {
			return err
		}
		return nil
	}
}

func TestAccAccessConditionsDataSource(t *testing.T) {
	t.Parallel()
	createFile, _ := os.ReadFile("../../tests/condition/data/TestAccAccessConditionsDataSource.tf")
	createFileConfig, _, _ := randomizeFileConfigs(
		string(createFile),
		"",
		"TF Acceptance Crowdstrike",
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: createFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify non-zero number of Access Conditions returned
					resource.TestCheckResourceAttrSet(
						testAccessConditionsDataSource,
						"access_conditions.#",
					),
					resource.TestCheckResourceAttr(
						testAccessConditionResource,
						"crowdstrike_conditions.match_mac_address",
						"true",
					),
					resource.TestCheckResourceAttr(
						testAccessConditionResource,
						"crowdstrike_conditions.match_local_ip",
						"true",
					),
					// Find newly created entry
					testFindAccessCondition(testAccessConditionResource),
				),
			},
		},
	})
}

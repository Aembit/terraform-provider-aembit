package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testAccessConditionResourceWiz string = "aembit_access_condition.wiz"
const testAccessConditionResourceCrowdstrike string = "aembit_access_condition.crowdstrike"

func testDeleteAccessCondition(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var rs *terraform.ResourceState
		var ok bool
		var err error
		if rs, ok = s.RootModule().Resources[resourceName]; !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if ok, err = testClient.DeleteAccessCondition(rs.Primary.ID, nil); !ok {
			return err
		}
		return nil
	}
}

func TestAccAccessConditionResource_Wiz(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/condition/wiz/TestAccAccessConditionResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/condition/wiz/TestAccAccessConditionResource.tfmod")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify AccessCondition Name
					resource.TestCheckResourceAttr(testAccessConditionResourceWiz, "name", "TF Acceptance Wiz"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testAccessConditionResourceWiz, "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(testAccessConditionResourceWiz, "id"),
				),
			},
			// Test Aembit API Removal causes re-create with non-empty plan
			{Config: string(createFile), Check: testDeleteAccessCondition(testAccessConditionResourceWiz), ExpectNonEmptyPlan: true},
			// Recreate the resource from the first test step
			{Config: string(createFile)},
			// ImportState testing
			{ResourceName: testAccessConditionResourceWiz, ImportState: true, ImportStateVerify: true},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(testAccessConditionResourceWiz, "name", "TF Acceptance Wiz - Modified"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccAccessConditionResource_Crowdstrike(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/condition/crowdstrike/TestAccAccessConditionResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/condition/crowdstrike/TestAccAccessConditionResource.tfmod")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify AccessCondition Name
					resource.TestCheckResourceAttr(testAccessConditionResourceCrowdstrike, "name", "TF Acceptance Crowdstrike"),
					// Verify Tags.
					resource.TestCheckResourceAttr(testAccessConditionResourceCrowdstrike, tagsCount, "2"),
					resource.TestCheckResourceAttr(testAccessConditionResourceCrowdstrike, tagsColor, "blue"),
					resource.TestCheckResourceAttr(testAccessConditionResourceCrowdstrike, tagsDay, "Sunday"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testAccessConditionResourceCrowdstrike, "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(testAccessConditionResourceCrowdstrike, "id"),
				),
			},
			// ImportState testing
			{ResourceName: testAccessConditionResourceCrowdstrike, ImportState: true, ImportStateVerify: true},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(testAccessConditionResourceCrowdstrike, "name", "TF Acceptance Crowdstrike - Modified"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

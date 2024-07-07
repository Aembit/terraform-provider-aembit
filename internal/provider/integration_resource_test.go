package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testIntegrationWiz string = "aembit_integration.wiz"
const testIntegrationCrowdstrike string = "aembit_integration.crowdstrike"

func testDeleteIntegration(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var rs *terraform.ResourceState
		var ok bool
		var err error
		if rs, ok = s.RootModule().Resources[resourceName]; !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if ok, err = testClient.DeleteIntegration(rs.Primary.ID, nil); !ok {
			return err
		}
		return nil
	}
}

func TestAccIntegrationResource_Wiz(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/integration/wiz/TestAccIntegrationResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/integration/wiz/TestAccIntegrationResource.tfmod")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Integration Name
					resource.TestCheckResourceAttr(testIntegrationWiz, "name", "TF Acceptance Wiz"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testIntegrationWiz, "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(testIntegrationWiz, "id"),
					// Verify Tags.
					resource.TestCheckResourceAttr(testIntegrationWiz, tagsCount, "2"),
					resource.TestCheckResourceAttr(testIntegrationWiz, tagsColor, "blue"),
					resource.TestCheckResourceAttr(testIntegrationWiz, tagsDay, "Sunday"),
				),
			},
			// Test Aembit API Removal causes re-create with non-empty plan
			{Config: string(createFile), Check: testDeleteIntegration(testIntegrationWiz), ExpectNonEmptyPlan: true},
			// Recreate the resource from the first test step
			{Config: string(createFile)},
			// ImportState testing
			{ResourceName: testIntegrationWiz, ImportState: true, ImportStateVerify: false},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(testIntegrationWiz, "name", "TF Acceptance Wiz - Modified"),
					// Verify Tags.
					resource.TestCheckResourceAttr(testIntegrationWiz, tagsCount, "2"),
					resource.TestCheckResourceAttr(testIntegrationWiz, tagsColor, "orange"),
					resource.TestCheckResourceAttr(testIntegrationWiz, tagsDay, "Tuesday"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccIntegrationResource_Crowdstrike(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/integration/crowdstrike/TestAccIntegrationResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/integration/crowdstrike/TestAccIntegrationResource.tfmod")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Integration Name
					resource.TestCheckResourceAttr(testIntegrationCrowdstrike, "name", "TF Acceptance Crowdstrike"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testIntegrationCrowdstrike, "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(testIntegrationCrowdstrike, "id"),
				),
			},
			// ImportState testing
			{ResourceName: testIntegrationCrowdstrike, ImportState: true, ImportStateVerify: false},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(testIntegrationCrowdstrike, "name", "TF Acceptance Crowdstrike - Modified"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

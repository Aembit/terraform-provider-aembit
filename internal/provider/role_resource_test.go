package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const resourceID string = "aembit_role.role"

func testDeleteRole(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var rs *terraform.ResourceState
		var ok bool
		var err error
		if rs, ok = s.RootModule().Resources[resourceName]; !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if ok, err = testClient.DeleteRole(rs.Primary.ID, nil); !ok {
			return err
		}
		return nil
	}
}

func TestAccRoleResource(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/roles/TestAccRoleResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/roles/TestAccRoleResource.tfmod")
	createFileConfig, modifyFileConfig, newName := randomizeFileConfigs(string(createFile), string(modifyFile), "TF Acceptance Role")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: createFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Trust Provider Name
					resource.TestCheckResourceAttr(resourceID, "name", newName),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(resourceID, "id"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(resourceID, "id"),
				),
			},
			// Test Aembit API Removal causes re-create with non-empty plan
			{Config: createFileConfig, Check: testDeleteRole(resourceID), ExpectNonEmptyPlan: true},
			// Recreate the resource from the first test step
			{Config: createFileConfig},
			// ImportState testing
			{ResourceName: resourceID, ImportState: true, ImportStateVerify: true},
			// Update and Read testing
			{
				Config: modifyFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(resourceID, "name", newName),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

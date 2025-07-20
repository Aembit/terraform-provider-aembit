package provider

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const (
	testRolesDataSource string = "data.aembit_roles.test"
	testRoleResource    string = "aembit_role.role"
)

func testFindRole(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var rs *terraform.ResourceState
		var err error
		var ok bool
		var notFound bool
		if rs, ok = s.RootModule().Resources[resourceName]; !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if _, err, notFound = testClient.GetRole(rs.Primary.ID, nil); notFound {
			return err
		}
		return nil
	}
}

func TestAccRolesDataSource(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/roles/data/TestAccRolesDataSource.tf")
	randID := rand.Intn(10000000)
	createFileConfig := strings.ReplaceAll(
		string(createFile),
		"TF Acceptance Role",
		fmt.Sprintf("TF Acceptance Role%d", randID),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: createFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of Roles returned
					resource.TestCheckResourceAttrSet(testRolesDataSource, "roles.#"),
					// Verify the attributes of the first three roles
					resource.TestCheckResourceAttr(
						testRolesDataSource,
						"roles.0.name",
						"SuperAdmin",
					),
					resource.TestCheckResourceAttr(
						testRolesDataSource,
						"roles.0.is_active",
						"true",
					),
					resource.TestCheckResourceAttr(testRolesDataSource, "roles.1.name", "Auditor"),
					resource.TestCheckResourceAttr(
						testRolesDataSource,
						"roles.1.is_active",
						"true",
					),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testRolesDataSource, "roles.0.id"),
					resource.TestCheckResourceAttrSet(testRolesDataSource, "roles.1.id"),
					resource.TestCheckResourceAttrSet(testRolesDataSource, "roles.2.id"),
					// Find newly created entry
					testFindRole(testRoleResource),
				),
			},
		},
	})
}

package provider

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testRolesDataSource string = "data.aembit_roles.test"

func TestAccRolesDataSource(t *testing.T) {

	createFile, _ := os.ReadFile("../../tests/roles/data/TestAccRolesDataSource.tf")
	randID := rand.Intn(10000000)
	createFileConfig := strings.ReplaceAll(string(createFile), "TF Acceptance Role", fmt.Sprintf("TF Acceptance Role%d", randID))

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
					resource.TestCheckResourceAttr(testRolesDataSource, "roles.0.name", "SuperAdmin"),
					resource.TestCheckResourceAttr(testRolesDataSource, "roles.0.is_active", "true"),
					resource.TestCheckResourceAttr(testRolesDataSource, "roles.1.name", "Auditor"),
					resource.TestCheckResourceAttr(testRolesDataSource, "roles.1.is_active", "true"),
					resource.TestCheckResourceAttr(testRolesDataSource, "roles.2.name", fmt.Sprintf("TF Acceptance Role%d", randID)),
					resource.TestCheckResourceAttr(testRolesDataSource, "roles.2.is_active", "false"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testRolesDataSource, "roles.0.id"),
					resource.TestCheckResourceAttrSet(testRolesDataSource, "roles.1.id"),
					resource.TestCheckResourceAttrSet(testRolesDataSource, "roles.2.id"),
				),
			},
		},
	})
}

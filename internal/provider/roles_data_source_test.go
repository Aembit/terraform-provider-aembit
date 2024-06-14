package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRolesDataSource(t *testing.T) {

	createFile, _ := os.ReadFile("../../tests/roles/data/TestAccRolesDataSource.tf")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of Roles returned
					resource.TestCheckResourceAttrSet("data.aembit_roles.test", "roles.#"),
					// Verify the attributes of the first three roles
					resource.TestCheckResourceAttr("data.aembit_roles.test", "roles.0.name", "SuperAdmin"),
					resource.TestCheckResourceAttr("data.aembit_roles.test", "roles.0.is_active", "true"),
					resource.TestCheckResourceAttr("data.aembit_roles.test", "roles.1.name", "Auditor"),
					resource.TestCheckResourceAttr("data.aembit_roles.test", "roles.1.is_active", "true"),
					resource.TestCheckResourceAttr("data.aembit_roles.test", "roles.2.name", "TF Acceptance Role"),
					resource.TestCheckResourceAttr("data.aembit_roles.test", "roles.2.is_active", "false"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("data.aembit_roles.test", "roles.0.id"),
					resource.TestCheckResourceAttrSet("data.aembit_roles.test", "roles.1.id"),
					resource.TestCheckResourceAttrSet("data.aembit_roles.test", "roles.2.id"),
				),
			},
		},
	})
}

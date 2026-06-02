package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccResourceSet(t *testing.T) {
	t.Parallel()
	createFile, err := os.ReadFile(
		"../../tests/resource_set/TestAccResourceSet.tf",
	)
	if err != nil {
		t.Fatalf("failed to read test config file: %v", err)
	}
	modifiedFile, err := os.ReadFile(
		"../../tests/resource_set/TestAccResourceSet.tfmod",
	)
	if err != nil {
		t.Fatalf("failed to read test config file: %v", err)
	}
	resourceName := "aembit_resource_set.crs"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: string(createFile),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName,
						"name",
						"TF Acceptance Custom ResourceSet",
					),
					resource.TestCheckResourceAttr(
						resourceName,
						"description",
						"TF Acceptance Custom ResourceSet",
					),
					resource.TestCheckResourceAttr(
						resourceName,
						"roles.#",
						"2",
					),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: string(modifiedFile),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName,
						"name",
						"TF Acceptance Custom ResourceSet - Modified",
					),
					resource.TestCheckResourceAttr(
						resourceName,
						"description",
						"TF Acceptance Custom ResourceSet - Modified",
					),
					resource.TestCheckResourceAttr(
						resourceName,
						"roles.#",
						"2",
					),
				),
			},
		},
	})
}

package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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

func TestAccResourceSetPolicy(t *testing.T) {
	t.Parallel()
	createFile, err := os.ReadFile(
		"../../tests/resource_set/TestAccResourceSetPolicy.tf",
	)
	if err != nil {
		t.Fatalf("failed to read test config file: %v", err)
	}
	modifiedFile, err := os.ReadFile(
		"../../tests/resource_set/TestAccResourceSetPolicy.tfmod",
	)
	if err != nil {
		t.Fatalf("failed to read test config file: %v", err)
	}
	resourceName := "aembit_access_policy.first_policy"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: string(createFile),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_set_id"),
					resource.TestCheckResourceAttr(resourceName, "name", "TF ResourceSet Policy"),
					resource.TestCheckResourceAttr(resourceName, "is_active", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "client_workload"),
					resource.TestCheckResourceAttr(resourceName, "trust_providers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "access_conditions.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "credential_provider"),
					resource.TestCheckResourceAttrSet(resourceName, "server_workload"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					rs, ok := state.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("Not found: %s", resourceName)
					}
					return fmt.Sprintf("%s,%s", rs.Primary.Attributes["resource_set_id"], rs.Primary.ID), nil
				},
			},
			{
				Config: string(modifiedFile),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_set_id"),
					resource.TestCheckResourceAttr(resourceName, "name", "TF ResourceSet Policy"),
					resource.TestCheckResourceAttr(resourceName, "is_active", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "client_workload"),
					resource.TestCheckResourceAttr(resourceName, "trust_providers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "access_conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "credential_providers.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(
						resourceName,
						"credential_providers.*",
						map[string]string{
							"mapping_type": "HttpHeader",
							"header_name":  "test_header_name",
							"header_value": "test_header_value",
						},
					),
					resource.TestCheckTypeSetElemNestedAttrs(
						resourceName,
						"credential_providers.*",
						map[string]string{
							"mapping_type":         "HttpBody",
							"httpbody_field_path":  "test_field_path",
							"httpbody_field_value": "test_field_value",
						},
					),
					resource.TestCheckResourceAttrSet(resourceName, "server_workload"),
				),
			},
		},
	})
}

package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testContentSecuritiesDataSource string = "data.aembit_content_securities.crowdstrike_datasource"

func testFindContentSecurity(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var rs *terraform.ResourceState
		var err error
		var ok bool
		var notFound bool
		if rs, ok = s.RootModule().Resources[resourceName]; !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		resourceSetID := rs.Primary.Attributes["resource_set_id"]

		if _, err, notFound = testClient.GetContentSecurity(rs.Primary.ID, nil, &resourceSetID); notFound {
			return err
		}
		return nil
	}
}

func TestAccContentSecurityDataSource(t *testing.T) {
	createFile1, _ := os.ReadFile("../../tests/content_security/data/TestAccContentSecurityDataSource_ProviderResourceSet.tf")
	createFile2, _ := os.ReadFile("../../tests/content_security/data/TestAccContentSecurityDataSource_ResourceSet.tf")
	createFile3, _ := os.ReadFile("../../tests/content_security/data/TestAccContentSecurityDataSource.tf")

	files := [3]string{string(createFile1), string(createFile2), string(createFile3)}

	for _, createFile := range files {
		resource.Test(t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				// Read testing
				{
					Config: createFile,
					Check: resource.ComposeAggregateTestCheckFunc(
						// Verify non-zero number of Content Securities returned
						resource.TestCheckResourceAttrSet(
							testContentSecuritiesDataSource,
							"content_securities.#",
						),
						// Verify dynamic values have any value set in the state.
						resource.TestCheckResourceAttrSet(
							testContentSecuritiesDataSource,
							"content_securities.0.id",
						),
						// Verify placeholder ID is set
						resource.TestCheckResourceAttrSet(
							testContentSecuritiesDataSource,
							"content_securities.0.id",
						),
						// Find newly created entry
						testFindContentSecurity(testContentSecurityResource),
					),
				},
			},
		})
	}
}

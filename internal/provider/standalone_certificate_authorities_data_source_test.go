package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const testStandaloneCertificatesDataSource string = "data.aembit_standalone_certificate_authorities.test"

func testFindStandaloneCertificate(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var rs *terraform.ResourceState
		var err error
		var ok bool
		var notFound bool
		if rs, ok = s.RootModule().Resources[resourceName]; !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if _, err, notFound = testClient.GetStandaloneCertificate(rs.Primary.ID, nil); notFound {
			return err
		}
		return nil
	}
}

func TestAccStandaloneCertificatesDataSource(t *testing.T) {

	createFile, _ := os.ReadFile("../../tests/standalone_certificate_authority/data/TestAccStandaloneCertificatesDataSource.tf")
	createFileConfig, _, _ := randomizeFileConfigs(string(createFile), "", "Unit Test 1")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: createFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify non-zero number of Standalone Certificate Authorities returned
					resource.TestCheckResourceAttrSet(testStandaloneCertificatesDataSource, "standalone_certificate_authorities.#"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet(testStandaloneCertificatesDataSource, "standalone_certificate_authorities.0.id"),
					resource.TestCheckResourceAttrSet(testStandaloneCertificatesDataSource, "standalone_certificate_authorities.0.not_before"),
					resource.TestCheckResourceAttrSet(testStandaloneCertificatesDataSource, "standalone_certificate_authorities.0.not_after"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(testStandaloneCertificatesDataSource, "standalone_certificate_authorities.0.id"),
					// Find newly created entry
					testFindStandaloneCertificate(testStandaloneCertificateResource),
				),
			},
		},
	})
}

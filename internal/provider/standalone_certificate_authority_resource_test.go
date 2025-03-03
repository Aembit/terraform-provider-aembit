package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testStandaloneCertificateResource string = "aembit_standalone_certificate_authority.test"

func TestAccStandaloneCertificateAuthorityResource(t *testing.T) {
	createFile, _ := os.ReadFile("../../tests/standalone_certificate_authority/TestAccStandaloneCertificateAuthorityResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/standalone_certificate_authority/TestAccStandaloneCertificateAuthorityResource.tfmod")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name
					resource.TestCheckResourceAttr(testStandaloneCertificateResource, "name", "Unit Test 1"),
					resource.TestCheckResourceAttr(testStandaloneCertificateResource, "is_active", "true"),
					// Verify Tags
					resource.TestCheckResourceAttr(testStandaloneCertificateResource, tagsCount, "2"),
					resource.TestCheckResourceAttr(testStandaloneCertificateResource, tagsColor, "blue"),
					resource.TestCheckResourceAttr(testStandaloneCertificateResource, tagsDay, "Sunday"),
					// Verify Metadata
					resource.TestCheckResourceAttrSet(testStandaloneCertificateResource, "not_before"),
					resource.TestCheckResourceAttrSet(testStandaloneCertificateResource, "not_after"),
					resource.TestCheckResourceAttr(testStandaloneCertificateResource, "leaf_lifetime", "1440"),
					resource.TestCheckResourceAttr(testStandaloneCertificateResource, "client_workload_count", "0"),
					resource.TestCheckResourceAttr(testStandaloneCertificateResource, "resource_set_count", "0"),
					// Verify placeholder ID is set
					resource.TestCheckResourceAttrSet(testStandaloneCertificateResource, "id"),
				),
			},
			// ImportState testing
			{ResourceName: testStandaloneCertificateResource, ImportState: true, ImportStateVerify: true},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name updated
					resource.TestCheckResourceAttr(testStandaloneCertificateResource, "name", "Unit Test 1 - Modified"),
					resource.TestCheckResourceAttr(testStandaloneCertificateResource, "is_active", "true"),
					// Verify Tags.
					resource.TestCheckResourceAttr(testStandaloneCertificateResource, tagsCount, "2"),
					resource.TestCheckResourceAttr(testStandaloneCertificateResource, tagsColor, "orange"),
					resource.TestCheckResourceAttr(testStandaloneCertificateResource, tagsDay, "Tuesday"),
					// Verify Metadata.
					resource.TestCheckResourceAttrSet(testStandaloneCertificateResource, "not_before"),
					resource.TestCheckResourceAttrSet(testStandaloneCertificateResource, "not_after"),
					resource.TestCheckResourceAttr(testStandaloneCertificateResource, "leaf_lifetime", "60"),
					resource.TestCheckResourceAttr(testStandaloneCertificateResource, "client_workload_count", "0"),
					resource.TestCheckResourceAttr(testStandaloneCertificateResource, "resource_set_count", "0"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

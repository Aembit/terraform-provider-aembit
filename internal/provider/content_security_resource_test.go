package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testContentSecurityResource string = "aembit_content_security.crowdstrike"

func TestAccContentSecurityResource(t *testing.T) {
	t.Parallel()
	createFile, _ := os.ReadFile(
		"../../tests/content_security/crowdstrike_aidr/TestAccContentSecurityResource.tf",
	)
	modifyFile, _ := os.ReadFile(
		"../../tests/content_security/crowdstrike_aidr/TestAccContentSecurityResource.tfmod",
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: string(createFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name
					resource.TestCheckResourceAttr(
						testContentSecurityResource,
						"name",
						"TF Acceptance CrowdStrike Content Security",
					),
					// Verify other properties
					resource.TestCheckResourceAttr(
						testContentSecurityResource,
						"is_active",
						"true",
					),
					resource.TestCheckResourceAttr(
						testContentSecurityResource,
						"type",
						"CrowdStrikeAIDR",
					),
					resource.TestCheckResourceAttr(
						testContentSecurityResource,
						"crowdstrike_falcon_aidr.base_url",
						"https://api.crowdstrike.com/auth",
					),
					resource.TestCheckResourceAttr(
						testContentSecurityResource,
						"crowdstrike_falcon_aidr.encrypted_token",
						"test_token",
					),
					resource.TestCheckResourceAttr(
						testContentSecurityResource,
						"crowdstrike_falcon_aidr.fail_open",
						"true",
					),
					resource.TestCheckResourceAttr(
						testContentSecurityResource,
						"crowdstrike_falcon_aidr.max_retries",
						"10",
					),
					resource.TestCheckResourceAttr(
						testContentSecurityResource,
						"crowdstrike_falcon_aidr.timeout_ms",
						"6000",
					),
					// Verify placeholder ID and resource_set_id are set
					resource.TestCheckResourceAttrSet(testContentSecurityResource, "id"),
					resource.TestCheckResourceAttrSet(testContentSecurityResource, "resource_set_id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      testContentSecurityResource,
				ImportState:       true,
				ImportStateVerify: false,
			},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name unchanged
					resource.TestCheckResourceAttr(
						testContentSecurityResource,
						"name",
						"TF Acceptance CrowdStrike Content Security - Modified",
					),
					resource.TestCheckResourceAttr(
						testContentSecurityResource,
						"is_active",
						"true",
					),
					resource.TestCheckResourceAttr(
						testContentSecurityResource,
						"type",
						"CrowdStrikeAIDR",
					),
					resource.TestCheckResourceAttr(
						testContentSecurityResource,
						"crowdstrike_falcon_aidr.base_url",
						"https://api.crowdstrike.com/auth",
					),
					resource.TestCheckResourceAttr(
						testContentSecurityResource,
						"crowdstrike_falcon_aidr.fail_open",
						"false",
					),
					// timeout_ms and max_retries should fall back to schema defaults
					resource.TestCheckResourceAttr(
						testContentSecurityResource,
						"crowdstrike_falcon_aidr.max_retries",
						"2",
					),
					resource.TestCheckResourceAttr(
						testContentSecurityResource,
						"crowdstrike_falcon_aidr.timeout_ms",
						"5000",
					),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

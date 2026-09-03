package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testMcpContentSecurityResource string = "aembit_content_security.mcp_access_control"
const testMcpAccessRuleResource string = "aembit_mcp_tool_access_rule.rule1"

func TestAccMcpContentSecurityResource(t *testing.T) {
	t.Parallel()
	createFile, _ := os.ReadFile(
		"../../tests/content_security/mcp_tool_access_control/TestAccContentSecurityResource.tf",
	)
	modifyFile, _ := os.ReadFile(
		"../../tests/content_security/mcp_tool_access_control/TestAccContentSecurityResource.tfmod",
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
						testMcpContentSecurityResource,
						"name",
						"TF Acceptance MCP Tool Access Control Content Security",
					),
					// Verify other properties
					resource.TestCheckResourceAttr(
						testMcpContentSecurityResource,
						"is_active",
						"true",
					),
					resource.TestCheckResourceAttr(
						testMcpContentSecurityResource,
						"type",
						"McpToolAccessControl",
					),
					resource.TestCheckResourceAttr(
						testMcpContentSecurityResource,
						"mcp_tool_access_control.mode",
						"Allow",
					),
					resource.TestCheckResourceAttr(
						testMcpContentSecurityResource,
						"mcp_tool_access_control.visibility",
						"AllowSpecific",
					),
					resource.TestCheckResourceAttr(
						testMcpContentSecurityResource,
						"mcp_tool_access_control.invocation",
						"BlockSpecific",
					),
					// Verify placeholder ID and resource_set_id are set
					resource.TestCheckResourceAttrSet(testMcpContentSecurityResource, "id"),
					resource.TestCheckResourceAttrSet(testMcpContentSecurityResource, "resource_set_id"),

					// Verify MCP access rule properties
					resource.TestCheckResourceAttr(
						testMcpAccessRuleResource,
						"tool_name",
						"test-tool",
					),
					resource.TestCheckResourceAttr(
						testMcpAccessRuleResource,
						"is_visible",
						"true",
					),
					resource.TestCheckResourceAttr(
						testMcpAccessRuleResource,
						"is_invocable",
						"false",
					),
					resource.TestCheckResourceAttrSet(testMcpAccessRuleResource, "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      testMcpContentSecurityResource,
				ImportState:       true,
				ImportStateVerify: false,
			},
			// Update and Read testing
			{
				Config: string(modifyFile),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify Name changed
					resource.TestCheckResourceAttr(
						testMcpContentSecurityResource,
						"name",
						"TF Acceptance MCP Tool Access Control Content Security Updated",
					),
					resource.TestCheckResourceAttr(
						testMcpContentSecurityResource,
						"is_active",
						"true",
					),
					resource.TestCheckResourceAttr(
						testMcpContentSecurityResource,
						"type",
						"McpToolAccessControl",
					),
					resource.TestCheckResourceAttr(
						testMcpContentSecurityResource,
						"mcp_tool_access_control.mode",
						"Block",
					),
					resource.TestCheckResourceAttr(
						testMcpContentSecurityResource,
						"mcp_tool_access_control.visibility",
						"BlockAll",
					),
					resource.TestCheckResourceAttr(
						testMcpContentSecurityResource,
						"mcp_tool_access_control.invocation",
						"BlockAll",
					),

					// Verify updated rule properties
					resource.TestCheckResourceAttr(
						testMcpAccessRuleResource,
						"tool_name",
						"test-tool-updated",
					),
					resource.TestCheckResourceAttr(
						testMcpAccessRuleResource,
						"is_visible",
						"false",
					),
					resource.TestCheckResourceAttr(
						testMcpAccessRuleResource,
						"is_invocable",
						"true",
					),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

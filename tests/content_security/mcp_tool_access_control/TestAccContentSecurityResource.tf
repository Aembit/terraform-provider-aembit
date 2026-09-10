provider "aembit" {
}

resource "aembit_content_security" "mcp_access_control" {
	name = "TF Acceptance MCP Tool Access Control Content Security"
	is_active = true
	type = "McpToolAccessControl"
	mcp_tool_access_control = {
		mode = "Allow"
		visibility = "AllowSpecific"
		invocation = "BlockSpecific"
	}
	tags = {
        color = "blue"
        day   = "Sunday"
    }
}

resource "aembit_mcp_tool_access_rule" "rule1" {
	content_security_id = aembit_content_security.mcp_access_control.id
	tool_name           = "test-tool"
	is_visible          = true
	is_invocable        = false
}

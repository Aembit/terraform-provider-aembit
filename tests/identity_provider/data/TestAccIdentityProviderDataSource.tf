provider "aembit" {
}

data "aembit_roles" "test" {}

locals {
  role_ids_by_name = { for role in data.aembit_roles.test.roles : role.name => role.id }
}

resource "aembit_identity_provider" "test_idp" {
	name = "Identity Provider for TF Acceptance Test"
	description = "Description of Identity Provider for TF Acceptance Test"
	is_active = true
    metadata_xml = "<md:EntityDescriptor xmlns:md=\"urn:oasis:names:tc:SAML:2.0:metadata\" entityID=\"https://mytest.com/saml2\"></md:EntityDescriptor>"
    saml_statement_role_mappings = [
        {
            attribute_name = "test-attribute-name"
            attribute_value = "test-attribute-value"
            roles = [local.role_ids_by_name["SuperAdmin"], local.role_ids_by_name["Auditor"]]
        }
    ]
}

data "aembit_identity_providers" "test_idps" {
    depends_on = [ aembit_identity_provider.test_idp ]
}
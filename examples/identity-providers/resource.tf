resource "aembit_identity_provider" "test_idp" {
  name         = "Sample Identity Provider Configuration"
  description  = "Identity provider description"
  is_active    = true
  metadata_xml = "<md:EntityDescriptor xmlns:md=\"urn:oasis:names:tc:SAML:2.0:metadata\" entityID=\"https://sample.test/path\"></md:EntityDescriptor>"
  sso_statement_role_mappings = [
    {
      attribute_name  = "test-attribute-name"
      attribute_value = "test-attribute-value"
      roles           = ["id_of_role1", "id_of_role2"]
    }
  ]
}
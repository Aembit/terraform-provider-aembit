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
    sso_statement_role_mappings = [
        {
            attribute_name = "test-attribute-name"
            attribute_value = "test-attribute-value"
            roles = [local.role_ids_by_name["SuperAdmin"], local.role_ids_by_name["Auditor"]]
        }
    ]
    saml = {
        metadata_xml = <<XML
<md:EntityDescriptor xmlns:md="urn:oasis:names:tc:SAML:2.0:metadata" entityID="https://aembit.test/saml" validUntil="2057-08-04T11:13:26.000Z">
    <md:IDPSSODescriptor WantAuthnRequestsSigned="false" protocolSupportEnumeration="urn:oasis:names:tc:SAML:2.0:protocol">
        <md:KeyDescriptor use="signing">
        <ds:KeyInfo xmlns:ds="http://www.w3.org/2000/09/xmldsig#">
            <ds:X509Data>
            <ds:X509Certificate>
                MIIDwjCCAqqgAwIBAgIQRtjIzhUI3pVK6FKOIQV4DDANBgkqhkiG9w0BAQsFADAX
                MRUwEwYDVQQDDAxhZW1iaXQubG9jYWwwIBcNMjQxMTEyMTUyMTQxWhgPMjA1NDEx
                MTIxNTMxNDBaMBcxFTATBgNVBAMMDGFlbWJpdC5sb2NhbDCCASIwDQYJKoZIhvcN
                AQEBBQADggEPADCCAQoCggEBAK4aB+6xe8YADZulXtpqG4KSavjBt2aYi13GfdiH
                UkFeak3oriA/0DecjF5T+PSkawMBIYh7WDRCQz3k7JGXvhzIgviV3zyxcZfcawH5
                ihXQlwp01T5IuiRvXMtqD8aIfioufZUcLNs1NPMbm6soF7OauD1dKFTmqM7rKP9v
                e9J1/Iutu4s00kKcyQ/ZLirfHfqQ/lLW8RiC8XEwhUMsPXl2pJroZSf4S2DMwvPx
                ixmIHTzcUDUFvt8Hp6ssZFV9AAFC0P/o1UQf3zMCY+FNRdncMADzbSGSwGCmgQz5
                HbC5Rp8FiWuH2nfZvmaXoGD8kWX/vzy54n90+JaR7bQIdPUCAwEAAaOCAQYwggEC
                MA4GA1UdDwEB/wQEAwIFoDAdBgNVHSUEFjAUBggrBgEFBQcDAgYIKwYBBQUHAwEw
                gbEGA1UdEQSBqTCBpoIMYWVtYml0LmxvY2Fsgg4qLmFlbWJpdC5sb2NhbIISKi5h
                cHAuYWVtYml0LmxvY2FsghIqLmFwaS5hZW1iaXQubG9jYWyCESouaWQuYWVtYml0
                LmxvY2FsghEqLmVjLmFlbWJpdC5sb2NhbIIbKi5pZC51cy1lYXN0LTEuYWVtYml0
                LmxvY2FsghsqLmVjLnVzLWVhc3QtMS5hZW1iaXQubG9jYWwwHQYDVR0OBBYEFNDd
                KCH4MSKtcmGo9P/TDKKvGWYWMA0GCSqGSIb3DQEBCwUAA4IBAQBw8C+h3G+cwZ3g
                z5dV5CIANdnPCFhnW/9HpYkyZm/LAcxT+xLXxOiHQoqHwvQt2m9qZ+XnMV/7wIys
                1Y9NdQgXEjwosnRSf3rmWIBrvJA7K9W2/6mv+CHF+3AqHNFcCJfBuNptuphng5vm
                jQU9OaOUg89O+st4RR9d5rxxVCHAqtyIkroD198tnduAqAYV7Sg74MxWS/zUANPm
                EdNeMIMACOc7FQZFR1RB55Kd3/sl2WEIwxojW66YEIHQ1+F36zX+jwYJ7eQTBiMu
                yPx+z13vQLRM3SfJiccOPTEwr90R5YpZVWGazxYOgCP35Qqho6hq5nvM0T+9Ii9B
                Gqg7Ey2L
            </ds:X509Certificate>
            </ds:X509Data>
        </ds:KeyInfo>
        </md:KeyDescriptor>
        <md:NameIDFormat>urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress</md:NameIDFormat>
        <md:SingleSignOnService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect" Location="https://aembit.test/saml"/>
        <md:SingleSignOnService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST" Location="https://aembit.test/saml"/>
    </md:IDPSSODescriptor>
</md:EntityDescriptor>
XML        
    }
}

data "aembit_identity_providers" "test_idps" {
    depends_on = [ aembit_identity_provider.test_idp ]
}
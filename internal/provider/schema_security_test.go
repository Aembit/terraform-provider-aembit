package provider

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testingReporter abstracts test error reporting, allowing production tests to pass *testing.T
// and negative validation tests to pass a recording mock.
type testingReporter interface {
	// Errorf formats and records a failure diagnostic.
	Errorf(format string, args ...any)
}

// testFailureRecorder captures error messages during negative test assertions.
type testFailureRecorder struct {
	// errors contains all recorded failure messages.
	errors []string
}

// Errorf formats and appends an error message to the recorder's slice.
func (r *testFailureRecorder) Errorf(format string, args ...any) {
	r.errors = append(r.errors, fmt.Sprintf(format, args...))
}

// discoveredAttribute stores schema metadata collected during recursive traversal.
type discoveredAttribute struct {
	// sourceKind specifies the origin category: "resource", "data_source", or "provider".
	sourceKind string
	// entityName specifies the resource or data source type name (e.g., "aembit_discovery_integration").
	entityName string
	// path specifies the dot-delimited attribute path within the schema (e.g., "wiz_integration.client_secret").
	path string
	// isSensitive indicates whether the attribute has Sensitive: true in its schema definition.
	isSensitive bool
	// isWriteOnly indicates whether the attribute has WriteOnly: true in its schema definition.
	isWriteOnly bool
	// isContainer indicates whether the attribute contains nested attributes or blocks rather than a primitive value.
	isContainer bool
	// description provides the attribute's human-readable documentation.
	description string
}

// schemaExemption documents a known non-sensitive attribute matching sensitive naming patterns.
type schemaExemption struct {
	// sourceKind indicates whether the exemption targets a "resource", "data_source", or "provider".
	sourceKind string
	// entityName indicates the schema type name (e.g., "aembit_credential_provider").
	entityName string
	// attributePath indicates the dot-delimited path of the exempted attribute.
	attributePath string
	// reason details the security and architectural rationale explaining why this attribute is non-sensitive.
	reason string
}

// schemaSecurityValidator coordinates schema inspection, sensitive keyword matching, and exemption verification.
type schemaSecurityValidator struct {
	// logger provides structured, context-aware diagnostic logging.
	logger *slog.Logger
	// sensitiveKeywordPattern matches attribute naming conventions that typically signify sensitive credentials.
	sensitiveKeywordPattern *regexp.Regexp
	// exemptions maps compound keys ("sourceKind:entityName:attributePath") to documented schema exemptions.
	exemptions map[string]schemaExemption
	// seenExemptions tracks how many times each exemption was matched during active schema traversal.
	seenExemptions map[string]int
}

// newSchemaSecurityValidator constructs and initializes a schemaSecurityValidator instance with standard rules.
// It maintains zero package-level globals by encapsulating regex compilation and the exemption registry.
func newSchemaSecurityValidator(logger *slog.Logger) *schemaSecurityValidator {
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	}

	// sensitiveKeywordPattern captures leaf names containing secret, password, private_key, api_key,
	// credential, or token.
	sensitiveKeywordPattern := regexp.MustCompile(`(?i)(secret|password|private_key|api_key|credential|token)`)

	exemptions := make(map[string]schemaExemption)

	// registerExemption is a local helper function to register documented exemptions into the validator.
	registerExemption := func(sourceKind string, entityName string, path string, reason string) {
		key := fmt.Sprintf("%s:%s:%s", sourceKind, entityName, path)
		exemptions[key] = schemaExemption{
			sourceKind:    sourceKind,
			entityName:    entityName,
			attributePath: path,
			reason:        reason,
		}
	}

	// -------------------------------------------------------------------------
	// Documented Schema Exemptions
	// Every exemption below must exist in the active schema. Any obsolete entry
	// will trigger a test failure during assertNoStaleExemptions.
	// -------------------------------------------------------------------------

	// AWS Secrets Manager: ARNs and JSON key names are structural identifiers, not secret content.
	registerExemption("resource", "aembit_credential_provider", "aws_secrets_manager_value.secret_arn",
		"AWS Secrets Manager ARN is an AWS resource identifier string, not the secret payload.")
	registerExemption("resource", "aembit_credential_provider", "aws_secrets_manager_value.secret_key_1",
		"Identifies the JSON dictionary key within the secret payload, not the secret payload itself.")
	registerExemption("resource", "aembit_credential_provider", "aws_secrets_manager_value.secret_key_2",
		"Identifies the second JSON dictionary key within the secret payload, not the secret payload itself.")
	registerExemption("data_source", "aembit_credential_providers", "credential_providers.aws_secrets_manager_value.secret_arn",
		"AWS Secrets Manager ARN is an AWS resource identifier string, not the secret payload.")
	registerExemption("data_source", "aembit_credential_providers", "credential_providers.aws_secrets_manager_value.secret_key_1",
		"Identifies the JSON dictionary key within the secret payload, not the secret payload itself.")
	registerExemption("data_source", "aembit_credential_providers", "credential_providers.aws_secrets_manager_value.secret_key_2",
		"Identifies the second JSON dictionary key within the secret payload, not the secret payload itself.")

	// Azure Key Vault: Secret names are vault identifiers, not secret content.
	registerExemption("resource", "aembit_credential_provider", "azure_key_vault_value.secret_name_1",
		"Identifies the secret name within the Azure Key Vault, not the secret value.")
	registerExemption("resource", "aembit_credential_provider", "azure_key_vault_value.secret_name_2",
		"Identifies the second secret name within the Azure Key Vault, not the secret value.")
	registerExemption("data_source", "aembit_credential_providers", "credential_providers.azure_key_vault_value.secret_name_1",
		"Identifies the secret name within the Azure Key Vault, not the secret value.")
	registerExemption("data_source", "aembit_credential_providers", "credential_providers.azure_key_vault_value.secret_name_2",
		"Identifies the second secret name within the Azure Key Vault, not the secret value.")

	// Credential Provider Integrations: Lists of authorized ARNs and secret names.
	registerExemption("resource", "aembit_credential_provider_integration", "aws_iam_role.fetch_secret_arns",
		"List of AWS IAM Role ARNs permitted to fetch secrets, not the secret payload.")
	registerExemption("data_source", "aembit_credential_provider_integrations", "credential_provider_integrations.aws_iam_role.fetch_secret_arns",
		"List of AWS IAM Role ARNs permitted to fetch secrets, not the secret payload.")
	registerExemption("resource", "aembit_credential_provider_integration", "azure_entra_federation.fetch_secret_names",
		"List of Azure Key Vault secret names permitted to be fetched, not the secret payload.")
	registerExemption("data_source", "aembit_credential_provider_integrations", "credential_provider_integrations.azure_entra_federation.fetch_secret_names",
		"List of Azure Key Vault secret names permitted to be fetched, not the secret payload.")

	// Entity References: IDs linking policies and providers to other entities.
	registerExemption("resource", "aembit_credential_provider", "aws_secrets_manager_value.credential_provider_integration_id",
		"Entity reference UUID linking to the integration, not credential data.")
	registerExemption("resource", "aembit_credential_provider", "azure_key_vault_value.credential_provider_integration_id",
		"Entity reference UUID linking to the integration, not credential data.")
	registerExemption("resource", "aembit_credential_provider", "managed_gitlab_account.credential_provider_integration_id",
		"Entity reference UUID linking to the integration, not credential data.")
	registerExemption("data_source", "aembit_credential_providers", "credential_providers.aws_secrets_manager_value.credential_provider_integration_id",
		"Entity reference UUID linking to the integration, not credential data.")
	registerExemption("data_source", "aembit_credential_providers", "credential_providers.azure_key_vault_value.credential_provider_integration_id",
		"Entity reference UUID linking to the integration, not credential data.")
	registerExemption("data_source", "aembit_credential_providers", "credential_providers.managed_gitlab_account.credential_provider_integration_id",
		"Entity reference UUID linking to the integration, not credential data.")
	registerExemption("resource", "aembit_access_policy", "credential_provider",
		"Entity reference UUID linking an access policy to a credential provider.")
	registerExemption("resource", "aembit_access_policy", "credential_providers.credential_provider_id",
		"Entity reference UUID linking an access policy to a credential provider.")
	registerExemption("data_source", "aembit_access_policies", "access_policies.credential_provider",
		"Entity reference UUID linking an access policy to a credential provider.")
	registerExemption("data_source", "aembit_access_policies", "access_policies.credential_providers.credential_provider_id",
		"Entity reference UUID linking an access policy to a credential provider.")

	// OAuth Configuration Enums: Request encoding style configuration.
	registerExemption("resource", "aembit_credential_provider", "oauth_client_credentials.credential_style",
		"Configuration enum specifying parameter passing method (InHeader vs InParams), not credential material.")
	registerExemption("data_source", "aembit_credential_providers", "credential_providers.oauth_client_credentials.credential_style",
		"Configuration enum specifying parameter passing method (InHeader vs InParams), not credential material.")

	// Token URLs: Public or corporate endpoint URLs.
	registerExemption("resource", "aembit_credential_provider", "mcp_ema.token_url",
		"OAuth token endpoint URL, not the token itself.")
	registerExemption("data_source", "aembit_credential_providers", "credential_providers.mcp_ema.token_url",
		"OAuth token endpoint URL, not the token itself.")
	registerExemption("resource", "aembit_credential_provider", "mcp_user_based_access_token.oauth_token_url",
		"OAuth token endpoint URL, not the token itself.")
	registerExemption("data_source", "aembit_credential_providers", "credential_providers.mcp_user_based_access_token.oauth_token_url",
		"OAuth token endpoint URL, not the token itself.")
	registerExemption("resource", "aembit_credential_provider", "oauth_authorization_code.oauth_token_url",
		"OAuth token endpoint URL, not the token itself.")
	registerExemption("data_source", "aembit_credential_providers", "credential_providers.oauth_authorization_code.oauth_token_url",
		"OAuth token endpoint URL, not the token itself.")
	registerExemption("resource", "aembit_credential_provider", "oauth_client_credentials.token_url",
		"OAuth token endpoint URL, not the token itself.")
	registerExemption("data_source", "aembit_credential_providers", "credential_providers.oauth_client_credentials.token_url",
		"OAuth token endpoint URL, not the token itself.")
	registerExemption("resource", "aembit_discovery_integration", "wiz_integration.token_url",
		"Wiz OAuth token endpoint URL, not the token itself.")
	registerExemption("data_source", "aembit_discovery_integrations", "discovery_integrations.wiz_integration.token_url",
		"Wiz OAuth token endpoint URL, not the token itself.")
	registerExemption("resource", "aembit_integration", "oauth_client_credentials.token_url",
		"OAuth token endpoint URL, not the token itself.")
	registerExemption("data_source", "aembit_integrations", "integrations.oauth_client_credentials.token_url",
		"OAuth token endpoint URL, not the token itself.")

	// Token Lifetimes & Durations: Integer time values.
	registerExemption("resource", "aembit_credential_provider", "aembit_access_token.absolute_token_lifetime",
		"Token validity duration in seconds/minutes, not the token itself.")
	registerExemption("data_source", "aembit_credential_providers", "credential_providers.aembit_access_token.absolute_token_lifetime",
		"Token validity duration in seconds/minutes, not the token itself.")
	registerExemption("resource", "aembit_credential_provider", "oidc_id_token.absolute_token_lifetime",
		"Token validity duration in seconds/minutes, not the token itself.")
	registerExemption("data_source", "aembit_credential_providers", "credential_providers.oidc_id_token.absolute_token_lifetime",
		"Token validity duration in seconds/minutes, not the token itself.")
	registerExemption("resource", "aembit_log_stream", "gcs_bucket.token_lifetime",
		"Token validity duration integer, not the token itself.")
	registerExemption("data_source", "aembit_log_streams", "log_streams.gcs_bucket.token_lifetime",
		"Token validity duration integer, not the token itself.")

	// Token Audiences: Public audience identifier strings.
	registerExemption("resource", "aembit_credential_provider", "aws_sts.token_audience",
		"Target audience string for AWS STS token minting, not the token itself.")
	registerExemption("data_source", "aembit_credential_providers", "credential_providers.aws_sts.token_audience",
		"Target audience string for AWS STS token minting, not the token itself.")
	registerExemption("resource", "aembit_credential_provider_integration", "aws_iam_role.token_audience",
		"Target audience string for AWS IAM Role federation, not the token itself.")
	registerExemption("data_source", "aembit_credential_provider_integrations", "credential_provider_integrations.aws_iam_role.token_audience",
		"Target audience string for AWS IAM Role federation, not the token itself.")

	// Boolean Token Flags: Configuration booleans.
	registerExemption("resource", "aembit_trust_provider", "kubernetes_service_account.is_aembit_tenant_oidc_token",
		"Boolean flag indicating whether the OIDC token is minted by Aembit tenant, not token payload.")
	registerExemption("data_source", "aembit_trust_providers", "trust_providers.kubernetes_service_account.is_aembit_tenant_oidc_token",
		"Boolean flag indicating whether the OIDC token is minted by Aembit tenant, not token payload.")
	registerExemption("resource", "aembit_trust_provider", "oidc_id_token.is_aembit_tenant_oidc_token",
		"Boolean flag indicating whether the OIDC token is minted by Aembit tenant, not token payload.")
	registerExemption("data_source", "aembit_trust_providers", "trust_providers.oidc_id_token.is_aembit_tenant_oidc_token",
		"Boolean flag indicating whether the OIDC token is minted by Aembit tenant, not token payload.")

	return &schemaSecurityValidator{
		logger:                  logger,
		sensitiveKeywordPattern: sensitiveKeywordPattern,
		exemptions:              exemptions,
		seenExemptions:          make(map[string]int),
	}
}

// walkSchemaAttributes recursively inspects a slice of tfprotov6.SchemaAttribute elements.
func (v *schemaSecurityValidator) walkSchemaAttributes(
	sourceKind string,
	entityName string,
	prefix string,
	attrs []*tfprotov6.SchemaAttribute,
	out *[]discoveredAttribute,
) {
	for _, attr := range attrs {
		if attr == nil {
			continue
		}
		fullPath := attr.Name
		if prefix != "" {
			fullPath = prefix + "." + attr.Name
		}

		isContainer := attr.NestedType != nil

		*out = append(*out, discoveredAttribute{
			sourceKind:  sourceKind,
			entityName:  entityName,
			path:        fullPath,
			isSensitive: attr.Sensitive,
			isWriteOnly: attr.WriteOnly,
			isContainer: isContainer,
			description: attr.Description,
		})

		// Recurse into nested object attributes if present.
		if attr.NestedType != nil {
			v.walkSchemaAttributes(sourceKind, entityName, fullPath, attr.NestedType.Attributes, out)
		}
	}
}

// walkSchemaBlock recursively inspects a tfprotov6.SchemaBlock, its attributes, and nested blocks.
func (v *schemaSecurityValidator) walkSchemaBlock(
	sourceKind string,
	entityName string,
	prefix string,
	block *tfprotov6.SchemaBlock,
	out *[]discoveredAttribute,
) {
	if block == nil {
		return
	}

	v.walkSchemaAttributes(sourceKind, entityName, prefix, block.Attributes, out)

	for _, nb := range block.BlockTypes {
		if nb == nil {
			continue
		}
		fullPath := nb.TypeName
		if prefix != "" {
			fullPath = prefix + "." + nb.TypeName
		}
		v.walkSchemaBlock(sourceKind, entityName, fullPath, nb.Block, out)
	}
}

// validateAttribute inspects a single discovered attribute against sensitive keyword rules and exemptions.
func (v *schemaSecurityValidator) validateAttribute(reporter testingReporter, a discoveredAttribute) {
	// Structural containers (objects or blocks holding children) do not hold primitive values.
	if a.isContainer {
		return
	}

	parts := strings.Split(a.path, ".")
	leafName := parts[len(parts)-1]

	// Check if leaf attribute matches sensitive naming patterns.
	if !v.sensitiveKeywordPattern.MatchString(leafName) {
		return
	}

	// If marked Sensitive or WriteOnly, the credential is secure in Terraform plans and CLI outputs.
	if a.isSensitive || a.isWriteOnly {
		v.logger.LogAttrs(
			context.Background(),
			slog.LevelDebug,
			"verified sensitive attribute tagged",
			slog.String("kind", a.sourceKind),
			slog.String("entity", a.entityName),
			slog.String("path", a.path),
		)
		return
	}

	// Check if this attribute is an authorized exemption.
	key := fmt.Sprintf("%s:%s:%s", a.sourceKind, a.entityName, a.path)
	if _, ok := v.exemptions[key]; ok {
		v.seenExemptions[key]++
		v.logger.LogAttrs(
			context.Background(),
			slog.LevelDebug,
			"verified authorized schema exemption",
			slog.String("key", key),
		)
		return
	}

	// Failure: Attribute matches sensitive keywords but is neither sensitive, write-only, nor exempted.
	reporter.Errorf(
		"Untagged sensitive attribute detected:\n"+
			"  Source:    %s\n"+
			"  Entity:    %s\n"+
			"  Path:      %s\n"+
			"  Leaf:      %s\n"+
			"  Issue:     Attribute matches sensitive pattern '%s' but has Sensitive=false and WriteOnly=false.\n"+
			"  Remediation:\n"+
			"    1. If this attribute holds secret, credential, or sensitive token material, set 'Sensitive: true' (or WriteOnly: true) in the schema definition.\n"+
			"    2. If this attribute is a non-sensitive identifier, configuration enum, or metadata, register an exemption with justification in newSchemaSecurityValidator().",
		a.sourceKind,
		a.entityName,
		a.path,
		leafName,
		v.sensitiveKeywordPattern.String(),
	)
}

// assertNoStaleExemptions verifies that every registered exemption corresponds to an active schema attribute.
func (v *schemaSecurityValidator) assertNoStaleExemptions(reporter testingReporter) {
	var staleKeys []string
	for key := range v.exemptions {
		if v.seenExemptions[key] == 0 {
			staleKeys = append(staleKeys, key)
		}
	}
	sort.Strings(staleKeys)

	if len(staleKeys) > 0 {
		reporter.Errorf(
			"Found %d obsolete or misconfigured exemptions in schemaSecurityValidator that do not match any active attribute in the provider schema:\n%s",
			len(staleKeys),
			strings.Join(staleKeys, "\n"),
		)
	}
}

// TestUnitSchema_SensitiveAttributesTagging validates all provider resources, data sources,
// and provider configurations to ensure every secret, credential, or token attribute is
// explicitly tagged with Sensitive: true or WriteOnly: true.
func TestUnitSchema_SensitiveAttributesTagging(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	validator := newSchemaSecurityValidator(logger)

	// Obtain schema from provider server implementation.
	server, err := providerserver.NewProtocol6WithError(New("test", "unittest")())()
	require.NoError(t, err, "failed to instantiate protocol 6 provider server")

	resp, err := server.GetProviderSchema(ctx, &tfprotov6.GetProviderSchemaRequest{})
	require.NoError(t, err, "GetProviderSchema failed")
	require.Empty(t, resp.Diagnostics, "GetProviderSchema returned diagnostics: %v", resp.Diagnostics)

	var allAttrs []discoveredAttribute

	// 1. Inspect Provider Configuration Schema
	if resp.Provider != nil {
		t.Run("provider/aembit", func(t *testing.T) {
			var providerAttrs []discoveredAttribute
			validator.walkSchemaBlock("provider", "aembit", "", resp.Provider.Block, &providerAttrs)
			for _, a := range providerAttrs {
				validator.validateAttribute(t, a)
			}
			allAttrs = append(allAttrs, providerAttrs...)
		})
	}

	// 2. Inspect Managed Resources
	for resName, resSchema := range resp.ResourceSchemas {
		resName := resName
		resSchema := resSchema
		t.Run("resource/"+resName, func(t *testing.T) {
			var resAttrs []discoveredAttribute
			validator.walkSchemaBlock("resource", resName, "", resSchema.Block, &resAttrs)
			for _, a := range resAttrs {
				validator.validateAttribute(t, a)
			}
			allAttrs = append(allAttrs, resAttrs...)
		})
	}

	// 3. Inspect Data Sources
	for dsName, dsSchema := range resp.DataSourceSchemas {
		dsName := dsName
		dsSchema := dsSchema
		t.Run("data_source/"+dsName, func(t *testing.T) {
			var dsAttrs []discoveredAttribute
			validator.walkSchemaBlock("data_source", dsName, "", dsSchema.Block, &dsAttrs)
			for _, a := range dsAttrs {
				validator.validateAttribute(t, a)
			}
			allAttrs = append(allAttrs, dsAttrs...)
		})
	}

	// 4. Assert Zero Stale Exemptions
	t.Run("assert_no_stale_exemptions", func(t *testing.T) {
		validator.assertNoStaleExemptions(t)
	})

	t.Logf("Schema security audit completed successfully. Total attributes verified: %d", len(allAttrs))
}

// TestUnitSchema_SensitiveAttributesTagging_Negative executes adversarial test cases
// asserting that unmarked sensitive attributes and obsolete exemptions are caught and reported.
func TestUnitSchema_SensitiveAttributesTagging_Negative(t *testing.T) {
	t.Parallel()

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	// Case 1: Untagged sensitive attribute must trigger a failure diagnostic.
	t.Run("untagged_client_secret_triggers_failure", func(t *testing.T) {
		validator := newSchemaSecurityValidator(logger)
		recorder := &testFailureRecorder{}

		unmarkedAttr := discoveredAttribute{
			sourceKind:  "resource",
			entityName:  "mock_integration",
			path:        "oauth_config.client_secret",
			isSensitive: false,
			isWriteOnly: false,
			isContainer: false,
			description: "Mock OAuth client secret",
		}

		validator.validateAttribute(recorder, unmarkedAttr)

		require.Len(t, recorder.errors, 1)
		assert.Contains(t, recorder.errors[0], "Untagged sensitive attribute detected")
		assert.Contains(t, recorder.errors[0], "mock_integration")
		assert.Contains(t, recorder.errors[0], "oauth_config.client_secret")
		assert.Contains(t, recorder.errors[0], "Sensitive: true")
	})

	// Case 2: Untagged password must trigger a failure diagnostic.
	t.Run("untagged_password_triggers_failure", func(t *testing.T) {
		validator := newSchemaSecurityValidator(logger)
		recorder := &testFailureRecorder{}

		unmarkedAttr := discoveredAttribute{
			sourceKind:  "resource",
			entityName:  "mock_database",
			path:        "credentials.admin_password",
			isSensitive: false,
			isWriteOnly: false,
			isContainer: false,
			description: "Mock database admin password",
		}

		validator.validateAttribute(recorder, unmarkedAttr)

		require.Len(t, recorder.errors, 1)
		assert.Contains(t, recorder.errors[0], "Untagged sensitive attribute detected")
		assert.Contains(t, recorder.errors[0], "admin_password")
	})

	// Case 3: Untagged personal access token must trigger a failure diagnostic.
	t.Run("untagged_access_token_triggers_failure", func(t *testing.T) {
		validator := newSchemaSecurityValidator(logger)
		recorder := &testFailureRecorder{}

		unmarkedAttr := discoveredAttribute{
			sourceKind:  "resource",
			entityName:  "mock_service",
			path:        "auth.personal_access_token",
			isSensitive: false,
			isWriteOnly: false,
			isContainer: false,
			description: "Mock personal access token",
		}

		validator.validateAttribute(recorder, unmarkedAttr)

		require.Len(t, recorder.errors, 1)
		assert.Contains(t, recorder.errors[0], "personal_access_token")
	})

	// Case 4: Correctly tagged Sensitive: true attribute must pass without error.
	t.Run("tagged_sensitive_attribute_passes", func(t *testing.T) {
		validator := newSchemaSecurityValidator(logger)
		recorder := &testFailureRecorder{}

		taggedAttr := discoveredAttribute{
			sourceKind:  "resource",
			entityName:  "mock_integration",
			path:        "oauth_config.client_secret",
			isSensitive: true,
			isWriteOnly: false,
			isContainer: false,
		}

		validator.validateAttribute(recorder, taggedAttr)
		assert.Empty(t, recorder.errors)
	})

	// Case 5: Correctly tagged WriteOnly: true attribute must pass without error.
	t.Run("tagged_write_only_attribute_passes", func(t *testing.T) {
		validator := newSchemaSecurityValidator(logger)
		recorder := &testFailureRecorder{}

		writeOnlyAttr := discoveredAttribute{
			sourceKind:  "resource",
			entityName:  "mock_integration",
			path:        "credentials.api_key",
			isSensitive: false,
			isWriteOnly: true,
			isContainer: false,
		}

		validator.validateAttribute(recorder, writeOnlyAttr)
		assert.Empty(t, recorder.errors)
	})

	// Case 6: Asserting stale exemptions triggers failure when an exemption is never encountered in the schema.
	t.Run("stale_exemption_triggers_assertion_failure", func(t *testing.T) {
		validator := newSchemaSecurityValidator(logger)
		recorder := &testFailureRecorder{}

		// Reset exemptions to isolate the stale check to the simulated exemption.
		validator.exemptions = make(map[string]schemaExemption)
		validator.exemptions["resource:fake_resource:fake_secret_name"] = schemaExemption{
			sourceKind:    "resource",
			entityName:    "fake_resource",
			attributePath: "fake_secret_name",
			reason:        "Simulated stale exemption for negative test assertion.",
		}

		validator.assertNoStaleExemptions(recorder)

		require.Len(t, recorder.errors, 1)
		assert.Contains(t, recorder.errors[0], "Found 1 obsolete or misconfigured exemptions")
		assert.Contains(t, recorder.errors[0], "resource:fake_resource:fake_secret_name")
	})
}

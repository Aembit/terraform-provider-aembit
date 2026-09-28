package provider

import (
	"os"
	"regexp"
	"testing"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"terraform-provider-aembit/internal/provider/models"
)

const (
	AccessPolicyPathFirst           string = "aembit_access_policy.first_policy"
	AccessPolicyPathSecond          string = "aembit_access_policy.second_policy"
	AccessPolicyPathMultiCPFirst    string = "aembit_access_policy.multi_cp_first_policy"
	AccessPolicyPathMultiCPSecond   string = "aembit_access_policy.multi_cp_second_policy"
	AccessPolicyPathMultiStsCPFirst string = "aembit_access_policy.multi_sts_cp_first_policy"
	CredentialProvidersCount        string = "credential_providers.#"
)

var accessPolicyChecks = []resource.TestCheckFunc{
	// Verify values for First Policy.
	resource.TestCheckResourceAttrSet(AccessPolicyPathFirst, "id"),
	resource.TestCheckResourceAttrSet(AccessPolicyPathFirst, "client_workload"),
	resource.TestCheckResourceAttrSet(AccessPolicyPathFirst, "credential_provider"),
	resource.TestCheckResourceAttrSet(AccessPolicyPathFirst, "server_workload"),
}

func TestAccAccessPolicyResource(t *testing.T) {
	t.Parallel()
	createFile, _ := os.ReadFile("../../tests/policy/TestAccAccessPolicyResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/policy/TestAccAccessPolicyResource.tfmod")
	createFileConfig, modifyFileConfig, _ := randomizeFileConfigs(
		string(createFile),
		string(modifyFile),
		"clientworkloadNamespace",
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: createFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					accessPolicyChecks...,
				),
			},
			// ImportState testing
			{ResourceName: AccessPolicyPathFirst, ImportState: true, ImportStateVerify: true},
			// Update and Read testing
			{
				Config: modifyFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					accessPolicyChecks...,
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

var basicAccessPolicyChecks = []resource.TestCheckFunc{
	// Verify values for First Policy.
	resource.TestCheckResourceAttrSet(AccessPolicyPathFirst, "id"),
	resource.TestCheckResourceAttrSet(AccessPolicyPathFirst, "client_workload"),
	resource.TestCheckResourceAttr(AccessPolicyPathFirst, "trust_providers.#", "0"),
	resource.TestCheckResourceAttr(AccessPolicyPathFirst, "access_conditions.#", "0"),
	resource.TestCheckResourceAttr(AccessPolicyPathFirst, "content_security.#", "0"),
	resource.TestCheckResourceAttrSet(AccessPolicyPathFirst, "server_workload"),

	// Verify values for Second Policy.
	resource.TestCheckResourceAttrSet(AccessPolicyPathSecond, "id"),
	resource.TestCheckResourceAttrSet(AccessPolicyPathSecond, "client_workload"),
	resource.TestCheckResourceAttrSet(AccessPolicyPathSecond, "credential_provider"),
	resource.TestCheckResourceAttrSet(AccessPolicyPathSecond, "server_workload"),
}

func TestAccBasicAccessPolicyResource(t *testing.T) {
	t.Parallel()
	createFile, _ := os.ReadFile("../../tests/policy/TestAccBasicAccessPolicyResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/policy/TestAccBasicAccessPolicyResource.tfmod")
	createFileConfig, modifyFileConfig, _ := randomizeFileConfigs(
		string(createFile),
		string(modifyFile),
		"clientworkloadNamespace",
	)
	createFileConfig, modifyFileConfig, _ = randomizeFileConfigs(
		createFileConfig,
		modifyFileConfig,
		"secondClientWorkloadNamespace",
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: createFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					append(
						basicAccessPolicyChecks,
						resource.TestCheckResourceAttr(
							AccessPolicyPathFirst,
							"name",
							"TF First Policy",
						),
						resource.TestCheckResourceAttr(
							AccessPolicyPathSecond,
							"name",
							"TF Second Policy",
						),
					)...,
				),
			},
			// ImportState testing
			{ResourceName: AccessPolicyPathFirst, ImportState: true, ImportStateVerify: true},
			// Update and Read testing
			{
				Config: modifyFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					append(
						basicAccessPolicyChecks,
						resource.TestCheckResourceAttr(
							AccessPolicyPathFirst,
							"name",
							"Placeholder",
						),
						resource.TestCheckResourceAttr(
							AccessPolicyPathSecond,
							"name",
							"Placeholder",
						),
					)...,
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccMultipleCredentialProviders_AccessPolicyResource(t *testing.T) {
	t.Parallel()
	createFile, _ := os.ReadFile("../../tests/policy/TestAccMultipleCPAccessPolicyResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/policy/TestAccMultipleCPAccessPolicyResource.tfmod")
	createFileConfig, modifyFileConfig, _ := randomizeFileConfigs(
		string(createFile),
		string(modifyFile),
		"clientworkloadNamespace",
	)
	createFileConfig, modifyFileConfig, _ = randomizeFileConfigs(
		createFileConfig,
		modifyFileConfig,
		"secondClientWorkloadNamespace",
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: createFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						AccessPolicyPathMultiCPFirst,
						"name",
						"TF Multi CP First Policy",
					),
					resource.TestCheckResourceAttr(
						AccessPolicyPathMultiCPFirst,
						CredentialProvidersCount,
						"2",
					),

					resource.TestCheckResourceAttr(
						AccessPolicyPathMultiCPSecond,
						"name",
						"TF Multi CP Second Policy",
					),
					resource.TestCheckResourceAttr(
						AccessPolicyPathMultiCPSecond,
						CredentialProvidersCount,
						"3",
					),
				),
			},
			// ImportState testing
			{
				ResourceName:      AccessPolicyPathMultiCPFirst,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: modifyFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						AccessPolicyPathMultiCPFirst,
						"name",
						"TF Multi CP First Policy Updated",
					),
					resource.TestCheckResourceAttr(
						AccessPolicyPathMultiCPFirst,
						CredentialProvidersCount,
						"2",
					),

					resource.TestCheckResourceAttr(
						AccessPolicyPathMultiCPSecond,
						"name",
						"TF Multi CP Second Policy Updated",
					),
					resource.TestCheckResourceAttr(
						AccessPolicyPathMultiCPSecond,
						CredentialProvidersCount,
						"2",
					),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccMultipleCPAccessPolicyResource_ErrorDuplicateMappings_Create(t *testing.T) {
	t.Parallel()
	createFile, _ := os.ReadFile("../../tests/policy/TestAccAccessPolicyDuplicateMappings.tf")
	createFileConfig, _, _ := randomizeFileConfigs(
		string(createFile),
		"",
		"clientworkloadNamespace",
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: createFileConfig,
				ExpectError: regexp.MustCompile(
					`duplicate credential provider mapping already exists`,
				),
			},
		},
	})
}

func TestAccMultipleSTSCredentialProviders_AccessPolicyResource(t *testing.T) {
	t.Parallel()
	createFile, _ := os.ReadFile("../../tests/policy/TestAccMultipleCPAccessPolicyResource.tf")
	modifyFile, _ := os.ReadFile("../../tests/policy/TestAccMultipleCPAccessPolicyResource.tfmod")
	createFileConfig, modifyFileConfig, _ := randomizeFileConfigs(
		string(createFile),
		string(modifyFile),
		"clientworkloadNamespace",
	)
	createFileConfig, modifyFileConfig, _ = randomizeFileConfigs(
		createFileConfig,
		modifyFileConfig,
		"secondClientWorkloadNamespace",
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: createFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						AccessPolicyPathMultiStsCPFirst,
						"name",
						"TF Multi STS CP First Policy",
					),
					resource.TestCheckResourceAttr(
						AccessPolicyPathMultiStsCPFirst,
						CredentialProvidersCount,
						"2",
					),
				),
			},
			// ImportState testing
			{
				ResourceName:      AccessPolicyPathMultiCPFirst,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: modifyFileConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						AccessPolicyPathMultiStsCPFirst,
						"name",
						"TF Multi STS CP First Policy Updated",
					),
					resource.TestCheckResourceAttr(
						AccessPolicyPathMultiStsCPFirst,
						CredentialProvidersCount,
						"2",
					),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestConvertAccessPolicyModelToPolicyDTO_CredentialProvidersNullable(t *testing.T) {
	t.Parallel()

	model := models.AccessPolicyResourceModel{
		Name:           types.StringValue("Test Policy"),
		IsActive:       types.BoolValue(true),
		ClientWorkload: types.StringValue("cw-1"),
		ServerWorkload: types.StringValue("sw-1"),
		CredentialProviders: []*models.PolicyCredentialMappingModel{
			{
				CredentialProviderId: types.StringValue("cp-1"),
				MappingType:          types.StringValue("AccessKeyId"),
				AccessKeyId:          types.StringValue("ACCESSKEY1"),
				AccountName:          types.StringNull(),
				HeaderName:           types.StringNull(),
				HeaderValue:          types.StringNull(),
				HttpbodyFieldPath:    types.StringNull(),
				HttpbodyFieldValue:   types.StringNull(),
			},
			{
				CredentialProviderId: types.StringValue("cp-2"),
				MappingType:          types.StringValue("HttpHeader"),
				AccessKeyId:          types.StringNull(),
				AccountName:          types.StringNull(),
				HeaderName:           types.StringValue("X-Test-Header"),
				HeaderValue:          types.StringValue("test-val"),
				HttpbodyFieldPath:    types.StringNull(),
				HttpbodyFieldValue:   types.StringNull(),
			},
		},
	}

	dto := convertAccessPolicyModelToPolicyDTO(model, nil)

	require.Len(t, dto.CredentialProviders, 2)

	// First mapping (AccessKeyId)
	assert.Equal(t, "cp-1", dto.CredentialProviders[0].CredentialProviderId)
	assert.Equal(t, "AccessKeyId", dto.CredentialProviders[0].MappingType)
	require.NotNil(t, dto.CredentialProviders[0].AccessKeyId)
	assert.Equal(t, "ACCESSKEY1", *dto.CredentialProviders[0].AccessKeyId)
	assert.Nil(t, dto.CredentialProviders[0].AccountName)
	assert.Nil(t, dto.CredentialProviders[0].HeaderName)
	assert.Nil(t, dto.CredentialProviders[0].HeaderValue)
	assert.Nil(t, dto.CredentialProviders[0].HttpbodyFieldPath)
	assert.Nil(t, dto.CredentialProviders[0].HttpbodyFieldValue)

	// Second mapping (HttpHeader)
	assert.Equal(t, "cp-2", dto.CredentialProviders[1].CredentialProviderId)
	assert.Equal(t, "HttpHeader", dto.CredentialProviders[1].MappingType)
	assert.Nil(t, dto.CredentialProviders[1].AccessKeyId)
	assert.Nil(t, dto.CredentialProviders[1].AccountName)
	require.NotNil(t, dto.CredentialProviders[1].HeaderName)
	assert.Equal(t, "X-Test-Header", *dto.CredentialProviders[1].HeaderName)
	require.NotNil(t, dto.CredentialProviders[1].HeaderValue)
	assert.Equal(t, "test-val", *dto.CredentialProviders[1].HeaderValue)
	assert.Nil(t, dto.CredentialProviders[1].HttpbodyFieldPath)
	assert.Nil(t, dto.CredentialProviders[1].HttpbodyFieldValue)
}

func TestConvertAccessPolicyDTOToModel_CredentialProvidersNullable(t *testing.T) {
	t.Parallel()

	accessKey := "ACCESSKEY1"
	accountName := "my-account"
	dto := aembit.CreatePolicyDTO{
		AccessEntityDTO: aembit.AccessEntityDTO{
			EntityDTO: aembit.EntityDTO{
				ExternalID: "policy-1",
				Name:       "Test Policy",
				IsActive:   true,
			},
		},
		ClientWorkload: "cw-1",
		ServerWorkload: "sw-1",
		CredentialProviders: []aembit.PolicyCredentialMappingDTO{
			{
				CredentialProviderId: "cp-1",
				MappingType:          "AccessKeyId",
				AccessKeyId:          &accessKey,
				AccountName:          nil,
				HeaderName:           nil,
				HeaderValue:          nil,
				HttpbodyFieldPath:    nil,
				HttpbodyFieldValue:   nil,
			},
			{
				CredentialProviderId: "cp-2",
				MappingType:          "AccountName",
				AccessKeyId:          nil,
				AccountName:          &accountName,
				HeaderName:           nil,
				HeaderValue:          nil,
				HttpbodyFieldPath:    nil,
				HttpbodyFieldValue:   nil,
			},
		},
	}

	plan := models.AccessPolicyResourceModel{}
	model := convertAccessPolicyDTOToModel(plan, dto)

	require.Len(t, model.CredentialProviders, 2)

	// First mapping: access_key_id should be string value, others should be StringNull()
	assert.Equal(t, types.StringValue("cp-1"), model.CredentialProviders[0].CredentialProviderId)
	assert.Equal(t, types.StringValue("AccessKeyId"), model.CredentialProviders[0].MappingType)
	assert.Equal(t, types.StringValue("ACCESSKEY1"), model.CredentialProviders[0].AccessKeyId)
	assert.True(t, model.CredentialProviders[0].AccountName.IsNull(), "expected AccountName to be null")
	assert.True(t, model.CredentialProviders[0].HeaderName.IsNull(), "expected HeaderName to be null")
	assert.True(t, model.CredentialProviders[0].HeaderValue.IsNull(), "expected HeaderValue to be null")
	assert.True(t, model.CredentialProviders[0].HttpbodyFieldPath.IsNull(), "expected HttpbodyFieldPath to be null")
	assert.True(t, model.CredentialProviders[0].HttpbodyFieldValue.IsNull(), "expected HttpbodyFieldValue to be null")

	// Second mapping: account_name should be string value, others should be StringNull()
	assert.Equal(t, types.StringValue("cp-2"), model.CredentialProviders[1].CredentialProviderId)
	assert.Equal(t, types.StringValue("AccountName"), model.CredentialProviders[1].MappingType)
	assert.Equal(t, types.StringValue("my-account"), model.CredentialProviders[1].AccountName)
	assert.True(t, model.CredentialProviders[1].AccessKeyId.IsNull(), "expected AccessKeyId to be null")
	assert.True(t, model.CredentialProviders[1].HeaderName.IsNull(), "expected HeaderName to be null")
	assert.True(t, model.CredentialProviders[1].HeaderValue.IsNull(), "expected HeaderValue to be null")
	assert.True(t, model.CredentialProviders[1].HttpbodyFieldPath.IsNull(), "expected HttpbodyFieldPath to be null")
	assert.True(t, model.CredentialProviders[1].HttpbodyFieldValue.IsNull(), "expected HttpbodyFieldValue to be null")
}

func TestConvertAccessPolicyExternalDTOToModel_CredentialProvidersNullable(t *testing.T) {
	t.Parallel()

	dto := aembit.GetPolicyDTO{
		AccessEntityDTO: aembit.AccessEntityDTO{
			EntityDTO: aembit.EntityDTO{
				ExternalID: "policy-1",
				Name:       "Test Policy",
				IsActive:   true,
			},
		},
		ClientWorkload: aembit.EntityMetaDTO{ExternalID: "cw-1"},
		ServerWorkload: aembit.EntityMetaDTO{ExternalID: "sw-1"},
		CredentialProviders: []aembit.EntityMetaDTO{
			{ExternalID: "cp-1"},
			{ExternalID: "cp-2"},
		},
	}

	headerName := "X-Custom"
	headerVal := "custom-val"
	accName := "acc-1"
	mappings := []aembit.PolicyCredentialMappingDTO{
		{
			CredentialProviderId: "cp-1",
			MappingType:          "HttpHeader",
			HeaderName:           &headerName,
			HeaderValue:          &headerVal,
		},
		{
			CredentialProviderId: "cp-2",
			MappingType:          "AccountName",
			AccountName:          &accName,
		},
	}

	model := convertAccessPolicyExternalDTOToModel(dto, mappings)

	require.Len(t, model.CredentialProviders, 2)
	assert.Equal(t, types.StringValue("X-Custom"), model.CredentialProviders[0].HeaderName)
	assert.Equal(t, types.StringValue("custom-val"), model.CredentialProviders[0].HeaderValue)
	assert.True(t, model.CredentialProviders[0].AccessKeyId.IsNull())
	assert.True(t, model.CredentialProviders[0].AccountName.IsNull())
	assert.True(t, model.CredentialProviders[0].HttpbodyFieldPath.IsNull())
	assert.True(t, model.CredentialProviders[0].HttpbodyFieldValue.IsNull())

	assert.Equal(t, types.StringValue("acc-1"), model.CredentialProviders[1].AccountName)
	assert.True(t, model.CredentialProviders[1].AccessKeyId.IsNull())
	assert.True(t, model.CredentialProviders[1].HeaderName.IsNull())
	assert.True(t, model.CredentialProviders[1].HeaderValue.IsNull())
	assert.True(t, model.CredentialProviders[1].HttpbodyFieldPath.IsNull())
	assert.True(t, model.CredentialProviders[1].HttpbodyFieldValue.IsNull())
}

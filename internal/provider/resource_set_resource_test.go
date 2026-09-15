package provider

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// sweepResourceSets deletes any existing resource sets whose names start with any of the provided prefixes.
// This ensures that orphaned resource sets from previously failed or aborted test runs do not cause naming collisions.
func sweepResourceSets(t *testing.T, prefixes ...string) {
	if testClient == nil {
		return
	}
	rss, err := testClient.GetResourceSets(nil)
	if err != nil {
		t.Logf("sweepResourceSets: unable to list resource sets: %v", err)
		return
	}
	ctx := context.Background()
	for _, rs := range rss {
		for _, prefix := range prefixes {
			if strings.HasPrefix(rs.Name, prefix) {
				t.Logf("sweepResourceSets: cleaning up orphaned resource set: %s (id: %s)", rs.Name, rs.ExternalID)
				_, _ = testClient.DeleteResourceSet(ctx, rs.ExternalID, nil)
				break
			}
		}
	}
}

// TestAccResourceSet tests the end-to-end lifecycle of a custom resource set.
func TestAccResourceSet(t *testing.T) {
	t.Parallel()
	createFile, err := os.ReadFile(
		"../../tests/resource_set/TestAccResourceSet.tf",
	)
	if err != nil {
		t.Fatalf("failed to read test config file: %v", err)
	}
	modifiedFile, err := os.ReadFile(
		"../../tests/resource_set/TestAccResourceSet.tfmod",
	)
	if err != nil {
		t.Fatalf("failed to read test config file: %v", err)
	}

	sweepResourceSets(t, "TF Acceptance Custom ResourceSet")

	randID := rand.Intn(10000000)
	rsName := fmt.Sprintf("TF Acceptance Custom ResourceSet %d", randID)
	rsModifiedName := fmt.Sprintf("TF Acceptance Custom ResourceSet %d - Modified", randID)

	createConfig := strings.ReplaceAll(string(createFile), "TF Acceptance Custom ResourceSet", rsName)
	modifiedConfig := strings.ReplaceAll(string(modifiedFile), "TF Acceptance Custom ResourceSet - Modified", rsModifiedName)
	modifiedConfig = strings.ReplaceAll(modifiedConfig, "TF Acceptance Custom ResourceSet", rsName)

	t.Cleanup(func() {
		sweepResourceSets(t, rsName, rsModifiedName)
	})

	resourceName := "aembit_resource_set.crs"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: createConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName,
						"name",
						rsName,
					),
					resource.TestCheckResourceAttr(
						resourceName,
						"description",
						rsName,
					),
					resource.TestCheckResourceAttr(
						resourceName,
						"roles.#",
						"2",
					),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: modifiedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName,
						"name",
						rsModifiedName,
					),
					resource.TestCheckResourceAttr(
						resourceName,
						"description",
						rsModifiedName,
					),
					resource.TestCheckResourceAttr(
						resourceName,
						"roles.#",
						"2",
					),
				),
			},
		},
	})
}

// TestAccResourceSetPolicy tests the end-to-end lifecycle of an access policy within a custom resource set.
func TestAccResourceSetPolicy(t *testing.T) {
	t.Parallel()
	createFile, err := os.ReadFile(
		"../../tests/resource_set/TestAccResourceSetPolicy.tf",
	)
	if err != nil {
		t.Fatalf("failed to read test config file: %v", err)
	}
	modifiedFile, err := os.ReadFile(
		"../../tests/resource_set/TestAccResourceSetPolicy.tfmod",
	)
	if err != nil {
		t.Fatalf("failed to read test config file: %v", err)
	}

	sweepResourceSets(t, "TF Acceptance Custom Policy ResourceSet")

	randID := rand.Intn(10000000)
	rsName := fmt.Sprintf("TF Acceptance Custom Policy ResourceSet %d", randID)
	rsModifiedName := fmt.Sprintf("TF Acceptance Custom Policy ResourceSet %d - Modified", randID)

	createConfig := strings.ReplaceAll(string(createFile), "TF Acceptance Custom Policy ResourceSet", rsName)
	modifiedConfig := strings.ReplaceAll(string(modifiedFile), "TF Acceptance Custom Policy ResourceSet - Modified", rsModifiedName)
	modifiedConfig = strings.ReplaceAll(modifiedConfig, "TF Acceptance Custom Policy ResourceSet", rsName)

	t.Cleanup(func() {
		sweepResourceSets(t, rsName, rsModifiedName)
	})

	resourceName := "aembit_access_policy.first_policy"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: createConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_set_id"),
					resource.TestCheckResourceAttr(resourceName, "name", "TF ResourceSet Policy"),
					resource.TestCheckResourceAttr(resourceName, "is_active", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "client_workload"),
					resource.TestCheckResourceAttr(resourceName, "trust_providers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "access_conditions.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "credential_provider"),
					resource.TestCheckResourceAttrSet(resourceName, "server_workload"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(state *terraform.State) (string, error) {
					rs, ok := state.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("Not found: %s", resourceName)
					}
					return fmt.Sprintf("%s,%s", rs.Primary.Attributes["resource_set_id"], rs.Primary.ID), nil
				},
			},
			{
				Config: modifiedConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "resource_set_id"),
					resource.TestCheckResourceAttr(resourceName, "name", "TF ResourceSet Policy"),
					resource.TestCheckResourceAttr(resourceName, "is_active", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "client_workload"),
					resource.TestCheckResourceAttr(resourceName, "trust_providers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "access_conditions.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "credential_providers.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(
						resourceName,
						"credential_providers.*",
						map[string]string{
							"mapping_type": "HttpHeader",
							"header_name":  "test_header_name",
							"header_value": "test_header_value",
						},
					),
					resource.TestCheckTypeSetElemNestedAttrs(
						resourceName,
						"credential_providers.*",
						map[string]string{
							"mapping_type":         "HttpBody",
							"httpbody_field_path":  "test_field_path",
							"httpbody_field_value": "test_field_value",
						},
					),
					resource.TestCheckResourceAttrSet(resourceName, "server_workload"),
				),
			},
		},
	})
}

// TestAccResourceSetNoRoles tests the creation of a resource set with no assigned roles.
func TestAccResourceSetNoRoles(t *testing.T) {
	t.Parallel()
	config, err := os.ReadFile(
		"../../tests/resource_set/TestAccResourceSetNoRoles.tf",
	)
	if err != nil {
		t.Fatalf("failed to read test config file: %v", err)
	}

	sweepResourceSets(t, "TF Acceptance ResourceSet No Roles")

	randID := rand.Intn(10000000)
	rsName := fmt.Sprintf("TF Acceptance ResourceSet No Roles %d", randID)
	testConfig := strings.ReplaceAll(string(config), "TF Acceptance ResourceSet No Roles", rsName)

	t.Cleanup(func() {
		sweepResourceSets(t, rsName)
	})

	resourceName := "aembit_resource_set.no_roles"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(
						resourceName,
						"name",
						rsName,
					),
					resource.TestCheckResourceAttr(
						resourceName,
						"roles.#",
						"2",
					),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

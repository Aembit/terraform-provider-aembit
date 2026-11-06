package provider

import (
	"context"
	"encoding/json"
	"slices"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// mockPrivateState provides an in-memory implementation of the privateStateData interface for unit testing.
type mockPrivateState struct {
	data map[string][]byte
}

// newMockPrivateState constructs an initialized mockPrivateState instance.
func newMockPrivateState() *mockPrivateState {
	return &mockPrivateState{
		data: make(map[string][]byte),
	}
}

// GetKey retrieves a private state value by key from the in-memory map.
func (m *mockPrivateState) GetKey(_ context.Context, key string) ([]byte, diag.Diagnostics) {
	val, ok := m.data[key]
	if !ok {
		return nil, nil
	}
	return val, nil
}

// SetKey stores a private state value by key in the in-memory map.
func (m *mockPrivateState) SetKey(_ context.Context, key string, value []byte) diag.Diagnostics {
	m.data[key] = value
	return nil
}

// getStringSliceFromSet extracts a slice of strings from a framework types.Set.
func getStringSliceFromSet(ctx context.Context, t *testing.T, set types.Set) []string {
	t.Helper()
	if set.IsNull() || set.IsUnknown() {
		return nil
	}
	var elements []string
	diags := set.ElementsAs(ctx, &elements, false)
	if diags.HasError() {
		t.Fatalf("unexpected error converting set to string slice: %v", diags)
	}
	return elements
}

// TestReconcileRolesOnApply_OmittedRoles verifies reconciliation when planned roles is unknown (omitted).
// It verifies that backend-returned roles are preserved and implicit_roles is set to an empty array.
func TestReconcileRolesOnApply_OmittedRoles(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	mockState := newMockPrivateState()

	planRoles := types.SetUnknown(types.StringType)
	backendRoles := []string{"superadmin-id", "auditor-id"}

	reconciledSet, diags := reconcileRolesOnApply(ctx, planRoles, backendRoles, mockState)
	if diags.HasError() {
		t.Fatalf("reconcileRolesOnApply returned unexpected error diags: %v", diags)
	}

	resultRoles := getStringSliceFromSet(ctx, t, reconciledSet)
	if len(resultRoles) != 2 {
		t.Fatalf("expected 2 reconciled roles, got %d: %v", len(resultRoles), resultRoles)
	}
	for _, expectedRole := range backendRoles {
		if !slices.Contains(resultRoles, expectedRole) {
			t.Errorf("expected reconciled set to contain role %q, got: %v", expectedRole, resultRoles)
		}
	}

	storedBytes, getDiags := mockState.GetKey(ctx, resourceSetImplicitRolesKey)
	if getDiags.HasError() {
		t.Fatalf("mockState.GetKey returned unexpected error diags: %v", getDiags)
	}
	if string(storedBytes) != "[]" {
		t.Errorf("expected implicit_roles in private state to be %q, got %q", "[]", string(storedBytes))
	}
}

// TestReconcileRolesOnApply_NullRoles verifies reconciliation when planned roles is explicit null (roles = null).
// It verifies that null is treated like omitted/unknown: backend-returned roles are preserved and implicit_roles is set to [].
func TestReconcileRolesOnApply_NullRoles(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	mockState := newMockPrivateState()

	planRoles := types.SetNull(types.StringType)
	backendRoles := []string{"superadmin-id", "auditor-id"}

	reconciledSet, diags := reconcileRolesOnApply(ctx, planRoles, backendRoles, mockState)
	if diags.HasError() {
		t.Fatalf("reconcileRolesOnApply returned unexpected error diags: %v", diags)
	}

	resultRoles := getStringSliceFromSet(ctx, t, reconciledSet)
	if len(resultRoles) != 2 {
		t.Fatalf("expected 2 reconciled roles, got %d: %v", len(resultRoles), resultRoles)
	}
	for _, expectedRole := range backendRoles {
		if !slices.Contains(resultRoles, expectedRole) {
			t.Errorf("expected reconciled set to contain role %q, got: %v", expectedRole, resultRoles)
		}
	}

	storedBytes, getDiags := mockState.GetKey(ctx, resourceSetImplicitRolesKey)
	if getDiags.HasError() {
		t.Fatalf("mockState.GetKey returned unexpected error diags: %v", getDiags)
	}
	if string(storedBytes) != "[]" {
		t.Errorf("expected implicit_roles in private state to be %q, got %q", "[]", string(storedBytes))
	}
}

// TestReconcileRolesOnApply_ExplicitEmptyRoles verifies reconciliation when roles is explicitly configured as empty ([]).
// It verifies that reconciled roles is empty (0 elements, non-null) and implicit roles are stored in private state.
func TestReconcileRolesOnApply_ExplicitEmptyRoles(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	mockState := newMockPrivateState()

	emptySet, setDiags := types.SetValue(types.StringType, []attr.Value{})
	if setDiags.HasError() {
		t.Fatalf("types.SetValue returned unexpected error diags: %v", setDiags)
	}

	backendRoles := []string{"superadmin-id", "auditor-id"}

	reconciledSet, diags := reconcileRolesOnApply(ctx, emptySet, backendRoles, mockState)
	if diags.HasError() {
		t.Fatalf("reconcileRolesOnApply returned unexpected error diags: %v", diags)
	}

	if reconciledSet.IsNull() {
		t.Fatal("expected reconciled set to be non-null")
	}
	resultRoles := getStringSliceFromSet(ctx, t, reconciledSet)
	if len(resultRoles) != 0 {
		t.Fatalf("expected 0 reconciled roles, got %d: %v", len(resultRoles), resultRoles)
	}

	storedBytes, getDiags := mockState.GetKey(ctx, resourceSetImplicitRolesKey)
	if getDiags.HasError() {
		t.Fatalf("mockState.GetKey returned unexpected error diags: %v", getDiags)
	}

	var storedImplicitRoles []string
	if err := json.Unmarshal(storedBytes, &storedImplicitRoles); err != nil {
		t.Fatalf("failed to unmarshal stored implicit roles: %v", err)
	}

	if len(storedImplicitRoles) != 2 {
		t.Fatalf("expected 2 stored implicit roles, got %d: %v", len(storedImplicitRoles), storedImplicitRoles)
	}
	for _, expectedRole := range backendRoles {
		if !slices.Contains(storedImplicitRoles, expectedRole) {
			t.Errorf("expected stored implicit roles to contain %q, got: %v", expectedRole, storedImplicitRoles)
		}
	}
}

// TestReconcileRolesOnApply_CustomRoles verifies reconciliation when custom roles are explicitly planned.
// It verifies that reconciled roles contains only the custom role and backend-injected roles are persisted as implicit.
func TestReconcileRolesOnApply_CustomRoles(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	mockState := newMockPrivateState()

	customRole := "custom-1"
	customSet, setDiags := types.SetValueFrom(ctx, types.StringType, []string{customRole})
	if setDiags.HasError() {
		t.Fatalf("types.SetValueFrom returned unexpected error diags: %v", setDiags)
	}

	backendRoles := []string{"superadmin-id", "auditor-id", customRole}

	reconciledSet, diags := reconcileRolesOnApply(ctx, customSet, backendRoles, mockState)
	if diags.HasError() {
		t.Fatalf("reconcileRolesOnApply returned unexpected error diags: %v", diags)
	}

	resultRoles := getStringSliceFromSet(ctx, t, reconciledSet)
	if len(resultRoles) != 1 {
		t.Fatalf("expected 1 reconciled role, got %d: %v", len(resultRoles), resultRoles)
	}
	if resultRoles[0] != customRole {
		t.Errorf("expected reconciled role %q, got %q", customRole, resultRoles[0])
	}

	storedBytes, getDiags := mockState.GetKey(ctx, resourceSetImplicitRolesKey)
	if getDiags.HasError() {
		t.Fatalf("mockState.GetKey returned unexpected error diags: %v", getDiags)
	}

	var storedImplicitRoles []string
	if err := json.Unmarshal(storedBytes, &storedImplicitRoles); err != nil {
		t.Fatalf("failed to unmarshal stored implicit roles: %v", err)
	}

	expectedImplicit := []string{"superadmin-id", "auditor-id"}
	if len(storedImplicitRoles) != len(expectedImplicit) {
		t.Fatalf("expected %d stored implicit roles, got %d: %v", len(expectedImplicit), len(storedImplicitRoles), storedImplicitRoles)
	}
	for _, expectedRole := range expectedImplicit {
		if !slices.Contains(storedImplicitRoles, expectedRole) {
			t.Errorf("expected stored implicit roles to contain %q, got: %v", expectedRole, storedImplicitRoles)
		}
	}
}

// TestReconcileRolesOnApply_PlannedRoleMissingInBackend verifies that if a planned role is missing in the backend response,
// the discrepancy is not masked, implicit_roles in private state is cleared, and full backend roles are returned.
func TestReconcileRolesOnApply_PlannedRoleMissingInBackend(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	mockState := newMockPrivateState()

	// Pre-seed private state with prior implicit roles to verify they are cleared on discrepancy
	mockState.SetKey(ctx, resourceSetImplicitRolesKey, []byte(`["superadmin-id","auditor-id"]`))

	plannedRolesList := []string{"custom-1", "custom-missing"}
	planRoles, setDiags := types.SetValueFrom(ctx, types.StringType, plannedRolesList)
	if setDiags.HasError() {
		t.Fatalf("types.SetValueFrom returned unexpected error diags: %v", setDiags)
	}

	backendRoles := []string{"superadmin-id", "auditor-id", "custom-1"}

	reconciledSet, diags := reconcileRolesOnApply(ctx, planRoles, backendRoles, mockState)
	if diags.HasError() {
		t.Fatalf("reconcileRolesOnApply returned unexpected error diags: %v", diags)
	}

	resultRoles := getStringSliceFromSet(ctx, t, reconciledSet)
	if len(resultRoles) != len(backendRoles) {
		t.Fatalf("expected %d roles returned (full backend roles), got %d: %v", len(backendRoles), len(resultRoles), resultRoles)
	}
	for _, bRole := range backendRoles {
		if !slices.Contains(resultRoles, bRole) {
			t.Errorf("expected reconciled set to contain backend role %q, got: %v", bRole, resultRoles)
		}
	}

	storedBytes, getDiags := mockState.GetKey(ctx, resourceSetImplicitRolesKey)
	if getDiags.HasError() {
		t.Fatalf("mockState.GetKey returned unexpected error diags: %v", getDiags)
	}
	if string(storedBytes) != "[]" {
		t.Errorf("expected implicit_roles in private state to be cleared to %q, got %q", "[]", string(storedBytes))
	}
}

// TestFilterRolesOnRead_WithImplicitRoles verifies that filterRolesOnRead removes implicit system roles from the backend roles.
func TestFilterRolesOnRead_WithImplicitRoles(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	mockState := newMockPrivateState()

	implicitRoles := []string{"superadmin-id", "auditor-id"}
	marshaled, err := json.Marshal(implicitRoles)
	if err != nil {
		t.Fatalf("failed to marshal implicit roles: %v", err)
	}
	mockState.SetKey(ctx, resourceSetImplicitRolesKey, marshaled)

	backendRoles := []string{"superadmin-id", "auditor-id", "custom-1"}

	filteredSet, diags := filterRolesOnRead(ctx, backendRoles, mockState)
	if diags.HasError() {
		t.Fatalf("filterRolesOnRead returned unexpected error diags: %v", diags)
	}

	resultRoles := getStringSliceFromSet(ctx, t, filteredSet)
	if len(resultRoles) != 1 {
		t.Fatalf("expected 1 filtered role, got %d: %v", len(resultRoles), resultRoles)
	}
	if resultRoles[0] != "custom-1" {
		t.Errorf("expected filtered role %q, got %q", "custom-1", resultRoles[0])
	}
}

// TestFilterRolesOnRead_DetectOutOfBandDrift verifies that if a role was added out-of-band in the UI,
// it is retained in the filtered roles set so Terraform can detect state drift.
func TestFilterRolesOnRead_DetectOutOfBandDrift(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	mockState := newMockPrivateState()

	// Resource created with roles = [] (implicit roles are superadmin and auditor)
	implicitRoles := []string{"superadmin-id", "auditor-id"}
	marshaled, err := json.Marshal(implicitRoles)
	if err != nil {
		t.Fatalf("failed to marshal implicit roles: %v", err)
	}
	mockState.SetKey(ctx, resourceSetImplicitRolesKey, marshaled)

	// Backend now returns an out-of-band added role in UI
	backendRoles := []string{"superadmin-id", "auditor-id", "custom-added-in-ui"}

	filteredSet, diags := filterRolesOnRead(ctx, backendRoles, mockState)
	if diags.HasError() {
		t.Fatalf("filterRolesOnRead returned unexpected error diags: %v", diags)
	}

	resultRoles := getStringSliceFromSet(ctx, t, filteredSet)
	if len(resultRoles) != 1 {
		t.Fatalf("expected 1 filtered role (the drift role), got %d: %v", len(resultRoles), resultRoles)
	}
	if resultRoles[0] != "custom-added-in-ui" {
		t.Errorf("expected filtered role %q, got %q", "custom-added-in-ui", resultRoles[0])
	}
}

// TestFilterRolesOnRead_NoImplicitRoles verifies that if private state has no key or an empty array,
// the full backend roles slice is returned unchanged.
func TestFilterRolesOnRead_NoImplicitRoles(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("nil private state", func(t *testing.T) {
		t.Parallel()
		backendRoles := []string{"superadmin-id", "auditor-id"}
		filteredSet, diags := filterRolesOnRead(ctx, backendRoles, nil)
		if diags.HasError() {
			t.Fatalf("filterRolesOnRead returned unexpected error diags: %v", diags)
		}
		resultRoles := getStringSliceFromSet(ctx, t, filteredSet)
		if len(resultRoles) != 2 {
			t.Fatalf("expected 2 roles, got %d: %v", len(resultRoles), resultRoles)
		}
	})

	t.Run("empty private state key", func(t *testing.T) {
		t.Parallel()
		mockState := newMockPrivateState()
		backendRoles := []string{"superadmin-id", "auditor-id"}
		filteredSet, diags := filterRolesOnRead(ctx, backendRoles, mockState)
		if diags.HasError() {
			t.Fatalf("filterRolesOnRead returned unexpected error diags: %v", diags)
		}
		resultRoles := getStringSliceFromSet(ctx, t, filteredSet)
		if len(resultRoles) != 2 {
			t.Fatalf("expected 2 roles, got %d: %v", len(resultRoles), resultRoles)
		}
	})

	t.Run("empty array in private state", func(t *testing.T) {
		t.Parallel()
		mockState := newMockPrivateState()
		mockState.SetKey(ctx, resourceSetImplicitRolesKey, []byte("[]"))
		backendRoles := []string{"superadmin-id", "auditor-id"}
		filteredSet, diags := filterRolesOnRead(ctx, backendRoles, mockState)
		if diags.HasError() {
			t.Fatalf("filterRolesOnRead returned unexpected error diags: %v", diags)
		}
		resultRoles := getStringSliceFromSet(ctx, t, filteredSet)
		if len(resultRoles) != 2 {
			t.Fatalf("expected 2 roles, got %d: %v", len(resultRoles), resultRoles)
		}
	})
}

// TestRoleResource_SchemaMetadata verifies the schema descriptions and attribute metadata
// for aembit_role, ensuring documentation for Default Resource Set assignments, permission
// requirements, and single-sided ownership is maintained alongside unchanged default/computed behavior.
func TestRoleResource_SchemaMetadata(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	res := NewRoleResource()
	var resp resource.SchemaResponse
	res.Schema(ctx, resource.SchemaRequest{}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema returned unexpected error diagnostics: %v", resp.Diagnostics)
	}

	// 1. Top-level Description
	expectedTopDesc := "Resource to create and manage Roles in an Aembit Cloud tenant."
	if resp.Schema.Description != expectedTopDesc {
		t.Errorf("expected top-level Schema.Description %q, got %q", expectedTopDesc, resp.Schema.Description)
	}

	// 2. resource_sets_assignments attribute metadata
	attrVal, ok := resp.Schema.Attributes["resource_sets_assignments"]
	if !ok {
		t.Fatalf("missing required schema attribute resource_sets_assignments")
	}
	setAttr, ok := attrVal.(schema.SetAttribute)
	if !ok {
		t.Fatalf("expected resource_sets_assignments to be schema.SetAttribute, got %T", attrVal)
	}

	if !setAttr.Optional {
		t.Errorf("expected resource_sets_assignments.Optional to be true")
	}
	if !setAttr.Computed {
		t.Errorf("expected resource_sets_assignments.Computed to be true")
	}

	// Description checks
	expectedDescSnippets := []string{
		"ffffffff-ffff-ffff-ffff-ffffffffffff",
		"[]",
		"resource_sets = { read = true, ... }",
		"Manage Role-to-ResourceSet associations on either aembit_role.resource_sets_assignments or aembit_resource_set.roles, not both",
	}
	for _, snippet := range expectedDescSnippets {
		if !strings.Contains(setAttr.Description, snippet) {
			t.Errorf("expected resource_sets_assignments Description to contain %q, but got %q", snippet, setAttr.Description)
		}
	}

	// Default value verification
	if setAttr.Default == nil {
		t.Fatalf("expected setAttr.Default to be configured, got nil")
	}
	var defaultResp defaults.SetResponse
	setAttr.Default.DefaultSet(ctx, defaults.SetRequest{
		Path: path.Root("resource_sets_assignments"),
	}, &defaultResp)
	if defaultResp.Diagnostics.HasError() {
		t.Fatalf("DefaultSet returned unexpected error diagnostics: %v", defaultResp.Diagnostics)
	}

	expectedDefaultSet := types.SetValueMust(types.StringType, []attr.Value{
		types.StringValue("ffffffff-ffff-ffff-ffff-ffffffffffff"),
	})
	if !defaultResp.PlanValue.Equal(expectedDefaultSet) {
		t.Errorf("expected default set value %v, got %v", expectedDefaultSet, defaultResp.PlanValue)
	}
}

// TestResourceSetResource_SchemaMetadata verifies the schema descriptions and attribute metadata
// for aembit_resource_set, ensuring documentation for roles, effective_roles, and single-sided
// ownership is maintained alongside unchanged computed behavior.
func TestResourceSetResource_SchemaMetadata(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	res := NewResourceSetResource()
	var resp resource.SchemaResponse
	res.Schema(ctx, resource.SchemaRequest{}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema returned unexpected error diagnostics: %v", resp.Diagnostics)
	}

	// 1. roles attribute metadata
	rolesAttrVal, ok := resp.Schema.Attributes["roles"]
	if !ok {
		t.Fatalf("missing required schema attribute roles")
	}
	rolesSetAttr, ok := rolesAttrVal.(schema.SetAttribute)
	if !ok {
		t.Fatalf("expected roles to be schema.SetAttribute, got %T", rolesAttrVal)
	}

	if !rolesSetAttr.Optional {
		t.Errorf("expected roles.Optional to be true")
	}
	if !rolesSetAttr.Computed {
		t.Errorf("expected roles.Computed to be true")
	}

	expectedRolesSnippets := []string{
		"SuperAdmin",
		"Auditor",
		"effective_roles",
		"Manage Role-to-ResourceSet associations on either aembit_resource_set.roles or aembit_role.resource_sets_assignments, not both",
	}
	for _, snippet := range expectedRolesSnippets {
		if !strings.Contains(rolesSetAttr.Description, snippet) {
			t.Errorf("expected roles Description to contain %q, but got %q", snippet, rolesSetAttr.Description)
		}
	}

	// 2. effective_roles attribute metadata
	effRolesAttrVal, ok := resp.Schema.Attributes["effective_roles"]
	if !ok {
		t.Fatalf("missing required schema attribute effective_roles")
	}
	effRolesSetAttr, ok := effRolesAttrVal.(schema.SetAttribute)
	if !ok {
		t.Fatalf("expected effective_roles to be schema.SetAttribute, got %T", effRolesAttrVal)
	}

	if !effRolesSetAttr.Computed {
		t.Errorf("expected effective_roles.Computed to be true")
	}
	if effRolesSetAttr.Optional {
		t.Errorf("expected effective_roles.Optional to be false")
	}
}

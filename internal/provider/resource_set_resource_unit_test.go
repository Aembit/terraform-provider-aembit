package provider

import (
	"context"
	"encoding/json"
	"slices"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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

// TestReconcileRolesOnApply_OmittedOrNullRoles verifies reconciliation when planned roles is unknown (omitted/null).
// It verifies that backend-returned roles are preserved and implicit_roles is set to an empty array.
func TestReconcileRolesOnApply_OmittedOrNullRoles(t *testing.T) {
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
// the discrepancy is not masked and full backend roles are returned.
func TestReconcileRolesOnApply_PlannedRoleMissingInBackend(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	mockState := newMockPrivateState()

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

package modifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Custom plan modifier to default a Set to empty
func UseEmptySet() planmodifier.Set {
	return useEmptySetModifier{}
}

type useEmptySetModifier struct{}

func (m useEmptySetModifier) Description(ctx context.Context) string {
	return "Defaults to an empty set if not configured"
}

func (m useEmptySetModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m useEmptySetModifier) PlanModifySet(
	ctx context.Context,
	req planmodifier.SetRequest,
	resp *planmodifier.SetResponse,
) {
	// If value is null and unknown, default to empty set
	if req.PlanValue.IsNull() || req.PlanValue.IsUnknown() {
		resp.PlanValue = types.SetValueMust(req.PlanValue.ElementType(ctx), []attr.Value{})
	}
}

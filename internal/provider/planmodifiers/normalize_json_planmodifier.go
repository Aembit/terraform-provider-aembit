package planmodifiers

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type normalizeJSONStringModifier struct{}

func (m normalizeJSONStringModifier) Description(_ context.Context) string {
	return "Normalizes JSON string (e.g., removes whitespace)"
}

func (m normalizeJSONStringModifier) MarkdownDescription(_ context.Context) string {
	return m.Description(context.Background())
}

func (m normalizeJSONStringModifier) PlanModifyString(
	ctx context.Context,
	req planmodifier.StringRequest,
	resp *planmodifier.StringResponse,
) {
	// Don't change null or unknown values
	if req.PlanValue.IsNull() || req.PlanValue.IsUnknown() {
		return
	}

	original := req.PlanValue.ValueString()
	normalized, err := normalizeJSON(original)
	if err != nil {
		resp.Diagnostics.AddError("Invalid JSON", err.Error())
		return
	}

	resp.PlanValue = types.StringValue(normalized)
}

// The actual normalization function
func normalizeJSON(input string) (string, error) {
	input = strings.TrimSpace(input)
	var buf bytes.Buffer
	err := json.Compact(&buf, []byte(input))
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// Exported for use in schema
func NormalizeJSONString() planmodifier.String {
	return normalizeJSONStringModifier{}
}

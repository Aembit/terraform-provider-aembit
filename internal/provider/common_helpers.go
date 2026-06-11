package provider

import (
	"context"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func newHTTPHeadersModel(ctx context.Context, headers []aembit.KeyValuePair) types.Map {
	respMap := make(map[string]string)

	if len(headers) > 0 {
		tflog.Debug(ctx, "newHTTPHeadersModel: static headers found.")
		for _, headerEntry := range headers {
			respMap[headerEntry.Key] = headerEntry.Value
		}
		headersMap, _ := types.MapValueFrom(ctx, types.StringType, respMap)
		return headersMap
	}

	return types.MapNull(types.StringType)
}

func modifyPlanForResourceSetId(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
	client *aembit.CloudClient,
) {
	// Exit if the resource is being destroyed
	if req.Plan.Raw.IsNull() {
		return
	}

	// Inspect the CONFIG, not the Plan.
	// This tells us exactly what the user wrote in their HCL file.
	var configVal types.String
	diags := req.Config.GetAttribute(ctx, path.Root("resource_set_id"), &configVal)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If the config value is Unknown, the user explicitly populated this field
	// with a reference to another resource that hasn't been created yet.
	// We MUST exit early and let Terraform Core resolve the dependency.
	if configVal.IsUnknown() {
		return
	}

	// If the config value is Null, it means the user completely omitted the
	// attribute from their HCL block. Now it is safe to apply our defaulting logic.
	if configVal.IsNull() {
		// Determine our target fallback ID
		targetResourceSetId := DEFAULT_RESOURCESET_ID
		if client != nil && client.ResourceSetId != "" {
			targetResourceSetId = client.ResourceSetId
		}

		var stateVal types.String
		diags := req.State.GetAttribute(ctx, path.Root("resource_set_id"), &stateVal)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		// If it's already in the state and matches our target provider ID,
		// pass the state back into the plan to declare "No Changes".
		if !stateVal.IsNull() && !stateVal.IsUnknown() {
			if stateVal.ValueString() == targetResourceSetId {
				resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("resource_set_id"), stateVal)...)
				return
			}
		}

		// Otherwise (creation or provider ID changed), fill the plan with the target ID
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("resource_set_id"), types.StringValue(targetResourceSetId))...)
	}

	// If configVal is a known string, the user explicitly provided a hardcoded ID.
	// We do nothing and let Terraform handle standard value matching.
}

func getResourceSetId(resourceSetId types.String, client *aembit.CloudClient) string {
	rsId := resourceSetId.ValueString()
	if resourceSetId.IsNull() || resourceSetId.IsUnknown() || rsId == "" {
		rsId = DEFAULT_RESOURCESET_ID
	}

	if client != nil && client.ResourceSetId != "" && client.ResourceSetId != DEFAULT_RESOURCESET_ID {
		rsId = client.ResourceSetId
	}

	return rsId
}

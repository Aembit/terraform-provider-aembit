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
	// Exit if the plan or state is null (destroy scenario)
	if req.Plan.Raw.IsNull() {
		return
	}

	if client == nil || client.ResourceSetId == "" || client.ResourceSetId == DEFAULT_RESOURCESET_ID {
		return
	}

	var resourceSetId types.String
	diags := req.Plan.GetAttribute(ctx, path.Root("resource_set_id"), &resourceSetId)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !resourceSetId.IsUnknown() && resourceSetId.ValueString() != client.ResourceSetId {
		resp.Diagnostics.Append(resp.Plan.SetAttribute(ctx, path.Root("resource_set_id"), client.ResourceSetId)...)
	}
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

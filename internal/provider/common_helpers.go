package provider

import (
	"context"

	"aembit.io/aembit"
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

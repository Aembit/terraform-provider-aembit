package provider

import (
	"context"
	"os"
	"testing"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func newTagsModel(ctx context.Context, tags []aembit.TagDTO) types.Map {
	respMap := make(map[string]string)

	if len(tags) > 0 {
		tflog.Debug(ctx, "newTagsModel: tags found.")
		for _, tagEntry := range tags {
			respMap[tagEntry.Key] = tagEntry.Value
		}
		tagsMap, _ := types.MapValueFrom(ctx, types.StringType, respMap)
		return tagsMap
	}

	return types.MapNull(types.StringType)
}

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

// skipNotCI can be used to skip tests which can ONLY run on GitHub.
func skipNotCI(t *testing.T) {
	if os.Getenv("CI") == "" {
		t.Skip("Skipping testing in CI environment")
	}
}

func getTenantId() string {
	tenantId := os.Getenv("AEMBIT_TENANT_ID")

	if len(tenantId) == 0 { // get the tenant from clientId
		tenantId = getAembitTenantId(os.Getenv("AEMBIT_CLIENT_ID"))
	}

	return tenantId
}

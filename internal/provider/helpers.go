package provider

import (
	"context"
	"math/big"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func getStringAttr(ctx context.Context, tfObject types.Object, name string) string {
	var value string
	tfValue, _ := tfObject.Attributes()[name].ToTerraformValue(ctx)
	tfValue.As(&value)
	return value
}

func getInt32Attr(ctx context.Context, tfObject types.Object, name string) int32 {
	var value big.Float
	tfValue, _ := tfObject.Attributes()[name].ToTerraformValue(ctx)
	tfValue.As(&value)
	result, _ := value.Int64()
	return int32(result)
}

func getBoolAttr(ctx context.Context, tfObject types.Object, name string) bool {
	var value bool
	tfValue, _ := tfObject.Attributes()[name].ToTerraformValue(ctx)
	tfValue.As(&value)
	return value
}

func getSetObjectAttr(ctx context.Context, tfObject types.Object, name string) []types.Object {
	var objSlice []types.Object
	tfsdk.ValueAs(ctx, tfObject.Attributes()[name], &objSlice)

	objects := make([]types.Object, len(objSlice))
	for i, val := range objSlice {
		objects[i] = val
	}

	return objects
}

package provider

import (
	"context"

	"aembit.io/aembit"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var TagsMapAttribute = func() schema.MapAttribute {
	return schema.MapAttribute{
		Description: "Tags are key-value pairs.",
		ElementType: types.StringType,
		Optional:    true,
	}
}

var TagsAllMapAttribute = func() schema.MapAttribute {
	return schema.MapAttribute{
		Description: "A map of all tags that are associated with the resource, including both user-defined tags and any provider-level default tags that are automatically applied. Changes to provider default tags will be reflected in this attribute after the next apply or refresh.",
		ElementType: types.StringType,
		Computed:    true,
		Optional:    true,
	}
}

func newTagsModel(ctx context.Context, tags []aembit.TagDTO) types.Map {
	respMap := make(map[string]string)

	if len(tags) > 0 {
		for _, tagEntry := range tags {
			respMap[tagEntry.Key] = tagEntry.Value
		}
		tagsMap, _ := types.MapValueFrom(ctx, types.StringType, respMap)
		return tagsMap
	}

	return types.MapNull(types.StringType)
}

func newTagsModelFromPlan(ctx context.Context, tags types.Map) types.Map {
	if !tags.IsNull() {
		planTags := []aembit.TagDTO{}
		tagsMap := make(map[string]string)
		_ = tags.ElementsAs(ctx, &tagsMap, true)

		for key, value := range tagsMap {
			planTags = append(planTags, aembit.TagDTO{
				Key:   key,
				Value: value,
			})
		}

		return newTagsModel(ctx, planTags)
	}

	return types.MapNull(types.StringType)
}

func collectAllTagsDto(
	ctx context.Context,
	defaultTags map[string]string,
	resourceTags types.Map,
) []aembit.TagDTO {
	merged_tags := make(map[string]string)

	for k, v := range defaultTags {
		merged_tags[k] = v
	}
	resourceTagsMap := make(map[string]string)
	_ = resourceTags.ElementsAs(ctx, &resourceTagsMap, true)
	for k, v := range resourceTagsMap {
		merged_tags[k] = v
	}

	dtoTags := []aembit.TagDTO{}
	for k, v := range merged_tags {
		dtoTags = append(dtoTags, aembit.TagDTO{
			Key:   k,
			Value: v,
		})
	}
	return dtoTags
}

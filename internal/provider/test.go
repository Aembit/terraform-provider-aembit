package provider

import (
	"fmt"
	"math/rand"
	"strings"
)

func randomizeFileConfig(config, startValue string) (string, string) {
	randID := rand.Intn(10000000)

	endValue := fmt.Sprintf("%s%d", startValue, randID)
	return strings.ReplaceAll(config, startValue, endValue), endValue
}

func updateFileConfig(config, startValue, endValue string) string {
	return strings.ReplaceAll(config, startValue, endValue)
}

func randomizeFileConfigs(newConfig, modifyConfig, startValue string) (string, string, string) {
	randID := rand.Intn(10000000)

	endValue := fmt.Sprintf("%s%d", startValue, randID)
	return strings.ReplaceAll(newConfig, startValue, endValue), strings.ReplaceAll(modifyConfig, startValue, endValue), endValue
}

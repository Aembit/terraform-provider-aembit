package provider

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

var maxRand = big.NewInt(10000000)

const tagsCount = "tags.%"
const tagsColor = "tags.color"
const tagsDay = "tags.day"

//func randomizeFileConfig(config, startValue string) (string, string) {
//	randID, _ := rand.Int(rand.Reader, maxRand)
//
//	endValue := fmt.Sprintf("%s%d", startValue, randID)
//	return strings.ReplaceAll(config, startValue, endValue), endValue
//}

//func updateFileConfig(config, startValue, endValue string) string {
//	return strings.ReplaceAll(config, startValue, endValue)
//}

func randomizeFileConfigs(newConfig, modifyConfig, startValue string) (string, string, string) {
	randID, _ := rand.Int(rand.Reader, maxRand)

	endValue := fmt.Sprintf("%s%d", startValue, randID)
	return strings.ReplaceAll(newConfig, startValue, endValue), strings.ReplaceAll(modifyConfig, startValue, endValue), endValue
}

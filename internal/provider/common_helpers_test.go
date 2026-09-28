package provider

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// skipNotCI can be used to skip tests which can ONLY run on GitHub.
//
//nolint:unused
func skipNotCI(t *testing.T) {
	if os.Getenv("CI") == "" {
		t.Skip("Skipping testing in non CI environment")
	}
}

func getTerraformVersion() string {
	cmd := exec.Command("terraform", "version")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error executing command: %v\n", err)
		return "v1.6" // return the lowest version if something goes wrong
	}

	terraformVersion := strings.Split(strings.TrimSpace(string(output)), "\n")[0]
	fmt.Println(terraformVersion)

	return terraformVersion
}

func Test_NewStringSetModel(t *testing.T) {
	t.Parallel()

	// empty slice -> types.SetNull(types.StringType)
	res := newStringSetModel(context.Background(), nil)
	assert.True(t, res.IsNull())

	res = newStringSetModel(context.Background(), []string{})
	assert.True(t, res.IsNull())

	// non-empty slice -> types.Set with elements
	res = newStringSetModel(context.Background(), []string{"a", "b"})
	assert.False(t, res.IsNull())
	assert.Equal(t, 2, len(res.Elements()))
}

func Test_StringToTypesString(t *testing.T) {
	t.Parallel()

	// empty string -> types.StringNull()
	res := stringToTypesString("")
	assert.True(t, res.IsNull())

	// non-empty string -> types.StringValue("hello")
	res = stringToTypesString("hello")
	assert.False(t, res.IsNull())
	assert.Equal(t, types.StringValue("hello"), res)
}

func Test_StringPointerToTypesString(t *testing.T) {
	t.Parallel()

	// nil pointer -> types.StringNull()
	res := stringPointerToTypesString(nil)
	assert.True(t, res.IsNull())

	// empty string pointer -> types.StringNull()
	empty := ""
	res = stringPointerToTypesString(&empty)
	assert.True(t, res.IsNull())

	// non-empty string pointer -> types.StringValue("hello")
	val := "hello"
	res = stringPointerToTypesString(&val)
	assert.False(t, res.IsNull())
	assert.Equal(t, types.StringValue("hello"), res)
}

func Test_TypesStringToStringPointer(t *testing.T) {
	t.Parallel()

	// Null string -> nil
	p := typesStringToStringPointer(types.StringNull())
	assert.Nil(t, p)

	// Unknown string -> nil
	p = typesStringToStringPointer(types.StringUnknown())
	assert.Nil(t, p)

	// Empty string -> nil
	p = typesStringToStringPointer(types.StringValue(""))
	assert.Nil(t, p)

	// Non-empty string -> pointer to value
	p = typesStringToStringPointer(types.StringValue("hello"))
	require.NotNil(t, p)
	assert.Equal(t, "hello", *p)
}

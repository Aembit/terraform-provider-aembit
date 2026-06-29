package provider

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// skipNotCI can be used to skip tests which can ONLY run on GitHub.
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

// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"flag"
	"log"

	"terraform-provider-aembit/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// Run "go generate" to format example terraform files and generate the docs for the registry/website

// If you do not have terraform installed, you can remove the formatting command, but its suggested to
// ensure the documentation is formatted properly.
//go:generate terraform fmt -recursive ./examples/

// Run the docs generation tool, check its repository for more information on how it works and how docs
// can be customized.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

// The following values for version and date are based on the default GoReleaser functionality: https://goreleaser.com/cookbooks/using-main.version/

// version is set to the release version by goreleaser. It is "dev" when the provider is built and run locally.
var version string = "dev"

// releaseTime is set to the time of the release, in ISO 8601 format. It is "unknown" when the provider is built and run locally.
var date string = "unknown"

func main() {
	var debug bool

	flag.BoolVar(
		&debug,
		"debug",
		false,
		"set to true to run the provider with support for debuggers like delve",
	)
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/aembit/aembit",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), provider.New(version, date), opts)
	if err != nil {
		log.Fatal(err.Error())
	}
}

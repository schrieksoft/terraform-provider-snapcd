// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	snapcd "terraform-provider-snapcd/client"
	"terraform-provider-snapcd/internal/provider"
)

var (
	// these will be set by the goreleaser configuration
	// to appropriate values for the compiled binary.
	version string = "dev"

	// goreleaser can pass other information to the main package, such as the specific commit
	// https://goreleaser.com/cookbooks/using-main.version/
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/schrieksoft/snapcd",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), provider.New(version), opts)

	// go-plugin kills the process about two seconds after asking it to stop.
	snapcd.ReleaseModuleHolds(1500 * time.Millisecond)

	if err != nil {
		log.Fatal(err.Error())
	}
}

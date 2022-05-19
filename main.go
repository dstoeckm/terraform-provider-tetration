package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/tetration-exchange/terraform-provider/tetration"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: tetration.Provider,
	})
}

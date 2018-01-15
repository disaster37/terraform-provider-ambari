package main

import (
	"github.com/disaster37/terraform-provider-ambari/ambari"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: ambari.Provider,
	})
}

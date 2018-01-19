package main

import (
	"github.com/disaster37/terraform-provider-ambari/ambari"
	"github.com/hashicorp/terraform/plugin"
	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

func main() {

	formatter := new(prefixed.TextFormatter)
	formatter.FullTimestamp = true
	formatter.ForceFormatting = true
	log.SetFormatter(formatter)
	log.SetLevel(log.DebugLevel)

	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: ambari.Provider,
	})
}

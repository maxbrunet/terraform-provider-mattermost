package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/maxbrunet/terraform-provider-mattermost/mattermost"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: mattermost.Provider})
}

package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/andrewn3wman7/terraform-provider-statuscake/statuscake"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: statuscake.Provider})
}

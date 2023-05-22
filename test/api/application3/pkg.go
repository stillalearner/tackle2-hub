package application3

import (
	"github.com/konveyor/tackle2-hub/binding"
	"github.com/konveyor/tackle2-hub/test/api/client"
)

var (
	Client *binding.Client
	RichClient *binding.RichClient
	Application binding.Application
)


func init() {
	// Prepare RichClient and login to Hub API (configured from env variables).
	RichClient = client.PrepareRichClient()

	// Access REST client directly (some test API call need it)
	Client = RichClient.Client()

	// Shortcut for Application-related RichClient methods.
	Application = RichClient.Application
}

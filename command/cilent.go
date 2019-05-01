package command

import (
	"fmt"
	"github.com/hashicorp/vault/api"
)

// NewClient return a new vault client and an error
func NewClient(config *api.Config) (*api.Client, error) {
	if config == nil {
		config = api.DefaultConfig()
	}
	if err := config.ReadEnvironment(); err != nil {
		return nil, fmt.Errorf("%s failed to read environment", err)
	}
	client, err := api.NewClient(config)
	return client, err
}

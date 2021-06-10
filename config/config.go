package config

import "github.com/hashicorp/go-hclog"

// Configurations wraps all the config variables required by the auth service
type Configurations struct {
}

func NewConfigurationsFromFile(configFilePath string, logger hclog.Logger) *Configurations {
	var confs Configurations
	return &confs
}

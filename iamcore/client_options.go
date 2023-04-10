package iamcore

import (
	"fmt"
	"log"
)

type Options struct {

	// IamcoreURL to access the Iamcore; "https://cloud.iamcore.io/" by default.
	IamcoreURL string
}

const (
	iamcoreURLKey     = "iamcore.url"
	iamcoreURLEnvKey  = "IAMCORE_URL"
	iamcoreDefaultURL = "https://cloud.iamcore.io/"
)

// ConfigProvider is the configuration data provider interface.
// The interface is explicitly defined to avoid dependency on go-config.
type ConfigProvider interface {
	GetString(key string) string

	EnvBinder
	DefaultSetter
}

// NewOptions returns Options based on the data provided by ConfigProvider.
// It is the responsibility of the client code to preload the config provider with data, but the function will bind env
// variables and set defaults.
// The function will return nil in case of errors.
func NewOptions(cfg ConfigProvider) *Options {
	bindEnvVariables(cfg)
	setDefaults(cfg)

	return &Options{
		IamcoreURL: cfg.GetString(iamcoreURLKey),
	}
}

// Validate validates necessary Authentication client properties and return error when validation failed.
func (o *Options) Validate() error {
	switch {
	case o == nil:
		return fmt.Errorf("nil options")
	case o.IamcoreURL == "":
		return fmt.Errorf("iamcoreURL is a required")
	}

	return nil
}

type EnvBinder interface{ BindEnv(input ...string) error }

type DefaultSetter interface {
	SetDefault(key string, value interface{})
}

func bindEnvVariables(cfg EnvBinder) {
	if err := cfg.BindEnv(iamcoreURLKey, iamcoreURLEnvKey); err != nil {
		log.Printf("Error binding %q environment variable: %v", iamcoreURLKey, err)
	}
}

func setDefaults(cfg DefaultSetter) {
	cfg.SetDefault(iamcoreURLKey, iamcoreDefaultURL)
}

package config

import (
	"fmt"
	"net/url"

	"github.com/l-vitaly/goenv"
)

var urlNil = url.URL{}

const errPattern = "could not set %s"

// env name constants
const (
	BindAddrEnvName = "IMG_PROC_BIND_ADDR"
	SavePathEnvName = "IMG_PROC_SAVE_PATH"
)

// Config config.
type Config struct {
	BindAddr string
	SavePath string
}

// Parse env config vars.
func Parse() (*Config, error) {
	cfg := &Config{}

	goenv.StringVar(&cfg.BindAddr, BindAddrEnvName, "0.0.0.0:9000")
	goenv.StringVar(&cfg.SavePath, SavePathEnvName, "./images")

	goenv.Parse()

	if cfg.BindAddr == "" {
		return nil, fmt.Errorf(errPattern, BindAddrEnvName)
	}
	if cfg.SavePath == "" {
		return nil, fmt.Errorf(errPattern, SavePathEnvName)
	}

	return cfg, nil
}

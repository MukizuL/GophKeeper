package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

type config struct {
	Addr string
	TLS  bool
	Cert string
}

func newConfig() (*config, error) {
	cfg := config{}
	flag.StringVar(&cfg.Addr, "a", "localhost:8080", "Server address")

	flag.BoolVar(&cfg.TLS, "tls", false, "Enable TLS")

	flag.StringVar(&cfg.Cert, "cert", "", "Sets certificate file path.")

	flag.Parse()

	return &cfg, cfg.check()
}

func (cfg *config) check() error {
	if cfg.TLS {
		if cfg.Cert == "" {
			return fmt.Errorf("cert file is required")
		}

		if _, err := os.Stat(cfg.Cert); errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("cert file %s does not exist", cfg.Cert)
		}
	}
	return nil
}

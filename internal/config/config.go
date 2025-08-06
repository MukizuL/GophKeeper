package config

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/MukizuL/GophKeeper/internal/errs"
	"github.com/caarlos0/env/v11"
	"go.uber.org/fx"
)

// Config holds all application configuration.
type Config struct {
	GRPCPort string `env:"GRPC_PORT" json:"grpc_port"`

	TLS  bool   `env:"ENABLE_TLS" json:"enable_TLS"`
	Cert string `env:"CERT_PATH" json:"cert_path"`
	PK   string `env:"PK_PATH" json:"pk_path"`

	Config         string `env:"CONFIG" json:"config"`
	Filepath       string `env:"FILE_STORAGE_PATH" json:"file_storage_path"`
	DSN            string `env:"DATABASE_DSN" json:"database_dsn"`
	MasterPassword string `env:"MASTER_PASSWORD" json:"master_password"`

	Debug bool `env:"DEBUG" json:"debug"`
}

// newConfig fetches parameters, firstly from env variables, secondly from flags, then from file.
func newConfig() (*Config, error) {
	resultCfg := &Config{}

	envCfg, err := envConfig()
	if err != nil {
		flag.Usage()
		return nil, fmt.Errorf("error loading config from env: %w", err)
	}

	flagCfg, err := flagConfig()
	if err != nil {
		flag.Usage()
		return nil, fmt.Errorf("error loading config from flag: %w", err)
	}

	if envCfg.Config != "" || flagCfg.Config != "" {
		fileCfg, err := fileConfig("")
		if err != nil {
			return nil, fmt.Errorf("error loading config from file: %w", err)
		}
		mergeConfig(resultCfg, fileCfg)
	}

	mergeConfig(resultCfg, flagCfg)

	mergeConfig(resultCfg, envCfg)

	err = checkParams(resultCfg)
	if err != nil {
		flag.Usage()
		return nil, err
	}

	return resultCfg, nil
}

func checkParams(cfg *Config) error {
	if cfg.GRPCPort != "" {
		port, err := strconv.Atoi(strings.TrimPrefix(cfg.GRPCPort, ":"))
		if err != nil || port < 0 || port > 65535 {
			return fmt.Errorf("grpc port must be between 0 and 65535")
		}
	} else {
		return fmt.Errorf("grpc port must be set")
	}

	if cfg.TLS {
		if cfg.Cert == "" {
			return errs.ErrNoCert
		}

		if cfg.PK == "" {
			return errs.ErrNoPK
		}

		err := checkFiles(cfg.Cert, cfg.PK)
		if err != nil {
			return err
		}
	}

	if cfg.Filepath != "" {
		// cfg.Filepath will be an absolute path
		if !filepath.IsAbs(cfg.Filepath) {
			temp, err := filepath.Abs(cfg.Filepath)
			if err != nil {
				return fmt.Errorf("error getting absolute path for file storage: %w", err)
			}

			cfg.Filepath = temp
		}
	} else {
		return fmt.Errorf("storage filepath must be set")
	}

	if cfg.DSN == "" {
		return fmt.Errorf("dsn must be set")
	}

	if cfg.MasterPassword == "" {
		return fmt.Errorf("master password must be set")
	}

	//setSwagger(cfg)

	return nil
}

// checkFiles checks whether certificate and private key files exist
func checkFiles(cert, pk string) error {
	if _, err := os.Stat(cert); errors.Is(err, os.ErrNotExist) {
		return errs.ErrNoCert
	}

	if _, err := os.Stat(pk); errors.Is(err, os.ErrNotExist) {
		return errs.ErrNoPK
	}

	return nil
}

//func setSwagger(cfg *Config) {
//	docs.SwaggerInfo.Host = cfg.Addr
//	docs.SwaggerInfo.BasePath = cfg.Base
//}

// envConfig populates Config from environment
func envConfig() (*Config, error) {
	cfg := &Config{}

	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// flagConfig populates Config from program arguments
func flagConfig() (*Config, error) {
	cfg := &Config{}

	flag.StringVar(&cfg.GRPCPort, "p", "", "Sets GRPC server port (e.g.: :8080).")

	flag.StringVar(&cfg.Filepath, "s", "./storage", "Sets server storage folder path.")

	flag.StringVar(&cfg.Config, "c", "", "Sets server config file name.")

	flag.StringVar(&cfg.DSN, "d", "", "Sets server DSN.")

	flag.BoolVar(&cfg.TLS, "tls", false, "Turns on HTTPS. Requires cert and pk to be set.")

	flag.StringVar(&cfg.Cert, "tls-cert", "", "Sets certificate file path.")

	flag.StringVar(&cfg.PK, "tls-pk", "", "Sets private key file path.")

	flag.BoolVar(&cfg.Debug, "debug", false, "Sets server debug mode.")

	flag.Parse()

	return cfg, nil
}

// fileConfig populates Config from a file
func fileConfig(name string) (*Config, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg Config
	err = json.NewDecoder(file).Decode(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

// mergeConfig merges source config to destination config
func mergeConfig(dst, src *Config) {
	if src == nil {
		return
	}

	if src.Config != "" {
		dst.Config = src.Config
	}
	if src.Filepath != "" {
		dst.Filepath = src.Filepath
	}
	if src.DSN != "" {
		dst.DSN = src.DSN
	}
	if src.MasterPassword != "" {
		dst.MasterPassword = src.MasterPassword
	}
	if src.Cert != "" {
		dst.Cert = src.Cert
	}
	if src.PK != "" {
		dst.PK = src.PK
	}
	if src.GRPCPort != "" {
		dst.GRPCPort = src.GRPCPort
	}
	// Booleans: only overwrite if true to preserve priority
	if src.TLS {
		dst.TLS = true
	}
	if src.Debug {
		dst.Debug = true
	}
}

func Provide() fx.Option {
	return fx.Provide(newConfig)
}

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"

	"ban7er.xscotophilic.art/internal/jsonlog"
	vault "github.com/hashicorp/vault/api"
)

var (
	buildTime string
	version   string
)

type config struct {
	port int
	env  string
	cors struct {
		trustedOrigins []string
	}
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
	api struct {
		key    string
		secret string
	}
	vault struct {
		address string
		token   string
	}
}

type application struct {
	config      config
	logger      *jsonlog.Logger
	wg          sync.WaitGroup
	vaultClient *vault.Client
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	flag.Func(
		"cors-trusted-origins",
		"Trusted CORS origins (space separated)",
		func(val string) error {
			cfg.cors.trustedOrigins = strings.Fields(val)
			return nil
		},
	)

	flag.StringVar(&cfg.api.key, "api-key", "", "Api Key")
	flag.StringVar(&cfg.api.secret, "api-secret", "", "Api Secret")

	flag.StringVar(&cfg.vault.address, "vault-address", "", "Vault Address")
	flag.StringVar(&cfg.vault.token, "vault-token", "", "Vault Token")

	displayVersion := flag.Bool("version", false, "Display version and exit")

	flag.Parse()

	if *displayVersion {
		fmt.Printf("Version:\t%s\n", version)
		fmt.Printf("Build time:\t%s\n", buildTime)
		os.Exit(0)
	}

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	vaultClient, err := (&VaultModel{
		address:            cfg.vault.address,
		token:              cfg.vault.token,
		insecureSkipVerify: cfg.env == "development",
	}).NewClient()

	if err != nil {
		logger.PrintError(
			err,
			map[string]string{
				"message": "Error creating vault client",
			},
		)
		os.Exit(0)
	}

	app := &application{
		config:      cfg,
		logger:      logger,
		vaultClient: vaultClient,
	}

	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}

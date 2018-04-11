package main

import (
	"io"
	"os"
)

//GollumClient holds configuration for the app
type GollumClient struct {
	vaultGatekeeperURL string
	vaultURL           string
	vaultKeyPath       string
	consulURL          string
	consulKeyPath      string
	appVersion         bool
}

//GollumProvider for passing funcs to interface
type GollumProvider struct {
}

//Gollum interface
type Gollum interface {
	validateRequiredCLIArgs(gollum GollumClient) bool
	importFlags() (GollumClient, error)
}

// This is used so we can test for stdout
var out io.Writer = os.Stdout

// Version should be set at build with -X main.Version 0.0.0
var Version = "No version provided"

func main() {
	config := parseCLIArgs(GollumProvider{})

	if config.consulURL != "" {
		outputKVs(ConsulProvider{Config: config})
	}

	if config.vaultGatekeeperURL != "" {
		outputSecrets(SecretsProvider{Config: config})
	}
}

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
}

//GollumProvider for passing funcs to interface
type GollumProvider struct {
}

//Gollum interface
type Gollum interface {
	validateRequiredCLIArgs(gollum GollumClient) error
	importFlags() (GollumClient, error)
}

// This is used so we can test for stdout
var out io.Writer = os.Stdout

func main() {
	config := parseCLIArgs(GollumProvider{})

	if config.consulURL != "" {
		outputKVs(ConsulProvider{Config: config})
	}

	if config.vaultGatekeeperURL != "" {
		outputSecrets(SecretsProvider{Config: config})
	}
}

package main

import (
	//"errors"
	"flag"
	"fmt"
	"os"
)

func parseCLIArgs(g Gollum) GollumClient {

	client, err := g.importFlags()
	if err != nil {
		fmt.Println("something errored")
	}

	return client
}

func (g GollumProvider) importFlags() (GollumClient, error) {
	gatekeeperurl := flag.String("gatekeeperurl", "", "URL for vault gatekeeper")
	vaulturl := flag.String("vaulturl", "", "URL for vault")
	vaultkeypath := flag.String("vaultkeypath", "", "Vault path to key with secrets")
	consulurl := flag.String("consulurl", "", "URL for consul cluster")
	keypath := flag.String("keypath", "", "Path to consul key values for configuration")
	version := flag.Bool("version", false, "prints current app version")
	flag.Parse()

	vaultGatekeeperURL := *gatekeeperurl
	vaultURL := *vaulturl
	vaultKeyPath := *vaultkeypath
	consulURL := *consulurl
	consulKeyPath := *keypath
	appVersion := *version

	client := GollumClient{
		vaultGatekeeperURL: vaultGatekeeperURL,
		vaultURL:           vaultURL,
		vaultKeyPath:       vaultKeyPath,
		consulURL:          consulURL,
		consulKeyPath:      consulKeyPath,
		appVersion:         appVersion,
	}

	validated := g.validateRequiredCLIArgs(client)
	if validated != true {
		flag.PrintDefaults()
		os.Exit(1)
	}

	return client, nil
}

func (g GollumProvider) validateRequiredCLIArgs(flags GollumClient) bool {
	//make sure we have minimum amount of CLI arguments
	if len(os.Args) < 2 {
		return false
	}

	if flags.appVersion {
		fmt.Println(Version)
		return false
	}

	if flags.vaultGatekeeperURL != "" && flags.vaultURL == "" {
		return false

	}

	if flags.vaultGatekeeperURL == "" && flags.vaultURL != "" {
		return false
	}

	if flags.vaultGatekeeperURL != "" && flags.vaultURL != "" && flags.vaultKeyPath == "" {
		return false
	}

	if flags.consulURL != "" && flags.consulKeyPath == "" {
		return false
	}

	if flags.consulURL == "" && flags.consulKeyPath != "" {
		return false
	}

	return true
}

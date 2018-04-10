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

	err := g.validateRequiredCLIArgs(client)
	if err != nil {
		fmt.Println("some error with validateRequiredCLIargs")
	}

	return client, nil
}

func (g GollumProvider) validateRequiredCLIArgs(flags GollumClient) error {
	//make sure we have minimum amount of CLI arguments
	if len(os.Args) < 2 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if flags.appVersion {
		fmt.Println(Version)
		os.Exit(0)
	}

	if flags.vaultGatekeeperURL != "" && flags.vaultURL == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if flags.vaultGatekeeperURL == "" && flags.vaultURL != "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if flags.consulURL != "" && flags.consulKeyPath == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if flags.consulURL == "" && flags.consulKeyPath != "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	return nil
}

package main

import (
	//"flag"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestParseCLIArgs(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"cmd", "-consulurl=https://consulurl", "-keypath=/some/path/to/keys", "-gatekeeperurl=https://gatekeeperurl", "-vaulturl=https://vaulturl:8200", "-vaultkeypath=secrets/someapp/someenv/secrets"}

	g := &GollumProvider{}

	client := parseCLIArgs(g)

	assert.Equal(t, client.consulURL, "https://consulurl", "consulurl should be https://consulurl")
	assert.Equal(t, client.consulKeyPath, "/some/path/to/keys", "keypath should be /some/path/to/keys")
	assert.Equal(t, client.vaultGatekeeperURL, "https://gatekeeperurl", "gatekeeper url should be https://gatekeeperurl")
	assert.Equal(t, client.vaultURL, "https://vaulturl:8200", "vault url should be https://vaulturl:8200")
	assert.Equal(t, client.vaultKeyPath, "secrets/someapp/someenv/secrets", "vault keypath should be secrets/someapp/someenv/secrets")

}

// need to fix this

func TestFlagValidation(t *testing.T) {
	/*
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()
		os.Args = []string{"cmd", "-consulurl=https://consulurl", "-keypath=/some/path/to/keys", "-gatekeeperurl=https://gatekeeperurl", "-vaulturl=https://vaulturl:8200", "-vaultkeypath=secrets/someapp/someenv/secrets"}
	*/
	client := GollumClient{
		vaultGatekeeperURL: "",
		vaultURL:           "",
		vaultKeyPath:       "",
		consulURL:          "",
		consulKeyPath:      "",
		appVersion:         false,
	}

	g := &GollumProvider{}
	client.vaultURL = "https://vault:8200"

	assert.False(t, g.validateRequiredCLIArgs(client), "It should invalidate with only vaultURL")

	client.vaultURL = ""
	client.vaultGatekeeperURL = "https://vault-gatekeeper"
	assert.False(t, g.validateRequiredCLIArgs(client), "It should invalidate with only vaultGatekeeperURL")

	client.vaultURL = "https://vault:8200"
	assert.False(t, g.validateRequiredCLIArgs(client), "It should invalidate with only vaultGatekeeperURL and vaultURL")

	client.vaultKeyPath = "securet/appname/productin/secrets"
	client.consulURL = "https://consul"
	assert.False(t, g.validateRequiredCLIArgs(client), "It should invalidate with with only consulURL set")

	client.consulURL = ""
	client.consulKeyPath = "/some/key"
	assert.False(t, g.validateRequiredCLIArgs(client), "It should invalidate with with only consulKeyPath set")

	client.appVersion = true
	assert.False(t, g.validateRequiredCLIArgs(client), "It should invalidate when appVersion is true")

	//set < 2 os args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"cmd"}
	assert.False(t, g.validateRequiredCLIArgs(client), "It should invalidate when command line arguments are less than 2")
}

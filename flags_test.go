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
/*
func TestFlagValidation(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"cmd", "-consul=true", "-consulurl=https://consulurl", "-keypath=/some/path/to/keys", "-gatekeeper=true", "-gatekeeperurl=https://gatekeeperurl"}
	client := GollumClient{
		vaultGatekeeperEnabled: false,
		vaultGatekeeperURL:     "",
		consulEnabled:          true,
		consulURL:              "",
		consulKeyPath:          "/some/path/to/keys",
	}

	g := &GollumProvider{}
	err := g.validateRequiredCLIArgs(client)
	assert.NotNil(t, err)
}
*/

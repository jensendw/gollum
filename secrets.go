package main

import (
	"crypto/x509"
	"errors"
	"fmt"
	"github.com/channelmeter/vault-gatekeeper-mesos/gatekeeper"
	vaultapi "github.com/hashicorp/vault/api"
	"os"
)

//SecretsProvider for passing funcs to interface
type SecretsProvider struct {
	Config GollumClient
}

//Secrets interface
type Secrets interface {
	getVaultTokenFromGatekeeper() (*string, error)
	getSecrets() (map[string]interface{}, error)
	getTaskID() (*string, error)
}

func (s SecretsProvider) getTaskID() (*string, error) {
	taskID := os.Getenv("MESOS_TASK_ID")
	if taskID == "" {
		return nil, errors.New("Mesos task ID not found")
	}
	return &taskID, nil
}

func (s SecretsProvider) getVaultTokenFromGatekeeper() (*string, error) {
	taskID, err := s.getTaskID()
	if err != nil {
		return nil, err
	}
	//In order to not use env variables we need this in order to initialize a NewClient for vault gatekeeper lib
	var rootCas *x509.CertPool
	client, err := gatekeeper.NewClient(s.Config.vaultURL, s.Config.vaultGatekeeperURL, rootCas)
	if err != nil {
		return nil, err
	}
	token, err := client.RequestVaultToken(*taskID)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (s SecretsProvider) getSecrets() (map[string]interface{}, error) {
	client, err := vaultapi.NewClient(nil)
	if err != nil {
		return nil, err
	}
	vaultToken, err := s.getVaultTokenFromGatekeeper()
	if err != nil {
		return nil, err
	}
	client.SetToken(*vaultToken)

	keys, err := client.Logical().Read(s.Config.vaultKeyPath)
	if err != nil {
		return nil, err
	}

	return keys.Data, nil
}

func outputSecrets(s Secrets) error {
	secrets, err := s.getSecrets()
	if err != nil {
		fmt.Fprint(out, err)
	}

	for key, value := range secrets {
		fmt.Fprintf(out, "export %v=%v\n", key, value)
	}
	return nil
}

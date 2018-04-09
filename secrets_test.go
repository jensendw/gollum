package main

import (
	"bytes"
	"errors"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

type FakeSecretsProvider struct {
	Error bool
}

func (f FakeSecretsProvider) getSecrets() (map[string]interface{}, error) {
	secretsMap := map[string]interface{}{
		"KEY1": "FAKEVALUE1",
		"KEY2": "FAKEVALUE2",
	}
	if f.Error == true {
		return nil, errors.New("Unable to retrieve vault secrets")
	}
	return secretsMap, nil
}

func (f FakeSecretsProvider) getVaultTokenFromGatekeeper() (*string, error) {
	someToken := "abc123TOKEN123abc"
	return &someToken, nil
}

func (f FakeSecretsProvider) getTaskID() (*string, error) {
	taskID := "myFakeTaskID"
	return &taskID, nil
}

func TestOutputSecrets(t *testing.T) {
	bak := out
	out = new(bytes.Buffer)
	defer func() { out = bak }()
	fakeSecrets := FakeSecretsProvider{}
	outputSecrets(fakeSecrets)
	assert.Equal(t, "export KEY1=FAKEVALUE1\nexport KEY2=FAKEVALUE2\n", out.(*bytes.Buffer).String(), "outputSeccrets should return export for vault keys and values")

}

func TestOutputSecretsError(t *testing.T) {
	bak := out
	out = new(bytes.Buffer)
	defer func() { out = bak }()
	fakeSecrets := FakeSecretsProvider{Error: true}
	outputSecrets(fakeSecrets)
	assert.Equal(t, "Unable to retrieve vault secrets", out.(*bytes.Buffer).String(), "If getSecrets produces errors it should be emitted")
}

func TestGetTaskID(t *testing.T) {
	provider := SecretsProvider{}
	// Should error if there is no env variable set
	_, err := provider.getTaskID()
	assert.NotNil(t, err)

	//set the env variable and make sure we get the right value back
	os.Setenv("MESOS_TASK_ID", "dantest.46c3aab0-3c3d-11e8-9418-02426e7d70f7")
	taskID, _ := provider.getTaskID()
	assert.Equal(t, "dantest.46c3aab0-3c3d-11e8-9418-02426e7d70f7", *taskID, "It should retrieve the task ID from MESOS_TASK_ID env variable")
}

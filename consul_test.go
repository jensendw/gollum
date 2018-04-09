package main

import (
	"bytes"
	"errors"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

type FakeConsulProvider struct {
	Error bool
}

func (f FakeConsulProvider) getKVs() ([]map[string]string, error) {
	keysMap := []map[string]string{
		{
			"key":   "fakekey1",
			"value": "fakevalue1",
		},
		{
			"key":   "fakekey2",
			"value": "fakevalue2",
		},
	}
	if f.Error == true {
		return nil, errors.New("error when getting key values")
	}
	return keysMap, nil
}

func (f FakeConsulProvider) splitKeyName(key string) string {
	splitKey := strings.Split(key, "/")
	return splitKey[len(splitKey)-1] //get last element of slice
}

func TestOutputKVs(t *testing.T) {
	bak := out
	out = new(bytes.Buffer)
	defer func() { out = bak }()
	fakeConsul := FakeConsulProvider{}
	outputKVs(fakeConsul)
	assert.Equal(t, "export fakekey1=fakevalue1\nexport fakekey2=fakevalue2\n", out.(*bytes.Buffer).String(), "outputKVs should return map of consul keys and values")
}

func TestOutputKVsError(t *testing.T) {
	bak := out
	out = new(bytes.Buffer)
	defer func() { out = bak }()
	fakeConsulError := FakeConsulProvider{Error: true}
	outputKVs(fakeConsulError)
	assert.Equal(t, "error when getting key values", out.(*bytes.Buffer).String(), "outputKVs should return map of consul keys and values")
}

func TestSplitKeyName(t *testing.T) {
	provider := ConsulProvider{}
	assert.Equal(t, "key", provider.splitKeyName("/my_awesome/_folder_/key"), "outputKVs should return map of consul keys and values")
	assert.Equal(t, "key", provider.splitKeyName("key"), "should return proper name for key without path")
	assert.Equal(t, "key", provider.splitKeyName("/s0me/str@nge/path$/key"), "should return proper key name for strange characters")
}

package main

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"strings"
)

//ConsulProvider for passing funcs to interface
type ConsulProvider struct {
	Config GollumClient
}

//Consul interface
type Consul interface {
	getKVs() ([]map[string]string, error)
	splitKeyName(string) string
}

//http://techblog.zeomega.com/devops/golang/2015/06/09/consul-kv-api-in-golang.html

func (c ConsulProvider) getKVs() ([]map[string]string, error) {
	config := consulapi.DefaultConfig()
	config.Address = c.Config.consulURL
	config.Scheme = "https"
	consul, _ := consulapi.NewClient(config)
	kv := consul.KV()

	kvp, _, err := kv.List(c.Config.consulKeyPath, nil)
	if err != nil {
		return nil, err
	}

	keyValues := []map[string]string{}

	for _, k := range kvp {
		keysMap := map[string]string{
			"key":   string(k.Key),
			"value": string(k.Value),
		}
		keyValues = append(keyValues, keysMap)
	}
	return keyValues, nil

}

func (c ConsulProvider) splitKeyName(key string) string {
	splitKey := strings.Split(key, "/")
	return splitKey[len(splitKey)-1] //get last element of slice
}

func outputKVs(c Consul) {
	kvs, err := c.getKVs()
	if err != nil {
		fmt.Fprint(out, err)
	}
	for _, k := range kvs {
		fmt.Fprintf(out, "export %v=%v\n", c.splitKeyName(k["key"]), k["value"])
	}
}

package config

import (
	"io/ioutil"
)

var cfg *Config

func init() {
	file, err := ioutil.ReadFile("./servercfg.yml")
    if err != nil {
		panic(err)
    }
    yamlString := string(file)
	cfg, err = ParseYaml(yamlString)
    if err != nil {
		panic(err)
    }
}

func GetInstancePtr() *Config{
	return cfg
}


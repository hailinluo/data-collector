package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	LogType string                       `yaml:"log-type"`
	DbUri   string                       `yaml:"db-uri"`
	Tasks   map[string]map[string]string `yaml:"tasks"`
}

var Server = &Config{
	LogType: "stdout",
	DbUri:   "root:root@tcp(127.0.0.1:3306)/data-fundcompany?charset=utf8&clientFoundRows=true",
	Tasks: map[string]map[string]string{
		"fundcompany": {
			"spec":         "0 * * * * *",
			"resource-url": "http://fund.eastmoney.com/company",
		},
	},
}

func InitConfig(path string) error {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return err
	}

	Server = &config
	return nil
}

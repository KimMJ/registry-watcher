package models

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	Webhook []Webhook `yaml:"webhook"`
}

func (c *Config) ReadConfig(filepath string) error {
	fmt.Println(os.Getwd())
	yamlFile, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Printf("yamlFile.Get err #%v", err)
		return err
	}
	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
		return err
	}
	return err
}

package models

type DockerRegistry struct {
	Endpoint         string   `yaml:"endpoint"`
	Username         string   `yaml:"username"`
	Password         string   `yaml:"password"`
	InsecureRegistry bool     `yaml:"insecure-registry"`
	Images           []string `yaml:"images"`
}

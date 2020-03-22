package models

type DockerRegistry struct {
	EndPoint         string
	UserName         string
	Password         string
	InsecureRegistry bool
	Images           []string
}

package models

type Webhook struct {
	Name       string     `yaml:"name"`
	Type       string     `yaml:"type"`
	EndPoint   string     `yaml:"endPoint"`
	Registries Registries `yaml:"registries"`
}

type Registries struct {
	DockerRegistry []DockerRegistry `yaml:"dockerRegistry"`
}

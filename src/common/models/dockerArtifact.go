package models

type DockerArtifact struct {
	CustomKind bool   `json:"customKind"`
	Reference  string `json:"reference"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	Version    string `json:"version"`
}

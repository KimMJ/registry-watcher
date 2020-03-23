package models

type Artifact struct {
	Artifacts []DockerArtifact `json:"artifacts"`
}

func (artifacts *Artifact) AddItem(item DockerArtifact) []DockerArtifact {
	artifacts.Artifacts = append(artifacts.Artifacts, item)
	return artifacts.Artifacts
}

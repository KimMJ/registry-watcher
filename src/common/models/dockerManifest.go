package models

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"time"
)

type DockerManifest struct {
	SchemaVersion int       `json:"schemaVersion"`
	Name          string    `json:"name"`
	Tag           string    `json:"tag"`
	Architecture  string    `json:"architecture"`
	History       []History `json:"history"`
}

type History struct {
	V1Compatibility string `json:"v1Compatibility"`
}

type V1Compatibility struct {
	Created time.Time `json:"created"`
}

func (d *DockerManifest) GetCreationDate() time.Time {
	var history V1Compatibility
	err := json.Unmarshal([]byte(d.History[0].V1Compatibility), &history)
	if err != nil {
		log.Error(err)
		return time.Time{}
	}
	return history.Created
}

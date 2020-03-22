package models

import "time"

type ImageManifest struct {
	Tag          string    `json:"tag"`
	Digest       string    `json:"digest"`
	CreationDate time.Time `json:"creationDate"`
}

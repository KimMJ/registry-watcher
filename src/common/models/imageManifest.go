package models

type ImageManifest struct {
	Tag    string `json:"tag"`
	Digest string `json:"digest"`
	//CreationDate string `json:"creationDate"`
}

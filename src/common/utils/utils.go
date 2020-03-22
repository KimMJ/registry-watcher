package utils

import (
	"encoding/json"
	"log"
)

// Print JSON struct with indent
func PrettyPrintJSON(v interface{}) string {
	s, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		log.Printf("Pretty Print got err #%v", err)
	}
	return string(s)
}

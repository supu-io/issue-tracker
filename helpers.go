package main

import (
	"encoding/json"
	"log"
)

// ToJSON Represents an object as json
func ToJSON(i interface{}) *[]byte {
	json, err := json.Marshal(i)
	if err != nil {
		log.Println(err)
	}
	return &json
}

package main

import (
	"encoding/json"
	"log"
)

func ToJSON(i interface{}) *[]byte {
	json, err := json.Marshal(i)
	if err != nil {
		log.Println(err)
	}
	return &json
}

package main

import (
	"encoding/json"
	"log"

	"github.com/iv-p/mapaccess"
)

func main() {
	j := []byte(`{
		"id": "9b92b11b-b57f-4fa6-af5e-e35a290dc764",	
		"name": "John Doe",
		"friends": [
			{
				"name": "Jaime Mckinney"
			},
			{
				"name": "Evangeline Alvarado"
			},
			{
				"name": "Beth Cantrell"
			}
		]
	}`)
	var deserialised interface{}
	err := json.Unmarshal(j, &deserialised)
	if err != nil {
		panic(err)
	}

	bestFriendName, err := mapaccess.Get(deserialised, "friends[0].name")
	if err != nil {
		panic(err)
	}

	myName, err := mapaccess.Get(deserialised, "name")
	if err != nil {
		panic(err)
	}

	log.Printf("My name is %s and my best friend's name is %s",
		myName, bestFriendName)
}

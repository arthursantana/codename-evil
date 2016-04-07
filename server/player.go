package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Player struct {
	Name   string `json:"name"`
	Color  string `json:"color"`
	Points int    `json:"points"`
}

func playersJSON(w http.ResponseWriter, players []Player) {
	fmt.Fprintf(w, "{\"players\": [")
	separator := ""
	for i := 0; i < len(players); i++ {
		p := players[i]
		pJson, err := json.Marshal(p)
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Fprintf(w, "%v%v", separator, string(pJson))
		separator = ",\n"
	}
	fmt.Fprintf(w, "]}")
}

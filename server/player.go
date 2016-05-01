package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
)

type Player struct {
	Name   string `json:"name"`
	Color  string `json:"color"`
	Points int    `json:"points"`
}

type PlayerList []Player

var players PlayerList

func (players PlayerList) writeJSON(w io.WriteCloser) {
	fmt.Fprintf(w, "\"players\": [")
	separator := ""
	for i := range players {
		p := players[i]
		pJson, err := json.Marshal(p)
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Fprintf(w, "%v%v", separator, string(pJson))
		separator = ",\n"
	}
	fmt.Fprintf(w, "]")
}

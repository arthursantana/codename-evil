package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
)

type Unit struct {
	Id       int    `json:"id"`
	OwnerId  int    `json:"ownerId"`
	PlanetId int    `json:"planetId"`
	ShipId   int    `json:"shipId"`
	Type     string `json:"type"`
}

type UnitList []Unit

var units UnitList

func (unit UnitList) writeJSON(w io.WriteCloser) {
	fmt.Fprintf(w, "\"units\": [")
	separator := ""
	for i := range unit {
		u := unit[i]
		uJson, err := json.Marshal(u)
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Fprintf(w, "%v%v", separator, string(uJson))
		separator = ",\n"
	}
	fmt.Fprintf(w, "]")
}

func (u *Unit) hits() bool {
	effectiveness := 0.008

	bingo := rand.Float64()

	return (bingo < effectiveness)
}

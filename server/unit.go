package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
)

const (
	workerCostPerSoldierUnit   = 25000
	obtaniumCostPerSoldierUnit = 500
)

type Unit struct {
	Id       int    `json:"id"`
	OwnerId  int    `json:"ownerId"`
	PlanetId int    `json:"planetId"`
	ShipId   int    `json:"shipId"`
	Type     string `json:"type"`
	Name     string `json:"name"`
}

func unitJSON(w http.ResponseWriter, unit []Unit) {
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
	effectiveness := 0.005

	bingo := rand.Float64()

	return (bingo < effectiveness)
}

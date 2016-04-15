package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
)

type Ship struct {
	OwnerId  int    `json:"ownerId"`
	PlanetId int    `json:"planetId"`
	Type     string `json:"type"`
	Name     string `json:"name"`

	Position    [2]float64 `json:"position"`
	Destination [2]float64 `json:"destination"`

	// cargo
	Workers  int `json:"workers"`
	Cattle   int `json:"cattle"`
	Obtanium int `json:"obtanium"`
}

func shipsJSON(w http.ResponseWriter, ships []Ship) {
	fmt.Fprintf(w, "\"ships\": [")
	separator := ""
	for i := range ships {
		s := ships[i]
		sJson, err := json.Marshal(s)
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Fprintf(w, "%v%v", separator, string(sJson))
		separator = ",\n"
	}
	fmt.Fprintf(w, "]")
}

func (s *Ship) move() {
	if s.PlanetId != -1 { // ship is docked
		return
	}

	speed := 2.0

	vector := [2]float64{s.Destination[0] - s.Position[0], s.Destination[1] - s.Position[1]}
	norm := math.Sqrt(vector[0]*vector[0] + vector[1]*vector[1])

	if norm == 0 {
		return
	}

	vector = [2]float64{vector[0] / norm * speed, vector[1] / norm * speed}

	s.Position = [2]float64{s.Position[0] + vector[0], s.Position[1] + vector[1]}
}

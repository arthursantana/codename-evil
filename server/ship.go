package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
)

const (
	workerCostPerColonizer   = 1000
	cattleCostPerColonizer   = 1000
	obtaniumCostPerColonizer = 1000
)

type Ship struct {
	Id       int    `json:"id"`
	OwnerId  int    `json:"ownerId"`
	PlanetId int    `json:"planetId"`
	Type     string `json:"type"`
	Name     string `json:"name"`

	Position    [2]int  `json:"position"`
	Destination *Planet `json:"destination"`
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
	speed := 2.5

	vector := [2]float64{float64(s.Destination.Position[0] - s.Position[0]), float64(s.Destination.Position[1] - s.Position[1])}
	norm := math.Sqrt(vector[0]*vector[0] + vector[1]*vector[1])

	if norm <= speed+1 { // reached planet
		if s.Destination.OwnerId == -1 { // unhabited planet, colonize
			s.Destination.Workers = workerCostPerColonizer
			s.Destination.Cattle = cattleCostPerColonizer
			s.Destination.Obtanium = obtaniumCostPerColonizer

			s.Destination.OwnerId = s.OwnerId
			s.Destination.DockSpace--

			s.PlanetId = -1
			s.OwnerId = -1
		} else { // habited planet, dock (later will behave differently if planet is someone else's)
			if s.Destination.DockSpace > 0 {
				s.PlanetId = s.Destination.Id
				s.OwnerId = s.Destination.OwnerId
				s.Destination.DockSpace--
			} else {
				s.PlanetId = -1
				s.OwnerId = -1
			}
		}
	} else {
		vector = [2]float64{vector[0] / norm * speed, vector[1] / norm * speed}

		s.Position = [2]int{s.Position[0] + int(vector[0]), s.Position[1] + int(vector[1])}
	}
}

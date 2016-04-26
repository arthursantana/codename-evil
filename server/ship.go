package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
)

const (
	workerCostPerColonizer   = 50000
	cattleCostPerColonizer   = 20000
	obtaniumCostPerColonizer = 5000

	workerCostPerTrojan   = 0
	cattleCostPerTrojan   = 0
	obtaniumCostPerTrojan = 2000
)

type Ship struct {
	Id       int    `json:"id"`
	OwnerId  int    `json:"ownerId"`
	PlanetId int    `json:"planetId"`
	Type     string `json:"type"`
	Name     string `json:"name"`

	Position    [2]float64 `json:"position"`
	Origin      *Planet    `json:"-"`
	Destination *Planet    `json:"-"`

	UnitSpace int `json:"unitSpace"`
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
	speed := 8.0

	vector := [2]float64{s.Destination.Position[0] - s.Position[0], s.Destination.Position[1] - s.Position[1]}
	norm := math.Sqrt(vector[0]*vector[0] + vector[1]*vector[1])

	if norm <= speed+1 { // reached planet
		switch s.Type {
		case "colonizer":
			if s.Destination.OwnerId == -1 { // unhabited planet, colonize
				s.Destination.Workers = workerCostPerColonizer
				s.Destination.Cattle = cattleCostPerColonizer
				s.Destination.Obtanium = obtaniumCostPerColonizer

				s.Destination.OwnerId = s.OwnerId
				s.Destination.DockSpace--

				s.PlanetId = -1
				s.OwnerId = -1
			} else {
				if s.Destination.OwnerId == s.OwnerId { // my planet, dock
					if s.Destination.DockSpace > 0 {
						s.PlanetId = s.Destination.Id
						s.OwnerId = s.Destination.OwnerId
						s.Destination.DockSpace--
					} else {
						s.PlanetId = s.Destination.Id

						s.OwnerId = s.Destination.OwnerId // temporary hack to avoid ships lost in limbo
					}
				} else { // already colonized, go back
					s.Destination = s.Origin
				}
			}
		case "trojan":
			if s.Destination.OwnerId == -1 { // empty planet, go back
				s.Destination = s.Origin

				return
			} else { // drop them soldiers
				if s.Destination.OwnerId == s.OwnerId { // my planet, dock
					for i := range units {
						if s.Destination.UnitSpace == 0 {
							break
						} else {
							if units[i].ShipId == s.Id {
								units[i].PlanetId = s.Destination.Id
								units[i].ShipId = -1
								s.Destination.UnitSpace--
								s.UnitSpace++
							}
						}
					}

					s.PlanetId = s.Destination.Id
				} else {
					for i := range units {
						if s.Destination.EnemyUnitSpace == 0 {
							break
						} else {
							if units[i].ShipId == s.Id {
								units[i].PlanetId = s.Destination.Id
								units[i].ShipId = -1
								s.Destination.EnemyUnitSpace--
								s.UnitSpace++
							}
						}
					}

					s.Destination = s.Origin // not my planet, go back; temporary
				}
			}
		}
	} else {
		vector = [2]float64{vector[0] / norm * speed, vector[1] / norm * speed}

		s.Position = [2]float64{s.Position[0] + vector[0], s.Position[1] + vector[1]}
	}
}

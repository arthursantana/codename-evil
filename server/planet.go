package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
)

type Planet struct {
	Id      int    `json:"id"`
	OwnerId int    `json:"ownerId"`
	Name    string `json:"name"`

	Position [2]float64 `json:"position"`
	R        int        `json:"r"`

	// resources
	Workers  int `json:"workers"`
	Cattle   int `json:"cattle"`
	Energy   int `json:"energy"`
	Obtanium int `json:"obtanium"`

	BusyWorkers int `json:"busyWorkers"`
	BusyCattle  int `json:"busyCattle"`
	BusyEnergy  int `json:"busyEnergy"`

	NotEnoughWorkers bool `json:"notEnoughWorkers"`
	NotEnoughCattle  bool `json:"notEnoughCattle"`
	NotEnoughEnergy  bool `json:"notEnoughEnergy"`
	NotEnoughFarms   bool `json:"notEnoughFarms"`

	Buildings [][]Building `json:"buildings"`

	DockSpace      int `json:"dockSpace"`
	UnitSpace      int `json:"unitSpace"`
	EnemyUnitSpace int `json:"enemyUnitSpace"`
}

var defaultNames = []string{"Big Rock", "Mercury", "Venus", "Earth", "Mars", "Jupiter", "Saturn", "Neptune", "Uranus", "Pluto", "Urectum", "Tessia", "Sur'Kesh", "Tuchanka", "Omega", "Palaven", "Rannoch", "3834 Zappafrank", "Omicron Persei 8", "Planet 9 from Outer Space"}

func (p *Planet) randomize() {
	index := rand.Intn(len(defaultNames))
	p.Name = defaultNames[index]
	defaultNames = append(defaultNames[:index], defaultNames[index+1:]...)

	p.Position[0] = 50 + float64(rand.Intn(650))
	p.Position[1] = 50 + float64(rand.Intn(650))
	p.R = 4 // + rand.Intn(5)

	p.Workers = 10000
	p.Cattle = 10000
	p.Energy = 0
	p.Obtanium = 2000

	p.DockSpace = 10
	p.UnitSpace = 10
	p.EnemyUnitSpace = 10

	p.Buildings = make([][]Building, p.R)
	for i := range p.Buildings {
		p.Buildings[i] = make([]Building, p.R)

		for j := range p.Buildings[i] {
			p.Buildings[i][j] = Building{"", true}
		}
	}

	p.Buildings[0][0].Type = "hq"
	p.Buildings[0][0].Operational = true

	p.OwnerId = -1
}

func planetsJSON(w http.ResponseWriter, planets []Planet) {
	fmt.Fprintf(w, "\"planets\": [")
	separator := ""
	for i := range planets {
		p := planets[i]
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

func (p *Planet) combat() {
	homeFrags := 0
	awayFrags := 0

	for i := range units {
		if units[i].PlanetId == p.Id {
			if units[i].OwnerId == p.OwnerId {
				if units[i].hits() {
					homeFrags++
				}
			} else {
				if units[i].hits() {
					awayFrags++
				}
			}
		}
	}

	anybodyFromHome := false
	anybodyFromAway := false
	for i := range units {
		if units[i].PlanetId == p.Id {
			if units[i].OwnerId == p.OwnerId {
				if awayFrags > 0 {
					units[i].PlanetId = -1
					units[i].OwnerId = -1
					awayFrags--
				} else {
					anybodyFromHome = true
				}
			} else {
				if homeFrags > 0 {
					units[i].PlanetId = -1
					units[i].OwnerId = -1
					homeFrags--
				} else {
					anybodyFromAway = true
				}
			}
		}
	}

	if anybodyFromAway && !anybodyFromHome { // conquest!
		for i := range units {
			if units[i].PlanetId == p.Id {
				if units[i].OwnerId != p.OwnerId {
					p.OwnerId = units[i].OwnerId
					break
				}
			}
		}
	}
}

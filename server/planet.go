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
}

func (p *Planet) randomize() {
	p.Name = "Big Rock"

	p.Position[0] = 50 + float64(rand.Intn(650))
	p.Position[1] = 50 + float64(rand.Intn(650))
	p.R = 4 //5 + float64(rand.Intn(5))

	p.Workers = 1000
	p.Cattle = 45000
	p.Energy = 0
	p.Obtanium = 0

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

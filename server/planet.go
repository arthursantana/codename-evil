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
	Population int `json:"population"`
	Cattle     int `json:"cattle"`
	Energy     int `json:"energy"`
	Obtanium   int `json:"obtanium"`

	Buildings [][]Building `json:"buildings"`
}

func (p *Planet) randomize() {
	bigString := "akdfya87w36c4iungt2673tc5unyedc7g2j97sa6ged7f6cnrgnydgf7awgcj57g62cnybfubwhe897r6gc7c63k84r"
	pos := rand.Intn(65)
	p.Name = bigString[pos : pos+rand.Intn(10)+5]

	p.Position[0] = 100 + float64(rand.Intn(550))
	p.Position[1] = 100 + float64(rand.Intn(550))
	p.R = 4 //5 + float64(rand.Intn(5))

	p.Population = 1000 + rand.Intn(100000)
	p.Cattle = 1000 + rand.Intn(100000)
	p.Energy = 1000 + rand.Intn(100000)
	p.Obtanium = 1000 + rand.Intn(100000)

	p.Buildings = make([][]Building, p.R)
	for i := range p.Buildings {
		p.Buildings[i] = make([]Building, p.R)

		for j := range p.Buildings[i] {
			p.Buildings[i][j] = Building{}
		}

	}

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

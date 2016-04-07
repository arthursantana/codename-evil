package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
)

type Planet struct {
	Id       int        `json:"id"`
	Name     string     `json:"name"`
	Position [3]float64 `json:"position"`
	R        float64    `json:"r"`
	OwnerId  int        `json:"ownerId"`
}

func (p *Planet) randomize() {
	bigString := "akdfya87w36c4iungt2673tc5unyedc7g2j97sa6ged7f6cnrgnydgf7awgcj57g62cnybfubwhe897r6gc7c63k84r"
	pos := rand.Intn(65)
	p.Name = bigString[pos : pos+rand.Intn(10)+5]

	p.randomizePosition()
	p.randomizeRadius()

	p.OwnerId = -1
}

func (p *Planet) randomizeRadius() {
	p.R = 5 + float64(rand.Intn(100))
}

func (p *Planet) randomizePosition() {
	p.Position[0] = 100 + float64(rand.Intn(800))
	p.Position[1] = 100 + float64(rand.Intn(800))
	p.Position[2] = 100 + float64(rand.Intn(800))
}

func planetsJSON(w http.ResponseWriter, planets []Planet) {
	fmt.Fprintf(w, "{\"planets\": [")
	separator := ""
	for i := len(planets) - 1; i >= 0; i-- {
		p := planets[i]
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

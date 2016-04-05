package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
)

type Planet struct {
	Name    string  `json:"name"`
	X       float64 `json:"x"`
	Y       float64 `json:"y"`
	R       int     `json:"r"`
	OwnerId int     `json:"ownerId"`
}

func (p *Planet) randomize() {
	bigString := "akdfya87w36c4iungt2673tc5unyedc7g2j97sa6ged7f6cnrgnydgf7awgcj57g62cnybfubwhe897r6gc7c63k84r"
	pos := rand.Intn(65)
	p.Name = bigString[pos : pos+rand.Intn(10)+5]

	p.X = 100 + float64(rand.Intn(800))
	p.Y = 100 + float64(rand.Intn(800))
	p.R = 5 + rand.Intn(50)

	p.OwnerId = -1
}

func (p *Planet) rotate(angle float64) {
	x := p.X - 500
	y := p.Y - 500

	p.X = x*math.Cos(angle) - y*math.Sin(angle) + 500
	p.Y = x*math.Sin(angle) + y*math.Cos(angle) + 500
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

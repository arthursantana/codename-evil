package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type Planet struct {
	name string
	x, y int
	r    int
}

func (p *Planet) randomize() {
	bigString := "akdfya87w36c4iungt2673tc5unyedc7g2j97sa6ged7f6cnrgnydgf7awgcj57g62cnybfubwhe897r6gc7c63k84r"
	pos := rand.Intn(65)
	p.name = bigString[pos : pos+rand.Intn(10)+5]

	p.x = 100 + rand.Intn(800)
	p.y = 100 + rand.Intn(800)
	p.r = 5 + rand.Intn(50)
}

type Player struct {
	name  string
	color string
}

func main() {
	var (
		nPlanets = flag.Int("planets", 10, "number of planets")
		nPlayers = flag.Int("players", 1, "number of players")
	)

	flag.Parse()

	// VALIDATE NUMBERS HERE
	// END

	//fmt.Printf("%d players in %d planets\n", *nPlayers, *nPlanets)

	planets := make([]Planet, *nPlanets)
	players := make([]Player, *nPlayers)
	_ = players

	// GENERATE RANDOM STUFF
	rand.Seed(time.Now().UTC().UnixNano())

	for i := 0; i < *nPlanets; i++ {
		planets[i].randomize()
	}

	// SERVE
	fs := http.FileServer(http.Dir("static"))
	fsHandler := http.StripPrefix("/static/", fs)

	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		fsHandler.ServeHTTP(w, r)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/static/index.htm", http.StatusMovedPermanently)
	})
	http.HandleFunc("/data/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "{\"planets\": [")
		separator := ""
		for i := 0; i < *nPlanets; i++ {
			p := planets[i]
			fmt.Fprintf(w, "%v{\n\"name\": \"%v\",\n\"r\": %v,\n\"x\": %v,\n\"y\": %v\n}", separator, p.name, p.r, p.x, p.y)
			separator = ",\n"
		}
		fmt.Fprintf(w, "]}")
	})

	http.ListenAndServe(":8085", nil)
}

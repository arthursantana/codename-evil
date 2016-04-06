package main

import (
	"flag"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	var (
		nPlanets       = flag.Int("planets", 10, "number of planets")
		nPlayers       = flag.Int("players", 1, "number of players")
		lastTick int64 = 0
	)

	flag.Parse()

	// VALIDATE NUMBERS HERE
	// END

	players := make([]Player, *nPlayers)
	planets := make([]Planet, *nPlanets+*nPlayers)

	// GENERATE RANDOM STUFF
	rand.Seed(time.Now().UTC().UnixNano())

	for i := 0; i < len(planets); i++ {
		planets[i].randomize()
	}

	for i := 0; i < len(players); i++ {
		players[i].randomize()
		planets[i].OwnerId = i
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

	http.HandleFunc("/planets/", func(w http.ResponseWriter, r *http.Request) {
		planetsJSON(w, planets)
	})
	http.HandleFunc("/players/", func(w http.ResponseWriter, r *http.Request) {
		playersJSON(w, players)
	})

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	http.HandleFunc("/ws/", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)

		var lastConnTick int64 = 0

		if err != nil {
			log.Println(err)
			return
		}

		for {
			if lastConnTick < lastTick {
				if err = conn.WriteMessage(websocket.TextMessage, []byte("reload")); err != nil {
					log.Println(err)
					return
				}
				lastConnTick = lastTick
			} else {
				time.Sleep(8 * time.Millisecond)
			}
		}
	})

	go func() {
		for {
			time.Sleep(16 * time.Millisecond)
			tick(players, planets)
			lastTick = time.Now().UTC().UnixNano()
		}
	}()

	http.ListenAndServe(":8081", nil)
}

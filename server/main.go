package main

import (
	"encoding/json"
	"flag"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	var (
		nPlanets       = flag.Int("planets", 10, "number of planets")
		speed          = flag.Float64("speed", 10, "speed")
		lastTick int64 = 0
	)

	flag.Parse()

	// VALIDATE NUMBERS HERE
	// END

	planets := make([]Planet, *nPlanets)
	players := make([]Player, 0)

	// GENERATE RANDOM STUFF
	rand.Seed(time.Now().UTC().UnixNano())

	for i := 0; i < len(planets); i++ {
		planets[i].randomize()
		planets[i].Id = i
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

		go func() {
			for {
				if lastConnTick < lastTick {
					if err = conn.WriteMessage(websocket.TextMessage, []byte("tick")); err != nil {
						log.Println(err)
						return
					}
					lastConnTick = lastTick
				} else {
					time.Sleep(8 * time.Millisecond)
				}
			}
		}()

		playerId := -1

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			} else if playerId == -1 { // try to register player
				answer := ""

				if len(players) >= len(planets) {
					answer = "server full"
				} else {
					// FALTA EVITAR CONCORRÊNCIA AQUI
					p := Player{}
					json.Unmarshal(message, &p)
					log.Printf("New player: %v (%v)\n", p.Name, p.Color)

					playerId = len(players)
					players = append(players, p)
					planets[playerId].OwnerId = playerId
					answer = strconv.Itoa(playerId)
				}

				log.Println(string(answer))
				if err = conn.WriteMessage(websocket.TextMessage, []byte(answer)); err != nil {
					log.Println(err)
					return
				}
			} else { // planet was clicked
				planetId, err := strconv.Atoi(string(message))
				if err != nil {
					log.Println(err)
				} else {
					p := &planets[planetId]

					if p.OwnerId != -1 {
						p.R -= 2
						if p.R < 0 {
							p.R = 0.5
						}
					} else {
						p.R += 2
						if p.R > 150 {
							p.R = 1
						}
					}
				}
			}
		}
	})

	go func() {
		for {
			time.Sleep(16 * time.Millisecond)
			tick(players, planets, *speed)
			lastTick = time.Now().UTC().UnixNano()
		}
	}()

	http.ListenAndServe(":8081", nil)
}

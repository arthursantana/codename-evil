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

var rotX, rotY, rotZ int
var dirX, dirY, dirZ int

func main() {
	var (
		nPlanets       = flag.Int("planets", 10, "number of planets")
		speed          = flag.Float64("speed", 5, "speed")
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

				p := Player{}
				json.Unmarshal(message, &p)
				log.Printf("New player attempt: %v (%v)\n", p.Name, p.Color)

				if p.Name == "___reconnect___" {
					playerId, _ = strconv.Atoi(p.Color[1:]) // BWAHAHAHA
					answer = strconv.Itoa(playerId)
				} else {
					if len(players) >= len(planets) {
						answer = "server full"
						log.Printf("Server full.\n")
					} else {
						// FALTA EVITAR CONCORRÊNCIA AQUI
						log.Printf("Ok.\n")
						playerId = len(players)
						players = append(players, p)
						planets[playerId].OwnerId = playerId
						answer = strconv.Itoa(playerId)
					}
				}

				//log.Println(string(answer))
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
						p.R -= 10
						if p.R < 10 {
							p.randomizePosition()
							p.randomizeRadius()

							players[p.OwnerId].Points++
						}
					} else {
						p.R += 20
						if p.R > 150 {
							p.R = 1
						}
					}
				}
			}
		}
	})

	rotX = 2000
	rotY = 2000
	rotZ = 2000
	dirX = 1
	dirY = 1
	dirZ = 1

	go func() {
		for {
			time.Sleep(16 * time.Millisecond)
			tick(players, planets, *speed)
			lastTick = time.Now().UTC().UnixNano()
		}
	}()

	http.ListenAndServe(":8081", nil)
}

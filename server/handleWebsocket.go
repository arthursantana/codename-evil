package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func handleWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	var lastConnTick int64 = 0

	if err != nil {
		log.Println(err)
		return
	}

	go func() {
		for {
			if lastConnTick < lastDataUpdate {
				if err = conn.WriteMessage(websocket.TextMessage, []byte("tick")); err != nil {
					log.Println(err)
					return
				}
				lastConnTick = lastDataUpdate
			} else {
				time.Sleep(100 * time.Millisecond)
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
					planets[playerId].Buildings[0][0].Type = "hq"
					answer = strconv.Itoa(playerId)
				}
			}

			if err = conn.WriteMessage(websocket.TextMessage, []byte(answer)); err != nil {
				log.Println(err)
				return
			}
		} else {
			m := Message{}
			json.Unmarshal(message, &m)

			switch m.Command {
			case "build":
				planets[m.ParamsBuild.Id].Buildings[m.ParamsBuild.I][m.ParamsBuild.J].Type = m.ParamsBuild.Type
			case "changePlanetName":
				planetId := m.ParamsChangePlanetName.Id
				planets[planetId].Name = m.ParamsChangePlanetName.Name
			}
		}
	}
}

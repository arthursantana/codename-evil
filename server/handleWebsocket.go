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
				if planets[m.ParamsBuild.PlanetId].OwnerId == playerId {
					if planets[m.ParamsBuild.PlanetId].Buildings[m.ParamsBuild.I][m.ParamsBuild.J].Type == "" {
						obtaniumCost := 0
						switch m.ParamsBuild.Type {
						case "generator":
							obtaniumCost = 8
						case "nasa":
							obtaniumCost = 40
						case "vale":
							obtaniumCost = 80
						}

						if planets[m.ParamsBuild.PlanetId].Obtanium >= obtaniumCost {
							planets[m.ParamsBuild.PlanetId].Buildings[m.ParamsBuild.I][m.ParamsBuild.J].Type = m.ParamsBuild.Type
							planets[m.ParamsBuild.PlanetId].Buildings[m.ParamsBuild.I][m.ParamsBuild.J].Operational = false
							planets[m.ParamsBuild.PlanetId].Obtanium -= obtaniumCost
						} else {
							// error: not enough obtanium
						}
					} else {
						// error: there's already something there
					}
				} else {
					// error: trying to build in somebody else's planet
				}
			case "buildShip":
				if planets[m.ParamsBuildShip.PlanetId].OwnerId == playerId {
					obtaniumCost := 0

					switch m.ParamsBuildShip.Type {
					case "colonizer":
						obtaniumCost = 100
					}

					if planets[m.ParamsBuildShip.PlanetId].DockSpace > 0 && planets[m.ParamsBuildShip.PlanetId].Obtanium >= obtaniumCost {
						s := Ship{playerId, m.ParamsBuildShip.PlanetId, m.ParamsBuildShip.Type, m.ParamsBuildShip.Name, planets[m.ParamsBuildShip.PlanetId].Position, [2]float64{0, 0}, 0, 0, 0}
						ships = append(ships, s)

						planets[m.ParamsBuildShip.PlanetId].Obtanium -= obtaniumCost
						planets[m.ParamsBuildShip.PlanetId].DockSpace--
					} else {
						// error: not enough space or obtanium
					}
				} else {
					// error: trying to build in somebody else's planet
				}
			case "setVoyage":
				ship := &ships[m.ParamsSetVoyage.ShipId]

				if ship.OwnerId == playerId {
					if ship.PlanetId != -1 {
						if planets[ship.PlanetId].OwnerId == playerId {
							ship.PlanetId = -1

							// needs to test and subtract
							ship.Workers = m.ParamsSetVoyage.Workers
							ship.Cattle = m.ParamsSetVoyage.Cattle
							ship.Obtanium = m.ParamsSetVoyage.Obtanium

							ship.Destination = planets[m.ParamsSetVoyage.DestinationId].Position
						} else {
							// error: how the hell did this happen? (probably custom JSON sent to API)
						}
					} else {
						// error: can't order ships that are not docked
					}
				} else {
					// error: trying to fly somebody else's ship
				}
			case "changePlanetName":
				planetId := m.ParamsChangePlanetName.Id
				planets[planetId].Name = m.ParamsChangePlanetName.Name
			}
		}
	}
}
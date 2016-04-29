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

	go func() { // send tick warnings to client
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
						ticksToBuild := 0
						switch m.ParamsBuild.Type {
						case "farm":
							obtaniumCost = obtaniumCostPerFarm
							ticksToBuild = ticksToBuildFarm
						case "generator":
							obtaniumCost = obtaniumCostPerGenerator
							ticksToBuild = ticksToBuildGenerator
						case "nasa":
							obtaniumCost = obtaniumCostPerNasa
							ticksToBuild = ticksToBuildNasa
						case "vale":
							obtaniumCost = obtaniumCostPerVale
							ticksToBuild = ticksToBuildVale
						case "lockheed":
							obtaniumCost = obtaniumCostPerLockheed
							ticksToBuild = ticksToBuildLockheed
						}

						if planets[m.ParamsBuild.PlanetId].Obtanium >= obtaniumCost {
							planets[m.ParamsBuild.PlanetId].Buildings[m.ParamsBuild.I][m.ParamsBuild.J].Type = m.ParamsBuild.Type
							planets[m.ParamsBuild.PlanetId].Buildings[m.ParamsBuild.I][m.ParamsBuild.J].Operational = false
							planets[m.ParamsBuild.PlanetId].Buildings[m.ParamsBuild.I][m.ParamsBuild.J].TicksUntilDone = ticksToBuild
							planets[m.ParamsBuild.PlanetId].Buildings[m.ParamsBuild.I][m.ParamsBuild.J].TicksUntilProductionDone = 0
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
			case "sellBuilding":
				planetId := m.ParamsSellBuilding.PlanetId
				x := m.ParamsSellBuilding.I
				y := m.ParamsSellBuilding.J

				if planets[planetId].OwnerId == playerId && planets[planetId].Buildings[x][y].Type != "hq" {
					obtaniumCost := 0
					switch planets[planetId].Buildings[x][y].Type {
					case "farm":
						obtaniumCost = obtaniumCostPerFarm
					case "generator":
						obtaniumCost = obtaniumCostPerGenerator
					case "nasa":
						obtaniumCost = obtaniumCostPerNasa
					case "vale":
						obtaniumCost = obtaniumCostPerVale
					case "lockheed":
						obtaniumCost = obtaniumCostPerLockheed
					}

					if planets[planetId].Buildings[x][y].TicksUntilDone > 0 {
						planets[planetId].Obtanium += obtaniumCost
					} else {
						planets[planetId].Obtanium += obtaniumCost / 2
					}
					planets[planetId].Buildings[x][y].Type = ""
					planets[planetId].Buildings[x][y].Operational = true
					planets[planetId].Buildings[x][y].TicksUntilDone = 0
				} else {
					// error: not your planet, dagnabbit!
				}
			case "train":
				if planets[m.ParamsTrain.PlanetId].OwnerId == playerId {
					var relevantSpace *int = nil
					switch m.ParamsTrain.Type {
					case "colonizer", "trojan":
						relevantSpace = &planets[m.ParamsTrain.PlanetId].DockSpace
					case "soldier":
						relevantSpace = &planets[m.ParamsTrain.PlanetId].UnitSpace
					}

					if *relevantSpace > 0 &&
						planets[m.ParamsTrain.PlanetId].Workers >= stats[m.ParamsTrain.Type].workerCost &&
						planets[m.ParamsTrain.PlanetId].Cattle >= stats[m.ParamsTrain.Type].cattleCost &&
						planets[m.ParamsTrain.PlanetId].Obtanium >= stats[m.ParamsTrain.Type].obtaniumCost {

						if len(planets[m.ParamsTrain.PlanetId].Buildings[m.ParamsTrain.I][m.ParamsTrain.J].ProductionQueue) == 0 {
							planets[m.ParamsTrain.PlanetId].Buildings[m.ParamsTrain.I][m.ParamsTrain.J].TicksUntilProductionDone = stats[m.ParamsTrain.Type].ticksToBuild
						}
						planets[m.ParamsTrain.PlanetId].Buildings[m.ParamsTrain.I][m.ParamsTrain.J].ProductionQueue = append(planets[m.ParamsTrain.PlanetId].Buildings[m.ParamsTrain.I][m.ParamsTrain.J].ProductionQueue, m.ParamsTrain.Type)
						(*relevantSpace)--

						planets[m.ParamsTrain.PlanetId].Workers -= stats[m.ParamsTrain.Type].workerCost
						planets[m.ParamsTrain.PlanetId].Cattle -= stats[m.ParamsTrain.Type].cattleCost
						planets[m.ParamsTrain.PlanetId].Obtanium -= stats[m.ParamsTrain.Type].obtaniumCost
					} else {
						// error: not enough space or obtanium
					}
				} else {
					// error: trying to build in somebody else's planet
				}
			case "setDestination":
				ship := &ships[m.ParamsSetDestination.ShipId]

				if ship.OwnerId == playerId {
					if ship.PlanetId != -1 {
						if planets[ship.PlanetId].OwnerId == playerId {
							planets[ship.PlanetId].DockSpace++

							ship.Destination = &planets[m.ParamsSetDestination.DestinationId]

							ship.PlanetId = -1
						} else {
							// error: how the hell did this happen? (probably custom JSON sent to API)
						}
					} else {
						// error: can't order ships that are not docked
					}
				} else {
					// error: trying to fly somebody else's ship
				}
			case "boardShip":
				unit := &units[m.ParamsBoardShip.UnitId]
				ship := &ships[m.ParamsBoardShip.ShipId]

				if ship.OwnerId == playerId {
					if unit.OwnerId == playerId {
						if ship.PlanetId == unit.PlanetId {
							if ship.UnitSpace > 0 {
								unit.PlanetId = -1
								unit.ShipId = ship.Id
								planets[ship.PlanetId].UnitSpace++
								ship.UnitSpace--
							} else {
								// error: no space in ship
							}
						} else {
							// error: ship and unit not at the same place
						}
					} else {
						// error: trying to board with somebody else's units
					}
				} else {
					// error: trying to board somebody else's ship
				}
			case "changePlanetName":
				planetId := m.ParamsChangePlanetName.Id

				if planets[planetId].OwnerId == playerId {
					planets[planetId].Name = m.ParamsChangePlanetName.Name
				} else {
					// error: not your planet, dagnabbit!
				}
			}
		}
	}
}

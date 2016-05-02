package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var mu sync.Mutex

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
				mu.Lock()
				s, err := conn.NextWriter(websocket.TextMessage)
				if err != nil {
					log.Println(err)
					return
				}

				// write data to stream
				fmt.Fprintf(s, "{\"type\": \"dataUpdate\",")
				fmt.Fprintf(s, "\"timestamp\": %v,", lastDataUpdate)
				planets.writeJSON(s)
				fmt.Fprintf(s, ",")
				players.writeJSON(s)
				fmt.Fprintf(s, ",")
				ships.writeJSON(s)
				fmt.Fprintf(s, ",")
				units.writeJSON(s)
				fmt.Fprintf(s, "}")

				if err := s.Close(); err != nil {
					log.Println(err)
					mu.Unlock()
					return
				}

				mu.Unlock()

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

			mu.Lock()
			if err = conn.WriteMessage(websocket.TextMessage, []byte(answer)); err != nil {
				log.Println(err)
				mu.Unlock()
				return
			}
			mu.Unlock()
		} else {
			m := Message{}
			json.Unmarshal(message, &m)

			switch m.Command {
			case "build":
				if planets[m.ParamsBuild.PlanetId].OwnerId == playerId {
					if planets[m.ParamsBuild.PlanetId].Buildings[m.ParamsBuild.I][m.ParamsBuild.J].Type == "" {
						switch m.ParamsBuild.Type {
						case "farm", "generator", "nasa", "vale", "lockheed":
						default:
							// should throw an error here
						}

						if planets[m.ParamsBuild.PlanetId].Obtanium >= buildingStats[m.ParamsBuild.Type].ObtaniumCost {
							planets[m.ParamsBuild.PlanetId].Buildings[m.ParamsBuild.I][m.ParamsBuild.J].Type = m.ParamsBuild.Type
							planets[m.ParamsBuild.PlanetId].Buildings[m.ParamsBuild.I][m.ParamsBuild.J].Operational = false
							planets[m.ParamsBuild.PlanetId].Buildings[m.ParamsBuild.I][m.ParamsBuild.J].TicksUntilDone = buildingStats[m.ParamsBuild.Type].TicksToBuild
							planets[m.ParamsBuild.PlanetId].Buildings[m.ParamsBuild.I][m.ParamsBuild.J].TicksUntilProductionDone = 0
							planets[m.ParamsBuild.PlanetId].Obtanium -= buildingStats[m.ParamsBuild.Type].ObtaniumCost
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
					switch planets[planetId].Buildings[x][y].Type {
					case "farm", "generator", "nasa", "vale", "lockheed":
					default:
						// should throw an error here
					}

					if planets[planetId].Buildings[x][y].TicksUntilDone > 0 {
						planets[planetId].Obtanium += buildingStats[planets[planetId].Buildings[x][y].Type].ObtaniumCost
					} else {
						planets[planetId].Obtanium += buildingStats[planets[planetId].Buildings[x][y].Type].ObtaniumCost / 2
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
					default:
						// should throw an error here
					}

					if *relevantSpace > 0 &&
						planets[m.ParamsTrain.PlanetId].Workers >= stats[m.ParamsTrain.Type].WorkerCost &&
						planets[m.ParamsTrain.PlanetId].Cattle >= stats[m.ParamsTrain.Type].CattleCost &&
						planets[m.ParamsTrain.PlanetId].Obtanium >= stats[m.ParamsTrain.Type].ObtaniumCost {

						if len(planets[m.ParamsTrain.PlanetId].Buildings[m.ParamsTrain.I][m.ParamsTrain.J].ProductionQueue) == 0 {
							planets[m.ParamsTrain.PlanetId].Buildings[m.ParamsTrain.I][m.ParamsTrain.J].TicksUntilProductionDone = stats[m.ParamsTrain.Type].TicksToBuild
						}
						planets[m.ParamsTrain.PlanetId].Buildings[m.ParamsTrain.I][m.ParamsTrain.J].ProductionQueue = append(planets[m.ParamsTrain.PlanetId].Buildings[m.ParamsTrain.I][m.ParamsTrain.J].ProductionQueue, m.ParamsTrain.Type)
						(*relevantSpace)--

						planets[m.ParamsTrain.PlanetId].Workers -= stats[m.ParamsTrain.Type].WorkerCost
						planets[m.ParamsTrain.PlanetId].Cattle -= stats[m.ParamsTrain.Type].CattleCost
						planets[m.ParamsTrain.PlanetId].Obtanium -= stats[m.ParamsTrain.Type].ObtaniumCost
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

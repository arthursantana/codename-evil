package main

import (
	"math"
)

const (
	workerFertility                            = 0.02
	cattleFertility                            = 0.04
	mealsPerCow                                = 20
	cattleReproductionToFoodRateAtStableGrowth = 2
)

func tick() {
	for i := range planets {
		if planets[i].OwnerId == -1 { // nothing happens on ownerless planets
			continue
		}

		freeWorkers := planets[i].Workers

		planets[i].BusyWorkers = 0
		planets[i].BusyCattle = 0
		planets[i].BusyEnergy = 0
		planets[i].NotEnoughWorkers = false
		planets[i].NotEnoughCattle = false
		planets[i].NotEnoughEnergy = false
		planets[i].NotEnoughFarms = false

		// build buildings
	Loop:
		for j := range planets[i].Buildings {
			for k := range planets[i].Buildings[j] {
				if planets[i].Buildings[j][k].TicksUntilDone > 0 { // still being built
					planets[i].Buildings[j][k].TicksUntilDone--
					break Loop
				}
			}
		}

		// count number of operating generators
		generators := 0
		for j := range planets[i].Buildings {
			for k := range planets[i].Buildings[j] {
				if planets[i].Buildings[j][k].TicksUntilDone > 0 { // still being built
					continue
				}

				if planets[i].Buildings[j][k].Type == "generator" {
					planets[i].BusyWorkers -= buildingStats["generator"].Operators
					if freeWorkers > -1*buildingStats["generator"].Operators {
						generators++
						freeWorkers += buildingStats["generator"].Operators
						planets[i].Buildings[j][k].Operational = true
					} else {
						planets[i].Buildings[j][k].Operational = false
						planets[i].NotEnoughWorkers = true
					}
				}
			}
		}

		// count number of buildings
		planets[i].Energy = buildingStats["generator"].EnergyPerTick*generators + buildingStats["hq"].EnergyPerTick
		freeEnergy := planets[i].Energy
		farms := 0
		vales := 0
		nasas := 0
		dummy := 0
		count := &dummy
		for j := range planets[i].Buildings {
			for k := range planets[i].Buildings[j] {
				if planets[i].Buildings[j][k].TicksUntilDone > 0 { // still being built
					continue
				}

				switch planets[i].Buildings[j][k].Type {
				case "farm":
					count = &farms
				case "lockheed":
					count = &dummy
				case "vale":
					count = &vales
				case "nasa":
					count = &nasas
				case "hq", "generator":
					continue
				default:
					// should thrown an error here
				}

				planets[i].BusyWorkers -= buildingStats[planets[i].Buildings[j][k].Type].Operators
				planets[i].BusyEnergy -= buildingStats[planets[i].Buildings[j][k].Type].EnergyPerTick
				if freeWorkers >= -1*buildingStats[planets[i].Buildings[j][k].Type].EnergyPerTick && freeEnergy >= -1*buildingStats[planets[i].Buildings[j][k].Type].EnergyPerTick {
					*count++
					freeWorkers += buildingStats[planets[i].Buildings[j][k].Type].Operators
					freeEnergy += buildingStats[planets[i].Buildings[j][k].Type].EnergyPerTick
					planets[i].Buildings[j][k].Operational = true
				} else {
					if freeWorkers < -1*buildingStats[planets[i].Buildings[j][k].Type].Operators {
						planets[i].NotEnoughWorkers = true
					} else {
						planets[i].NotEnoughEnergy = true
					}
					planets[i].Buildings[j][k].Operational = false
				}
			}
		}

		// build units and ships
		for j := range planets[i].Buildings {
			for k := range planets[i].Buildings[j] {
				if planets[i].Buildings[j][k].Operational && planets[i].Buildings[j][k].TicksUntilProductionDone > 0 {
					planets[i].Buildings[j][k].TicksUntilProductionDone--

					if planets[i].Buildings[j][k].TicksUntilProductionDone == 0 {
						// deliver unit built

						switch planets[i].Buildings[j][k].Type {
						case "lockheed":
							u := Unit{len(units), planets[i].OwnerId, planets[i].Id, -1, planets[i].Buildings[j][k].ProductionQueue[0]}
							units = append(units, u)
						case "nasa":
							s := Ship{len(ships), planets[i].OwnerId, planets[i].Id, planets[i].Buildings[j][k].ProductionQueue[0], planets[i].Position, &planets[i], nil, 3}
							ships = append(ships, s)
						}
						planets[i].Buildings[j][k].ProductionQueue = planets[i].Buildings[j][k].ProductionQueue[1:]

						if len(planets[i].Buildings[j][k].ProductionQueue) > 0 {
							planets[i].Buildings[j][k].TicksUntilProductionDone = stats[planets[i].Buildings[j][k].ProductionQueue[0]].TicksToBuild
						} else {
							planets[i].Buildings[j][k].TicksUntilProductionDone = 0 // not building anything
						}
					}
				}
			}
		}

		// WORKERS AND CATTLE
		// births and cattle consumption
		newWorkers := int(math.Ceil(float64(planets[i].Workers) * workerFertility))
		newCattle := int(math.Ceil(float64(planets[i].Cattle) * cattleFertility))
		planets[i].Workers += newWorkers

		// limit cattle workers
		cattleLimit := (farms)*buildingStats["farm"].MaxCattlePop + buildingStats["hq"].MaxCattlePop
		planets[i].Cattle += newCattle
		if planets[i].Cattle > cattleLimit {
			planets[i].Cattle = cattleLimit
			planets[i].NotEnoughFarms = true
		}

		// consume food
		consumptionRate := float64(newCattle*mealsPerCow) / (cattleReproductionToFoodRateAtStableGrowth * float64(planets[i].Workers))
		cattleConsumption := int(math.Ceil(float64(newCattle) * consumptionRate))

		if cattleConsumption*mealsPerCow > planets[i].Workers {
			cattleConsumption = planets[i].Workers / mealsPerCow
		} else {
			planets[i].Workers = cattleConsumption * mealsPerCow
			planets[i].NotEnoughCattle = true
		}
		planets[i].BusyCattle = planets[i].Cattle - newCattle
		planets[i].Cattle -= cattleConsumption

		// mining
		planets[i].Obtanium += buildingStats["vale"].ObtaniumPerTick*vales + buildingStats["hq"].ObtaniumPerTick

		// move ships
		for i := range ships {
			if ships[i].PlanetId == -1 && ships[i].OwnerId != -1 { // ship is not docked and is alive
				ships[i].move()
			}
		}

		// combat
		for i := range planets {
			if planets[i].OwnerId != -1 {
				planets[i].combat()
			}
		}
	}
}

package main

import (
	"math"
)

const (
	workerFertility                            = 0.02
	cattleFertility                            = 0.04
	mealsPerCow                                = 20
	cattleReproductionToFoodRateAtStableGrowth = 2

	maxCattlePerHQ   = 150000
	maxCattlePerFarm = 50000

	energyPerHQ = 2000

	operatorsPerFarm    = -1000
	energyPerFarm       = -100
	obtaniumCostPerFarm = 100
	ticksToBuildFarm    = 10

	operatorsPerGenerator    = -5000
	energyPerGenerator       = 1000
	obtaniumCostPerGenerator = 200
	ticksToBuildGenerator    = 10

	operatorsPerLockheed    = -15000
	energyPerLockheed       = -500
	obtaniumCostPerLockheed = 200
	ticksToBuildLockheed    = 10

	operatorsPerVale    = -20000
	energyPerVale       = -2000
	obtaniumCostPerVale = 1000
	ticksToBuildVale    = 20

	operatorsPerNasa    = -30000
	energyPerNasa       = -1000
	obtaniumCostPerNasa = 500
	ticksToBuildNasa    = 30

	obtaniumPerVale = 15
	obtaniumPerHQ   = 5
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
					planets[i].BusyWorkers -= operatorsPerGenerator
					if freeWorkers > -1*operatorsPerGenerator {
						generators++
						freeWorkers += operatorsPerGenerator
						planets[i].Buildings[j][k].Operational = true
					} else {
						planets[i].Buildings[j][k].Operational = false
						planets[i].NotEnoughWorkers = true
					}
				}
			}
		}

		// count number of buildings
		planets[i].Energy = energyPerGenerator*generators + energyPerHQ
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

				operatorsPerThing := 0
				energyPerThing := 0
				switch planets[i].Buildings[j][k].Type {
				default:
					continue
				case "farm":
					operatorsPerThing = operatorsPerFarm
					energyPerThing = energyPerFarm
					count = &farms
				case "lockheed":
					operatorsPerThing = operatorsPerLockheed
					energyPerThing = energyPerLockheed
					count = &dummy
				case "vale":
					operatorsPerThing = operatorsPerVale
					energyPerThing = energyPerVale
					count = &vales
				case "nasa":
					operatorsPerThing = operatorsPerNasa
					energyPerThing = energyPerNasa
					count = &nasas
				}

				planets[i].BusyWorkers -= operatorsPerThing
				planets[i].BusyEnergy -= energyPerThing
				if freeWorkers >= -1*operatorsPerThing && freeEnergy >= -1*energyPerThing {
					*count++
					freeWorkers += operatorsPerThing
					freeEnergy += energyPerThing
					planets[i].Buildings[j][k].Operational = true
				} else {
					if freeWorkers < -1*operatorsPerThing {
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
							planets[i].Buildings[j][k].TicksUntilProductionDone = stats[planets[i].Buildings[j][k].ProductionQueue[0]].ticksToBuild
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
		cattleLimit := (farms)*maxCattlePerFarm + maxCattlePerHQ
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
		planets[i].Obtanium += obtaniumPerVale*vales + obtaniumPerHQ

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

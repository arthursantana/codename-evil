package main

import (
	"math"
)

const (
	workerFertility                            = 0.1
	cattleFertility                            = 0.2
	mealsPerCow                                = 20
	cattleReproductionToFoodRateAtStableGrowth = 2

	maxCattlePerHQ   = 50000
	maxCattlePerFarm = 100000

	energyPerHQ        = 100
	energyPerGenerator = 700
	energyPerFarm      = -50
	energyPerNasa      = -1000
	energyPerVale      = -2000

	obtaniumPerVale = 15
	obtaniumPerHQ   = 1

	operatorsPerFarm      = -1000
	operatorsPerGenerator = -5000
	operatorsPerVale      = -50000
	operatorsPerNasa      = -100000
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

		// count number of operating generators
		generators := 0

		for j := range planets[i].Buildings {
			for k := range planets[i].Buildings[j] {
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

		// count number of buildings (except generators and HQ)
		farms := 0
		vales := 0
		nasas := 0
		planets[i].Energy = energyPerGenerator*generators + energyPerHQ
		freeEnergy := planets[i].Energy
		dummy := 0
		count := &dummy
		for j := range planets[i].Buildings {
			for k := range planets[i].Buildings[j] {
				operatorsPerThing := 0
				energyPerThing := 0
				switch planets[i].Buildings[j][k].Type {
				default:
					continue
				case "farm":
					operatorsPerThing = operatorsPerFarm
					energyPerThing = energyPerFarm
					count = &farms
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
		for i := 0; i < len(ships); i++ {
			if ships[i].PlanetId == -1 && ships[i].OwnerId != -1 { // ship is not docked and is alive
				ships[i].move()
				ships[i].tick()
			}
		}
	}
}

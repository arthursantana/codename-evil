package main

import (
	"log"
	"math"
)

const (
	popFertility                               = 0.1
	cattleFertility                            = 0.02
	mealsPerCow                                = 20
	cattleReproductionToFoodRateAtStableGrowth = 2

	maxCattlePerHQ   = 100000
	maxCattlePerFarm = 100000

	energyPerHQ        = 10
	energyPerGenerator = 30
	energyPerFarm      = -1
	energyPerNasa      = -10
	energyPerVale      = -20

	obtaniumPerVale = 15
	obtaniumPerHQ   = 1

	operatorsPerFarm      = 1000
	operatorsPerGenerator = 5000
	operatorsPerVale      = 10000
	operatorsPerNasa      = 20000
)

func tick() {
	for i := range planets {
		if planets[i].OwnerId == -1 { // temporary: nothing happens on ownerless planets
			continue
		}

		freeWorkers := planets[i].Population

		// count number of operating generators
		generators := 0

		for j := range planets[i].Buildings {
			for k := range planets[i].Buildings[j] {
				if planets[i].Buildings[j][k].Type == "generator" {
					if freeWorkers > operatorsPerGenerator {
						generators++
						freeWorkers -= operatorsPerGenerator
						planets[i].Buildings[j][k].Operational = true
					} else {
						planets[i].Buildings[j][k].Operational = false
					}
				}
			}
		}

		// count number of buildings (except generators and HQ)
		farms := 0
		vales := 0
		nasas := 0
		freeEnergy := energyPerGenerator*generators + energyPerHQ
		dummy := 0
		count := &dummy
		for j := range planets[i].Buildings {
			for k := range planets[i].Buildings[j] {
				operatorsPerThing := 0
				energyPerThing := 0
				switch planets[i].Buildings[j][k].Type {
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
				default:
					continue
				}

				if freeWorkers > operatorsPerThing && freeEnergy > energyPerThing {
					*count++
					freeWorkers -= operatorsPerThing
					freeEnergy += energyPerThing
					planets[i].Buildings[j][k].Operational = true
				} else {
					planets[i].Buildings[j][k].Operational = false
				}
			}
		}

		// population and cattle
		{
			// births and cattle consumption
			newPeople := int(math.Ceil(float64(planets[i].Population) * popFertility))
			newCattle := int(math.Ceil(float64(planets[i].Cattle) * cattleFertility))
			planets[i].Population += newPeople

			// consume food
			consumptionRate := float64(newCattle*mealsPerCow) / (cattleReproductionToFoodRateAtStableGrowth * float64(planets[i].Population))
			cattleConsumption := int(math.Ceil(float64(newCattle) * consumptionRate))

			if cattleConsumption*mealsPerCow > planets[i].Population {
				cattleConsumption = planets[i].Population / mealsPerCow
			} else {
				planets[i].Population = cattleConsumption * mealsPerCow
			}
			planets[i].Cattle -= cattleConsumption

			// limit cattle population
			cattleLimit := (farms)*maxCattlePerFarm + maxCattlePerHQ
			planets[i].Cattle += newCattle
			if planets[i].Cattle > cattleLimit {
				planets[i].Cattle = cattleLimit
			}
		}

		// building effects
		planets[i].Obtanium += obtaniumPerVale*vales + obtaniumPerHQ
		planets[i].Energy = freeEnergy

		log.Printf("generators: %v, farms: %v, nasas: %v, vales: %v\n", generators, farms, nasas, vales)
		log.Printf("free workers: %v, free energy: %v", freeWorkers, freeEnergy)
	}
}

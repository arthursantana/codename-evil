package main

func tick(players []Player, planets []Planet) {
	for i := range planets {
		planets[i].Population = int(float64(planets[i].Population) * 1.01)
	}
}

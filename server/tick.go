package main

import "math"

func tick(players []Player, planets []Planet) {
	for i := range planets {
		planets[i].rotate(math.Pi / 600)
	}
}

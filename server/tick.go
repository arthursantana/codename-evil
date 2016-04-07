package main

import "math"

func tick(players []Player, planets []Planet, speed float64) {
	for i := range planets {
		planets[i].R += 0.005
		planets[i].rotate(math.Pi/2000*speed, math.Pi/4000*speed, math.Pi/3000*speed)
	}
}

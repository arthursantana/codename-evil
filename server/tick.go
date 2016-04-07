package main

import (
	"math"
	"math/rand"
)

func tick(players []Player, planets []Planet, speed float64) {
	for i := range planets {
		planets[i].R += 0.005
		planets[i].rotate(math.Pi/float64(rotX+1000)*speed/5, math.Pi/float64(rotY+1000)*speed/5, math.Pi/float64(rotZ+1000)*speed/5)
	}

	limit := 5000
	q := 50
	if rotX > limit {
		rotX = limit
		dirX *= -1
	} else if rotX < -1*limit {
		rotX = -1 * limit
		dirX *= -1
	}

	if rotY > limit {
		rotY = limit
		dirY *= -1
	} else if rotY < -1*limit {
		rotY = -1 * limit
		dirY *= -1
	}

	if rotZ > limit {
		rotZ = limit
		dirZ *= -1
	} else if rotZ < -1*limit {
		rotZ = -1 * limit
		dirZ *= -1
	}

	rotX += rand.Intn(q) * dirX
	rotY += rand.Intn(q) * dirY
	rotZ += rand.Intn(q) * dirZ
}

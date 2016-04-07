package main

import (
	"math"
)

func (p *Planet) rotateX(angle float64) {
	y := p.Position[1] - 500
	z := p.Position[2] - 500

	p.Position[1] = y*math.Cos(angle) - z*math.Sin(angle) + 500
	p.Position[2] = y*math.Sin(angle) + z*math.Cos(angle) + 500
}

func (p *Planet) rotateY(angle float64) {
	x := p.Position[0] - 500
	z := p.Position[2] - 500

	p.Position[2] = z*math.Cos(angle) - x*math.Sin(angle) + 500
	p.Position[0] = z*math.Sin(angle) + x*math.Cos(angle) + 500
}

func (p *Planet) rotateZ(angle float64) {
	x := p.Position[0] - 500
	y := p.Position[1] - 500

	p.Position[0] = x*math.Cos(angle) - y*math.Sin(angle) + 500
	p.Position[1] = x*math.Sin(angle) + y*math.Cos(angle) + 500
}

func (p *Planet) rotate(ax, ay, az float64) {
	p.rotateX(ax)
	p.rotateY(ay)
	p.rotateZ(az)

	if math.IsNaN(p.Position[0]) || math.IsNaN(p.Position[1]) || math.IsNaN(p.Position[2]) {
		p.randomizePosition()
	}
}

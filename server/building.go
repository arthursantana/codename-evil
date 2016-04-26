package main

type Building struct {
	Type                     string   `json:"type"`
	Operational              bool     `json:"operational"`
	TicksUntilDone           int      `json:"ticksUntilDone"`
	TicksUntilProductionDone int      `json:"ticksUntilProductionDone"`
	ProductionQueue          []string `json:"productionQueue"`
}

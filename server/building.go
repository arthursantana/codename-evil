package main

type Building struct {
	Type           string `json:"type"`
	Operational    bool   `json:"operational"`
	TicksUntilDone int    `json:"ticksUntilDone"`
}

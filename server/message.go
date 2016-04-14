package main

type Message struct {
	Command                string                 `json:"command"`
	ParamsChangePlanetName ParamsChangePlanetName `json:"paramsChangePlanetName"`
	ParamsBuild            ParamsBuild            `json:"paramsBuild"`
}

type ParamsChangePlanetName struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type ParamsBuild struct {
	Type     string `json:"type"`
	PlanetId int    `json:"planetId"`
	I        int    `json:"i"`
	J        int    `json:"j"`
}

package main

type Message struct {
	Command                string                 `json:"command"`
	ParamsChangePlanetName ParamsChangePlanetName `json:"paramsChangePlanetName"`
	ParamsBuild            ParamsBuild            `json:"paramsBuild"`
	ParamsSellBuilding     ParamsSellBuilding     `json:"paramsSellBuilding"`
	ParamsBuildShip        ParamsBuildShip        `json:"paramsBuildShip"`
	ParamsSetVoyage        ParamsSetVoyage        `json:"paramsSetVoyage"`
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

type ParamsSellBuilding struct {
	PlanetId int `json:"planetId"`
	I        int `json:"i"`
	J        int `json:"j"`
}

type ParamsBuildShip struct {
	Type     string `json:"type"`
	Name     string `json:"name"`
	PlanetId int    `json:"planetId"`
}

type ParamsSetVoyage struct {
	ShipId        int `json:"shipId"`
	DestinationId int `json:"destinationId"`
}

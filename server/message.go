package main

type Message struct {
	Command                string                 `json:"command"`
	ParamsChangePlanetName ParamsChangePlanetName `json:"paramsChangePlanetName"`
	ParamsBuild            ParamsBuild            `json:"paramsBuild"`
	ParamsSellBuilding     ParamsSellBuilding     `json:"paramsSellBuilding"`
	ParamsTrain            ParamsTrain            `json:"paramsTrain"`
	ParamsSetDestination   ParamsSetDestination   `json:"paramsSetDestination"`
	ParamsBoardShip        ParamsBoardShip        `json:"paramsBoardShip"`
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

type ParamsTrain struct {
	Type     string `json:"type"`
	PlanetId int    `json:"planetId"`
	I        int    `json:"i"`
	J        int    `json:"j"`
}

type ParamsSetDestination struct {
	ShipId        int `json:"shipId"`
	DestinationId int `json:"destinationId"`
}

type ParamsBoardShip struct {
	UnitId int `json:"unitId"`
	ShipId int `json:"shipId"`
}

package main

type Stats struct {
	WorkerCost   int
	CattleCost   int
	ObtaniumCost int
	TicksToBuild int
}

var stats = func() map[string]Stats {
	var s = make(map[string]Stats)

	s["soldier"] = Stats{5000, 0, 500, 10}
	s["colonizer"] = Stats{50000, 20000, 5000, 30}
	s["trojan"] = Stats{0, 0, 2000, 20}

	return s
}()

type BuildingStats struct {
	Operators       int
	EnergyPerTick   int
	ObtaniumPerTick int
	ObtaniumCost    int
	TicksToBuild    int
	MaxCattlePop    int
}

var buildingStats = func() map[string]BuildingStats {
	var s = make(map[string]BuildingStats)

	s["hq"] = BuildingStats{0, 2000, 10, 0, 0, 150000}
	s["farm"] = BuildingStats{-1000, -100, 0, 100, 10, 50000}
	s["generator"] = BuildingStats{-5000, 1000, 0, 200, 10, 0}
	s["lockheed"] = BuildingStats{-15000, -500, 0, 200, 10, 0}
	s["vale"] = BuildingStats{-20000, -2000, 15, 1000, 20, 0}
	s["nasa"] = BuildingStats{-30000, -1000, 0, 500, 30, 0}

	return s
}()

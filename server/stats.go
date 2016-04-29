package main

type Stats struct {
	workerCost   int
	cattleCost   int
	obtaniumCost int
	ticksToBuild int
}

var stats = makeStats()

func makeStats() map[string]Stats {
	var s = make(map[string]Stats)

	// workerCost, cattleCost, obtaniumCost, ticksToBuild
	s["soldier"] = Stats{5000, 0, 500, 10}
	s["colonizer"] = Stats{50000, 20000, 5000, 30}
	s["trojan"] = Stats{0, 0, 2000, 20}

	return s
}

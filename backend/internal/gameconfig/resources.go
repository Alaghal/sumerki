package gameconfig

type ResourceValues struct {
	Gold       int64
	Food       int64
	Wood       int64
	Stone      int64
	Population int64
}

var StartingResources = ResourceValues{
	Gold:       600,
	Food:       400,
	Wood:       400,
	Stone:      300,
	Population: 120,
}

var BaseProductionPerHour = ResourceValues{
	Gold:       20,
	Food:       30,
	Wood:       25,
	Stone:      15,
	Population: 1,
}

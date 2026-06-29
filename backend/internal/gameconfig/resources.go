package gameconfig

type ResourceValues struct {
	Gold       int64
	Food       int64
	Wood       int64
	Stone      int64
	Population int64
}

var StartingResources = ResourceValues{
	Gold:       500,
	Food:       300,
	Wood:       300,
	Stone:      200,
	Population: 100,
}

var BaseProductionPerHour = ResourceValues{
	Gold:       20,
	Food:       30,
	Wood:       25,
	Stone:      15,
	Population: 1,
}

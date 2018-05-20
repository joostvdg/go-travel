package model

type Trip struct {
	Travels          []Travel          `json:"Travels"`
	Activities       []Activity        `json:"Activities"`
	PointsOfInterest []PointOfInterest `json:"PointsOfInterest"`
	Itineraries      []Trip            `json:"Itineraries"`
	Location         string            `json:"Location"`
	LocationType     string            `json:"LocationType"`
	TripType         string            `json:"TripType"`
	Name             string            `json:"Name"`
}

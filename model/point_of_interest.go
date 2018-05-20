package model

type PointOfInterest struct {
	Name             string `json:"Name"`
	Location         string `json:"Location"`
	Description      string `json:"Description"`
	ReasonOfInterest string `json:"ReasonOfInterest"`
	Order            int    `json:"Order"`
}

package model

type Activity struct {
	Name         string `json:"Name"`
	Location     string `json:"Location"`
	Description  string `json:"Description"`
	ActivityType string `json:"Type"`
	Order        int    `json:"Order"`
}

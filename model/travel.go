package model

type Travel struct {
	From  string `json:"From"`
	Too   string `json:"Too"`
	Mode  string `json:"Mode"`
	Order int    `json:"Order"`
}

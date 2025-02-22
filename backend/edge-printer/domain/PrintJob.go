package domain

type PrintJob struct {
	ShopName          string         `json:"jobName"`
	AccountHolderName string         `json:"accountHolderName"`
	Description       string         `json:"description"`
	OrderPositions    map[string]int `json:"orderPositions"`
}

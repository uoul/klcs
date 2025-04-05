package domain

type PrintJob struct {
	TransactionId     string         `json:"transactionId"`
	Timestamp         string         `json:"timestamp"`
	ShopName          string         `json:"jobName"`
	AccountHolderName string         `json:"accountHolderName"`
	Description       string         `json:"description"`
	OrderPositions    map[string]int `json:"orderPositions"`
}

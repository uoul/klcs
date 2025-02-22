package domain

type PrintJob struct {
	ShopName          string
	AccountHolderName string
	Description       string
	OrderPositions    map[string]int
}

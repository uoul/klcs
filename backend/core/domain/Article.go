package domain

type Article struct {
	Id          string
	Name        string
	Description string
	Price       int
	Category    string `json:"-"`
	StockAmount *int
}

package domain

type ArticleDetails struct {
	Id          string
	Name        string
	Description string
	Price       int
	Category    string
	StockAmount *int
	Printer     *Printer
}

package domain

type HistoryArticle struct {
	Id          string
	Name        string
	Description string
	Pieces      int
	PrinterAck  bool `db:"printer_ack"`
}

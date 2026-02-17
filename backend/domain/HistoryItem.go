package domain

import "time"

type HistoryItem struct {
	TransactionId string `db:"id"`
	Timestamp     time.Time
	Description   *string
	AccountHolder *string `db:"holder_name"`
	Articles      []HistoryArticle
}

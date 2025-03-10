package domain

import "time"

type HistoryItem struct {
	TransactionId string
	Timestamp     time.Time
	Description   *string
	AccountHolder *string
	Articles      []HistoryArtilce
}

package domain

import (
	"time"
)

type Transaction struct {
	Id          string
	Timestamp   time.Time
	Type        string
	Amount      int
	Description string
}

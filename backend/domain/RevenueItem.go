package domain

type RevenueItem struct {
	Article string  `db:"name"`
	Amount  int     `db:"amount"`
	Sum     float32 `db:"sum"`
}

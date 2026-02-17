package dal

import (
	"context"

	db "github.com/uoul/go-dbx"
	"github.com/uoul/klcs/backend/core/domain"
)

type HistoryDao struct{}

// GetHistoryForUser implements [IHistoryDao].
func (h *HistoryDao) GetHistoryForUser(ctx context.Context, s db.IDbSession, username string, length int) ([]domain.HistoryItem, error) {
	// Get last x history items (transactions)
	history, err := db.Query[domain.HistoryItem](
		ctx, s,
		`
			SELECT t.id, t.timestamp, t.description, ac.holder_name
			FROM klcs.article_transaction at
				JOIN klcs.transaction t ON (at.transaction_id = t.id)
				JOIN klcs.user u ON (t.user_id = u.id)
				LEFT JOIN klcs.account ac ON (t.account_id = ac.id)
			WHERE u.username ILIKE $1
			GROUP BY t.id, t.timestamp, ac.holder_name, t.description
			ORDER BY t.timestamp DESC
			LIMIT $2
		`,
		username, length,
	)
	if err != nil {
		return nil, err
	}
	// Get Articles for history
	for i, item := range history {
		articles, err := db.Query[domain.HistoryArticle](
			ctx, s,
			`
				SELECT a.id, a.name, a.description, at.pieces, at.printer_ack
				FROM klcs.article_transaction at
					JOIN klcs.article a ON (at.article_id = a.id)
				WHERE at.transaction_id = $1
			`,
			item.TransactionId,
		)
		if err != nil {
			return nil, err
		}
		history[i].Articles = articles
	}
	// Return history
	return history, nil
}

func NewHistoryDao() *HistoryDao {
	return &HistoryDao{}
}

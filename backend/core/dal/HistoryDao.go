package dal

import (
	"database/sql"

	"github.com/uoul/go-common/async"
	"github.com/uoul/go-common/db"
	"github.com/uoul/klcs/backend/core/domain"
)

// ------------------------------------------------------------------------------------------
// Type
// ------------------------------------------------------------------------------------------

type HistoryDao struct{}

// ------------------------------------------------------------------------------------------
// Public
// ------------------------------------------------------------------------------------------

// GetHistoryForUser implements IHistoryDao.
func (h *HistoryDao) GetHistoryForUser(tx *sql.Tx, username string, length int) chan async.ActionResult[[]domain.HistoryItem] {
	r := make(chan async.ActionResult[[]domain.HistoryItem])
	go func() {
		items := <-h.getHistoryItems(tx, length, username)
		if items.Error != nil {
			r <- async.NewErrorActionResult[[]domain.HistoryItem](items.Error)
			return
		}
		for i, item := range items.Result {
			a := <-h.getHistoryArticlesForTransaction(tx, item.TransactionId)
			if a.Error != nil {
				r <- async.NewErrorActionResult[[]domain.HistoryItem](a.Error)
				return
			}
			items.Result[i].Articles = a.Result
		}
		r <- async.ActionResult[[]domain.HistoryItem]{
			Result: items.Result,
			Error:  nil,
		}
	}()
	return r
}

// ------------------------------------------------------------------------------------------
// Private
// ------------------------------------------------------------------------------------------
func (h *HistoryDao) getHistoryItems(tx *sql.Tx, length int, username string) chan async.ActionResult[[]domain.HistoryItem] {
	sql := `
		SELECT t.id, t.timestamp, t.description, ac.holder_name
		FROM klcs.user_article_transaction uat
			JOIN klcs.transaction t ON (uat.transaction_id = t.id)
			JOIN klcs.user u ON (uat.user_id = u.id)
			LEFT JOIN klcs.account ac ON (t.account_id = ac.id)
		WHERE u.username ILIKE $1
		GROUP BY t.id, t.timestamp, ac.holder_name, t.description
		ORDER BY t.timestamp DESC
		LIMIT $2
	`
	return db.QueryStatementTx(
		tx,
		func() ([]any, *domain.HistoryItem) {
			v := domain.HistoryItem{}
			return []any{&v.TransactionId, &v.Timestamp, &v.Description, &v.AccountHolder}, &v
		},
		sql,
		username,
		length,
	)
}

func (h *HistoryDao) getHistoryArticlesForTransaction(tx *sql.Tx, transactionId string) chan async.ActionResult[[]domain.HistoryArtilce] {
	sql := `
		SELECT a.id, a.name, a.description, uat.pieces
		FROM klcs.user_article_transaction uat
			JOIN klcs.article a ON (uat.article_id = a.id)
		WHERE uat.transaction_id = $1
	`
	return db.QueryStatementTx(
		tx,
		func() ([]any, *domain.HistoryArtilce) {
			v := domain.HistoryArtilce{}
			return []any{&v.Id, &v.Name, &v.Description, &v.Pieces}, &v
		},
		sql,
		transactionId,
	)
}

// ------------------------------------------------------------------------------------------
// Constructor
// ------------------------------------------------------------------------------------------

func NewHistoryDao() IHistoryDao {
	return &HistoryDao{}
}

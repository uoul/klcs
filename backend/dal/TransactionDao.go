package dal

import (
	"context"
	"fmt"

	db "github.com/uoul/go-dbx"
	"github.com/uoul/klcs/backend/core/domain"
)

type TransactionDao struct{}

// CreateTransaction implements [ITransactionDao].
func (d *TransactionDao) CreateTransaction(ctx context.Context, s db.IDbSession, userId string, accountId *string, articles map[string]int, transaction domain.Transaction, printDisabled bool) ([]domain.Transaction, error) {
	// Create transaction
	t, err := db.Query[domain.Transaction](
		ctx, s,
		`
			INSERT INTO klcs.transaction (type,amount,description,account_id,user_id)
			VALUES ($1,$2,$3,$4,$5)
			RETURNING id,timestamp,type,amount,description
		`,
		transaction.Type, transaction.Amount, transaction.Description, accountId, userId,
	)
	if err != nil {
		return nil, err
	}
	if len(t) <= 0 {
		return nil, fmt.Errorf("failed to create transaction")
	}
	// Add Reference Article <-> Transcation
	for articleId, amount := range articles {
		// Check if article has printer
		hasPrinter, err := db.Query[string](
			ctx, s,
			`
				SELECT a.id
				FROM klcs.article a
				WHERE a.id = $1 AND a.printer_id IS NOT NULL
			`,
			articleId,
		)
		if err != nil {
			return nil, err
		}
		// Add Reference
		_, err = db.Query[any](
			ctx, s,
			`
			INSERT INTO klcs.article_transaction (article_id, transaction_id, pieces, printer_ack)
			VALUES ($1,$2,$3,$4)
			`,
			articleId, t[0].Id, amount, printDisabled || len(hasPrinter) <= 0,
		)
		if err != nil {
			return nil, err
		}
	}
	// Return transaction
	return t, nil
}

// GetAccountBalance implements [ITransactionDao].
func (d *TransactionDao) GetAccountBalance(ctx context.Context, s db.IDbSession, accountId string) ([]int, error) {
	return db.Query[int](
		ctx, s,
		`
			SELECT coalesce(sum(t.amount), 0)
			FROM klcs.transaction t
			WHERE t.account_id = $1
		`,
		accountId,
	)
}

// GetTranscation implements [ITransactionDao].
func (d *TransactionDao) GetTranscation(ctx context.Context, s db.IDbSession, transactionId string) ([]domain.Transaction, error) {
	return db.Query[domain.Transaction](
		ctx, s,
		`
			SELECT t.id,t.timestamp,t.type,t.amount,t.description
			FROM klcs.transaction t
			WHERE t.id = $1
		`,
		transactionId,
	)
}

func NewTransactionDao() *TransactionDao {
	return &TransactionDao{}
}

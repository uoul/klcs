package dal

import (
	"database/sql"

	"github.com/uoul/go-common/async"
	"github.com/uoul/go-common/db"
	"github.com/uoul/klcs/backend/core/domain"
)

type TransactionDao struct{}

// CreateTransaction implements ITransactionDao.
func (d *TransactionDao) CreateTransaction(tx *sql.Tx, userId string, accountId *string, articles map[string]int, transaction *domain.Transaction) chan async.ActionResult[domain.Transaction] {
	retVal := make(chan async.ActionResult[domain.Transaction])
	go func() {
		t := <-createTransaction(tx, userId, accountId, transaction)
		if t.Error != nil {
			retVal <- async.ActionResult[domain.Transaction]{
				Result: *new(domain.Transaction),
				Error:  t.Error,
			}
		}
		err := createUserArticleReferenceForTransaction(tx, t.Result.Id, articles)
		retVal <- async.ActionResult[domain.Transaction]{
			Result: t.Result,
			Error:  err,
		}
	}()
	return retVal
}

// GetAccountBalance implements ITransactionDao.
func (d *TransactionDao) GetAccountBalance(tx *sql.Tx, accountId *string) chan async.ActionResult[int] {
	sql := `
		SELECT coalesce(sum(t.amount), 0)
		FROM klcs.transaction t
		WHERE t.account_id = $1
	`
	return db.QuerySingleTx(
		tx,
		func() ([]any, *int) {
			v := 0
			return []any{&v}, &v
		},
		sql,
		accountId,
	)
}

// GetTranscation implements ITransactionDao.
func (d *TransactionDao) GetTranscation(tx *sql.Tx, transactionId string) chan async.ActionResult[domain.Transaction] {
	sql := `
		SELECT t.id,t.timestamp,t.type,t.amount,t.description
		FROM klcs.transaction t
		WHERE t.id = $1
	`
	return db.QuerySingleTx(
		tx,
		transactionMapper,
		sql,
		transactionId,
	)
}

func createTransaction(tx *sql.Tx, userId string, accountId *string, transaction *domain.Transaction) chan async.ActionResult[domain.Transaction] {
	sql := `
		INSERT INTO klcs.transaction (type,amount,description,account_id,user_id)
		VALUES ($1,$2,$3,$4,$5)
		RETURNING id,timestamp,type,amount,description
	`
	return db.QuerySingleTx(
		tx,
		transactionMapper,
		sql,
		transaction.Type,
		transaction.Amount,
		transaction.Description,
		accountId,
		userId,
	)
}

func createUserArticleReferenceForTransaction(tx *sql.Tx, transactionId string, articles map[string]int) error {
	sql := `
		INSERT INTO klcs.article_transaction (article_id, transaction_id, pieces)
		VALUES ($1,$2,$3)
	`
	for articleId, amount := range articles {
		r := <-db.ExecStatementTx(
			tx,
			sql,
			articleId,
			transactionId,
			amount,
		)
		if r.Error != nil {
			return r.Error
		}
	}
	return nil
}

func transactionMapper() ([]any, *domain.Transaction) {
	v := domain.Transaction{}
	return []any{&v.Id, &v.Timestamp, &v.Type, &v.Amount, &v.Description}, &v
}

func NewTransactionDao() ITransactionDao {
	return &TransactionDao{}
}

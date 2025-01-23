package dal

import (
	"database/sql"

	"github.com/uoul/go-common/async"
	"github.com/uoul/klcs/backend/oos-core/domain"
)

type ITransactionDao interface {
	GetTranscation(tx *sql.Tx, transactionId string) chan async.ActionResult[domain.Transaction]
	CreateTransaction(tx *sql.Tx, userId string, accountId *string, articles map[string]int, transaction *domain.Transaction) chan async.ActionResult[domain.Transaction]

	GetAccountBalance(tx *sql.Tx, accountId *string) chan async.ActionResult[int]
}

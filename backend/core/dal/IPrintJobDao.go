package dal

import (
	"database/sql"

	"github.com/uoul/go-common/async"
	"github.com/uoul/klcs/backend/core/domain"
)

type IPrintJobDao interface {
	GetPrintOpenJobsForTransaction(tx *sql.Tx, transactionId string) chan async.ActionResult[map[string]domain.PrintJob]
}

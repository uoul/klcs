package dal

import (
	"database/sql"

	"github.com/uoul/go-common/async"
	"github.com/uoul/klcs/backend/core/domain"
)

type IHistoryDao interface {
	GetHistoryForUser(tx *sql.Tx, username string, length int) chan async.ActionResult[[]domain.HistoryItem]
}

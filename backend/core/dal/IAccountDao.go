package dal

import (
	"database/sql"

	"github.com/uoul/go-common/async"
	"github.com/uoul/go-common/db"
	"github.com/uoul/klcs/backend/oos-core/domain"
)

type IAccountDao interface {
	GetAccount(tx *sql.Tx, accountId string) chan async.ActionResult[domain.Account]
	CreateAccount(tx *sql.Tx, account *domain.Account) chan async.ActionResult[domain.Account]
	DeleteAccount(tx *sql.Tx, accountId string) chan async.ActionResult[db.EffectedRows]
	UpdateAccount(tx *sql.Tx, account *domain.Account) chan async.ActionResult[db.EffectedRows]

	GetAll(tx *sql.Tx) chan async.ActionResult[[]domain.Account]
}

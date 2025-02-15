package dal

import (
	"database/sql"

	"github.com/uoul/go-common/async"
	"github.com/uoul/go-common/db"
	"github.com/uoul/klcs/backend/core/domain"
)

type IPrinterDao interface {
	GetPrinter(tx *sql.Tx, printerId string) chan async.ActionResult[domain.Printer]
	CreatePrinter(tx *sql.Tx, shopId string, printer *domain.Printer) chan async.ActionResult[domain.Printer]
	DeletePrinter(tx *sql.Tx, printerId string) chan async.ActionResult[db.EffectedRows]
	UpdatePrinter(tx *sql.Tx, shopId string, printer *domain.Printer) chan async.ActionResult[db.EffectedRows]

	GetPrintersForShop(tx *sql.Tx, shopId string) chan async.ActionResult[[]domain.Printer]
	GetPrinterForArticle(tx *sql.Tx, articleId string) chan async.ActionResult[domain.Printer]
}

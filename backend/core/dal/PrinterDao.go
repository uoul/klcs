package dal

import (
	"database/sql"

	"github.com/uoul/go-common/async"
	"github.com/uoul/go-common/db"
	"github.com/uoul/klcs/backend/core/domain"
)

type PrinterDao struct{}

// GetPrinterForArticle implements IPrinterDao.
func (p *PrinterDao) GetPrinterForArticle(tx *sql.Tx, articleId string) chan async.ActionResult[domain.Printer] {
	sql := `
		SELECT p.id, p.name
		FROM klcs.printer p
			JOIN klcs.article a ON (a.printer_id = p.id)
		WHERE a.id = $1
	`
	return db.QuerySingleTx(
		tx,
		printerMapper,
		sql,
		articleId,
	)
}

// CreatePrinter implements IPrinterDao.
func (p *PrinterDao) CreatePrinter(tx *sql.Tx, shopId string, printer *domain.Printer) chan async.ActionResult[domain.Printer] {
	sql := `
		INSERT INTO klcs.printer (name, shop_id)
		VALUES ($1,$2)
		RETURNING id,name
	`
	return db.QuerySingleTx(
		tx,
		printerMapper,
		sql,
		printer.Name,
		shopId,
	)
}

// DeletePrinter implements IPrinterDao.
func (p *PrinterDao) DeletePrinter(tx *sql.Tx, printerId string) chan async.ActionResult[db.EffectedRows] {
	sql := `
		DELETE FROM klcs.printer WHERE id = $1
	`
	return db.ExecStatementTx(
		tx,
		sql,
		printerId,
	)
}

// GetPrinter implements IPrinterDao.
func (p *PrinterDao) GetPrinter(tx *sql.Tx, printerId string) chan async.ActionResult[domain.Printer] {
	sql := `
		SELECT p.id, p.name
		FROM klcs.printer p
		WHERE p.id = $1
	`
	return db.QuerySingleTx(
		tx,
		printerMapper,
		sql,
		printerId,
	)
}

// GetPrintersForShop implements IPrinterDao.
func (p *PrinterDao) GetPrintersForShop(tx *sql.Tx, shopId string) chan async.ActionResult[[]domain.Printer] {
	sql := `
		SELECT p.id, p.name
		FROM klcs.printer p
		WHERE p.shop_id = $1
	`
	return db.QueryStatementTx(
		tx,
		printerMapper,
		sql,
		shopId,
	)
}

// UpdatePrinter implements IPrinterDao.
func (p *PrinterDao) UpdatePrinter(tx *sql.Tx, shopId string, printer *domain.Printer) chan async.ActionResult[db.EffectedRows] {
	sql := `
		UPDATE klcs.printer
		SET name=$2, shop_id=$3
		WHERE id = $1
	`
	return db.ExecStatementTx(
		tx,
		sql,
		printer.Id,
		printer.Name,
		shopId,
	)
}

func printerMapper() ([]any, *domain.Printer) {
	v := domain.Printer{}
	return []any{&v.Id, &v.Name}, &v
}

func NewPrinterDao() IPrinterDao {
	return &PrinterDao{}
}

package dal

import (
	"context"

	db "github.com/uoul/go-dbx"
	"github.com/uoul/klcs/backend/core/domain"
)

type PrinterDao struct{}

// CreatePrinter implements [IPrinterDao].
func (p *PrinterDao) CreatePrinter(ctx context.Context, s db.IDbSession, shopId string, printer domain.Printer) ([]domain.Printer, error) {
	return db.Query[domain.Printer](
		ctx, s,
		`
			INSERT INTO klcs.printer (name, shop_id)
			VALUES ($1,$2)
			RETURNING id,name
		`,
		printer.Name, shopId,
	)
}

// DeletePrinter implements [IPrinterDao].
func (p *PrinterDao) DeletePrinter(ctx context.Context, s db.IDbSession, printerId string) error {
	_, err := db.Query[any](
		ctx, s,
		`DELETE FROM klcs.printer WHERE id = $1`,
		printerId,
	)
	return err
}

// GetPrinter implements [IPrinterDao].
func (p *PrinterDao) GetPrinter(ctx context.Context, s db.IDbSession, printerId string) ([]domain.Printer, error) {
	return db.Query[domain.Printer](
		ctx, s,
		`
			SELECT p.id, p.name
			FROM klcs.printer p
			WHERE p.id = $1
		`,
		printerId,
	)
}

// GetPrinterForArticle implements [IPrinterDao].
func (p *PrinterDao) GetPrinterForArticle(ctx context.Context, s db.IDbSession, articleId string) ([]domain.Printer, error) {
	return db.Query[domain.Printer](
		ctx, s,
		`
			SELECT p.id, p.name
			FROM klcs.printer p
				JOIN klcs.article a ON (a.printer_id = p.id)
			WHERE a.id = $1
		`,
		articleId,
	)
}

// GetPrintersForShop implements [IPrinterDao].
func (p *PrinterDao) GetPrintersForShop(ctx context.Context, s db.IDbSession, shopId string) ([]domain.Printer, error) {
	return db.Query[domain.Printer](
		ctx, s,
		`
			SELECT p.id, p.name
			FROM klcs.printer p
			WHERE p.shop_id = $1
		`,
		shopId,
	)
}

// UpdatePrinter implements [IPrinterDao].
func (p *PrinterDao) UpdatePrinter(ctx context.Context, s db.IDbSession, shopId string, printer domain.Printer) error {
	_, err := db.Query[any](
		ctx, s,
		`
			UPDATE klcs.printer
			SET name=$2, shop_id=$3
			WHERE id = $1
		`,
		printer.Id,
		printer.Name,
		shopId,
	)
	return err
}

func NewPrinterDao() *PrinterDao {
	return &PrinterDao{}
}

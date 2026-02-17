package dal

import (
	"context"

	db "github.com/uoul/go-dbx"
	"github.com/uoul/klcs/backend/core/domain"
)

type PrintJobDao struct{}

type printJobRow struct {
	PrinterId     string `db:"printer_id"`
	Timestamp     string
	Description   string
	ShopName      string `db:"shop_name"`
	ArticleName   string `db:"article_name"`
	AccountHolder string `db:"holder_name"`
	Amount        int    `db:"pieces"`
}

// AcknowledgeByTransactionId implements [IPrintJobDao].
func (p *PrintJobDao) AcknowledgeByTransactionId(ctx context.Context, s db.IDbSession, printerId string, transactionId string) error {
	_, err := db.Query[any](
		ctx, s,
		`
			UPDATE klcs.article_transaction
			SET printer_ack=true 
			from klcs.article a
			where a.printer_id = $1 and transaction_id = $2 and article_id = a.id 
		`,
		printerId, transactionId,
	)
	return err
}

// GetPrintOpenJobsForTransaction implements [IPrintJobDao].
func (p *PrintJobDao) GetPrintOpenJobsForTransaction(ctx context.Context, s db.IDbSession, transactionId string) (map[string]domain.PrintJob, error) {
	printers := map[string]domain.PrintJob{}
	// Get all articles, that needs to be printed for transactionId
	printJobs, err := db.Query[printJobRow](
		ctx, s,
		`
			select a.printer_id, t.timestamp, coalesce(t.description, '') AS description, s.name AS shop_name, a.name AS article_name, coalesce(ac.holder_name, '') AS holder_name, at.pieces
			from klcs."transaction" t 
				join klcs.article_transaction at on (t.id = at.transaction_id)
				join klcs.article a on (at.article_id = a.id)
				join klcs.shop s on (a.shop_id = s.id)
				left outer join klcs.account ac on (ac.id = t.account_id)
			where a.printer_id is not null and at.printer_ack = false and t.id = $1
		`,
		transactionId,
	)
	if err != nil {
		return nil, err
	}
	// Add open order positions
	for _, pj := range printJobs {
		// Ensure printer is already in result
		if _, exists := printers[pj.PrinterId]; !exists {
			printers[pj.PrinterId] = domain.PrintJob{
				TransactionId:     transactionId,
				Timestamp:         pj.Timestamp,
				ShopName:          pj.ShopName,
				Description:       pj.Description,
				AccountHolderName: pj.AccountHolder,
				OrderPositions:    map[string]int{},
			}
		}
		// Add open order positions
		printers[pj.PrinterId].OrderPositions[pj.ArticleName] += pj.Amount
	}
	// Return Printjobs
	return printers, nil
}

func NewPrintJobDao() *PrintJobDao {
	return &PrintJobDao{}
}

package dal

import (
	"database/sql"

	"github.com/uoul/go-common/async"
	"github.com/uoul/go-common/db"
	"github.com/uoul/klcs/backend/core/domain"
)

type PrintJobDao struct{}

type printJobRow struct {
	PrinterId     string
	Description   string
	ShopName      string
	ArticleName   string
	AccountHolder string
	Amount        int
}

// GetPrintOpenJobsForTransaction implements IPrintJobDao.
func (p *PrintJobDao) GetPrintOpenJobsForTransaction(tx *sql.Tx, transactionId string) chan async.ActionResult[map[string]domain.PrintJob] {
	r := make(chan async.ActionResult[map[string]domain.PrintJob])
	go func() {
		printers := map[string]domain.PrintJob{}
		// Get all Articles, that needs to be printed for transactionId
		rows := <-p.getJobRows(tx, transactionId)
		if rows.Error != nil {
			r <- async.NewErrorActionResult[map[string]domain.PrintJob](rows.Error)
			return
		}
		for _, r := range rows.Result {
			// create new entry, if printer not in map
			if _, exists := printers[r.PrinterId]; !exists {
				printers[r.PrinterId] = domain.PrintJob{
					TransactionId:     transactionId,
					ShopName:          r.ShopName,
					Description:       r.Description,
					AccountHolderName: r.AccountHolder,
					OrderPositions:    map[string]int{},
				}
			}
			// Add order position
			printers[r.PrinterId].OrderPositions[r.ArticleName] += r.Amount
		}
		r <- async.ActionResult[map[string]domain.PrintJob]{
			Result: printers,
			Error:  nil,
		}
	}()
	return r
}

func (p *PrintJobDao) getJobRows(tx *sql.Tx, transactionId string) chan async.ActionResult[[]printJobRow] {
	sql := `
		select a.printer_id, coalesce(t.description, ''), s.name, a.name, coalesce(ac.holder_name, ''), at.pieces
		from klcs."transaction" t 
			join klcs.article_transaction at on (t.id = at.transaction_id)
			join klcs.article a on (at.article_id = a.id)
			join klcs.shop s on (a.shop_id = s.id)
			left outer join klcs.account ac on (ac.id = t.account_id)
		where a.printer_id is not null and at.printer_ack = false and t.id = $1
	`
	return db.QueryStatementTx(
		tx,
		func() ([]any, *printJobRow) {
			v := printJobRow{}
			return []any{&v.PrinterId, &v.Description, &v.ShopName, &v.ArticleName, &v.AccountHolder, &v.Amount}, &v
		},
		sql,
		transactionId,
	)
}

func NewPrintJobDao() IPrintJobDao {
	return &PrintJobDao{}
}

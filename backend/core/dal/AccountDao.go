package dal

import (
	"database/sql"

	"github.com/uoul/go-common/async"
	"github.com/uoul/go-common/db"
	"github.com/uoul/klcs/backend/oos-core/domain"
)

type AccountDao struct{}

// GetAll implements IAccountDao.
func (a *AccountDao) GetAll(tx *sql.Tx) chan async.ActionResult[[]domain.Account] {
	sql := `
		SELECT a.id, a.holder_name, a.locked
		FROM klcs.account a
	`
	return db.QueryStatementTx(
		tx,
		accountMapper,
		sql,
	)
}

// CreateAccount implements IAccountDao.
func (a *AccountDao) CreateAccount(tx *sql.Tx, account *domain.Account) chan async.ActionResult[domain.Account] {
	sql := `
		INSERT INTO klcs.account (holder_name, locked) 
		VALUES ($1,$2)
		RETURNING id, holder_name, locked
	`
	return db.QuerySingleTx(
		tx,
		accountMapper,
		sql,
		account.HolderName,
		account.Locked,
	)
}

// DeleteAccount implements IAccountDao.
func (a *AccountDao) DeleteAccount(tx *sql.Tx, accountId string) chan async.ActionResult[db.EffectedRows] {
	sql := `
		DELETE FROM klcs.account WHERE id = $1
	`
	return db.ExecStatementTx(
		tx,
		sql,
		accountId,
	)
}

// GetAccount implements IAccountDao.
func (a *AccountDao) GetAccount(tx *sql.Tx, accountId string) chan async.ActionResult[domain.Account] {
	sql := `
		SELECT a.id, a.holder_name, a.locked
		FROM klcs.account a
		WHERE a.id = $1
	`
	return db.QuerySingleTx(
		tx,
		accountMapper,
		sql,
		accountId,
	)
}

// UpdateAccount implements IAccountDao.
func (a *AccountDao) UpdateAccount(tx *sql.Tx, account *domain.Account) chan async.ActionResult[db.EffectedRows] {
	sql := `
		UPDATE klcs.account
		SET holder_name=$2,locked=$3
		WHERE id = $1
	`
	return db.ExecStatementTx(
		tx,
		sql,
		account.Id,
		account.HolderName,
		account.Locked,
	)
}

func accountMapper() ([]any, *domain.Account) {
	v := domain.Account{}
	return []any{&v.Id, &v.HolderName, &v.Locked}, &v
}

func NewAccountDao() IAccountDao {
	return &AccountDao{}
}

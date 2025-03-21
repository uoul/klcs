package dal

import (
	"database/sql"

	"github.com/uoul/go-common/async"
	"github.com/uoul/go-common/db"
	"github.com/uoul/klcs/backend/core/domain"
)

type AccountDao struct{}

// GetAll implements IAccountDao.
func (a *AccountDao) GetAll(tx *sql.Tx) chan async.ActionResult[[]domain.Account] {
	sql := `
		SELECT a.id, a.holder_name, a.locked, a.external_id
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
	if len(account.Id) > 0 {
		sql := `
			INSERT INTO klcs.account (id, holder_name, locked, external_id) 
			VALUES ($1,$2,$3,$4)
			RETURNING id, holder_name, locked, external_id
		`
		return db.QuerySingleTx(
			tx,
			accountMapper,
			sql,
			account.Id,
			account.HolderName,
			account.Locked,
			account.ExternalId,
		)
	} else {
		sql := `
			INSERT INTO klcs.account (holder_name, locked, external_id) 
			VALUES ($1,$2,$3)
			RETURNING id, holder_name, locked, external_id
		`
		return db.QuerySingleTx(
			tx,
			accountMapper,
			sql,
			account.HolderName,
			account.Locked,
			account.ExternalId,
		)
	}
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
		SELECT a.id, a.holder_name, a.locked, a.external_id
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

// GetAccountByExternalId implements IAccountDao.
func (a *AccountDao) GetAccountsByExternalId(tx *sql.Tx, externalId string) chan async.ActionResult[[]domain.Account] {
	sql := `
		SELECT a.id, a.holder_name, a.locked, a.external_id
		FROM klcs.account a
		WHERE a.external_id = $1 
	`
	return db.QueryStatementTx(
		tx,
		accountMapper,
		sql,
		externalId,
	)
}

// UpdateAccount implements IAccountDao.
func (a *AccountDao) UpdateAccount(tx *sql.Tx, account *domain.Account) chan async.ActionResult[db.EffectedRows] {
	sql := `
		UPDATE klcs.account
		SET holder_name=$2,locked=$3,external_id=$4
		WHERE id = $1
	`
	return db.ExecStatementTx(
		tx,
		sql,
		account.Id,
		account.HolderName,
		account.Locked,
		account.ExternalId,
	)
}

func accountMapper() ([]any, *domain.Account) {
	v := domain.Account{}
	return []any{&v.Id, &v.HolderName, &v.Locked, &v.ExternalId}, &v
}

func NewAccountDao() IAccountDao {
	return &AccountDao{}
}

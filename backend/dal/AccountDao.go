package dal

import (
	"context"

	db "github.com/uoul/go-dbx"
	"github.com/uoul/klcs/backend/core/domain"
)

type AccountDao struct{}

// CreateAccount implements [IAccountDao].
func (a *AccountDao) CreateAccount(ctx context.Context, s db.IDbSession, account domain.Account) ([]domain.Account, error) {
	if len(account.Id) > 0 {
		return db.Query[domain.Account](
			ctx, s,
			`
				INSERT INTO klcs.account (id, holder_name, locked, external_id) 
				VALUES ($1,$2,$3,$4)
				RETURNING id, holder_name, locked, external_id
			`,
			account.Id, account.HolderName, account.Locked, account.ExternalId,
		)
	} else {
		return db.Query[domain.Account](
			ctx, s,
			`
				INSERT INTO klcs.account (holder_name, locked, external_id) 
				VALUES ($1,$2,$3)
				RETURNING id, holder_name, locked, external_id
			`,
			account.HolderName, account.Locked, account.ExternalId,
		)
	}
}

// DeleteAccount implements [IAccountDao].
func (a *AccountDao) DeleteAccount(ctx context.Context, s db.IDbSession, accountId string) error {
	_, err := db.Query[any](
		ctx, s,
		`DELETE FROM klcs.account WHERE id = $1`,
		accountId,
	)
	return err
}

// GetAccount implements [IAccountDao].
func (a *AccountDao) GetAccount(ctx context.Context, s db.IDbSession, accountId string) ([]domain.Account, error) {
	return db.Query[domain.Account](
		ctx, s,
		`
			SELECT a.id, a.holder_name, a.locked, a.external_id
			FROM klcs.account a
			WHERE a.id = $1
		`,
		accountId,
	)
}

// GetAccountsByExternalId implements [IAccountDao].
func (a *AccountDao) GetAccountsByExternalId(ctx context.Context, s db.IDbSession, externalId string) ([]domain.Account, error) {
	return db.Query[domain.Account](
		ctx, s,
		`
			SELECT a.id, a.holder_name, a.locked, a.external_id
			FROM klcs.account a
			WHERE a.external_id = $1
		`,
		externalId,
	)
}

// GetAll implements [IAccountDao].
func (a *AccountDao) GetAll(ctx context.Context, s db.IDbSession) ([]domain.Account, error) {
	return db.Query[domain.Account](
		ctx, s,
		`
			SELECT a.id, a.holder_name, a.locked, a.external_id
			FROM klcs.account a
		`,
	)
}

// UpdateAccount implements [IAccountDao].
func (a *AccountDao) UpdateAccount(ctx context.Context, s db.IDbSession, account domain.Account) error {
	_, err := db.Query[domain.Account](
		ctx, s,
		`
			UPDATE klcs.account
			SET holder_name=$2,locked=$3,external_id=$4
			WHERE id = $1
		`,
		account.Id, account.HolderName, account.Locked, account.ExternalId,
	)
	return err
}

func NewAccountDao() *AccountDao {
	return &AccountDao{}
}

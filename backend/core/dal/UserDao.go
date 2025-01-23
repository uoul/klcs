package dal

import (
	"database/sql"

	"github.com/uoul/go-common/async"
	"github.com/uoul/go-common/db"
	"github.com/uoul/klcs/backend/oos-core/domain"
)

type UserDao struct{}

// GetAll implements IUserDao.
func (u *UserDao) GetAll(tx *sql.Tx) chan async.ActionResult[[]domain.User] {
	sql := `
		SELECT u.id, u.username, u.name
		FROM klcs."user" u
	`
	return db.QueryStatementTx(
		tx,
		userMapper,
		sql,
	)
}

// GetUserByUsername implements IUserDao.
func (u *UserDao) GetUserByUsername(tx *sql.Tx, username string) chan async.ActionResult[domain.User] {
	sql := `
		SELECT u.id, u.username, u.name
		FROM klcs."user" u
		WHERE u.username like $1
	`
	return db.QuerySingleTx(
		tx,
		userMapper,
		sql,
		username,
	)
}

// AssignUserShopRole implements IUserDao.
func (u *UserDao) AssignUserShopRole(tx *sql.Tx, userId string, shopId string, roleId string) chan async.ActionResult[db.EffectedRows] {
	sql := `
		INSERT INTO klcs.user_shop_role (user_id, shop_id, role_id) VALUES ($1,$2,$3) ON CONFLICT DO NOTHING
	`
	return db.ExecStatementTx(
		tx,
		sql,
		userId,
		shopId,
		roleId,
	)
}

// GetUsersForShop implements IUserDao.
func (u *UserDao) GetUsersForShop(tx *sql.Tx, shopId string) chan async.ActionResult[[]domain.User] {
	sql := `
		SELECT u.id, u.username, u.name
		FROM klcs."user" u
			JOIN klcs.user_shop_role usr ON (usr.user_id = u.id)
		WHERE usr.shop_id = $1
		GROUP BY u.id, u.username, u.name
	`
	return db.QueryStatementTx(
		tx,
		userMapper,
		sql,
		shopId,
	)
}

// UnassignUserShopRole implements IUserDao.
func (u *UserDao) UnassignUserShopRole(tx *sql.Tx, userId string, shopId string, roleId string) chan async.ActionResult[db.EffectedRows] {
	sql := `
		DELETE FROM klcs.user_shop_role WHERE user_id = $1 AND shop_id = $2 AND role_id = $3
	`
	return db.ExecStatementTx(
		tx,
		sql,
		userId,
		shopId,
		roleId,
	)
}

// CreateOrUpdateUser implements IUserDao.
func (u *UserDao) CreateOrUpdateUser(tx *sql.Tx, user *domain.User) chan async.ActionResult[domain.User] {
	retVal := make(chan async.ActionResult[domain.User])
	go func() {
		dbUser := <-getUserByUsername(tx, user.Username)
		switch dbUser.Error {
		case sql.ErrNoRows:
			retVal <- <-createUser(tx, user)
		case nil:
			if *user == dbUser.Result {
				retVal <- dbUser
			}
			r := <-updateUser(tx, user)
			retVal <- async.ActionResult[domain.User]{
				Error:  r.Error,
				Result: dbUser.Result,
			}
		default:
			retVal <- async.ActionResult[domain.User]{
				Error:  dbUser.Error,
				Result: dbUser.Result,
			}
		}
	}()
	return retVal
}

func getUserByUsername(tx *sql.Tx, username string) chan async.ActionResult[domain.User] {
	sql := `
		SELECT u.id, u.username, u.name
		FROM klcs."user" u
		WHERE u.username = $1
	`
	return db.QuerySingleTx[domain.User](
		tx,
		userMapper,
		sql,
		username,
	)
}

func createUser(tx *sql.Tx, user *domain.User) chan async.ActionResult[domain.User] {
	sql := `
		INSERT INTO klcs."user" (username, name)
		VALUES ($1,$2)
		RETURNING id, username, name
	`
	return db.QuerySingleTx[domain.User](
		tx,
		userMapper,
		sql,
		user.Username,
		user.Name,
	)
}

func updateUser(tx *sql.Tx, user *domain.User) chan async.ActionResult[db.EffectedRows] {
	sql := `
		UPDATE klcs."user"
		SET name=$2
		WHERE username = $1
	`
	return db.ExecStatementTx(
		tx,
		sql,
		user.Username,
		user.Name,
	)
}

func userMapper() ([]any, *domain.User) {
	v := domain.User{}
	return []any{&v.Id, &v.Username, &v.Name}, &v
}

func NewUserDao() IUserDao {
	return &UserDao{}
}

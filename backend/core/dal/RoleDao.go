package dal

import (
	"database/sql"

	"github.com/uoul/go-common/async"
	"github.com/uoul/go-common/db"
	"github.com/uoul/klcs/backend/oos-core/domain"
)

type RoleDao struct{}

// GetRoles implements IRoleDao.
func (r *RoleDao) GetRoles(tx *sql.Tx) chan async.ActionResult[[]domain.Role] {
	sql := `
		SELECT r.id, r.name
		FROM klcs."role" r
	`
	return db.QueryStatementTx(
		tx,
		roleMapper,
		sql,
	)
}

// GetRoleByName implements IRoleDao.
func (r *RoleDao) GetRoleByName(tx *sql.Tx, roleName string) chan async.ActionResult[domain.Role] {
	sql := `
		SELECT r.id, r.name
		FROM klcs."role" r
		WHERE r.name ilike $1
	`
	return db.QuerySingleTx(
		tx,
		roleMapper,
		sql,
		roleName,
	)
}

// GetUserRolesForShop implements IRoleDao.
func (r *RoleDao) GetUserRolesForShop(tx *sql.Tx, username string, shopId string) chan async.ActionResult[[]domain.Role] {
	sql := `
		SELECT r.id, r.name
		FROM klcs."role" r
			JOIN klcs.user_shop_role usr ON (usr.role_id = r.id)
			JOIN klcs."user" u ON (usr.user_id = u.id)
		WHERE u.username ilike $1
			AND usr.shop_id = $2
	`
	return db.QueryStatementTx[domain.Role](
		tx,
		roleMapper,
		sql,
		username,
		shopId,
	)
}

func roleMapper() ([]any, *domain.Role) {
	v := domain.Role{}
	return []any{&v.Id, &v.Name}, &v
}

func NewRoleDao() IRoleDao {
	return &RoleDao{}
}

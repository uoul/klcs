package dal

import (
	"database/sql"

	"github.com/uoul/go-common/async"
	"github.com/uoul/klcs/backend/core/domain"
)

type IRoleDao interface {
	GetRoles(tx *sql.Tx) chan async.ActionResult[[]domain.Role]
	GetRoleByName(tx *sql.Tx, roleName string) chan async.ActionResult[domain.Role]
	GetUserRolesForShop(tx *sql.Tx, username string, shopId string) chan async.ActionResult[[]domain.Role]
}

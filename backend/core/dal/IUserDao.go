package dal

import (
	"database/sql"

	"github.com/uoul/go-common/async"
	"github.com/uoul/go-common/db"
	"github.com/uoul/klcs/backend/oos-core/domain"
)

type IUserDao interface {
	GetAll(tx *sql.Tx) chan async.ActionResult[[]domain.User]
	GetUserByUsername(tx *sql.Tx, username string) chan async.ActionResult[domain.User]
	CreateOrUpdateUser(tx *sql.Tx, user *domain.User) chan async.ActionResult[domain.User]
	GetUsersForShop(tx *sql.Tx, shopId string) chan async.ActionResult[[]domain.User]
	AssignUserShopRole(tx *sql.Tx, userId, shopId, roleId string) chan async.ActionResult[db.EffectedRows]
	UnassignUserShopRole(tx *sql.Tx, userId, shopId, roleId string) chan async.ActionResult[db.EffectedRows]
}

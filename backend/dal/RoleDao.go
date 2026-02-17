package dal

import (
	"context"

	db "github.com/uoul/go-dbx"
	"github.com/uoul/klcs/backend/core/domain"
)

type RoleDao struct{}

// GetRoleByName implements [IRoleDao].
func (r *RoleDao) GetRoleByName(ctx context.Context, s db.IDbSession, roleName string) ([]domain.Role, error) {
	return db.Query[domain.Role](
		ctx, s,
		`
			SELECT r.id, r.name
			FROM klcs."role" r
			WHERE r.name ilike $1
		`,
		roleName,
	)
}

// GetRoles implements [IRoleDao].
func (r *RoleDao) GetRoles(ctx context.Context, s db.IDbSession) ([]domain.Role, error) {
	return db.Query[domain.Role](
		ctx, s,
		`
			SELECT r.id, r.name
			FROM klcs."role" r
		`,
	)
}

// GetUserRolesForShop implements [IRoleDao].
func (r *RoleDao) GetUserRolesForShop(ctx context.Context, s db.IDbSession, username string, shopId string) ([]domain.Role, error) {
	return db.Query[domain.Role](
		ctx, s,
		`
			SELECT r.id, r.name
			FROM klcs."role" r
				JOIN klcs.user_shop_role usr ON (usr.role_id = r.id)
				JOIN klcs."user" u ON (usr.user_id = u.id)
			WHERE u.username ilike $1
				AND usr.shop_id = $2
		`,
		username, shopId,
	)
}

func NewRoleDao() *RoleDao {
	return &RoleDao{}
}

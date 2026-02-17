package dal

import (
	"context"

	db "github.com/uoul/go-dbx"
	"github.com/uoul/klcs/backend/core/domain"
)

type UserDao struct{}

// AssignUserShopRole implements [IUserDao].
func (u *UserDao) AssignUserShopRole(ctx context.Context, s db.IDbSession, userId string, shopId string, roleId string) error {
	_, err := db.Query[any](
		ctx, s,
		`INSERT INTO klcs.user_shop_role (user_id, shop_id, role_id) VALUES ($1,$2,$3) ON CONFLICT DO NOTHING`,
		userId, shopId, roleId,
	)
	return err
}

// UnassignUserShopRole implements [IUserDao].
func (u *UserDao) UnassignUserShopRole(ctx context.Context, s db.IDbSession, userId string, shopId string, roleId string) error {
	_, err := db.Query[any](
		ctx, s,
		`DELETE FROM klcs.user_shop_role WHERE user_id = $1 AND shop_id = $2 AND role_id = $3`,
		userId, shopId, roleId,
	)
	return err
}

// GetAll implements [IUserDao].
func (u *UserDao) GetAll(ctx context.Context, s db.IDbSession) ([]domain.User, error) {
	return db.Query[domain.User](
		ctx, s,
		`
			SELECT u.id, u.username, u.name
			FROM klcs."user" u
		`,
	)
}

// GetUserByUsername implements [IUserDao].
func (u *UserDao) GetUserByUsername(ctx context.Context, s db.IDbSession, username string) ([]domain.User, error) {
	return db.Query[domain.User](
		ctx, s,
		`
			SELECT u.id, u.username, u.name
			FROM klcs."user" u
			WHERE u.username like $1
		`,
		username,
	)
}

// GetUsersForShop implements [IUserDao].
func (u *UserDao) GetUsersForShop(ctx context.Context, s db.IDbSession, shopId string) ([]domain.User, error) {
	return db.Query[domain.User](
		ctx, s,
		`
			SELECT u.id, u.username, u.name
			FROM klcs."user" u
				JOIN klcs.user_shop_role usr ON (usr.user_id = u.id)
			WHERE usr.shop_id = $1
			GROUP BY u.id, u.username, u.name
		`,
		shopId,
	)
}

// CreateOrUpdateUser implements [IUserDao].
func (u *UserDao) CreateOrUpdateUser(ctx context.Context, s db.IDbSession, user domain.User) ([]domain.User, error) {
	// Get user by username
	users, err := db.Query[domain.User](
		ctx, s,
		`
			SELECT u.id, u.username, u.name
			FROM klcs."user" u
			WHERE u.username = $1
		`,
		user.Username,
	)
	if err != nil {
		return nil, err
	}
	// Create or update user
	if len(users) <= 0 {
		// Create User
		return db.Query[domain.User](
			ctx, s,
			`
				INSERT INTO klcs."user" (username, name)
				VALUES ($1,$2)
				RETURNING id, username, name
			`,
			user.Username, user.Name,
		)
	} else {
		return db.Query[domain.User](
			ctx, s,
			`
				UPDATE klcs."user"
				SET name=$2
				WHERE username = $1
				RETURNING id, username, name
			`,
			user.Username, user.Name,
		)
	}
}

func NewUserDao() *UserDao {
	return &UserDao{}
}

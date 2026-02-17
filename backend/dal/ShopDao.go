package dal

import (
	"context"

	"github.com/lib/pq"
	db "github.com/uoul/go-dbx"
	"github.com/uoul/klcs/backend/core/domain"
)

type ShopDao struct{}

// CreateShop implements [IShopDao].
func (*ShopDao) CreateShop(ctx context.Context, s db.IDbSession, shop domain.Shop) ([]domain.Shop, error) {
	return db.Query[domain.Shop](
		ctx, s,
		`INSERT INTO klcs.shop (name) VALUES ($1) RETURNING id, name`,
		shop.Name,
	)
}

// DeleteShop implements [IShopDao].
func (*ShopDao) DeleteShop(ctx context.Context, s db.IDbSession, shopId string) error {
	_, err := db.Query[any](
		ctx, s,
		`DELETE FROM klcs.shop WHERE id = $1`,
		shopId,
	)
	return err
}

// GetAll implements [IShopDao].
func (*ShopDao) GetAll(ctx context.Context, s db.IDbSession) ([]domain.Shop, error) {
	return db.Query[domain.Shop](
		ctx, s,
		`
			SELECt s.id, s.name
			FROM klcs.shop s
		`,
	)
}

// GetShop implements [IShopDao].
func (*ShopDao) GetShop(ctx context.Context, s db.IDbSession, shopId string) ([]domain.Shop, error) {
	return db.Query[domain.Shop](
		ctx, s,
		`
			SELECT s.id, s.name
			FROM klcs.shop s
			WHERE s.id = $1
		`,
		shopId,
	)
}

// GetShopForArticle implements [IShopDao].
func (*ShopDao) GetShopForArticle(ctx context.Context, s db.IDbSession, articleId string) ([]domain.Shop, error) {
	return db.Query[domain.Shop](
		ctx, s,
		`
			SELECT s.id, s.name
			FROM klcs.shop s
				JOIN klcs.article a ON (a.shop_id = s.id)
			WHERE a.id = $1
		`,
		articleId,
	)
}

// GetShopForPrinter implements [IShopDao].
func (*ShopDao) GetShopForPrinter(ctx context.Context, s db.IDbSession, printerId string) ([]domain.Shop, error) {
	return db.Query[domain.Shop](
		ctx, s,
		`
			SELECT s.id, s.name
			FROM klcs.shop s
				JOIN klcs.printer p ON (s.id = p.shop_id)
			WHERE p.id = $1
		`,
		printerId,
	)
}

// GetShopsForArticles implements [IShopDao].
func (*ShopDao) GetShopsForArticles(ctx context.Context, s db.IDbSession, articleIds []string) ([]domain.Shop, error) {
	return db.Query[domain.Shop](
		ctx, s,
		`
			SELECT s.id, s.name
			FROM klcs.shop s
				JOIN klcs.article a ON (a.shop_id = s.id)
			WHERE a.id = ANY($1)
		`,
		pq.Array(articleIds),
	)
}

// GetShopsForUser implements [IShopDao].
func (*ShopDao) GetShopsForUser(ctx context.Context, s db.IDbSession, username string) ([]domain.Shop, error) {
	return db.Query[domain.Shop](
		ctx, s,
		`
			SELECT s.id, s.name
			FROM klcs.shop s
				JOIN klcs.user_shop_role usr ON (usr.shop_id = s.id)
				JOIN klcs."user" u ON (usr.user_id = u.id)
			WHERE u.username ilike $1
			GROUP BY s.id, s.name
		`,
		username,
	)
}

// UpdateShop implements [IShopDao].
func (*ShopDao) UpdateShop(ctx context.Context, s db.IDbSession, shop domain.Shop) error {
	_, err := db.Query[any](
		ctx, s,
		`UPDATE klcs.shop SET name = $1 WHERE id = $2`,
		shop.Name, shop.Id,
	)
	return err
}

func NewShopDao() *ShopDao {
	return &ShopDao{}
}

package dal

import (
	"database/sql"

	"github.com/lib/pq"
	"github.com/uoul/go-common/async"
	"github.com/uoul/go-common/db"
	"github.com/uoul/klcs/backend/oos-core/domain"
)

type ShopDao struct{}

// GetShopForPrinter implements IShopDao.
func (s *ShopDao) GetShopForPrinter(tx *sql.Tx, printerId string) chan async.ActionResult[domain.Shop] {
	sql := `
		SELECT s.id, s.name
		FROM klcs.shop s
			JOIN klcs.printer p ON (s.id = p.shop_id)
		WHERE p.id = $1
	`
	return db.QuerySingleTx(
		tx,
		shopMapper,
		sql,
		printerId,
	)
}

// GetAll implements IShopDao.
func (s *ShopDao) GetAll(tx *sql.Tx) chan async.ActionResult[[]domain.Shop] {
	sql := `
		SELECt s.id, s.name
		FROM klcs.shop s
	`
	return db.QueryStatementTx(
		tx,
		shopMapper,
		sql,
	)
}

// GetShopsForArticles implements IShopDao.
func (s *ShopDao) GetShopsForArticles(tx *sql.Tx, articleIds []string) chan async.ActionResult[[]domain.Shop] {
	sql := `
		SELECT s.id, s.name
		FROM klcs.shop s
			JOIN klcs.article a ON (a.shop_id = s.id)
		WHERE a.id = ANY($1)
	`
	return db.QueryStatementTx(
		tx,
		shopMapper,
		sql,
		pq.Array(articleIds),
	)
}

// GetShopForArticle implements IShopDao.
func (s *ShopDao) GetShopForArticle(tx *sql.Tx, articleId string) chan async.ActionResult[domain.Shop] {
	sql := `
		SELECT s.id, s.name
		FROM klcs.shop s
			JOIN klcs.article a ON (a.shop_id = s.id)
		WHERE a.id = $1
	`
	return db.QuerySingleTx(
		tx,
		shopMapper,
		sql,
		articleId,
	)
}

// GetShop implements IShopDao.
func (s *ShopDao) GetShop(tx *sql.Tx, shopId string) chan async.ActionResult[domain.Shop] {
	sql := `
		SELECT s.id, s.name
		FROM klcs.shop s
		WHERE s.id = $1
	`
	return db.QuerySingleTx[domain.Shop](
		tx,
		shopMapper,
		sql,
		shopId,
	)
}

// CreateShop implements IShopDao.
func (s *ShopDao) CreateShop(tx *sql.Tx, shop *domain.Shop) chan async.ActionResult[domain.Shop] {
	sql := `
		INSERT INTO klcs.shop (name) VALUES ($1) RETURNING id, name
	`
	return db.QuerySingleTx[domain.Shop](
		tx,
		shopMapper,
		sql,
		shop.Name,
	)
}

// DeleteShop implements IShopDao.
func (s *ShopDao) DeleteShop(tx *sql.Tx, shopId string) chan async.ActionResult[db.EffectedRows] {
	sql := `
		DELETE FROM klcs.shop WHERE id = $1
	`
	return db.ExecStatementTx(
		tx,
		sql,
		shopId,
	)
}

// UpdateShop implements IShopDao.
func (s *ShopDao) UpdateShop(tx *sql.Tx, shop *domain.Shop) chan async.ActionResult[db.EffectedRows] {
	sql := `
		UPDATE klcs.shop SET name = $1 WHERE id = $2
	`
	return db.ExecStatementTx(
		tx,
		sql,
		shop.Name,
		shop.Id,
	)
}

// GetShopsForUser implements IShopDao.
func (s *ShopDao) GetShopsForUser(tx *sql.Tx, username string) chan async.ActionResult[[]domain.Shop] {
	sql := `
		SELECT s.id, s.name
		FROM klcs.shop s
			JOIN klcs.user_shop_role usr ON (usr.shop_id = s.id)
			JOIN klcs."user" u ON (usr.user_id = u.id)
		WHERE u.username ilike $1
		GROUP BY s.id, s.name
	`
	return db.QueryStatementTx[domain.Shop](
		tx,
		func() ([]any, *domain.Shop) {
			v := domain.Shop{}
			return []any{&v.Id, &v.Name}, &v
		},
		sql,
		username,
	)
}

func shopMapper() ([]any, *domain.Shop) {
	v := domain.Shop{}
	return []any{&v.Id, &v.Name}, &v
}

func NewShopDao() IShopDao {
	return &ShopDao{}
}

package dal

import (
	"database/sql"

	"github.com/uoul/go-common/async"
	"github.com/uoul/go-common/db"
	"github.com/uoul/klcs/backend/core/domain"
)

type IShopDao interface {
	GetAll(tx *sql.Tx) chan async.ActionResult[[]domain.Shop]
	GetShop(tx *sql.Tx, shopId string) chan async.ActionResult[domain.Shop]
	CreateShop(tx *sql.Tx, shop *domain.Shop) chan async.ActionResult[domain.Shop]
	DeleteShop(tx *sql.Tx, shopId string) chan async.ActionResult[db.EffectedRows]
	UpdateShop(tx *sql.Tx, shop *domain.Shop) chan async.ActionResult[db.EffectedRows]

	GetShopsForUser(tx *sql.Tx, username string) chan async.ActionResult[[]domain.Shop]
	GetShopForArticle(tx *sql.Tx, articleId string) chan async.ActionResult[domain.Shop]
	GetShopsForArticles(tx *sql.Tx, articleIds []string) chan async.ActionResult[[]domain.Shop]
	GetShopForPrinter(tx *sql.Tx, printerId string) chan async.ActionResult[domain.Shop]
}

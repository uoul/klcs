package dal

import (
	"database/sql"

	"github.com/uoul/go-common/async"
	"github.com/uoul/go-common/db"
	"github.com/uoul/klcs/backend/core/domain"
)

type IArticleDao interface {
	GetArticle(tx *sql.Tx, articleId string) chan async.ActionResult[domain.Article]
	CreateArticle(tx *sql.Tx, article *domain.Article, shopId string, printerId *string) chan async.ActionResult[domain.Article]
	DeleteArticle(tx *sql.Tx, articleId string) chan async.ActionResult[db.EffectedRows]
	UpdateArticle(tx *sql.Tx, article *domain.Article) chan async.ActionResult[db.EffectedRows]

	GetArticlesForShop(tx *sql.Tx, shopId string) chan async.ActionResult[[]domain.Article]
	GetArticlesIn(tx *sql.Tx, articleIds []string) chan async.ActionResult[[]domain.Article]
	SetPrinterForArticle(tx *sql.Tx, articleId string, printerId *string) chan async.ActionResult[db.EffectedRows]
}

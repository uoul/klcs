package dal

import (
	"database/sql"

	"github.com/lib/pq"
	"github.com/uoul/go-common/async"
	"github.com/uoul/go-common/db"
	"github.com/uoul/klcs/backend/oos-core/domain"
)

type ArticleDao struct{}

// SetPrinterForArticle implements IArticleDao.
func (a *ArticleDao) SetPrinterForArticle(tx *sql.Tx, articleId string, printerId *string) chan async.ActionResult[db.EffectedRows] {
	sql := `
		UPDATE klcs.article
		SET printer_id=$2
		WHERE id = $1
	`
	return db.ExecStatementTx(
		tx,
		sql,
		articleId,
		printerId,
	)
}

// GetArticlesIn implements IArticleDao.
func (a *ArticleDao) GetArticlesIn(tx *sql.Tx, articleIds []string) chan async.ActionResult[[]domain.Article] {
	sql := `
		SELECT a.id, a.name, a.description, a.category, a.price, a.stock_amount
		FROM klcs.article a
		WHERE a.id = ANY ($1)
	`
	return db.QueryStatementTx(
		tx,
		articleMapper,
		sql,
		pq.Array(articleIds),
	)
}

// CreateArticle implements IArticleDao.
func (a *ArticleDao) CreateArticle(tx *sql.Tx, article *domain.Article, shopId string, printerId *string) chan async.ActionResult[domain.Article] {
	sql := `
		INSERT INTO klcs.article (name, description, price, category, stock_amount, shop_id, printer_id)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		RETURNING id, name, description, category, price, stock_amount
	`
	return db.QuerySingleTx(
		tx,
		articleMapper,
		sql,
		article.Name,
		article.Description,
		article.Price,
		article.Category,
		article.StockAmount,
		shopId,
		printerId,
	)
}

// DeleteArticle implements IArticleDao.
func (a *ArticleDao) DeleteArticle(tx *sql.Tx, articleId string) chan async.ActionResult[db.EffectedRows] {
	sql := `
		DELETE FROM klcs.article WHERE id = $1
	`
	return db.ExecStatementTx(
		tx,
		sql,
		articleId,
	)
}

// GetArticle implements IArticleDao.
func (a *ArticleDao) GetArticle(tx *sql.Tx, articleId string) chan async.ActionResult[domain.Article] {
	sql := `
		SELECT a.id, a.name, a.description, a.category, a.price, a.stock_amount
		FROM klcs.article a
		WHERE a.id = $1
	`
	return db.QuerySingleTx(
		tx,
		articleMapper,
		sql,
		articleId,
	)
}

// UpdateArticle implements IArticleDao.
func (a *ArticleDao) UpdateArticle(tx *sql.Tx, article *domain.Article) chan async.ActionResult[db.EffectedRows] {
	sql := `
		UPDATE klcs.article
		SET name=$1, description=$2, price=$3, category=$4, stock_amount=$5
		WHERE id = $6
	`
	return db.ExecStatementTx(
		tx,
		sql,
		article.Name,
		article.Description,
		article.Price,
		article.Category,
		article.StockAmount,
		article.Id,
	)
}

func (a *ArticleDao) GetArticlesForShop(tx *sql.Tx, shopId string) chan async.ActionResult[[]domain.Article] {
	sql := `
		SELECT a.id, a.name, a.description, a.category, a.price, a.stock_amount
		FROM klcs.article a
			JOIN klcs.shop s ON (a.shop_id = s.id)
		WHERE s.id = $1
	`
	return db.QueryStatementTx[domain.Article](
		tx,
		articleMapper,
		sql,
		shopId,
	)
}

func articleMapper() ([]any, *domain.Article) {
	v := domain.Article{}
	return []any{&v.Id, &v.Name, &v.Description, &v.Category, &v.Price, &v.StockAmount}, &v
}

func NewArticleDao() IArticleDao {
	return &ArticleDao{}
}

package dal

import (
	"context"

	"github.com/lib/pq"
	db "github.com/uoul/go-dbx"
	"github.com/uoul/klcs/backend/core/domain"
)

type ArticleDao struct{}

// CreateArticle implements [IArticleDao].
func (a *ArticleDao) CreateArticle(ctx context.Context, s db.IDbSession, article domain.Article, shopId string, printerId *string) ([]domain.Article, error) {
	return db.Query[domain.Article](
		ctx, s,
		`
			INSERT INTO klcs.article (name, description, price, category, stock_amount, shop_id, printer_id)
			VALUES ($1,$2,$3,$4,$5,$6,$7)
			RETURNING id, name, description, category, price, stock_amount
		`,
		article.Name, article.Description, article.Price, article.Category, article.StockAmount, shopId, printerId,
	)
}

// DeleteArticle implements [IArticleDao].
func (a *ArticleDao) DeleteArticle(ctx context.Context, s db.IDbSession, articleId string) error {
	_, err := db.Query[any](
		ctx, s,
		`DELETE FROM klcs.article WHERE id = $1`,
		articleId,
	)
	return err
}

// GetArticle implements [IArticleDao].
func (a *ArticleDao) GetArticle(ctx context.Context, s db.IDbSession, articleId string) ([]domain.Article, error) {
	return db.Query[domain.Article](
		ctx, s,
		`
			SELECT a.id, a.name, a.description, a.category, a.price, a.stock_amount
			FROM klcs.article a
			WHERE a.id = $1
		`,
		articleId,
	)
}

// GetArticlesForShop implements [IArticleDao].
func (a *ArticleDao) GetArticlesForShop(ctx context.Context, s db.IDbSession, shopId string) ([]domain.Article, error) {
	return db.Query[domain.Article](
		ctx, s,
		`
			SELECT a.id, a.name, a.description, a.category, a.price, a.stock_amount
			FROM klcs.article a
				JOIN klcs.shop s ON (a.shop_id = s.id)
			WHERE s.id = $1
		`,
		shopId,
	)
}

// GetArticlesIn implements [IArticleDao].
func (a *ArticleDao) GetArticlesIn(ctx context.Context, s db.IDbSession, articleIds []string) ([]domain.Article, error) {
	return db.Query[domain.Article](
		ctx, s,
		`
			SELECT a.id, a.name, a.description, a.category, a.price, a.stock_amount
			FROM klcs.article a
			WHERE a.id = ANY ($1)
		`,
		pq.Array(articleIds),
	)
}

// SetPrinterForArticle implements [IArticleDao].
func (a *ArticleDao) SetPrinterForArticle(ctx context.Context, s db.IDbSession, articleId string, printerId *string) error {
	_, err := db.Query[any](
		ctx, s,
		`
			UPDATE klcs.article
			SET printer_id=$2
			WHERE id = $1
		`,
		articleId, printerId,
	)
	return err
}

// UpdateArticle implements [IArticleDao].
func (a *ArticleDao) UpdateArticle(ctx context.Context, s db.IDbSession, article domain.Article) error {
	_, err := db.Query[any](
		ctx, s,
		`
			UPDATE klcs.article
			SET name=$1, description=$2, price=$3, category=$4, stock_amount=$5
			WHERE id = $6
		`,
		article.Name, article.Description, article.Price, article.Category, article.StockAmount, article.Id,
	)
	return err
}

func NewArticleDao() *ArticleDao {
	return &ArticleDao{}
}

package logic

import (
	"context"

	"github.com/uoul/klcs/backend/core/domain"
)

type ILogic interface {
	// KLCS-Admin
	CreateShop(ctx context.Context, username string, shop *domain.Shop) (*domain.Shop, error)
	GetShops(ctx context.Context) ([]domain.Shop, error)
	UpdateShop(ctx context.Context, shop *domain.Shop) (*domain.Shop, error)
	DeleteShop(ctx context.Context, shopId string) error

	// User registration
	UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error)

	// Seller
	GetShopsForUser(ctx context.Context, username string) ([]domain.Shop, error)
	GetShopDetailsForUser(ctx context.Context, username string, shopId string) (*domain.ShopDetails, error)
	Checkout(ctx context.Context, username string, order *domain.Order) (*domain.Order, error)
	GetAccountDetails(ctx context.Context, accountId string) (*domain.AccountDetails, error)
	GetHistory(ctx context.Context, username string, length int) ([]domain.HistoryItem, error)

	// Shop-Admin
	GetArticlesForShop(ctx context.Context, username string, shopId string) ([]domain.Article, error)
	CreateArticle(ctx context.Context, username string, shopId string, article *domain.ArticleDetails) (*domain.ArticleDetails, error)
	GetArticle(ctx context.Context, username string, articleId string) (*domain.ArticleDetails, error)
	UpdateArticle(ctx context.Context, username string, article *domain.ArticleDetails) (*domain.ArticleDetails, error)
	DeleteArticle(ctx context.Context, username string, articleId string) error

	GetPrintersForShop(ctx context.Context, username string, shopId string) ([]domain.Printer, error)
	CreatePrinter(ctx context.Context, username, shopId string, printer *domain.Printer) (*domain.Printer, error)
	DeletePrinter(ctx context.Context, username, printerId string) error

	GetUsers(ctx context.Context) ([]domain.User, error)
	GetShopUsers(ctx context.Context, username, shopId string) (map[domain.User][]domain.Role, error)
	GetRoles(ctx context.Context) ([]domain.Role, error)
	AddUserRole(ctx context.Context, username, shopId, userId, roleId string) error
	DeleteUserRole(ctx context.Context, username, shopId, userId, roleId string) error

	// Account-Manager
	GetAllAccounts(ctx context.Context) ([]domain.Account, error)
	CreateAccount(ctx context.Context, account *domain.Account) (*domain.Account, error)
	UpdateAccount(ctx context.Context, account *domain.Account) (*domain.Account, error)
	CloseAccount(ctx context.Context, username, accountId string) (*domain.AccountDetails, error)
	PostToAccount(ctx context.Context, username, accountId string, amount int) (*domain.AccountDetails, error)
}

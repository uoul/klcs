package logic

import (
	"context"
	"database/sql"
	"fmt"
	"maps"
	"slices"

	"github.com/uoul/go-common/db"
	"github.com/uoul/klcs/backend/oos-core/dal"
	"github.com/uoul/klcs/backend/oos-core/domain"
	appError "github.com/uoul/klcs/backend/oos-core/error"
)

const (
	ADMIN_ROLE = "ADMIN"
)

type Logic struct {
	cf             db.IConnectionFactory
	shopDao        dal.IShopDao
	userDao        dal.IUserDao
	articleDao     dal.IArticleDao
	roleDao        dal.IRoleDao
	printerDao     dal.IPrinterDao
	accountDao     dal.IAccountDao
	transactionDao dal.ITransactionDao
}

// DeletePrinter implements ILogic.
func (l *Logic) DeletePrinter(ctx context.Context, username string, printerId string) error {
	_, err := db.ExecInTransactionContext(
		ctx,
		l.cf,
		func(ctx context.Context, tx *sql.Tx) (any, error) {
			shop := <-l.shopDao.GetShopForPrinter(tx, printerId)
			if shop.Error != nil {
				return nil, shop.Error
			}
			err := l.checkUserRole(tx, username, shop.Result.Id, ADMIN_ROLE)
			if err != nil {
				return nil, err
			}
			r := <-l.printerDao.DeletePrinter(tx, printerId)
			return r.Result, r.Error
		},
	)
	return err
}

// CloseAccount implements ILogic.
func (l *Logic) CloseAccount(ctx context.Context, username, accountId string) (*domain.AccountDetails, error) {
	return db.ExecInTransactionContext(
		ctx,
		l.cf,
		func(ctx context.Context, tx *sql.Tx) (*domain.AccountDetails, error) {
			accountDetails, err := l.getAccountDetails(tx, accountId)
			if err != nil {
				return nil, err
			}
			if accountDetails.Balance != 0 {
				user := <-l.userDao.GetUserByUsername(tx, username)
				if user.Error != nil {
					return nil, appError.NewPermissionError(user.Error)
				}
				t := <-l.transactionDao.CreateTransaction(
					tx,
					user.Result.Id,
					&accountId,
					nil,
					&domain.Transaction{
						Type:        "CARD",
						Amount:      -accountDetails.Balance,
						Description: "Account closed",
					},
				)
				if t.Error != nil {
					return nil, t.Error
				}
			}
			return accountDetails, nil
		},
	)
}

// CreateAccount implements ILogic.
func (l *Logic) CreateAccount(ctx context.Context, account *domain.Account) (*domain.Account, error) {
	return db.ExecInTransactionContext(
		ctx,
		l.cf,
		func(ctx context.Context, tx *sql.Tx) (*domain.Account, error) {
			account := <-l.accountDao.CreateAccount(tx, account)
			return &account.Result, account.Error
		},
	)
}

// PostToAccount implements ILogic.
func (l *Logic) PostToAccount(ctx context.Context, username, accountId string, amount int) (*domain.AccountDetails, error) {
	return db.ExecInTransactionContext(
		ctx,
		l.cf,
		func(ctx context.Context, tx *sql.Tx) (*domain.AccountDetails, error) {
			user := <-l.userDao.GetUserByUsername(tx, username)
			if user.Error != nil {
				return nil, appError.NewPermissionError(user.Error)
			}
			transaction := <-l.transactionDao.CreateTransaction(
				tx,
				user.Result.Id,
				&accountId,
				nil,
				&domain.Transaction{
					Type:        "CARD",
					Amount:      amount,
					Description: "Credit top-up",
				},
			)
			if transaction.Error != nil {
				return nil, transaction.Error
			}
			return l.getAccountDetails(tx, accountId)
		},
	)
}

// UpdateAccount implements ILogic.
func (l *Logic) UpdateAccount(ctx context.Context, account *domain.Account) (*domain.Account, error) {
	return db.ExecInTransactionContext(
		ctx,
		l.cf,
		func(ctx context.Context, tx *sql.Tx) (*domain.Account, error) {
			r := <-l.accountDao.UpdateAccount(tx, account)
			return account, r.Error
		},
	)
}

// GetAccountDetails implements ILogic.
func (l *Logic) GetAccountDetails(ctx context.Context, accountId string) (*domain.AccountDetails, error) {
	return db.ExecInTransactionContext(
		ctx,
		l.cf,
		func(ctx context.Context, tx *sql.Tx) (*domain.AccountDetails, error) {
			return l.getAccountDetails(tx, accountId)
		},
	)
}

// GetRoles implements ILogic.
func (l *Logic) GetRoles(ctx context.Context) ([]domain.Role, error) {
	return db.ExecInTransactionContext(
		ctx,
		l.cf,
		func(ctx context.Context, tx *sql.Tx) ([]domain.Role, error) {
			roles := <-l.roleDao.GetRoles(tx)
			return roles.Result, roles.Error
		},
	)
}

// AssignShopAdmin implements ILogic.
func (l *Logic) AssignShopAdmin(ctx context.Context, shopId string, userId string) error {
	_, err := db.ExecInTransactionContext(
		ctx,
		l.cf,
		func(ctx context.Context, tx *sql.Tx) (*any, error) {
			role := <-l.roleDao.GetRoleByName(tx, ADMIN_ROLE)
			if role.Error != nil {
				return nil, role.Error
			}
			r := <-l.userDao.AssignUserShopRole(tx, userId, shopId, role.Result.Id)
			return nil, r.Error
		},
	)
	return err
}

// GetUsers implements ILogic.
func (l *Logic) GetUsers(ctx context.Context) ([]domain.User, error) {
	return db.ExecInTransactionContext(
		ctx,
		l.cf,
		func(ctx context.Context, tx *sql.Tx) ([]domain.User, error) {
			users := <-l.userDao.GetAll(tx)
			return users.Result, users.Error
		},
	)
}

// Checkout implements ILogic.
func (l *Logic) Checkout(ctx context.Context, username string, order *domain.Order) (*domain.Order, error) {
	return db.ExecInTransactionContext(
		ctx,
		l.cf,
		func(ctx context.Context, tx *sql.Tx) (*domain.Order, error) {
			// validate order
			err := l.validateOrder(order)
			if err != nil {
				return nil, err
			}
			// check user permissions for all articles
			err = l.checkUserPermissionsForArticles(tx, username, slices.Collect(maps.Keys(order.Articles)))
			if err != nil {
				return nil, err
			}
			// calculate order sum + update stock
			orderSum, err := l.updateStockAmountAndCalculateSumOfOrder(tx, order)
			if err != nil {
				return nil, err
			}
			// Check account for card payment
			err = l.checkAccountConditionsForCheckOutWithCard(tx, order, orderSum)
			if err != nil {
				return nil, err
			}
			// Create transaction in database
			err = l.createTransactionForCheckout(tx, username, order, orderSum)
			if err != nil {
				return nil, err
			}
			sum := float32(orderSum)
			order.Sum = &sum
			return order, nil
		},
	)
}

// AddUserRole implements ILogic.
func (l *Logic) AddUserRole(ctx context.Context, username string, shopId string, userId string, roleId string) error {
	_, err := db.ExecInTransactionContext(
		ctx,
		l.cf,
		func(ctx context.Context, tx *sql.Tx) (*any, error) {
			err := l.checkUserRole(tx, username, shopId, ADMIN_ROLE)
			if err != nil {
				return nil, err
			}
			r := <-l.userDao.AssignUserShopRole(tx, userId, shopId, roleId)
			return nil, r.Error
		},
	)
	return err
}

// DeleteUserRole implements ILogic.
func (l *Logic) DeleteUserRole(ctx context.Context, username string, shopId string, userId string, roleId string) error {
	_, err := db.ExecInTransactionContext(
		ctx,
		l.cf,
		func(ctx context.Context, tx *sql.Tx) (*any, error) {
			err := l.checkUserRole(tx, username, shopId, ADMIN_ROLE)
			if err != nil {
				return nil, err
			}
			r := <-l.userDao.UnassignUserShopRole(tx, userId, shopId, roleId)
			return nil, r.Error
		},
	)
	return err
}

// CreateArticle implements ILogic.
func (l *Logic) CreateArticle(ctx context.Context, username string, shopId string, article *domain.ArticleDetails) (*domain.ArticleDetails, error) {
	return db.ExecInTransactionContext(
		ctx,
		l.cf,
		func(ctx context.Context, tx *sql.Tx) (*domain.ArticleDetails, error) {
			err := l.checkUserRole(tx, username, shopId, ADMIN_ROLE)
			if err != nil {
				return nil, err
			}
			var printerId *string = nil
			if article.Printer != nil {
				p := <-l.printerDao.GetPrinter(tx, article.Printer.Id)
				switch p.Error {
				case nil:
				case sql.ErrNoRows:
					return nil, appError.NewValidationError(fmt.Errorf("given printer with id %s does not exist", article.Id))
				default:
					return nil, p.Error
				}
				shopForPrinter := <-l.shopDao.GetShopForPrinter(tx, article.Printer.Id)
				if shopForPrinter.Error != nil {
					return nil, shopForPrinter.Error
				}
				if shopId != shopForPrinter.Result.Id {
					return nil, appError.NewValidationError(fmt.Errorf("given printer and given article does not belong to same shop"))
				}
				printerId = &article.Printer.Id
			}
			a := <-l.articleDao.CreateArticle(tx, &domain.Article{
				Name:        article.Name,
				Description: article.Description,
				Price:       article.Price,
				Category:    article.Category,
				StockAmount: article.StockAmount,
			}, shopId, printerId)
			if a.Error != nil {
				return nil, a.Error
			}
			return &domain.ArticleDetails{
				Id:          a.Result.Id,
				Name:        a.Result.Name,
				Description: a.Result.Description,
				Price:       a.Result.Price,
				Category:    a.Result.Category,
				StockAmount: a.Result.StockAmount,
				Printer:     article.Printer,
			}, nil
		},
	)
}

// CreatePrinterForShop implements ILogic.
func (l *Logic) CreatePrinter(ctx context.Context, username string, shopId string, printer *domain.Printer) (*domain.Printer, error) {
	return db.ExecInTransactionContext(
		ctx,
		l.cf,
		func(ctx context.Context, tx *sql.Tx) (*domain.Printer, error) {
			err := l.checkUserRole(tx, username, shopId, ADMIN_ROLE)
			if err != nil {
				return nil, err
			}
			p := <-l.printerDao.CreatePrinter(tx, shopId, printer)
			return &p.Result, p.Error
		},
	)
}

// CreateShop implements ILogic.
func (l *Logic) CreateShop(ctx context.Context, username string, shop *domain.Shop) (*domain.Shop, error) {
	return db.ExecInTransactionContext(
		ctx,
		l.cf,
		func(ctx context.Context, tx *sql.Tx) (*domain.Shop, error) {
			s := <-l.shopDao.CreateShop(tx, shop)
			if s.Error != nil {
				return nil, s.Error
			}
			user := <-l.userDao.GetUserByUsername(tx, username)
			if user.Error != nil {
				return nil, user.Error
			}
			role := <-l.roleDao.GetRoleByName(tx, ADMIN_ROLE)
			if role.Error != nil {
				return nil, role.Error
			}
			r := <-l.userDao.AssignUserShopRole(tx, user.Result.Id, s.Result.Id, role.Result.Id)
			return &s.Result, r.Error
		},
	)
}

// DeleteArticle implements ILogic.
func (l *Logic) DeleteArticle(ctx context.Context, username string, articleId string) error {
	_, err := db.ExecInTransactionContext(
		ctx,
		l.cf,
		func(ctx context.Context, tx *sql.Tx) (*any, error) {
			shop := <-l.shopDao.GetShopForArticle(tx, articleId)
			if shop.Error != nil {
				return nil, shop.Error
			}
			err := l.checkUserRole(tx, username, shop.Result.Id, ADMIN_ROLE)
			if err != nil {
				return nil, err
			}
			a := <-l.articleDao.DeleteArticle(tx, articleId)
			return nil, a.Error
		},
	)
	return err
}

// DeleteShop implements ILogic.
func (l *Logic) DeleteShop(ctx context.Context, shopId string) error {
	_, err := db.ExecInTransactionContext(
		ctx,
		l.cf,
		func(ctx context.Context, tx *sql.Tx) (*any, error) {
			r := <-l.shopDao.DeleteShop(tx, shopId)
			return nil, r.Error
		},
	)
	return err
}

// GetArticle implements ILogic.
func (l *Logic) GetArticle(ctx context.Context, username string, articleId string) (*domain.ArticleDetails, error) {
	return db.ExecInTransactionContext(
		ctx,
		l.cf,
		func(ctx context.Context, tx *sql.Tx) (*domain.ArticleDetails, error) {
			shopOfArticle := <-l.shopDao.GetShopForArticle(tx, articleId)
			switch shopOfArticle.Error {
			case nil:
			case sql.ErrNoRows:
				return nil, appError.NewNotFoundError(fmt.Errorf("shop for Article not found"))
			default:
				return nil, shopOfArticle.Error
			}
			err := l.checkUserMemberOfShop(tx, username, shopOfArticle.Result.Id)
			if err != nil {
				return nil, err
			}
			var printer *domain.Printer = nil
			p := <-l.printerDao.GetPrinterForArticle(tx, articleId)
			switch p.Error {
			case nil:
				printer = &p.Result
			case sql.ErrNoRows:
			default:
				return nil, err
			}
			a := <-l.articleDao.GetArticle(tx, articleId)
			switch a.Error {
			case nil:
				return &domain.ArticleDetails{
					Id:          a.Result.Id,
					Name:        a.Result.Name,
					Description: a.Result.Description,
					Price:       a.Result.Price,
					Category:    a.Result.Category,
					StockAmount: a.Result.StockAmount,
					Printer:     printer,
				}, nil
			case sql.ErrNoRows:
				return nil, appError.NewNotFoundError(fmt.Errorf("article not found"))
			default:
				return nil, a.Error
			}
		},
	)
}

// GetArticlesForShop implements ILogic.
func (l *Logic) GetArticlesForShop(ctx context.Context, username string, shopId string) ([]domain.Article, error) {
	return db.ExecInTransactionContext(
		ctx,
		l.cf,
		func(ctx context.Context, tx *sql.Tx) ([]domain.Article, error) {
			err := l.checkUserMemberOfShop(tx, username, shopId)
			if err != nil {
				return nil, err
			}
			a := <-l.articleDao.GetArticlesForShop(tx, shopId)
			return a.Result, a.Error
		},
	)
}

// GetPrintersForShop implements ILogic.
func (l *Logic) GetPrintersForShop(ctx context.Context, username string, shopId string) ([]domain.Printer, error) {
	return db.ExecInTransactionContext(
		ctx,
		l.cf,
		func(ctx context.Context, tx *sql.Tx) ([]domain.Printer, error) {
			err := l.checkUserMemberOfShop(tx, username, shopId)
			if err != nil {
				return nil, err
			}
			p := <-l.printerDao.GetPrintersForShop(tx, shopId)
			return p.Result, p.Error
		},
	)
}

// GetShop implements ILogic.
func (l *Logic) GetShops(ctx context.Context) ([]domain.Shop, error) {
	return db.ExecInTransactionContext(
		ctx,
		l.cf,
		func(ctx context.Context, tx *sql.Tx) ([]domain.Shop, error) {
			r := <-l.shopDao.GetAll(tx)
			return r.Result, r.Error
		},
	)
}

// GetShopUsers implements ILogic.
func (l *Logic) GetShopUsers(ctx context.Context, username string, shopId string) (map[domain.User][]domain.Role, error) {
	return db.ExecInTransactionContext(
		ctx,
		l.cf,
		func(ctx context.Context, tx *sql.Tx) (map[domain.User][]domain.Role, error) {
			err := l.checkUserRole(tx, username, shopId, ADMIN_ROLE)
			if err != nil {
				return nil, err
			}
			users := <-l.userDao.GetAll(tx)
			if users.Error != nil {
				return nil, users.Error
			}
			userMapping := make(map[domain.User][]domain.Role)
			for _, user := range users.Result {
				l := <-l.roleDao.GetUserRolesForShop(tx, user.Username, shopId)
				if l.Error != nil {
					return nil, l.Error
				}
				userMapping[user] = l.Result
			}
			return userMapping, nil
		},
	)
}

// UpdateArticle implements ILogic.
func (l *Logic) UpdateArticle(ctx context.Context, username string, article *domain.ArticleDetails) (*domain.ArticleDetails, error) {
	return db.ExecInTransactionContext(
		ctx,
		l.cf,
		func(ctx context.Context, tx *sql.Tx) (*domain.ArticleDetails, error) {
			shopForArticle := <-l.shopDao.GetShopForArticle(tx, article.Id)
			if shopForArticle.Error != nil {
				return nil, shopForArticle.Error
			}
			err := l.checkUserRole(tx, username, shopForArticle.Result.Id, ADMIN_ROLE)
			if err != nil {
				return nil, err
			}
			a := <-l.articleDao.UpdateArticle(tx, &domain.Article{
				Id:          article.Id,
				Name:        article.Name,
				Description: article.Description,
				Price:       article.Price,
				Category:    article.Category,
				StockAmount: article.StockAmount,
			})
			if a.Error != nil {
				return nil, a.Error
			}
			if article.Printer != nil {
				p := <-l.printerDao.GetPrinter(tx, article.Printer.Id)
				switch p.Error {
				case nil:
				case sql.ErrNoRows:
					return nil, appError.NewValidationError(fmt.Errorf("given printer with id %s does not exist", article.Id))
				default:
					return nil, p.Error
				}
				shopForPrinter := <-l.shopDao.GetShopForPrinter(tx, article.Printer.Id)
				if shopForPrinter.Error != nil {
					return nil, shopForPrinter.Error
				}
				if shopForArticle.Result.Id != shopForPrinter.Result.Id {
					return nil, appError.NewValidationError(fmt.Errorf("given printer and given article does not belong to same shop"))
				}
				r := <-l.articleDao.SetPrinterForArticle(tx, article.Id, &article.Printer.Id)
				if r.Error != nil {
					return nil, r.Error
				}
			} else {
				r := <-l.articleDao.SetPrinterForArticle(tx, article.Id, nil)
				if r.Error != nil {
					return nil, r.Error
				}
			}
			return article, a.Error
		},
	)
}

// UpdateShop implements ILogic.
func (l *Logic) UpdateShop(ctx context.Context, shop *domain.Shop) (*domain.Shop, error) {
	return db.ExecInTransactionContext(
		ctx,
		l.cf,
		func(ctx context.Context, tx *sql.Tx) (*domain.Shop, error) {
			r := <-l.shopDao.UpdateShop(tx, shop)
			return shop, r.Error
		},
	)
}

// UpdateUser implements ILogic.
func (l *Logic) UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	return db.ExecInTransactionContext(
		ctx,
		l.cf,
		func(ctx context.Context, tx *sql.Tx) (*domain.User, error) {
			u := <-l.userDao.CreateOrUpdateUser(tx, user)
			return &u.Result, u.Error
		},
	)
}

// GetShopDetailsForUser implements ILogic.
func (l *Logic) GetShopDetailsForUser(ctx context.Context, username string, shopId string) (*domain.ShopDetails, error) {
	return db.ExecInTransactionContext(
		ctx,
		l.cf,
		func(ctx context.Context, tx *sql.Tx) (*domain.ShopDetails, error) {
			shop := <-l.shopDao.GetShop(tx, shopId)
			if shop.Error != nil {
				return nil, shop.Error
			}
			userRoles := <-l.roleDao.GetUserRolesForShop(tx, username, shopId)
			if userRoles.Error != nil {
				return nil, userRoles.Error
			}
			if len(userRoles.Result) <= 0 {
				return nil, appError.NewPermissionError(fmt.Errorf("%s is no member of %s", username, shop.Result.Name))
			}
			articles := <-l.articleDao.GetArticlesForShop(tx, shopId)
			if articles.Error != nil {
				return nil, articles.Error
			}
			shopDetails := &domain.ShopDetails{
				Id:         shop.Result.Id,
				Name:       shop.Result.Name,
				UserRoles:  convertUserRoles(userRoles.Result),
				Categories: convertArticles(articles.Result),
			}
			return shopDetails, nil
		},
	)
}

// GetShopsForUser implements ILogic.
func (l *Logic) GetShopsForUser(ctx context.Context, username string) ([]domain.Shop, error) {
	return db.ExecInTransactionContext(
		ctx,
		l.cf,
		func(ctx context.Context, tx *sql.Tx) ([]domain.Shop, error) {
			shops := <-l.shopDao.GetShopsForUser(tx, username)
			return shops.Result, shops.Error
		},
	)
}

func convertUserRoles(roles []domain.Role) []string {
	r := []string{}
	for _, role := range roles {
		r = append(r, role.Name)
	}
	return r
}

func convertArticles(articles []domain.Article) map[string][]domain.Article {
	r := make(map[string][]domain.Article)
	for _, article := range articles {
		if _, exists := r[article.Category]; !exists {
			r[article.Category] = []domain.Article{}
		}
		r[article.Category] = append(r[article.Category], article)
	}
	return r
}

func (l *Logic) checkUserRole(tx *sql.Tx, username string, shopId string, role string) error {
	userRoles := <-l.roleDao.GetUserRolesForShop(tx, username, shopId)
	if userRoles.Error != nil {
		return userRoles.Error
	}
	if len(userRoles.Result) <= 0 {
		return appError.NewPermissionError(fmt.Errorf("user %s is no member of shop %s", username, shopId))
	}
	for _, r := range userRoles.Result {
		if r.Name == role {
			return nil
		}
	}
	return appError.NewPermissionError(fmt.Errorf("user %s does not have role %s at shop %s", username, role, shopId))
}

func (l *Logic) checkUserMemberOfShop(tx *sql.Tx, username string, shopId string) error {
	userRoles := <-l.roleDao.GetUserRolesForShop(tx, username, shopId)
	if userRoles.Error != nil {
		return userRoles.Error
	}
	if len(userRoles.Result) <= 0 {
		return appError.NewPermissionError(fmt.Errorf("user %s is no member of shop %s", username, shopId))
	}
	return nil
}

func (l *Logic) getCurrentStockForArticles(tx *sql.Tx, articleIds []string) (map[string]domain.Article, error) {
	stockList := <-l.articleDao.GetArticlesIn(tx, articleIds)
	if stockList.Error != nil {
		return nil, stockList.Error
	}
	stock := make(map[string]domain.Article)
	for _, article := range stockList.Result {
		stock[article.Id] = article
	}
	return stock, nil
}

func (l *Logic) updateStockAmountAndCalculateSumOfOrder(tx *sql.Tx, order *domain.Order) (int, error) {
	articleIds := slices.Collect(maps.Keys(order.Articles))
	stock, err := l.getCurrentStockForArticles(tx, articleIds)
	if err != nil {
		return 0, err
	}
	// Check StockAmount
	ordersum := 0
	for articleId, amount := range order.Articles {
		if stock[articleId].StockAmount != nil && *stock[articleId].StockAmount < amount {
			return 0, appError.NewValidationError(fmt.Errorf("current StockAmount of article %s to low for order - need: %v, current: %v", articleId, amount, stock[articleId].StockAmount))
		}
		ordersum += stock[articleId].Price * amount
		// update StockAmount
		if stock[articleId].StockAmount != nil {
			*stock[articleId].StockAmount -= amount
			a := stock[articleId]
			r := <-l.articleDao.UpdateArticle(tx, &a)
			if r.Error != nil {
				return 0, r.Error
			}
		}
	}
	return ordersum, nil
}

func (l *Logic) validateOrder(order *domain.Order) error {
	if order.Type != "CARD" && order.Type != "CASH" {
		return appError.NewValidationError(fmt.Errorf("invalid order type %s", order.Type))
	}
	return nil
}

func (l *Logic) checkAccountConditionsForCheckOutWithCard(tx *sql.Tx, order *domain.Order, sumOfOrder int) error {
	if order.Type == "CARD" {
		accountId := *order.AccountId
		account := <-l.accountDao.GetAccount(tx, accountId)
		if account.Error != nil {
			return account.Error
		}
		if account.Result.Locked {
			return appError.NewValidationError(fmt.Errorf("account %s is currently locked", *order.AccountId))
		}
		accountBalance := <-l.transactionDao.GetAccountBalance(tx, order.AccountId)
		if accountBalance.Error != nil {
			return accountBalance.Error
		}
		if accountBalance.Result < sumOfOrder {
			return appError.NewValidationError(fmt.Errorf("account %s does not have neccessary balance - need: %v current: %v", *order.AccountId, sumOfOrder, accountBalance.Result))
		}
	}
	return nil
}

func (l *Logic) createTransactionForCheckout(tx *sql.Tx, username string, order *domain.Order, sumOfOrder int) error {
	user := <-l.userDao.GetUserByUsername(tx, username)
	if user.Error != nil {
		return user.Error
	}
	transaction := <-l.transactionDao.CreateTransaction(
		tx,
		user.Result.Id,
		order.AccountId,
		order.Articles,
		&domain.Transaction{
			Type:        order.Type,
			Amount:      -sumOfOrder,
			Description: order.Description,
		},
	)
	if transaction.Error != nil {
		return transaction.Error
	}
	return nil
}

func (l *Logic) checkUserPermissionsForArticles(tx *sql.Tx, username string, articleIds []string) error {
	shops := <-l.shopDao.GetShopsForArticles(tx, articleIds)
	if shops.Error != nil {
		return shops.Error
	}
	for _, shop := range shops.Result {
		err := l.checkUserMemberOfShop(tx, username, shop.Id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *Logic) getAccountDetails(tx *sql.Tx, accountId string) (*domain.AccountDetails, error) {
	account := <-l.accountDao.GetAccount(tx, accountId)
	switch account.Error {
	case nil:
		balance := <-l.transactionDao.GetAccountBalance(tx, &accountId)
		if balance.Error != nil {
			return nil, balance.Error
		}
		return &domain.AccountDetails{
			Id:         account.Result.Id,
			HolderName: account.Result.HolderName,
			Locked:     account.Result.Locked,
			Balance:    balance.Result,
		}, nil
	case sql.ErrNoRows:
		return nil, appError.NewNotFoundError(fmt.Errorf("account with id %s does not exist", accountId))
	default:
		return nil, account.Error
	}
}

func NewLogic(cf db.IConnectionFactory) ILogic {
	return &Logic{
		cf: cf,

		shopDao:        dal.NewShopDao(),
		userDao:        dal.NewUserDao(),
		articleDao:     dal.NewArticleDao(),
		roleDao:        dal.NewRoleDao(),
		printerDao:     dal.NewPrinterDao(),
		accountDao:     dal.NewAccountDao(),
		transactionDao: dal.NewTransactionDao(),
	}
}

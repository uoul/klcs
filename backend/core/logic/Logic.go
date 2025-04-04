package logic

import (
	"context"
	"database/sql"
	"maps"
	"slices"

	"github.com/uoul/go-common/db"
	"github.com/uoul/go-common/log"
	"github.com/uoul/klcs/backend/core/dal"
	"github.com/uoul/klcs/backend/core/domain"
	appError "github.com/uoul/klcs/backend/core/error"
	"github.com/uoul/klcs/backend/core/services"
)

const (
	ADMIN_ROLE = "ADMIN"
)

type Logic struct {
	cf             db.IConnectionFactory
	logger         log.ILogger
	printService   *services.PrintService
	shopDao        dal.IShopDao
	userDao        dal.IUserDao
	articleDao     dal.IArticleDao
	roleDao        dal.IRoleDao
	printerDao     dal.IPrinterDao
	accountDao     dal.IAccountDao
	transactionDao dal.ITransactionDao
	historyDao     dal.IHistoryDao
	printJobDao    dal.IPrintJobDao
}

// GetAccountsByExternalId implements ILogic.
func (l *Logic) GetAccountsByExternalId(ctx context.Context, externalId string) ([]domain.Account, error) {
	return db.ExecInTransactionContext(
		ctx,
		l.cf,
		func(ctx context.Context, tx *sql.Tx) ([]domain.Account, error) {
			accounts := <-l.accountDao.GetAccountsByExternalId(tx, externalId)
			if accounts.Error != nil {
				return nil, appError.NewErrDataAccess("failed to obtain accounts for given externalId(%s) - %v", externalId, accounts.Error)
			}
			return accounts.Result, nil
		},
	)
}

// GetHistory implements ILogic.
func (l *Logic) GetHistory(ctx context.Context, username string, length int) ([]domain.HistoryItem, error) {
	return db.ExecInTransactionContext(
		ctx,
		l.cf,
		func(ctx context.Context, tx *sql.Tx) ([]domain.HistoryItem, error) {
			history := <-l.historyDao.GetHistoryForUser(tx, username, length)
			if history.Error != nil {
				return nil, appError.NewErrDataAccess("failed to get history for user(%s) - %v", username, history.Error)
			}
			return history.Result, nil
		},
	)
}

// GetAllAccounts implements ILogic.
func (l *Logic) GetAllAccounts(ctx context.Context) ([]domain.Account, error) {
	return db.ExecInTransactionContext(
		ctx,
		l.cf,
		func(ctx context.Context, tx *sql.Tx) ([]domain.Account, error) {
			accounts := <-l.accountDao.GetAll(tx)
			if accounts.Error != nil {
				return nil, appError.NewErrDataAccess("failed to get accouts - %v", accounts.Error)
			}
			return accounts.Result, nil
		},
	)
}

// DeletePrinter implements ILogic.
func (l *Logic) DeletePrinter(ctx context.Context, username string, printerId string) error {
	_, err := db.ExecInTransactionContext(
		ctx,
		l.cf,
		func(ctx context.Context, tx *sql.Tx) (any, error) {
			shop := <-l.shopDao.GetShopForPrinter(tx, printerId)
			if shop.Error != nil {
				if shop.Error == sql.ErrNoRows {
					return nil, appError.NewErrNotFound("shop for printer(%s) not found - %v", printerId, shop.Error)
				}
				return nil, appError.NewErrDataAccess("failed to get shop for printer(%s) - %v", printerId, shop.Error)
			}
			err := l.checkUserRole(tx, username, shop.Result.Id, ADMIN_ROLE)
			if err != nil {
				return nil, err
			}
			r := <-l.printerDao.DeletePrinter(tx, printerId)
			if r.Error != nil {
				if r.Error == sql.ErrNoRows {
					return nil, appError.NewErrNotFound("printer(%s) not found - %v", printerId, r.Error)
				}
				return nil, appError.NewErrDataAccess("failed to delete printer(%s) - %v", printerId, r.Error)
			}
			return r.Result, nil
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
				if err == sql.ErrNoRows {
					return nil, appError.NewErrNotFound("account(%s) not found - %v", accountId, err)
				}
				return nil, appError.NewErrDataAccess("failed to get balance for account(%s) - %v", accountId, err)
			}
			if accountDetails.Locked {
				return nil, appError.NewErrValidation("cannot close account(%s), that is locked", accountDetails.Id)
			}
			if accountDetails.Balance != 0 {
				user := <-l.userDao.GetUserByUsername(tx, username)
				if user.Error != nil {
					if user.Error == sql.ErrNoRows {
						return nil, appError.NewErrNotFound("user(%s) not found - %v", username, user.Error)
					}
					return nil, appError.NewErrDataAccess("failed to get user(%s) - %v", username, user.Error)
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
					return nil, appError.NewErrDataAccess("failed to create new transaction - %v", t.Error)
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
			a := <-l.accountDao.CreateAccount(tx, account)
			if a.Error != nil {
				return nil, appError.NewErrDataAccess("failed to create account(%s) - %v", account.HolderName, a.Error)
			}
			return &a.Result, nil
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
				return nil, appError.NewErrDataAccess("failed to get user(%s) - %v", username, user.Error)
			}
			account := <-l.accountDao.GetAccount(tx, accountId)
			if account.Error != nil {
				if account.Error == sql.ErrNoRows {
					return nil, appError.NewErrNotFound("accout(%s) not found - %v", accountId, account.Error)
				}
				return nil, appError.NewErrDataAccess("failed to get account(%s) - %v", accountId, account.Error)
			}
			if account.Result.Locked {
				return nil, appError.NewErrValidation("cannot charge locked account(%s)", accountId)
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
				return nil, appError.NewErrDataAccess("failed to create transaction - %v", transaction.Error)
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
			if r.Error != nil {
				if r.Error == sql.ErrNoRows {
					return nil, appError.NewErrNotFound("account(%s) not found - %v", account.Id, r.Error)
				}
				return nil, appError.NewErrDataAccess("failed to update account(%s) - %v", account.Id, r.Error)
			}
			return account, nil
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
			if roles.Error != nil {
				return nil, appError.NewErrDataAccess("failed to get roles - %v", roles.Error)
			}
			return roles.Result, nil
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
				if role.Error == sql.ErrNoRows {
					return nil, appError.NewErrNotFound("role(%s) not found - %v", ADMIN_ROLE, role.Error)
				}
				return nil, appError.NewErrDataAccess("failed to get role(%s) - %v", ADMIN_ROLE, role.Error)
			}
			r := <-l.userDao.AssignUserShopRole(tx, userId, shopId, role.Result.Id)
			if r.Error != nil {
				return nil, appError.NewErrDataAccess("failed to assign role(%s) for user(%s) to shop(%s) - %v", ADMIN_ROLE, userId, shopId, r.Error)
			}
			return nil, nil
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
			if users.Error != nil {
				return nil, appError.NewErrDataAccess("failed to get users - %v", users.Error)
			}
			return users.Result, nil
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
			transaction, err := l.createTransactionForCheckout(tx, username, order, orderSum)
			if err != nil {
				return nil, err
			}
			// Generate PrintJobs for order
			printJobs := <-l.printJobDao.GetPrintOpenJobsForTransaction(tx, transaction.Id)
			if printJobs.Error != nil {
				return nil, appError.NewErrDataAccess("failed to get printjobs or transaction - %v", printJobs.Error)
			}
			// Print
			for printerId, job := range printJobs.Result {
				err := l.printService.PrintJob(printerId, job)
				if err != nil {
					l.logger.Warningf("failed to send printjob to printer(%s) - %v", printerId, err)
				}
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
			if r.Error != nil {
				return nil, appError.NewErrDataAccess("failed to add userrole - %v", r.Error)
			}
			return nil, nil
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
			if r.Error != nil {
				return nil, appError.NewErrDataAccess("failed to unassign userrole - %v", r.Error)
			}
			return nil, nil
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
					return nil, appError.NewErrNotFound("given printer with id %s does not exist", article.Id)
				default:
					return nil, appError.NewErrDataAccess("failed to get printer(%s) - %v", article.Printer.Id, p.Error)
				}
				shopForPrinter := <-l.shopDao.GetShopForPrinter(tx, article.Printer.Id)
				if shopForPrinter.Error != nil {
					return nil, appError.NewErrDataAccess("failed to get shop for printer(%v) - %v", printerId, shopForPrinter.Error)
				}
				if shopId != shopForPrinter.Result.Id {
					return nil, appError.NewErrValidation("given printer and given article does not belong to same shop")
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
				return nil, appError.NewErrDataAccess("failed to create article - %v", a.Error)
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
			if p.Error != nil {
				return nil, appError.NewErrDataAccess("failed to create printer - %v", p.Error)
			}
			return &p.Result, nil
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
				return nil, appError.NewErrDataAccess("failed to create shop(%s) - %v", shop.Name, s.Error)
			}
			user := <-l.userDao.GetUserByUsername(tx, username)
			if user.Error != nil {
				return nil, appError.NewErrDataAccess("failed to get user(%s) - %v", username, user.Error)
			}
			role := <-l.roleDao.GetRoleByName(tx, ADMIN_ROLE)
			if role.Error != nil {
				return nil, appError.NewErrDataAccess("failed to get role(%s) - %v", ADMIN_ROLE, role.Error)
			}
			r := <-l.userDao.AssignUserShopRole(tx, user.Result.Id, s.Result.Id, role.Result.Id)
			if r.Error != nil {
				return nil, appError.NewErrDataAccess("failed to assign role(%s) to user(%s) for shop(%s) - %v", ADMIN_ROLE, username, shop.Name, r.Error)
			}
			return &s.Result, nil
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
				return nil, appError.NewErrDataAccess("failed to get shop for article(%s), - %v", articleId, shop.Error)
			}
			err := l.checkUserRole(tx, username, shop.Result.Id, ADMIN_ROLE)
			if err != nil {
				return nil, err
			}
			a := <-l.articleDao.DeleteArticle(tx, articleId)
			if a.Error != nil {
				return nil, appError.NewErrDataAccess("failed to delete article - %v", a.Error)
			}
			return nil, nil
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
			if r.Error != nil {
				return nil, appError.NewErrDataAccess("failed to delete shop(%s) - %v", shopId, r.Error)
			}
			return nil, nil
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
			if shopOfArticle.Error != nil {
				if shopOfArticle.Error == sql.ErrNoRows {
					return nil, appError.NewErrNotFound("shop for article(%s) not found - %v", articleId, shopOfArticle.Error)
				}
				return nil, appError.NewErrDataAccess("failed to get shop for article(%s) - %v", articleId, shopOfArticle.Error)
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
				return nil, appError.NewErrDataAccess("failed to get printer for article(%s) - %v", articleId, p.Error)
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
				return nil, appError.NewErrNotFound("article(%s) not found - %v", articleId, a.Error)
			default:
				return nil, appError.NewErrDataAccess("failed to get article(%s) - %v", articleId, a.Error)
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
			if a.Error != nil {
				if a.Error == sql.ErrNoRows {
					return nil, appError.NewErrNotFound("shop(%s) not found - %v", shopId, a.Error)
				}
				return nil, appError.NewErrDataAccess("failed to get articles for shop(%s) - %v", shopId, a.Error)
			}
			return a.Result, nil
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
			if p.Error != nil {
				if p.Error == sql.ErrNoRows {
					return nil, appError.NewErrNotFound("shop(%s) not found - %v", shopId, p.Error)
				}
				return nil, appError.NewErrDataAccess("failed to get printers for shop(%s) - %v", shopId, p.Error)
			}
			return p.Result, nil
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
			if r.Error != nil {
				return nil, appError.NewErrDataAccess("failed to get shops - %v", r.Error)
			}
			return r.Result, nil
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
				return nil, appError.NewErrDataAccess("failed to get users - %v", users.Error)
			}
			userMapping := make(map[domain.User][]domain.Role)
			for _, user := range users.Result {
				l := <-l.roleDao.GetUserRolesForShop(tx, user.Username, shopId)
				if l.Error != nil {
					return nil, appError.NewErrDataAccess("failed to get roles of user(%s) for shop(%s) - %v", username, shopId, l.Error)
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
				return nil, appError.NewErrDataAccess("failed to get shop for article(%s) - %v", article.Name, shopForArticle.Error)
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
				return nil, appError.NewErrDataAccess("failed to update article(%s) - %v", article.Name, a.Error)
			}
			if article.Printer != nil {
				p := <-l.printerDao.GetPrinter(tx, article.Printer.Id)
				if p.Error != nil {
					if p.Error == sql.ErrNoRows {
						return nil, appError.NewErrNotFound("printer(%s) does not exist - %v", article.Printer.Id, p.Error)
					}
					return nil, appError.NewErrDataAccess("failed to get printer(%s) - %v", article.Printer.Id, p.Error)
				}
				shopForPrinter := <-l.shopDao.GetShopForPrinter(tx, article.Printer.Id)
				if shopForPrinter.Error != nil {
					return nil, appError.NewErrDataAccess("failed to get shop for printer(%s) - %v", article.Printer.Id, shopForPrinter.Error)
				}
				if shopForArticle.Result.Id != shopForPrinter.Result.Id {
					return nil, appError.NewErrValidation("given printer(%s) and given article(%s) does not belong to same shop", article.Printer.Id, article.Id)
				}
				r := <-l.articleDao.SetPrinterForArticle(tx, article.Id, &article.Printer.Id)
				if r.Error != nil {
					return nil, appError.NewErrDataAccess("failed to set printer(%s) for article(%s) - %v", article.Printer.Id, article.Id, r.Error)
				}
			} else {
				r := <-l.articleDao.SetPrinterForArticle(tx, article.Id, nil)
				if r.Error != nil {
					return nil, appError.NewErrDataAccess("failed to remove Printer from article(%s) - %v", article.Id, r.Error)
				}
			}
			return article, nil
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
			if r.Error != nil {
				return nil, appError.NewErrDataAccess("failed to update shop(%s)", shop.Id)
			}
			return shop, nil
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
			if u.Error != nil {
				return nil, appError.NewErrDataAccess("failed to update user - %v", u.Error)
			}
			return &u.Result, nil
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
				if shop.Error == sql.ErrNoRows {
					return nil, appError.NewErrNotFound("shop(%s) does not exist - %v", shopId, shop.Error)
				}
				return nil, appError.NewErrDataAccess("failed to get shop(%s) - %v", shopId, shop.Error)
			}
			userRoles := <-l.roleDao.GetUserRolesForShop(tx, username, shopId)
			if userRoles.Error != nil {
				return nil, appError.NewErrDataAccess("failed to get roles of user(%s) for shop(%s) - %v", username, shopId, userRoles.Error)
			}
			if len(userRoles.Result) <= 0 {
				return nil, appError.NewErrForbidden("%s is no member of %s", username, shop.Result.Name)
			}
			articles := <-l.articleDao.GetArticlesForShop(tx, shopId)
			if articles.Error != nil {
				return nil, appError.NewErrDataAccess("failed to get articles for shop(%s) - %v", shopId, articles.Error)
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
			if shops.Error != nil {
				return nil, appError.NewErrDataAccess("failed to get shops for user(%s) - %v", username, shops.Error)
			}
			return shops.Result, nil
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
		return appError.NewErrDataAccess("failed to get shop roles for user - %v", userRoles.Error)
	}
	if len(userRoles.Result) <= 0 {
		return appError.NewErrForbidden("user %s is no member of shop %s", username, shopId)
	}
	for _, r := range userRoles.Result {
		if r.Name == role {
			return nil
		}
	}
	return appError.NewErrForbidden("user %s does not have role %s at shop %s", username, role, shopId)
}

func (l *Logic) checkUserMemberOfShop(tx *sql.Tx, username string, shopId string) error {
	userRoles := <-l.roleDao.GetUserRolesForShop(tx, username, shopId)
	if userRoles.Error != nil {
		return appError.NewErrDataAccess("failed to get shop roles for user - %v", userRoles.Error)
	}
	if len(userRoles.Result) <= 0 {
		return appError.NewErrForbidden("user %s is no member of shop %s", username, shopId)
	}
	return nil
}

func (l *Logic) getCurrentStockForArticles(tx *sql.Tx, articleIds []string) (map[string]domain.Article, error) {
	stockList := <-l.articleDao.GetArticlesIn(tx, articleIds)
	if stockList.Error != nil {
		return nil, appError.NewErrDataAccess("failed to get articles - %v", stockList.Error)
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
			return 0, appError.NewErrValidation("current StockAmount of article %s to low for order - need: %v, current: %v", articleId, amount, *stock[articleId].StockAmount)
		}
		ordersum += stock[articleId].Price * amount
		// update StockAmount
		if stock[articleId].StockAmount != nil {
			*stock[articleId].StockAmount -= amount
			a := stock[articleId]
			r := <-l.articleDao.UpdateArticle(tx, &a)
			if r.Error != nil {
				return 0, appError.NewErrDataAccess("failed to update stock amount for article(%s) - %v", articleId, r.Error)
			}
		}
	}
	return ordersum, nil
}

func (l *Logic) validateOrder(order *domain.Order) error {
	if order.Type != "CARD" && order.Type != "CASH" {
		return appError.NewErrValidation("invalid order type %s", order.Type)
	}
	return nil
}

func (l *Logic) checkAccountConditionsForCheckOutWithCard(tx *sql.Tx, order *domain.Order, sumOfOrder int) error {
	if order.Type == "CARD" {
		accountId := *order.AccountId
		account := <-l.accountDao.GetAccount(tx, accountId)
		if account.Error != nil {
			if account.Error == sql.ErrNoRows {
				return appError.NewErrNotFound("account(%s) not found - %v", accountId, account.Error)
			}
			return appError.NewErrDataAccess("failed to get account(%s) - %v", accountId, account.Error)
		}
		if account.Result.Locked {
			return appError.NewErrValidation("account %s is currently locked", *order.AccountId)
		}
		accountBalance := <-l.transactionDao.GetAccountBalance(tx, order.AccountId)
		if accountBalance.Error != nil {
			return appError.NewErrDataAccess("failed to get account(%s) balance - %v", accountId, accountBalance.Error)
		}
		if accountBalance.Result < sumOfOrder {
			return appError.NewErrValidation("account %s does not have neccessary balance - need: %v current: %v", *order.AccountId, sumOfOrder, accountBalance.Result)
		}
	}
	return nil
}

func (l *Logic) createTransactionForCheckout(tx *sql.Tx, username string, order *domain.Order, sumOfOrder int) (*domain.Transaction, error) {
	user := <-l.userDao.GetUserByUsername(tx, username)
	if user.Error != nil {
		if user.Error == sql.ErrNoRows {
			return nil, appError.NewErrNotFound("user(%s) not found - %v", username, user.Error)
		}
		return nil, appError.NewErrDataAccess("failed to get user(%s) - %v", username, user.Error)
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
		return nil, appError.NewErrDataAccess("failed to create transaction - %v", transaction.Error)
	}
	return &transaction.Result, nil
}

func (l *Logic) checkUserPermissionsForArticles(tx *sql.Tx, username string, articleIds []string) error {
	shops := <-l.shopDao.GetShopsForArticles(tx, articleIds)
	if shops.Error != nil {
		return appError.NewErrDataAccess("failed to get shops for articles(%s) - %v", articleIds, shops.Error)
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
			return nil, appError.NewErrDataAccess("failed to get balance for account(%s) - %v", accountId, balance.Error)
		}
		return &domain.AccountDetails{
			Id:         account.Result.Id,
			HolderName: account.Result.HolderName,
			Locked:     account.Result.Locked,
			ExternalId: account.Result.ExternalId,
			Balance:    balance.Result,
		}, nil
	case sql.ErrNoRows:
		return nil, appError.NewErrNotFound("account with id %s does not exist", accountId)
	default:
		return nil, appError.NewErrDataAccess("failed to get account(%s) - %v", accountId, account.Error)
	}
}

func NewLogic(cf db.IConnectionFactory, logger log.ILogger, printService *services.PrintService) ILogic {
	return &Logic{
		cf: cf,

		printService: printService,
		logger:       logger,

		shopDao:        dal.NewShopDao(),
		userDao:        dal.NewUserDao(),
		articleDao:     dal.NewArticleDao(),
		roleDao:        dal.NewRoleDao(),
		printerDao:     dal.NewPrinterDao(),
		accountDao:     dal.NewAccountDao(),
		transactionDao: dal.NewTransactionDao(),
		historyDao:     dal.NewHistoryDao(),
		printJobDao:    dal.NewPrintJobDao(),
	}
}

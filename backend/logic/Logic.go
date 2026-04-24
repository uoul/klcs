package logic

import (
	"context"
	"database/sql"
	"log/slog"
	"maps"
	"slices"

	db "github.com/uoul/go-dbx"
	"github.com/uoul/klcs/backend/core/apperror"
	"github.com/uoul/klcs/backend/core/dal"
	"github.com/uoul/klcs/backend/core/domain"
	"github.com/uoul/klcs/backend/core/services"
)

const (
	SHOP_ADMIN_ROLE = "ADMIN"
)

type Logic struct {
	dbConn       db.IDbConnection
	printService *services.PrintService

	shopDao        *dal.ShopDao
	userDao        *dal.UserDao
	articleDao     *dal.ArticleDao
	roleDao        *dal.RoleDao
	printerDao     *dal.PrinterDao
	accountDao     *dal.AccountDao
	transactionDao *dal.TransactionDao
	historyDao     *dal.HistoryDao
	printJobDao    *dal.PrintJobDao
}

// ReprintOpenTransaction implements ILogic.
func (l *Logic) Reprint(ctx context.Context, transactionId string) error {
	// Get print jobs
	jobs, err := l.printJobDao.GetPrintOpenJobsForTransaction(ctx, l.dbConn, transactionId)
	if err != nil {
		return apperror.NewErrDataAccess(err, "failed to get open printjobs for transaction(%s)", transactionId)
	}
	// Print Jobs
	for printerId, job := range jobs {
		if e := l.printService.PrintJob(printerId, job); e != nil {
			// Log Error but keep trying other jobs
			slog.Warn("failed to print job", slog.String("printer", printerId), slog.String("transaction", transactionId), slog.Any("error", e))
			err = e
		}
	}
	// Return last error
	return err
}

// AcknowledgePrintJob implements ILogic.
func (l *Logic) AcknowledgePrintJob(ctx context.Context, printerId string, transactionId string) error {
	if err := l.printJobDao.AcknowledgeByTransactionId(ctx, l.dbConn, printerId, transactionId); err != nil {
		return apperror.NewErrDataAccess(err, "failed to store acknowledgement for transaction(%s) and printer(%s)", transactionId, printerId)
	}
	return nil
}

// GetAccountsByExternalId implements ILogic.
func (l *Logic) GetAccountsByExternalId(ctx context.Context, externalId string) ([]domain.Account, error) {
	accounts, err := l.accountDao.GetAccountsByExternalId(ctx, l.dbConn, externalId)
	if err != nil {
		return nil, apperror.NewErrDataAccess(err, "failed to get accounts for externalId(%s)", externalId)
	}
	return accounts, nil
}

// GetHistory implements ILogic.
func (l *Logic) GetHistory(ctx context.Context, username string, length int) ([]domain.HistoryItem, error) {
	return db.ExecuteInTransaction(
		ctx, l.dbConn,
		func(ctx context.Context, tx *sql.Tx) ([]domain.HistoryItem, error) {
			// Check if user exists
			if _, err := l.getUser(ctx, tx, username); err != nil {
				return nil, err
			}
			// Get History
			history, err := l.historyDao.GetHistoryForUser(ctx, l.dbConn, username, length)
			if err != nil {
				return nil, apperror.NewErrDataAccess(err, "failed to get history for user(%s)", username)
			}
			return history, nil
		},
	)
}

// GetAllAccounts implements ILogic.
func (l *Logic) GetAllAccounts(ctx context.Context) ([]domain.Account, error) {
	accounts, err := l.accountDao.GetAll(ctx, l.dbConn)
	if err != nil {
		return nil, apperror.NewErrDataAccess(err, "failed to get accounts")
	}
	return accounts, nil
}

// DeletePrinter implements ILogic.
func (l *Logic) DeletePrinter(ctx context.Context, username string, printerId string) error {
	_, err := db.ExecuteInTransaction(
		ctx, l.dbConn,
		func(ctx context.Context, tx *sql.Tx) (any, error) {
			// Get shop for printer
			shops, err := l.shopDao.GetShopForPrinter(ctx, tx, printerId)
			if err != nil {
				return nil, apperror.NewErrDataAccess(err, "failed to get shop for printer(%s)", printerId)
			}
			if len(shops) <= 0 {
				return nil, apperror.NewErrNoPrinterShopRelation(nil, "printer(%s) does not belong to a shop", printerId)
			}
			// Check if user has priveledges on shop
			if err := l.checkUserRole(ctx, tx, username, shops[0].Id, SHOP_ADMIN_ROLE); err != nil {
				return nil, err
			}
			// Delete printer
			if err := l.printerDao.DeletePrinter(ctx, tx, printerId); err != nil {
				return nil, apperror.NewErrDataAccess(err, "failed to delete printer(%s)", printerId)
			}
			return nil, nil
		},
	)
	return err
}

// CloseAccount implements ILogic.
func (l *Logic) CloseAccount(ctx context.Context, username, accountId string) (domain.AccountDetails, error) {
	return db.ExecuteInTransaction(
		ctx, l.dbConn,
		func(ctx context.Context, tx *sql.Tx) (domain.AccountDetails, error) {
			// Get account details
			accounts, err := l.accountDao.GetAccount(ctx, tx, accountId)
			if err != nil {
				return domain.AccountDetails{}, apperror.NewErrDataAccess(err, "failed to get account(%s)", accountId)
			}
			if len(accounts) <= 0 {
				return domain.AccountDetails{}, apperror.NewErrAccountNotFound(nil, "account(%s) not found", accountId)
			}
			if accounts[0].Locked {
				return domain.AccountDetails{}, apperror.NewErrCloseLockedAccount(nil, "cannot close locked account(%s)", accountId)
			}
			// Get account balance
			balance, err := l.transactionDao.GetAccountBalance(ctx, tx, accountId)
			if err != nil || len(balance) <= 0 {
				return domain.AccountDetails{}, apperror.NewErrDataAccess(err, "failed to get account balance for %s", accountId)
			}
			// Create inverse transaction
			if balance[0] != 0 {
				// Get user by username
				user, err := l.getUser(ctx, tx, username)
				if err != nil {
					return domain.AccountDetails{}, err
				}
				// Create transaction
				transactions, err := l.transactionDao.CreateTransaction(ctx, tx, user.Id, &accountId, nil, domain.Transaction{
					Type:        "CARD",
					Amount:      -balance[0],
					Description: "Account closed",
				}, true)
				if err != nil || len(transactions) <= 0 {
					return domain.AccountDetails{}, apperror.NewErrDataAccess(err, "failed to create inverse transaction for closing account(%s)", accountId)
				}
			}
			return domain.AccountDetails{
				Account: accounts[0],
				Balance: balance[0],
			}, nil
		},
	)
}

// CreateAccount implements ILogic.
func (l *Logic) CreateAccount(ctx context.Context, account domain.Account) (domain.Account, error) {
	a, err := l.accountDao.CreateAccount(ctx, l.dbConn, account)
	if err != nil || len(a) <= 0 {
		return domain.Account{}, apperror.NewErrDataAccess(err, "failed to create account(%s)", account.HolderName)
	}
	return a[0], nil
}

// PostToAccount implements ILogic.
func (l *Logic) PostToAccount(ctx context.Context, username, accountId string, amount int) (domain.AccountDetails, error) {
	return db.ExecuteInTransaction(
		ctx, l.dbConn,
		func(ctx context.Context, tx *sql.Tx) (domain.AccountDetails, error) {
			// Get user, that initiates the post
			user, err := l.getUser(ctx, tx, username)
			if err != nil {
				return domain.AccountDetails{}, err
			}
			// Get account
			accounts, err := l.accountDao.GetAccount(ctx, tx, accountId)
			if err != nil {
				return domain.AccountDetails{}, apperror.NewErrDataAccess(err, "failed to get account(%s)", accountId)
			}
			if len(accounts) <= 0 {
				return domain.AccountDetails{}, apperror.NewErrAccountNotFound(nil, "account(%s) not found", accountId)
			}
			// Check if account is locked
			if accounts[0].Locked {
				return domain.AccountDetails{}, apperror.NewErrPostLockedAccount(nil, "cannot charge locked account(%s)", accountId)
			}
			// Create charging transaction
			transactions, err := l.transactionDao.CreateTransaction(ctx, tx, user.Id, &accountId, nil, domain.Transaction{
				Type:        "CARD",
				Amount:      amount,
				Description: "Credit top-up",
			}, true)
			if err != nil || len(transactions) <= 0 {
				return domain.AccountDetails{}, apperror.NewErrDataAccess(err, "failed to create transaction")
			}
			// Get Updated account balance
			balance, err := l.transactionDao.GetAccountBalance(ctx, tx, accountId)
			if err != nil || len(balance) <= 0 {
				return domain.AccountDetails{}, apperror.NewErrDataAccess(err, "failed to get account(%s) balance", accountId)
			}
			// Return updated account
			return domain.AccountDetails{
				Account: accounts[0],
				Balance: balance[0],
			}, nil
		},
	)
}

// UpdateAccount implements ILogic.
func (l *Logic) UpdateAccount(ctx context.Context, account domain.Account) (domain.Account, error) {
	return db.ExecuteInTransaction(
		ctx, l.dbConn,
		func(ctx context.Context, tx *sql.Tx) (domain.Account, error) {
			// Check if account exists
			accounts, err := l.accountDao.GetAccount(ctx, tx, account.Id)
			if err != nil {
				return domain.Account{}, apperror.NewErrDataAccess(err, "failed to get account(%s)", account.Id)
			}
			if len(accounts) <= 0 {
				return domain.Account{}, apperror.NewErrAccountNotFound(nil, "account(%s) not found", account.Id)
			}
			// Update Account
			if err := l.accountDao.UpdateAccount(ctx, l.dbConn, account); err != nil {
				return domain.Account{}, apperror.NewErrDataAccess(err, "failed to update account")
			}
			// Success
			return account, nil
		},
	)
}

// GetAccountDetails implements ILogic.
func (l *Logic) GetAccountDetails(ctx context.Context, accountId string) (domain.AccountDetails, error) {
	return db.ExecuteInTransaction(
		ctx, l.dbConn,
		func(ctx context.Context, tx *sql.Tx) (domain.AccountDetails, error) {
			// Get account
			accounts, err := l.accountDao.GetAccount(ctx, tx, accountId)
			if err != nil {
				return domain.AccountDetails{}, apperror.NewErrDataAccess(err, "failed to get account(%s)", accountId)
			}
			if len(accounts) <= 0 {
				return domain.AccountDetails{}, apperror.NewErrAccountNotFound(nil, "account(%s) not found", accountId)
			}
			// Get account balance
			balance, err := l.transactionDao.GetAccountBalance(ctx, tx, accountId)
			if err != nil || len(balance) <= 0 {
				return domain.AccountDetails{}, apperror.NewErrDataAccess(err, "failed to get account(%s) balance", accountId)
			}
			// Return account details
			return domain.AccountDetails{
				Account: accounts[0],
				Balance: balance[0],
			}, nil
		},
	)
}

// GetRoles implements ILogic.
func (l *Logic) GetRoles(ctx context.Context) ([]domain.Role, error) {
	roles, err := l.roleDao.GetRoles(ctx, l.dbConn)
	if err != nil {
		return nil, apperror.NewErrDataAccess(err, "failed to get roles")
	}
	return roles, nil
}

// AssignShopAdmin implements ILogic.
func (l *Logic) AssignShopAdmin(ctx context.Context, shopId string, userId string) error {
	_, err := db.ExecuteInTransaction(
		ctx, l.dbConn,
		func(ctx context.Context, tx *sql.Tx) (any, error) {
			// Get Admin role
			roles, err := l.roleDao.GetRoleByName(ctx, tx, SHOP_ADMIN_ROLE)
			if err != nil || len(roles) <= 0 {
				return nil, apperror.NewErrDataAccess(err, "failed to get admin role")
			}
			// Assign role to user
			if err := l.userDao.AssignUserShopRole(ctx, tx, userId, shopId, roles[0].Id); err != nil {
				return nil, apperror.NewErrDataAccess(err, "failed to assign admin role(%s) to user(%s)", roles[0].Id, userId)
			}
			return nil, nil
		},
	)
	return err
}

// GetUsers implements ILogic.
func (l *Logic) GetUsers(ctx context.Context) ([]domain.User, error) {
	users, err := l.userDao.GetAll(ctx, l.dbConn)
	if err != nil {
		return nil, apperror.NewErrDataAccess(err, "failed to get users")
	}
	return users, nil
}

// Checkout implements ILogic.
func (l *Logic) Checkout(ctx context.Context, username string, printDisabled bool, order domain.Order) (domain.Order, error) {
	return db.ExecuteInTransaction(
		ctx, l.dbConn,
		func(ctx context.Context, tx *sql.Tx) (domain.Order, error) {
			// Validate Order
			if err := l.validateOrder(order); err != nil {
				return domain.Order{}, nil
			}
			// Check user permissions for all articles
			if err := l.checkUserPermissionsForArticles(ctx, tx, username, slices.Collect(maps.Keys(order.Articles))); err != nil {
				return domain.Order{}, err
			}
			// Calculate order sum + update stock
			orderSum, err := l.updateStockAmountAndCalculateSumOfOrder(ctx, tx, order)
			if err != nil {
				return domain.Order{}, err
			}
			// Check Account for Card Payment
			if err := l.checkAccountConditionsForCheckOutWithCard(ctx, tx, order, orderSum); err != nil {
				return domain.Order{}, err
			}
			// Create Transaction
			transaction, err := l.createTransactionForCheckout(ctx, tx, username, printDisabled, order, orderSum)
			if err != nil {
				return domain.Order{}, err
			}
			// Generate PrintJobs for order
			if !printDisabled {
				printJobs, err := l.printJobDao.GetPrintOpenJobsForTransaction(ctx, tx, transaction.Id)
				if err != nil {
					return domain.Order{}, apperror.NewErrDataAccess(err, "failed to get printjobs for transaction(%s)", transaction.Id)
				}
				// Print
				for printerId, job := range printJobs {
					err := l.printService.PrintJob(printerId, job)
					if err != nil {
						slog.Warn("failed to send printjob", slog.String("printer", printerId), slog.Any("error", err))
					}
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
	_, err := db.ExecuteInTransaction(
		ctx, l.dbConn,
		func(ctx context.Context, tx *sql.Tx) (any, error) {
			// Check if requesting user has admin role on shop
			if err := l.checkUserRole(ctx, tx, username, shopId, SHOP_ADMIN_ROLE); err != nil {
				return nil, err
			}
			// Assign user role
			if err := l.userDao.AssignUserShopRole(ctx, tx, userId, shopId, roleId); err != nil {
				return nil, apperror.NewErrDataAccess(err, "failed to assign user(%s) role(%s) for shop(%s)", userId, roleId, shopId)
			}
			return nil, nil
		},
	)
	return err
}

// DeleteUserRole implements ILogic.
func (l *Logic) DeleteUserRole(ctx context.Context, username string, shopId string, userId string, roleId string) error {
	_, err := db.ExecuteInTransaction(
		ctx, l.dbConn,
		func(ctx context.Context, tx *sql.Tx) (any, error) {
			// Check if requesting user has admin role on shop
			if err := l.checkUserRole(ctx, tx, username, shopId, SHOP_ADMIN_ROLE); err != nil {
				return nil, err
			}
			// Assign user role
			if err := l.userDao.UnassignUserShopRole(ctx, tx, userId, shopId, roleId); err != nil {
				return nil, apperror.NewErrDataAccess(err, "failed to unassign user(%s) role(%s) for shop(%s)", userId, roleId, shopId)
			}
			return nil, nil
		},
	)
	return err
}

// CreateArticle implements ILogic.
func (l *Logic) CreateArticle(ctx context.Context, username string, shopId string, article domain.ArticleDetails) (domain.ArticleDetails, error) {
	return db.ExecuteInTransaction(
		ctx, l.dbConn,
		func(ctx context.Context, tx *sql.Tx) (domain.ArticleDetails, error) {
			// Check user permissions on shop
			if err := l.checkUserRole(ctx, tx, username, shopId, SHOP_ADMIN_ROLE); err != nil {
				return domain.ArticleDetails{}, err
			}
			// Get printer if is set
			var printerId *string
			if article.Printer != nil {
				// Get Printer
				printers, err := l.printerDao.GetPrinter(ctx, tx, article.Printer.Id)
				if err != nil {
					return domain.ArticleDetails{}, apperror.NewErrDataAccess(err, "failed to get printer(%s)", article.Printer.Id)
				}
				if len(printers) <= 0 {
					return domain.ArticleDetails{}, apperror.NewErrPrinterNotFound(nil, "printer(%s) not found", article.Printer.Id)
				}
				// Check if printer corresponds to shop
				shops, err := l.shopDao.GetShopForPrinter(ctx, tx, printers[0].Id)
				if err != nil {
					return domain.ArticleDetails{}, apperror.NewErrDataAccess(err, "failed to get shop for printer(%s)", printers[0].Id)
				}
				if len(shops) <= 0 {
					return domain.ArticleDetails{}, apperror.NewErrNoPrinterShopRelation(nil, "shop for printer(%s) not found", printers[0].Id)
				}
				if shopId != shops[0].Id {
					return domain.ArticleDetails{}, apperror.NewErrPrinterAndArticleNotSameShop(nil, "given printer(%s) and given article(%s) does not belong to same shop", article.Printer.Id, article.Id)
				}
				printerId = &article.Printer.Id
			}
			// Create article
			articles, err := l.articleDao.CreateArticle(ctx, tx, article.Article, shopId, printerId)
			if err != nil || len(articles) <= 0 {
				return domain.ArticleDetails{}, apperror.NewErrDataAccess(err, "failed to create article")
			}
			// Return result
			return domain.ArticleDetails{
				Article: articles[0],
				Printer: article.Printer,
			}, nil
		},
	)
}

// CreatePrinterForShop implements ILogic.
func (l *Logic) CreatePrinter(ctx context.Context, username string, shopId string, printer domain.Printer) (domain.Printer, error) {
	return db.ExecuteInTransaction(
		ctx, l.dbConn,
		func(ctx context.Context, tx *sql.Tx) (domain.Printer, error) {
			// Check user permissions on shop
			if err := l.checkUserRole(ctx, tx, username, shopId, SHOP_ADMIN_ROLE); err != nil {
				return domain.Printer{}, err
			}
			// Create printer
			printers, err := l.printerDao.CreatePrinter(ctx, tx, shopId, printer)
			if err != nil || len(printers) <= 0 {
				return domain.Printer{}, apperror.NewErrDataAccess(err, "failed to create printer")
			}
			// Return new printer
			return printers[0], nil
		},
	)
}

// CreateShop implements ILogic.
func (l *Logic) CreateShop(ctx context.Context, username string, shop domain.Shop) (domain.Shop, error) {
	return db.ExecuteInTransaction(
		ctx, l.dbConn,
		func(ctx context.Context, tx *sql.Tx) (domain.Shop, error) {
			// Create shop
			shops, err := l.shopDao.CreateShop(ctx, tx, shop)
			if err != nil || len(shops) <= 0 {
				return domain.Shop{}, apperror.NewErrDataAccess(err, "failed to create shop")
			}
			// Get current user
			users, err := l.userDao.GetUserByUsername(ctx, tx, username)
			if err != nil || len(users) <= 0 {
				return domain.Shop{}, apperror.NewErrUserForShopCreationNotFound(err, "failed to get current user(%s)", username)
			}
			// Get admin role
			roles, err := l.roleDao.GetRoleByName(ctx, tx, SHOP_ADMIN_ROLE)
			if err != nil || len(roles) <= 0 {
				return domain.Shop{}, apperror.NewErrDataAccess(err, "failed to get admin role")
			}
			// Assign shop admin role to current user
			if err := l.userDao.AssignUserShopRole(ctx, tx, users[0].Id, shops[0].Id, roles[0].Id); err != nil {
				return domain.Shop{}, apperror.NewErrDataAccess(err, "failed to assign admin role(%s) for shop(%s) to user(%s)", roles[0].Id, shops[0].Id, users[0].Id)
			}
			// Return new shop
			return shops[0], nil
		},
	)
}

// DeleteArticle implements ILogic.
func (l *Logic) DeleteArticle(ctx context.Context, username string, articleId string) error {
	_, err := db.ExecuteInTransaction(
		ctx, l.dbConn,
		func(ctx context.Context, tx *sql.Tx) (any, error) {
			// Get shop for article
			shops, err := l.shopDao.GetShopForArticle(ctx, tx, articleId)
			if err != nil {
				return nil, apperror.NewErrDataAccess(err, "failed to get shop for article(%s)", articleId)
			}
			if len(shops) <= 0 {
				return nil, apperror.NewErrArticleShopRelation(nil, "article(%s) does not belong to any shop", articleId)
			}
			// Check user has permission
			if err := l.checkUserRole(ctx, tx, username, shops[0].Id, SHOP_ADMIN_ROLE); err != nil {
				return nil, err
			}
			// Delete article
			if err := l.articleDao.DeleteArticle(ctx, tx, articleId); err != nil {
				return nil, apperror.NewErrDataAccess(err, "failed to delete article(%s)", articleId)
			}
			// Success
			return nil, nil
		},
	)
	return err
}

// DeleteShop implements ILogic.
func (l *Logic) DeleteShop(ctx context.Context, shopId string) error {
	if err := l.shopDao.DeleteShop(ctx, l.dbConn, shopId); err != nil {
		return apperror.NewErrDataAccess(err, "failed to delete shop(%s)", shopId)
	}
	return nil
}

// GetArticle implements ILogic.
func (l *Logic) GetArticle(ctx context.Context, username string, articleId string) (domain.ArticleDetails, error) {
	return db.ExecuteInTransaction(
		ctx, l.dbConn,
		func(ctx context.Context, tx *sql.Tx) (domain.ArticleDetails, error) {
			// Get shop for article
			shops, err := l.shopDao.GetShopForArticle(ctx, tx, articleId)
			if err != nil {
				return domain.ArticleDetails{}, apperror.NewErrDataAccess(err, "failed to get shop for article(%s)", articleId)
			}
			if len(shops) <= 0 {
				return domain.ArticleDetails{}, apperror.NewErrArticleShopRelation(nil, "shop for article(%s) not found", articleId)
			}
			// Check if user belongs to shop
			if err := l.checkUserMemberOfShop(ctx, tx, username, shops[0].Id); err != nil {
				return domain.ArticleDetails{}, err
			}
			// Get Printer
			printers, err := l.printerDao.GetPrinterForArticle(ctx, tx, articleId)
			if err != nil {
				return domain.ArticleDetails{}, apperror.NewErrDataAccess(err, "failed to get printer for article(%s)", articleId)
			}
			var printer *domain.Printer = nil
			if len(printers) > 0 {
				printer = &printers[0]
			}
			// Get Article
			articles, err := l.articleDao.GetArticle(ctx, tx, articleId)
			if err != nil {
				return domain.ArticleDetails{}, apperror.NewErrDataAccess(err, "failed to get article(%s)", articleId)
			}
			if len(articles) <= 0 {
				return domain.ArticleDetails{}, apperror.NewErrArtilceNotFound(nil, "article(%s) not found", articleId)
			}
			// Return ArticleDetails
			return domain.ArticleDetails{
				Article: articles[0],
				Printer: printer,
			}, nil
		},
	)
}

// GetArticlesForShop implements ILogic.
func (l *Logic) GetArticlesForShop(ctx context.Context, username string, shopId string) ([]domain.Article, error) {
	return db.ExecuteInTransaction(
		ctx, l.dbConn,
		func(ctx context.Context, tx *sql.Tx) ([]domain.Article, error) {
			// Check user belongs to shop
			if err := l.checkUserMemberOfShop(ctx, tx, username, shopId); err != nil {
				return nil, err
			}
			// Get Articles
			articles, err := l.articleDao.GetArticlesForShop(ctx, tx, shopId)
			if err != nil {
				return nil, apperror.NewErrDataAccess(err, "failed to get articles for shop(%s) - %v", shopId, err)
			}
			// Return Articles
			return articles, nil
		},
	)
}

// GetPrintersForShop implements ILogic.
func (l *Logic) GetPrintersForShop(ctx context.Context, username string, shopId string) ([]domain.Printer, error) {
	return db.ExecuteInTransaction(
		ctx, l.dbConn,
		func(ctx context.Context, tx *sql.Tx) ([]domain.Printer, error) {
			// Check user belongs to shop
			if err := l.checkUserMemberOfShop(ctx, tx, username, shopId); err != nil {
				return nil, err
			}
			// Get Printers for shop
			printers, err := l.printerDao.GetPrintersForShop(ctx, tx, shopId)
			if err != nil {
				return nil, apperror.NewErrDataAccess(err, "failed to get printers for shop(%s)", shopId)
			}
			// Get Currently connected printers
			connected := l.printService.GetConnectedPrinters()
			for i, p := range printers {
				if slices.ContainsFunc(connected, func(connectedPrinterId string) bool { return connectedPrinterId == p.Id }) {
					printers[i].Connected = true
				}
			}
			// Return Printers
			return printers, nil
		},
	)
}

// GetShop implements ILogic.
func (l *Logic) GetShops(ctx context.Context) ([]domain.Shop, error) {
	shops, err := l.shopDao.GetAll(ctx, l.dbConn)
	if err != nil {
		return nil, apperror.NewErrDataAccess(err, "failed to get shops")
	}
	return shops, nil
}

// GetShopUsers implements ILogic.
func (l *Logic) GetShopUsers(ctx context.Context, username string, shopId string) (map[domain.User][]domain.Role, error) {
	return db.ExecuteInTransaction(
		ctx, l.dbConn,
		func(ctx context.Context, tx *sql.Tx) (map[domain.User][]domain.Role, error) {
			// Check user permissions
			if err := l.checkUserRole(ctx, tx, username, shopId, SHOP_ADMIN_ROLE); err != nil {
				return nil, err
			}
			// Get users
			users, err := l.userDao.GetAll(ctx, tx)
			if err != nil {
				return nil, apperror.NewErrDataAccess(err, "failed to get users")
			}
			// Create userMapping
			userMapping := make(map[domain.User][]domain.Role)
			for _, user := range users {
				roles, err := l.roleDao.GetUserRolesForShop(ctx, tx, user.Username, shopId)
				if err != nil {
					return nil, apperror.NewErrDataAccess(err, "failed to get userroles for %s on shop(%s)", user.Username, shopId)
				}
				userMapping[user] = roles
			}
			// Return Result
			return userMapping, nil
		},
	)
}

// UpdateArticle implements ILogic.
func (l *Logic) UpdateArticle(ctx context.Context, username string, article domain.ArticleDetails) (domain.ArticleDetails, error) {
	return db.ExecuteInTransaction(
		ctx, l.dbConn,
		func(ctx context.Context, tx *sql.Tx) (domain.ArticleDetails, error) {
			// Get Shop for Article
			shopForArticle, err := l.shopDao.GetShopForArticle(ctx, tx, article.Id)
			if err != nil {
				return domain.ArticleDetails{}, apperror.NewErrDataAccess(err, "failed to get shop for article(%s)", article.Id)
			}
			if len(shopForArticle) <= 0 {
				return domain.ArticleDetails{}, apperror.NewErrArticleShopRelation(nil, "shop for article(%s) not found", article.Id)
			}
			// Check user permissions
			if err := l.checkUserRole(ctx, tx, username, shopForArticle[0].Id, SHOP_ADMIN_ROLE); err != nil {
				return domain.ArticleDetails{}, err
			}
			// Update Article
			if err := l.articleDao.UpdateArticle(ctx, tx, article.Article); err != nil {
				return domain.ArticleDetails{}, apperror.NewErrDataAccess(err, "failed to update article(%s)", article.Id)
			}
			// Update Printer
			if article.Printer != nil {
				// Get Printer
				printers, err := l.printerDao.GetPrinter(ctx, tx, article.Printer.Id)
				if err != nil {
					return domain.ArticleDetails{}, apperror.NewErrDataAccess(err, "failed to get printer(%s)", article.Printer.Id)
				}
				if len(printers) <= 0 {
					return domain.ArticleDetails{}, apperror.NewErrPrinterNotFound(nil, "printer(%s) not found", article.Printer.Id)
				}
				// Get shop for printer
				shopForPrinter, err := l.shopDao.GetShopForPrinter(ctx, tx, printers[0].Id)
				if err != nil {
					return domain.ArticleDetails{}, apperror.NewErrDataAccess(err, "failed to get shop for printer(%s)", printers[0].Id)
				}
				if len(shopForPrinter) <= 0 {
					return domain.ArticleDetails{}, apperror.NewErrNoPrinterShopRelation(nil, "shop for printer(%s) not found", printers[0].Id)
				}
				// Check if Article and Printer belongs to the same shop
				if shopForPrinter[0].Id != shopForArticle[0].Id {
					return domain.ArticleDetails{}, apperror.NewErrArticlePrinterRelation(nil, "printer(%s) and article(%s) does not belong to same shop", printers[0].Id, article.Id)
				}
				// Set Printer for Article
				if err := l.articleDao.SetPrinterForArticle(ctx, tx, article.Id, &article.Printer.Id); err != nil {
					return domain.ArticleDetails{}, apperror.NewErrDataAccess(err, "failed to set printer(%s) for article(%s)", article.Printer.Id, article.Id)
				}
			} else {
				if err := l.articleDao.SetPrinterForArticle(ctx, tx, article.Id, nil); err != nil {
					return domain.ArticleDetails{}, apperror.NewErrDataAccess(err, "failed to remove printer from article(%s)", article.Id)
				}
			}
			return article, nil
		},
	)
}

// UpdateShop implements ILogic.
func (l *Logic) UpdateShop(ctx context.Context, shop domain.Shop) (domain.Shop, error) {
	err := l.shopDao.UpdateShop(ctx, l.dbConn, shop)
	if err != nil {
		return domain.Shop{}, apperror.NewErrDataAccess(err, "failed to update shop(%s)", shop.Id)
	}
	return shop, nil
}

// UpdateUser implements ILogic.
func (l *Logic) UpdateUser(ctx context.Context, user domain.User) (domain.User, error) {
	users, err := l.userDao.CreateOrUpdateUser(ctx, l.dbConn, user)
	if err != nil || len(users) <= 0 {
		return domain.User{}, apperror.NewErrDataAccess(err, "failed to create or update user(%s)", user.Username)
	}
	return users[0], nil
}

// GetShopDetailsForUser implements ILogic.
func (l *Logic) GetShopDetailsForUser(ctx context.Context, username string, shopId string) (domain.ShopDetails, error) {
	return db.ExecuteInTransaction(
		ctx, l.dbConn,
		func(ctx context.Context, tx *sql.Tx) (domain.ShopDetails, error) {
			// Get Shop
			shops, err := l.shopDao.GetShop(ctx, tx, shopId)
			if err != nil {
				return domain.ShopDetails{}, apperror.NewErrDataAccess(err, "failed to get shop(%s)", shopId)
			}
			if len(shops) <= 0 {
				return domain.ShopDetails{}, apperror.NewErrShopNotFound(nil, "shop(%s) not found", shopId)
			}
			// Get UserRoles for Shop
			userRoles, err := l.roleDao.GetUserRolesForShop(ctx, tx, username, shopId)
			if err != nil {
				return domain.ShopDetails{}, apperror.NewErrDataAccess(err, "failed to get roles for user(%s) on shop(%s)", username, shopId)
			}
			// Check if user member of shop
			if len(userRoles) <= 0 {
				return domain.ShopDetails{}, apperror.NewErrUserNoMemberOfShop(nil, "%s is no member of shop(%s)", username, shopId)
			}
			// Get Articles for Shop
			articles, err := l.articleDao.GetArticlesForShop(ctx, tx, shopId)
			if err != nil {
				return domain.ShopDetails{}, apperror.NewErrDataAccess(err, "failed to get articles for shop(%s)", shopId)
			}
			// Return ShopDetails
			return domain.ShopDetails{
				Shop:       shops[0],
				UserRoles:  convertUserRoles(userRoles),
				Categories: convertArticles(articles),
			}, nil
		},
	)
}

// GetShopsForUser implements ILogic.
func (l *Logic) GetShopsForUser(ctx context.Context, username string) ([]domain.Shop, error) {
	shops, err := l.shopDao.GetShopsForUser(ctx, l.dbConn, username)
	if err != nil {
		return nil, apperror.NewErrDataAccess(err, "failed to get shops for user(%s)", username)
	}
	return shops, nil
}

// -----------------------------------------------------------------------------------------------------------
// Helper functions
// -----------------------------------------------------------------------------------------------------------

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

func (l *Logic) checkUserRole(ctx context.Context, tx *sql.Tx, username string, shopId string, role string) error {
	// Get User Roles for Shop
	roles, err := l.roleDao.GetUserRolesForShop(ctx, tx, username, shopId)
	if err != nil {
		return apperror.NewErrDataAccess(err, "failed to get shop roles for user")
	}
	// Check if User belongs to shop
	if len(roles) <= 0 {
		return apperror.NewErrUserNoMemberOfShop(nil, "user %s is no member of shop %s", username, shopId)
	}
	// Check if user has Role on Shop
	for _, r := range roles {
		if r.Name == role {
			return nil
		}
	}
	return apperror.NewErrUserShopPermission(nil, "user %s does not have role %s at shop %s", username, role, shopId)
}

func (l *Logic) checkUserMemberOfShop(ctx context.Context, tx *sql.Tx, username string, shopId string) error {
	// Get roles for user on shop
	userRoles, err := l.roleDao.GetUserRolesForShop(ctx, tx, username, shopId)
	if err != nil {
		return apperror.NewErrDataAccess(err, "failed to get roles for %s on shop(%s)", username, shopId)
	}
	// Check if user has at least one role
	if len(userRoles) <= 0 {
		return apperror.NewErrUserNoMemberOfShop(nil, "user %s is no member of shop %s", username, shopId)
	}
	return nil
}

func (l *Logic) getCurrentStockForArticles(ctx context.Context, tx *sql.Tx, articleIds []string) (map[string]domain.Article, error) {
	// Get Articles
	stockList, err := l.articleDao.GetArticlesIn(ctx, tx, articleIds)
	if err != nil {
		return nil, apperror.NewErrDataAccess(err, "failed to get articles")
	}
	// Convert Article list to map
	stock := make(map[string]domain.Article)
	for _, article := range stockList {
		stock[article.Id] = article
	}
	return stock, nil
}

func (l *Logic) updateStockAmountAndCalculateSumOfOrder(ctx context.Context, tx *sql.Tx, order domain.Order) (int, error) {
	articleIds := slices.Collect(maps.Keys(order.Articles))
	stock, err := l.getCurrentStockForArticles(ctx, tx, articleIds)
	if err != nil {
		return 0, err
	}
	// Check StockAmount
	ordersum := 0
	for articleId, amount := range order.Articles {
		if stock[articleId].StockAmount != nil && *stock[articleId].StockAmount < amount {
			return 0, apperror.NewErrStockAmount(nil, "current StockAmount of article %s to low for order - need: %v, have: %v", articleId, amount, *stock[articleId].StockAmount)
		}
		ordersum += stock[articleId].Price * amount
		// update StockAmount
		if stock[articleId].StockAmount != nil {
			*stock[articleId].StockAmount -= amount
			a := stock[articleId]
			if err := l.articleDao.UpdateArticle(ctx, tx, a); err != nil {
				return 0, apperror.NewErrDataAccess(err, "failed to update stock amount for article(%s)", articleId)
			}
		}
	}
	return ordersum, nil
}

func (l *Logic) validateOrder(order domain.Order) error {
	if order.Type != "CARD" && order.Type != "CASH" {
		return apperror.NewErrPaymentTypeNotSupported(nil, "invalid order type %s", order.Type)
	}
	return nil
}

func (l *Logic) checkAccountConditionsForCheckOutWithCard(ctx context.Context, tx *sql.Tx, order domain.Order, sumOfOrder int) error {
	// Only check card orders
	if order.Type == "CARD" {
		// Check if AccountId given
		if order.AccountId == nil {
			return apperror.NewErrNoAccountIdForCartPayment(nil, "Orders using CARD as payment has to provide accountId")
		}
		// Get Account
		accounts, err := l.accountDao.GetAccount(ctx, tx, *order.AccountId)
		if err != nil {
			return apperror.NewErrDataAccess(err, "failed to get account(%s)", *order.AccountId)
		}
		if len(accounts) <= 0 {
			return apperror.NewErrAccountNotFound(nil, "account(%s) not found", *order.AccountId)
		}
		// Check if account is locked
		if accounts[0].Locked {
			return apperror.NewErrCheckoutLockedAccount(nil, "cannot checkout using account %s that is currently locked", *order.AccountId)
		}
		// Check account balance
		accountBalance, err := l.transactionDao.GetAccountBalance(ctx, tx, *order.AccountId)
		if err != nil || len(accountBalance) <= 0 {
			return apperror.NewErrDataAccess(nil, "failed to get account(%s) balance", *order.AccountId)
		}
		if accountBalance[0] < sumOfOrder {
			return apperror.NewErrNotEnoughBalance(nil, "account %s does not have neccessary balance - need: %v have: %v", *order.AccountId, sumOfOrder, accountBalance[0])
		}
	}
	return nil
}

func (l *Logic) createTransactionForCheckout(ctx context.Context, tx *sql.Tx, username string, printDisabled bool, order domain.Order, sumOfOrder int) (domain.Transaction, error) {
	// Get User
	user, err := l.getUser(ctx, tx, username)
	if err != nil {
		return domain.Transaction{}, err
	}
	// Create Transaction
	transactions, err := l.transactionDao.CreateTransaction(ctx, tx, user.Id, order.AccountId, order.Articles, domain.Transaction{
		Type:        order.Type,
		Amount:      -sumOfOrder,
		Description: order.Description,
	}, printDisabled)
	if err != nil || len(transactions) <= 0 {
		return domain.Transaction{}, apperror.NewErrDataAccess(err, "failed to create transaction")
	}
	return transactions[0], nil
}

func (l *Logic) checkUserPermissionsForArticles(ctx context.Context, tx *sql.Tx, username string, articleIds []string) error {
	// Get Shops for Articles
	shops, err := l.shopDao.GetShopsForArticles(ctx, tx, articleIds)
	if err != nil {
		return apperror.NewErrDataAccess(err, "failed to get shops for articles(%s)", articleIds)
	}
	// Check if user is member of all Shops for given Articles
	for _, shop := range shops {
		err := l.checkUserMemberOfShop(ctx, tx, username, shop.Id)
		if err != nil {
			return err
		}
	}
	return nil
}

func (l *Logic) getUser(ctx context.Context, s db.IDbSession, username string) (domain.User, error) {
	users, err := l.userDao.GetUserByUsername(ctx, s, username)
	if err != nil {
		return domain.User{}, apperror.NewErrDataAccess(err, "failed to get user %s", username)
	}
	if len(users) <= 0 {
		return domain.User{}, apperror.NewErrUserNotFound(nil, "user %s not found", username)
	}
	return users[0], nil
}

func NewLogic(dbConn db.IDbConnection, printService *services.PrintService) *Logic {
	return &Logic{
		dbConn:       dbConn,
		printService: printService,

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

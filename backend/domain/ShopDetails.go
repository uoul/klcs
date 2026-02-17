package domain

type ShopDetails struct {
	Shop
	UserRoles  []string
	Categories map[string][]Article
}

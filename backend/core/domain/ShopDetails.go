package domain

type ShopDetails struct {
	Id         string
	Name       string
	UserRoles  []string
	Categories map[string][]Article
}

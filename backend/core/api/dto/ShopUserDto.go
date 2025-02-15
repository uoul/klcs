package dto

import "github.com/uoul/klcs/backend/core/domain"

type ShopUserDto struct {
	Id        string
	Username  string
	Name      string
	ShopRoles []domain.Role
}

package domain

type Account struct {
	Id         string
	HolderName string
	Locked     bool
	ExternalId *string
}

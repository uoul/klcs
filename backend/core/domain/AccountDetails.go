package domain

type AccountDetails struct {
	Id         string
	HolderName string
	Locked     bool
	ExternalId *string
	Balance    int
}

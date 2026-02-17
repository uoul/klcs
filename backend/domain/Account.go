package domain

type Account struct {
	Id         string
	HolderName string `db:"holder_name"`
	Locked     bool
	ExternalId *string `db:"external_id"`
}

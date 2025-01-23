package domain

type Order struct {
	Type        string
	Description string
	AccountId   *string  `json:"AccountId,omitempty"`
	Sum         *float32 `json:"Sum,omitempty"`
	Articles    map[string]int
}

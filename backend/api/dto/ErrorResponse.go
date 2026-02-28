package dto

type ErrorResponse struct {
	Code         string
	DebugMessage string `json:"DebugMessage,omitempty"`
}

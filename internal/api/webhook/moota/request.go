package moota

type CreateWebhookMootaRequest struct {
	Account_number int64  `json:"account_number"`
	Date           string `json:"date"`
	Description    string `json:"description"`
	Amount         int64  `json:"amount"`
	Type           string `json:"type"`
	Balance        int64  `json:"balance"`
	Updated_at     string `json:"updated_at"`
	Created_at     string `json:"created_at"`
	Mutation_id    string `json:"mutation_id"`
	Token          string `json:"token"`
	Bank_id        string `json:"bank_id"`
}

type IdWebhookMootaRequest struct {
	Id string `json:"id" validate:"required"`
}

type UpdateWebhookMootaRequest struct {
	IdWebhookMootaRequest
	CreateWebhookMootaRequest
}

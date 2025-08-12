package transaction

type CreateAdminTransactionRequest struct {
	AccountID   string `json:"account_id" validate:"required"`
	TypeCash    string `json:"type_cash" validate:"required"`
	TypeInOut   string `json:"type_in_out" validate:"required"`
	Amount      int64  `json:"amount" validate:"required"`
	Category    string `json:"category" validate:"required"`
	Method      string `json:"method" validate:"required"`
	Description string `json:"description"`
}

type IdAdminTransactionRequest struct {
	Id string `json:"id" validate:"required"`
}

type UpdateAdminTransactionRequest struct {
	IdAdminTransactionRequest
	CreateAdminTransactionRequest
}

type SearchAdminTransactionRequest struct {
	TypeCash string `json:"type_cash"`
	Type     string `json:"type"`
}

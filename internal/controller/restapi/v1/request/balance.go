package request

type DeductBalance struct {
	UserID int64   `json:"user_id" validate:"required,gt=0"`
	Amount float64 `json:"amount" validate:"required,gt=0"`
}

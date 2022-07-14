package entity

type User struct {
	ID       string `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
	UserBalance
}

type UserBalance struct {
	Balance float64 `json:"current"`
	Spent   float64 `json:"withdrawn"`
}

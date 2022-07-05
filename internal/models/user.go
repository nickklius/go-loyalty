package models

type User struct {
	ID        string `json:"id"`
	FirstName string `json:"first-name"`
	LastName  string `json:"last-name"`
	Login     string `json:"login"`
	Password  string `json:"password"`
}

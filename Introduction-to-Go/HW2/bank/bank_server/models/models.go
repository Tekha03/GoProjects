package models

type Account struct {
	ID  	string 	`json:"id"`
	Name 	string 	`json:"name"`
	Balance float64 `json:"balance"`
}

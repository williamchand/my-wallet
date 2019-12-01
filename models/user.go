package models

// User represent the user model
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Customer struct {
	ID string `json:"customer_id"`
}

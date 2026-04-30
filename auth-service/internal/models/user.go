package models

type User struct {
	UserID   int64  `json:"UserID"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

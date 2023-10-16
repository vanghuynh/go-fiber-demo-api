package models

type User struct {
	ID        uint   `json:"id" gorm:"primaryKey" swaggerignore:"true"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UserName  string `json:"user_name"`
	Password  string `json:"password"`
	UserRole  string `json:"user_role"`
}

package types

type User struct {
	Base
	FullName    string `json:"full_name"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email" gorm:"unique"`
}

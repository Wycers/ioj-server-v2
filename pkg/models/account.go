package models

type Account struct {
	Model

	Name string `json:"name" gorm:"unique_index:account_name"`
	Nickname string `json:"nickname" gorm:"unique_index:account_nickname"`

	Locale string `json:"locale" gorm:"default: 'en'"`
	Avatar string `json:"avatar"`
	Email  string `json:"email"`
	Gender string `json:"gender"`
}

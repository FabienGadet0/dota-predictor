package models

// User model
type User struct {
	UserID      int `gorm:"primary_key";"AUTO_INCREMENT"`
	AccessToken string
	GrantLvl    int `gorm:"default:1"`
	NBCallsLeft int
}

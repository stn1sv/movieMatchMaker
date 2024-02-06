package models

type User struct {
	UserName string `json:"userName"`
	IsAdmin  bool   `json:"isAdmin" default:"false"`
}

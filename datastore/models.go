package datastore

import (
	"gorm.io/gorm"
)

// schema for user table

// User is the data type for user object
type User struct {
	gorm.Model
	ID         string
	Email      string `json:"email" validate:"required" sql:"email" gorm:"primary_key"`
	Password   string `json:"password" validate:"required" sql:"password"`
	ClientId   string `json:"clientId" validate:"required" gorm:"column:clientid;not null;primary_key"`
	Roles      string `json:"roles"  gorm:"column:roles";DEFAULT:default`
	TokenHash  string `json:"tokenhash" gorm:"column:tokenhash;not null;"`
	IsVerified bool   `json:"isverified" gorm:"column:isverified"`
}

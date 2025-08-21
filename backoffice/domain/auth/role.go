package auth

import (
	"backend-poc/backoffice/domain"
)

type Role struct {
	domain.Entity
	Title       string       `json:"title" gorm:"type:varchar(100);not null"`
	Description string       `json:"description" gorm:"type:text"`
	Permissions []Permission `json:"permissions" gorm:"many2many:role_permission"`
	Users       []User       `json:"users" gorm:"many2many:user_role"`
}

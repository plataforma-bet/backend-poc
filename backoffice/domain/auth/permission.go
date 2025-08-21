package auth

import "backend-poc/backoffice/domain"

type Permission struct {
	domain.Entity
	Name        string `json:"name" gorm:"type:varchar(100);not null"`
	Title       string `json:"title" gorm:"type:varchar(100);not null"`
	Description string `json:"description" gorm:"type:text"`
	Method      string `json:"method" gorm:"type:varchar(10);not null"`
	Path        string `json:"path" gorm:"type:varchar(100);not null"`
}

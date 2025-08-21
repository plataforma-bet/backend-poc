package domain

import (
	"time"

	"gorm.io/gorm"
)

type Entity struct {
	Id        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `gorm:"<-:create" json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

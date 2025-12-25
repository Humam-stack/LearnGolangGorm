package belajargolanggorm

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"type:varchar(250);uniqueIndex;not null"`
	Description string `gorm:"type:text"`
	CreatedAt   time.Time
	Products    []Product `gorm:"foreignKey:CategoryID"`
}

type Product struct {
	ID          uint    `gorm:"primaryKey"`
	Name        string  `gorm:"type:varchar(200);not null"`
	Description string  `gorm:"type:text"`
	Price       float64 `gorm:"type:decimal(10,2);not null"`
	Stock       int     `gorm:"default:0"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	CategoryID  uint           `gorm:"index;not null"`
	Category    Category       `gorm:"foreignKey:CategoryID"`
}

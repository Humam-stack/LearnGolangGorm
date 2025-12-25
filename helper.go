package belajargolanggorm

import "gorm.io/gorm"

func seedData(db *gorm.DB) error {
	var v int64
	db.Model(&Category{}).Where("name = ?", "Elektronik").Count(&v)
	if v > 0 {
		return nil
	}
	category := Category{
		Name:        "Elektronik",
		Description: "Barang Elektronik",
		Products: []Product{
			{Name: "HP", Price: 550000, Stock: 20},
			{Name: "Laptop", Price: 230000, Stock: 15},
			{Name: "Smartwatch", Price: 434343, Stock: 30},
		},
	}

	return db.Create(&category).Error
}

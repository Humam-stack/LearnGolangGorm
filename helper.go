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
	}

	if err := db.Create(&category).Error; err != nil {
		return err
	}

	products := []Product{
		{Name: "Android", Price: 550000, Stock: 20, CategoryID: category.ID},
		{Name: "Laptop", Price: 230000, Stock: 15, CategoryID: category.ID},
		{Name: "Smartwatch", Price: 434343, Stock: 30, CategoryID: category.ID},
	}

	return db.Create(&products).Error
}

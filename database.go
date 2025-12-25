package belajargolanggorm

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func SetupDB() *gorm.DB {
	dsn := "root@tcp(localhost:3306)/golanggormcrud?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	db.Migrator().DropTable(&Product{}, &Category{})
	db.AutoMigrate(&Category{}, &Product{})

	return db
}

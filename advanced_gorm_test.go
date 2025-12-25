package belajargolanggorm

import (
	"fmt"
	"testing"
)

func TestFilterData(t *testing.T) {
	db := SetupDB()
	seedData(db)

	fmt.Println("===Kondisi Tunggal===")
	var products []Product
	db.Where("price > ?", 10000).Find(&products)
	for _, p := range products {
		fmt.Println("Produk Nama : ", p.Name)
	}

	//multiple kondisi
	fmt.Println("===Multiple Kondisi===")
	db.Where("price > ? AND stock <?", 10000, 30).Find(&products)
	for _, a := range products {
		fmt.Println("Produk Nama : ", a.Name)
	}

	//dengan struct
	var pr Product
	fmt.Println("===Dengan Strcut===")
	db.Where(&Product{CategoryID: 1}).Find(&pr)
	fmt.Println("Produk Nama : ", pr.Name)

	// Dengan Select
	var pro []Product
	fmt.Println("===Dengan Select===")
	db.Select("name,price").Find(&pro)
	for _, b := range pro {
		fmt.Println("Produk Nama : ", b.Name)
		fmt.Println("Produk Price : ", b.Price)
	}

	//Order, Limit, OFFSET
	fmt.Println("===Dengan ORDER,LIMIT,OFFSET===")
	db.Order("price ASC").Limit(5).Find(&pro)
	for _, b := range pro {
		fmt.Println("Produk Nama : ", b.Name)
		fmt.Println("Produk Price : ", b.Price)
	}

	// Join Table pada gorm
	type Result struct {
		ProductName  string
		CategoryName string
	}

	var result []Result
	fmt.Println("===Join Table pada GORM===")
	db.Table("products").
		Select("products.name as product_name, categories.name as category_name").
		Joins("JOIN categories ON categories.id = products.category_id").
		Scan(&result)
	for _, res := range result {
		fmt.Printf("Product : %s | Category : %s\n", res.ProductName, res.CategoryName)
	}

	// GroupBY pada golang

	type Result2 struct {
		CategoryID uint
		Total      int64
	}
	fmt.Println("===Join Table pada GORM===")
	var result2 []Result2
	db.Model(&Product{}).Select("category_id, COUNT(*) as total").Group("category_id").Scan(&result2)
	for _, res2 := range result2 {
		fmt.Printf("CategoryID : %d | Category : %d\n", res2.CategoryID, res2.Total)
	}
}

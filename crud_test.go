package belajargolanggorm

import (
	"fmt"
	"testing"

	"gorm.io/gorm"
)

func TestCreateCategory(t *testing.T) {
	db := SetupDB()

	category := Category{
		Name:        "Elektronik",
		Description: "Barang Elektronik",
	}

	result := db.Create(&category)

	if result.Error != nil {
		t.Fatal("Error creating category : ", result.Error)
	}

	fmt.Println("Category Berhasil dibuat!")
	fmt.Println("ID:", category.ID)
	fmt.Println("Name : ", category.Name)
}

func TestCreateProduct(t *testing.T) {
	db := SetupDB()

	category := Category{Name: "Elektronik"}
	db.Create(&category)

	product := Product{
		Name:        "Laptop",
		Description: "Laptop Asus AAAA",
		Price:       1000000,
		Stock:       10,
		CategoryID:  category.ID,
	}

	result := db.Create(&product)

	if result.Error != nil {
		t.Fatal("Error Input Data Product", result.Error)
	}

	fmt.Println("Product Created!")
	fmt.Println("ID : ", product.ID)
	fmt.Println("Name : ", product.Name)
	fmt.Println("Price : ", product.Price)

}

func TestCreateProductNestedCateogry(t *testing.T) {
	db := SetupDB()
	category := Category{
		Name:        "Elektronik",
		Description: "Barang Elektronik",
		Products: []Product{
			{Name: "HP", Price: 550000, Stock: 20},
			{Name: "Laptop", Price: 230000, Stock: 15},
			{Name: "Smartwatch", Price: 434343, Stock: 30},
		},
	}

	result := db.Create(&category)

	if result.Error != nil {
		t.Fatal("Error :", result.Error)
	}

	fmt.Println("Category + Product Created!")
	fmt.Println("Category ID : ", category.ID)
	fmt.Println("Category Name : ", category.Name)
	fmt.Println("Total Product", len(category.Products))

	for i, product := range category.Products {
		fmt.Printf("%d. %s (ID: %d, Price : %.0f)\n",
			i+1, product.Name, product.ID, product.Price)
	}
}

func TestSearchGorm(t *testing.T) {
	db := SetupDB()
	seedData(db)

	var categoryRes Category
	var productID Product
	var ProductRest Product

	db.Find(&categoryRes, 1)
	fmt.Println("ID Category : ", categoryRes.ID)
	fmt.Println("Name Category : ", categoryRes.Name)
	fmt.Println("Description Category : ", categoryRes.Description)

	db.Find(&productID, 1)
	fmt.Println("ID Product : ", productID.ID)
	fmt.Println("Name Product : ", productID.Name)
	fmt.Println("Price Product : ", productID.Price)

	//Mencari Produk Dengan Preload

	db.Preload("Category").Find(&productID, 1)
	fmt.Println("Product Name : ", productID.Name)
	fmt.Println("Product Price : ", productID.Price)
	fmt.Println("Product Stock : ", productID.Stock)
	fmt.Println("Product Category : ", productID.Category.Name)

	//Mencari Produk Dengan Preload dan Kondisi tertentu dengan Where
	db.Debug().Preload("Category").Where("price < ?", 550000.00).First(&ProductRest)
	fmt.Println("Product Name : ", ProductRest.Name)
	fmt.Println("Product Price : ", ProductRest.Price)
	fmt.Println("Product Stock : ", ProductRest.Stock)
	fmt.Println("Product Category : ", ProductRest.Category.Name)
}

// Mencari Produk Dengan First Take / FIND
func TestSearchGormFirstTakeFind(t *testing.T) {
	db := SetupDB()
	seedData(db)

	var product1 Product
	err := db.Where("price > ?", 0).First(&product1).Error
	fmt.Println("First :", err, product1.Name)

	var product2 Product
	err2 := db.Where("price > ?", 0).Take(&product2).Error
	fmt.Println("Take : ", err2, product2.Name)

	var product3 []Product
	err3 := db.Where("price > ?", 0).Find(&product3).Error
	fmt.Println("Find :", err3, "Total:", len(product3))
}

// Update semua field Data product dengan Save
// Harus diisi semua, karena bila ada yang kosong maka akan NULL
//

func TestGormSaveUpdate(t *testing.T) {
	db := SetupDB()
	seedData(db)

	var p Product
	db.First(&p, 1)

	t.Logf("Data sebelum : %s - Rp%.0f", p.Name, p.Price)

	p.Name = "TV"
	p.Price = 600000
	p.Stock = 20

	res := db.Save(&p)
	if res.Error != nil {
		t.Fatal("Error update", res.Error)
	}

	//Verifikasi
	var updated Product
	db.First(&updated, 1)
	t.Logf("Data Setelah di update : %s - Rp%.0f ", updated.Name, updated.Price)

	if updated.Name != "TV" {
		t.Error("Nama tidak teupdate")
	}

	if updated.Price != 600000 {
		t.Error("Harga tidak terupdate")
	}

}

// Update hanya field yang disebutkan
// bisa update menjadi 0, false, atau string kosong
func TestUpdatedDenganUpdatesMap(t *testing.T) {
	db := SetupDB()
	seedData(db)

	res := db.Model(&Product{}).Where("id = ?", 1).Updates(map[string]interface{}{
		"price": 650000,
		"stock": 300,
	})

	if res.Error != nil {
		t.Fatal("Error :", res.Error)
	}

	//data setelah diupdated

	var product Product
	db.First(&product, 1)
	fmt.Println("Updated Price : ", product.Price)
	fmt.Println("Updated Stock :", product.Stock)

}

// Cara Updated Data ke nilai 0!
func TestZeroValueUpdatedGorm(t *testing.T) {
	db := SetupDB()
	seedData(db)

	//menggunakan map
	db.Model(&Product{}).Where("id = ?", 1).Updates(map[string]interface{}{
		"stock": 0,
	})

	//menggunakan select
	db.Model(&Product{}).Where("id = ?", 2).Select("stock").Updates(
		Product{Stock: 0},
	)
}

// Update batch data

func TestUpdateBatchData(t *testing.T) {
	db := SetupDB()
	seedData(db)

	res := db.Model(&Product{}).Where("price < ?", 50000000).Updates(map[string]interface{}{
		"stock": 100,
	})

	if res.Error != nil {
		t.Fatal("Error :", res.Error)
	}

}

// Upate dengan Expression SQL

func TestUpdateExpressionSQL(t *testing.T) {
	db := SetupDB()
	seedData(db)

	var sebelum Product
	db.First(&sebelum, 1)
	t.Logf("Stock sebelumnya : %d", sebelum.Stock)

	//tambah stock 5

	db.Model(&Product{}).Where("id = ?", 1).Update("stock", gorm.Expr("stock + ?", 5))

	var setelah Product
	db.First(&setelah, 1)
	t.Logf("Stock setelah di update : %d", setelah.Stock)
}

func TestSoftDelete(t *testing.T) {
	db := SetupDB()
	seedData(db)

	var sebelum int64
	db.Model(&Product{}).Count(&sebelum)
	t.Logf("Product sebelum di delete : %d", sebelum)

	//delte dengan id 1
	res := db.Delete(&Product{}, 1)
	if res.Error != nil {
		t.Fatal("Error : ", res.Error)
	}

	//delete semua product dengan dengan where
	// res := db.Where("stock < ?", 20).Delete(&Product{})

	// delete product dengan spesifik id
	// rest := db.Delete(&Product{}, []int{1,2,3})

	//hard delete
	// rest := db.Unscoped().Delete(&Product{},1)

	//cek data
	var sesudah int64
	db.Model(&Product{}).Count(&sesudah)
	t.Logf("Product setelah di delete : %d", sesudah)

	if sesudah != sebelum-1 {
		t.Error("Product belum didelete")
	}

	var p Product
	db.Unscoped().First(&p, 1)
	t.Logf("DeletedAt : %v", p.DeletedAt)

	//restore soft delete
	//rest := db.Unscoped().Model(&Product{}).Where("id = ? ",1).Update("deleted_at",nil)

}

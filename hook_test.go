package belajargolanggorm

import (
	"fmt"
	"testing"
)

func TestHook(t *testing.T) {
	db := SetupDB()
	seedData(db)

	fmt.Println("Testing Hook")
	var category Category
	db.Where("name = ?", "Elektronik").First(&category)

	product := Product{
		Name:       "Mouse Gaming",
		Price:      350000,
		Stock:      24,
		CategoryID: category.ID,
	}

	result := db.Create(&product)
	if result.Error != nil {
		t.Fatal(result.Error)
	}

	fmt.Printf("Product telah terbuat dengan nama : %s", product.Name)
	fmt.Printf("ID : %d", product.ID)
	fmt.Printf("Slug : %s", product.Slug)
}

func TestHookUpdates(t *testing.T) {
	db := SetupDB()
	seedData(db)

	var product Product
	if err := db.First(&product).Error; err != nil {
		t.Fatal(err)
	}

	product.Name = "Handphone"
	product.Price = 400000

	if err := db.Save(&product).Error; err != nil {
		t.Fatal("Update gagal : ", err)
	}

	if err := db.First(&product, product.ID).Error; err != nil {
		t.Fatal("Reload produk gagal : ", err)
	}

	fmt.Printf("Update : %s, price : RP%.0f", product.Name, product.Price)
	fmt.Printf("Slug : %s", product.Slug)

}

func TestHookDeleteStockOnValue(t *testing.T) {
	db := SetupDB()
	seedData(db)

	var product Product
	db.First(&product, 1)

	result := db.Delete(&product)

	if result.Error == nil {
		t.Error("Tidak bisa menghapus produk dengan stok 0")
	} else {
		t.Logf("Berhasil : %v", result.Error)
	}

}

func TestHooksDeleteStockZero(t *testing.T) {
	db := SetupDB()
	seedData(db)

	var product Product
	db.First(&product, 1)

	db.Model(&product).Update("stock", 0)

	result := db.Delete(&product)

	if result.Error != nil {
		t.Fatal("delete error,", result.Error)
	}

	fmt.Println("Sukses")
}

func TestHookFind(t *testing.T) {
	db := SetupDB()
	seedData(db)

	var product []Product
	db.Limit(2).Find(&product)
	fmt.Println("Produk (setelah find trigger) : ", len(product))

}

func TestHookValidation(t *testing.T) {
	db := SetupDB()
	seedData(db)

	product1 := Product{
		Name:       "Inavalid",
		Price:      -100000,
		CategoryID: 1,
	}

	result := db.Create(&product1)
	if result.Error == nil {
		t.Error("Fail : Harga Neggative!")
	} else {
		t.Logf("validation block karena harga negatice %v", result.Error)
	}

	product2 := Product{
		Name:       "Inavlid2",
		Price:      100000,
		Stock:      -10,
		CategoryID: 1,
	}

	result2 := db.Create(&product2)
	if result2.Error == nil {
		t.Error("X : harga negatice")
	} else {
		t.Logf("Harga diblock validation karena negative proce : %v", result2.Error)
	}

}

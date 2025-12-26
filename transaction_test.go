package belajargolanggorm

import (
	"errors"
	"testing"

	"gorm.io/gorm"
)

func TestBasicTransactionGorm(t *testing.T) {
	db := SetupDB()
	seedData(db)

	// 1. Begin Transaction GORM
	tx := db.Begin()

	// Operasinya
	category := Category{
		Name:        "Buku",
		Description: "Buku dan Majalah",
	}

	if err := tx.Create(&category).Error; err != nil {
		tx.Rollback()
		t.Fatal("Error Membuat Category :", err)
	}

	product := Product{
		Name:       "Novel Tere Liye",
		Price:      150000,
		Stock:      50,
		CategoryID: category.ID,
	}

	if err := tx.Create(&product).Error; err != nil {
		tx.Rollback() //Rollback kalau error
		t.Fatal("Error Buat product :", err)
	}

	//commit
	if err := tx.Commit().Error; err != nil {
		t.Fatal("Error Commit : ", err)
	}

	t.Log("Transcation commit berhasil")

	//verifikasi data tersimpan
	var save Product
	db.First(&save, product.ID)
	t.Logf("Product : %s", save.Name)

}

func TestTransactionDefer(t *testing.T) {
	db := SetupDB()
	seedData(db)

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	//Operasi1
	category := Category{
		Name: "Mainan",
	}
	if err := tx.Create(&category).Error; err != nil {
		t.Fatal(err)
	}

	// Operasi2
	product := Product{
		Name:       "Buzz Lighyear",
		Price:      250000,
		Stock:      10,
		CategoryID: category.ID,
	}

	if err := tx.Create(&product).Error; err != nil {
		t.Fatal(err)
	}

	// commit
	if err := tx.Commit().Error; err != nil {
		t.Fatal(err)
	}

	t.Log("Transaction berhasil!")

}

func TestTransactionCallback(t *testing.T) {
	db := SetupDB()
	seedData(db)

	err := db.Transaction(func(tx *gorm.DB) error {
		//Semua transaksi di dalam function ini
		// jika hasilnya nil maka akan auto commit
		// jika error maka auto rollback

		category := Category{
			Name:        "Olahraga",
			Description: "Peralatan Olahraga",
		}

		if err := tx.Create(&category).Error; err != nil {
			return err
		}

		product := []Product{
			{Name: "Bola", Price: 200000, Stock: 30, CategoryID: category.ID},
			{Name: "Basket", Price: 250000, Stock: 35, CategoryID: category.ID},
		}

		if err := tx.Create(&product).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		t.Fatal("Transaction gagal :", err)
	}

	t.Log("Transaction sukses")

	//verifikasi data
	var c int64
	db.Model(&Product{}).Where("category_id = ?", 2).Count(&c)
	t.Logf("%d product", c)
}

func TestNestedTransaction(t *testing.T) {
	//FUNGSI YANG MASIH SAYA PIKIR MEMBINGUNGKANNNNNNNN!!! HAHA!!!
	db := SetupDB()
	seedData(db)

	err := db.Transaction(func(tx *gorm.DB) error {
		//main

		//create category
		category := Category{
			Name:        "Dapur",
			Description: "Peralatan Dapur",
		}
		if err := tx.Create(&category).Error; err != nil {
			return err
		}

		//NesTED! (BERSARANG BOYYYYY)
		err := tx.Transaction(func(tx2 *gorm.DB) error {
			product1 := Product{
				Name:       "Pisau",
				Price:      15000,
				Stock:      10,
				CategoryID: category.ID,
			}
			if err := tx2.Create(&product1).Error; err != nil {
				return err
			}
			return nil
		})

		if err != nil {
			t.Log("Nested transaction telah di rollback, tapi yang main bisa dilanjutkan")
		}

		product2 := Product{
			Name:       "Wajan",
			Price:      200000,
			Stock:      20,
			CategoryID: category.ID,
		}

		if err := tx.Create(&product2).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		t.Fatal("Main transaction gagal :", err)
	}

	t.Log("Tansaction sukses")
}

func TestTransferStock(t *testing.T) {
	db := SetupDB()
	seedData(db)

	//Transfer 5 stock dari product 1 ke product 2
	sourceID := uint(1)
	targetID := uint(2)
	amount := 5

	err := db.Transaction(func(tx *gorm.DB) error {
		// Get source product
		var source Product
		if err := tx.First(&source, sourceID).Error; err != nil {
			return err
		}

		if source.Stock < amount {
			return errors.New("Stock Kurang")
		}

		//mengurangi stock source
		if err := tx.Model(&source).Update("stock", source.Stock-amount).Error; err != nil {
			return err
		}

		// menambahkan stock target
		if err := tx.Model(&Product{}).Where("id = ?", targetID).Update("stock", gorm.Expr("stock + ? ", amount)).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		t.Fatal("Transfer gagal", err)
	}

	//verifikasi

	var source, target Product
	db.First(&source, sourceID)
	db.First(&target, targetID)

	t.Logf("Transfer Sukses")
	t.Logf("Source stock : %d", source.Stock)
	t.Logf("Target stock : %d", target.Stock)
}

func TestMultipleOperationTransaction(t *testing.T) {
	db := SetupDB()
	seedData(db)

	err := db.Transaction(func(tx *gorm.DB) error {
		// Update category
		if err := tx.Model(&Category{}).Where("id = ?", 1).Update("Description", "Description ter-update").Error; err != nil {
			return err
		}

		// membuat produk baru
		product := Product{
			Name:       "Product 1",
			Price:      400000,
			Stock:      15,
			CategoryID: 1,
		}

		if err := tx.Create(&product).Error; err != nil {
			return err
		}

		// Update stock prodcut lain

		if err := tx.Model(&Product{}).Where("id = ?", 1).Update("stock", gorm.Expr("stock + ?", 10)).Error; err != nil {
			return err
		}

		// Delete Product
		if err := tx.Delete(&Product{}, 3).Error; err != nil {
			return nil
		}

		// menghitung total product
		var c int64
		if err := tx.Model(&Product{}).Count(&c).Error; err != nil {
			return err
		}

		t.Logf("Total Product setelah operasi diatas : %d", c)

		return nil
	})

	if err != nil {
		t.Fatal("transaction gagal : ", err)
	}

	t.Log("Semua operasi sukses!")
}

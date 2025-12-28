package belajargolanggorm

import (
	"fmt"
	"strings"
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
	Slug        string  `gorm:"type:varchar(250);uniqueIndex"` // akan auto digenerate!
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	CategoryID  uint           `gorm:"index;not null"`
	Category    Category       `gorm:"foreignKey:CategoryID"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) error {
	fmt.Println("Sebelum hook trigger")

	// melakukan validasi
	if p.Price < 0 {
		return fmt.Errorf("Harga tidak boleh negative")
	}

	if p.Stock < 0 {
		return fmt.Errorf("Error : stok tidak boleh negative")
	}

	// auto generate slug dari kolom name
	if p.Slug == "" {
		p.Slug = strings.ToLower(strings.ReplaceAll(p.Name, " ", "-"))
	}

	// set default values
	if p.Stock == 0 {
		p.Stock = 1 // default stock
	}

	fmt.Printf("Slug Generated! : %s\n", p.Slug)
	return nil
}

// After Delete

func (p *Product) AfterCreate(tx *gorm.DB) error {
	fmt.Println("Setelah hook trigger")
	fmt.Printf("Produk dibuat dengan ID : %d\n", p.ID)

	return nil
}

//Befor Update

func (p *Product) BeforeUpdate(tx *gorm.DB) error {
	fmt.Println("Before Update ter-terigger")

	//validasi perubahan
	if tx.Statement.Changed("price") {
		if p.Price < 0 {
			return fmt.Errorf("price tidak boleh negative")
		}
		fmt.Printf("Price berhasild dirubah : Rp%.0f", p.Price)
	}

	if tx.Statement.Changed("stock") {
		if p.Stock < 0 {
			return fmt.Errorf("Price tidak boleh negative")
		}
		fmt.Printf("Stock telah dirubah menjadi : %d\n", p.Stock)
	}

	// Update slug jika name berubah!
	p.Slug = strings.ToLower(strings.ReplaceAll(p.Name, " ", "-"))

	return nil
}

func (p *Product) AfterUpdate(tx *gorm.DB) error {
	fmt.Println("AfterUpdate Ter-Trigger")
	fmt.Println("Perubahan Berhasil !")

	return nil
}

func (p *Product) BeforeDelete(tx *gorm.DB) error {
	fmt.Println("Befor Delete ter-terigger!")

	// validasi : tidak boleh menghapus product dengan stock > 0
	if p.Stock > 0 {
		return fmt.Errorf("Gagal menghapus product karena stock nya masih :%d", p.Stock)
	}

	fmt.Printf("Product dengan ID : %d, telah terhapus", p.ID)
	return nil
}

func (p *Product) AfterDelete(tx *gorm.DB) error {
	fmt.Println("After Delete Trigger!")
	fmt.Println("Produk Berhasil di hapus!")
	return nil
}

func (p *Product) BeforeSave(tx *gorm.DB) error {
	fmt.Println("BeforeSave ter-trrigger!")

	// Validasi yang umum untuk create dan update

	if len(p.Name) < 3 {
		return fmt.Errorf("minimal 3 huruf untuk nama dair produk")
	}

	return nil
}

func (p *Product) AfterSave(tx *gorm.DB) error {
	fmt.Println("Produk berhasil di save")
	return nil
}

func (p *Product) AfterFind(tx *gorm.DB) error {
	fmt.Printf("Produk dengan Nama : %s\n", p.Name)
	return nil
}

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

func AvailableProduct(db *gorm.DB) *gorm.DB {
	return db.Where("stock > ?", 0)
}

//scope dengan minimal harga

func minPrice(minPrice float64) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("price >= ?", minPrice)
	}
}

// scope dengan maxharga
func maxPrice(maxPrice float64) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("price <= ?", maxPrice)
	}
}

// scope dengan order by price
func OrderByPrice(ascending bool) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if ascending {
			return db.Order("price ASC")
		}
		return db.Order("price DESC")
	}
}

func Paginate(page, pageSize int) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func SearchByName(keyword string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if keyword == "" {
			return db
		}
		return db.Where("name LIKE ?", keyword)
	}
}

func SearchByCategory(cid uint) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if cid == 0 {
			return db
		}
		return db.Where("category_id = ?", cid)
	}
}

func FilterProducts(filter map[string]interface{}) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		query := db

		if minPrice, ok := filter["min_price"].(float64); ok && minPrice > 0 {
			query = db.Where("price >= ?", minPrice)
		}

		if maxPrice, ok := filter["max_price"].(float64); ok && maxPrice > 0 {
			query = db.Where("price < ?", maxPrice)
		}

		if keyword, ok := filter["keyword"].(string); ok && keyword != "" {
			query = db.Where("name LIKE ? ", "%"+keyword+"%")
		}

		if categoryId, ok := filter["category_id"].(uint); ok && categoryId > 0 {
			query = db.Where("category_id = ?", categoryId)
		}

		return query
	}
}

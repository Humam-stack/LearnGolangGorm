package belajargolanggorm

import (
	"fmt"
	"testing"
)

func TestWhereGorm(t *testing.T) {
	db := SetupDB()
	seedData(db)

	var products []Product

	// Where saja
	db.Where("price > ?", 500000).Find(&products)
	t.Logf("Produk dengan price > 500rb : %d", len(products))

	// Where dengan banyak kondisi
	db.Where("price > ? AND stock > ?", 300000, 20).Find(&products)
	t.Logf("Produk dengan harga > 300000 dan stock > 20 : %d", len(products))

	// where dengan OR
	db.Where("price < ?", 300000).Or("stock > ?", 25).Find(&products)
	t.Logf("produk dengan harga < 300k atau stock itu > 25 : %d", len(products))

	// where dengan NOT
	db.Not("name = ? ", "HP").Find(&products)
	t.Logf("Produk yang bukan HP : %d", len(products))

	// where dengan IN
	db.Where("id IN ?", []int{1, 2, 3}).Find(&products)
	t.Logf("Produk dengan ID in 1,2,3 : %d", len(products))

	// where dengan LIKE
	db.Where("name LIKE ?", "%Laptop%").Find(&products)
	t.Logf("Produk yang ada nama laptopnya : %d", len(products))

	// where dengan between
	db.Where("price Between ? and ?", 200000, 500000).Find(&products)
	t.Logf("Produk dengan harga diantara 200rb-500rb : %d", len(products))

	//  where dengan struct
	db.Where(&Product{Stock: 20}).Find(&products)
	t.Logf("Produk dengan stok 20 : %d", len(products))

	// where dengan MAP
	db.Where(map[string]interface{}{
		"stock": 30,
		"price": 434343,
	}).Find(&products)
	t.Logf("produk dengan stok 30 dan harganya 434343 : %d", len(products))
}

func TestOrderLimitPagination(t *testing.T) {
	db := SetupDB()
	seedData(db)

	var products []Product

	// Order By price DESC( mengurutkan produk dengan harga yang termahal)
	db.Order("price DESC").Find(&products)
	t.Log("===Harga Termahal===")
	for _, p := range products {
		t.Logf("%s - Rp%.0f", p.Name, p.Price)
	}

	// Order by Name ASC (dari alphabet A-Z)
	db.Order("name ASC").Find(&products)
	for _, s := range products {
		t.Logf("%s", s.Name)
	}

	// Limit (mengambil 2 / ?  data pertama)
	db.Limit(2).Find(&products)
	t.Logf("\n2 Data pertama dari produk : %d", len(products))

	// OFFSet(skip data pertama , dan mengmabil sisanya)
	db.Offset(1).Find(&products)
	t.Logf("Skip data pertama dan ambil : %d", len(products))

	// Pagination (page 1 dengan size 2 data)
	page := 1
	pageSize := 2
	db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&products)
	t.Logf("\nPage %d (size : %d) : %d products", page, pageSize, len(products))

	page = 2
	db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&products)
	t.Logf("Page %d dengan size %d : %d produk", page, pageSize, len(products))

	// kombinasi : order _ limit (mencari top 2 termahal karena data yang ada di seed data hanya 3)
	db.Order("price DESC").Limit(2).Find(&products)
	t.Log("\n 2 Produk Termahal")
	for i, c := range products {
		t.Logf("%d. %s - Rp%.0f", i+1, c.Name, c.Price)
	}

}

func TestJoinGorm(t *testing.T) {
	db := SetupDB()
	seedData(db)

	var products []Product

	//Inner Join
	db.Joins("Category").Find(&products)
	t.Log("===Join Category===")
	for _, p := range products {
		t.Logf("%s - Category : %s", p.Name, p.Category.Name)
	}

	//Join degan Where Pada table join

	db.Joins("Category").Where("categories.name = ?", "Elektronik").Find(&products)
	t.Logf("\n Produk di kategori Elektronik ada : %d", len(products))

	//joins dengan select untuk kolom yang spesifik
	db.Select("products.*, categories.name as category_name").Joins("Category").Find(&products)

	// Left Join
	db.Joins("Left JOIN categories ON categories.id = products.category_id").Find(&products)
}

func TestPreloadVsJoins(t *testing.T) {
	db := SetupDB()
	seedData(db)

	var products1, products2 []Product
	// digunakan unuk semua data relasi (has many)
	db.Preload("Category").Find(&products1)
	for _, p := range products1 {
		t.Logf("%s - Category : %s", p.Name, p.Category.Name)
	}

	//join
	//digunakan kalau anda memfilter by relasi atau belongsto
	db.Joins("Category").Find(&products2)

	for _, p := range products2 {
		t.Logf("%s - Category : %s", p.Name, p.Category.Name)
	}
}

func TestGroupByHaving(t *testing.T) {
	db := SetupDB()
	seedData(db)

	type Result struct {
		CategoryID   uint
		CategoryName string
		TotalProduct int64
	}

	var results []Result
	db.Model(&Product{}).Select(
		"category_id, categories.name as category_name, COUNT(*) as total_product").
		Joins("LEFT JOIN categories ON categories.id = products.category_id").
		Group("category_id").Scan(&results)

	t.Log("\n===Product per Category===")
	for _, r := range results {
		t.Logf("Category :%s - Total : %d produk", r.CategoryName, r.TotalProduct)
	}

	// GroupBY dengan having
	db.Model(&Product{}).
		Select("category_id, COUNT(*) as total").
		Group("category_id").Having("COUNT (*) > ?", 1).Scan(&results)

	t.Logf("\n Categories dengan >1 produk : %d", len(results))
}

func TestAggregateGorm(t *testing.T) {
	db := SetupDB()
	seedData(db)

	// Count
	var c int64
	db.Model(&Product{}).Count(&c)
	t.Logf("Total Product : %d", c)

	// count dengan where
	db.Model(&Product{}).Where("price > ?", 300000).Count(&c)
	t.Logf("Product dengan harga > 300rb : %d", c)

	// sum
	var total float64
	db.Model(&Product{}).Select("SUM (price * stock)").Scan(&total)
	t.Logf("Total Semua Harga di dalam data : RP%.0f", total)

	//avg
	var avg float64
	db.Model(&Product{}).Select("AVG(price)").Scan(&avg)
	t.Logf("Rata rata Harga : Rp%.0f", avg)

	// max
	var max float64
	db.Model(&Product{}).Select("MAX(price)").Scan(&max)
	t.Logf("Data dengan harga paling Mahal : RP%.0f", max)

	// min
	var min float64
	db.Model(&Product{}).Select("MIn(price)").Scan(&min)
	t.Logf("Paling murah : RP%.0f", min)

	// Multiple Agregate
	type Stats struct {
		Total int64
		Avg   float64
		Max   float64
		Min   float64
	}

	var status Stats
	db.Model(&Product{}).Select("Count(*) as total, AVG(price) as avg, MAX(price) as max, MIN(price) as min").Scan(&status)
	fmt.Println("\n===Statistik Produk===")
	fmt.Printf("Total : %d", status.Total)
	fmt.Printf("Avg : %0.f", status.Avg)
	fmt.Printf("Max : %0.f", status.Max)
	fmt.Printf("Min : %0.f", status.Min)
}

func TestSubqueries(t *testing.T) {
	db := SetupDB()
	seedData(db)

	var product []Product

	//Subquery di Where
	db.Where("price > (?)", db.Model(&Product{}).Select("AVG(price)")).Find(&product)
	for _, p := range product {
		fmt.Printf("\n%s - Rp%0.f ", p.Name, p.Price)
	}

	// Subquery dengan IN
	// mencari produk di category dengan total produk > 1

	db.Where("category_id IN (?)", db.Model(&Product{}).Select("category_id").Group("category_id").Having("Count(*) > ?", 1)).Find(&product)

	fmt.Printf("\nProduk : %d ", len(product))

}

package belajargolanggorm

import (
	"fmt"
	"testing"
)

func TestScopeBasic(t *testing.T) {
	db := SetupDB()
	seedData(db)

	var product []Product

	db.Scopes(AvailableProduct).Find(&product)
	fmt.Printf("Produk dengan stock yang masih ada, ada : %d produk", len(product))

	var min []Product

	db.Scopes(minPrice(300000)).Find(&min)
	fmt.Printf("Produk dengan harga > 300k %d", len(min))

	var max []Product
	db.Scopes(maxPrice(300000)).Find(&max)
	fmt.Printf("Produk dengan harga < 300k ada : %d", len(max))

	var filter []Product
	db.Scopes(AvailableProduct, minPrice(300000), maxPrice(500000)).Find(&filter)
	fmt.Printf("Produk yang ada dan minimal price nya 300k dan maximal pricenya 500k ada : %d", len(filter))
}

func TestPaginationScope(t *testing.T) {
	db := SetupDB()
	seedData(db)

	page := 1
	pageSize := 2

	var product []Product
	db.Scopes(Paginate(page, pageSize)).Find(&product)
	fmt.Printf("Page %d (pageSizenya : %d) : %d produk", page, pageSize, len(product))

	page = 2
	db.Scopes(Paginate(page, pageSize)).Find(&product)
	fmt.Printf("\nPage %d (pageSize %d) : %d produk", page, pageSize, len(product))
}

func TestScopesSearch(t *testing.T) {
	db := SetupDB()
	seedData(db)

	var product []Product
	db.Scopes(SearchByName("Laptop")).Find(&product)
	fmt.Printf("Produk bernama laptop terdapat : %d\n", len(product))

	db.Scopes(SearchByName("Laptop"), SearchByCategory(1), Paginate(1, 2)).Find(&product)

	fmt.Printf("Category 1, page 1 : %d result", len(product))
}

func TestDynamicScopes(t *testing.T) {
	db := SetupDB()
	seedData(db)

	filters := map[string]interface{}{
		"min_price": 200000.0,
		"max_price": 600000.0,
		"keyword":   "Laptop",
	}

	var product []Product
	db.Scopes(FilterProducts(filters)).Find(&product)

	for _, p := range product {
		fmt.Printf("- %s : RP.%.0f", p.Name, p.Price)
	}
}

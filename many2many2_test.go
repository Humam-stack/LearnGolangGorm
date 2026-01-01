package belajargolanggorm

import (
	"fmt"
	"testing"
)

func TestManyToMany(t *testing.T) {
	db := SetupDB()
	db.Migrator().DropTable(&User{}, &Role{})
	db.AutoMigrate(&User{}, &Role{})

	adminRole := Role{Name: "Admin"}
	staffRole := Role{Name: "Staff"}
	userRole := Role{Name: "User"}

	db.Create(&adminRole)
	db.Create(&staffRole)
	db.Create(&userRole)

	user := []User{
		{
			Name:  "admin",
			Email: "admin@admin.com",
			Roles: []*Role{&adminRole},
		},
		{
			Name:  "staff",
			Email: "staff@staff.com",
			Roles: []*Role{&staffRole},
		},
		{
			Name:  "user",
			Email: "user@user.com",
			Roles: []*Role{&userRole},
		},
	}

	for i := range user {
		db.Create(&user[i])
	}

	var u []User
	db.Preload("Roles").Find(&u)

	for _, p := range u {
		fmt.Printf("Username : %s\n", p.Name)
		for _, r := range p.Roles {
			fmt.Printf("Roles : %s \n", r.Name)
		}
	}

	//association// "menambah role baru" contoh kita akan menggunakan append!
	var findUser User
	db.First(&findUser, 1)

	// mengambil staff role dari DB
	var staffRoleFromDB Role
	db.Where("name = ?", "Staff").Find(&staffRoleFromDB)
	db.Model(&findUser).Association("Roles").Append(&staffRoleFromDB)

	//Setelah di append ! :
	db.Preload("Roles").Find(&u)
	fmt.Println("\nSETELAH ASSOCIATION APPEND PADA USER DENGAN ID 1")
	for _, p := range u {
		fmt.Printf("Username : %s\n", p.Name)
		for _, r := range p.Roles {
			fmt.Printf("Roles : %s\n", r.Name)
		}
	}

	//Hapus spesific ROLE!
	db.First(&findUser, 1)
	db.Model(&findUser).Association("Roles").Delete(&staffRoleFromDB)
	//Setelah di DELETE ! :
	db.Preload("Roles").Find(&u)
	fmt.Println("\nSETELAH ASSOCIATION DELETE PADA USER DENGAN ID 1")
	for _, p := range u {
		fmt.Printf("Username : %s\n", p.Name)
		for _, r := range p.Roles {
			fmt.Printf("Roles : %s\n", r.Name)
		}
	}
}

package belajargolanggorm

import (
	"fmt"
	"log"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name    string
	Profile Profile
	Post    []Post
}

type Profile struct {
	gorm.Model
	Address string
	UserID  uint
}

type Post struct {
	gorm.Model
	Title       string
	Description string
	UserID      uint
	Tag         []Tag `gorm:"many2many:tag_post;"`
}

type Tag struct {
	gorm.Model
	Name string
	Post []Post `gorm:"many2many:tag_post;"`
}

func TestGormChallange(t *testing.T) {
	dsn := "root@tcp(localhost:3306)/belajargolanggorm?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}

	db.AutoMigrate(&User{}, &Profile{}, &Post{}, &Tag{})

	// user := User{
	// 	Name: "Humam",
	// 	Profile: Profile{
	// 		Address: "Indonesia",
	// 	},
	// 	Post: []Post{
	// 		{
	// 			Title:       "Post1",
	// 			Description: "Description Post 1",
	// 			Tag: []Tag{
	// 				{Name: "News"},
	// 				{Name: "Tech"},
	// 			},
	// 		},
	// 	},
	// }
	// db.Create(&user)

	var Users User
	db.Preload("Profile").Preload("Post.Tag").First(&Users, "name = ?", "Humam")
	fmt.Println(" User : ", Users.Name)
	fmt.Println(" Address : ", Users.Profile.Address)
	for _, result := range Users.Post {
		fmt.Println("List Post : ", result.Title)

		for _, tag := range result.Tag {
			fmt.Println("  Tag : ", tag.Name)
		}
	}

}

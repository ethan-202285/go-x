package auth

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
)

var (
	auth       *Auth
	repository *Repository
)

func init() {
	dsn := fmt.Sprintf("host=db port=5432 user=app dbname=app password=app sslmode=disable")
	db, err := gorm.Open("postgres", dsn)
	if err != nil {
		log.Fatal("openDB: ", err)
	}
	db.LogMode(true)

	secretKey := []byte("aasdfkjksjdfaaasdfkjksjdfa123405") // 32 bytes
	auth = New(Config{DB: db, SecretKey: secretKey})
	repository = auth.Repository
}

package auth

import (
	"fmt"
	"log"
	"time"

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

	// 缩短时间，方便测试
	cleanupInterval = 50 * time.Millisecond
	DefaultTokenLife = 50 * time.Millisecond

	secretKey := []byte("aasdfkjksjdfaaasdfkjksjdfa123405") // 32 bytes
	auth = New(Config{DB: db, SecretKey: secretKey})
	repository = auth.Repository
}

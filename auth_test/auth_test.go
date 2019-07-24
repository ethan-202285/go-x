package auth_test

import (
	"fmt"
	"time"

	"github.com/goodwong/go-x/auth"
	"github.com/jinzhu/gorm"
)

var (
	secretKey  = []byte("aasdfkjksjdfaaasdfkjksjdfa123405") // 32 bytes
	auths      *auth.Auth
	repository *auth.Repository
	db         *gorm.DB
)

func init() {
	dsn := fmt.Sprintf("host=db port=5432 user=app dbname=app password=app sslmode=disable")
	var err error
	db, err = gorm.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	//db.LogMode(true)

	// 缩短时间，方便测试
	auth.CleanupInterval = 50 * time.Millisecond
	auth.DefaultTokenLife = 50 * time.Millisecond

	auths = auth.New(auth.Config{DB: db, SecretKey: secretKey})
	repository = auths.Repository
}

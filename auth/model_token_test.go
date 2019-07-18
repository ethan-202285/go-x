package auth

import (
	"log"
	"testing"
)

func TestToken(t *testing.T) {
	secretKey := []byte("aasdfkjksjdfaaasdfkjksjdfa123405")
	// 生成token
	token := newToken(Token{UserID: 9999, Device: "deviceID", Remark: "remark"}, secretKey)
	token.ID = 1234567899 // 模拟保存数据库

	// 发送给客户端，由客户端保存，tokenString内含不可复原的nonce
	tokenString := token.TokenString()
	log.Println(tokenString)

	// 验证
	parsedToken, err := parseTokenString(tokenString, secretKey)
	if err != nil {
		t.Fatalf("解析失败：%s", err)
	}
	parsedToken.Hash = token.Hash // 模拟读取数据库
	log.Printf("%d,  %t\n", parsedToken.ID, parsedToken.Verify())
}

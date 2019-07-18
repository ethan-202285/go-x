package auth

import (
	"testing"
)

var (
	loginUser *User
)

func init() {
	// 注册password登陆
	passwordProvider := NewPasswordProvider(auth)
	auth.Service.RegisterProvider(passwordProvider)

	// 准备数据
	auth.db().Delete(&User{}, "username IN (?)", []string{"testpassword"})
	auth.db().Delete(&UserIdentity{}, "open_id IN (?)", []string{"testpassword"})
	username, password := "testpassword", "testpassword"
	var err error
	loginUser, err = passwordProvider.Register(username, password)
	if err != nil {
		panic(err)
	}
}

func TestLogin(t *testing.T) {
	credentials := []byte(`{"username":"testpassword", "password":"testpassword"}`)
	tokens, err := auth.Service.Login("password", credentials, true, "gotest")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", tokens)

	// 再登陆一次
	tokens, err = auth.Service.Login("password", credentials, true, "gotest")
	if err != nil {
		t.Fatal(err)
	}
	if tokens.RefreshToken == nil {
		t.Fatal("RefreshToken生成失败")
	}
	t.Logf("%+v", tokens)

	// 多次登陆，只保留最新的
	count := 0
	auth.db().Model(&Token{}).Where(&Token{Device: "gotest", UserID: loginUser.ID}).Count(&count)
	if count != 1 {
		t.Fatal("多次登陆，理应只有一个有效Token，现在却是：", count)
	}

	// Renew
	tokens, err = auth.Service.Renew(*tokens.RefreshToken)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", tokens)
	if tokens.RefreshToken != nil {
		t.Fatal("Renew不应该生成RefreshToken")
	}
	if tokens.Token == "" {
		t.Fatal("Renew Token失败")
	}

	// Logout
}

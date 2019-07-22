package auth

import (
	"testing"
	"time"
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
	auth.db().Unscoped().Delete(&Token{}, "device IN (?)", []string{"test", "gotest"})
	username, password := "testpassword", "testpassWord123,"
	var err error
	loginUser, err = passwordProvider.Register(username, password)
	if err != nil {
		panic(err)
	}
}

func TestLogin(t *testing.T) {
	credentials := []byte(`{"username":"testpassword", "password":"testpassWord123,"}`)
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
	_, tokens, err = auth.Service.Renew(*tokens.RefreshToken)
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
	jwtToken, err := auth.jwtauth.Decode(tokens.Token)
	if err != nil {
		t.Fatal(err)
	}
	if needRelogin := auth.Service.Validate(jwtToken); needRelogin {
		t.Fatal("不需要重新登录")
	}
	if err := auth.Service.Logout(loginUser, "gotest"); err != nil {
		t.Fatal(err)
	}
	if needRelogin := auth.Service.Validate(jwtToken); !needRelogin {
		t.Fatal("理应需要重新登录")
	}
}

func TestClearLogoutsLoop(t *testing.T) {
	count := 0

	// 第一次测量
	auth.Service.logouts.Range(func(k, v interface{}) bool {
		count++
		return true
	})
	if count != 1 {
		t.Fatalf("logouts条目应该是1条，测试的数据%d条，不对\n", count)
	}

	// 清理后的测量
	time.Sleep(2 * cleanupInterval)
	count = 0
	auth.Service.logouts.Range(func(k, v interface{}) bool {
		count++
		return true
	})
	if count != 0 {
		t.Fatalf("清理后的logouts条目应该是0条，测试的数据%d条，不对\n", count)
	}
}

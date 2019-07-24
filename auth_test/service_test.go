package auth_test

import (
	"testing"

	"github.com/goodwong/go-x/auth"
	"github.com/goodwong/go-x/auth/providers/password"
)

var (
	loginUser *auth.User
)

func init() {
	// 注册password登陆
	passwords := password.NewProvider(&password.Config{
		Auth:      auths,
		SecretKey: secretKey,
	})
	auths.RegisterProvider(passwords)

	// 准备数据
	db.Delete(&auth.User{}, "username IN (?)", []string{"testpassword"})
	db.Delete(&auth.UserIdentity{}, "open_id IN (?)", []string{"testpassword"})
	db.Unscoped().Delete(&auth.Token{}, "device IN (?)", []string{"test", "gotest"})
	username, password := "testpassword", "testpassWord123,"
	var err error
	loginUser, err = passwords.Register(username, password)
	if err != nil {
		panic(err)
	}
}

func TestLogin(t *testing.T) {
	credentials := []byte(`{"username":"testpassword", "password":"testpassWord123,"}`)
	tokens, err := auths.Service.Login("password", credentials, true, "gotest")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", tokens)

	// 再登陆一次
	tokens, err = auths.Service.Login("password", credentials, true, "gotest")
	if err != nil {
		t.Fatal(err)
	}
	if tokens.RefreshToken == nil {
		t.Fatal("RefreshToken生成失败")
	}
	t.Logf("%+v", tokens)

	// 多次登陆，只保留最新的
	count := 0
	db.Model(&auth.Token{}).Where(&auth.Token{Device: "gotest", UserID: loginUser.ID}).Count(&count)
	if count != 1 {
		t.Fatal("多次登陆，理应只有一个有效Token，现在却是：", count)
	}

	// Renew
	_, tokens, err = auths.Service.Renew(*tokens.RefreshToken)
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
	//jwtToken, err := auths.jwtauth.Decode(tokens.Token)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//if needRelogin := auths.Service.JwtInvalid(jwtToken); needRelogin {
	//	t.Fatal("不需要重新登录")
	//}
	//if err := auths.Service.Logout(loginUser, "gotest"); err != nil {
	//	t.Fatal(err)
	//}
	//if needRelogin := auths.Service.JwtInvalid(jwtToken); !needRelogin {
	//	t.Fatal("理应需要重新登录")
	//}
}

func TestClearLogoutsLoop(t *testing.T) {
	//count := 0
	//
	//// 第一次测量
	//auths.Service.logouts.Range(func(k, v interface{}) bool {
	//	count++
	//	return true
	//})
	//if count != 1 {
	//	t.Fatalf("logouts条目应该是1条，测试的数据%d条，不对\n", count)
	//}
	//
	//// 清理后的测量
	//time.Sleep(2 * auths.CleanupInterval)
	//count = 0
	//auths.Service.logouts.Range(func(k, v interface{}) bool {
	//	count++
	//	return true
	//})
	//if count != 0 {
	//	t.Fatalf("清理后的logouts条目应该是0条，测试的数据%d条，不对\n", count)
	//}
}

package auth_test

import (
	"net/http/httptest"
	"testing"

	"github.com/goodwong/go-x/auth"
)

func TestContextRepository(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost/api/login", nil)
	ctx := auth.NewContext(req.Context())
	// 测试 WithUserID
	ctx.WithUserID(125)
	if userID := ctx.UserID(); userID != 125 {
		t.Fatal("ctx.UserID() 获取userID失败")
	}

	// 测试 WithUser
	ctx.WithUser(&auth.User{ID: 15})
	user := ctx.User()
	if user == nil {
		t.Fatal("ctx.User() 获取user失败")
	}
	t.Logf("user from ctx: %+v\n", user)

	user.Name = "老小王"
	user2 := ctx.User() // 模拟数据库
	t.Logf("user from ctx: %+v\n", user2)

}

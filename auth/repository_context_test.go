package auth

import (
	"net/http/httptest"
	"testing"
)

func TestContextRepository(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost/api/login", nil)
	context := auth.NewContext(req)
	// 测试 WithUserID
	context.WithUserID(125)
	if userID := context.UserID(); userID != 125 {
		t.Fatal("context.UserID() 获取userID失败")
	}

	// 测试 WithUser
	context.WithUser(&User{ID: 15})
	user := context.User()
	if user == nil {
		t.Fatal("context.User() 获取user失败")
	}
	t.Logf("user from context: %+v\n", user)

	user.Name = "老小王"
	user2 := context.User() // 模拟数据库
	t.Logf("user from context: %+v\n", user2)

}

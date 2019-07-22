package auth

import (
	"testing"
)

var (
	user *User
	err  error
)

func init() {
	// 准备
	auth.db().Delete(&User{}, "username IN (?)", []string{"test_user", "new_username", "test"})
	user, err = repository.Create("test", "测试号")
}

func TestUserCreate(t *testing.T) {

	// 创建
	user, err = repository.Create("test_user", "老小王")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(user)
	// 重复创建
	user, err = repository.Create("test_user", "老小王")
	if err == nil {
		t.Logf("错误的保存user: % v\n", user)
		t.Fatal("重复用户名理应报错，却没有报")
	}
	// 第三个参数
	user, err = repository.Create("william", "威廉", "https://static.dingtalk.com/media/lADPDgQ9qSqL8CjNAtjNAu4_750_728.jpg_120x120q90.jpg")
	if err != nil {
		t.Fatal(err)
	}
}

func TestUserUpdate(t *testing.T) {
	// 更新
	err = repository.Update(user, User{Username: "bad_username", Name: "小威廉"})
	if err != nil {
		t.Fatal(err)
	}
	if user.Username != "william" {
		t.Fatal("Update理应只允许更新Name和Avatar的，这里Username也被改了")
	}
	// 更新username
	err = repository.UpdateUserName(user, "new_username")
	if err != nil {
		t.Fatal(err)
	}
	if user.Username != "new_username" {
		t.Fatal("更新Username失败")
	}
	// 更新已存在的username
	err = repository.UpdateUserName(user, "test_user")
	if err == nil {
		t.Fatal("错误更新重复的Username")
	}
	//! todo @gorm bug
	//!if user.Username == "test_user" {
	//!	t.Error("错误更新重复的Username")
	//!}
}

func TestUserFind(t *testing.T) {
	// 按ID
	user, err = repository.Find(user.ID)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Find: % v\n", user)

	// 按username
	user, err = repository.FindByUsername("test_user")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("FindByUsername: % v\n", user)
}

func TestUserList(t *testing.T) {
	offset := 0
	limit := 10
	users, err := repository.List(offset, limit, "id < ?", 10)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("List() %d 人:\n", len(users))

	count, err := repository.Count("id < ?", 10)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Count() %d 人:\n", count)
}

package auth_test

import (
	"testing"

	"github.com/goodwong/go-x/auth"
)

var (
	identity *auth.UserIdentity
)

func init() {
	// 准备
	db.Delete(&auth.UserIdentity{}, "provider IN (?)", []string{"test"})
}

func TestIdentityCreate(t *testing.T) {
	// 创建
	userID, provider, openID := user.ID, "test", "test_user"
	data := map[string]string{
		"password": "123456",
		"爱好":       "打球",
	}
	identity, err = repository.CreateIdentity(userID, provider, openID, data)
	if err != nil {
		t.Fatal(err)
	}

	// 重复主键
	_, err = repository.CreateIdentity(userID, provider, openID)
	if err == nil {
		t.Fatal(err)
	}
}

func TestIdentityUpdate(t *testing.T) {
	repository.UpdateIdentityData(identity, nil)
	repository.UpdateIdentityUser(identity, &auth.User{ID: 1})
}

func TestIdentityFind(t *testing.T) {
	provider, openID := "test", "test_user"
	identity, err = repository.FindIdentity(provider, openID)
	if err != nil {
		t.Fatal(err)
	}

	userID := uint64(1)
	identity, err = repository.FindIdentityByUser(userID, provider)
	if err != nil {
		t.Fatal(err)
	}
}

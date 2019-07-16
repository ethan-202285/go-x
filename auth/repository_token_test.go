package auth

import (
	"testing"
)

func TestTokenInRepository(t *testing.T) {
	// Create
	userID, device, remark := uint64(1), "test", "test"
	_, tokenString, err := repository.CreateToken(userID, device, remark)
	if err != nil {
		t.Fatal(err)
	}

	// Find
	token, err := repository.FindToken(tokenString)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("token: % v\n", token)

	// Delete
	err = repository.DeleteToken(userID, device)
	if err != nil {
		t.Fatal(err)
	}
}

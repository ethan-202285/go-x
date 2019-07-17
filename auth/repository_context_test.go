package auth

import (
	"net/http/httptest"
	"testing"
)

func TestContextRepository(t *testing.T) {
	req := httptest.NewRequest("GET", "http://localhost/api/login", nil)
	context := auth.NewContext(req)
	context.WithUserID(125)
	if userID := context.UserID(); userID != 125 {
		t.FailNow()
	}
}

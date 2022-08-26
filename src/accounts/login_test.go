package accounts

import (
	"testing"

	"git.xenonstack.com/xs-onboarding/accounts/src/methods"
)

func TestLoginEndPoint(t *testing.T) {
	_, code := RegisterAccount("XSTestV1.testaccount@example.com", methods.HashForNewPassword("Test#123"), "someone", "123456789")
	if code != 200 {
		t.Errorf("expected %d , found %d", 200, code)
	}
	//case 1: when account does not exist
	_, code = LoginEndPoint("test2@example.com", "93824935405")
	if code != 404 {
		t.Errorf("expected %d , found %d", 404, code)
	}
	//case 2 : when users account is blocked
	_, code = AccountStatus("testaccount@example.com", "blocked")
	if code != 200 {
		t.Errorf("expected %d , found %d", 200, code)
	}
	_, code = LoginEndPoint("testaccount@example.com", "12345789")
	if code != 400 {
		t.Errorf("expected %d , found %d", 400, code)
	}
	//case 3: if password entered is incorrect
	_, code = AccountStatus("testaccount@example.com", "active")
	if code != 200 {
		t.Errorf("expected %d , found %d", 200, code)
	}
	_, code = LoginEndPoint("testaccount@example.com", "98765432")
	if code != 400 {
		t.Errorf("expected %d , found %d", 400, code)
	}
	//case 4: when password and email entered are correct
	mapd, code := LoginEndPoint("testaccount@example.com", "Test#123")
	if code != 200 {
		t.Error(mapd)
		t.Errorf("expected %d , found %d", 200, code)
	}
}

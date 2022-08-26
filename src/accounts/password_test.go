package accounts

import (
	"testing"

	"git.xenonstack.com/xs-onboarding/accounts/src/methods"
)

func TestForgotPassword(t *testing.T) {
	_, code := RegisterAccount("XSTestV1.forgetpassword@example.com", methods.HashForNewPassword("Test#123"), "someone", "123456789")
	if code != 200 {
		t.Errorf("expected %d , found %d", 200, code)
	}
	//case 1 : when account does not exist
	_, code = ForgotPassword("someemail@example.com")
	if code != 400 {
		t.Errorf("expected %d , found %d", 400, code)
	}
	//case 2: when account is not active
	_, code = AccountStatus("forgetpassword@example.com", "blocked")
	if code != 200 {
		t.Errorf("expected %d , found %d", 200, code)
	}
	_, code = ForgotPassword("forgetpassword@example.com")
	if code != 400 {
		t.Errorf("expected %d , found %d", 400, code)
	}
	//case 3: when account is active
	_, code = AccountStatus("forgetpassword@example.com", "active")
	if code != 200 {
		t.Errorf("expected %d , found %d", 200, code)
	}
	_, code = ForgotPassword("forgetpassword@example.com")
	if code != 500 {
		t.Errorf("expected %d , found %d", 500, code)
	}
}

func TestSetNewPassword(t *testing.T) {
	_, code := RegisterAccount("XSTestV1.setpassword@example.com", methods.HashForNewPassword("Test#123"), "someone", "123456789")
	if code != 200 {
		t.Errorf("expected %d , found %d", 200, code)
	}
	//case 1 : when account does not exist
	_, code = SetNewPassword("someemail@example.com", "9876543210")
	if code != 400 {
		t.Errorf("expected %d , found %d", 400, code)
	}
	//case 2: when account is not active
	_, code = AccountStatus("setpassword@example.com", "blocked")
	if code != 200 {
		t.Errorf("expected %d , found %d", 200, code)
	}
	_, code = SetNewPassword("setpassword@example.com", "9876543210")
	if code != 400 {
		t.Errorf("expected %d , found %d", 400, code)
	}
	//case 3: when account is active
	_, code = AccountStatus("setpassword@example.com", "active")
	if code != 200 {
		t.Errorf("expected %d , found %d", 200, code)
	}
	_, code = SetNewPassword("setpassword@example.com", "9876543210")
	if code != 200 {
		t.Errorf("expected %d , found %d", 500, code)
	}
}

func TestChangePassword(t *testing.T) {
	_, code := RegisterAccount("XSTestV1.testchangepassword@example.com", methods.HashForNewPassword("Test#123"), "someone", "123456789")
	if code != 200 {
		t.Errorf("expected %d , found %d", 200, code)
	}
	//case 1: when old password entered is right
	_, code = ChangePassword("Test#123", "test#123", "testchangepassword@example.com")
	if code != 200 {
		t.Errorf("expected %d , found %d", 200, code)
	}

	//case 2: when account does not exist
	_, code = ChangePassword("123456789", "9876543210", "randomemail@example.com")
	if code != 400 {
		t.Errorf("expected %d , found %d", 400, code)
	}

	//case 3: when old password entered is wrong
	_, code = ChangePassword("123456789", "9876543210", "testchangepassword@example.com")
	if code != 400 {
		t.Errorf("expected %d , found %d", 400, code)
	}
}

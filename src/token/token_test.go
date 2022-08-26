package token

import (
	"net/http/httptest"
	"testing"

	"git.xenonstack.com/xs-onboarding/accounts/config"
	"github.com/gin-gonic/gin"
)

func TestToken(t *testing.T) {
	config.Conf.Service.InviteLinkTimeout = 300
	config.Conf.Service.VerifyLinkTimeout = 10
	//Generate Token
	tok := GenerateToken("someone", "invite", "L1")
	//when token exists
	_, err := VerifyToken(tok, "invite")
	if err != nil {
		t.Error(err)
	}
	//token does not exists
	DeleteToken(tok)
	_, err = VerifyToken(tok, "invite")
	if err == nil {
		t.Error(err)
	}
	//reset-password
	tok = GenerateToken("someone", "reset-password", "L1")
	_, err = VerifyToken(tok, "reset-password")
	if err != nil {
		t.Error(err)
	}
	//DeleteToken
	DeleteToken(tok)
	//SaveSession
	var start, end int64
	start = 987
	end = 123
	info := make(map[string]interface{})
	info["start"] = start
	info["end"] = end
	SaveSessions(1, "sometoken", info)

}
func TestUnauthorisedFunction(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	unauthorizedFunc(c, 200, "some msg")
}

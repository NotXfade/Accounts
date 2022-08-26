package token

import (
	"errors"
	"log"
	"strconv"
	"time"

	"git.xenonstack.com/xs-onboarding/accounts/config"
	"git.xenonstack.com/xs-onboarding/accounts/database"
	"git.xenonstack.com/xs-onboarding/accounts/src/methods"
)

//==================================== GENERATE TOKEN =====================================

// GenerateToken is a method to generate new verification token and add to database
func GenerateToken(email, task,level string) string {
	// connecting to db
	db := config.DB
	token := database.Tokens{}
	// setting token user id
	token.Email = email
	token.Token = methods.RandomString(32)
	token.TokenTask = task
	token.Status = "active"
	token.Level = level
	token.Timestamp = time.Now().Unix()
	token.CreateAt = time.Now()
	token.UpdatedAt = time.Now()

	// save data in db
	db.Create(&token)
	return token.Token
}

//====================================== VERIFY TOKEN =========================================

//VerifyToken is a method to check token is valid or not
func VerifyToken(token, task string) (database.Tokens, error) {
	// connecting to db
	db := config.DB
	// fetch token details
	tok := []database.Tokens{}
	db.Where("token=? AND status=? AND token_task=?", token, "active", task).Find(&tok)
	//token not found
	if len(tok) == 0 {
		log.Println("token not found")
		return database.Tokens{}, errors.New("Invalid or expired token")
	}
	if task == "reset-password" {
		//expired token
		if (time.Now().Unix() - tok[0].Timestamp) > config.Conf.Service.VerifyLinkTimeout {
			log.Println("expired token")
			return database.Tokens{}, errors.New("Invalid or expired token")
		}
	}
	//expired token
	if (time.Now().Unix() - tok[0].Timestamp) > config.Conf.Service.InviteLinkTimeout {
		log.Println("expired token")
		return database.Tokens{}, errors.New("Invalid or expired token")
	}
	return tok[0], nil
}

//====================================== DELETE TOKEN ======================================

//DeleteToken is a method used to delete token entry from database
func DeleteToken(token string) {
	// connecting to db
	db := config.DB
	// delete used token
	row := db.Where("token=?", token).Delete(&database.Tokens{}).RowsAffected
	log.Println(row)
	// delete expired tokens
	row = db.Where("timestamp < ?", strconv.FormatInt((time.Now().Unix()-config.Conf.Service.InviteLinkTimeout), 10)).RowsAffected
	log.Println(row)
}

//=================================== GENERATE JWT TOKENS ================================================

//GenerateJwtToken is used to generate JWT Token with Accounts info claims
func GenerateJwtToken(account database.Accounts) map[string]interface{} {
	claims := make(map[string]interface{})
	claims["id"] = account.ID
	claims["name"] = account.Name
	claims["email"] = account.Email
	claims["role"] = account.Role
	claims["level"] = account.Level
	mapd, info := GinJwtToken(claims)
	// check token is empty or not
	if mapd["token"].(string) == "" {
		return mapd
	}
	// remove all other sessions from session storage and save this session
	SaveSessions(account.ID, mapd["token"].(string), info)
	return mapd
}

//==================================== SAVE SESSION OF TOKEN ===================================

// SaveSessions is a method for saving session details in database
func SaveSessions(id int, newSessToken string, info map[string]interface{}) {

	db := config.DB
	// deleting other active sessions of that user
	if config.Conf.Service.IsLogoutOthers == "true" {
		// fetch active session from dbs
		var actses []database.ActiveSessions
		db.Where("userid=?", id).Find(&actses)

		// delete all session from db
		db.Where("userid=?", id).Delete(&database.ActiveSessions{})
	}

	// creating one active session

	db.Create(&database.ActiveSessions{
		Userid:    id,
		SessionID: newSessToken,
		Start:     info["start"].(int64),
		End:       info["end"].(int64)})
}

//============================== JWT REFRESH TOKEN ========================================

// JwtRefreshToken is a method for save old claims in a token
// and also save sessions in cockroach database and redis database
func JwtRefreshToken(claims map[string]interface{}) map[string]interface{} {

	// generate jwt token, expiration time and extra info like (expire jwt time, start and end time)
	mapd, info := GinJwtToken(claims)

	// check token is empty or not
	if mapd["token"].(string) == "" {
		return mapd
	}
	// remove all other sessions from session storage and save this session
	SaveSessions(methods.ConvertID(claims["id"]), mapd["token"].(string), info)

	return mapd
}

//======================== Delete JWT SESSION ====================================================

// DeleteTokenFromDb is a method to delete saved jwt token from db
func DeleteTokenFromDb(token string) {
	db := config.DB
	db.Exec("delete from active_sessions where session_id= '" + token + "';")
}

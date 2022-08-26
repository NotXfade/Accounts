package googlecalendar

import (
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"

	"git.xenonstack.com/xs-onboarding/accounts/src/methods"

	"git.xenonstack.com/xs-onboarding/accounts/config"
	"git.xenonstack.com/xs-onboarding/accounts/database"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

var (
	googleConfig *oauth2.Config
)

//State sturcture for login and redirect
type State struct {
	Redirect string `json:"redirect"`
}

//GoogleLogin is used to login using google
func GoogleLogin(state State) string {
	googleConfig = &oauth2.Config{
		ClientID:     config.Conf.Google.ClientID,
		ClientSecret: config.Conf.Google.ClientKey,
		Scopes:       strings.Split(config.Conf.Google.Scopes, ","),
		RedirectURL:  config.Conf.Google.Redirect,
		Endpoint:     google.Endpoint,
	}
	data, _ := json.Marshal(state)
	url := googleConfig.AuthCodeURL(string(data), oauth2.AccessTypeOffline)
	return url
}

//GoogleCallback is used to
func GoogleCallback(code, state string) string {
	//log.Println(code, state)
	oToken, err := googleConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Println(err)
		return ""
	}
	if oToken.RefreshToken != "" {
		db := config.DB
		token := database.SocialToken{}
		err := db.Where("id=?", 1).Find(&token).Error
		if err != nil {
			log.Println(err)
		}
		token.AccessToken = oToken.AccessToken
		token.TokenType = oToken.TokenType
		token.RefreshToken = oToken.RefreshToken
		token.Expiry = oToken.Expiry
		token.CreatedAt = time.Now()
		token.UpdatedAt = time.Now()
		err = db.Save(&token).Error
		if err != nil {
			log.Println(err)
		}
	}
	return "https://" + config.Conf.Service.HostAddr
}

//GetToken is used to get token from db
func GetToken() (*calendar.Service, error) {
	db := config.DB
	token := database.SocialToken{}
	err := db.Where("id=?", 1).Find(&token).Error
	if err != nil {
		return &calendar.Service{}, err
	}
	if token.RefreshToken == "" {
		return &calendar.Service{}, errors.New("Credentials not found")
	}
	oToken := &oauth2.Token{
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
	}
	ctx := context.Background()
	googleConfig = &oauth2.Config{
		ClientID:     config.Conf.Google.ClientID,
		ClientSecret: config.Conf.Google.ClientKey,
		Scopes:       strings.Split(config.Conf.Google.Scopes, ","),
		RedirectURL:  config.Conf.Google.Redirect,
		Endpoint:     google.Endpoint,
	}
	srv, err := calendar.NewService(ctx, option.WithTokenSource(googleConfig.TokenSource(ctx, oToken)))
	return srv, err
}

func newBool(val bool) *bool {
	return &val
}

//CreateInvite is used to create calendar Invite
func CreateInvite(inviteData database.InviteEmail) (map[string]interface{}, int) {
	mapd := make(map[string]interface{})
	srv, err := GetToken()
	if err != nil {
		mapd["error"] = true
		mapd["message"] = "Internal server error"
		return mapd, 500
	}
	//log.Println(srv)
	attendees := []*calendar.EventAttendee{}
	for i := 0; i < len(inviteData.Email); i++ {
		attendees = append(attendees, &calendar.EventAttendee{
			Email: inviteData.Email[i],
		})
	}
	event := &calendar.Event{
		Summary:     inviteData.Title,
		Description: inviteData.Description,
		Start: &calendar.EventDateTime{
			DateTime: time.Unix(inviteData.StartTime, 0).Format(time.RFC3339),
		},

		End: &calendar.EventDateTime{
			DateTime: time.Unix(inviteData.EndTime, 0).Format(time.RFC3339),
		},
		Attendees: attendees,
		ConferenceData: &calendar.ConferenceData{
			CreateRequest: &calendar.CreateConferenceRequest{
				ConferenceSolutionKey: &calendar.ConferenceSolutionKey{
					Type: "hangoutsMeet",
				},
				RequestId: methods.RandomStringIntegerOnly(10),
			},
			EntryPoints: []*calendar.EntryPoint{
				&calendar.EntryPoint{
					EntryPointType: "video",
				},
			},
		},
		GuestsCanSeeOtherGuests: newBool(false),
		GuestsCanModify:         false,
		GuestsCanInviteOthers:   newBool(false),
	}
	calendarID := "primary"
	//log.Println(oToken, event)
	event, err = srv.Events.Insert(calendarID, event).SendUpdates("all").ConferenceDataVersion(1).Do()
	if err != nil {
		log.Printf("Unable to create event. %v\n", err)
		mapd["error"] = true
		mapd["message"] = "Unable to create calendar invite, Please try again with correct data."
		return mapd, 400
	}

	mapd["error"] = false
	mapd["message"] = "Calendar invite created successfully"
	return mapd, 200
}

//CheckCredentials is used to check if credentials exist
func CheckCredentials() (map[string]interface{}, int) {
	mapd := make(map[string]interface{})
	_, err := GetToken()
	if err != nil {
		mapd["error"] = true
		mapd["message"] = "Please signin to continue"
		return mapd, 404
	}
	mapd["error"] = false
	mapd["message"] = "Operation Successful"
	return mapd, 200
}

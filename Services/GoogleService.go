package Services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"SimpleGo_xpns.go/Models"
	"SimpleGo_xpns.go/Utilities"
	Executer "SimpleGo_xpns.go/Utilities/APIExecuter"
	dbConnector "SimpleGo_xpns.go/Utilities/DbConnector"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/go-chi/chi"

	gmail "google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

var client *http.Client

func RegisterGoogleAPIs(r chi.Router) {
	r.Get("/auth/callback", GoogleCallback)
	r.Get("/googleIn", RedirectGoogle)
	r.Get("/getLabels", GetGmailLabels)
}

func GoogleCallback(w http.ResponseWriter, r *http.Request) {
	response := Utilities.GetResponse()
	oauthConfig := initializeOauthConfig()
	raw := r.FormValue("code")
	token, err := oauthConfig.Exchange(context.Background(), raw)
	if err != nil {
		response.JSON(w, http.StatusInternalServerError, fmt.Sprintf("Not able to serve the request."))
		return
	}
	client = oauthConfig.Client(context.Background(), token)

	//Getting user info and saving in the db
	user := getUserInfo(token.AccessToken)
	if user != nil {
		//Mapping token in diff model struct
		tokenData := Utilities.MapTokenResponse(user.ID, *token)
		if tokenData != nil {
			if dbConnector.InsertUserData(*user, *tokenData) {
				//returning valid response to channel for further use of token
				response.JSON(w, http.StatusOK, fmt.Sprintf("%v", Models.SignInResponse{Id: user.ID, Name: user.Name, Token: token.AccessToken, Email: user.Email}))
				return
			}
		}
	} else {
		response.JSON(w, http.StatusInternalServerError, fmt.Sprintf("Not able to Get user Info"))
		return
	}
	return
}

func initializeOauthConfig() *oauth2.Config {
	b, err := os.ReadFile("GoogleAuth.json") //can be replace with viper get config and convert it into json.
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope, "https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile")
	return config
}

func RedirectGoogle(w http.ResponseWriter, r *http.Request) {
	oauthConfig := initializeOauthConfig()
	url := oauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, 301)

}

func GetGmailLabels(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	user := "me"
	rp, _ := srv.Users.Labels.List(user).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve labels: %v", err)
	}
	if len(rp.Labels) == 0 {
		fmt.Println("No labels found.")
		return
	}
	fmt.Println("Labels:")
	for _, l := range rp.Labels {
		fmt.Printf("- %s\n", l.Name)
	}
}

func getUserInfo(accessToken string) *Models.User {
	var req Executer.APIRequest
	req.BaseURL = "https://www.googleapis.com/"
	req.Action = fmt.Sprintf("oauth2/v1/userinfo?alt=json&access_token=%v", accessToken)
	resp := Executer.GET(req)
	if resp.Status == 200 {
		var user *Models.User
		err := json.Unmarshal([]byte(resp.Response), &user)
		if err == nil {
			return user
		}
	}
	return nil
}

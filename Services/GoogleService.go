package Services

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"SimpleGo_xpns/Models"
	"SimpleGo_xpns/Utilities"
	dbConnector "SimpleGo_xpns/Utilities/DbConnector"
	workflow "SimpleGo_xpns/Workflows"

	"github.com/go-chi/chi"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	gmail "google.golang.org/api/gmail/v1"
)

func RegisterGoogleAPIs(r chi.Router) {
	r.Get("/auth/callback", GoogleCallback)
	r.Get("/googleIn", RedirectGoogle)
	//r.Get("/SyncUp", SyncMail)
}

var oauthConfig *oauth2.Config

func initializeOauthConfig() *oauth2.Config {
	b, err := os.ReadFile("GoogleAuth.json") //can be replace with viper get config and convert it into json.
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope, "https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile")
	return config
}

//"/auth/callback"
func GoogleCallback(w http.ResponseWriter, r *http.Request) {
	response := Utilities.GetResponse()
	oauthConfig = initializeOauthConfig()
	raw := r.FormValue("code")
	token, err := oauthConfig.Exchange(context.Background(), raw)
	if err != nil {
		response.JSON(w, http.StatusInternalServerError, fmt.Sprintf("Not able to serve the request."))
		return
	}

	//Getting user info and saving in the db
	user := workflow.GetUserInfo(token.AccessToken)
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

//"/googleIn"
func RedirectGoogle(w http.ResponseWriter, r *http.Request) {
	oauthConfig := initializeOauthConfig()
	url := oauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, 301)
}

// /"/SyncUp"
// func SyncMail(w http.ResponseWriter, r *http.Request) {
// 	var client *http.Client

// 	var label string
// 	if r.URL.Query().Get("labels") != "" {
// 		label = r.URL.Query().Get("labels")
// 	}

// 	var accessToken string
// 	if r.Header.Get("authToken") != "" {
// 		accessToken = r.Header.Get("authToken")
// 	}

// 	if workflow.VerifyIdToken(accessToken) {
// 		token := dbConnector.GetUserToken(accessToken)
// 		client = oauthConfig.Client(context.Background(), token)

// 		ctx := context.Background()

// 		srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
// 		if err != nil {
// 			log.Fatalf("Unable to retrieve Gmail client: %v", err)
// 		}

// 		//fetching the latest db historyId

// 		lasthistoryid := dbConnector.GetlastHistoryByLabel(label)
// 		var k *gmail.ListThreadsResponse
// 		if len(lasthistoryid) == 0 {
// 			//list all history and fetch all latest messages
// 			workflow.FetchAllThread(label, srv)
// 			return

// 		} else {
// 			//get lastest history and fetch lastest messages
// 			//use users.history.list
// 		}

// 		if err != nil {
// 			log.Fatalf("Unable to retrieve labels: %v", err)
// 		}
// 		user := "me"
// 		threadservice := srv.Users.Threads.List(user)
// 		threadservice.Q(fmt.Sprintf("label:%v", label))
// 		k, err = threadservice.Do()
// 		//looping through threads
// 		for _, x := range k.Threads {
// 			//getting one thread  using id
// 			threadservice := srv.Users.Threads.Get(user, x.Id)
// 			th, _ := threadservice.Do()

// 			//looping over mails
// 			for _, msg := range th.Messages {

// 				msgId := msg.Id
// 				attId := msg.Payload.Parts[1].Body.AttachmentId

// 				msgService := srv.Users.Messages.Attachments.Get(user, msgId, attId)
// 				u, _ := msgService.Do()
// 				decoded, _ := base64.URLEncoding.DecodeString(u.Data)
// 				d := bytes.NewReader(decoded)
// 				req := workflow.SBIparser("", d)
// 				fmt.Printf("%v", req)
// 				//send data to postgres
// 				// isSuccess := dbConnector.SendDataToPostgres(req)
// 				// if isSuccess {
// 				// 	fmt.Println("sucessfully inserted user record to db")
// 				// }
// 			}
// 		}
// 	} else {
// 		http.Redirect(w, r, "localhost:9005/googleIn", http.StatusAccepted)
// 	}
// }

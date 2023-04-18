package Services

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"SimpleGo_xpns/Models"
	"SimpleGo_xpns/Utilities"
	dbConnector "SimpleGo_xpns/Utilities/DbConnector"
	workflow "SimpleGo_xpns/Workflows"

	"github.com/go-chi/chi"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	gmail "google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

func RegisterGoogleAPIs(r chi.Router) {
	r.Get("/auth/callback", GoogleCallback)
	r.Get("/SignIn", RedirectGoogle)
	r.Get("/SyncMail", SyncMail)
}

var oauthConfig *oauth2.Config
var user = "me"

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
				response.JSON(w, http.StatusOK, Models.SignInResponse{Id: user.ID, Name: user.Name, Token: token.AccessToken, Email: user.Email})
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
	fmt.Println("Redirection eurl from config : ", oauthConfig.RedirectURL)
	url := oauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	http.Redirect(w, r, url, 301)
}

//"/SyncUp"
func SyncMail(w http.ResponseWriter, r *http.Request) {
	var client *http.Client
	response := Utilities.GetResponse()
	var label, from, to, query string
	var accessToken string

	if r.Header.Get("authToken") != "" {
		accessToken = r.Header.Get("authToken")
	}

	if r.URL.Query().Get("label") != "" {
		label = r.URL.Query().Get("label")
	} else {
		response.JSON(w, http.StatusBadRequest, Models.BaseResponse{Status: false, Error: "Please provide label to fetch"})
		return
	}
	from = r.URL.Query().Get("from")
	to = r.URL.Query().Get("to")

	if workflow.VerifyIdToken(accessToken) {
		token := dbConnector.GetUserToken(accessToken)
		client = oauthConfig.Client(context.Background(), token)

		ctx := context.Background()
		var k *gmail.ListThreadsResponse

		srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
		if err != nil {
			log.Fatalf("Unable to retrieve Gmail client: %v", err)
		}

		if err != nil {
			log.Fatalf("Unable to retrieve labels: %v", err)
		}

		threadservice := srv.Users.Threads.List(user)

		if from != "" && to != "" {
			query = fmt.Sprintf("label:%v after:%v before:%v", label, to, from)
		} else {
			//fetch all
			query = fmt.Sprintf("label:%v", label)
		}

		threadservice.Q(query)
		k, err = threadservice.Do()

		//getting all the msgs encoded data
		encodedReq, timestampmap := getDecodeFromThread(k, srv)
		decodedData := workflow.GetDataForbase64(encodedReq)

		for _, f := range decodedData {
			f.ETime = timestampmap[f.TransactionId] //Note : using pointer to updating struct
		}

		//send data to db with merging SBI and HDFC records in postgres
		//TODO
		//Currently creating another table
		if decodedData != nil {
			failedTxns := dbConnector.SendHDFCToPostgres(decodedData)
			if len(failedTxns) > 0 {
				response.JSON(w, http.StatusOK, len(failedTxns))
				return
			} else {
				response.JSON(w, http.StatusOK, "sucesss")
				return
			}
		}

	} else {
		http.Redirect(w, r, "localhost:9005/SignIn", http.StatusAccepted)
	}
}

func getDecodeFromThread(k *gmail.ListThreadsResponse, srv *gmail.Service) ([]Models.GetEncodedDataReq, map[string]time.Time) {
	var encodedReq []Models.GetEncodedDataReq
	timestampmap := make(map[string]time.Time)

	for _, x := range k.Threads {
		//getting one thread  using id
		threadservice := srv.Users.Threads.Get(user, x.Id)
		th, _ := threadservice.Do()

		//looping over mails over a thread
		for _, msg := range th.Messages {
			var encodedData string
			if len(msg.Payload.Parts) > 0 {
				encodedData = msg.Payload.Parts[0].Body.Data
			} else {
				encodedData = msg.Payload.Body.Data
			}
			encodedReq = append(encodedReq, Models.GetEncodedDataReq{
				MsgEncodedData: encodedData,
				MsgId:          msg.Id,
			})
			loc, _ := time.LoadLocation("Asia/Kolkata")
			t := time.UnixMilli(msg.InternalDate).In(loc).Format("2-Jan-06 15:04:05")
			timestampmap[msg.Id], _ = time.Parse("2-Jan-06 15:04:05", t)
		}
	}
	return encodedReq, timestampmap
}

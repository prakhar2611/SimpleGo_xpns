package Services

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	gmail "google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

var client *http.Client

func RegisterGoogleAPIs(r chi.Router) {
	r.Get("/auth/callback", GoogleCallback)
	//r.Get("/auth", GoogleSignIn)
	//r.Get("/", GoogleIndex)
	r.Get("/newauth", RedirectGoogle)
	r.Get("/getLabels", GetGmailLabels)
}

// func GoogleSignIn(res http.ResponseWriter, req *http.Request) {
// 	IntializeConfigs()
// 	gothic.BeginAuthHandler(res, req)
// }

// func GoogleIndex(res http.ResponseWriter, req *http.Request) {
// 	IntializeConfigs()
// 	t, _ := template.ParseFiles("templates/index.html")
// 	t.Execute(res, false)
// }

func GoogleCallback(w http.ResponseWriter, r *http.Request) {
	oauthConfig := initializeOauthConfig()
	raw := r.FormValue("code")
	token, err := oauthConfig.Exchange(context.Background(), raw)
	if err != nil {
		// Handle error
	}
	client = oauthConfig.Client(context.Background(), token)

	fmt.Println(raw)
	fmt.Println("========================")
	fmt.Println(token)

}

func initializeOauthConfig() *oauth2.Config {
	// oauthConfig := &oauth2.Config{
	// 	ClientID:     "64464811543-fee5m8plhj94lpv9vgcei91r15189b45.apps.googleusercontent.com",
	// 	ClientSecret: "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",
	// 	RedirectURL:  "http://localhost:9005/auth/callback?provider=google",
	// 	Scopes: []string{
	// 		"https://www.googleapis.com/auth/userinfo.email",
	// 		"https://www.googleapis.com/auth/userinfo.profile",
	// 		"https://www.googleapis.com/auth/gmail.readonly",
	// 	},
	// 	Endpoint: google.Endpoint,
	// }
	// := context.Background()
	b, err := os.ReadFile("GoogleAuth.json")
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

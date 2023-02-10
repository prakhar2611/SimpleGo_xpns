package Services

import (
	"SimpleGo_xpns/Utilities"
	"context"
	"fmt"
	"net/http"
	"text/template"

	"github.com/go-chi/chi"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"go.mongodb.org/mongo-driver/bson"
)

func RegisterGoogleAPIs(r chi.Router) {
	r.Get("/auth/callback", GoogleCallback)
	r.Get("/auth", GoogleSignIn)
	r.Get("/", GoogleIndex)
}

func GoogleSignIn(res http.ResponseWriter, req *http.Request) {
	IntializeConfigs()
	gothic.BeginAuthHandler(res, req)
}

func GoogleIndex(res http.ResponseWriter, req *http.Request) {
	IntializeConfigs()
	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(res, false)
}

func GoogleCallback(res http.ResponseWriter, req *http.Request) {
	IntializeConfigs()
	user, err := gothic.CompleteUserAuth(res, req)
	if err != nil {
		fmt.Fprintln(res, err)
		return
	}
	client, ctx := Utilities.ConectToMongo()
	col := client.Database("ExpenseDb").Collection("Users")
	r, err := col.InsertOne(context.Background(), bson.M{"Id": user.UserID, "Email": user.Email})
	if err != nil {
		fmt.Println("Error occured while Inserting data -", err, ctx)
	}
	id := r.InsertedID
	fmt.Println("Data Inserted with index -", id)
	t, _ := template.ParseFiles("templates/success.html")
	t.Execute(res, user)
}

func IntializeConfigs() {
	key := "Secret-session-key" // Replace with your SESSION_SECRET or similar
	maxAge := 86400 * 30        // 30 days
	isProd := false             // Set to true when serving over https

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true // HttpOnly should always be enabled
	store.Options.Secure = isProd

	gothic.Store = store

	goth.UseProviders(
		google.New("64464811543-fee5m8plhj94lpv9vgcei91r15189b45.apps.googleusercontent.com", "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", "http://localhost:9005/auth/callback?provider=google", "email", "profile"),
	)
}

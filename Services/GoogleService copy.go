package Services

// import (
// 	"github.com/go-chi/chi"
// )

// func RegisterGoogleAPIs(r chi.Router) {
// 	//r.Get("/auth/callback", GoogleCallback)
// 	//r.Get("/auth", GoogleSignIn)
// 	//r.Get("/", GoogleIndex)
// 	//r.Get("/newauth", RedirectGoogle)
// }

// // func GoogleSignIn(res http.ResponseWriter, req *http.Request) {
// // 	IntializeConfigs()
// // 	gothic.BeginAuthHandler(res, req)
// // }

// // func GoogleIndex(res http.ResponseWriter, req *http.Request) {
// // 	IntializeConfigs()
// // 	t, _ := template.ParseFiles("templates/index.html")
// // 	t.Execute(res, false)
// // }

// func GoogleCallback(res http.ResponseWriter, req *http.Request) {
// 	IntializeConfigs()
// 	user, err := gothic.CompleteUserAuth(res, req)

// 	if err != nil {
// 		fmt.Fprintln(res, err)
// 		return
// 	}
// 	client, ctx := Utilities.ConectToMongo()
// 	col := client.Database("ExpenseDb").Collection("Users")
// 	r, err := col.InsertOne(context.Background(), bson.M{"Id": user.UserID, "Email": user.Email})
// 	if err != nil {
// 		fmt.Println("Error occured while Inserting data -", err, ctx)
// 	}
// 	id := r.InsertedID
// 	fmt.Println("Data Inserted with index -", id)
// 	t, _ := template.ParseFiles("templates/success.html")
// 	t.Execute(res, user)
// }

// func IntializeConfigs() {
// 	key := "Secret-session-key" // Replace with your SESSION_SECRET or similar
// 	maxAge := 86400 * 30        // 30 days
// 	isProd := false             // Set to true when serving over https

// 	store := sessions.NewCookieStore([]byte(key))
// 	store.MaxAge(maxAge)
// 	store.Options.Path = "/"
// 	store.Options.HttpOnly = true // HttpOnly should always be enabled
// 	store.Options.Secure = isProd

// 	gothic.Store = store

// 	goth.UseProviders(
// 	//google.New("64464811543-fee5m8plhj94lpv9vgcei91r15189b45.apps.googleusercontent.com", "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", "http://localhost:9005/auth/callback?provider=google", "email", "profile", "https://www.googleapis.com/auth/gmail.readonly"),
// 	)
// }

// // func initializeOauthConfig() *oauth2.Config {
// // 	oauthConfig := &oauth2.Config{
// // 		ClientID:     "64464811543-fee5m8plhj94lpv9vgcei91r15189b45.apps.googleusercontent.com",
// // 		ClientSecret: "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",
// // 		RedirectURL:  "http://localhost:9005/auth/callback?provider=google",
// // 		Scopes: []string{
// // 			"https://www.googleapis.com/auth/userinfo.email",
// // 			"https://www.googleapis.com/auth/userinfo.profile",
// // 			"https://www.googleapis.com/auth/gmail.readonly",
// // 		},
// // 		Endpoint: google.Endpoint,
// // 	}

// // 	return oauthConfig
// // }

// // func RedirectGoogle(w http.ResponseWriter, r *http.Request) {
// // 	oauthConfig := initializeOauthConfig()
// // 	url := oauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

// // 	var req Utilities.APIRequest
// // 	req.BaseURL = url
// // 	res := Utilities.GET(req)

// // 	fmt.Println(res)

// // 	http.Redirect(w, r, url, http.StatusFound)

// // 	token, err := oauthConfig.Exchange(context.Background(), r.FormValue("code"))
// // 	if err != nil {
// // 		// Handle error
// // 	}

// // 	client := oauthConfig.Client(context.Background(), token)
// // 	srv, err := gmail.NewService(context.Background(), option.WithHTTPClient(client))
// // 	if err != nil {
// // 		log.Fatalf("Unable to retrieve Gmail client: %v", err)
// // 	}

// // 	user := "me"
// // 	ri, err := srv.Users.Labels.List(user).Do()
// // 	if err != nil {
// // 		log.Fatalf("Unable to retrieve labels: %v", err)
// // 	}
// // 	if len(ri.Labels) == 0 {
// // 		fmt.Println("No labels found.")
// // 		return
// // 	}
// // 	fmt.Println("Labels:")
// // 	for _, l := range ri.Labels {
// // 		fmt.Printf("- %s\n", l.Name)
// // 	}
// // }

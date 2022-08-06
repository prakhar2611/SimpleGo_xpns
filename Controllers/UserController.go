package Controllers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

func RegisterUserAPI(r chi.Router) {
	r.Get("/api/User/v1/GetUser", GetUserDetails)
	r.Handle("/api/testuser", http.HandlerFunc(TestUser))
}

func GetUserDetails(w http.ResponseWriter, r *http.Request) {
	//var d []User
	//d := UserModel.User{}
}

func TestUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("the server is running fine in local host")
}

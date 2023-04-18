package Controllers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

func RegisterUserAPI(r chi.Router) {
	r.Get("/api/User/v1/GetUser", GetUserDetails)
	r.Post("/api/User/v1/Signin", SignIn)
}

func GetUserDetails(w http.ResponseWriter, r *http.Request) {
	//var d []User
	//d := UserModel.User{}
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	fmt.Println("the server is running fine in local host")

}

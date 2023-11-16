package Controllers

import (
	"SimpleGo_xpns/Models"
	"SimpleGo_xpns/Utilities"
	dbConnector "SimpleGo_xpns/Utilities/DbConnector"
	workflow "SimpleGo_xpns/Workflows"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"golang.org/x/oauth2"
)

func RegisterUserAPI(r chi.Router) {
	r.Get("/api/User/v1/GetUserProfile", GetUserProfile)
	r.Post("/api/User/v1/Signin", SignIn)
}

func GetUserProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	response := Utilities.GetResponse()
	var userProfile Models.UserProfile
	var token string
	if r.Header.Get("token") != "" {
		token = r.Header.Get("token")
	}
	if token != "" {
		user := workflow.GetUserInfo(token)
		if user != nil {
			userProfile.BaseResponse = Models.BaseResponse{Status: true, Error: ""}
			userProfile.Email = user.Email
			userProfile.Picture = user.Picture
			userProfile.Name = user.Name
			userProfile.Userid = user.ID

			response.JSON(w, http.StatusOK, userProfile)
			return
		} else {
			userProfile.BaseResponse = Models.BaseResponse{Status: false, Error: "Unable to retrive data"}
			response.JSON(w, http.StatusUnauthorized, userProfile)
			return
		}
	} else {
		userProfile.BaseResponse = Models.BaseResponse{Status: false, Error: "Unable to retrive data"}
		response.JSON(w, http.StatusBadRequest, userProfile)
		return
	}
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	response := Utilities.GetResponse()

	var token *oauth2.Token
	err := json.NewDecoder(r.Body).Decode(&token)
	if err == nil {
		user := workflow.GetUserInfo(token.AccessToken)
		isCache := Utilities.SetKey(user.ID, token)
		if user != nil {
			if dbConnector.InsertUserData(*user) != 0 && isCache {
				//returning valid response to channel for further use of token
				Utilities.GetKeyValue(user.ID)
				//0 - error
				//1 - new user
				//2 - existing
				if dbConnector.InsertUserData(*user) == 2 {
					response.JSON(w, http.StatusOK, Models.BaseResponse{Status: true, IsNewUser: false, Error: ""})
					return
				} else {
					response.JSON(w, http.StatusOK, Models.BaseResponse{Status: true, IsNewUser: true, Error: ""})
					return
				}
			}
		}
	}
	response.JSON(w, http.StatusInternalServerError, Models.BaseResponse{Status: false, Error: "Unable to login at server"})
	return
}

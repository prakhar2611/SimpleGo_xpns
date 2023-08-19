package Controllers

import (
	"SimpleGo_xpns/Models"
	"SimpleGo_xpns/Utilities"
	dbConnector "SimpleGo_xpns/Utilities/DbConnector"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func RegisterDocsAPI(r chi.Router) {
	r.Get("/api/v1/listDocs", GetDocs)
	r.Post("/api/v1/save", SaveDoc)

}

func SaveDoc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	response := Utilities.GetResponse()
	var request *Models.DocsMeta

	// var token string
	// if r.Header.Get("token") != "" {
	// 	token = r.Header.Get("token")
	// }

	var new string
	if r.URL.Query().Get("new") != "" {
		new = r.URL.Query().Get("new")
	}
	token := ""
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		response.JSON(w, http.StatusUnauthorized, Models.BaseResponse{Status: false, Error: "Unauthorized, please login"})
	}
	// if r.URL.Query().Get("label") != "" {
	// 	label = r.URL.Query().Get("label")
	// } else {
	// 	response.JSON(w, http.StatusBadRequest, Models.BaseResponse{Status: false, Error: "Please provide label to fetch"})
	// 	return
	// }
	// from = r.URL.Query().Get("from")
	// to = r.URL.Query().Get("to")
	if request.Title == "" {
		response.JSON(w, http.StatusBadRequest, Models.BaseResponse{Status: false, Error: "Please provide title"})

	}
	if token == "" {
		// user := workflow.GetUserInfo(token)
		//creating uniq uuid for doc

		if new == "true" && request.ID == "" {
			uuid := uuid.New().String()
			request.ID = uuid
			if dbConnector.CreateDoc(request) {
				response.JSON(w, http.StatusOK, Models.BaseResponse{Status: true})
			} else {
				response.JSON(w, http.StatusInternalServerError, Models.BaseResponse{Status: false, Error: "not implemented !"})

			}
		}

		return
	} else {
		response.JSON(w, http.StatusUnauthorized, Models.BaseResponse{Status: false, Error: "Unauthorized, please login"})
		return
	}
}

func GetDocs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	response := Utilities.GetResponse()
	// var request *Models.DocsMeta

	// var token string
	// if r.Header.Get("token") != "" {
	// 	token = r.Header.Get("token")
	// }
	token := ""
	// err := json.NewDecoder(r.Body).Decode(&request)
	// if err != nil {
	// 	response.JSON(w, http.StatusUnauthorized, Models.BaseResponse{Status: false, Error: "Unauthorized, please login"})
	// }
	// if r.URL.Query().Get("label") != "" {
	// 	label = r.URL.Query().Get("label")
	// } else {
	// 	response.JSON(w, http.StatusBadRequest, Models.BaseResponse{Status: false, Error: "Please provide label to fetch"})
	// 	return
	// }
	// from = r.URL.Query().Get("from")
	// to = r.URL.Query().Get("to")
	// if request.Title == "" {
	// 	response.JSON(w, http.StatusBadRequest, Models.BaseResponse{Status: false, Error: "Please provide title"})

	// }
	if token == "" {
		// user := workflow.GetUserInfo(token)
		//creating uniq uuid for doc
		// uuid := uuid.New().String()
		// request.ID = uuid
		data := dbConnector.GetAllDoc()
		if data != nil {
			var res Models.GetAllDocsResponse
			res.Data = data
			res.BaseResponse.Status = true
			response.JSON(w, http.StatusOK, res)
		}
		return
	} else {
		response.JSON(w, http.StatusUnauthorized, Models.BaseResponse{Status: false, Error: "Unauthorized, please login"})
		return
	}
}

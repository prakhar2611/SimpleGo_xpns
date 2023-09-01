package Controllers

import (
	"SimpleGo_xpns/Models"
	"SimpleGo_xpns/Utilities"
	dbConnector "SimpleGo_xpns/Utilities/DbConnector"
	workflow "SimpleGo_xpns/Workflows"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func RegisterDocsAPI(r chi.Router) {
	r.Get("/api/v1/listDocs", GetDocs)
	r.Post("/api/v1/save", SaveDoc)
	r.Post("/api/v1/getDocsMeta", GetDocsMeta)

}

func SaveDoc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	response := Utilities.GetResponse()
	var request *Models.DocsMeta
	// var dbPayload *Models.DocsMeta

	var token string
	if r.Header.Get("token") != "" {
		token = r.Header.Get("token")
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		response.JSON(w, http.StatusUnauthorized, Models.BaseResponse{Status: false, Error: "Unauthorized, please login"})
	}

	if token != "" {
		user := workflow.GetUserInfo(token)
		//creating uniq uuid for doc

		if user != nil {

			if len(request.Title) > 0 && len(request.MetaData) > 0 {
				uuid := uuid.New().String()
				//mapping request to db interface
				request.ID = uuid
				request.UserId = user.ID
				if len(request.Folder) == 0 {
					request.Folder = "Others"
				}
			}

			if dbConnector.CreateDoc(request, user.ID) {
				response.JSON(w, http.StatusOK, Models.BaseResponse{Status: true})
				return
			} else {
				response.JSON(w, http.StatusInternalServerError, Models.BaseResponse{Status: false, Error: "not implemented !"})
				return
			}
		}
		response.JSON(w, http.StatusUnauthorized, Models.BaseResponse{Status: false, Error: "Unauthorized, please login"})

		return
	} else {
		response.JSON(w, http.StatusUnauthorized, Models.BaseResponse{Status: false, Error: "Unauthorized, please login"})
		return
	}
}

func GetDocs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	response := Utilities.GetResponse()

	var token string
	if r.Header.Get("token") != "" {
		token = r.Header.Get("token")
	}

	user := workflow.GetUserInfo(token)

	if user != nil {
		data := dbConnector.GetAllDoc(user.ID)
		if data != nil {
			d := mapDirectory(data)
			var res Models.GetAllDocsResponse
			res.Data = d
			res.BaseResponse.Status = true
			response.JSON(w, http.StatusOK, res)
		}
		return
	} else {
		response.JSON(w, http.StatusUnauthorized, Models.BaseResponse{Status: false, Error: "Unauthorized, please login"})
		return
	}
}

func GetDocsMeta(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	var request Models.GetDocMetaRequest
	response := Utilities.GetResponse()
	// var request *Models.DocsMeta

	var token string
	if r.Header.Get("token") != "" {
		token = r.Header.Get("token")
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, Models.BaseResponse{Status: false, Error: "request invalid"})
	}

	user := workflow.GetUserInfo(token)
	if user != nil {
		data := dbConnector.GetDocMeta(request, user.ID)
		if data != nil {
			var res Models.DocsMetaResponse
			res.Data = *data
			res.BaseResponse.Status = true
			response.JSON(w, http.StatusOK, res)
		} else {
			response.JSON(w, http.StatusUnauthorized, Models.BaseResponse{Status: false, Error: "no data found"})
			return
		}
		return
	}
	response.JSON(w, http.StatusInternalServerError, Models.BaseResponse{Status: false, Error: "Unable to fecth the data"})
	return
}

func mapDirectory(data []Models.DocsMeta) []Models.Directory {
	var directory []Models.Directory
	f := make(map[string][]string)
	for _, x := range data {
		f[x.Folder] = append(f[x.Folder], x.Title)
	}

	for key, value := range f {
		var d Models.Directory
		var cd []Models.Children

		for _, x := range value {
			var c Models.Children
			c.Title = x
			c.IsLeaf = true
			c.Folder = key
			cd = append(cd, c)
		}
		d.Title = key
		d.Children = cd

		directory = append(directory, d)
	}
	return directory
}

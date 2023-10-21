package Controllers

import (
	"SimpleGo_xpns/Models"
	Model "SimpleGo_xpns/Models"
	"SimpleGo_xpns/Utilities"
	dbConnector "SimpleGo_xpns/Utilities/DbConnector"
	"SimpleGo_xpns/Workflows"
	workflow "SimpleGo_xpns/Workflows"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

func RegisterDataAPI(r chi.Router) {
	r.Post("/api/v1/sendData", SendToMongo) //not using it for external use, same is implemented in the workflow check there
	r.Get("/api/v1/getExpense", GetXpns)
	r.Get("/api/v1/getXpnsByVpa", GetXpnsByVpa)
	//r.Get("/api/v1/getExpense", )
	r.Post("/api/v1/importFromFile", ImportFormFile)
	r.Post("/api/v1/Update", UpdateCategory)
	r.Post("/api/v1/UpdateVpaMapping", UpdateVpaMapping)
	r.Post("/api/v1/UpdatePockets", UpdatePockets)

	//r.Get("/api/v1/GetFromEmailByMonth", GetGmailByMonth)
}

func SendToMongo(w http.ResponseWriter, r *http.Request) {
	var payload Model.PayloadToMongo
	//var payload string
	response := Utilities.GetResponse()

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&payload); err != nil {
		response.JSON(w, 400, "Bad Request")
		return
	}

	c, _ := Utilities.ConectToMongo()
	col := c.Database("SimpleGo").Collection("Expenses")
	for _, x := range payload.Data {
		// t, _ := time.Parse("2006-01-02", x.Date)
		// x.Day = t.Weekday().String()
		// x.Month = t.Month().String()
		if id, err := col.InsertOne(context.Background(), x); err == nil {
			fmt.Print(id) //successfull case
		} else {
			response.JSON(w, http.StatusInternalServerError, fmt.Sprintf("Error Occured while Inserting data %s", err))
			return
		}
	}
	response.JSON(w, http.StatusOK, fmt.Sprintf("Data Imported !!"))
}

func UpdateCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	response := Utilities.GetResponse()

	var token string

	var request *Models.UpdatecategoryPayload
	err := json.NewDecoder(r.Body).Decode(&request)
	if r.Header.Get("token") != "" && request != nil {
		token = r.Header.Get("token")
	} else {
		response.JSON(w, http.StatusBadRequest, Models.BaseResponse{Status: false, Error: "Bad request"})
		return
	}
	if err == nil {
		user := workflow.GetUserInfo(token)
		if user != nil {
			x, failureIds := dbConnector.UpdateCategory(request)
			if x {
				response.JSON(w, http.StatusOK, Models.UpdateCategoryResponse{FailureMsgId: failureIds})
				return
			}
		}
	}
	response.JSON(w, http.StatusInternalServerError, Models.BaseResponse{Status: false, Error: "Unable to serve the req"})
	return
}

func GetXpns(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	response := Utilities.GetResponse()
	var xpnsData []*Models.B64decodedResponse
	var label, from, to string

	var token string
	if r.Header.Get("token") != "" {
		token = r.Header.Get("token")
	}

	if r.URL.Query().Get("label") != "" {
		label = r.URL.Query().Get("label")
	} else {
		response.JSON(w, http.StatusBadRequest, Models.BaseResponse{Status: false, Error: "Please provide label to fetch"})
		return
	}
	from = r.URL.Query().Get("from")
	to = r.URL.Query().Get("to")

	if token != "" {
		user := workflow.GetUserInfo(token)
		if user != nil && label == "HDFC" {
			xpnsData = dbConnector.GetXpnsFromPostgres(from, to, user.ID)
		}
		if xpnsData != nil {
			response.JSON(w, http.StatusOK, xpnsData)
			return
		}
		response.JSON(w, http.StatusInternalServerError, Models.BaseResponse{Status: false, Error: "Unable to fecth the data"})
		return
	} else {
		response.JSON(w, http.StatusUnauthorized, Models.BaseResponse{Status: false, Error: "Unauthorized, please login"})
		return
	}
}

func GetXpnsByVpa(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	response := Utilities.GetResponse()
	var xpnsData []Models.VpaMapping
	var label string
	var limit string
	var offset string

	var token string
	if r.Header.Get("token") != "" {
		token = r.Header.Get("token")
	}

	if r.URL.Query().Get("label") != "" {
		label = r.URL.Query().Get("label")
	} else {
		response.JSON(w, http.StatusBadRequest, Models.BaseResponse{Status: false, Error: "Please provide label to fetch"})
		return
	}

	if r.URL.Query().Get("limit") != "" {
		limit = r.URL.Query().Get("limit")
	} else {
		response.JSON(w, http.StatusBadRequest, Models.BaseResponse{Status: false, Error: "Please provide label to fetch"})
		return
	}

	if r.URL.Query().Get("limit") != "" {
		offset = r.URL.Query().Get("offset")
	} else {
		response.JSON(w, http.StatusBadRequest, Models.BaseResponse{Status: false, Error: "Please provide label to fetch"})
		return
	}

	if token != "" {
		user := workflow.GetUserInfo(token)
		if user != nil && label == "HDFC" {
			xpnsData = dbConnector.GetGroupedVpa(limit, offset)
		}
		if len(xpnsData) > 0 {
			response.JSON(w, http.StatusOK, xpnsData)
			return
		}
		response.JSON(w, http.StatusInternalServerError, Models.BaseResponse{Status: false, Error: "Unable to fecth the data"})
		return
	} else {
		response.JSON(w, http.StatusUnauthorized, Models.BaseResponse{Status: false, Error: "Unauthorized, please login"})
		return
	}
}

func UpdateVpaMapping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	response := Utilities.GetResponse()

	var token string

	var request *Models.UpdatecategoryPayload
	err := json.NewDecoder(r.Body).Decode(&request)
	if r.Header.Get("token") != "" && request != nil {
		token = r.Header.Get("token")
	} else {
		response.JSON(w, http.StatusBadRequest, Models.BaseResponse{Status: false, Error: "Bad request"})
		return
	}
	if err == nil {
		user := workflow.GetUserInfo(token)
		if user != nil {
			//mapping vpa with
			x := workflow.UpdateVpaMapping(request, user.ID)
			if x {
				response.JSON(w, http.StatusOK, "TRANSACTION_VPA_LABEL_UPDATED")
				return
			}
		}
	}
	response.JSON(w, http.StatusInternalServerError, Models.BaseResponse{Status: false, Error: "TRANSACTION_VPA_LABEL_FAILURE"})
	return
}

func UpdatePockets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	response := Utilities.GetResponse()

	var token string

	var request *Models.UpdatePocketsPayload
	err := json.NewDecoder(r.Body).Decode(&request)
	if r.Header.Get("token") != "" && request != nil {
		token = r.Header.Get("token")
	} else {
		response.JSON(w, http.StatusBadRequest, Models.BaseResponse{Status: false, Error: "Bad request"})
		return
	}
	if err == nil {
		user := workflow.GetUserInfo(token)
		if user != nil {
			//mapping vpa with
			x := workflow.UpdatePockets(request, user.ID)
			if x {
				response.JSON(w, http.StatusOK, "TRANSACTION_POCKETS_UPDATED")
				return
			}
		}
	}
	response.JSON(w, http.StatusInternalServerError, Models.BaseResponse{Status: false, Error: "TRANSACTION_POCKETS_FAILURE"})
	return
}

//latest used code to db
func ImportFormFile(w http.ResponseWriter, r *http.Request) {
	var payload []Model.ExpenseBO
	var res Model.BaseResponse
	response := Utilities.GetResponse()

	if r.URL.Query().Get("files") != "" {
		files := r.URL.Query().Get("files")
		accountType := r.URL.Query().Get("accountType")
		toFormat := r.URL.Query().Get("toFormat")
		token := r.Header.Get("token")

		//getting user data
		user := Workflows.GetUserInfo(token)
		if user != nil {
			//based on the account type will call different workflow for the extractor
			//1.Sbiparser workflow
			//2.HDFCparser workflow
			if accountType == "SBI" {
				payload = Workflows.SBIparser(files, user.ID, nil)
			} else if accountType == "HDFC" {
				//Workflows.HDFCparser()
				response.JSON(w, http.StatusUnprocessableEntity, "Not implemented")
				return
			} else {
				res := Model.BaseResponse{
					Status: false,
					Error:  "Please select Account Name as SBI or HDFC",
				}
				response.JSON(w, http.StatusInternalServerError, res)
				return
			}
			status := false
			if len(payload) > 0 {
				//status := Workflows.SendToMongo(payload)
				if toFormat == "db" {
					status = Workflows.SendPayloadTodb(token, payload)
				} else if toFormat == "exl" {
					status = Workflows.SendPayloadToExcel(token, payload)
				}
				if status {
					response.JSON(w, http.StatusOK, fmt.Sprintf("Successfully Inserted Excel : ", files, accountType))
					return
				} else {
					response.JSON(w, http.StatusOK, fmt.Sprintf("Already Inserted  : ", files, accountType))
				}
			} else {
				res := Model.BaseResponse{
					Status: false,
					Error:  "payload is Empty",
				}
				response.JSON(w, http.StatusUnprocessableEntity, res)
			}
		} else {
			res := Model.BaseResponse{
				Status: false,
				Error:  "User Token Invalid",
			}
			response.JSON(w, http.StatusUnauthorized, res)
		}

		//fmt.Println(req)
		return
	}
	response.JSON(w, http.StatusOK, res)
	return
}

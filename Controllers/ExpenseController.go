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
	//r.Get("/api/v1/getExpense", )
	r.Post("/api/v1/importFromFile", ImportFormFile)
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

//need to check this fuction flow
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
			xpnsData = dbConnector.GetXpnsFromPostgres(from, to)
		}
		if xpnsData != nil {
			response.JSON(w, http.StatusOK, xpnsData)
			return
		}
		response.JSON(w, http.StatusInternalServerError, fmt.Sprintf("Unable to fecth the data"))
		return
	} else {
		response.JSON(w, http.StatusUnauthorized, fmt.Sprintf("Unauthorized, please login"))
		return
	}
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

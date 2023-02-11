package Controllers

import (
	Model "SimpleGo_xpns/Models"
	Utilities "SimpleGo_xpns/Utilities"
	Workflows "SimpleGo_xpns/Workflows"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

func RegisterDataAPI(r chi.Router) {
	//r.Post("/api/User/v1/SendData", SendToMongo)
	//r.Get("/api/User/v1/GetExpense", GetXpns)
	r.Post("/api/User/v1/ImportFromFile", ImportFormFile)
}

// func SendToMongo(w http.ResponseWriter, r *http.Request) {
// 	var payload Model.PayloadToMongo
// 	//var payload string
// 	response := Utilities.GetResponse()

// 	decoder := json.NewDecoder(r.Body)
// 	if err := decoder.Decode(&payload); err != nil {
// 		response.JSON(w, 400, "Bad Request")
// 		return
// 	}

// 	c, _ := Utilities.ConectToMongo()
// 	col := c.Database("SimpleGo").Collection("Expenses")
// 	for _, x := range payload.Data {
// 		t, _ := time.Parse("2006-01-02", x.Date)
// 		x.Day = t.Weekday().String()
// 		x.Month = t.Month().String()
// 		if id, err := col.InsertOne(context.Background(), x); err == nil {
// 			fmt.Print(id) //successfull case
// 		} else {
// 			response.JSON(w, http.StatusInternalServerError, fmt.Sprintf("Error Occured while Inserting data %s", err))
// 			return
// 		}
// 	}
// 	response.JSON(w, http.StatusOK, fmt.Sprintf("Data Imported !!"))

// }

// func GetXpns(w http.ResponseWriter, r *http.Request) {
// 	var userId string
// 	var res Model.PayloadToMongo
// 	if r.URL.Query().Get("UserId") != "" {
// 		userId = r.URL.Query().Get("UserId")
// 	}
// 	u, _ := strconv.ParseInt(userId, 10, 32)

// 	response := Utilities.GetResponse()
// 	c, _ := Utilities.ConectToMongo()
// 	col := c.Database("SimpleGo").Collection("Expenses")
// 	if cur, err := col.Find(context.Background(), bson.D{{"userId", u}}); err == nil {
// 		//fmt.Print(id)
// 		for cur.Next(context.TODO()) {
// 			//Create a value into which the single document can be decoded
// 			var elem Model.ExpenseBO
// 			err := cur.Decode(&elem)
// 			if err != nil {
// 				log.Fatal(err)
// 			}

// 			res.Data = append(res.Data, elem)

// 		}

// 		response.JSON(w, http.StatusInternalServerError, res)
// 		return
// 	} else {
// 		response.JSON(w, http.StatusInternalServerError, fmt.Sprintf("Error Occured while Inserting data %s", err))
// 		return
// 	}
// }

func ImportFormFile(w http.ResponseWriter, r *http.Request) {
	var reqToMongo Model.PayloadToMongo
	var res Model.BaseResponse
	response := Utilities.GetResponse()

	if r.URL.Query().Get("FileName") != "" {
		fileName := r.URL.Query().Get("FileName")
		fileName = "SBIFiles/" + fileName
		accountType := r.URL.Query().Get("AccountName")

		//based on the account type will call different workflow for the extractor
		//1.Sbiparser workflow
		//2.HDFCparser workflow
		if accountType == "SBI" {
			reqToMongo = Workflows.SBIparser(fileName)
		} else if accountType == "HDFC" {
			//Workflows.HDFCparser()
		} else {
			res := Model.BaseResponse{
				Status: false,
				Error:  "Please select Account Name as SBI or HDFC",
			}
			response.JSON(w, http.StatusInternalServerError, res)
			return
		}
		if len(reqToMongo.Data) > 0 {
			status := Workflows.SendToMongo(reqToMongo)
			if status {
				response.JSON(w, http.StatusOK, fmt.Sprintf("Successfully Inserted Excel : ", fileName, accountType))
				return
			}
		} else {
			res := Model.BaseResponse{
				Status: false,
				Error:  "Mongo payload is Empty",
			}
			response.JSON(w, http.StatusUnprocessableEntity, res)
		}

		//fmt.Println(req)
		return
	}
	response.JSON(w, http.StatusOK, res)
	return
}

package Controllers

import (
	Model "SimpleGo_xpns/Models"
	Utilities "SimpleGo_xpns/Utilities"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/bson"
)

func RegisterDataAPI(r chi.Router) {
	r.Post("/api/User/v1/SendData", SendToMongo)
	r.Get("/api/User/v1/GetExpense", GetXpns)
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
		t, _ := time.Parse("2006-01-02", x.Date)
		x.Day = t.Weekday().String()
		x.Month = t.Month().String()
		if id, err := col.InsertOne(context.Background(), x); err == nil {
			fmt.Print(id) //successfull case
		} else {
			response.JSON(w, http.StatusInternalServerError, fmt.Sprintf("Error Occured while Inserting data %s", err))
			return
		}
	}
	response.JSON(w, http.StatusOK, fmt.Sprintf("Data Imported !!"))

}

func GetXpns(w http.ResponseWriter, r *http.Request) {
	var userId string
	var res Model.PayloadToMongo
	if r.URL.Query().Get("UserId") != "" {
		userId = r.URL.Query().Get("UserId")
	}
	u, _ := strconv.ParseInt(userId, 10, 32)

	response := Utilities.GetResponse()
	c, _ := Utilities.ConectToMongo()
	col := c.Database("SimpleGo").Collection("Expenses")
	if cur, err := col.Find(context.Background(), bson.D{{"userId", u}}); err == nil {
		//fmt.Print(id)
		for cur.Next(context.TODO()) {
			//Create a value into which the single document can be decoded
			var elem Model.ExpenseBO
			err := cur.Decode(&elem)
			if err != nil {
				log.Fatal(err)
			}

			res.Data = append(res.Data, elem)

		}

		response.JSON(w, http.StatusInternalServerError, res)
		return
	} else {
		response.JSON(w, http.StatusInternalServerError, fmt.Sprintf("Error Occured while Inserting data %s", err))
		return
	}
}

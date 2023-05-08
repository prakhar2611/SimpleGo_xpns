package Utilities

import (
	"fmt"
	"time"

	"SimpleGo_xpns/Models"
	"SimpleGo_xpns/Utilities"

	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func GetDbConnection() bool {
	var username, password, dbname, host string
	username = viper.Get("Postgresql.userName").(string)
	password = viper.Get("Postgresql.password").(string)
	dbname = viper.Get("Postgresql.dbName").(string)
	host = viper.Get("Postgresql.host").(string)
	if db == nil {
		dns := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=5432 sslmode=disable", host, username, password, dbname)
		dbInstance, err := gorm.Open(postgres.New(postgres.Config{
			DSN:                  dns,
			PreferSimpleProtocol: false, // disables implicit prepared statement usage
		}), &gorm.Config{})
		if err != nil {
			return false
		} else {
			db = dbInstance
			//will us it when creating our own token with diff login system - user token
			//db.AutoMigrate(&Models.User{}, &Models.UserToken{}, &Models.GmailExcelThreadSnapShot{}, &Models.ExpenseBO{}, &Models.B64decodedResponse{})
			db.AutoMigrate(&Models.User{}, &Models.ExpenseBO{}, &Models.B64decodedResponse{})
			return true
		}
	}
	return true
}

func InsertUserData(user Models.User) bool {
	if GetDbConnection() {
		var dbUser Models.User
		resp := db.Where("Email = ?", user.Email).First(&dbUser)
		if resp.RowsAffected > 0 {
			//UpdateAuthToken(user, token) skipping saving the auth token now, saving in local cache
			return true
			//user already inserted in the db no need to insert user token
		} else {
			resp = db.Create(&user)
			if resp != nil && resp.RowsAffected > 0 {
				//UpdateAuthToken(user, token)
				return true
			} else {
				fmt.Printf("Getting error while inserting user data : %v", user.Email)
				return false
			}
		}

	}
	return false
}

func UpdateCategory(payload *Models.UpdatecategoryPayload) (bool, []string) {
	var failureMsgId []string
	if GetDbConnection() {
		for key, value := range *&payload.Data {
			r := db.Model(&Models.B64decodedResponse{}).Where("transaction_id = ?", key).Update("category", value)
			if r.RowsAffected == 0 {
				failureMsgId = append(failureMsgId, key)
			}
		}
		return true, failureMsgId
	}
	return false, failureMsgId
}

func UpdateAuthToken(user Models.User, token Models.UserToken) bool {
	var dbToken Models.UserToken
	if GetDbConnection() {
		token.UpdatedAt = time.Now()
		resp := db.Where("ID = ? ", user.ID).First(&dbToken)
		if resp.RowsAffected > 0 {
			r := db.Where("ID = ? ", user.ID).Updates(token)
			if r.RowsAffected > 0 {
				fmt.Printf("Access Token updated for user : %v", user.Email)
				return true
			} else {
				fmt.Printf("Error while updating user token : %v", user.Email)

			}
		} else {
			r := db.Create(&token)
			if r.RowsAffected > 0 {
				fmt.Printf("Access Token Created for user : %v", user.Email)
				return true
			} else {
				fmt.Printf("Error while creating user token : %v", user.Email)
			}
		}
	}
	return false
}

func GetUserToken(token string) *oauth2.Token {
	var t *oauth2.Token
	var y Models.UserToken
	if GetDbConnection() {
		resp := db.Where("access_token = ?", token).First(&y)
		if resp.RowsAffected > 0 {
			t = Utilities.MapToken(y)
		}
	}
	return t
}

func GetlastHistoryByLabel(label string) string {
	if GetDbConnection() {
		var snapShot Models.GmailExcelThreadSnapShot
		resp := db.Where("label = ?", label).First(&snapShot)
		if resp != nil && resp.RowsAffected > 0 {
			return snapShot.LastHistoryId
		} else if resp.RowsAffected == 0 {
			return ""
		}
	}
	return ""
}

func SendDataToPostgres(req []Models.ExpenseBO) bool {
	if GetDbConnection() {
		resp := db.Create(&req)
		if resp != nil && resp.RowsAffected > 0 {
			return true
		} else {
			//throwing integrity check if any comes as counting excel sheet as whole body
			fmt.Printf("Getting error while inserting err : %v", resp.Error.Error())
			return false
		}
	}
	return false
}

func SendHDFCToPostgres(req []*Models.B64decodedResponse) []string {
	var failure []string
	if GetDbConnection() {
		for _, x := range req {
			resp := db.Create(&x)
			if resp != nil && resp.RowsAffected == 0 {
				failure = append(failure, x.TransactionId)
			}
		}
	}
	return failure
}

func GetXpnsFromPostgres(from string, to string) []*Models.B64decodedResponse {
	if GetDbConnection() {
		var data []*Models.B64decodedResponse
		resp := db.Where("e_time >= ? AND e_time <= ?", from, to).Order("e_time DESC").Find(&data)
		if resp != nil && resp.RowsAffected > 0 {
			return data
		} else if resp.RowsAffected == 0 {
			return nil
		}
	}
	return nil
}

// func InsertUserData(user *Models.User) bool {
// 	if GetDbConnection() {
// 		var dbUser Models.User
// 		resp := db.Where("Email = ?", user.Email).First(dbUser)
// 		if resp.RowsAffected > 0 {
// 			//UpdateAuthToken(*user, token)
// 			fmt.Println("user is already inserted")
// 			return true
// 			//user already inserted in the db no need to insert user token
// 		} else {
// 			resp = db.Create(&user)
// 			if resp != nil && resp.RowsAffected > 0 {
// 				fmt.Println("user has been created in db")
// 				return true
// 			} else {
// 				fmt.Printf("Getting error while inserting user data : %v", user.Email)
// 				return false
// 			}
// 		}

// 	}
// 	return false
// }

//not using it currently
func Migrate(m interface{}) {

	//d.db.DropTable(&models.UserAuth{})
	//d.db.CreateTable(&models.UserAuth{})

	// Do it the hard way
	//if d.db.HasTable(&m) == false {
	// Create table for model `User`
	//  d.db.CreateTable(&m)
	//  d.logThis.Info(fmt.Sprintf("%s %s with error %s", logEntry, "Failed", d.db.Error))
	//}

	// Migrate the schema
	d := db.AutoMigrate(&m)
	if d != nil && d.Error != nil {
		//We have an error
		fmt.Printf(d.Error())
	}
}

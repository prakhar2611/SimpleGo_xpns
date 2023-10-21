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
			db.AutoMigrate(&Models.User{}, &Models.ExpenseBO{}, &Models.B64decodedResponse{}, &Models.VpaMapping{}, &Models.PocketsMappingDbo{}, &Models.DocsMeta{})
			return true
		}
	}
	return true
}

func InsertUserData(user Models.User) int {
	if GetDbConnection() {
		var dbUser Models.User
		resp := db.Where("Email = ?", user.Email).First(&dbUser)
		if resp.RowsAffected > 0 {
			//UpdateAuthToken(user, token) skipping saving the auth token now, saving in local cache
			return 2
			//user already inserted in the db no need to insert user token
		} else {
			resp = db.Create(&user)
			if resp != nil && resp.RowsAffected > 0 {
				//UpdateAuthToken(user, token)
				return 1
			} else {
				fmt.Printf("Getting error while inserting user data : %v", user.Email)
				return 0
			}
		}

	}
	return 0
}

func UpdateCategory(payload *Models.UpdatecategoryPayload) (bool, []string) {
	var failureMsgId []string
	if GetDbConnection() {
		for key, value := range *&payload.Data {
			r := db.Model(&Models.B64decodedResponse{}).Where("transaction_id = ?", key).Updates(map[string]interface{}{"category": value.Category, "label": value.Label})
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
			//throwing integrity check if any comes as counting excel sheet as whole body
			fmt.Printf("Getting error while inserting err : %v", resp.Error.Error())
			return false
		}
	}
	return false
}

func SendHDFCToPostgres(req []Models.B64decodedResponse) []string {
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

//db helper to get the xpns from the db
func GetXpnsFromPostgres(from string, to string, userId string) []*Models.B64decodedResponse {
	if GetDbConnection() {
		var data []*Models.B64decodedResponse
		resp := db.Where("e_time >= ? AND e_time <= ? AND to_account != '' AND user_id = ? ", from, to, userId).Order("e_time DESC").Find(&data)
		if resp != nil && resp.RowsAffected > 0 {
			return data
		} else if resp.RowsAffected == 0 {
			return nil
		}
	}
	return nil
}

//db helper to get the group VPA based on top number of resp vpa txn and total amount value
func GetGroupedVpa(limit string, offset string) []Models.VpaMapping {
	if GetDbConnection() {
		id := 1
		var data []Models.VpaMapping
		rows, err := db.Raw("select to_account as vpa ,SUM(amount_debited) as totalAmount ,count(*) as totalTxn, label, category from b64decoded_responses group by to_account,label, category order by totaltxn DESC limit ? offset ?", limit, offset).Rows()
		defer rows.Close()
		if err == nil {
			for rows.Next() {
				var raw Models.VpaMapping
				rows.Scan(&raw.Vpa, &raw.TotalAmount, &raw.TotalTxn, &raw.Label, &raw.Category)
				raw.Id = id
				id += 1
				data = append(data, raw)
			}
		}
		if len(data) > 0 {
			return data
		} else if len(data) == 0 {
			return nil
		}
	}
	return nil

}

func PushVPAMappingToDb(req []Models.VpaMapping) int {
	var count int
	if GetDbConnection() {
		for _, x := range req {
			resp := db.Create(&x)
			if resp != nil && resp.RowsAffected == 0 {
				count += 1
			}
		}
	}
	return count
}

// db helper to create new pocket else update the initial value based on distinct user
func CreateAndUpdatePocketDb(payload *Models.UpdatePocketsPayload, userId string) bool {

	if GetDbConnection() {
		for _, req := range payload.Data {
			req.UserId = userId
			resp := db.Where("pocket = ? && userId = ? ", req.Pocket, userId).First(&Models.PocketsMappingDbo{})
			if resp.RowsAffected > 0 {
				resp := db.Model(&Models.PocketsMappingDbo{}).Where("pocket= ? and user_id = ?", req.Pocket, userId).Updates(map[string]interface{}{"labels": req.Labels})
				if resp != nil && resp.RowsAffected > 0 {
					return true
				} else {
					fmt.Printf("Getting error while updating pocket lables, pocket : %v, user id : %v", req.Pocket, req.UserId)
					return false
				}
			} else {
				resp := db.Create(req)
				if resp != nil && resp.RowsAffected > 0 {
					return true
				} else {
					fmt.Printf("Getting error while Creating new pocket, pocket : %v, user id : %v", req.Pocket, req.UserId)
					return false
				}
			}
		}

	}
	return false
}

//db helper to update txns across the table with pocket
func UpdatePocketTxnLevel(payload *Models.UpdatePocketsPayload, userId string) (bool, []string) {
	var failureMsgId []string
	if GetDbConnection() {
		for _, req := range payload.Data {
			r := db.Model(&Models.B64decodedResponse{}).Where("label in ? and user_id = ?", req.Labels, userId).Updates(map[string]interface{}{"pocket": req.Pocket})
			if r.RowsAffected == 0 {
				failureMsgId = append(failureMsgId, req.Pocket)
			}
		}
		return true, failureMsgId
	}
	return false, failureMsgId
}

//db helper to update txns across the table with same VPA name
func UpdateVPATxnLevel(payload *Models.UpdatecategoryPayload, userId string) (bool, []string) {
	var failureMsgId []string
	if GetDbConnection() {
		for key, value := range *&payload.Data {
			r := db.Model(&Models.B64decodedResponse{}).Where("to_account= ? and user_id = ?", key, userId).Updates(map[string]interface{}{"category": value.Category, "label": value.Label})
			if r.RowsAffected == 0 {
				failureMsgId = append(failureMsgId, key)
			}
		}
		return true, failureMsgId
	}
	return false, failureMsgId
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

func CreateDoc(payload *Models.DocsMeta, userId string) bool {
	if GetDbConnection() {
		var dbDoc Models.DocsMeta
		resp := db.Where("title = ? and user_id = ? and folder = ?", payload.Title, userId, payload.Folder).First(&dbDoc)
		if resp.RowsAffected > 0 {
			resp := db.Model(Models.DocsMeta{}).Where("title = ? and user_id = ? and folder = ?", payload.Title, userId, payload.Folder).Update("meta_data", payload.MetaData)
			if resp != nil && resp.RowsAffected > 0 {
				return true
			} else {
				fmt.Printf("Getting error while Updating user data : %v", payload.ID)
				return false
			}
		} else {
			resp := db.Create(&payload)
			if resp != nil && resp.RowsAffected > 0 {
				return true
			} else {
				fmt.Printf("Getting error while inserting user data : %v", payload.ID)
				return false
			}

		}
	}
	return false
}

func GetAllDoc(user_id string) []Models.DocsMeta {
	if GetDbConnection() {
		var data []Models.DocsMeta

		// rows, err := db.Raw("select id,meta_Data,user_id,title,folder,created_at from docs_meta").Rows()

		rows, err := db.Raw("select id,user_id,title,folder,created_at from docs_meta where user_id = ?", user_id).Rows()
		defer rows.Close()
		if err == nil {
			for rows.Next() {
				var raw Models.DocsMeta
				rows.Scan(&raw.ID, &raw.UserId, &raw.Title, &raw.Folder, &raw.CreatedAt)
				data = append(data, raw)
			}
		}
		if len(data) > 0 {
			return data
		} else if len(data) == 0 {
			return nil
		}
	}
	return nil
}

func GetDocMeta(request Models.GetDocMetaRequest, user_id string) *Models.DocsMeta {
	if GetDbConnection() {
		var raw Models.DocsMeta
		rows, err := db.Raw("select id,meta_Data,user_id,title,folder,created_at from docs_meta where title = ? and user_id = ? and folder = ?", request.Title, user_id, request.Folder).Rows()
		defer rows.Close()
		if err == nil {
			for rows.Next() {
				rows.Scan(&raw.ID, &raw.MetaData, &raw.UserId, &raw.Title, &raw.Folder, &raw.CreatedAt)
			}
		}
		if len(raw.Title) > 0 {
			return &raw
		} else {
			return nil
		}
	}
	return nil
}

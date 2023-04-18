// // Not using it currently
package Utilities

// import (
// 	"database/sql"
// 	"fmt"

// 	_ "github.com/go-sql-driver/mysql"
// 	_ "github.com/lib/pq"
// )

// var db *sql.DB
// var passwrod string = "expensedb1"

// func ConnectDB() {

// 	// var dbName string = "expensedb"
// 	// var dbUser string = "Admin"
// 	// var dbHost string = "expensedb.cxmpaamkcbyh.ap-southeast-1.rds.amazonaws.com"
// 	// var dbPort int = 1433
// 	// var dbEndpoint string = fmt.Sprintf("%s:%d", dbHost, dbPort)
// 	// var region string = "us-east-1"

// 	// cfg, err := config.LoadDefaultConfig(context.TODO())
// 	// if err != nil {
// 	// 	panic("configuration error: " + err.Error())
// 	// }

// 	// authenticationToken, err := auth.BuildAuthToken(
// 	// 	context.TODO(), dbEndpoint, region, dbUser, cfg.Credentials)
// 	// if err != nil {
// 	// 	panic("failed to create authentication token: " + err.Error())
// 	// }

// 	// dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?tls=true&allowCleartextPasswords=true&allowCleartextPasswords=true",
// 	// 	dbUser, authenticationToken, dbEndpoint, dbName,
// 	// )

// 	// db, err := sql.Open("mysql", dsn)
// 	// if err != nil {
// 	// 	panic(err)
// 	// }
// 	fmt.Println("all good 1")
// 	var err error
// 	//authToken, err := rdsutils.BuildAuthToken("expensedb.cxmpaamkcbyh.ap-southeast-1.rds.amazonaws.com", "ap-southeast-1", "Admin", passwrod)
// 	db, err = sql.Open("mysql", "Admin:expensedb1@tcp(expensedb.cxmpaamkcbyh.ap-southeast-1.rds.amazonaws.com:1433)/expensedb?tls=true&allowCleartextPasswords=true")

// 	if err != nil {
// 		fmt.Print(err.Error())
// 	}
// 	fmt.Println("all good 2")

// 	q := "SELECT [FirstName] FROM [dbo].[User]"

// 	h, err := db.Query(q)
// 	var name string
// 	err2 := h.Scan(&name)
// 	if err2 == nil {
// 		fmt.Println(name)
// 	}
// }

// func GetFromDB() {

// }

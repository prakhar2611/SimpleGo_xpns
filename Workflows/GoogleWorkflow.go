package Workflows

import (
	Executer "SimpleGo_xpns/Utilities/APIExecuter"
	"encoding/json"
	"fmt"

	"SimpleGo_xpns/Models"

	"google.golang.org/api/gmail/v1"
)

func GetUserInfo(accessToken string) *Models.User {
	var req Executer.APIRequest
	req.BaseURL = "https://www.googleapis.com/"
	req.Action = fmt.Sprintf("oauth2/v1/userinfo?alt=json&access_token=%v", accessToken)
	resp := Executer.GET(req)
	if resp.Status == 200 {
		var user *Models.User
		err := json.Unmarshal([]byte(resp.Response), &user)
		if err == nil {
			return user
		}
	}
	return nil
}

func VerifyIdToken(idToken string) string {

	userid := GetUserInfo(idToken)
	if userid == nil {
		return ""
	}
	return userid.ID
}

func FetchAllThread(label string, srv *gmail.Service) bool {
	user := "me"
	var k *gmail.ListThreadsResponse
	threadservice := srv.Users.Threads.List(user)
	threadservice.Q(fmt.Sprintf("label:%v", label))
	k, _ = threadservice.Do()
	fmt.Println(len(k.Threads))
	//looping through threads
	for _, x := range k.Threads {
		//getting one thread  using id
		threadservice := srv.Users.Threads.Get(user, x.Id)
		th, _ := threadservice.Do()

		fmt.Println(th.HistoryId)
		fmt.Println(len(th.Messages))

		// //looping over mails
		// for _, msg := range th.Messages {

		// 	msgId := msg.Id
		// 	attId := msg.Payload.Parts[1].Body.AttachmentId

		// 	msgService := srv.Users.Messages.Attachments.Get(user, msgId, attId)
		// 	u, _ := msgService.Do()
		// 	decoded, _ := base64.URLEncoding.DecodeString(u.Data)
		// 	d := bytes.NewReader(decoded)
		// 	req := Workflows.SBIparser("", d)
		// 	fmt.Printf("%v", req)
		//send data to postgres
		// isSuccess := dbConnector.SendDataToPostgres(req)
		// if isSuccess {
		// 	fmt.Println("sucessfully inserted user record to db")
		// }
		//}
	}
	return true
}

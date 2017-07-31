package remote

import (
	"net/url"
	"fmt"
	"github.com/pharosnet/logs"
	"liulishuo/somechat/log"
)

type ChatApiNotifyContactAddRequestAcceptForm struct {
	Id string	`form:"id"`
	Owner string	`form:"owner"`
	UserId string	`form:"userId"`
	UserName string	`form:"userName"`
	UserEmail string	`form:"userEmail"`
	ReqId string	`form:"reqId"`
}

const chat_api_notify_contact_add_req_accept = "/api/notify/contact/add/accept"

func ChatApiNotifyContactAddRequestAccept(form ChatApiNotifyContactAddRequestAcceptForm) (*responseResult, error) {
	remoteForm := url.Values{}
	remoteForm.Set("id", form.Id)
	remoteForm.Set("owner", form.Owner)
	remoteForm.Set("userId", form.UserId)
	remoteForm.Set("userEmail", form.UserEmail)
	remoteForm.Set("userName", form.UserName)
	remoteForm.Set("reqId", form.ReqId)
	respBody, respErr := ChatApiHtppPost(chat_api_notify_contact_add_req_accept, remoteForm)
	if respErr != nil {
		err := fmt.Errorf("new http request failed, chat remote api, error = %v", respErr)
		log.Log().Println(logs.Error(err).Extra(logs.F{"remote", "chat"}).Trace())
		return nil, err
	}
	result, resultErr := newResponseResult(respBody)
	if resultErr != nil {
		err := fmt.Errorf("read http response failed, chat remote api, error = %v", respErr)
		log.Log().Println(logs.Error(err).Extra(logs.F{"remote", "chat"}).Trace())
		return nil, err
	}
	return result, nil
}

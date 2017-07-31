package remote

import (
	"net/url"
	"fmt"
	"github.com/pharosnet/logs"
	"liulishuo/somechat/log"
)

type ChatApiNotifyContactAddRequestNewForm struct {
	ReqId string	`form:"reqId"`
	FromId string	`form:"fromId"`
	FromName string	`form:"fromName"`
	FromEmail string	`form:"fromEmail"`
	ToId string	`form:"toId"`
	ToName string	`form:"toName"`
	ToEmail string	`form:"toEmail"`
}

const chat_api_notify_contact_add_req_new = "/api/notify/contact/add/new"

func ChatApiNotifyContactAddRequestNew(form ChatApiNotifyContactAddRequestNewForm) (*responseResult, error) {
	remoteForm := url.Values{}
	remoteForm.Set("reqId", form.ReqId)
	remoteForm.Set("fromId", form.FromId)
	remoteForm.Set("fromName", form.FromName)
	remoteForm.Set("fromEmail", form.FromEmail)
	remoteForm.Set("toId", form.ToId)
	remoteForm.Set("toName", form.ToName)
	remoteForm.Set("toEmail", form.ToEmail)
	respBody, respErr := ChatApiHtppPost(chat_api_notify_contact_add_req_new, remoteForm)
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

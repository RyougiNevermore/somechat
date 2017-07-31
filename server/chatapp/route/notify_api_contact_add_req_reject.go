package route

import (
	"github.com/kataras/iris/context"
	"fmt"
	"github.com/pharosnet/logs"
	"liulishuo/somechat/server/chatapp/agent"
	"liulishuo/somechat/log"
)

type notifyContactAddRequestRejectForm struct {
	ReqId string	`form:"reqId"`
	FromId string	`form:"fromId"`
	FromName string	`form:"fromName"`
	FromEmail string	`form:"fromEmail"`
	ToId string	`form:"toId"`
	ToName string	`form:"toName"`
	ToEmail string	`form:"toEmail"`
}

func notifyApiContactAddRequestReject(ctx context.Context)  {
	form := new(notifyContactAddRequestRejectForm)
	if formErr := ctx.ReadForm(form); formErr != nil {
		err := fmt.Errorf("form read failed, %v", formErr)
		log.Log().Println(logs.Error(err).Extra(logs.F{"http", ctx.RequestPath(true)}).Trace())
		ctx.JSON(result{Success:false, Message:"bad form"})
		return
	}
	msg := new(agent.Message)
	msg.Head = make(map[string]string)
	msg.Head["fromId"] = form.FromId
	msg.Head["fromName"] = form.FromName
	msg.Head["fromEmail"] = form.FromEmail
	msg.Head["toId"] = form.ToId
	msg.Head["toName"] = form.ToName
	msg.Head["toEmail"] = form.ToEmail
	msg.SetContentType("contact_add_req")
	msg.Body = fmt.Sprintf(`{"reqId": "%s", "action":"reject", "message":"%s reject your contact add request"}`, form.ReqId, form.ToEmail)
	err := gateway.Notify(form.FromId, ws_topic_notify, msg)
	if err != nil {
		err = fmt.Errorf("notify contact add request reject failed, to = %s, msg = %v", form.FromId, msg)
		log.Log().Println(logs.Error(err).Extra(logs.F{"http", ctx.RequestPath(true)}).Trace())
		ctx.JSON(result{Success:false, Message:"notify contact add request reject failed, to = %s" + form.FromId})
		return
	}
	ctx.JSON(result{Success:true, Message:""})
	return
}

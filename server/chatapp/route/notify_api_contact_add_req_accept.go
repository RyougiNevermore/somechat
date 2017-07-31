package route

import (
	"github.com/kataras/iris/context"
	"fmt"
	"github.com/pharosnet/logs"
	"liulishuo/somechat/log"
	"liulishuo/somechat/server/chatapp/agent"
)

type notifyContactAddRequestAcceptForm struct {
	Id string	`form:"id"`
	Owner string	`form:"owner"`
	UserId string	`form:"userId"`
	UserName string	`form:"userName"`
	UserEmail string	`form:"userEmail"`
	ReqId string	`form:"reqId"`
}

func notifyApiContactAddRequestAccept(ctx context.Context)  {
	form := new(notifyContactAddRequestAcceptForm)
	if formErr := ctx.ReadForm(form); formErr != nil {
		err := fmt.Errorf("form read failed, %v", formErr)
		log.Log().Println(logs.Error(err).Extra(logs.F{"http", ctx.RequestPath(true)}).Trace())
		ctx.JSON(result{Success:false, Message:"bad form"})
		return
	}
	msg := new(agent.Message)
	msg.Head = make(map[string]string)
	msg.Head["id"] = form.Id
	msg.Head["owner"] = form.Owner
	msg.Head["userId"] = form.UserId
	msg.Head["userName"] = form.UserName
	msg.Head["userEmail"] = form.UserEmail
	msg.Head["reqId"] = form.ReqId
	msg.SetContentType("contact_add_req")
	msg.Body = fmt.Sprintf(`{"reqId": "%s", "action":"accept", "message":"%s accpet your contact add request"}`, form.Id, form.UserEmail)
	err := gateway.Notify(form.Owner, ws_topic_notify, msg)
	if err != nil {
		err = fmt.Errorf("notify contact add request accept failed, to = %s, msg = %v", form.Id, msg)
		log.Log().Println(logs.Error(err).Extra(logs.F{"http", ctx.RequestPath(true)}).Trace())
		ctx.JSON(result{Success:false, Message:"notify contact add request accept failed, to = %s" + form.Id})
		return
	}
	ctx.JSON(result{Success:true, Message:""})
	return
}

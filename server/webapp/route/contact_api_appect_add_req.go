package route

import (
	"github.com/kataras/iris/context"
	"fmt"
	"github.com/pharosnet/logs"
	"liulishuo/somechat/log"
	"strings"
	"liulishuo/somechat/server/webapp/service"
	"liulishuo/somechat/server/webapp/remote"
)

type contactAcceptForm struct {
	Id string `form:"reqId"`
	OpUserId string `form:"userId"`
}

func contactApiAcceptAddRequest(ctx context.Context) {
	form := new(contactAcceptForm)
	if formErr := ctx.ReadForm(form); formErr != nil {
		err := fmt.Errorf("form read failed, %v", formErr)
		log.Log().Println(logs.Error(err).Extra(logs.F{"http", ctx.RequestPath(true)}).Trace())
		ctx.JSON(result{Success:false, Message:"bad form"})
		return
	}
	reqId := strings.TrimSpace(form.Id)
	opUserId := strings.TrimSpace(form.OpUserId)
	if reqId == "" || opUserId == "" {
		err := fmt.Errorf("form is bad failed, req id = %v, user id = %v", form.Id, form.OpUserId)
		log.Log().Println(logs.Error(err).Extra(logs.F{"http", ctx.RequestPath(true)}).Trace())
		ctx.JSON(result{Success:false, Message:"bad form"})
		return
	}
	contact, acceptErr := service.ContactAddRequestAccept(reqId, opUserId)
	if acceptErr != nil {
		err := fmt.Errorf("accept failed, req id = %v, user id = %v, error = %v", form.Id, form.OpUserId, acceptErr)
		log.Log().Println(logs.Error(err).Extra(logs.F{"http", ctx.RequestPath(true)}).Trace())
		ctx.JSON(result{Success:false, Message:"accept failed."})
		return
	}
	ctx.JSON(
		result{
			Success:true, Message:"accept successed.",
			Data:map[string]interface{}{"userId": contact.UserId, "userName": contact.UserName, "userEmail": contact.UserEmail},
		})
	notifyContactAddRequestAccept(reqId, contact)
	return
}

func notifyContactAddRequestAccept(reqId string, contact *service.Contact)  {
	go func(reqId string, contact *service.Contact) {
		if reqId == "" {
			return
		}
		row, rowGetErr := service.ContactGetOneByOwnerAndUserId(contact.UserId, contact.Owner)
		if rowGetErr != nil {
			err := fmt.Errorf("chat remote failed, get contact add req failed, error = %v", rowGetErr)
			log.Log().Println(logs.Error(err).Extra(logs.F{"remote", "chat"}).Trace())
			return
		}
		_, remoteErr := remote.ChatApiNotifyContactAddRequestAccept(
			remote.ChatApiNotifyContactAddRequestAcceptForm{
				Id: row.Id,
				Owner: row.Owner,
				UserId: row.UserId,
				UserName: row.UserName,
				UserEmail: row.UserEmail,
				ReqId:reqId,
			},
		)
		if remoteErr != nil {
			err := fmt.Errorf("chat remote failed, error = %v", remoteErr)
			log.Log().Println(logs.Error(err).Extra(logs.F{"remote", "chat"}).Trace())
			return
		}
	}(reqId, contact)
}
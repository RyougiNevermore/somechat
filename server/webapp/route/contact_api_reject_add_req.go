package route

import (
	"github.com/kataras/iris/context"
	"strings"
	"fmt"
	"github.com/pharosnet/logs"
	"liulishuo/somechat/log"
	"liulishuo/somechat/server/webapp/service"
	"liulishuo/somechat/server/webapp/remote"
)

type contactRejectForm struct {
	Id string `form:"reqId"`
	OpUserId string `form:"userId"`
}

func contactApiRejectAddRequest(ctx context.Context)  {
	form := new(contactRejectForm)
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
	//
	req, reqGetErr := service.ContactAddRequestGetById(reqId)
	if reqGetErr != nil {
		err := fmt.Errorf("reject failed, get contact add req failed, id = %s, error = %v", reqId, reqGetErr)
		log.Log().Println(logs.Error(err).Extra(logs.F{"http", ctx.RequestPath(true)}).Trace())
		ctx.JSON(result{Success:false, Message:"can not find request, id = " + reqId})
		return
	}
	rejectErr := service.ContactAddRequestReject(reqId, opUserId)
	if rejectErr != nil {
		err := fmt.Errorf("reject failed, req id = %v, user id = %v, error = %v", form.Id, form.OpUserId, rejectErr)
		log.Log().Println(logs.Error(err).Extra(logs.F{"http", ctx.RequestPath(true)}).Trace())
		ctx.JSON(result{Success:false, Message:"reject failed."})
		return
	}
	ctx.JSON(
		result{
			Success:true, Message:"reject successed.",
		})
	notifyContactAddRequestReject(req)
	return
}

func notifyContactAddRequestReject(row *service.ContactAddRequest)  {
	go func(row *service.ContactAddRequest) {
		_, remoteErr := remote.ChatApiNotifyContactAddRequestReject(
			remote.ChatApiNotifyContactAddRequestRejectForm{
				ReqId: row.Id,
				FromId: row.FromId,
				FromName: row.FromName,
				FromEmail: row.FromEmail,
				ToId: row.ToId,
				ToName: row.ToName,
				ToEmail: row.ToEmail,
			},
		)
		if remoteErr != nil {
			err := fmt.Errorf("chat remote failed, error = %v", remoteErr)
			log.Log().Println(logs.Error(err).Extra(logs.F{"remote", "chat"}).Trace())
			return
		}
	}(row)
}
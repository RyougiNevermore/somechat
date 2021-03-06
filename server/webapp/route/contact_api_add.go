package route

import (
	"github.com/kataras/iris/context"
	"fmt"
	"github.com/pharosnet/logs"
	"liulishuo/somechat/log"
	"liulishuo/somechat/server/webapp/service"
	"strings"
	"liulishuo/somechat/server/webapp/remote"
)

type contactAddForm struct {
	UserId string	`form:"userId"`
	Email string	`form:"email"`
}

func contactApiAdd(ctx context.Context)  {
	form := new(contactAddForm)
	if formErr := ctx.ReadForm(form); formErr != nil {
		err := fmt.Errorf("form read failed, %v", formErr)
		log.Log().Println(logs.Error(err).Extra(logs.F{"http", ctx.RequestPath(true)}).Trace())
		ctx.JSON(result{Success:false, Message:"bad form"})
		return
	}
	form.Email = strings.TrimSpace(form.Email)
	form.UserId = strings.TrimSpace(form.UserId)
	if form.Email == "" || form.UserId == "" {
		err := fmt.Errorf("form read failed, email = %s, user id = %s", form.Email, form.UserId)
		log.Log().Println(logs.Error(err).Extra(logs.F{"http", ctx.RequestPath(true)}).Trace())
		ctx.JSON(result{Success:false, Message:"bad form"})
		return
	}
	toUser, toUserGetErr := service.UserGetByEmail(form.Email)
	if toUserGetErr != nil {
		err := fmt.Errorf("get user failed, email = %s, error = %v", form.Email, toUserGetErr)
		log.Log().Println(logs.Error(err).Extra(logs.F{"http", ctx.RequestPath(true)}).Trace())
		ctx.JSON(result{Success:false, Message:"can not find user by email"})
		return
	}
	req, contactAddReqErr := service.ContactAddRequestNew(form.UserId, toUser.Id)
	if contactAddReqErr != nil {
		err := fmt.Errorf("make contact add req failed, error = %v", contactAddReqErr)
		log.Log().Println(logs.Error(err).Extra(logs.F{"http", ctx.RequestPath(true)}).Trace())
		ctx.JSON(result{Success:false, Message:"make contact add request failed."})
		return
	}
	ctx.JSON(result{Success:true,
		Message:"make contact add request successed. wait for accepting.",
		Data:map[string]interface{}{
			"id" : req.Id,
			"toUserName": req.ToName,
			"toUserEmail": req.ToEmail,
			"toUserId": req.ToId,
		}},
	)
	notifyContactAddRequestNew(req)
	return
}

func notifyContactAddRequestNew(row *service.ContactAddRequest)  {
	go func(row *service.ContactAddRequest) {
		_, remoteErr := remote.ChatApiNotifyContactAddRequestNew(
			remote.ChatApiNotifyContactAddRequestNewForm{
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
package route

import (
	"github.com/kataras/iris/context"
	"fmt"
	"github.com/pharosnet/logs"
	"liulishuo/somechat/log"
	"liulishuo/somechat/server/webapp/service"
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
	_, contactAddReqErr := service.ContactAddRequestNew(form.UserId, toUser.Id)
	if contactAddReqErr != nil {
		err := fmt.Errorf("make contact add req failed, error = %v", contactAddReqErr)
		log.Log().Println(logs.Error(err).Extra(logs.F{"http", ctx.RequestPath(true)}).Trace())
		ctx.JSON(result{Success:false, Message:"make contact add request failed."})
		return
	}
	ctx.JSON(result{Success:true, Message:"make contact add request successed. wait for accepting."})
	return
}

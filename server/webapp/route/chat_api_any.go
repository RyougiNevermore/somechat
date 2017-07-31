package route

import (
	"github.com/kataras/iris/context"
	"fmt"
	"github.com/pharosnet/logs"
	"strings"
	"liulishuo/somechat/log"
	"liulishuo/somechat/server/webapp/service"
)

type chatAnyForm struct {
	UserId string	`form:"userId"`
	Email string	`form:"email"`
}


func chatApiAny(ctx context.Context)  {
	form := new(chatAnyForm)
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
	// check in contact
	contact, contactGetErr := service.ContactGetOneByOwnerAndUserEmail(form.UserId, form.Email)
	if contactGetErr != nil {
		err := fmt.Errorf("get contact failed, email = %s, user id = %s", form.Email, form.UserId)
		log.Log().Println(logs.Error(err).Extra(logs.F{"http", ctx.RequestPath(true)}).Trace())
		ctx.JSON(result{Success:false, Message:"check user contact failed, email=" + form.Email})
		return
	}
	if contact != nil && contact.Id != "" {
		ctx.JSON(result{
			Success:true,
			Message:fmt.Sprintf("This email (%s) is in your contact.", form.Email),
			Data:map[string]interface{}{"userId":contact.UserId},
		})
		return
	}
	// no contact to new contact add request
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
	ctx.JSON(result{
		Success:true,
		Message:fmt.Sprintf("This email (%s) is not in your contact, and send a request to add a contact.", form.Email),
		Data:map[string]interface{}{"userId":req.ToId},
	})
	return
}

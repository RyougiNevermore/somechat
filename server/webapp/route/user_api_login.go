package route

import (
	"github.com/kataras/iris/context"
	"liulishuo/somechat/log"
	"github.com/pharosnet/logs"
	"fmt"
	"liulishuo/somechat/server/webapp/service"
)

type loginForm struct {
	Email string	`form:"email"`
	Password string	`form:"password"`
}

func userApiLogin(ctx context.Context)  {
	form := new(loginForm)
	if formErr := ctx.ReadForm(form); formErr != nil {
		err := fmt.Errorf("form read failed, %v", formErr)
		log.Log().Println(logs.Error(err).Extra(logs.F{"http", ctx.RequestPath(true)}).Trace())
		ctx.JSON(result{Success:false, Message:"bad form"})
		return
	}
	if form.Email == "" || form.Password == "" {
		err := fmt.Errorf("form read failed, email = %s, password = %s", form.Email, form.Password)
		log.Log().Println(logs.Error(err).Extra(logs.F{"http", ctx.RequestPath(true)}).Trace())
		ctx.JSON(result{Success:false, Message:"bad form"})
		return
	}
	user, userGetErr := service.UserGetByEmail(form.Email)
	if userGetErr != nil {
		err := fmt.Errorf("user get failed, %v", userGetErr)
		log.Log().Println(logs.Error(err).Extra(logs.F{"http", ctx.RequestPath(true)}).Trace())
		ctx.JSON(result{Success:false, Message:"get user failed."})
		return
	}
	if user.Password != form.Password {
		ctx.JSON(result{Success:false, Message:"failed. bad password."})
		return
	}
	ctx.JSON(result{Success:true, Message:"Success.", Data:map[string]interface{}{"userId":user.Id}})
	return
}

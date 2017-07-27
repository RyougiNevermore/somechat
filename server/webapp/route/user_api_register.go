package route

import (
	"github.com/kataras/iris/context"
	"fmt"
	"github.com/pharosnet/logs"
	"liulishuo/somechat/log"
	"liulishuo/somechat/server/webapp/service"
)

type registerForm struct {
	Name string `form:"name"`
	Email string	`form:"email"`
	Password string	`form:"password"`
}

func userApiRegister(ctx context.Context)  {
	form := new(registerForm)
	if formErr := ctx.ReadForm(form); formErr != nil {
		err := fmt.Errorf("form read failed, %v", formErr)
		log.Log().Println(logs.Error(err).Extra(logs.F{"http", ctx.RequestPath(true)}).Trace())
		ctx.JSON(result{Success:false, Message:"bad form"})
		return
	}
	if form.Email == "" || form.Password == "" || form.Name == "" {
		err := fmt.Errorf("form read failed, email = %s, password = %s, name = %s", form.Email, form.Password, form.Name)
		log.Log().Println(logs.Error(err).Extra(logs.F{"http", ctx.RequestPath(true)}).Trace())
		ctx.JSON(result{Success:false, Message:"bad form"})
		return
	}
	user, userAddErr := service.UserAdd(form.Name, form.Email, form.Password)
	if userAddErr != nil {
		err := fmt.Errorf("user add failed, %v", userAddErr)
		log.Log().Println(logs.Error(err).Extra(logs.F{"http", ctx.RequestPath(true)}).Trace())
		ctx.JSON(result{Success:false, Message:"add user failed."})
		return
	}
	ctx.JSON(result{Success:true, Message:"Success.", Data:map[string]interface{}{"userId":user.Id}})
	return
}

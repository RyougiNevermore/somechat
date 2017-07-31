package app

import (
	"github.com/kataras/iris"
	"liulishuo/somechat/server/chatapp/conf"
	"liulishuo/somechat/server/chatapp/route"
)

func StartUp() error {
	app := iris.New()
	route.Register(app)
	return app.Run(iris.Addr(conf.Conf.Web.Port), iris.WithCharset("UTF-8"))
}

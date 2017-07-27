package app

import (
	"github.com/kataras/iris"
	"liulishuo/somechat/server/webapp/conf"
	"github.com/kataras/iris/view"
	"liulishuo/somechat/server/webapp/route"
)

func StartUp() error {
	app := iris.New()
	app.StaticWeb("/static", conf.Conf.Web.Static)
	app.AttachView(view.HTML(conf.Conf.Web.Tpl, ".html").Reload(true))
	route.Register(app)
	return app.Run(iris.Addr(conf.Conf.Web.Port), iris.WithCharset("UTF-8"))
}

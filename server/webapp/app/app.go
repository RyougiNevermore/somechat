package app

import (
	"github.com/kataras/iris"
	"liulishuo/somechat/server/webapp/conf"
	"github.com/kataras/iris/view"
	"liulishuo/somechat/server/webapp/route"
	"github.com/kataras/iris/websocket"
	"github.com/kataras/iris/context"
	"html/template"
)

func StartUp() error {
	app := iris.New()
	app.StaticWeb("/static", conf.Conf.Web.Static)

	viewer := view.HTML(conf.Conf.Web.Tpl, ".html").Reload(true)
	viewer.AddFunc(
		"unescaped",
		func (x string) interface{} {
			return template.HTML(x)
		},
	)
	app.AttachView(viewer)
	app.Any("/iris-ws.js", func(ctx context.Context) {
		ctx.Write(websocket.ClientSource)
	})
	route.Register(app)
	return app.Run(iris.Addr(conf.Conf.Web.Port), iris.WithCharset("UTF-8"))
}

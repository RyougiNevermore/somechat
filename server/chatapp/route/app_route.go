package route

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/websocket"
	"github.com/kataras/iris/context"
	"liulishuo/somechat/server/chatapp/agent"
	"sync"
)



func Register(app *iris.Application)  {
	// agent
	wsAgentRef = make(map[string]string)
	lock = sync.Mutex{}
	bus = agent.NewAgentLocalBus()
	store := agent.MessageStore{}
	gateway = agent.NewGateway(bus, &store)
	// ws
	app.Any("/iris-ws.js", func(ctx context.Context) {
		ctx.Write(websocket.ClientSource)
	})
	ws := websocket.New(websocket.Config{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	})
	ws.OnConnection(webSocketHandler)
	app.Get("/msg", ws.Handler())
	// api
	app.Post("/api/notify/contact/add/new", notifyApiContactAddRequestNew)
	app.Post("/api/notify/contact/add/accept", notifyApiContactAddRequestAccept)
	app.Post("/api/notify/contact/add/reject", notifyApiContactAddRequestReject)
}

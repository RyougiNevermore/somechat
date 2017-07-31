package route

import "github.com/kataras/iris"

func Register(app *iris.Application)  {
	app.Get("/", userViewLogin)
	app.Get("/register", userViewRegister)
	app.Get("/contacts", contactsView)
	app.Get("/chat", chatView)


	// user api
	app.Post("/api/user/login", userApiLogin)
	app.Post("/api/user/register", userApiRegister)

	// contact
	app.Post("/api/contact/add", contactApiAdd)
	app.Post("/api/contact/add/request/accept", contactApiAcceptAddRequest)
	app.Post("/api/contact/add/request/reject", contactApiRejectAddRequest)

	// chat
	app.Post("/api/chat/message/id/new", chatApiMessageIdNew)
	app.Post("/api/chat/message/remove", chatApiMessageRemove)
	app.Post("/api/chat/any", chatApiAny)
	app.Post("/api/chat/history", chatApiHistory)

}

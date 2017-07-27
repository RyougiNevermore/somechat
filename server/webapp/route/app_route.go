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


}

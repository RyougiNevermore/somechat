package route

import "github.com/kataras/iris/context"

func userViewLogin(ctx context.Context)  {
	ctx.View(view_login)
}

package route

import "github.com/kataras/iris/context"

func userViewRegister(ctx context.Context)  {
	ctx.View(view_register)
}

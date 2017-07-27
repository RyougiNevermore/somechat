package route

import "github.com/kataras/iris/context"

func chatView(ctx context.Context)  {
	ctx.View(view_chat)
}

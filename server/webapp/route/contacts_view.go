package route

import "github.com/kataras/iris/context"

func contactsView(ctx context.Context)  {
	// TODO LOAD CONTACTS AND UNREAD
	ctx.View(view_contacts)
}

package route

import (
	"github.com/kataras/iris/context"
	"liulishuo/somechat/server/webapp/service"
	"fmt"
	"github.com/pharosnet/logs"
	"liulishuo/somechat/log"
)

func chatApiMessageRemove(ctx context.Context)  {
	messageId := ctx.FormValue("msgId")
	userId := ctx.FormValue("userId")
	err := service.ChatHistoryRemove(messageId, userId)
	if err != nil {
		err = fmt.Errorf("chat history remove failed, id = %s, user = %s, error = %v", messageId, userId, err)
		log.Log().Println(logs.Error(err).Extra(logs.F{"http", ctx.RequestPath(true)}).Trace())
		ctx.JSON(result{Success:false, Message:"remove failed."})
		return
	}
	ctx.JSON(result{Success:true, Message:"remove successed."})
	return
}

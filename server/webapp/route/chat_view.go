package route

import (
	"github.com/kataras/iris/context"
	"liulishuo/somechat/server/webapp/service"
	"fmt"
	"github.com/pharosnet/logs"
	"liulishuo/somechat/log"
	"liulishuo/somechat/server/webapp/conf"
	"strings"
	"math"
)

type chatViewData struct {
	Error string
	User userViewData
	ToUser userViewData
	MessageList []ChatMessage
	ChatAppServer string
	Index int64
}



func chatView(ctx context.Context)  {
	data := chatViewData{}
	from := ctx.FormValue("from")
	to := ctx.FormValue("to")
	if from == "" || to == "" {
		data.Error = "miss from or to"
		ctx.ViewData("view", data)
		ctx.View(view_chat)
		return
	}
	fromUser, fromUserGetErr := service.UserGetById(from)
	if fromUserGetErr != nil {
		err := fmt.Errorf("get user failed, from = %s, error = %v", from, fromUserGetErr)
		log.Log().Println(logs.Error(err).Extra(logs.F{"http", ctx.RequestPath(true)}).Trace())
		data.Error = err.Error()
		ctx.ViewData("view", data)
		ctx.View(view_chat)
		return
	}
	toUser, toUserErr := service.UserGetById(to)
	if toUserErr != nil {
		err := fmt.Errorf("get user failed, to = %s, error = %v", to, toUserErr)
		log.Log().Println(logs.Error(err).Extra(logs.F{"http", ctx.RequestPath(true)}).Trace())
		data.Error = err.Error()
		ctx.ViewData("view", data)
		ctx.View(view_chat)
		return
	}
	data.User = userViewData{Id:fromUser.Id, Name: fromUser.Name, Email: fromUser.Email}
	data.ToUser = userViewData{Id:toUser.Id, Name: toUser.Name, Email: toUser.Email}
	data.ChatAppServer = conf.Conf.Remote.Chat
	rows, rowsErr := service.ChatHistory(fromUser.Id, toUser.Id, math.MaxInt64, int64(0), int64(10))
	if rowsErr != nil {
		err := fmt.Errorf("get chat history failed, from = %s, to = %s, error = %v", fromUser.Id, toUser.Id, rowsErr)
		log.Log().Println(logs.Error(err).Extra(logs.F{"http", ctx.RequestPath(true)}).Trace())
		data.Error = err.Error()
		ctx.ViewData("view", data)
		ctx.View(view_chat)
		return
	}
	history := make([]ChatMessage, 0, len(rows))
	for _, row := range rows {
		content := strings.Replace(row.Content, "\r\n", "\n", -1)
		content = strings.Replace(row.Content, "\n\r", "\n", -1)
		content = strings.Replace(row.Content, "\r", "\n", -1)
		content = strings.Replace(row.Content, "\n", "<br>", -1)
		history = append(history, ChatMessage{
			Id: row.Id,
			Room:row.Room,
			UserId:row.UserId,
			UserName:row.UserName,
			UserEmail:row.UserEmail,
			Content:content,
			Index:row.Index,
			CreateTime:row.CreateTime.Time().Format("2006-01-02 15:04:05"),
			My:row.UserId == fromUser.Id,
		})
	}
	data.MessageList = history
	if len(history) > 0 {
		data.Index = history[0].Index
	}
	ctx.ViewData("view", data)
	// remove unread
	clearErr := service.ChatUnReadClearByFromTo(toUser.Id, fromUser.Id)
	if clearErr != nil {
		err := fmt.Errorf("unread clear failed, from = %s, to = %s, error = %v", fromUser.Id, toUser.Id, clearErr)
		log.Log().Println(logs.Error(err).Extra(logs.F{"http", ctx.RequestPath(true)}).Trace())
	}
	ctx.View(view_chat)
	return
}

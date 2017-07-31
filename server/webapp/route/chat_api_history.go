package route

import (
	"github.com/kataras/iris/context"
	"fmt"
	"github.com/pharosnet/logs"
	"strings"
	"liulishuo/somechat/log"
	"liulishuo/somechat/server/webapp/service"
)

type chatHistoryForm struct {
	FromId string	`form:"fromId"`
	ToId string	`form:"toId"`
	Offset int64	`form:"offset"`
	Limit int64	`form:"limit"`
	Index int64 	`form:"index"`

}

func chatApiHistory(ctx context.Context)  {
	form := new(chatHistoryForm)
	if formErr := ctx.ReadForm(form); formErr != nil {
		err := fmt.Errorf("form read failed, %v", formErr)
		log.Log().Println(logs.Error(err).Extra(logs.F{"http", ctx.RequestPath(true)}).Trace())
		ctx.JSON(result{Success:false, Message:"bad form"})
		return
	}
	form.FromId = strings.TrimSpace(form.FromId)
	form.ToId = strings.TrimSpace(form.ToId)
	if form.FromId == "" || form.ToId == "" {
		err := fmt.Errorf("form read failed, from = %s, to = %s", form.FromId, form.ToId)
		log.Log().Println(logs.Error(err).Extra(logs.F{"http", ctx.RequestPath(true)}).Trace())
		ctx.JSON(result{Success:false, Message:"bad form"})
		return
	}
	if form.Offset < int64(0) {
		form.Offset = int64(0)
	}
	if form.Limit < int64(1) {
		form.Limit = int64(10)
	}
	rows, rowsErr := service.ChatHistory(form.FromId, form.ToId, form.Index, form.Offset, form.Limit)
	if rowsErr != nil {
		err := fmt.Errorf("get chat history failed, from = %s, to = %s, error = %v", form.FromId, form.ToId, rowsErr)
		log.Log().Println(logs.Error(err).Extra(logs.F{"http", ctx.RequestPath(true)}).Trace())
		ctx.JSON(result{Success:false, Message:"bad form"})
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
			My:row.UserId == form.FromId,
		})
	}
	ctx.JSON(result{Success:true, Message:"", Data:map[string]interface{}{"list":history}})
	return
}

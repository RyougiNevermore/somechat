package service

import (
	"liulishuo/somechat/core/data"
	"fmt"
	"github.com/pharosnet/logs"
	"liulishuo/somechat/log"
	"liulishuo/somechat/core/common"
)

func ChatHistory(room string, offset, limit int64) ([]ChatMessage, error) {
	if room == "" {
		err := fmt.Errorf("get chat history failed. room is empty, room = %s", room)
		log.Log().Println(logs.Error(err).Extra(logs.F{"service", "ChatHistory"}).Trace())
		return nil, err
	}
	if offset < int64(0) {
		offset = int64(0)
	}
	if limit < int64(1) {
		limit = int64(10)
	}
	rows, rowsErr := data.ChatHistoryListByRoom(room, offset, limit)
	if rowsErr != nil {
		err := fmt.Errorf("get chat history failed. room = %s, error = %v", room, rowsErr)
		log.Log().Println(logs.Error(err).Extra(logs.F{"service", "ChatHistory"}).Trace())
		return nil, err
	}
	messages := make([]ChatMessage, 0, len(rows))
	for _, row := range rows {
		messages = append(messages, ChatMessage{
			Room:row.Room,
			UserId:row.UserId,
			UserName:row.UserName,
			UserEmail:row.UserEmail,
			Content:row.Content,
			Index:row.Index,
			CreateTime:common.JsonTimed(row.CreateTime),
		})
	}
	return messages, nil
}

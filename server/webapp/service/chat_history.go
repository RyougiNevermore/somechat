package service

import (
	"liulishuo/somechat/core/data"
	"fmt"
	"github.com/pharosnet/logs"
	"liulishuo/somechat/log"
	"liulishuo/somechat/core/common"
)

func ChatHistory(from, to string, index,  offset, limit int64) ([]ChatMessage, error) {
	room, roomErr := data.ChatHistoryGetRoom(from, to)
	if roomErr != nil {
		err := fmt.Errorf("get chat history failed. room is bad, from = %s, to = %s error = %s", from, to, roomErr)
		log.Log().Println(logs.Error(err).Extra(logs.F{"service", "ChatHistory"}).Trace())
		return nil, err
	}
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
	rows, rowsErr := data.ChatHistoryListByRoom(room, index, offset, limit)
	if rowsErr != nil {
		err := fmt.Errorf("get chat history failed. room = %s, error = %v", room, rowsErr)
		log.Log().Println(logs.Error(err).Extra(logs.F{"service", "ChatHistory"}).Trace())
		return nil, err
	}
	size := len(rows)
	for i := 0 ; i < size ; i ++ {
		for j := 0 ; j < size ; j ++  {
			if rows[i].Index < rows[j].Index {
				tmp := rows[i]
				rows[i] = rows[j]
				rows[j] = tmp
			}
		}
	}
	messages := make([]ChatMessage, 0, len(rows))
	for _, row := range rows {
		messages = append(messages, ChatMessage{
			Id: row.Id,
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

func ChatHistoryRemove(messageId, userId string) (error) {
	row, rowGetErr := data.ChatHistoryGetById(messageId)
	if rowGetErr != nil {
		err := fmt.Errorf("get chat history failed. id = %s, error = %v", messageId, rowGetErr)
		log.Log().Println(logs.Error(err).Extra(logs.F{"service", "ChatHistoryRemove"}).Trace())
		return err
	}
	if row.UserId != userId {
		err := fmt.Errorf("remove chat history failed. id = %s, message user id = %s, who = %s", messageId, row.UserId, userId)
		log.Log().Println(logs.Error(err).Extra(logs.F{"service", "ChatHistoryRemove"}).Trace())
		return err
	}
	tx, txBegErr := data.DAL().BeginTx()
	if txBegErr != nil {
		err := fmt.Errorf("remove chat history failed. tx begin failed, %v", txBegErr)
		log.Log().Println(logs.Error(err).Extra(logs.F{"service", "ChatHistoryRemove"}).Trace())
		return err
	}
	deleteAffected, deleteErr := data.ChatHistoryDelete(tx, row)
	if deleteErr != nil || deleteAffected == int64(0) {
		err := fmt.Errorf("remove chat history failed. table delete failed, row affected = %v, error = %v, tx roll back = %v", deleteAffected, deleteErr, tx.Rollback())
		log.Log().Println(logs.Error(err).Extra(logs.F{"service", "ChatHistoryRemove"}).Trace())
		return err
	}
	if cmtErr := tx.Commit(); cmtErr != nil {
		err := fmt.Errorf("remove chat history failed. tx commit failed, %v, tx roll back = %v", cmtErr, tx.Rollback())
		log.Log().Println(logs.Error(err).Extra(logs.F{"service", "ChatHistoryRemove"}).Trace())
		return err
	}
	return nil
}
package service

import (
	"liulishuo/somechat/core/data"
	"fmt"
	"github.com/pharosnet/logs"
	"liulishuo/somechat/log"
)

func ChatUnReadList(userId string) ([]ChatUnRead, error) {
	if userId == "" {
		err := fmt.Errorf("get chat unread failed. userId is empty, userId = %s", userId)
		log.Log().Println(logs.Error(err).Extra(logs.F{"service", "ChatUnRead"}).Trace())
		return nil, err
	}
	rows, rowsErr := data.ChatHistoryUnReadGetByUser(userId)
	if rowsErr != nil {
		err := fmt.Errorf("get chat unread failed. userId = %s, error = %v", userId, rowsErr)
		log.Log().Println(logs.Error(err).Extra(logs.F{"service", "ChatUnRead"}).Trace())
		return nil, err
	}
	items := make([]ChatUnRead, 0, len(rows))
	for _, row := range rows {
		items = append(items, ChatUnRead{
			Room:row.Room,
			FromUserId:row.FromUserId,
			ToUserId:row.ToUserId,
			Number:row.Number,
		})
	}
	return items, nil
}

func ChatUnReadClear(room string, toUserId string) (error) {
	if room == "" || toUserId == "" {
		err := fmt.Errorf("clear chat unread failed. room is empty, room = %s", room)
		log.Log().Println(logs.Error(err).Extra(logs.F{"service", "ChatUnRead"}).Trace())
		return err
	}
	row, rowGetErr := data.ChatHistoryUnReadGetByRoomAndToUserId(room, toUserId)
	if rowGetErr != nil {
		err := fmt.Errorf("clear chat unread failed. room = %s, error = %v", room, rowGetErr)
		log.Log().Println(logs.Error(err).Extra(logs.F{"service", "ChatUnRead"}).Trace())
		return  err
	}
	if row == nil || row.Id == "" {
		return nil
	}
	tx, txBegErr := data.DAL().BeginTx()
	if txBegErr != nil {
		err := fmt.Errorf("clear chat unread  failed. tx begin failed, %v", txBegErr)
		log.Log().Println(logs.Error(err).Extra(logs.F{"service", "ChatUnRead"}).Trace())
		return err
	}
	deleteResult, deleteErr := data.ChatHistoryUnReadDelete(tx, row)
	if deleteErr != nil || deleteResult == int64(0) {
		err := fmt.Errorf("clear chat unread  failed. table delete failed, row affected = %v, error = %v, tx roll back = %v", deleteResult, deleteErr, tx.Rollback())
		log.Log().Println(logs.Error(err).Extra(logs.F{"service", "ChatUnRead"}).Trace())
		return err
	}
	if cmtErr := tx.Commit(); cmtErr != nil {
		err := fmt.Errorf("clear chat unread  failed. tx commit failed, %v, tx roll back = %v", cmtErr, tx.Rollback())
		log.Log().Println(logs.Error(err).Extra(logs.F{"service", "ChatUnRead"}).Trace())
		return err
	}
	return nil
}

func ChatUnReadClearByFromTo(from string, to string) (error) {
	room, roomErr := data.ChatHistoryGetRoom(from, to)
	if roomErr != nil {
		return roomErr
	}
	return ChatUnReadClear(room, to)
}
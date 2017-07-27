package data

import (
	"github.com/niflheims-io/qb"
	"fmt"
	"github.com/pharosnet/logs"
	"liulishuo/somechat/log"
)

type ChatHistoryUnRead struct {
	Room string	`pk:"ROOM"` // from user id + to user id | group ??
	FromUserId string `col:"FROM_USER_ID"`
	ToUserId string	`col:"TO_USER_ID"`
	Number int64	`col:"NUMBER"`
	Version    int64     `version:"VERSION"`
}

func (r ChatHistoryUnRead) TableName() string {
	return "CHAT_HISTORY_UNREAD"
}

func ChatHistoryUnReadInsert(tx *qb.Tx, rows ...*ChatHistoryUnRead) (int64, error) {
	if tx == nil {
		err := fmt.Errorf("chat history unread insert failed, tx is nil, tx = %v", tx)
		return int64(0), err
	}
	affected := int64(0)
	for _, row := range rows {
		affectedOne, err := tx.Insert(row)
		if err != nil || affectedOne == int64(0) {
			err = fmt.Errorf("chat history unread insert failed, affected=%d, error = %v", affectedOne, err)
			log.Log().Println(logs.Error(err).Extra(logs.F{"sql", "CHAT_HISTORY_UNREAD"}).Trace())
			return int64(0), err
		}
		affected = affected + affectedOne
	}
	return affected, nil
}

func ChatHistoryUnReadUpdate(tx *qb.Tx, rows ...*ChatHistoryUnRead) (int64, error) {
	if tx == nil {
		err := fmt.Errorf("chat history unread update failed, tx is nil, tx = %v", tx)
		return int64(0), err
	}
	affected := int64(0)
	for _, row := range rows {
		affectedOne, err := tx.Update(row)
		if err != nil || affectedOne == int64(0) {
			err = fmt.Errorf("chat history unread update failed, affected=%d, error = %v", affectedOne, err)
			log.Log().Println(logs.Error(err).Extra(logs.F{"sql", "CHAT_HISTORY_UNREAD"}).Trace())
			return int64(0), err
		}
		affected = affected + affectedOne
	}
	return affected, nil
}

func ChatHistoryUnReadDelete(tx *qb.Tx, rows ...*ChatHistoryUnRead) (int64, error) {
	if tx == nil {
		err := fmt.Errorf("chat history unread delete failed, tx is nil, tx = %v", tx)
		return int64(0), err
	}
	affected := int64(0)
	for _, row := range rows {
		affectedOne, err := tx.Delete(row)
		if err != nil || affectedOne == int64(0) {
			err = fmt.Errorf("chat history unread delete failed, affected=%d, error = %v", affectedOne, err)
			log.Log().Println(logs.Error(err).Extra(logs.F{"sql", "CHAT_HISTORY_UNREAD"}).Trace())
			return int64(0), err
		}
		affected = affected + affectedOne
	}
	return affected, nil
}

func ChatHistoryUnReadGetByRoom(room string) (*ChatHistoryUnRead, error) {
	var one ChatHistoryUnRead
	if err := DAL().Query(`SELECT * FROM "CHAT_HISTORY_UNREAD" WHERE "ROOM" = $1 `, &room).One(&one); err != nil {
		err := fmt.Errorf("chat history unread get failed, can not find by room = %s, error = %v", room, err)
		return  nil, err
	}
	if one.Room == "" {
		err := fmt.Errorf("chat history unread get failed, can not find by room = %s", room)
		log.Log().Println(logs.Error(err).Extra(logs.F{"sql", "CHAT_HISTORY_UNREAD"}).Trace())
		return  nil, err
	}
	return &one, nil
}

func ChatHistoryUnReadGetByUser(userId string) ([]ChatHistoryUnRead, error) {
	var list []ChatHistoryUnRead
	if err := DAL().Query(`SELECT * FROM "CHAT_HISTORY_UNREAD" WHERE "TO_USER_ID" = $1`, &userId).List(&list); err != nil {
		err := fmt.Errorf("chat history unread get failed, can not find by user id = %s, error = %v", userId, err)
		return  nil, err
	}
	return list, nil
}
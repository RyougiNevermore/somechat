package data

import (
	"fmt"
	"github.com/niflheims-io/qb"
)

type ChatHistoryIndex struct {
	Room string	`pk:"ROOM"` // from user id + to user id | group ??
	LastIndex int64	`col:"LAST_INDEX"`
}

func (r ChatHistoryIndex) TableName() string {
	return "CHAT_HISTORY_INDEX"
}

func ChatHistoryIndexInsert(tx *qb.Tx, rows ...*ChatHistoryIndex) (int64, error) {
	if tx == nil {
		err := fmt.Errorf("chat history index insert failed, tx is nil, tx = %v", tx)
		return int64(0), err
	}
	affected := int64(0)
	for _, row := range rows {
		affectedOne, err := tx.Insert(row)
		if err != nil || affectedOne == int64(0) {
			err = fmt.Errorf("chat history index insert failed, affected=%d, error = %v", affectedOne, err)
			return int64(0), err
		}
		affected = affected + affectedOne
	}
	return affected, nil
}

func ChatHistoryIndexIncre(tx *qb.Tx, id string) (int64, error) {
	if tx == nil {
		err := fmt.Errorf("chat history index incre failed, tx is nil, tx = %v", tx)
		return int64(0), err
	}
	affected, err := tx.Exec(`UPDATE "CHAT_HISTORY_INDEX" SET "LAST_INDEX" = "LAST_INDEX" + 1 WHERE "ROOM" = $1`, id)
	if err != nil || affected == int64(0) {
		err = fmt.Errorf("chat history index incre failed, affected=%d, error = %v", affected, err)
		return int64(0), err
	}
	return affected, nil
}

func ChatHistoryIndexDelete(tx *qb.Tx, rows ...ChatHistoryIndex) (int64, error) {
	if tx == nil {
		err := fmt.Errorf("chat history index delete failed, tx is nil, tx = %v", tx)
		return int64(0), err
	}
	affected := int64(0)
	for _, row := range rows {
		affectedOne, err := tx.Delete(row)
		if err != nil || affectedOne == int64(0) {
			err = fmt.Errorf("chat history index delete failed, affected=%d, error = %v", affectedOne, err)
			return int64(0), err
		}
		affected = affected + affectedOne
	}
	return affected, nil
}
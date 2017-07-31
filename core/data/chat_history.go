package data

import (
	"time"
	"fmt"
	"github.com/niflheims-io/qb"
	"github.com/pharosnet/logs"
	"liulishuo/somechat/log"
	"crypto/md5"
	"encoding/hex"
)

type ChatHistory struct {
	Id string 	`pk:"ID"`
	Room string	`col:"ROOM"` // from user id + to user id | group ??
	UserId string	`col:"USER_ID"`
	UserName string	`col:"USER_NAME"`
	UserEmail string `col:"USER_EMAIL"`
	Content string	`col:"CONTENT"`
	Index int64	`col:"INDEX"`
	CreateTime time.Time	`col:"CREATE_TIME"`
}

func (r ChatHistory) TableName() string {
	return "CHAT_HISTORY"
}

func ChatHistoryRoomBuild(userIds ...string) string {
	room := ""
	for _, userId := range userIds {
		room = room + "," + userId
	}
	room = room[1:]
	md5Hash := md5.New()
	md5Hash.Write([]byte(room))
	room = hex.EncodeToString(md5Hash.Sum(nil))
	return room
}

func ChatHistoryGetIndex() (int64, error) {
	seq := int64(0)
	seqErr := DAL().Query(`SELECT nextval('chat_history_index_seq')`).Int(&seq)
	if seqErr != nil {
		return int64(0), seqErr
	}
	return seq, nil
}

func ChatHistoryGetRoom(from, to string) (string, error) {
	querySQL := `SELECT "MAGIC_CODE" FROM "CONTACT" WHERE "OWNER" = $1 AND "USER_ID" = $2 `
	fromMagic := ""
	fromMagicErr := DAL().Query(querySQL, &from, &to).String(&fromMagic)
	if fromMagicErr != nil {
		return "", fromMagicErr
	}
	if fromMagic == "" {
		return "", fmt.Errorf("from(%s) and to(%s) have no same magic code in contact", from, to)
	}
	toMagic := ""
	toMagicErr := DAL().Query(querySQL, &to, &from).String(&toMagic)
	if toMagicErr != nil {
		return "", toMagicErr
	}
	if toMagic == "" {
		return "", fmt.Errorf("from(%s) and to(%s) have no same magic code in contact", from, to)
	}
	if fromMagic != toMagic {
		return "", fmt.Errorf("from(%s) and to(%s) have no same magic code in contact", from, to)
	}
	return fromMagic, nil
}

func ChatHistoryInsert(tx *qb.Tx, rows ...*ChatHistory) (int64, error) {
	if tx == nil {
		err := fmt.Errorf("chat history insert failed, tx is nil, tx = %v", tx)
		return int64(0), err
	}
	affected := int64(0)
	for _, row := range rows {
		affectedOne, err := tx.Insert(row)
		if err != nil || affectedOne == int64(0) {
			err = fmt.Errorf("chat history insert failed, affected=%d, error = %v", affectedOne, err)
			log.Log().Println(logs.Error(err).Extra(logs.F{"sql", "CHAT_HISTORY"}).Trace())
			return int64(0), err
		}
		affected = affected + affectedOne
	}
	return affected, nil
}

func ChatHistoryDelete(tx *qb.Tx, rows ...*ChatHistory) (int64, error) {
	if tx == nil {
		err := fmt.Errorf("chat history delete failed, tx is nil, tx = %v", tx)
		return int64(0), err
	}
	affected := int64(0)
	for _, row := range rows {
		affectedOne, err := tx.Delete(row)
		if err != nil || affectedOne == int64(0) {
			err = fmt.Errorf("chat history delete failed, affected=%d, error = %v", affectedOne, err)
			log.Log().Println(logs.Error(err).Extra(logs.F{"sql", "CHAT_HISTORY"}).Trace())
			return int64(0), err
		}
		affected = affected + affectedOne
	}
	return affected, nil
}

func ChatHistoryListByRoom(room string, index, offset, limit int64) ([]ChatHistory, error) {
	var list []ChatHistory
	if err := DAL().Query(
		`SELECT * FROM "CHAT_HISTORY" WHERE "ROOM" = $1 AND "INDEX" < $2 ORDER BY "CREATE_TIME" DESC,"INDEX" DESC LIMIT $3 OFFSET 0`,
		&room, &index, &limit,
	).List(&list); err != nil {
		err := fmt.Errorf("chat history list failed, can not find by room = %s, offset = %v, limit = %v, error = %v", room, offset, limit, err)
		log.Log().Println(logs.Error(err).Extra(logs.F{"sql", "CHAT_HISTORY"}).Trace())
		return  nil, err
	}
	return list, nil
}

func ChatHistoryGetById(id string) (*ChatHistory, error) {
	var one ChatHistory
	if err := DAL().Query(`SELECT * FROM "CHAT_HISTORY" WHERE "ID" = $1`, &id).One(&one); err != nil {
		err := fmt.Errorf("chat history get by id failed, id = %s, error = %v", id, err)
		log.Log().Println(logs.Error(err).Extra(logs.F{"sql", "CHAT_HISTORY"}).Trace())
		return  nil, err
	}
	if one.Id == "" {
		err := fmt.Errorf("chat history get by id failed, id = %s", id)
		log.Log().Println(logs.Error(err).Extra(logs.F{"sql", "CHAT_HISTORY"}).Trace())
		return  nil, err
	}
	return &one, nil
}
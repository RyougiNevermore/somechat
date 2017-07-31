package agent

import (
	"encoding/json"
	"liulishuo/somechat/core/data"
	"fmt"
	"liulishuo/somechat/log"
	"github.com/pharosnet/logs"
	"time"
)

type Message struct {
	Head map[string]string        `json:"head"`
	Body string        `json:"body"`
}

func NewMessage(p []byte) (*Message, error) {
	msg := Message{}
	err := json.Unmarshal(p, &msg)
	return &msg, err
}

func (m *Message) Encode() ([]byte, error) {
	return json.Marshal(m)
}

func (m *Message) GetFrom() string {
	return m.Head["from"]
}

//content-type
func (m *Message) GetContentType() string {
	return m.Head["contentType"]
}

func (m *Message) SetContentType(val string) {
	m.Head["contentType"] = val
}

func (m *Message) GetTarget() string {
	return m.Head["target"]
}

func (m *Message) GetId() string {
	return m.Head["id"]
}

func (m *Message) GetTo() string {
	return m.Head["to"]
}

func (m *Message) GetKind() string {
	return m.Head["kind"]
}

type MessageStore struct {
	agentChanMap map[string] chan Message // FAKE NO IMPLEMENT
}

func (s *MessageStore) Save(from, target string, msg *Message) error {
	user, userGetErr := data.UserGetById(from)
	if userGetErr != nil {
		err := fmt.Errorf("message save failed, get from user failed, from = %s, error = %v", from, userGetErr)
		log.Log().Println(logs.Error(err).Extra(logs.F{"websocket", "MessageStore"}).Trace())
		return err
	}
	//msgBytes, msgEncodeErr := msg.Encode()
	//if msgEncodeErr != nil {
	//	err := fmt.Errorf("message save failed, encode message failed, msg = %v, error = %v", msg, msgEncodeErr)
	//	log.Log().Println(logs.Error(err).Extra(logs.F{"websocket", "MessageStore"}).Trace())
	//	return err
	//}
	index, indexErr := data.ChatHistoryGetIndex()
	if indexErr != nil {
		err := fmt.Errorf("message save failed, get index failed, error = %v", indexErr)
		log.Log().Println(logs.Error(err).Extra(logs.F{"websocket", "MessageStore"}).Trace())
		return err
	}
	room, roomErr := data.ChatHistoryGetRoom(from, target)
	if roomErr != nil {
		err := fmt.Errorf("message save failed, get room failed, error = %v", roomErr)
		log.Log().Println(logs.Error(err).Extra(logs.F{"websocket", "MessageStore"}).Trace())
		return err
	}
	row := new(data.ChatHistory)
	row.Id = msg.GetId()
	row.Room = room
	row.UserId = user.Id
	row.UserName = user.Name
	row.UserEmail = user.Email
	row.CreateTime = time.Now()
	row.Content = msg.Body
	row.Index = index
	// insert
	tx, txBegErr := data.DAL().BeginTx()
	if txBegErr != nil {
		err := fmt.Errorf("message save failed, tx begin failed, error = %v", txBegErr)
		log.Log().Println(logs.Error(err).Extra(logs.F{"websocket", "MessageStore"}).Trace())
		return err
	}
	affected, insertErr := data.ChatHistoryInsert(tx, row)
	if affected == int64(0) || insertErr != nil {
		err := fmt.Errorf("message save failed, insert chat history failed, row = %v, affected = %v, error = %v, tx roll back = %v",
			row, affected, insertErr, tx.Rollback())
		log.Log().Println(logs.Error(err).Extra(logs.F{"websocket", "MessageStore"}).Trace())
		return err
	}
	if cmtErr := tx.Commit(); cmtErr != nil {
		err := fmt.Errorf("message save failed, insert chat history failed. tx commit failed, %v, tx roll back = %v", cmtErr, tx.Rollback())
		log.Log().Println(logs.Error(err).Extra(logs.F{"websocket", "MessageStore"}).Trace())
		return err
	}
	return nil
}
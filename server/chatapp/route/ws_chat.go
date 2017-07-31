package route

import (
	"github.com/kataras/iris/websocket"
	"fmt"
	"github.com/pharosnet/logs"
	"liulishuo/somechat/server/chatapp/agent"
	"strings"
	"liulishuo/somechat/log"
	"liulishuo/somechat/core/data"
	"github.com/pharosnet/auid"
)

func webSocketChatHandler(ws websocket.Connection)  {
	ws.On(ws_topic_chat, func(str string) {
		log.Log().Println(logs.Debug(str).Extra(logs.F{"websocket", "chat"}).Trace())
		p := []byte(str)
		if len(p) == 0 {
			log.Log().Println(logs.Debug("topic = %s, recive msg is empty.", ws_topic_chat).Trace())
			return
		}
		msg, msgErr := agent.NewMessage(p)
		if msgErr != nil {
			err := fmt.Errorf("topic = %s, decode msg failed, msg  = %s, error = %v", ws_topic_chat, string(p), msgErr)
			log.Log().Println(logs.Error(err).Trace())
			return
		}
		from := strings.TrimSpace(msg.GetFrom())
		to := strings.TrimSpace(msg.GetTo())
		if from == "" || to == "" {
			err := fmt.Errorf("topic = %s, miss from or target, from = %s, to = %s, msg = %s", ws_topic_chat, from, to, string(p))
			log.Log().Println(logs.Error(err).Trace())
			return
		}
		send, sendErr := gateway.Send(from, to, msg, true)
		log.Log().Println(logs.Debugf("send from = %s, to = %s, %v", from, to, send).Extra(logs.F{"websocket", "chat"}).Trace())
		if sendErr != nil {
			err := fmt.Errorf("topic = %s, send msg failed, from = %s, to = %s error = %v", ws_topic_chat, from, to, sendErr)
			log.Log().Println(logs.Error(err).Trace())
			return
		}
		if !send {
			// save unread
			room, roomErr := data.ChatHistoryGetRoom(from, to)
			if roomErr != nil {
				err := fmt.Errorf("unread save failed, get room failed, error = %v", roomErr)
				log.Log().Println(logs.Error(err).Extra(logs.F{"websocket", "chat"}).Trace())
				return
			}
			unreadRow, unreadRowGetErr := data.ChatHistoryUnReadGetByRoomAndToUserId(room, to)
			if unreadRowGetErr != nil {
				err := fmt.Errorf("unread save failed, get unread failed, error = %v", unreadRowGetErr)
				log.Log().Println(logs.Error(err).Extra(logs.F{"websocket", "chat"}).Trace())
				return
			}

			tx, txBegErr := data.DAL().BeginTx()
			if txBegErr != nil {
				err := fmt.Errorf("unread save failed. tx begin failed, %v", txBegErr)
				log.Log().Println(logs.Error(err).Extra(logs.F{"websocket", "chat"}).Trace())
				return
			}

			if unreadRow == nil || unreadRow.Id == "" {
				unreadRow = new(data.ChatHistoryUnRead)
				unreadRow.Id = auid.NewAuid()
				unreadRow.Room = room
				unreadRow.Number = int64(1)
				unreadRow.FromUserId = from
				unreadRow.ToUserId = to
				unreadRow.Version = int64(1)
				insertAffect, insertErr := data.ChatHistoryUnReadInsert(tx, unreadRow)
				if insertErr != nil || insertAffect == int64(0) {
					err := fmt.Errorf("save chat unread  failed. table insert failed, row affected = %v, error = %v, tx roll back = %v", insertAffect, insertErr, tx.Rollback())
					log.Log().Println(logs.Error(err).Extra(logs.F{"websocket", "chat"}).Trace())
					return
				}
			} else {
				unreadRow.Number = unreadRow.Number + int64(1)
				updateAffect, updateErr := data.ChatHistoryUnReadUpdate(tx, unreadRow)
				if updateErr != nil || updateAffect == int64(0) {
					err := fmt.Errorf("save chat unread  failed. table update failed, row affected = %v, error = %v, tx roll back = %v", updateAffect, updateErr, tx.Rollback())
					log.Log().Println(logs.Error(err).Extra(logs.F{"websocket", "chat"}).Trace())
					return
				}
			}

			if cmtErr := tx.Commit(); cmtErr != nil {
				err := fmt.Errorf("save chat unread  failed. tx commit failed, %v, tx roll back = %v", cmtErr, tx.Rollback())
				log.Log().Println(logs.Error(err).Extra(logs.F{"websocket", "chat"}).Trace())
				return
			}

			// make notify
			notifyMSg := new(agent.Message)
			notifyMSg.Head = make(map[string]string)
			notifyMSg.SetContentType("unread")
			notifyMSg.Head["from"] = from
			notifyMSg.Body = `1`
			err := gateway.Notify(to, ws_topic_notify, notifyMSg)
			if err != nil {
				err = fmt.Errorf("send notify failed, error = %v", err)
				log.Log().Println(logs.Error(err).Trace())
				return
			}
		}
	})
}

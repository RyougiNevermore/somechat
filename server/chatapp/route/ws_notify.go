package route

import (
	"github.com/kataras/iris/websocket"
	"github.com/pharosnet/logs"
	"liulishuo/somechat/server/chatapp/agent"
	"fmt"
	"strings"
	"liulishuo/somechat/log"
)

func webSocketNotifyHandler(ws websocket.Connection)  {
	ws.On(ws_topic_notify, func(str string) {
		p := []byte(str)
		if len(p) == 0 {
			log.Log().Println(logs.Debug("topic = %s, recive msg is empty.", ws_topic_notify).Trace())
			return
		}
		msg, msgErr := agent.NewMessage(p)
		if msgErr != nil {
			err := fmt.Errorf("topic = %s, decode msg failed, msg  = %s, error = %v", ws_topic_notify, string(p), msgErr)
			log.Log().Println(logs.Error(err).Trace())
			return
		}
		if msg.GetContentType() == "connect" {
			notifyOnConnect(ws, msg)
			return
		}
		if msg.GetContentType() == "chat_closed" {
			notifyOnChatClose(ws, msg)
			return
		}
	})
}

func notifyOnConnect(ws websocket.Connection, msg *agent.Message) {
	from := strings.TrimSpace(msg.GetFrom())
	if from == ""  {
		err := fmt.Errorf("topic = %s, miss from or target, from = %s", ws_topic_notify, from)
		log.Log().Println(logs.Error(err).Trace())
		return
	}
	agentInstance := agent.NewAgent(ws_topic_notify, from, ws_topic_notify, ws)
	bus.Register(agentInstance)
	lock.Lock()
	wsAgentRef[ws.ID()] = agentInstance.Id
	lock.Unlock()
	log.Log().Println(logs.Debugf("topic = %s, register agent ok, ws id = %s, agent id = %s", ws_topic_notify, ws.ID(), agentInstance.Id))
}

func notifyOnChatClose(ws websocket.Connection, msg *agent.Message) {
	from := strings.TrimSpace(msg.GetFrom())
	if from == ""  {
		err := fmt.Errorf("topic = %s, miss from or target, from = %s", ws_topic_notify, from)
		log.Log().Println(logs.Error(err).Trace())
		return
	}
	err := gateway.Notify(from, ws_topic_notify, msg)
	if err != nil {
		err = fmt.Errorf("send notify failed, error = %v", err)
		log.Log().Println(logs.Error(err).Trace())
		return
	}
}
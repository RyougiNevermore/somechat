package route

import (
	"github.com/kataras/iris/websocket"
	"github.com/pharosnet/logs"
	"liulishuo/somechat/log"
	"liulishuo/somechat/server/chatapp/agent"
	"fmt"
)

func webSocketDisconnectHandler(ws websocket.Connection)  {
	ws.OnDisconnect(func() {
		lock.Lock()
		defer lock.Unlock()
		agentId := wsAgentRef[ws.ID()]
		if agentId == "" {
			return
		}
		delete(wsAgentRef, ws.ID())
		agentInstance, hasAgent := bus.Find(agentId)
		if !hasAgent {
			return
		}
		bus.UnRegister(agentId)
		log.Log().Println(logs.Debug("disconnect ws.id = %s, agent = %s", ws.ID(), agentId))
		if agentInstance.Topic == ws_topic_chat {
			from, to := agent.DecodeAgentId(agentId)
			if from == "" {
				return
			}
			msg := new(agent.Message)
			msg.Head = make(map[string]string)
			msg.SetContentType("chat_closed")
			msg.Head["from"] = from
			msg.Head["to"] = to
			notifyErr := gateway.Notify(from, ws_topic_notify, msg)
			if notifyErr != nil {
				err := fmt.Errorf("send notify failed, error = %v", notifyErr)
				log.Log().Println(logs.Error(err).Trace())
				return
			}
		}
	})
}

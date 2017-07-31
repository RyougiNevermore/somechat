package route

import (
	"github.com/kataras/iris/websocket"
	"github.com/pharosnet/logs"
	"strings"
	"fmt"
	"liulishuo/somechat/server/chatapp/agent"
	"liulishuo/somechat/log"
)

func webSocketRoomHandler(ws websocket.Connection)  {
	ws.On(ws_topic_room, func(str string) {
		p := []byte(str)
		if len(p) == 0 {
			log.Log().Println(logs.Debug("topic = %s, recive msg is empty.", ws_topic_room).Trace())
			return
		}
		msg, msgErr := agent.NewMessage(p)
		if msgErr != nil {
			err := fmt.Errorf("topic = %s, decode msg failed, msg  = %s, error = %v", ws_topic_room, string(p), msgErr)
			log.Log().Println(logs.Error(err).Trace())
			return
		}
		from := strings.TrimSpace(msg.GetFrom())
		to := strings.TrimSpace(msg.GetTo())
		if from == "" || to == "" {
			err := fmt.Errorf("topic = %s, miss from or target, from = %s, to = %s, msg = %s", ws_topic_room, from, to, string(p))
			log.Log().Println(logs.Error(err).Trace())
			return
		}
		agentInstance := agent.NewAgent(ws_topic_chat, from, to, ws)
		bus.Register(agentInstance)
		lock.Lock()
		wsAgentRef[ws.ID()] = agentInstance.Id
		lock.Unlock()
		log.Log().Println(logs.Debugf("topic = %s, register agent ok, ws id = %s, agent id = %s", ws_topic_room, ws.ID(), agentInstance.Id))
	})
}

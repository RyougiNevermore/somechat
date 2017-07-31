package agent

import (
	"github.com/kataras/iris/websocket"
	"liulishuo/somechat/log"
	"github.com/pharosnet/logs"
	"fmt"
	"strings"
)

type Agent struct {
	Id string
	Topic string
	Conn websocket.Connection
}

func NewAgent(topic, from, target string, conn websocket.Connection) *Agent {
	return &Agent{
		Id:NewAgentId(from, target),
		Conn:conn,
		Topic:topic,
	}
}

func NewAgentId(from, target string) string  {
	return fmt.Sprintf("%s:%s", from, target)
}

func DecodeAgentId(agentId string) (string, string) {
	items := strings.Split(agentId, ":")
	return items[0], items[1]
}

func (a *Agent) Emit(msg *Message) (int64, error) {
	p, encodeErr := msg.Encode()
	if encodeErr != nil {
		err := fmt.Errorf("message encode failed, error = %v, content = %s", encodeErr, string(p))
		log.Log().Println(logs.Error(err).Extra(logs.F{"websocket", "agent"}).Trace())
		return int64(0), err
	}
	if emitErr := a.Conn.Emit(a.Topic, p); emitErr != nil {
		err := fmt.Errorf("websocket emit failed, error = %v", emitErr)
		log.Log().Println(logs.Error(err).Extra(logs.F{"websocket", "agent"}).Trace())
		return 0, err
	}
	return int64(len(p)), nil
}

func (a *Agent) On(handle func (*Message))  {
	a.Conn.On(a.Topic, func(p []byte) {
		msg, msgErr := NewMessage(p)
		if msgErr != nil {
			err := fmt.Errorf("message decode failed, error = %v, content = %s", msgErr, string(p))
			log.Log().Println(logs.Error(err).Extra(logs.F{"websocket", "agent"}).Trace())
			if emitErr := a.Conn.Emit("notify", "bad message. contant = " + string(p)); emitErr != nil {
				err := fmt.Errorf("websocket emit failed, error = %v", emitErr)
				log.Log().Println(logs.Error(err).Extra(logs.F{"websocket", "agent"}).Trace())
			}
			return
		}
		handle(msg)
	})
}
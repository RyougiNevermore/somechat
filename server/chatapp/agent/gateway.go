package agent

import (
	"fmt"
	"liulishuo/somechat/log"
	"github.com/pharosnet/logs"
	"time"
)

type Gateway struct {
	bus *AgentLocalBus
	store *MessageStore
}

func NewGateway(bus *AgentLocalBus, store *MessageStore) (*Gateway) {
	g := Gateway{}
	g.bus = bus
	g.store = store
	return &g
}

// if hasAgent == false, then send to {target>>>}
func (g *Gateway) Send(from, target string, msg *Message, saveFlag bool) (bool, error) {
	agent, hasAgent := g.bus.Find(NewAgentId(target, from))
	log.Log().Println(logs.Debugf("to agent, has = %v, target = %s, from = %s", hasAgent, target, from).Extra(logs.F{"websocket", "chat"}).Trace())
	var sendFlag bool = false
	var err error
	if hasAgent {
		msg.Head["createTime"] = time.Now().Format("2006-01-02 15:04:05")
		_, emitErr := agent.Emit(msg)
		if emitErr != nil {
			err = fmt.Errorf("gateway send message failed, from = %s, target = %s, error = %v", from, target, emitErr)
			log.Log().Println(logs.Error(err).Extra(logs.F{"websocket", "Gateway"}).Trace())
			sendFlag = false
			if !saveFlag {
				return sendFlag, err
			}
		}
		sendFlag = true
		// to store message
	}
	if !saveFlag {
		return sendFlag, err
	}
	saveErr := g.store.Save(from, target, msg)
	if saveErr != nil {
		err = fmt.Errorf("gateway save message failed, from = %s, target = %s, error = %v", from, target, saveErr)
		log.Log().Println(logs.Error(err).Extra(logs.F{"websocket", "Gateway"}).Trace())
		if sendFlag {
			// TODO TRY SAVE AGAIN
			// TODO MAKE AN EDA TO SAVE MSG WITH CHANNEL (ONE AGENT ONE CHANNEL)
		}
	}
	return sendFlag, err
}

func (g *Gateway) Notify(to, target string, msg *Message) (error) {
	agent, hasAgent := g.bus.Find(NewAgentId(to, target))
	var err error
	if hasAgent {
		_, emitErr := agent.Emit(msg)
		if emitErr != nil {
			err = fmt.Errorf("gateway send notify failed, to = %s, target = %s, error = %v", to, target, emitErr)
			log.Log().Println(logs.Error(err).Extra(logs.F{"websocket", "Gateway"}).Trace())
		}
	}
	return err
}
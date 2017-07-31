package route

import (
	"liulishuo/somechat/server/chatapp/agent"
	"sync"
	"github.com/kataras/iris/websocket"
)

const (
	ws_topic_chat = "chat"
	ws_topic_notify = "notify"
	ws_topic_room = "room"
)

var bus *agent.AgentLocalBus
var gateway *agent.Gateway

var wsAgentRef map[string]string
var lock sync.Mutex


func webSocketHandler(ws websocket.Connection) {
	webSocketRoomHandler(ws)
	webSocketChatHandler(ws)
	webSocketNotifyHandler(ws)
	webSocketDisconnectHandler(ws)
}
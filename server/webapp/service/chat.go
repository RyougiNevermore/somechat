package service

import (
	"liulishuo/somechat/core/common"
)

type ChatMessage struct {
	Room string        `json:"room"`
	UserId string        `json:"userId"`
	UserName string        `json:"userName"`
	UserEmail string        `json:"userEmail"`
	Content string        `json:"content"`
	Index int64        `json:"index"`
	CreateTime common.JsonTimed        `json:"createTime"`
}

type ChatUnRead struct {
	Room string `json:"room"`
	FromUserId string `json:"fromUserId"`
	ToUserId string `json:"toUserId"`
	Number int64 `json:"number"`
}
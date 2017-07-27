package service

import (
	"liulishuo/somechat/core/common"
)

type ContactAddRequest struct {
	Id string `json:"id"`
	FromId string        `json:"fromId"`
	FromName string        `json:"fromName"`
	FromEmail string        `json:"fromEmail"`
	ToId string        `json:"toId"`
	ToName string        `json:"toName"`
	ToEmail string        `json:"toEmail"`
	CreateTime common.JsonTimed        `json:"createTime"`
}

type Contact struct {
	Id string `json:"id"`
	Owner string       `json:"owner"`
	UserId string       `json:"userId"`
	UserName string       `json:"userName"`
	UserEmail string      `json:"userEmail"`
}
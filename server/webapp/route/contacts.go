package route

import (
	"liulishuo/somechat/core/common"
	"liulishuo/somechat/server/webapp/service"
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

func contactAddRequestListCopyFromService(list []service.ContactAddRequest) []ContactAddRequest {
	results := make([]ContactAddRequest, 0, len(list))
	for _, item := range list {
		results = append(results, ContactAddRequest{
			Id: item.Id,
			FromId: item.FromId,
			FromName: item.FromName,
			FromEmail: item.FromEmail,
			ToId: item.ToId,
			ToName: item.ToName,
			ToEmail: item.ToEmail,
			CreateTime: item.CreateTime,
		})
	}
	return results
}

type Contact struct {
	Id string `json:"id"`
	Owner string       `json:"owner"`
	UserId string       `json:"userId"`
	UserName string       `json:"userName"`
	UserEmail string      `json:"userEmail"`
	Unread UnreadChatMessage        `json:"unread"`
}



func contactListCopyFromService(list []service.Contact) []Contact {
	results := make([]Contact, 0, len(list))
	for _, item := range list {
		results = append(results, Contact{
			Id: item.Id,
			Owner: item.Owner,
			UserId: item.UserId,
			UserName: item.UserName,
			UserEmail: item.UserEmail,
		})
	}
	return results
}

func contactListMergeUnreadMessageList(contacts []Contact, messages []service.ChatUnRead) []Contact {
	for i, contact := range contacts {
		for _, unread := range messages {
			if contact.UserId != unread.FromUserId {
				continue
			}
			contact.Unread = UnreadChatMessage{
				Room: unread.Room,
				Number: unread.Number,
			}
			contacts[i] = contact
		}
	}
	return contacts
}
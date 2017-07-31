package route


type UnreadChatMessage struct {
	Room string `json:"room"`
	Number int64 `json:"number"`
}

type ChatMessage struct {
	Id string        `json:"id"`
	Room string        `json:"room"`
	UserId string        `json:"userId"`
	UserName string        `json:"userName"`
	UserEmail string        `json:"userEmail"`
	Content string        `json:"content"`
	Index int64        `json:"index"`
	My 	bool `json:"my"`
	CreateTime string        `json:"createTime"`
} 
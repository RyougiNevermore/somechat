package route

const (
	view_login = "login.html"
	view_register = "register.html"
	view_contacts = "contacts.html"
	view_chat = "chat.html"
)

// common view data
type userViewData struct {
	Id string
	Name string
	Email string
}
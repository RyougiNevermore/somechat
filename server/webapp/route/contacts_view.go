package route

import (
	"github.com/kataras/iris/context"
	"fmt"
	"github.com/pharosnet/logs"
	"liulishuo/somechat/log"
	"liulishuo/somechat/server/webapp/service"
	"strings"
	"liulishuo/somechat/server/webapp/conf"
)

type contactsForm struct {
	UserId string `form:"userId"`
}

type contactsViewData struct {
	Error string
	User userViewData
	ContactAddRequestToUserList []ContactAddRequest
	ContactAddRequestFromUserList []ContactAddRequest
	ContactList []Contact
	ChatAppServer string
}

func contactsView(ctx context.Context)  {
	viewData := contactsViewData{}
	viewData.ChatAppServer = conf.Conf.Remote.Chat
	form := new(contactsForm)
	if formErr := ctx.ReadForm(form); formErr != nil {
		err := fmt.Errorf("form read failed, %v", formErr)
		log.Log().Println(logs.Error(err).Extra(logs.F{"http", ctx.RequestPath(true)}).Trace())
		viewData.Error = err.Error()
		ctx.ViewData("view", viewData)
		ctx.View(view_contacts)
		return
	}
	form.UserId = strings.TrimSpace(form.UserId)
	if form.UserId == "" {
		err := fmt.Errorf("can not find userId at form, userId = %s", form.UserId)
		log.Log().Println(logs.Error(err).Extra(logs.F{"http", ctx.RequestPath(true)}).Trace())
		viewData.Error = err.Error()
		ctx.ViewData("view", viewData)
		ctx.View(view_contacts)
		return
	}
	// get user by userId
	user, userGetErr := service.UserGetById(form.UserId)
	if userGetErr != nil {
		err := fmt.Errorf("can not find user by userId, userId = %s, error = %v", form.UserId, userGetErr)
		log.Log().Println(logs.Error(err).Extra(logs.F{"http", ctx.RequestPath(true)}).Trace())
		viewData.Error = "cant not find user, id = " + form.UserId
		ctx.ViewData("view", viewData)
		ctx.View(view_contacts)
		return
	}
	viewData.User = userViewData{Id:user.Id, Name:user.Name, Email:user.Email}
	// contact add request by to user,
	contactAddRequestToUserList, contactAddRequestToUserListErr := service.ContactAddRequestListByToUser(user.Id)
	if contactAddRequestToUserListErr != nil {
		err := fmt.Errorf("can not find contact add req of to user by userId, userId = %s, error = %v", form.UserId, contactAddRequestToUserListErr)
		log.Log().Println(logs.Error(err).Extra(logs.F{"http", ctx.RequestPath(true)}).Trace())
		viewData.Error = "cant not find contact add req that to me, id = " + form.UserId
		ctx.ViewData("view", viewData)
		ctx.View(view_contacts)
		return
	}
	viewData.ContactAddRequestToUserList = contactAddRequestListCopyFromService(contactAddRequestToUserList)
	// contact add request by from user
	contactAddRequestFromUserList, contactAddRequestFromUserListErr := service.ContactAddRequestListByFromUser(user.Id)
	if contactAddRequestFromUserListErr != nil {
		err := fmt.Errorf("can not find contact add req of from user by userId, userId = %s, error = %v", form.UserId, contactAddRequestFromUserListErr)
		log.Log().Println(logs.Error(err).Extra(logs.F{"http", ctx.RequestPath(true)}).Trace())
		viewData.Error = "cant not find contact add req that i send, id = " + form.UserId
		ctx.ViewData("view", viewData)
		ctx.View(view_contacts)
		return
	}
	viewData.ContactAddRequestFromUserList = contactAddRequestListCopyFromService(contactAddRequestFromUserList)
	// contact list
	contactList, contactListErr := service.ContactListByUserId(user.Id)
	if contactListErr != nil {
		err := fmt.Errorf("can not find contact list of user, userId = %s, error = %v", form.UserId, contactListErr)
		log.Log().Println(logs.Error(err).Extra(logs.F{"http", ctx.RequestPath(true)}).Trace())
		viewData.Error = "cant not find contact list, id = " + form.UserId
		ctx.ViewData("view", viewData)
		ctx.View(view_contacts)
		return
	}
	viewData.ContactList = contactListCopyFromService(contactList)
	// unread
	unreadMessages, unreadMessagesGetErr := service.ChatUnReadList(user.Id)
	if unreadMessagesGetErr != nil {
		err := fmt.Errorf("find unread chat message list of user failed, userId = %s, error = %v", form.UserId, unreadMessagesGetErr)
		log.Log().Println(logs.Error(err).Extra(logs.F{"http", ctx.RequestPath(true)}).Trace())
		viewData.Error = "find chat message unread failed, id = " + form.UserId
		ctx.ViewData("view", viewData)
		ctx.View(view_contacts)
		return
	}
	viewData.ContactList = contactListMergeUnreadMessageList(viewData.ContactList, unreadMessages)

	ctx.ViewData("view", viewData)
	if viewErr := ctx.View(view_contacts); viewErr != nil {
		log.Log().Println(logs.Error(viewErr).Trace())
	}
	return
}

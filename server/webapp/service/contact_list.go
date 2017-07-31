package service

import (
	"fmt"
	"github.com/pharosnet/logs"
	"liulishuo/somechat/core/data"
	"liulishuo/somechat/log"
)

func ContactListByUserId(userId string) ([]Contact, error) {
	if userId == "" {
		err := fmt.Errorf("user id is empty, userId = (%s)", userId)
		log.Log().Println(logs.Error(err).Extra(logs.F{"service", "ContactListByUserId"}).Trace())
		return nil, err
	}
	rows, rowsErr := data.ContactListByUserId(userId)
	if rowsErr != nil {
		err := fmt.Errorf("list contact failed, userId = %s, error = %v", userId, rowsErr)
		log.Log().Println(logs.Error(err).Extra(logs.F{"service", "ContactListByUserId"}).Trace())
		return nil, err
	}
	results := make([]Contact, 0, len(rows))
	for _, row := range rows {
		results = append(results, Contact{
			Id: row.Id,
			Owner:row.Owner,
			UserId:row.UserId,
			UserName:row.UserName,
			UserEmail:row.UserEmail,
		})
	}
	return results, nil

}

func ContactGetOneByOwnerAndUserId(owner, userId string) (*Contact, error) {
	if owner == "" || userId == "" {
		err := fmt.Errorf("get contact failed, owner = (%s) userId = (%s)", owner, userId)
		log.Log().Println(logs.Error(err).Extra(logs.F{"service", "ContactGetOneByOwnerAndUser"}).Trace())
		return nil, err
	}
	row, getErr := data.ContactGetByOwnerAndUserId(owner, userId)
	if getErr != nil {
		err := fmt.Errorf("get contact failed, owner = (%s) userId = (%s), error = %v", owner, userId, getErr)
		log.Log().Println(logs.Error(err).Extra(logs.F{"service", "ContactGetOneByOwnerAndUser"}).Trace())
		return nil, err
	}
	return &Contact{
		Id: row.Id,
		Owner:row.Owner,
		UserId:row.UserId,
		UserName:row.UserName,
		UserEmail:row.UserEmail,
	}, nil
}

func ContactGetOneByOwnerAndUserEmail(owner, email string) (*Contact, error) {
	if owner == "" || email == "" {
		err := fmt.Errorf("get contact failed, owner = (%s) email = (%s)", owner, email)
		log.Log().Println(logs.Error(err).Extra(logs.F{"service", "ContactGetOneByOwnerAndUser"}).Trace())
		return nil, err
	}
	row, getErr := data.ContactGetByOwnerAndUserEmail(owner, email)
	if getErr != nil {
		err := fmt.Errorf("get contact failed, owner = (%s) email = (%s), error = %v", owner, email, getErr)
		log.Log().Println(logs.Error(err).Extra(logs.F{"service", "ContactGetOneByOwnerAndUser"}).Trace())
		return nil, err
	}
	return &Contact{
		Id: row.Id,
		Owner:row.Owner,
		UserId:row.UserId,
		UserName:row.UserName,
		UserEmail:row.UserEmail,
	}, nil
}
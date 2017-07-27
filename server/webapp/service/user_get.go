package service

import (
	"fmt"
	"github.com/pharosnet/logs"
	"liulishuo/somechat/core/data"
	"liulishuo/somechat/log"
)

func UserGetByEmail(email string) (*User, error) {
	if email == "" {
		err := fmt.Errorf("get user failed, email is empty. email=%s",  email)
		log.Log().Println(logs.Error(err).Extra(logs.F{"service", "UserGetByEmail"}).Trace())
		return nil, err
	}
	userRow, userRowGetErr := data.UserGetByEmail(email)
	if userRowGetErr != nil {
		err := fmt.Errorf("get user failed, can not find user by email = %s, error = %v",  email, userRowGetErr)
		log.Log().Println(logs.Error(err).Extra(logs.F{"service", "UserGetByEmail"}).Trace())
		return nil, err
	}
	user := &User{
		Id:userRow.Id,
		Name:userRow.Name,
		Email:userRow.Email,
		Password:userRow.Password,
	}
	// TODO CACHE
	return user, nil
}

func UserGetById(id string) (*User, error) {
	if id == "" {
		err := fmt.Errorf("get user failed, id is empty. id=%s",  id)
		log.Log().Println(logs.Error(err).Extra(logs.F{"service", "UserGetById"}).Trace())
		return nil, err
	}
	userRow, userRowGetErr := data.UserGetById(id)
	if userRowGetErr != nil {
		err := fmt.Errorf("get user failed, can not find user by id = %s, error = %v",  id, userRowGetErr)
		log.Log().Println(logs.Error(err).Extra(logs.F{"service", "UserGetById"}).Trace())
		return nil, err
	}
	user := &User{
		Id:userRow.Id,
		Name:userRow.Name,
		Email:userRow.Email,
	}
	// TODO CACHE
	return user, nil
}
package service

import (
	"fmt"
	"github.com/pharosnet/logs"
	"liulishuo/somechat/core/data"
	"github.com/pharosnet/auid"
	"time"
	"liulishuo/somechat/log"
)

func UserAdd(name, email, password string) (*User, error) {
	if name == "" || email == "" || password == "" {
		err := fmt.Errorf("new user failed, args is bad. name=%s, email=%s, password=%s", name, email, password)
		log.Log().Println(logs.Error(err).Extra(logs.F{"service", "UserAdd"}).Trace())
		return nil, err
	}
	userRow := new(data.UserRow)
	userRow.Id = auid.NewAuid()
	userRow.Name = name
	userRow.Email = email
	userRow.Password = password
	userRow.CreateTime = time.Now()
	userRow.CreateUser = userRow.Id
	userRow.UpdateTime = userRow.CreateTime
	userRow.UpdateUser = userRow.CreateUser
	userRow.Version = int64(1)
	tx, txBegErr := data.DAL().BeginTx()
	if txBegErr != nil {
		err := fmt.Errorf("new user failed. tx begin failed, %v", txBegErr)
		log.Log().Println(logs.Error(err).Extra(logs.F{"service", "UserAdd"}).Trace())
		return nil, err
	}
	insertResult, insertErr := data.UserInsert(tx, userRow)
	if insertErr != nil || insertResult == int64(0) {
		err := fmt.Errorf("new user failed. table insert failed, row affected = %v, error = %v, tx roll back = %v", insertResult, insertErr, tx.Rollback())
		log.Log().Println(logs.Error(err).Extra(logs.F{"service", "UserAdd"}).Trace())
		return nil, err
	}
	if cmtErr := tx.Commit(); cmtErr != nil {
		err := fmt.Errorf("new user failed. tx commit failed, %v, tx roll back = %v", cmtErr, tx.Rollback())
		log.Log().Println(logs.Error(err).Extra(logs.F{"service", "UserAdd"}).Trace())
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

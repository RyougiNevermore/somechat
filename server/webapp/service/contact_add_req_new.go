package service

import (
	"fmt"
	"github.com/pharosnet/logs"
	"liulishuo/somechat/log"
	"liulishuo/somechat/core/data"
	"github.com/pharosnet/auid"
	"time"
	"liulishuo/somechat/core/common"
)

func ContactAddRequestNew(fromUserId, toUserId string) (*ContactAddRequest, error) {
	fromUser, fromUserGetErr := UserGetById(fromUserId)
	if fromUserGetErr != nil {
		err := fmt.Errorf("get from user failed. id = %s, error = %v", fromUserId, fromUserGetErr)
		log.Log().Println(logs.Error(err).Extra(logs.F{"service", "ContactAddRequestNew"}).Trace())
		return nil, err
	}
	toUser, toUserGetErr := UserGetById(toUserId)
	if toUserGetErr != nil {
		err := fmt.Errorf("get to user failed. id = %s, error = %v", toUserId, toUserGetErr)
		log.Log().Println(logs.Error(err).Extra(logs.F{"service", "ContactAddRequestNew"}).Trace())
		return nil, err
	}
	reqRow := &data.ContactAddRequestRow{
		Id:auid.NewAuid(),
		FromId:fromUser.Id,
		FromName:fromUser.Name,
		FromEmail:fromUser.Email,
		ToId:toUser.Id,
		ToName:toUser.Name,
		ToEmail:toUser.Email,
		CreateTime:time.Now(),
	}
	tx, txBegErr := data.DAL().BeginTx()
	if txBegErr != nil {
		err := fmt.Errorf("new contact add request failed. tx begin failed, %v", txBegErr)
		log.Log().Println(logs.Error(err).Extra(logs.F{"service", "ContactAddRequestNew"}).Trace())
		return nil, err
	}
	insertResult, insertErr := data.ContactAddRequestInsert(tx, reqRow)
	if insertErr != nil || insertResult == int64(0) {
		err := fmt.Errorf("new contact add request failed. table insert failed, row affected = %v, error = %v, tx roll back = %v", insertResult, insertErr, tx.Rollback())
		log.Log().Println(logs.Error(err).Extra(logs.F{"service", "ContactAddRequestNew"}).Trace())
		return nil, err
	}
	if cmtErr := tx.Commit(); cmtErr != nil {
		err := fmt.Errorf("new contact add request failed. tx commit failed, %v, tx roll back = %v", cmtErr, tx.Rollback())
		log.Log().Println(logs.Error(err).Extra(logs.F{"service", "ContactAddRequestNew"}).Trace())
		return nil, err
	}
	req := &ContactAddRequest{
		Id:reqRow.Id,
		FromId:reqRow.FromId,
		FromName:reqRow.FromName,
		FromEmail:reqRow.FromEmail,
		ToId:reqRow.ToId,
		ToName:reqRow.ToName,
		ToEmail:reqRow.ToEmail,
		CreateTime:common.JsonTimed(reqRow.CreateTime),
	}
	// TODO CACHE
	return req, nil
}


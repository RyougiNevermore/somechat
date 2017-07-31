package service

import (
	"fmt"
	"github.com/pharosnet/logs"
	"liulishuo/somechat/core/data"
	"github.com/pharosnet/auid"
	"time"
	"liulishuo/somechat/log"
)

func ContactAddRequestAccept(reqId, opUserId string) (*Contact, error) {
	if reqId == "" {
		err := fmt.Errorf("contact add req accept failed, reqid is empty, reqId = %s", reqId)
		log.Log().Println(logs.Error(err).Extra(logs.F{"service", "ContactAddRequestAccept"}).Trace())
		return nil, err
	}
	if opUserId == "" {
		err := fmt.Errorf("contact add req accept failed, userId is empty, userId = %s", opUserId)
		log.Log().Println(logs.Error(err).Extra(logs.F{"service", "ContactAddRequestAccept"}).Trace())
		return nil, err
	}
	reqRow, reqRowGetErr := data.ContactAddRequestGetById(reqId)
	if reqRowGetErr != nil {
		err := fmt.Errorf("contact add req accept failed, get req failed, id = %s, error = %v", reqId, reqRowGetErr)
		log.Log().Println(logs.Error(err).Extra(logs.F{"service", "ContactAddRequestAccept"}).Trace())
		return nil, err
	}
	magicCode := auid.NewAuid()
	contactToRow := &data.ContactRow{
		Id:auid.NewAuid(),
		Owner:reqRow.ToId,
		UserId:reqRow.FromId,
		UserEmail:reqRow.FromEmail,
		UserName:reqRow.FromName,
		CreateUser:opUserId,
		CreateTime:time.Now(),
		UpdateUser:opUserId,
		UpdateTime:time.Now(),
		Version:int64(1),
		MagicCode:magicCode,
	}
	contactFromRow := &data.ContactRow{
		Id:auid.NewAuid(),
		Owner:reqRow.FromId,
		UserId:reqRow.ToId,
		UserEmail:reqRow.ToEmail,
		UserName:reqRow.ToName,
		CreateUser:opUserId,
		CreateTime:time.Now(),
		UpdateUser:opUserId,
		UpdateTime:time.Now(),
		Version:int64(1),
		MagicCode:magicCode,
	}
	tx, txBegErr := data.DAL().BeginTx()
	if txBegErr != nil {
		err := fmt.Errorf("accept contact add request failed. tx begin failed, %v", txBegErr)
		log.Log().Println(logs.Error(err).Extra(logs.F{"service", "ContactAddRequestAccept"}).Trace())
		return nil, err
	}
	insertResult, insertErr := data.ContactInsert(tx, contactToRow, contactFromRow)
	if insertErr != nil || insertResult == int64(0) {
		err := fmt.Errorf("accept contact add request failed. table insert failed, row affected = %v, error = %v, tx roll back = %v", insertResult, insertErr, tx.Rollback())
		log.Log().Println(logs.Error(err).Extra(logs.F{"service", "ContactAddRequestAccept"}).Trace())
		return nil, err
	}
	// del req
	deleteResult, deleteErr := data.ContactAddRequestDelete(tx, reqRow)
	if deleteErr != nil || deleteResult == int64(0) {
		err := fmt.Errorf("accept contact add request failed. table delete failed, row affected = %v, error = %v, tx roll back = %v", deleteResult, deleteErr, tx.Rollback())
		log.Log().Println(logs.Error(err).Extra(logs.F{"service", "ContactAddRequestAccept"}).Trace())
		return nil, err
	}
	if cmtErr := tx.Commit(); cmtErr != nil {
		err := fmt.Errorf("accept contact add request failed. tx commit failed, %v, tx roll back = %v", cmtErr, tx.Rollback())
		log.Log().Println(logs.Error(err).Extra(logs.F{"service", "ContactAddRequestAccept"}).Trace())
		return nil, err
	}
	return &Contact{
		Id:contactToRow.Id,
		Owner:contactToRow.Owner,
		UserId:contactToRow.UserId,
		UserName:contactToRow.UserName,
		UserEmail:contactToRow.UserEmail,
	}, nil
}

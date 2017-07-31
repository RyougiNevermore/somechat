package service

import (
	"fmt"
	"github.com/pharosnet/logs"
	"liulishuo/somechat/core/data"
	"liulishuo/somechat/log"
)

func ContactAddRequestReject(reqId, opUserId string) (error) {
	if reqId == "" {
		err := fmt.Errorf("contact add req reject failed, reqid is empty, reqId = %s", reqId)
		log.Log().Println(logs.Error(err).Extra(logs.F{"service", "ContactAddRequestReject"}).Trace())
		return err
	}
	if opUserId == "" {
		err := fmt.Errorf("contact add req reject failed, userId is empty, userId = %s", opUserId)
		log.Log().Println(logs.Error(err).Extra(logs.F{"service", "ContactAddRequestReject"}).Trace())
		return err
	}
	reqRow, reqRowGetErr := data.ContactAddRequestGetById(reqId)
	if reqRowGetErr != nil {
		err := fmt.Errorf("contact add req reject failed, get req failed, id = %s, error = %v", reqId, reqRowGetErr)
		log.Log().Println(logs.Error(err).Extra(logs.F{"service", "ContactAddRequestReject"}).Trace())
		return err
	}
	if reqRow.ToId != opUserId {
		err := fmt.Errorf("contact add req reject failed, req to user(%v) is not op user(%v)", reqRow.FromId, opUserId)
		log.Log().Println(logs.Error(err).Extra(logs.F{"service", "ContactAddRequestReject"}).Trace())
		return err
	}
	tx, txBegErr := data.DAL().BeginTx()
	if txBegErr != nil {
		err := fmt.Errorf("reject contact add request failed. tx begin failed, %v", txBegErr)
		log.Log().Println(logs.Error(err).Extra(logs.F{"service", "ContactAddRequestReject"}).Trace())
		return err
	}
	// del req
	deleteResult, deleteErr := data.ContactAddRequestDelete(tx, reqRow)
	if deleteErr != nil || deleteResult == int64(0) {
		err := fmt.Errorf("reject contact add request failed. table delete failed, row affected = %v, error = %v, tx roll back = %v", deleteResult, deleteErr, tx.Rollback())
		log.Log().Println(logs.Error(err).Extra(logs.F{"service", "ContactAddRequestReject"}).Trace())
		return err
	}
	if cmtErr := tx.Commit(); cmtErr != nil {
		err := fmt.Errorf("accept contact add request failed. tx commit failed, %v, tx roll back = %v", cmtErr, tx.Rollback())
		log.Log().Println(logs.Error(err).Extra(logs.F{"service", "ContactAddRequestReject"}).Trace())
		return err
	}
	return nil
}

package data

import (
	"time"
	"github.com/niflheims-io/qb"
	"fmt"
	"github.com/pharosnet/logs"
	"liulishuo/somechat/log"
)

type ContactAddRequestRow struct {
	Id string 	`pk:"ID"`
	FromId string	`col:"FROM_ID"`
	FromName string	`col:"FROM_NAME"`
	FromEmail string	`col:"FROM_EMAIL"`
	ToId string	`col:"TO_ID"`
	ToName string	`col:"TO_NAME"`
	ToEmail string	`col:"TO_EMAIL"`
	CreateTime time.Time	`col:"CREATE_TIME"`
}

func (r ContactAddRequestRow) TableName() string {
	return "CONTACT_ADD_REQ"
}

func ContactAddRequestInsert(tx *qb.Tx, rows ...*ContactAddRequestRow) (int64, error) {
	if tx == nil {
		err := fmt.Errorf("contact add req insert failed, tx is nil, tx = %v", tx)
		log.Log().Println(logs.Error(err).Extra(logs.F{"sql", "CONTACT_ADD_REQ"}).Trace())
		return int64(0), err
	}
	affected := int64(0)
	for _, row := range rows {
		affectedOne, err := tx.Insert(row)
		if err != nil || affectedOne == int64(0) {
			err = fmt.Errorf("contact add req insert failed, affected=%d, error = %v", affectedOne, err)
			log.Log().Println(logs.Error(err).Extra(logs.F{"sql", "CONTACT_ADD_REQ"}).Trace())
			return int64(0), err
		}
		affected = affected + affectedOne
	}
	return affected, nil
}

func ContactAddRequestDelete(tx *qb.Tx, rows ...*ContactAddRequestRow) (int64, error) {
	if tx == nil {
		err := fmt.Errorf("contact add req delete failed, tx is nil, tx = %v", tx)
		log.Log().Println(logs.Error(err).Extra(logs.F{"sql", "CONTACT_ADD_REQ"}).Trace())
		return int64(0), err
	}
	affected := int64(0)
	for _, row := range rows {
		affectedOne, err := tx.Delete(row)
		if err != nil || affectedOne == int64(0) {
			err = fmt.Errorf("contact add req delete failed, affected=%d, error = %v", affectedOne, err)
			log.Log().Println(logs.Error(err).Extra(logs.F{"kind", "sql"}).Trace())
			return int64(0), err
		}
		affected = affected + affectedOne
	}
	return affected, nil
}

func ContactAddRequestGetById(id string) (*ContactAddRequestRow, error) {
	var one ContactAddRequestRow
	if err := DAL().Query(`SELECT * FROM "CONTACT_ADD_REQ" WHERE "ID" = $1`, &id).One(&one); err != nil {
		err := fmt.Errorf("contact add req get failed, can not find by id = %s, error = %v", id, err)
		log.Log().Println(logs.Error(err).Extra(logs.F{"sql", "CONTACT_ADD_REQ"}).Trace())
		return  nil, err
	}
	if one.Id == "" {
		err := fmt.Errorf("contact add req get failed, can not find by id = %s", id)
		log.Log().Println(logs.Error(err).Extra(logs.F{"sql", "CONTACT_ADD_REQ"}).Trace())
		return  nil, err
	}
	return &one, nil
}

func ContactAddRequestListByToUser(userId string) ([]ContactAddRequestRow, error) {
	var list []ContactAddRequestRow
	if err := DAL().Query(`SELECT * FROM "CONTACT_ADD_REQ" WHERE "TO_ID" = $1 ORDER BY "CREATE_TIME" DESC, "ID" DESC`, &userId).List(&list); err != nil {
		err := fmt.Errorf("contact add req by to user failed, can not find by userId = %s, error = %v", userId, err)
		log.Log().Println(logs.Error(err).Extra(logs.F{"sql", "CONTACT_ADD_REQ"}).Trace())
		return  nil, err
	}
	return list, nil
}

func ContactAddRequestListByFromUser(userId string) ([]ContactAddRequestRow, error) {
	var list []ContactAddRequestRow
	if err := DAL().Query(`SELECT * FROM "CONTACT_ADD_REQ" WHERE "FROM_ID" = $1 ORDER BY "CREATE_TIME" DESC, "ID" DESC`, &userId).List(&list); err != nil {
		err := fmt.Errorf("contact add req by from user failed, can not find by userId = %s, error = %v", userId, err)
		log.Log().Println(logs.Error(err).Extra(logs.F{"sql", "CONTACT_ADD_REQ"}).Trace())
		return  nil, err
	}
	return list, nil
}
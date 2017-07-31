package data

import (
	"time"
	"github.com/niflheims-io/qb"
	"fmt"
	"github.com/pharosnet/logs"
	"liulishuo/somechat/log"
)

type ContactRow struct {
	Id              string   `pk:"ID"`       //
	Owner           string   `col:"OWNER"`    //
	UserId   string   `col:"USER_ID"` //
	UserName   string   `col:"USER_NAME"` //
	UserEmail string `col:"USER_EMAIL"`
	// common
	CreateUser string    `col:"CREATE_USER"`
	CreateTime time.Time `col:"CREATE_TIME"`
	UpdateUser string    `col:"UPDATE_USER"`
	UpdateTime time.Time `col:"UPDATE_TIME"`
	Version    int64     `version:"VERSION"`
	//
	MagicCode string `col:"MAGIC_CODE"`
}

func (u ContactRow) TableName() string {
	return "CONTACT"
}

func ContactInsert(tx *qb.Tx, rows ...*ContactRow) (int64, error) {
	if tx == nil {
		err := fmt.Errorf("contact insert failed, tx is nil, tx = %v", tx)
		log.Log().Println(logs.Error(err).Extra(logs.F{"sql", "CONTACT"}).Trace())
		return int64(0), err
	}
	affected := int64(0)
	for _, row := range rows {
		affectedOne, err := tx.Insert(row)
		if err != nil || affectedOne == int64(0) {
			err = fmt.Errorf("contact insert failed, affected=%d, error = %v", affectedOne, err)
			log.Log().Println(logs.Error(err).Extra(logs.F{"sql", "CONTACT"}).Trace())
			return int64(0), err
		}
		affected = affected + affectedOne
	}
	return affected, nil
}

func ContactDelete(tx *qb.Tx, rows ...*ContactRow) (int64, error) {
	if tx == nil {
		err := fmt.Errorf("contact delete failed, tx is nil, tx = %v", tx)
		log.Log().Println(logs.Error(err).Extra(logs.F{"sql", "CONTACT"}).Trace())
		return int64(0), err
	}
	affected := int64(0)
	for _, row := range rows {
		affectedOne, err := tx.Delete(row)
		if err != nil || affectedOne == int64(0) {
			err = fmt.Errorf("contact delete failed, affected=%d, error = %v", affectedOne, err)
			log.Log().Println(logs.Error(err).Extra(logs.F{"sql", "CONTACT"}).Trace())
			return int64(0), err
		}
		affected = affected + affectedOne
	}
	return affected, nil
}

func ContactListByUserId(userId string) ([]ContactRow, error) {
	var list []ContactRow
	if err := DAL().Query(`SELECT * FROM "CONTACT" WHERE "OWNER" = $1 ORDER BY "USER_NAME" ASC`, &userId).List(&list); err != nil {
		err := fmt.Errorf("contact get failed, can not find by user id = %s, error = %v", userId, err)
		log.Log().Println(logs.Error(err).Extra(logs.F{"sql", "CONTACT"}).Trace())
		return  nil, err
	}
	return list, nil
}

func ContactGetById(id string) (*ContactRow, error) {
	var one ContactRow
	if err := DAL().Query(`SELECT * FROM "CONTACT" WHERE "ID" = $1`, &id).One(&one); err != nil {
		err := fmt.Errorf("contact get failed, can not find by id = %s, error = %v", id, err)
		log.Log().Println(logs.Error(err).Extra(logs.F{"sql", "CONTACT"}).Trace())
		return  nil, err
	}
	if one.Id == "" {
		err := fmt.Errorf("contact get failed, can not find by id = %s", id)
		log.Log().Println(logs.Error(err).Extra(logs.F{"sql", "CONTACT"}).Trace())
		return  nil, err
	}
	return &one, nil
}

func ContactGetByOwnerAndUserId(ownerId string, userId string) (*ContactRow, error) {
	var one ContactRow
	if err := DAL().Query(`SELECT * FROM "CONTACT" WHERE "OWNER" = $1 AND "USER_ID" = $2`, &ownerId, &userId).One(&one); err != nil {
		err := fmt.Errorf("contact get failed, can not find by owner(%s) and user id(%s), error = %v", ownerId, userId, err)
		log.Log().Println(logs.Error(err).Extra(logs.F{"sql", "CONTACT"}).Trace())
		return  nil, err
	}
	return &one, nil
}

func ContactGetByOwnerAndUserEmail(ownerId string, email string) (*ContactRow, error) {
	var one ContactRow
	if err := DAL().Query(`SELECT * FROM "CONTACT" WHERE "OWNER" = $1 AND "USER_EMAIL" = $2`, &ownerId, &email).One(&one); err != nil {
		err := fmt.Errorf("contact get failed, can not find by owner(%s) and user email(%s), error = %v", ownerId, email, err)
		log.Log().Println(logs.Error(err).Extra(logs.F{"sql", "CONTACT"}).Trace())
		return  nil, err
	}
	return &one, nil
}
package data

import (
	"time"
	"github.com/niflheims-io/qb"
	"fmt"
	"liulishuo/somechat/log"
	"github.com/pharosnet/logs"
)


type UserRow struct {
	Id       string `pk:"ID"`     // ID
	Name     string `col:"NAME"`  // NAME
	Email    string `col:"EMAIL"` // MOBILE
	Password string `col:"PASSWORD"`
	// common
	CreateUser string    `col:"CREATE_USER"`
	CreateTime time.Time `col:"CREATE_TIME"`
	UpdateUser string    `col:"UPDATE_USER"`
	UpdateTime time.Time `col:"UPDATE_TIME"`
	DeleteUser string    `col:"DELETE_USER"`
	DeleteTime time.Time `col:"DELETE_TIME"`
	Version    int64     `version:"VERSION"`
}

func (u UserRow) TableName() string {
	return "USER"
}

type UserCache struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
}

func NewUserCache(row *UserRow) UserCache {
	return UserCache{
		Id:row.Id,
		Name:row.Name,
		Email:row.Email,
	}
}

// user insert
func UserInsert(tx *qb.Tx, rows ...*UserRow) (int64, error) {
	if tx == nil {
		err := fmt.Errorf("user save failed, tx is nil, tx = %v", tx)
		log.Log().Println(logs.Error(err).Extra(logs.F{"sql", "USER"}).Trace())
		return int64(0), err
	}
	affected := int64(0)
	for _, row := range rows {
		affectedOne, err := tx.Insert(row)
		if err != nil || affectedOne == int64(0) {
			err = fmt.Errorf("user save failed, affected=%d, error = %v", affectedOne, err)
			log.Log().Println(logs.Error(err).Extra(logs.F{"sql", "USER"}).Trace())
			return int64(0), err
		}
		affected = affected + affectedOne
	}
	return affected, nil
}

func UserGetByEmail(email string) (*UserRow, error) {
	var user UserRow
	if err := DAL().Query(`SELECT * FROM "USER" WHERE "EMAIL" = $1 AND "DELETE_USER" = '' `, &email).One(&user); err != nil {
		err := fmt.Errorf("user get failed, email = %s, error = %v", email, err)
		log.Log().Println(logs.Error(err).Extra(logs.F{"sql", "USER"}).Trace())
		return  nil, err
	}
	if user.Id == "" {
		err := fmt.Errorf("user get failed, can not find user by email = %s.", email)
		log.Log().Println(logs.Error(err).Extra(logs.F{"sql", "USER"}).Trace())
		return  nil, err
	}
	return &user, nil
}

func UserGetById(id string) (*UserRow, error) {
	var user UserRow
	if err := DAL().Query(`SELECT * FROM "USER" WHERE "ID" = $1 AND "DELETE_USER" = '' `, &id).One(&user); err != nil {
		err := fmt.Errorf("user get failed, id = %s, error = %v", id, err)
		log.Log().Println(logs.Error(err).Extra(logs.F{"sql", "USER"}).Trace())
		return  nil, err
	}
	if user.Id == "" {
		err := fmt.Errorf("user get failed, can not find user by id = %s.", id)
		log.Log().Println(logs.Error(err).Extra(logs.F{"sql", "USER"}).Trace())
		return  nil, err
	}
	return &user, nil
}

func UserGetByName(name string) (*UserRow, error) {
	var user UserRow
	if err := DAL().Query(`SELECT * FROM "USER" WHERE "NAME" = $1 AND "DELETE_USER" = '' `, &name).One(&user); err != nil {
		err := fmt.Errorf("user get failed, name = %s, error = %v", name, err)
		log.Log().Println(logs.Error(err).Extra(logs.F{"sql", "USER"}).Trace())
		return  nil, err
	}
	if user.Id == "" {
		err := fmt.Errorf("user get failed, can not find user by name = %s.", name)
		log.Log().Println(logs.Error(err).Extra(logs.F{"sql", "USER"}).Trace())
		return  nil, err
	}
	return &user, nil
}
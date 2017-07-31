package service

import (
	"fmt"
	"liulishuo/somechat/log"
	"github.com/pharosnet/logs"
	"liulishuo/somechat/core/data"
	"liulishuo/somechat/core/common"
)

func ContactAddRequestListByToUser(userId string) ([]ContactAddRequest, error) {
	if userId == "" {
		err := fmt.Errorf("user id is empty, userId = (%s)", userId)
		log.Log().Println(logs.Error(err).Extra(logs.F{"service", "ContactAddRequestListByUserId"}).Trace())
		return nil, err
	}
	rows, rowsErr := data.ContactAddRequestListByToUser(userId)
	if rowsErr != nil {
		err := fmt.Errorf("list contact add req failed, userId = %s, error = %v", userId, rowsErr)
		log.Log().Println(logs.Error(err).Extra(logs.F{"service", "ContactAddRequestListByUserId"}).Trace())
		return nil, err
	}
	results := make([]ContactAddRequest, 0, len(rows))
	for _, row := range rows {
		results = append(results, ContactAddRequest{
			Id: row.Id,
			FromId:row.FromId,
			FromName:row.FromName,
			FromEmail:row.FromEmail,
			ToName:row.ToName,
			ToEmail:row.ToEmail,
			ToId:row.ToId,
			CreateTime:common.JsonTimed(row.CreateTime),
		})
	}
	return results, nil
}

func ContactAddRequestListByFromUser(userId string) ([]ContactAddRequest, error) {
	if userId == "" {
		err := fmt.Errorf("user id is empty, userId = (%s)", userId)
		log.Log().Println(logs.Error(err).Extra(logs.F{"service", "ContactAddRequestListByUserId"}).Trace())
		return nil, err
	}
	rows, rowsErr := data.ContactAddRequestListByFromUser(userId)
	if rowsErr != nil {
		err := fmt.Errorf("list contact add req failed, userId = %s, error = %v", userId, rowsErr)
		log.Log().Println(logs.Error(err).Extra(logs.F{"service", "ContactAddRequestListByUserId"}).Trace())
		return nil, err
	}
	results := make([]ContactAddRequest, 0, len(rows))
	for _, row := range rows {
		results = append(results, ContactAddRequest{
			Id: row.Id,
			FromId:row.FromId,
			FromName:row.FromName,
			FromEmail:row.FromEmail,
			ToName:row.ToName,
			ToEmail:row.ToEmail,
			ToId:row.ToId,
			CreateTime:common.JsonTimed(row.CreateTime),
		})
	}
	return results, nil
}

func ContactAddRequestGetById(id string) (*ContactAddRequest, error) {
	row, rowErr := data.ContactAddRequestGetById(id)
	if rowErr != nil {
		err := fmt.Errorf("get contact add req failed, id = %s, error = %v", id, rowErr)
		log.Log().Println(logs.Error(err).Extra(logs.F{"service", "ContactAddRequestGetById"}).Trace())
		return nil, err
	}
	return &ContactAddRequest{
		Id: row.Id,
		FromId:row.FromId,
		FromName:row.FromName,
		FromEmail:row.FromEmail,
		ToName:row.ToName,
		ToEmail:row.ToEmail,
		ToId:row.ToId,
		CreateTime:common.JsonTimed(row.CreateTime),
	}, nil
}
package model

import (
	"github.com/m0cchi/gfalcon"
)

const SQL_GET_SERVICE_BY_ID = "SELECT `iid`, `id` FROM `services` WHERE `id` = :service_id"

type Service struct {
	IID uint32 `db:"iid"`
	ID  string `db:"id"`
}

func GetService(db gfsql.DB, serviceID string) (*Service, error) {
	stmt, err := db.PrepareNamed(SQL_GET_SERVICE_BY_ID)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	service := &Service{}
	args := map[string]interface{}{"service_id": serviceID}
	err = stmt.Get(service, args)

	return service, err
}

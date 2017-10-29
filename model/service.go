package model

import (
	"github.com/go-sql-driver/mysql"
	"github.com/m0cchi/gfalcon"
)

const SQL_GET_SERVICE_BY_ID = "SELECT `iid`, `id` FROM `services` WHERE `id` = :service_id"
const SQL_CREATE_SERVICE = "INSERT INTO `services` (`id`) VALUE (:service_id)"

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

func CreateService(db gfsql.DB, serviceID string) (*Service, error) {
	stmt, err := db.PrepareNamed(SQL_CREATE_SERVICE)

	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	args := map[string]interface{}{"service_id": serviceID}
	result, err := stmt.Exec(args)
	service := &Service{
		IID: 0,
		ID:  serviceID}

	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == gfsql.ERR_CODE_DUPLICATE_ENTRY {
				return nil, ErrDuplicate
			}
		}
		return nil, err
	}

	if c, err := result.LastInsertId(); err == nil {
		service.IID = uint32(c)
	}

	return service, err
}

package model

import (
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/m0cchi/gfalcon"
)

const SQL_GET_SERVICE_BY_ID = "SELECT `iid`, `id` FROM `services` WHERE `id` = :service_id"
const SQL_CREATE_SERVICE = "INSERT INTO `services` (`id`) VALUE (:service_id)"

const SQL_DELETE_SERVICE_BY_IID = "DELETE FROM `services` WHERE `iid` = :service_iid"
const SQL_DELETE_SERVICE_BY_ID = "DELETE FROM `services` WHERE `id` = :service_id"

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

func DeleteServiceByIID(db gfsql.DB, serviceIID uint32) error {
	stmt, err := db.PrepareNamed(SQL_DELETE_SERVICE_BY_IID)
	if err != nil {
		return err
	}
	defer stmt.Close()

	args := map[string]interface{}{"service_iid": serviceIID}
	_, err = stmt.Exec(args)
	return err
}

func DeleteServiceByID(db gfsql.DB, serviceID string) error {
	stmt, err := db.PrepareNamed(SQL_DELETE_SERVICE_BY_ID)
	if err != nil {
		return err
	}
	defer stmt.Close()

	args := map[string]interface{}{"service_id": serviceID}
	_, err = stmt.Exec(args)
	return err
}

func (service *Service) Delete(db gfsql.DB) error {
	if service == nil {
		return errors.New("not specify service")
	}

	if service.IID != 0 {
		return DeleteServiceByIID(db, service.IID)
	} else if service.ID != "" {
		return DeleteServiceByID(db, service.ID)
	}

	return errors.New("not specify service")
}

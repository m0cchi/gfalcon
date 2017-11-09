package model

import (
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/m0cchi/gfalcon"
)

const SqlGetServiceByID = "SELECT `iid`, `id` FROM `services` WHERE `id` = :service_id"
const SqlCreateService = "INSERT INTO `services` (`id`) VALUE (:service_id)"

const SqlDeleteServiceByIID = "DELETE FROM `services` WHERE `iid` = :service_iid"
const SqlDeleteServiceByID = "DELETE FROM `services` WHERE `id` = :service_id"

type Service struct {
	IID uint32 `db:"iid"`
	ID  string `db:"id"`
}

func GetService(db gfsql.DB, serviceID string) (*Service, error) {
	stmt, err := db.PrepareNamed(SqlGetServiceByID)
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
	stmt, err := db.PrepareNamed(SqlCreateService)

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
			if mysqlErr.Number == gfsql.ErrCodeDuplicateEntry {
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
	stmt, err := db.PrepareNamed(SqlDeleteServiceByIID)
	if err != nil {
		return err
	}
	defer stmt.Close()

	args := map[string]interface{}{"service_iid": serviceIID}
	_, err = stmt.Exec(args)
	return err
}

func DeleteServiceByID(db gfsql.DB, serviceID string) error {
	stmt, err := db.PrepareNamed(SqlDeleteServiceByID)
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

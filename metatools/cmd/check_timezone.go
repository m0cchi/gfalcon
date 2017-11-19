// +build check_timezone.go

package main

import (
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"os"
	"time"
)

type Result struct {
	T time.Time `db:"t"`
}

func checkTimezone(source string) error {
	db, err := sqlx.Connect("mysql", source)
	if err != nil {
		return err
	}

	sql := `SELECT CURRENT_TIMESTAMP as t FROM dual`
	r := &Result{}
	err = db.Get(r, sql)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	if now.Day() == r.T.UTC().Day() && now.Hour() == r.T.UTC().Hour() {
		fmt.Println("ok")
		return nil
	}

	sub := now.Sub(r.T)

	return errors.New(fmt.Sprintf("diff: %v", sub))
}

func main() {
	if len(os.Args) > 1 {
		err := checkTimezone(os.Args[1])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		fmt.Println(`e.g.
go run check_timezone.go 'gfadmin:gfadmin@unix(/tmp/mysql.sock)/gfalcon?parseTime=true&loc=Asia%2FTokyo'`)
		os.Exit(1)
	}
	os.Exit(0)
}

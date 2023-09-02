package banco

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql" // implicitly imported
)

func Conectar() (*sql.DB, error) {
	stringConnection := "golang:go231!@/findfy?charset=utf8&parseTime=True&loc=Local"

	db, err := sql.Open("mysql", stringConnection)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

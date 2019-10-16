package db

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"sync"
)

var (
	persistentDB     *sql.DB
	persistentDBOnce sync.Once
)

type Config interface {
	DbServer() string
	DbUsername() string
	DbPassword() string
	DbName() string
}

func newDb(cfg Config) *sql.DB {
	dataSource := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true",
		cfg.DbUsername(), cfg.DbPassword(), cfg.DbServer(), cfg.DbName())
	db, err := sql.Open("mysql", dataSource)
	if err != nil {
		panic(err)
	}
	if err := db.Ping(); err != nil {
		panic(err)
	}
	return db
}

func InitConnection(cfg Config) {
	persistentDBOnce.Do(func() {
		persistentDB = newDb(cfg)
	})
}

func GetPersistentDB() *sql.DB {
	return persistentDB
}


type Execer interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type ExecerContext interface {
	Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

type Queryer interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

type QueryerContext interface {
	Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

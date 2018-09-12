package mysql

import (
	"database/sql"
	"sync"
	"context"
	"os"
	_ "github.com/go-sql-driver/mysql"
)

var (
	persistentDB     *sql.DB
	persistentDBOnce sync.Once
)

func NewDb() *sql.DB {
	db, _ := sql.Open("mysql", os.Getenv("DATA_SOURCE"))
	if err := db.Ping(); err != nil {
		panic(err)
	}
	return db
}

func GetPersistentDB() *sql.DB {
	persistentDBOnce.Do(func() {
		persistentDB = NewDb()
	})
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

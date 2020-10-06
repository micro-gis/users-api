package users_db

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/micro-gis/utils/logger"
	"log"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

const (
	mysql_users_username = "mysql_users_username"
	mysql_users_pw       = "mysql_users_pw"
	mysql_users_host     = "mysql_users_host"
	mysql_users_schema   = "mysql_users_schema"
)

var (
	Client   *sql.DB
	username = os.Getenv(mysql_users_username)
	password = os.Getenv(mysql_users_pw)
	host     = os.Getenv(mysql_users_host)
	schema   = os.Getenv(mysql_users_schema)
)

func init() {
	datasourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8",
		username,
		password,
		host,
		schema,
	)
	var connErr error
	Client, connErr = sql.Open("mysql", datasourceName)
	if connErr != nil {
		panic(connErr)
	}
	if connErr = Client.Ping(); connErr != nil {
		panic(connErr)
	}
	mysql.SetLogger(logger.GetLogger())
	log.Println("database successfully configured")

}

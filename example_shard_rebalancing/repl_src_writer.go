package main

import (
    "database/sql"
    "log"
    "fmt"
)

const SQL_DROP_DATABASE = `DROP DATABASE IF EXISTS adt_test`
const SQL_CREATE_DATABASE = `CREATE DATABASE adt_test`
const SQL_RESET_MASTER = `RESET MASTER`

const SQL_CREATE_TABLE_TEST1 = `
CREATE TABLE adt_test.test_1 (
  no int(11) NOT NULL,
  seq int(11) NOT NULL,
  uk int(11) NOT NULL,
  v text NOT NULL,
  c int(11) NOT NULL,
  modtime timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  regtime datetime NOT NULL,
  PRIMARY KEY (no),
  UNIQUE KEY ux_uk_no (uk),
  KEY ix_modtime (modtime),
  KEY ix_regtime (regtime)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
`

//const SQL_SHOW_DATABASES = `show databases;`
//const SQL_SHOW_TABLES = `show tables from adt_test`

type ReplSrcWriter struct {

}

func (w *ReplSrcWriter) start(dbUser, dbPass, dbHost string, dbPort int) {

    connUrl := fmt.Sprintf("%s:%s@tcp(%s:%d)/", dbUser, dbPass, dbHost, dbPort)
    log.Println(connUrl)

    db, err := sql.Open("mysql", connUrl)
    if err != nil {
        log.Fatal(err)
        panic(err.Error())
    }
    defer db.Close()

    execSql := func(sql string) {
        result, err := db.Exec(sql)
        log.Println(result.LastInsertId())
        log.Println(result.RowsAffected())
        log.Println(err)
        if err != nil {
            log.Fatal(err)
            panic(err.Error())
        }
    }

    execSql(SQL_DROP_DATABASE)
    execSql(SQL_CREATE_DATABASE)
    execSql(SQL_RESET_MASTER)
    execSql(SQL_CREATE_TABLE_TEST1)

}




package main

import (
    "database/sql"
    "log"
    "fmt"
)

type ReplSrcWriter struct {
    
}

func NewReplSrcWriter(dbUser, dbPass, dbHost string, dbPort int) *ReplSrcWriter {

    writer := &ReplSrcWriter{}


}

func (w *ReplSrcWriter) start() {

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

    execSql("SELECT 1")

}




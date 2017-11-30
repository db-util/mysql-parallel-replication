package main

import (
    "fmt"
    "os"
    "strconv"
    "time"
)

func main() {

    fmt.Println("Let's start sharding test.")

    for _, env := range os.Environ() {
        log.Println(env)
    }

    dbUser := os.Getenv("MYSQL_USER")
    dbPass := os.Getenv("MYSQL_PASSWORD")
    dbHost := os.Getenv("MYSQL_HOST")
    dbPort, err := strconv.Atoi(os.Getenv("MYSQL_PORT"))

    if err != nil {
        panic("Illegal MYSQL_PORT ")
    }

    conn := NewTestDBConn(dbUser, dbPass, dbHost, dbPort, 1000, 1000)
    defer conn.Close()
    conn.ResetMaster()

    //ch := dmlConns[0].Table1_Insert()
    stopChans := make([]chan bool, 0, 4)
    stopChans = append(stopChans, conn.Table1_InsertOnDupKeyUpdate())
    stopChans = append(stopChans, conn.Table1_Update())
    stopChans = append(stopChans, conn.Table1_Delete())
    stopChans = append(stopChans, conn.Table1_Replace())

    log.Println("DML start")
    time.Sleep(60 * time.Second)

    for _, ch := range stopChans {
        ch <- true
    }
    //
    //log.Println("DML stop")
    //
    //time.Sleep(5 * time.Second)



}

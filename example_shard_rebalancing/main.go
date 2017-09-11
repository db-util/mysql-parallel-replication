package main

import (
    "fmt"
    "os"
    "strconv"
    "log"
)

func main() {
    fmt.Println("Hello World")

    for _, env := range os.Environ() {
        log.Println(env)
    }

    dbUser := "root"
    dbPass := os.Getenv("MYSQL_ROOT_PASSWORD")
    dbHost := "127.0.0.1"
    dbPortStr := os.Getenv("REPL_SRC_PORT")
    dbPort, err := strconv.Atoi(dbPortStr)

    if err != nil {
        panic("Illegal REPL_SRC_PORT " + dbPortStr)
    }

    replSrc := &ReplSrcWriter{}
    replSrc.start(dbUser, dbPass, dbHost, dbPort);

}
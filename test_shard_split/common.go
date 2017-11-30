package main

import (
    "os"
    "runtime"
    _log "log"
    _rand "math/rand"
    "time"
    "encoding/json"
)

var log = _log.New(os.Stdout, "[TEST_SHARD] ", _log.LstdFlags | _log.Lshortfile)
var rand = _rand.New(_rand.NewSource(time.Now().UnixNano()))

func panicIfError(err error) {
    if err != nil {
        log.Println(err)
        panic(err)
    }
}

type _Trace struct {
    file string
    line int
    fn string
}

func Trace() *_Trace {
    pc := make([]uintptr, 10)  // at least 1 entry needed
    runtime.Callers(2, pc)
    fn := runtime.FuncForPC(pc[0])
    file, line := fn.FileLine(pc[0])
    //fmt.Printf("%s:%d %s\n", file, line, fn.Name())
    return &_Trace{file, line, fn.Name()}
}

func (t *_Trace) Fn() string {
    return t.fn
}

func prettify(a interface{}) string {
    bArray, err := json.Marshal(a)
    panicIfError(err)
    return string(bArray)
}
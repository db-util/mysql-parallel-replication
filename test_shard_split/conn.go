package main

import (
    "database/sql"
    "fmt"
    "time"
)

type TestDBConn struct {
    *sql.DB

    RangeNo int
    RangeUk int
    running bool
}

func NewTestDBConn(dbUser, dbPass, dbHost string, dbPort int, rangeNo, rangeUk int) *TestDBConn {

    connUrl := fmt.Sprintf("%s:%s@tcp(%s:%d)/", dbUser, dbPass, dbHost, dbPort)
    log.Println(connUrl)

    db, err := sql.Open("mysql", connUrl)
    panicIfError(err)
    db.SetMaxOpenConns(100)
    db.SetMaxIdleConns(100)

    result := &TestDBConn{db, rangeNo, rangeUk, false}

    row := result.QueryRow("SELECT 'hello world'")
    var selectHello string
    row.Scan(&selectHello)
    log.Println(selectHello)

    return result

}

func (c *TestDBConn) ResetMaster() {
    var err error
    _, err = c.Exec("RESET MASTER")
    panicIfError(err)

    _, err = c.Exec("TRUNCATE TABLE test_shard_src.t1")
    panicIfError(err)

    _, err = c.Exec("TRUNCATE TABLE test_shard_dst1.t1")
    panicIfError(err)

    _, err = c.Exec("TRUNCATE TABLE test_shard_dst2.t1")
    panicIfError(err)
}

/*func (c *TestDBConn) Table1_Insert() chan bool {

    stmt, err := c.Prepare("INSERT INTO test_shard_src.t1 SET " +
            "no = ?, " +
            "seq = ?, " +
            "uk = ?, " +
            "v = '012ABCabc가나다', " +
            "c = 0, " +
            "modtime = CURRENT_TIMESTAMP, " +
            "regtime = CURRENT_TIMESTAMP ;")
    panicIfError(err)

    stopCh := make(chan bool, 1)
    go func() {
        for {
            select {
            case <-stopCh: return
            default:
            }

            _, err := stmt.Exec(rand.Intn(c.RangeNo), 0, rand.Intn(c.RangeUk))

            if err != nil {
                log.Println(Trace().Fn(), "fail", err)
            } else {
                log.Println(Trace().Fn(), "success")
            }

        }
    }()

    return stopCh

}*/

func (c *TestDBConn) Table1_InsertOnDupKeyUpdate() chan bool {

    stmt := c.Prepare("INSERT INTO test_shard_src.t1 SET " +
            "no = ?, " +
            "seq = ?, " +
            "uk = ?, " +
            "v = '012ABCabc가나다', " +
            "c = 0, " +
            "regtime = CURRENT_TIMESTAMP " +
        "ON DUPLICATE KEY UPDATE " +
            "c = c + 1, " +
            "modtime = CURRENT_TIMESTAMP ;")

    caller := Trace().Fn()
    stopCh := make(chan bool, 1)

    go func() {
        for c.runnable(stopCh) {
            c.ExecStmt(caller, stmt, rand.Intn(c.RangeNo), 0, rand.Intn(c.RangeUk))
        }
    }()

    return stopCh
}

func (c *TestDBConn) Table1_Update() chan bool {

    stmt := c.Prepare("UPDATE test_shard_src.t1 SET " +
        "uk = ?, " +
        "c = c + 1, " +
        "modtime = CURRENT_TIMESTAMP " +
        "WHERE " +
        "no = ?")

    caller := Trace().Fn()
    stopCh := make(chan bool, 1)

    go func(){
        for c.runnable(stopCh) {
            c.ExecStmt(caller, stmt, rand.Intn(c.RangeUk), rand.Intn(c.RangeNo))
        }
    }()

    return stopCh
}

func (c *TestDBConn) Table1_Replace() chan bool {

    stmt := c.Prepare("REPLACE INTO test_shard_src.t1 SET " +
        "no = ?, " +
        "seq = ?, " +
        "uk = ?, " +
        "v = '012ABCabc가나다', " +
        "c = 0, " +
        "modtime = CURRENT_TIMESTAMP, " +
        "regtime = CURRENT_TIMESTAMP ;")

    caller := Trace().Fn()
    stopCh := make(chan bool, 1)

    go func() {
        for c.runnable(stopCh) {
            c.ExecStmt(caller, stmt, rand.Intn(c.RangeNo), 0, rand.Intn(c.RangeUk))
        }
    }()

    return stopCh
}

func (c *TestDBConn) Table1_Delete() chan bool {

    stmt := c.Prepare("DELETE FROM test_shard_src.t1 WHERE no = ?")

    caller := Trace().Fn()
    stopCh := make(chan bool, 1)

    go func() {
        for c.runnable(stopCh) {
            c.ExecStmt(caller, stmt, rand.Intn(c.RangeNo))
        }
    }()

    return stopCh

}

func (c *TestDBConn) Prepare(query string) *sql.Stmt {
    stmt, err := c.DB.Prepare(query)
    panicIfError(err)
    return stmt
}

func (c *TestDBConn) runnable(stopCh chan bool) bool {

    select {
    case <-stopCh:
        return false
    default:
        return true
    }
}

func (c *TestDBConn) ExecStmt(caller string, stmt *sql.Stmt, args ...interface{}) {
    _, err := stmt.Exec(args ...)

    time.Sleep(10 *time.Millisecond)

    if err != nil {
        log.Println(caller, "fail", err)
    } else {
        log.Println(caller, "success")
    }
}
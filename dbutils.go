package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const driverName string = "sqlite3"
const dataSourceName string = "./widmgo.sqlite"

func example() {
	/*
				func main() {
		        db, err := sql.Open("sqlite3", "./foo.db")
		        checkErr(err)

		        // insert
		        stmt, err := db.Prepare("INSERT INTO userinfo(username, departname, created) values(?,?,?)")
		        checkErr(err)

		        res, err := stmt.Exec("astaxie", "研发部门", "2012-12-09")
		        checkErr(err)

		        id, err := res.LastInsertId()
		        checkErr(err)

		        fmt.Println(id)
		        // update
		        stmt, err = db.Prepare("update userinfo set username=? where uid=?")
		        checkErr(err)

		        res, err = stmt.Exec("astaxieupdate", id)
		        checkErr(err)

		        affect, err := res.RowsAffected()
		        checkErr(err)

		        fmt.Println(affect)

		        // query
		        rows, err := db.Query("SELECT * FROM userinfo")
		        checkErr(err)
		        var uid int
		        var username string
		        var department string
		        var created time.Time

		        for rows.Next() {
		            err = rows.Scan(&uid, &username, &department, &created)
		            checkErr(err)
		            fmt.Println(uid)
		            fmt.Println(username)
		            fmt.Println(department)
		            fmt.Println(created)
		        }

		        rows.Close() //good habit to close

		        // delete
		        stmt, err = db.Prepare("delete from userinfo where uid=?")
		        checkErr(err)

		        res, err = stmt.Exec(id)
		        checkErr(err)

		        affect, err = res.RowsAffected()
		        checkErr(err)

		        fmt.Println(affect)

		        db.Close()

		    }

		    func checkErr(err error) {
		        if err != nil {
		            panic(err)
		        }
		    }
	*/
}

func transactionExample() {
	/*
			 trashSQL, err := database.Prepare("update task set is_deleted='Y',last_modified_at=datetime() where id=?")
		    if err != nil {
		        fmt.Println(err)
		    }
		    tx, err := database.Begin()
		    if err != nil {
		        fmt.Println(err)
		    }
		    _, err = tx.Stmt(trashSQL).Exec(id)
		    if err != nil {
		        fmt.Println("doing rollback")
		        tx.Rollback()
		    } else {
		        tx.Commit()
		    }
	*/
}

func insertRecs(recs []Rec) (int, error) {
	// connect to database
	conn, err := sql.Open(driverName, dataSourceName)
	defer conn.Close()

	if err != nil {
		return 0, err
	}

	// prepare insert statement
	stmt, err := conn.Prepare("INSERT INTO recs(id, linija, kelias, km, pk, m, siule, skodas, suvirino, operatorius, aparatas, tdata, kelintas, input) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		return 0, err
	}

	// start transaction
	tx, err := conn.Begin()
	if err != nil {
		return 0, err
	}

	// try to perform all inserts
	c := 0
	err = nil
	for _, rec := range recs {
		var ps []interface{}

		if rec.Id == 0 {
			ps[0] = nil
		} else {
			ps[0] = rec.Id
		}

		ps[1] = rec.Linija
		ps[2] = rec.Kelias
		ps[3] = rec.Km
		ps[4] = rec.Pk
		ps[5] = rec.M

		if rec.Siule.Valid {
			ps[6] = rec.Siule.Int
		} else {
			ps[6] = nil
		}

		ps[7] = rec.Skodas
		ps[8] = rec.Suvirino
		ps[9] = rec.Operatorius
		ps[10] = rec.Aparatas
		ps[11] = rec.TData.Format("2006-01-02")
		ps[12] = rec.Kelintas
		ps[13] = time.Now().Format("2006-01-02T15:04:05")

		// excecuting the statement
		_, err = tx.Stmt(stmt).Exec(ps)

		// if not success - rolling back transaction
		if err != nil {
			fmt.Println("doing rollback")
			rbErr := tx.Rollback()

			// if rollback fails
			if rbErr != nil {
				return 0, failRollback(recs[:c], rbErr)
			}

			// if rollback success
			return 0, err
		}

		// incrementing the counter
		c++
	}
	
	err = tx.Commit()
	if err != nil {
		return 0, failCommit(recs, err)
	}

	return c, nil
}

func failCommit([]Rec recs, err error) error {

}

func failRollback([]Rec recs, err error) error {
	//
}

func fetchRecs(sqlWhere string, sqlParamValue int) ([]Rec, error) {

	sqlProjection := "SELECT id, vietoskodas, defkodas, pavojingumas, dataaptikta, dataiki, pastaba, dataatlikta, budas"
	sqlFrom := " FROM Defektai INNER JOIN Meistrijos ON Defektai.meistrija_id = Meistrijos.id"

	// connect to database
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	// query
	rows, err := db.Query(sqlProjection+sqlFrom+sqlWhere, sqlParamValue)
	if err != nil {
		return nil, err
	}

	var recs = []Rec{}

	for rows.Next() {
		var rec Rec
		err = rows.Scan(&rec.Id, &rec.Linija, &rec.Kelias, &rec.Km, &rec.Pk, &rec.M, &rec.Siule, &rec.Skodas, &rec.Suvirino, &rec.Operatorius, &rec.Aparatas, &rec.TData, &rec.Kelintas)
		if err != nil {
			return nil, err
		}
		recs = append(recs, rec)
	}
	rows.Close() //good habit to close
	return recs, err
}

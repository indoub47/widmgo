package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const driverName string = "sqlite3"
const dataSourceName string = "./widmgo.sqlite"

func insertRecs(recs []Rec) (int, error) {
	// connect to database
	conn, err := sql.Open(driverName, dataSourceName)
	defer conn.Close()

	if err != nil {
		return 0, err
	}

	// prepare insert statement
	stmt, err := conn.Prepare("INSERT INTO recs(id, linija, kelias, km, pk, m, siule, skodas, suvirino, operatorius, aparatas, tdata, kelintas, input) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?)")
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
	for _, r := range recs {
		// creating args slice
		//ps := rec.toSQLArgs()

		// excecuting the statement
		_, err = tx.Stmt(stmt).Exec(r.ID, r.Linija, r.Kelias, r.Km, r.Pk, r.M, r.Siule, r.Skodas, r.Suvirino, r.Operatorius, r.Aparatas, r.TData.Format("2006-01-02"), r.Kelintas, time.Now().Format("2006-01-02T15:04:05"))

		// if not success - rolling back transaction
		if err != nil {
			fmt.Println("DB INSERT error", err)
			fmt.Println("doing rollback")
			tx.Rollback()
			return 0, err
		}

		// incrementing the counter
		c++
	}

	tx.Commit()

	return c, nil
}

func fetchRecs(sqlWhere string, sqlParamValue int) ([]Rec, error) {

	// connect to database
	conn, err := sql.Open(driverName, dataSourceName)
	defer conn.Close()
	if err != nil {
		return nil, err
	}

	// start transaction
	// tx, err := conn.Begin()
	// if err != nil {
	// 	return nil, err
	// }

	// query
	sql := "SELECT * FROM recs WHERE sent IS NULL ORDER BY operatorius, kelintas"
	rows, err := conn.Query(sql)
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	// convert records to []Rec
	var recs = []Rec{}

	for rows.Next() {
		var rec Rec
		err = rows.Scan(&rec.ID, &rec.Linija, &rec.Kelias, &rec.Km, &rec.Pk, &rec.M, &rec.Siule, &rec.Skodas, &rec.Suvirino, &rec.Operatorius, &rec.Aparatas, &rec.TData, &rec.Kelintas)
		if err != nil {
			return nil, err
		}
		recs = append(recs, rec)
	}

	// done with rows
	rows.Close()

	// prepare update statement
	sql = "UPDATE recs SET sent = ? WHERE sent IS NULL"
	stmt, err := conn.Prepare(sql)
	if err != nil {
		return nil, err
	}

	// perform update
	_, err = stmt.Exec(time.Now())
	if err != nil {
		return nil, err
	}

	return recs, nil
}

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

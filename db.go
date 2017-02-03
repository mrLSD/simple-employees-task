package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type Employees struct {
	Id   int
	Name string
}

type EmployeeRoles struct {
	EmployeeId, RoleId int
	Enabled            bool
}

type Attendance struct {
	EmployeeId, RoleId int
	ActionTime         int
}

const sqlStmt = `
	create table Employees (Id integer not null primary key, Name text);
	create table EmployeeRoles (EmployeeId integer, RoleId integer, enabled boolean);
	create table Attendance (EmployeeId integer, RoleId integer, ActionTime integer);
	`

func InitTables(db *sql.DB) error {
	/*sqlStmt := `
	create table Employees (Id integer not null primary key, Name text);
	`*/
	_, err := db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("insert into Employees(id, name) values(?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	for i := 0; i < 4; i++ {
		_, err = stmt.Exec(i, fmt.Sprintf("Employee%03d", i))
		if err != nil {
			return err
		}
	}

	stmt2, err := tx.Prepare("insert into EmployeeRoles(EmployeeId, RoleId, enabled) values(?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt2.Close()
	_, err = stmt2.Exec(1, 1, 0)
	if err != nil {
		return err
	}
	_, err = stmt2.Exec(1, 2, 1)
	if err != nil {
		return err
	}
	_, err = stmt2.Exec(2, 1, true)
	if err != nil {
		return err
	}
	_, err = stmt2.Exec(2, 3, true)
	if err != nil {
		return err
	}

	tx.Commit()
	return nil
}

func GetEmployees(db *sql.DB) ([]Employees, error) {
	var e []Employees
	rows, err := db.Query("select id, name from Employees")
	if err != nil {
		return e, err
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			return e, err
		}
		e = append(e, Employees{
			Id:   id,
			Name: name,
		})
	}
	err = rows.Err()
	if err != nil {
		return e, err
	}
	return e, nil
}

func GetEmployeeRoles(id int, db *sql.DB) ([]EmployeeRoles, error) {
	var e []EmployeeRoles
	rows, err := db.Query("select EmployeeId, RoleId, Enabled from EmployeeRoles where EmployeeId=? AND enabled=?", id, 1)
	if err != nil {
		fmt.Printf("%v", err)
		return e, err
	}
	defer rows.Close()
	for rows.Next() {
		var eid, rid int
		var enabled bool
		err = rows.Scan(&eid, &rid, &enabled)
		if err != nil {
			fmt.Printf("%v", err)
			return e, err
		}
		e = append(e, EmployeeRoles{
			EmployeeId: eid,
			RoleId:     rid,
			Enabled:    enabled,
		})
	}
	err = rows.Err()
	if err != nil {
		fmt.Printf("%v", err)
		return e, err
	}
	return e, nil
}

func SaveAttendance(a Attendance, db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		fmt.Printf("%v", err)
		return err
	}
	stmt, err := tx.Prepare("insert into Attendance(EmployeeId, RoleId, ActionTime) values(?, ?, ?)")
	if err != nil {
		fmt.Printf("%v", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(a.EmployeeId, a.RoleId, a.ActionTime)
	if err != nil {
		fmt.Printf("%v", err)
		return err
	}
	tx.Commit()
	return nil
}

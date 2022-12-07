package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func connectToDB() {
	sqlCon, err := sql.Open("**validDriverName**", "**validDataSourceName**")
	if err != nil {
		panic(err.Error())
	}
	if err = sqlCon.Ping(); err != nil {
		panic(err.Error())
	}
	fmt.Println("Connected to database")
	db = sqlCon
}

func dropTables() {
	_, err := db.Exec("SET FOREIGN_KEY_CHECKS = 0")
	if err != nil {
		println("failed to drop constraints")
	}
	_, err = db.Exec("DROP TABLE IF EXISTS users")
	if err != nil {
		println("failed to drop users table")
	}
	_, err = db.Exec("DROP TABLE IF EXISTS userChanges")
	if err != nil {
		println("failed to drop userChanges table")
	}
	println("dropped tables")
}

func initDB() {
	_, err := db.Exec("CREATE DATABASE IF NOT EXISTS **NAME OF DB**") //#####################################  add **name of db** here #####################################
	if err == nil {
		println("db created")
	} else {
		println("db creation failed")
	}
	userTable()
	UserChangesTable()
}

func userTable() {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS users (id VARCHAR(27) PRIMARY KEY, name VARCHAR(20), password VARCHAR(64), power INT, lastActiveAt DATETIME, createdAt DATETIME, registeredAt DATETIME NULL)")
	if err == nil {
		println("user table created")
	} else {
		println("user table creation failed")
	}
}

func UserChangesTable() {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS userChanges (id INT AUTO_INCREMENT PRIMARY KEY, userid VARCHAR(27), columnName VARCHAR(30), newValue VARCHAR(100), oldValue VARCHAR(100), changedAt DATETIME, FOREIGN KEY (userid) REFERENCES users(id))")
	if err == nil {
		println("userChanges table created")
	} else {
		fmt.Println(err)
	}
}

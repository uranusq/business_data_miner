package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Database ... Contains methods to work with data
type Database struct {
	gorm.DB
}

// PrintInfo ...
func (db *Database) PrintInfo() {
	industries := []Industries{}
	db.Find(&industries)
	fmt.Println("Industries in DB: ", len(industries))

	companies := []Companies{}
	db.Find(&companies)
	fmt.Println("Companies in DB: ", len(companies))
}

func initizalizeDB(path string) *Database {
	gdb, err := gorm.Open("sqlite3", path)
	if err != nil {
		panic("failed to connect database")
	}
	defer gdb.Close()

	// Migrate the schema
	gdb.AutoMigrate(&Industries{}, &Companies{})
	db := Database{DB: *gdb}
	fmt.Println(gdb)
	return &db
}

func main() {
	//db := Database{}
	db := initizalizeDB("test.db")
	db.PrintInfo()
}

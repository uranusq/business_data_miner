package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Database ... Contains methods to work with data
type Database struct {
	*gorm.DB
	busyCrawlIDs  []int
	busyCollyIDs  []int
	busyGoogleIDs []int
}

// OpenInitialize ... Open connection to DB and Initializes custom structure. Connection needs to be closed
func (db *Database) OpenInitialize(path string) {
	gdb, err := gorm.Open("sqlite3", path)
	if err != nil {
		panic("failed to connect database")
	}
	//defer gdb.Close()

	gdb.AutoMigrate(&Industries{}, &Companies{})
	db.DB = gdb
	db.Exec("PRAGMA foreign_keys = ON;")

	// Exclude 0 indexes, since they always have empty values in SQLite
	db.busyCollyIDs = []int{0}
	db.busyCrawlIDs = []int{0}
	db.busyGoogleIDs = []int{0}
}

// PrintInfo ... Prints basic info about items in database
func (db *Database) PrintInfo() {
	industries := []Industries{}
	db.DB.Find(&industries)
	fmt.Println("Industries in DB: ", len(industries))

	companies := []Companies{}
	db.Find(&companies)
	fmt.Println("Companies in DB: ", len(companies))
}

func (db *Database) GetCrawlURL() {

}

func (db *Database) GetCollyURL() (string, error) {
	company := Companies{}
	db.Where("id NOT IN (?)", db.busyCollyIDs).Where(&Companies{IsCollyCrawled: false}).Find(&company)
	if company.ID == 0 {
		return "", errors.New("[GetCollyURL] no URLs found")
	}
	db.busyCollyIDs = append(db.busyCollyIDs, company.ID)
	fmt.Println("busy colly: ", db.busyCollyIDs)
	return company.URL, nil
}

func (db *Database) FinishedCollyURL() {

}

func (db *Database) GetGoogleURL() {

}

func (db *Database) fillToDebug() {
	testIndustr := []Industries{
		Industries{Industry: "Internet Services"},
		Industries{Industry: "Software Developement"},
		Industries{Industry: "Education"},
		Industries{Industry: "Retail"},
	}

	for _, ind := range testIndustr {
		db.Create(&ind)
	}

	testCompanies := []Companies{
		Companies{URL: "https://domru.ru", Industry: "Internet Services"},
		Companies{URL: "https://innopolis.ru", Industry: "Education"},
		Companies{URL: "https://tattelecom.ru/", Industry: "Internet Services"},
		Companies{URL: "https://www.wikipedia.org", Industry: "Education"},
		Companies{URL: "https://kai.ru", Industry: "Education"},
		Companies{URL: "https://www.acronis.com", Industry: "Software Developement"},
		Companies{URL: "https://www.kaspersky.ru", Industry: "Software Developement"},
		Companies{URL: "http://pivoman-kazan.ru", Industry: "Retail"},
	}
	for _, comp := range testCompanies {
		isPossible := db.NewRecord(&comp)
		if !isPossible {
			log.Println("Cannot insert value: ", &comp)
		}
		db.Create(&comp)
	}
}

func main() {
	//os.Remove("test.db")
	db := Database{}
	db.OpenInitialize("test.db")
	db.fillToDebug()

	db.PrintInfo()
	collyURL, err := db.GetCollyURL()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("CollyURL: ", collyURL)

	collyURL, err = db.GetCollyURL()
	fmt.Println("CollyURL: ", collyURL)

	collyURL, err = db.GetCollyURL()
	fmt.Println("CollyURL: ", collyURL)

	db.Close()
}

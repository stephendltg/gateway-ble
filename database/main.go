package database

import (
	"log"
	"os"
	"time"

	debug "github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Table beacons
type Beacon struct {
	gorm.Model
	Mac  string `gorm:"unique;not null"`
	Name string
}

// Client Database
var db gorm.DB

// Database connection
func Connect(mode bool) {

	debugger := debug.WithFields(debug.Fields{"package": "SQL"})

	// Debug level
	levelDebug := logger.Silent
	if mode {
		levelDebug = logger.Info
	}

	// SQL logger
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  levelDebug,  // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)

	// SQL
	DB, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		debugger.Fatal("failed to connect database: ", err)
	} else {
		db = *DB
	}

	// Migrate the schema
	DB.AutoMigrate(&Beacon{})
}

// Add beacon
func AddBeacon(data Beacon) (uint, error) {
	result := db.Create(&data)
	return data.ID, result.Error
}

// Count
func Count(table string) (int64, error) {
	var count int64
	result := db.Table(table).Count(&count)
	return count, result.Error
}

package utils

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB1 *gorm.DB
var DB2 *gorm.DB
var DB *gorm.DB

// todo : replace last variable with spread notation "..."
func ConnectDatabase(DBDriver string, DBSource1 string, DBSource2 string) {

	var db1 *gorm.DB
	var driver gorm.Dialector
	var err error

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,       // Don't include params in the SQL log
			Colorful:                  false,       // Disable color
		},
	)

	gormConfig := &gorm.Config{
		Logger: newLogger,
	}

	if DBDriver == "postgres" {
		driver = postgres.Open(DBSource1)
	} else {
		driver = mysql.Open(DBSource1)
	}

	if db1, err = gorm.Open(driver, gormConfig); err != nil {
		panic("Failed to connect to database!")
	}

	DB1 = db1

	if DBSource2 != "" {
		var db2 *gorm.DB

		if DBDriver == "postgres" {
			driver = postgres.Open(DBSource2)
		} else {
			driver = mysql.Open(DBSource2)
		}

		if db2, err = gorm.Open(driver, gormConfig); err != nil {
			panic("Failed to connect to database!")
		}

		DB2 = db2
	}

	DB = db1
}

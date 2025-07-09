package db

import (
	"log"
	"os"
	"path"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/nakurai/gourl/utils"
)



var Db *gorm.DB

func Init() error {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			IgnoreRecordNotFoundError: true,           // Ignore ErrRecordNotFound error for logger
		},
	)
	dbFilePath := path.Join(utils.DataDirPath, "data.db")
	var err error
	Db, err = gorm.Open(sqlite.Open(dbFilePath), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return err
	}
	return nil
}

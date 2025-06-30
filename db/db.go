package db

import (
	"path"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/nakurai/gourl/utils"
)

var Db *gorm.DB

func Init() error {
	dbFilePath := path.Join(utils.DataDirPath, "data.db")
	var err error
	Db, err = gorm.Open(sqlite.Open(dbFilePath), &gorm.Config{})
	if err != nil {
		return err
	}
	return nil
}

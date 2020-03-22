package storage

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/hailinluo/data-collector/logger"
	"github.com/hailinluo/data-collector/storage/structs"
	"github.com/hailinluo/data-collector/utils"
	"io"
)

var db *xorm.Engine

func InitDB(dataSource string) (io.Closer, error) {
	var err error
	db, err = xorm.NewEngine("mysql", dataSource)
	if err != nil {
		return nil, err
	}

	err = db.Sync2(new(structs.CompanyInfo))
	if err != nil {
		logger.Errorf("sync company table failed. err: %s", err)
	}

	err = db.Sync2(new(structs.Fund))
	if err != nil {
		logger.Errorf("sync fund table failed. err: %s", err)
	}

	var closer utils.Closer
	closer.AppendCloser(db)
	return &closer, nil
}

func DbEngine() *xorm.Engine {
	if db == nil {
		logger.Errorf("db engine is null")
	}
	return db
}


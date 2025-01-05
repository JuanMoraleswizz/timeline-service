package database

import (
	"fmt"
	"sync"

	"github.com/labstack/gommon/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"uala.com/timeline-service/config"
)

type mariadbDatabase struct {
	db *gorm.DB
}

var (
	once       sync.Once
	dbInstance *mariadbDatabase
)

func NewMariadbDatabase(conf *config.Config) *mariadbDatabase {
	once.Do(func() {
		dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
			conf.Db.User, conf.Db.Password, conf.Db.Host, conf.Db.Port, conf.Db.DBName, conf.Db.Charset)
		db, err := gorm.Open(mysql.Open(dns), &gorm.Config{})
		fmt.Println(dns)
		if err != nil {
			panic("failed to connect database")
		}
		dbInstance = &mariadbDatabase{db: db}
		log.Info("Database connected")

	})
	return dbInstance
}

func (m *mariadbDatabase) GetDb() *gorm.DB {
	return m.db
}

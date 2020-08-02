package db

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/xhyonline/xutil/xlog"
)

var log = xlog.Get(false)

// Config 数据库配置
type Config struct {
	Host     string `default:"mysql"`
	Port     string `default:"3306"`
	User     string `default:"root"`
	Password string `default:"root"`
	Name     string
	Lifetime int `default:"3000"`
}

// NewDataBase 实例化一个数据库
func NewDataBase(c *Config) *gorm.DB {
	var db *gorm.DB
	var err error
	for {
		db, err = gorm.Open("mysql", c.User+":"+c.Password+
			"@tcp("+c.Host+":"+c.Port+")/"+c.Name+
			"?charset=utf8mb4&parseTime=True&loc=Local&timeout=90s")
		if err != nil {
			log.WithError(err).Warn("waiting for connect to db")
			time.Sleep(time.Second * 2)
			continue
		}
		db.DB().SetConnMaxLifetime(time.Duration(c.Lifetime) * time.Second)
		log.Info("Mysql connect successful.")
		break
	}
	return db
}

package db

import (
	"time"

	"github.com/xhyonline/xutil/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Config 数据库配置
type Config struct {
	Host     string `default:"mysql"`
	Port     string `default:"3306"`
	User     string `default:"root"`
	Password string `default:"root"`
	Name     string
	Lifetime int `default:"3000"`
	// 此连接池和 kv 包中有所不同,连接池中的连接你不需要主动获取和释放
	// 你仅仅只需要 NewDataBase() 一个实例,直接获取实例即可,当执行 orm 操作时,它内部自动会为你维护连接池。这一切都是无感知的
	// 你可以通过查看 mysql 终端中的连接数来看到这一现象
	MaxActiveConn int `default:"100"` // 设置数据库连接池最大连接数
	MaxIdleConn   int `default:"20"`  // 连接池最大允许的空闲连接数，如果没有sql任务需要执行的连接数大于20，超过的连接会被连接池关闭。

}

// NewDataBase 实例化一个数据库
func NewDataBase(c *Config) *gorm.DB {
	var db *gorm.DB
	var err error
	for {
		db, err = gorm.Open(mysql.Open(c.User+":"+c.Password+
			"@tcp("+c.Host+":"+c.Port+")/"+c.Name+
			"?charset=utf8mb4&parseTime=True&loc=Local&timeout=90s"), &gorm.Config{})
		if err != nil {
			logger.Errorf("waiting for connect to db")
			time.Sleep(time.Second * 2)
			continue
		}
		dbs, err := db.DB()
		if err != nil {
			logger.Errorf("waiting for connect to db")
			time.Sleep(time.Second * 2)
			continue
		}
		dbs.SetConnMaxLifetime(time.Duration(c.Lifetime) * time.Second)
		dbs.SetMaxOpenConns(c.MaxActiveConn)
		dbs.SetMaxIdleConns(c.MaxIdleConn)
		logger.Info("Mysql connect successful.")
		break
	}
	return db
}

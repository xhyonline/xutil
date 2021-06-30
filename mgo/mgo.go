// MongoDB 工具包
package mgo

import (
	"time"

	"github.com/xhyonline/xutil/xlog"

	"gopkg.in/mgo.v2"
)

var log = xlog.Get().Debugger()

type Config struct {
	Database string // 数据库
	User     string // 账号
	Password string // 密码
	Host     string
	Port     string `default:"27017"`
}

// New mongodb client
func New(config Config) (*mgo.Database, error) {
	m, err := mgo.Dial(config.Host + ":" + config.Port)
	// 连接远程,Ping 不通一直到 Ping 通为止
	for {
		if err != nil {
			log.Warn(err)
			time.Sleep(time.Second * 3)
			continue
		}
		err := m.Ping()
		if err != nil {
			log.Warn(err)
			time.Sleep(time.Second * 3)
			continue
		}
		break
	}
	log.Info("Mongo Connect Success......")
	// 当账号密码不为空时
	// 登录 MongoDB
	if config.User != "" || config.Password != "" {
		err = m.Login(&mgo.Credential{
			Username:    config.User,
			Password:    config.Password,
			Source:      config.Database,
			ServiceHost: config.Host,
		})
		if err != nil {
			return nil, err
		}
	}
	m.SetMode(mgo.Monotonic, true)
	// defer m.Clone()	可不能关闭
	return m.DB(config.Database), nil
}

// mail 邮件发送包

package mail

import "gopkg.in/gomail.v2"

// Config 定义实例
type Config struct {
	User string
	Pass string
	Host string
	Port int
}

// Validate 基本验证
func (c Config) Validate() bool {
	return c.User != "" && c.Pass != "" && c.Host != "" && c.Port != 0
}

// Mail 邮件实例
type Mail struct {
	config  Config
	message *gomail.Message
	dialer  *gomail.Dialer
}

// Send 发送邮件
// subject 主题 、 from 自定义 、body 自定义 toUser 目标邮件
// 示例"
// m.Send("测试主题", "来自服务告警", "<p style='color:red'>线上服务告警</p>", "xxxx@qq.com")
func (m *Mail) Send(subject, from, body string, toUser ...string) error {
	m.message.SetHeader("Subject", subject)
	m.message.SetHeader("From", m.message.FormatAddress(m.config.User, from))
	m.message.SetHeader("To", toUser...)
	m.message.SetBody("text/html", body)
	return m.dialer.DialAndSend(m.message)
}

// NewMail 获取一个邮件
// 使用示例: 获取一个 Mail 实例
//  mail.NewMail(mail.Config{
//		User: "xxxxx@qq.com",
//		Pass: "授权码",
//		Host: "smtp.qq.com",
//		Port: 465,
//	})
func NewMail(c Config) *Mail {
	if !c.Validate() {
		return nil
	}
	m := gomail.NewMessage()
	return &Mail{config: c, dialer: gomail.NewDialer(c.Host, c.Port, c.User, c.Pass), message: m}
}

package beego

type respData struct {
	Code int         `json:"code"`
	Msg  string      `json:"message"`
	Data interface{} `json:"data"`
}

// NewResponse 数据响应格式
func NewResponse(code int, msg string, data interface{}) *respData {
	return &respData{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

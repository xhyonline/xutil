package g

// R 即 Response 响应标准
// 配合 c.Json(200,g.R(1,"请求成功",data))
func R(code int, message string, data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"code":    code,
		"message": message,
		"data":    data,
	}
}

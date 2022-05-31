package zdpgo_nntp

/*
@Time : 2022/5/31 17:40
@Author : 张大鹏
@File : response.go
@Software: Goland2021.3.1
@Description:
*/

// Response 响应结果
type Response struct {
	StatusCode int    `json:"status_code"`
	Content    []byte `json:"content"`
	Text       string `json:"text"`
	Uuid       string `json:"uuid"`
}

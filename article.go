package zdpgo_nntp

import "time"

/*
@Time : 2022/5/31 19:52
@Author : 张大鹏
@File : article.go
@Software: Goland2021.3.1
@Description:
*/

type Article struct {
	Uuid       string    `json:"uuid"`
	Group      string    `json:"group"`
	Author     string    `json:"author"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	Date       int       `json:"date"`
	DateString string    `json:"date_string"`
	DateTime   time.Time `json:"date_time"`
	MessageId  string    `json:"message_id"`
}

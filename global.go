package zdpgo_nntp

/*
@Time : 2022/5/31 16:42
@Author : 张大鹏
@File : global.go
@Software: Goland2021.3.1
@Description:
*/

import (
	"github.com/zhangdapeng520/zdpgo_password"
)

var (
	auths    map[string]Auth
	password = zdpgo_password.New()
)

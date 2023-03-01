// Copyright (c) 2023 MiaoZhongLuo(罗妙忠) &lt;2383304714@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/miao-crispy-corner/go_gateway_new.

package errno

var (
	// ErrUserAlreadyExist 代表用户已经存在.
	ErrUserAlreadyExist = &Errno{HTTP: 400, Code: "FailedOperation.UserAlreadyExist", Message: "User already exist."}
)

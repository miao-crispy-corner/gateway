// Copyright (c) 2023 MiaoZhongLuo(罗妙忠) &lt;2383304714@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/miao-crispy-corner/go_gateway_new.

package main

import (
	"github/miao-crispy-corner/go_gateway_new/internal/gateway"
	_ "go.uber.org/automaxprocs"
	"os"
)

// Go 程序的默认入口函数(主函数).
func main() {
	command := gateway.NewGateWayCommand()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}

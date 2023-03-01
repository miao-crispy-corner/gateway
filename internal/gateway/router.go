// Copyright (c) 2023 MiaoZhongLuo(罗妙忠) &lt;2383304714@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/miao-crispy-corner/go_gateway_new.

package gateway

import (
	"github.com/gin-gonic/gin"
	"github.com/miao-crispy-corner/gateway/internal/gateway/controller/v1/user"
	"github.com/miao-crispy-corner/gateway/internal/gateway/store"
	"github.com/miao-crispy-corner/gateway/internal/pkg/core"
	"github.com/miao-crispy-corner/gateway/internal/pkg/errno"
	"github.com/miao-crispy-corner/gateway/internal/pkg/log"
)

func installRouters(g *gin.Engine) error {
	// 注册 404 Handler.
	g.NoRoute(func(c *gin.Context) {
		core.WriteResponse(c, errno.ErrPageNotFound, nil)
	})

	// 注册 /healthz handler.
	g.GET("/healthz", func(c *gin.Context) {
		log.C(c).Infow("Healthz function called")

		core.WriteResponse(c, nil, map[string]string{"status": "ok"})
	})

	uc := user.New(store.S)

	// 创建 v1 路由分组
	v1 := g.Group("/v1")
	{
		// 创建 users 路由分组
		userv1 := v1.Group("/users")
		{
			userv1.POST("", uc.Create)
		}
	}

	return nil
}

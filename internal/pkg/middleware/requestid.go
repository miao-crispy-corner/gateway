// Copyright (c) 2023 MiaoZhongLuo(罗妙忠) &lt;2383304714@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/miao-crispy-corner/go_gateway_new.

package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/miao-crispy-corner/gateway/internal/pkg/known"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查请求头中是否有 `X-Request-ID`，如果有则复用，没有则新建
		requestID := c.Request.Header.Get(known.XRequestIDKey)

		if requestID == "" {
			requestID = uuid.New().String()
		}

		// 将 RequestID 保存在 gin.Context 中，方便后边程序使用
		c.Set(known.XRequestIDKey, requestID)

		// 将 RequestID 保存在 HTTP 返回头中，Header 的键为 `X-Request-ID`
		c.Writer.Header().Set(known.XRequestIDKey, requestID)
		c.Next()
	}
}

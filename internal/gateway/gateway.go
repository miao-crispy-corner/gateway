// Copyright (c) 2023 MiaoZhongLuo(罗妙忠) &lt;2383304714@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/miao-crispy-corner/go_gateway_new.

package gateway

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/miao-crispy-corner/gateway/internal/pkg/log"
	"github.com/miao-crispy-corner/gateway/internal/pkg/middleware"
	"github.com/miao-crispy-corner/gateway/pkg/version/verflag"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var cfgFile string

func NewGateWayCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gateway",
		Short: "api gateway",
		Long: `Find more gateway information at:
"https://github.com/miao-crispy-corner/go_gateway_new"`,
		SilenceUsage: true,
		RunE: func(command *cobra.Command, args []string) error {
			// 如果 `--version=true`，则打印版本并退出
			verflag.PrintAndExitIfRequested()

			// 初始化日志
			log.Init(logOptions())
			defer log.Sync() // Sync 将缓存中的日志刷新到磁盘文件中
			return run()
		},
		Args: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}
			return nil
		},
	}

	cobra.OnInitialize(initConfig)
	cmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "The path to the gateway configuration file.")
	cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	// 添加 --version 标志
	verflag.AddFlags(cmd.PersistentFlags())
	return cmd
}

func run() error {
	// 初始化 store 层
	if err := initStore(); err != nil {
		return err
	}

	// 设置 Gin 模式
	gin.SetMode(viper.GetString("runmode"))

	// 创建 Gin 引擎
	g := gin.New()

	middlewares := []gin.HandlerFunc{gin.Recovery(), middleware.NoCache, middleware.Cors, middleware.Secure, middleware.RequestID()}

	g.Use(middlewares...)

	if err := installRouters(g); err != nil {
		return err
	}

	// 创建 HTTP Server 实例
	httpsrv := &http.Server{Addr: viper.GetString("addr"), Handler: g}

	// 运行 HTTP 服务器
	// 打印一条日志，用来提示 HTTP 服务已经起来，方便排障
	go func() {
		log.Infow("Start to listening the incoming requests on http address", "addr", viper.GetString("addr"))
		if err := httpsrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalw(err.Error())
		}
	}()

	// 等待中断信号优雅地关闭服务器（10 秒超时)。
	quit := make(chan os.Signal)
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的 CTRL + C 就是触发系统 SIGINT 信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	log.Infow("Shutting down server ...")

	// 创建 ctx 用于通知服务器 goroutine, 它有 10 秒时间完成当前正在处理的请求
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 10 秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过 10 秒就超时退出
	if err := httpsrv.Shutdown(ctx); err != nil {
		log.Errorw("Insecure Server forced to shutdown", "err", err)
		return err
	}

	log.Infow("Server exiting")

	return nil
}

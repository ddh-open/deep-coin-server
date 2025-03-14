// Package middleware Copyright 2021 jianfengye.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package middleware

import (
	"context"
	gin2 "devops-http/framework/gin"
	"fmt"
	"log"
	"time"
)

func Timeout(d time.Duration) gin2.HandlerFunc {
	// 使用函数回调
	return func(c *gin2.Context) {
		finish := make(chan struct{}, 1)
		panicChan := make(chan interface{}, 1)
		// 执行业务逻辑前预操作：初始化超时context
		durationCtx, cancel := context.WithTimeout(c.BaseContext(), d)
		defer cancel()

		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()
			// 使用next执行具体的业务逻辑
			c.Next()

			finish <- struct{}{}
		}()
		// 执行业务逻辑后操作
		select {
		case p := <-panicChan:
			c.DSetStatus(500).DJson("time out")
			log.Println(p)
		case <-finish:
			fmt.Println("finish")
		case <-durationCtx.Done():
			c.DSetStatus(500).DJson("time out")
		}
	}
}

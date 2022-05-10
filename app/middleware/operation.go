package middleware

import (
	"bytes"
	"devops-http/app/module/base/utils"
	"devops-http/app/module/sys/model/operation"
	"encoding/json"
	"github.com/ddh-open/gin/framework/contract"
	"github.com/ddh-open/gin/framework/gin"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func OperationRecord() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := c.MustMakeLog()
		db, err := c.MustMake(contract.ORMKey).(contract.ORMService).GetDB()
		if err == nil {
			var body []byte
			var username string
			if c.Request.Method != http.MethodGet {
				var err error
				body, err = ioutil.ReadAll(c.Request.Body)
				if err != nil {
					logger.Error("read body from request error:", zap.Error(err))
				} else {
					c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
				}
			} else {
				query := c.Request.URL.RawQuery
				query, _ = url.QueryUnescape(query)
				split := strings.Split(query, "&")
				m := make(map[string]string)
				for _, v := range split {
					kv := strings.Split(v, "=")
					if len(kv) == 2 {
						m[kv[0]] = kv[1]
					}
				}
				body, _ = json.Marshal(&m)
			}
			tokenUser, _ := utils.ParseToken(c)
			if tokenUser != nil && tokenUser.Id != 0 {
				username = tokenUser.Username
			}
			record := operation.DevopsSysOperationRecord{
				Ip:       c.ClientIP(),
				Method:   c.Request.Method,
				Path:     c.Request.URL.Path,
				Agent:    c.Request.UserAgent(),
				Body:     string(body),
				UserName: username,
			}
			// 上传文件时候 中间件日志进行裁断操作
			if strings.Index(c.GetHeader("Content-Type"), "multipart/form-data") > -1 {
				if len(record.Body) > 512 {
					record.Body = "File or Length out of limit"
				}
			}
			writer := responseBodyWriter{
				ResponseWriter: c.Writer,
				body:           &bytes.Buffer{},
			}
			c.Writer = writer
			now := time.Now()

			c.Next()

			latency := time.Since(now)
			record.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()
			record.Status = c.Writer.Status()
			record.Latency = latency
			if len(writer.body.String()) < 500 {
				record.Resp = writer.body.String()
			}
			if err = db.Model(&operation.DevopsSysOperationRecord{}).Create(&record).Error; err != nil {
				logger.Error("create operation record error:", zap.Error(err))
			}
		} else {
			logger.Error(err.Error())
		}
	}
}

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

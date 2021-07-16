package ginc

import (
	"encoding/json"
	"fmt"
	"os"
	"source/server/util/lang"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"gopkg.in/resty.v1"
)

// Contexter wh 接口
type Contexter interface {
	ParseReqbody(interface{}) error         // 解析body入参
	Success(map[string]interface{})         // 标识成功返回状态
	SuccessByte(b []byte, params ...string) // 原生数据
	Fail(interface{})                       // 标识失败返回状态
	FailErr(error)                          // 标识失败返回状态
	GinCtx() *gin.Context                   // 获取gin context
	ClientPost(api string, params interface{}, options ...map[string]string) (
		*resty.Response, error) // Post请求
	ClientGet(api string, params interface{}, options ...map[string]string) (
		*resty.Response, error) // Get 请求
	ParamRoute(string) string // 路由参数解析
}

// Context ...
type Context struct {
	Gin    *gin.Context
	Client *Client
	keys   map[string]interface{}
	m      sync.Mutex
}

// Hander ..
func Hander(fn func(c Contexter)) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := &Context{Gin: c}
		fn(ctx)
	}
}

// GinCtx gin context
func (c *Context) GinCtx() *gin.Context {
	return c.Gin
}

// ParamRoute ..
func (c *Context) ParamRoute(key string) string {
	return strings.Replace(c.Gin.Param("id"), "/", "", -1)
}

// ParseReqbody 解析网络body
func (c *Context) ParseReqbody(reqbody interface{}) error {
	if err := c.Gin.ShouldBindJSON(reqbody); err != nil {
		c.Fail(map[string]interface{}{
			"errorcode": 300400,
			"errormsg":  fmt.Sprintf("参数解析 :%s", err.Error()),
		})
		return err
	}
	return nil
}

// Success 成功
func (c *Context) Success(data map[string]interface{}) {
	if data["errorcode"] == nil {
		data["errorcode"] = 0
	}
	if b, err := json.Marshal(data); err != nil {
		panic(err)
	} else {
		c.Gin.Data(200, "application/json", b)
	}
}

// SuccessByte 成功
func (c *Context) SuccessByte(b []byte, params ...string) {
	contentType := append(params, "application/json")[0]
	c.Gin.Data(200, contentType, b)

}

// Fail 失败
func (c *Context) Fail(data interface{}) {
	if b, err := json.Marshal(data); err != nil {
		panic(err)
	} else {
		c.Gin.Data(200, "application/json", b)
	}
}

func (c *Context) FailErr(err error) {
	c.Fail(map[string]interface{}{
		"errorcode": 500,
		"errormsg":  err,
	})
}

// rget rewrite Gin
func (c *Context) rget(key string) (value interface{}, exists bool) {
	value, exists = c.keys[key]
	return
}

// rset rewrite Gin
func (c *Context) rset(key string, value interface{}) {
	if c.keys == nil {
		c.keys = make(map[string]interface{})
	}
	c.keys[key] = value
}

// rgetint rewrite Gin
func (c *Context) rgetint(key string) (i int) {
	if val, ok := c.rget(key); ok && val != nil {
		i, _ = val.(int)
	}
	return
}

// ClientPost method
func (c *Context) ClientPost(api string, params interface{}, options ...map[string]string) (
	*resty.Response, error) {
	return c.request(c.Client.PostRow)(api, params, options...)
}

// ClientGet method
func (c *Context) ClientGet(api string, params interface{}, options ...map[string]string) (
	*resty.Response, error) {
	return c.request(c.Client.GetRow)(api, params, options...)
}

type requestHander func(string, interface{}, ...map[string]string) (
	*resty.Response, map[string]interface{}, error)

type requestFunc func(string, interface{}, ...map[string]string) (
	*resty.Response, error)

// request 公共部分(锁操作)
func (c *Context) request(fn requestHander) requestFunc {
	return func(api string, params interface{}, options ...map[string]string) (*resty.Response, error) {
		// fmt.Printf("options===>%+v\n", options)

		// 内容服务传递监控headers
		if strings.Contains(api, os.Getenv("KONG_CONTENT_SERVER_INTERNAL")) {
			headers := c.Gin.GetStringMapString("headers")
			if headers != nil {
				if len(options) == 1 {
					options = []map[string]string{lang.MapAssign(headers, options[0])}
				} else {
					options = []map[string]string{headers}
				}
			}
		}

		// 过滤正版资源
		params = c.fillSourceFilter(params)
		resp, log, err := fn(api, params, options...)

		c.m.Lock()
		total := c.rgetint("REQUEST_COUNT")
		if total == 0 {
			c.rset("REQUEST_COUNT", 1)
		} else {
			c.rset("REQUEST_COUNT", total+1)
		}
		if sourceTemp, ok := c.rget("REQUEST_SOURCE"); ok {
			source := sourceTemp.([]map[string]interface{})
			source = append(source, log)
			c.rset("REQUEST_SOURCE", source)
		} else {
			c.rset("REQUEST_SOURCE", []map[string]interface{}{log})
		}

		c.m.Unlock()
		return resp, err
	}
}

func (c *Context) fillSourceFilter(params interface{}) interface{} {

	sourceFilter := c.Gin.GetString("sourceFilter")
	if sourceFilter == "1" {
		if newparams, ok := params.(map[string]string); ok {
			newparams["sourceFilter"] = sourceFilter
			return newparams
		}
		if newparams, ok := params.(map[string]interface{}); ok {
			newparams["sourceFilter"] = sourceFilter
			return newparams
		}
		// if b, err := json.Marshal(params); err == nil {
		// 	fmt.Println(string(b))
		// }
	}
	return params
}

// NewContext ..
func NewContext() Contexter {
	return &Context{Gin: &gin.Context{}}
}

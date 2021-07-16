package ginc

import (
	"errors"
	"fmt"
	"source/server/util/lang"
	"strconv"
	"time"

	"gopkg.in/resty.v1"
)

// MethodType 方法类型
type MethodType int

const (
	_ MethodType = iota
	// POST ..
	POST
	// GET ..
	GET
	// PUT ...
	PUT
)

// Client .
type Client struct{}

// PostRow method
func (c *Client) PostRow(api string, params interface{}, options ...map[string]string) (
	*resty.Response, map[string]interface{}, error) {
	headers, timeout := c.cargs(options...)
	return c.request(POST, api, params, headers, timeout)
}

// GetRow method
func (c *Client) GetRow(api string, params interface{}, options ...map[string]string) (
	*resty.Response, map[string]interface{}, error) {
	headers, timeout := c.cargs(options...)
	return c.request(GET, api, params, headers, timeout)
}

// Post Post method
func (c *Client) Post(api string, params interface{}, options ...map[string]string) (
	*resty.Response, error) {
	resp, _, err := c.PostRow(api, params, options...)
	return resp, err
}

// Get Post method
func (c *Client) Get(api string, params interface{}, options ...map[string]string) (
	*resty.Response, error) {
	resp, _, err := c.GetRow(api, params, options...)
	return resp, err
}

// cargs 转换扩展参数
func (c *Client) cargs(options ...map[string]string) (headers map[string]string, timeout time.Duration) {
	if len(options) == 1 {
		headers := options[0]
		timeoutKey := "TIME_OUT"
		if timeoutStr, ok := headers[timeoutKey]; ok {
			if timeout2, err := strconv.ParseFloat(timeoutStr, 10); err == nil {
				timeout = time.Duration(timeout2)
			}
			delete(headers, timeoutKey)
		}
		return headers, timeout
	}
	return nil, 0
}

func (c *Client) request(typer MethodType, api string, params interface{},
	headers map[string]string, timeout time.Duration) (resp *resty.Response, log map[string]interface{}, err error) {
	// 超时时间
	if timeout == 0 {
		timeout = time.Duration(5000)
	}
	// fmt.Println("timeout=====>", headers, time.Millisecond*timeout)
	req := resty.New().
		SetTimeout(time.Millisecond*timeout).
		R().
		SetHeader("Content-Type", "application/json")

	// if apienv := os.Getenv("API_ENV"); apienv != "" && util.Contains([]string{"LOCAL_BETA"}, apienv) {
	// 	req = req.SetHeader("Host", "kong.adam.svc.cluster.local")
	// }
	// 请求headers
	if headers != nil {
		req = req.SetHeaders(headers)
	}

	timer := time.Now()

	if typer == POST {
		resp, err = req.SetBody(params).Post(api)
	} else if typer == GET {
		if params == nil {
			resp, err = req.Get(api)
		} else if v, ok := params.(map[string]string); ok {
			resp, err = req.SetQueryParams(v).Get(api)
		}
	} else {
		return nil, nil, fmt.Errorf("NO *%s* METHOD", c.mtToStr(typer))
	}

	// 日志
	log = map[string]interface{}{
		"url":     api,
		"params":  params,
		"consume": time.Since(timer).Nanoseconds() / 1000000,
		"error":   lang.If3(resp.StatusCode() >= 500, true, false),
		"timeout": lang.If3(resp.StatusCode() == 0, true, false),
	}

	// fmt.Println("=====>", params, string(resp.Body()), resp.StatusCode(), err)
	if resp.StatusCode() == 0 {
		return nil, log, errors.New("error: timeout")
	}
	// 非期望状态
	if err != nil || 200 != resp.StatusCode() {
		return nil, log, errors.New("Not found or StatusCode !=200")
	}
	return resp, log, nil

}

// mtToStr 类型
func (c *Client) mtToStr(typer MethodType) string {
	var methods = map[MethodType]string{
		POST: "POST",
		GET:  "GET",
		PUT:  "PUT",
	}
	return methods[typer]
}

package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/Meland-Inc/logservice/internal/pkg/dapr/invoke"
	"github.com/rs/cors"
	"github.com/urfave/cli/v2"
)

// time: "1645371844"             // UTC时间戳, 精确到秒
// level: "info"                  // fatal, error, info, debug, trace
// msg: ""                        // 日志消息
// exception: ""                  // 错误详情和上下文堆栈, 没有则无此字段
// traceSpan: "mutaion_userLogin" // 跟踪的业务. 如果调用某个网络协议. 则表示该网络协议的耗时. 若非 trace 类型,则无该字段
// traceDuration: ""              // 跟踪的耗时.单位毫秒. 比如 1000 则表示1s 若非 trace 类型,则无该字段
// scope: "login"                 // 日志范围, 用于定位日志产生的业务模块, 方便针对一些特殊高要求的业务模块设置监控
// client: "Chrome v84"           // 客户端版本 Chrome 98.0.4758.102 / MelandClientV1.1
// userId: "5"                    // 如果有用户信息则带上, 若没有的场景则无这个字段
// version: "g5dfcf2e"
type Log struct {
	Time           string  `json:"time"`
	Level          string  `json:"level"`
	Msg            string  `json:"msg"`
	Exception      string  `json:"exception"`
	TraceSpan      string  `json:"traceSpan"`
	TraceDuration  string  `json:"traceDuration"`
	Scope          string  `json:"scope"`
	Client         string  `json:"client"`
	UserId         string  `json:"userId"`
	UserEthAddress *string `json:"userEthAddress,omitempty"`
	Version        string  `json:"version"`
}

type RequestX struct {
	Content []Log `json:"content"`
}

func printJSON(v interface{}) {
	b, err := json.Marshal(v)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(b))
}

func Serve(c *cli.Context) error {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	err := xinit(c)

	if err != nil {
		return err
	}

	mux := http.NewServeMux()

	/// All you need to do here is to print out the logs, which are automatically collected by the cluster collector
	mux.HandleFunc("/batch-logs", func(rw http.ResponseWriter, r *http.Request) {
		var c RequestX
		err := json.NewDecoder(r.Body).Decode(&c)

		if err != nil {
			printJSON(map[string]string{
				"error": err.Error(),
				"scope": "batch-logs",
			})
			return
		}

		for _, log := range c.Content {
			if log.UserId != "" {
				o, e := invoke.GetUserWeb3Profile(log.UserId)

				if e != nil {
					printJSON(map[string]string{
						"error": e.Error(),
						"scope": "batch-logs",
					})
				} else {
					log.UserEthAddress = &o.BlockchainAddress
				}
			}

			printJSON(log)
		}
	})

	// cors.Default() setup the middleware with default options being
	// all origins accepted with simple methods (GET, POST). See
	// documentation below for more options.
	handler := cors.Default().Handler(mux)
	http.ListenAndServe(fmt.Sprintf(":%s", port), handler)

	return nil
}

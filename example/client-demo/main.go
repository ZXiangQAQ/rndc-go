package main

import (
	"fmt"
	"github.com/ZXiangQAQ/rndc-go"
)

func main() {
	// 默认日志等级为 info
	rndc.SetLevelByString("error")

	// 创建 RnDC client
	client, err := rndc.NewRNDCClient("192.168.196.170:953", "hmac-sha256", "xRmH2XdFcDqWO91pYhiCwlZmWnaSO8EleBFu1uz8d3g=")
	if err != nil {
		panic(err)
	}

	// 请求 RnDC 服务器, 并同步获取结果
	resp, err := client.Call("sync test123.com")
	if err != nil {
		panic(err)
	}

	// 获取 bind 返回的原生 response
	fmt.Printf("response: %s\n", resp)

	// 获取指令内容
	fmt.Printf("cmd: %s\n", resp.Data.Type)

	// 获取返回码, 类型是 string
	fmt.Printf("result: %s\n", resp.Data.Result)

	// 获取返回内容, 类型是 string
	fmt.Printf("text: %s\n", resp.Data.Text)

	// 获取错误, 类型是 string
	fmt.Printf("err: %s\n", resp.Data.Err)
}

# rndc-go

> `bind9` 服务使用的 `rndc` 协议的 go 版本客户端，仅支持协议版本1

## Installation

```bash
go get github.com/ZXiangQAQ/rndc-go
```

## Quickstart

```go
import (
	"fmt"
	"github.com/ZXiangQAQ/rndc-go"
)

func main() {
	// 默认日志等级为 info
	rndc.SetLevelByString("error")

	// 创建 RnDC client
	client, err := rndc.NewRNDCClient("your bind9 server address", "algo", "secret")
	if err != nil {
    fmt.Printf("conn failed, err: %s\n", err.Error())
	}

	// 请求 RnDC 服务器, 并同步获取结果
	resp, err := client.Call("sync test123.com")
	if err != nil {
		fmt.Printf("client call err: %s\n", err.Error())
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
```

## Feature

- 目前使用硬编码进行序列化和反序列化，计划基于 `tag` 标签和 `rndc` 的编码规范，实现动态 `Marshal` 和 `Unmarshal`


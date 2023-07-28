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
	client, err := rndc.NewRNDCClient("your bind server", "algo", "secrect")
	if err != nil {
		panic(err)
	}
	resp, err := client.Call("sync test123.com")
	if err != nil {
		panic(err)
	}
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


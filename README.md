# NetProbe - 网络延迟探测工具

NetProbe 是一个简单的网络延迟探测工具，用于测试多个网站的访问延迟。该工具可以同时测试多个 URL 的响应时间，并以清晰的方式显示结果。

## 功能

- 支持同时测试多个网站的访问延迟
- 可配置超时时间和重试次数

## 使用方法

1. 创建配置文件 `config.json`，格式如下：

```json
{
    "urls": [
        "https://google.com",
        "https://github.com",
        "https://twitter.com",
        "https://www.bilibili.com",
        "https://www.baidu.com"
    ],
    "timeout_seconds": 5,
    "max_retries": 3
}
```

- `urls`: 要测试的网站 URL 列表
- `timeout_seconds`: 每个请求的超时时间（秒）
- `max_retries`: 最大重试次数

2. 运行程序：

```bash
go run main.go
```

程序会显示每个 URL 的测试结果，包括：

- 成功：显示访问延迟时间
- 失败：显示错误信息

## 编译

要编译该程序，请确保已安装 Go 语言环境。然后在项目目录下运行以下命令：

```bash
go build -o NetProbe.exe main.go
```

编译完成后，会在当前目录下生成 `NetProbe.exe` 可执行文件。

## 依赖要求

- Go 1.24 或更高版本

# dify-go

dify-go是一个用于与Dify API进行交互的Go客户端库。Dify是一个大型语言模型（LLM）应用开发平台，该库可帮助Go开发者快速集成和使用Dify的各种功能。

## 特性

- 支持流式和阻塞式对话模式
- 文件上传功能
- 支持停止正在进行的响应
- 获取建议问题列表
- 获取会话历史消息
- 支持用户身份识别
- 简单易用的API设计

## 安装

```bash
go get github.com/Davied-H/dify-go
```

## 环境配置

你可以使用环境变量或者.env文件来配置API密钥：

```bash
# .env文件示例
DIFY_API_KEY=your_api_key_here
```

使用godotenv加载环境变量：

```go
import (
    "github.com/joho/godotenv"
    "log"
)

func init() {
    // 加载.env文件中的环境变量
    err := godotenv.Load()
    if err != nil {
        log.Println("Error loading .env file")
    }
}
```

## 快速开始

### 初始化客户端

```go
import (
    "github.com/Davied-H/dify-go"
    "os"
    "context"
    "fmt"
    "net/http"
    "time"
)

// 创建客户端实例
client := dify.NewClient("https://api.dify.ai/v1")
```

### 配置客户端

每个API调用都需要传递API密钥：

```go
// 在每个请求中单独设置API密钥
resp, err := client.ChatMessage(context.TODO(), dify.ChatMessageOption{
    ApiKey: os.Getenv("DIFY_API_KEY"),
    // ... 其他选项
})
```

你也可以使用自定义的HTTP客户端：

```go
// 使用自定义配置
httpClient := &http.Client{
    Timeout: time.Second * 30,
}

config := dify.DefaultConfig("https://api.dify.ai/v1")
config.HttpClient = httpClient

client := dify.NewClientWithConfig(*config)
```

### 流式对话

```go
_, err := client.ChatMessage(context.TODO(), dify.ChatMessageOption{
    ApiKey: os.Getenv("DIFY_API_KEY"),
    OnEvent: func(ev dify.ChatMessageRespSSEData) {
        switch ev.Event {
        case "workflow_started":
            fmt.Printf("工作流开始执行\n")
        case "message":
            fmt.Printf("%s", ev.Answer)  // 实时接收回答内容
        case "workflow_finished":
            fmt.Println("\n工作流结束执行")
        }
    },
    RequestBody: dify.ChatMessageReq{
        Inputs: map[string]interface{}{
            "role": "助手",  // 自定义变量
        },
        Query:          "你好，请介绍一下自己",
        ResponseMode:   dify.ResponseModeStreaming,
        ConversationId: "",  // 留空表示新会话
        User:           "user_id",  // 用户标识
    },
})
```

### 阻塞式对话

```go
resp, err := client.ChatMessage(context.TODO(), dify.ChatMessageOption{
    ApiKey: os.Getenv("DIFY_API_KEY"),
    RequestBody: dify.ChatMessageReq{
        Query:          "你好，请介绍一下自己",
        ResponseMode:   dify.ResponseModeBlocking,
        User:           "user_id",
    },
})
fmt.Println("回答:", resp.Answer)
```

### 上传文件

```go
file, _ := os.Open("example.pdf")
defer file.Close()

resp, err := client.UploadFile(context.TODO(), dify.UploadFileOption{
    ApiKey: os.Getenv("DIFY_API_KEY"),
    RequestFormData: dify.UploadFileReq{
        File: file,
        User: "user_id",
    },
})
```

### 停止响应

```go
resp, err := client.StopTask(context.TODO(), dify.StopTaskOption{
    ApiKey: os.Getenv("DIFY_API_KEY"),
    TaskId: "task_id",  // 从ChatMessage响应中获取
    RequestBody: dify.StopTaskReq{
        User: "user_id",  // 必须与发送消息时相同
    },
})
```

### 获取建议问题

```go
resp, err := client.GetSuggested(context.TODO(), dify.GetSuggestedOption{
    ApiKey:    os.Getenv("DIFY_API_KEY"),
    MessageId: "message_id",  // 从ChatMessage响应中获取
    RequestParams: dify.GetSuggestedReq{
        User: "user_id",
    },
})
```

### 获取会话历史

```go
resp, err := client.GetMessages(context.TODO(), dify.GetMessagesOption{
    ApiKey: os.Getenv("DIFY_API_KEY"),
    RequestParams: dify.GetMessagesReq{
        ConversationId: "conversation_id",  // 从ChatMessage响应中获取
        User:           "user_id",
        Limit:          20,  // 每页记录数
    },
})
```

## 功能进度

- [x] 发送对话消息 /chat-messages
- [x] 上传文件 /files/upload
- [x] 停止响应 /chat-messages/:task_id/stop
- [ ] 消息反馈（点赞）
- [x] 获取下一轮建议问题列表 /messages/{message_id}/suggested
- [x] 获取会话历史消息 /messages
- [ ] 获取会话列表
- [ ] 删除会话
- [ ] 会话重命名
- [ ] 获取对话变量
- [ ] 语音转文字
- [ ] 文字转语音
- [ ] 获取应用基本信息
- [ ] 获取应用参数
- [ ] 获取应用Meta信息
- [ ] 获取标注列表
- [ ] 创建标注
- [ ] 更新标注
- [ ] 删除标注
- [ ] 标注回复初始设置
- [ ] 查询标注回复初始设置任务状态

## 贡献

欢迎提交问题和Pull Request。

## 许可证

MIT

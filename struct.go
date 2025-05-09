package dify

import (
	"mime/multipart"
	"os"
)

type ChatMessageOption struct {
	ApiKey      string `validate:"required"`
	OnEvent     func(ev ChatMessageRespSSEData)
	RequestBody ChatMessageReq
}
type ChatMessageReq struct {
	Inputs         map[string]interface{} `json:"inputs"`                    // 允许传入 App 定义的各变量值
	Query          string                 `json:"query" validate:"required"` // 用户输入/提问内容
	ResponseMode   string                 `json:"response_mode"`             // streaming: 流式模式, blocking: 阻塞模式
	ConversationId string                 `json:"conversation_id"`           // 会话 ID，需要基于之前的聊天记录继续对话
	User           string                 `json:"user" validate:"required"`  // 用户标识，可用于终止请求等
	Files          []struct {
		Type           string `json:"type"`
		TransferMethod string `json:"transfer_method"`
		Url            string `json:"url"`
	} `json:"files"`
}
type ChatMessageResp struct {
	Event          string `json:"event"`
	TaskId         string `json:"task_id"`
	Id             string `json:"id"`
	MessageId      string `json:"message_id"`
	ConversationId string `json:"conversation_id"`
	Mode           string `json:"mode"`
	Answer         string `json:"answer"`
	Metadata       struct {
		Usage struct {
			PromptTokens        int     `json:"prompt_tokens"`
			PromptUnitPrice     string  `json:"prompt_unit_price"`
			PromptPriceUnit     string  `json:"prompt_price_unit"`
			PromptPrice         string  `json:"prompt_price"`
			CompletionTokens    int     `json:"completion_tokens"`
			CompletionUnitPrice string  `json:"completion_unit_price"`
			CompletionPriceUnit string  `json:"completion_price_unit"`
			CompletionPrice     string  `json:"completion_price"`
			TotalTokens         int     `json:"total_tokens"`
			TotalPrice          string  `json:"total_price"`
			Currency            string  `json:"currency"`
			Latency             float64 `json:"latency"`
		} `json:"usage"`
		RetrieverResources []struct {
			Position     int     `json:"position"`
			DatasetId    string  `json:"dataset_id"`
			DatasetName  string  `json:"dataset_name"`
			DocumentId   string  `json:"document_id"`
			DocumentName string  `json:"document_name"`
			SegmentId    string  `json:"segment_id"`
			Score        float64 `json:"score"`
			Content      string  `json:"content"`
		} `json:"retriever_resources"`
	} `json:"metadata"`
	CreatedAt int `json:"created_at"`
}
type ChatMessageRespSSEData struct {
	Event                string   `json:"event"`
	ConversationId       string   `json:"conversation_id"`
	MessageId            string   `json:"message_id"`
	CreatedAt            int      `json:"created_at"`
	TaskId               string   `json:"task_id"`
	Id                   string   `json:"id"`
	Answer               string   `json:"answer"`
	FromVariableSelector []string `json:"from_variable_selector"`
}

type UploadFileOption struct {
	ApiKey          string        `validate:"required"`
	RequestFormData UploadFileReq `validate:"required"`
}
type UploadFileViaGinOption struct {
	ApiKey          string              `validate:"required"`
	RequestFormData UploadFileViaGinReq `validate:"required"`
}
type UploadFileReq struct {
	File *os.File `validate:"required"`
	User string   `validate:"required"`
}
type UploadFileViaGinReq struct {
	FormFile *multipart.FileHeader `validate:"required"`
	User     string                `validate:"required"`
}
type UploadFileResp struct {
	Id         string      `json:"id"`
	Name       string      `json:"name"`
	Size       int         `json:"size"`
	Extension  string      `json:"extension"`
	MimeType   string      `json:"mime_type"`
	CreatedBy  string      `json:"created_by"`
	CreatedAt  int         `json:"created_at"`
	PreviewUrl interface{} `json:"preview_url"`
}

type StopTaskOption struct {
	ApiKey      string `validate:"required"`
	TaskId      string `validate:"required"`
	RequestBody StopTaskReq
}
type StopTaskReq struct {
	User string ` json:"user" validate:"required"` // 用户标识，用于定义终端用户的身份，必须和发送消息接口传入 ChatMessageReq.User 保持一致。
}
type StopTaskResp struct {
	Result string `json:"result"`
}

type GetSuggestedOption struct {
	ApiKey        string `validate:"required"`
	MessageId     string `validate:"required"`
	RequestParams GetSuggestedReq
}
type GetSuggestedReq struct {
	User string `url:"user"`
}
type GetSuggestedResp struct {
	Result string   `json:"result"`
	Data   []string `json:"data"`
}

type GetMessagesOption struct {
	ApiKey        string `validate:"required"`
	RequestParams GetMessagesReq
}
type GetMessagesReq struct {
	ConversationId string `url:"conversation_id" validate:"required"` // 会话 ID (ChatMessageRespSSEData.ConversationId or ChatMessageResp.ConversationId)
	User           string `url:"user" validate:"required"`            // 用户标识
	FirstId        string `url:"first_id"`                            // 当前页第一条聊天记录的 ID
	Limit          int    `url:"limit"`                               // 一次请求返回多少条聊天记录，默认 20 条
}
type GetMessagesResp struct {
	Limit   int  `json:"limit"`
	HasMore bool `json:"has_more"`
	Data    []struct {
		Id             string `json:"id"`
		ConversationId string `json:"conversation_id"`
		Inputs         struct {
			Name string `json:"name"`
		} `json:"inputs"`
		Query              string        `json:"query"`
		Answer             string        `json:"answer"`
		MessageFiles       []interface{} `json:"message_files"`
		Feedback           interface{}   `json:"feedback"`
		RetrieverResources []struct {
			Position     int     `json:"position"`
			DatasetId    string  `json:"dataset_id"`
			DatasetName  string  `json:"dataset_name"`
			DocumentId   string  `json:"document_id"`
			DocumentName string  `json:"document_name"`
			SegmentId    string  `json:"segment_id"`
			Score        float64 `json:"score"`
			Content      string  `json:"content"`
		} `json:"retriever_resources"`
		CreatedAt int `json:"created_at"`
	} `json:"data"`
}

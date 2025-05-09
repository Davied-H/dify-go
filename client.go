package dify

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/go-querystring/query"
	"github.com/tmaxmax/go-sse"
)

var (
	ApiPathChatMessage  = "/chat-messages"
	ApiPathUploadFile   = "/files/upload"
	ApiPathStopTask     = "/chat-messages/%s/stop"
	ApiPathGetSuggested = "/messages/%s/suggested"
	ApiPathGetMessages  = "/messages"

	ResponseModeBlocking  = "blocking"
	ResponseModeStreaming = "streaming"
)

type ClientI interface {
	ChatMessage(ctx context.Context, option ChatMessageOption) (*ChatMessageResp, error)
	UploadFile(ctx context.Context, option UploadFileOption) (*UploadFileResp, error)
	UploadFileViaGin(ctx context.Context, option UploadFileViaGinOption) (*UploadFileResp, error)
	StopTask(ctx context.Context, option StopTaskOption) (*StopTaskResp, error)
	GetSuggested(ctx context.Context, option GetSuggestedOption) (*GetSuggestedResp, error)
	GetMessages(ctx context.Context, option GetMessagesOption) (*GetMessagesResp, error)
}

type Client struct {
	config ClientConfig
}

func NewClient(apiUrl string, opts ...Option) ClientI {
	config := DefaultConfig(apiUrl)

	for _, opt := range opts {
		opt(config)
	}

	return NewClientWithConfig(*config)
}

func NewClientWithConfig(config ClientConfig) ClientI {
	return &Client{
		config: config,
	}
}

type requestOption struct {
	Method          string
	ApiPath         string
	ApiKey          string
	RequestBody     interface{}
	RequestFormData requestOptionRequestFormData
	Headers         map[string]string
}
type requestOptionRequestFormData struct {
	Buffer *bytes.Buffer
	Writer *multipart.Writer
}

func (c *Client) request(ctx context.Context, option requestOption) (readCloser *http.Response, err error) {

	var body *strings.Reader

	contentType := option.Headers["Content-Type"]
	if contentType == "application/json" || contentType == "" {
		bodyBytes, marshalErr := json.Marshal(option.RequestBody)
		if marshalErr != nil {
			fmt.Printf("marshalErr: %s\n", marshalErr.Error())
			return
		}
		body = strings.NewReader(string(bodyBytes))
	}

	if strings.Contains(contentType, "multipart/form-data") {
		body = strings.NewReader(option.RequestFormData.Buffer.String())
	}

	request, newRequestErr := http.NewRequestWithContext(ctx, option.Method, c.config.ApiBaseUrl+option.ApiPath, body)
	if newRequestErr != nil {
		err = errors.New(fmt.Sprintf("newRequestErr: %s", newRequestErr.Error()))
		return
	}

	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", option.ApiKey))
	switch option.Method {
	case http.MethodPost:
		request.Header.Add("Content-Type", "application/json")
	}
	for k, v := range option.Headers {
		request.Header.Add(k, v)
	}

	readCloser, doErr := c.config.HttpClient.Do(request)
	if doErr != nil {
		err = errors.New(fmt.Sprintf("doResp: %s", doErr.Error()))
		return
	}
	return
}

// ChatMessage 发送对话消息
func (c *Client) ChatMessage(ctx context.Context, option ChatMessageOption) (resp *ChatMessageResp, err error) {

	// 校验参数
	validate := validator.New()
	validateErr := validate.Struct(option)
	if validateErr != nil {
		err = errors.New(fmt.Sprintf("validateErr: %s", validateErr.Error()))
		return
	}
	if option.RequestBody.ResponseMode == ResponseModeStreaming && option.OnEvent == nil {
		err = errors.New("when the response mode is streaming, OnEvent is required")
		return
	}

	// 发起请求
	response, requestErr := c.request(ctx, requestOption{
		Method:      http.MethodPost,
		ApiPath:     ApiPathChatMessage,
		ApiKey:      option.ApiKey,
		RequestBody: option.RequestBody,
		Headers:     nil,
	})
	if requestErr != nil {
		err = errors.New(fmt.Sprintf("requestErr: %s", requestErr.Error()))
		return
	}

	// 错误处理
	if response.StatusCode != http.StatusOK {
		all, _ := io.ReadAll(response.Body)
		err = errors.New(string(all))
		return
	}

	// 解析返回参
	if option.RequestBody.ResponseMode == "streaming" {
		for ev, sseReadErr := range sse.Read(response.Body, nil) {
			if sseReadErr != nil {
				fmt.Printf("Error reading SSE error: %s", sseReadErr.Error())
				break
			}

			var difySSEData ChatMessageRespSSEData
			unmarshalErr := json.Unmarshal([]byte(ev.Data), &difySSEData)
			if unmarshalErr != nil {
				fmt.Printf("unmarshalErr: %s\n", unmarshalErr.Error())
				continue
			}

			option.OnEvent(difySSEData)
		}
	} else {
		all, readAllErr := io.ReadAll(response.Body)
		if readAllErr != nil {
			err = errors.New(fmt.Sprintf("readAllErr: %s", readAllErr.Error()))
			return
		}

		unmarshalErr := json.Unmarshal(all, &resp)
		if unmarshalErr != nil {
			err = errors.New(fmt.Sprintf("unmarshalErr: %s", unmarshalErr.Error()))
			return
		}
	}

	return
}

// UploadFile 上传文件
func (c *Client) UploadFile(ctx context.Context, option UploadFileOption) (resp *UploadFileResp, err error) {
	// 校验参数
	validate := validator.New()
	validateErr := validate.Struct(option)
	if validateErr != nil {
		err = errors.New(fmt.Sprintf("validateErr: %s", validateErr.Error()))
		return
	}

	// 发起请求
	buffer := &bytes.Buffer{}
	writer := multipart.NewWriter(buffer)
	writeFieldErr := writer.WriteField("user", option.RequestFormData.User)
	if writeFieldErr != nil {
		err = errors.New(fmt.Sprintf("writeFieldErr: %s", writeFieldErr.Error()))
		return
	}
	part2, createFormFileErr := writer.CreateFormFile("file", filepath.Base(option.RequestFormData.File.Name()))
	if createFormFileErr != nil {
		err = errors.New(fmt.Sprintf("createFormFileErr: %s", createFormFileErr.Error()))
		return
	}
	_, copyErr := io.Copy(part2, option.RequestFormData.File)
	if copyErr != nil {
		err = errors.New(fmt.Sprintf("copyErr: %s", copyErr.Error()))
		return
	}
	_ = writer.Close()

	requestResp, requestErr := c.request(ctx, requestOption{
		Method:      http.MethodPost,
		ApiPath:     ApiPathUploadFile,
		ApiKey:      option.ApiKey,
		RequestBody: nil,
		RequestFormData: requestOptionRequestFormData{
			Buffer: buffer,
			Writer: writer,
		},
		Headers: map[string]string{
			"Content-Type": writer.FormDataContentType(),
		},
	})
	if requestErr != nil {
		err = errors.New(fmt.Sprintf("requestErr: %s", requestErr.Error()))
		return
	}

	// 解析返回参
	all, readAllErr := io.ReadAll(requestResp.Body)
	if readAllErr != nil {
		err = errors.New(fmt.Sprintf("readAllErr: %s", readAllErr.Error()))
		return
	}
	unmarshalErr := json.Unmarshal(all, &resp)
	if unmarshalErr != nil {
		err = errors.New(fmt.Sprintf("unmarshalErr: %s", unmarshalErr.Error()))
		return
	}

	return
}

// UploadFileViaGin 上传文件通过gin
func (c *Client) UploadFileViaGin(ctx context.Context, option UploadFileViaGinOption) (resp *UploadFileResp, err error) {
	// 校验参数
	validate := validator.New()
	validateErr := validate.Struct(option)
	if validateErr != nil {
		err = errors.New(fmt.Sprintf("validateErr: %s", validateErr.Error()))
		return
	}

	// 发起请求
	buffer := &bytes.Buffer{}
	writer := multipart.NewWriter(buffer)
	writeFieldErr := writer.WriteField("user", option.RequestFormData.User)
	if writeFieldErr != nil {
		err = errors.New(fmt.Sprintf("writeFieldErr: %s", writeFieldErr.Error()))
		return
	}
	file, openFileErr := option.RequestFormData.FormFile.Open()
	if openFileErr != nil {
		err = errors.New(fmt.Sprintf("openFileErr: %s", openFileErr.Error()))
		return
	}
	defer func(file multipart.File) {
		_ = file.Close()
	}(file)
	part2, createFormFileErr := writer.CreateFormFile("file", option.RequestFormData.FormFile.Filename)
	if createFormFileErr != nil {
		err = errors.New(fmt.Sprintf("createFormFileErr: %s", createFormFileErr.Error()))
		return
	}
	_, copyErr := io.Copy(part2, file)
	if copyErr != nil {
		err = errors.New(fmt.Sprintf("copyErr: %s", copyErr.Error()))
		return
	}
	_ = writer.Close()

	requestResp, requestErr := c.request(ctx, requestOption{
		Method:      http.MethodPost,
		ApiPath:     ApiPathUploadFile,
		ApiKey:      option.ApiKey,
		RequestBody: nil,
		RequestFormData: requestOptionRequestFormData{
			Buffer: buffer,
			Writer: writer,
		},
		Headers: map[string]string{
			"Content-Type": writer.FormDataContentType(),
		},
	})
	if requestErr != nil {
		err = errors.New(fmt.Sprintf("requestErr: %s", requestErr.Error()))
		return
	}

	// 解析返回参
	all, readAllErr := io.ReadAll(requestResp.Body)
	if readAllErr != nil {
		err = errors.New(fmt.Sprintf("readAllErr: %s", readAllErr.Error()))
		return
	}
	unmarshalErr := json.Unmarshal(all, &resp)
	if unmarshalErr != nil {
		err = errors.New(fmt.Sprintf("unmarshalErr: %s", unmarshalErr.Error()))
		return
	}

	return
}

// StopTask 停止响应
func (c *Client) StopTask(ctx context.Context, option StopTaskOption) (resp *StopTaskResp, err error) {
	// 校验参数
	validate := validator.New()
	validateErr := validate.Struct(option)
	if validateErr != nil {
		err = errors.New(fmt.Sprintf("validateErr: %s", validateErr.Error()))
		return
	}

	// 发起请求
	requestResp, requestErr := c.request(ctx, requestOption{
		Method:      http.MethodPost,
		ApiPath:     fmt.Sprintf(ApiPathStopTask, option.TaskId),
		ApiKey:      option.ApiKey,
		RequestBody: option.RequestBody,
		Headers:     nil,
	})
	if requestErr != nil {
		err = errors.New(fmt.Sprintf("requestErr: %s", requestErr.Error()))
		return
	}

	// 解析返回参
	all, readAllErr := io.ReadAll(requestResp.Body)
	if readAllErr != nil {
		err = errors.New(fmt.Sprintf("readAllErr: %s", readAllErr.Error()))
		return
	}
	unmarshalErr := json.Unmarshal(all, &resp)
	if unmarshalErr != nil {
		err = errors.New(fmt.Sprintf("unmarshalErr: %s", unmarshalErr.Error()))
		return
	}

	return
}

// GetSuggested 获取下一轮建议问题列表
func (c *Client) GetSuggested(ctx context.Context, option GetSuggestedOption) (resp *GetSuggestedResp, err error) {
	// 校验参数
	validate := validator.New()
	validateErr := validate.Struct(option)
	if validateErr != nil {
		err = errors.New(fmt.Sprintf("validateErr: %s", validateErr.Error()))
		return
	}

	// 发起请求
	values, _ := query.Values(option.RequestParams)
	params := values.Encode()
	requestResp, requestErr := c.request(ctx, requestOption{
		Method:      http.MethodGet,
		ApiPath:     fmt.Sprintf(ApiPathGetSuggested, option.MessageId) + "?" + params,
		ApiKey:      option.ApiKey,
		RequestBody: nil,
		Headers:     nil,
	})
	if requestErr != nil {
		err = errors.New(fmt.Sprintf("requestErr: %s", requestErr.Error()))
		return
	}

	// 解析返回参
	all, readAllErr := io.ReadAll(requestResp.Body)
	if readAllErr != nil {
		err = errors.New(fmt.Sprintf("readAllErr: %s", readAllErr.Error()))
		return
	}
	unmarshalErr := json.Unmarshal(all, &resp)
	if unmarshalErr != nil {
		err = errors.New(fmt.Sprintf("unmarshalErr: %s", unmarshalErr.Error()))
		return
	}

	return
}

// GetMessages 获取会话历史消息
func (c *Client) GetMessages(ctx context.Context, option GetMessagesOption) (resp *GetMessagesResp, err error) {
	// 校验参数
	validate := validator.New()
	validateErr := validate.Struct(option)
	if validateErr != nil {
		err = errors.New(fmt.Sprintf("validateErr: %s", validateErr.Error()))
		return
	}

	// 发起请求
	values, _ := query.Values(option.RequestParams)
	params := values.Encode()
	requestResp, requestErr := c.request(ctx, requestOption{
		Method:  http.MethodGet,
		ApiPath: ApiPathGetMessages + "?" + params,
		ApiKey:  option.ApiKey,
	})
	if requestErr != nil {
		err = errors.New(fmt.Sprintf("requestErr: %s", requestErr.Error()))
		return
	}

	// 解析返回参
	all, readAllErr := io.ReadAll(requestResp.Body)
	if readAllErr != nil {
		err = errors.New(fmt.Sprintf("readAllErr: %s", readAllErr.Error()))
		return
	}
	unmarshalErr := json.Unmarshal(all, &resp)
	if unmarshalErr != nil {
		err = errors.New(fmt.Sprintf("unmarshalErr: %s", unmarshalErr.Error()))
		return
	}

	return
}

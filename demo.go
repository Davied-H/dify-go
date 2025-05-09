package dify

import (
	"context"
	"fmt"
	"os"

	"github.com/duke-git/lancet/v2/formatter"
	"github.com/joho/godotenv"
)

func init() {
	loadEnvErr := godotenv.Load(".env")
	if loadEnvErr != nil {
		panic(loadEnvErr)
	}
}

func chatMessageStreamDemo(query string) {
	client := NewClient("http://dify.hubs.org.cn/v1")
	user := "dong"
	_, chatMessageErr := client.ChatMessage(context.TODO(), ChatMessageOption{
		ApiKey: os.Getenv("DIFY_API_KEY"),
		OnEvent: func(ev ChatMessageRespSSEData) {
			switch ev.Event {
			case "workflow_started":
				fmt.Printf("工作流开始执行, \n"+
					"\t任务ID: %s\n"+
					"\t会话ID: %s\n"+
					"\t用户标识: %s\n"+
					"\t消息ID: %s\n", ev.TaskId, ev.ConversationId, user, ev.MessageId)
			case "node_started":
			case "node_finished":
			case "message":
				fmt.Printf("%s", ev.Answer)
			case "message_end":
			case "workflow_finished":
				fmt.Println("\n工作流结束执行")
			}
		},
		RequestBody: ChatMessageReq{
			Inputs: map[string]interface{}{
				"role": "唐老鸭",
			},
			Query:          query,
			ResponseMode:   ResponseModeStreaming,
			ConversationId: "",
			User:           user,
			Files:          nil,
		},
	})
	if chatMessageErr != nil {
		fmt.Println("chatMessageErr: ", chatMessageErr.Error())
	}
}

func chatMessageBlockDemo(query string) {
	client := NewClient("http://dify.hubs.org.cn/v1")
	chatMessageResp, chatMessageErr := client.ChatMessage(context.TODO(), ChatMessageOption{
		ApiKey: os.Getenv("DIFY_API_KEY"),
		RequestBody: ChatMessageReq{
			Inputs: map[string]interface{}{
				"role": "唐老鸭",
			},
			Query:          query,
			ResponseMode:   ResponseModeBlocking,
			ConversationId: "",
			User:           "dong",
			Files:          nil,
		},
	})
	if chatMessageErr != nil {
		fmt.Println("chatMessageErr: ", chatMessageErr.Error())
	}
	fmt.Println("chatMessageResp: ", chatMessageResp)
}

func stopTaskDemo(taskId string) {
	client := NewClient("http://dify.hubs.org.cn/v1")
	stopTaskResp, stopTaskErr := client.StopTask(context.TODO(), StopTaskOption{
		ApiKey: os.Getenv("DIFY_API_KEY"),
		TaskId: taskId,
		RequestBody: StopTaskReq{
			User: "dong",
		},
	})
	if stopTaskErr != nil {
		fmt.Println("stopTaskErr: ", stopTaskErr.Error())
		return
	}
	fmt.Println("stopTaskResp: ", stopTaskResp)
}

func getSuggestedDemo(messageId string) {
	client := NewClient("http://dify.hubs.org.cn/v1")
	getSuggested, getSuggestedErr := client.GetSuggested(context.TODO(), GetSuggestedOption{
		ApiKey:    os.Getenv("DIFY_API_KEY"),
		MessageId: messageId,
		RequestParams: GetSuggestedReq{
			User: "dong",
		},
	})
	if getSuggestedErr != nil {
		fmt.Println("getSuggestedErr: ", getSuggestedErr.Error())
		return
	}
	fmt.Println("getSuggested: ", getSuggested)
}

func getMessagesDemo(conversationId string) {
	client := NewClient("http://dify.hubs.org.cn/v1")
	var getMessagesResp, getMessageErr = client.GetMessages(context.TODO(), GetMessagesOption{
		ApiKey: os.Getenv("DIFY_API_KEY"),
		RequestParams: GetMessagesReq{
			ConversationId: conversationId,
			User:           "dong",
			FirstId:        "",
			Limit:          20,
		},
	})
	if getMessageErr != nil {
		fmt.Println("getMessageErr: ", getMessageErr.Error())
		return
	}
	fmt.Println("getMessagesResp: ", getMessagesResp)
}

func uploadFileDemo(f *os.File) {
	client := NewClient("http://dify.hubs.org.cn/v1")
	uploadFileResp, uploadFileErr := client.UploadFile(context.TODO(), UploadFileOption{
		ApiKey: os.Getenv("DIFY_API_KEY"),
		RequestFormData: UploadFileReq{
			File: f,
			User: "dong",
		},
	})
	if uploadFileErr != nil {
		fmt.Println("uploadFileErr: ", uploadFileErr.Error())
		return
	}
	pretty, _ := formatter.Pretty(uploadFileResp)
	fmt.Println("uploadFileResp: ", pretty)
}

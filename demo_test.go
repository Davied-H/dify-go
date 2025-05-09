package dify

import (
	"os"
	"testing"
)

func Test_chatMessageStreamDemo(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"chatMessageStreamDemo"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chatMessageStreamDemo("我是6个月宝宝的爸爸，我该做什么")
		})
	}
}

func Test_chatMessageBlockDemo(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"chatMessageBlockDemo"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chatMessageBlockDemo("生成个爸爸带娃指南，要求字数5000字，markdown形式，简单易懂")
		})
	}
}

func Test_stopTaskDemo(t *testing.T) {
	type args struct {
		taskId string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "stopTaskDemo",
			args: args{
				taskId: "4e274f67-37e2-4380-a3c8-e441c819b59a",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stopTaskDemo(tt.args.taskId)
		})
	}
}

func Test_getSuggestedDemo(t *testing.T) {
	type args struct {
		messageId string
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "getSuggestedDemo", args: args{
			messageId: "7069bf38-0966-4f48-980a-b5fcc854c380",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getSuggestedDemo(tt.args.messageId)
		})
	}
}

func Test_getMessagesDemo(t *testing.T) {
	type args struct {
		conversationId string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "getMessagesDemo",
			args: args{
				conversationId: "1a419dba-2151-4bbd-8252-722741c78a5d",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getMessagesDemo(tt.args.conversationId)
		})
	}
}

func Test_uploadFileDemo(t *testing.T) {
	type args struct {
		f *os.File
	}

	openFile, err := os.Open("./test.txt")
	if err != nil {
		panic(err)
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "uploadFileDemo",
			args: args{
				f: openFile,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uploadFileDemo(tt.args.f)
		})
	}
}

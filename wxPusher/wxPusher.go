package wxPusher

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type Loli struct {
	// 密钥
	AppToken string `json:"appToken"`
	// 内容
	Content string `json:"content"`
	// 消息摘要
	Summary string `json:"summary"`
	// 内容类型
	ContentType int `json:"contentType"`
	// 主题id 为空不转发
	TopicIds []int `json:"topicIds"`
	// 个人id 为空不转发
	Uids []string `json:"uids"`
}

func Send(appToken string, content string, summary string, contentType int, topicIds []int, uIds []string) ([]byte, error) {
	url := "https://wxpusher.zjiecode.com/api/send/message"
	loli := Loli{appToken, content, summary, contentType, topicIds, uIds}
	data, _ := json.Marshal(loli)
	resp, err := http.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, _ = io.ReadAll(resp.Body)
	return data, nil
}

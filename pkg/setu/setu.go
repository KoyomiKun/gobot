package setu

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension"
	"github.com/wdvxdr1123/ZeroBot/message"
)

const (
	baseUrl       = "https://api.lolicon.app/setu/v2"
	msgInfoFormat = "\n%s : %s,"
)

func init() {
	zero.OnCommand("色图").Handle(func(ctx *zero.Ctx) {
		parseSetu(ctx, false)
	})
	zero.OnCommand("R18色图").Handle(func(ctx *zero.Ctx) {
		parseSetu(ctx, true)
	})
}

type Data struct {
	Pid        int               `json:"pid,omitempty"`
	P          int               `json:"p,omitempty"`
	Uid        int               `json:"uid,omitempty"`
	Title      string            `json:"title,omitempty"`
	Author     string            `json:"author,omitempty"`
	R18        bool              `json:"r18,omitempty"`
	Width      int               `json:"width,omitempty"`
	Height     int               `json:"height,omitempty"`
	Tags       []string          `json:"tags,omitempty"`
	Ext        string            `json:"ext,omitempty"`
	UploadData int               `json:"uploadData,omitempty"`
	Urls       map[string]string `json:"urls,omitempty"`
}

type Resp struct {
	Err  string  `json:"error,omitempty"`
	Data []*Data `json:"data,omitempty"`
}

func parseSetu(ctx *zero.Ctx, r18 bool) {

	var cmd extension.CommandModel
	err := ctx.Parse(&cmd)
	if err != nil {
		log.Errorf("fail parse ctx: %s", ctx.Event.RawMessage)
		ctx.Send(message.Text("命令解析出错，请联系管理员"))
		return
	}

	if cmd.Args == "" {
		ctx.Send(message.Text("请附带关键字参数"))
		return
	}

	go ctx.Send(
		message.Message{
			message.At(ctx.Event.UserID),
			message.Text("开始查找..."),
		},
	)

	setu := getSetu(cmd.Args, r18)
	if setu == nil {
		ctx.Send(message.Text("查找图片失败"))
		return
	}
	if len(setu) == 0 {
		ctx.Send(message.Text("该关键词下搜不到图片"))
		return
	}

	topSetu := setu[0]
	picUrl, ok := topSetu.Urls["original"]
	if !ok {
		ctx.Send(message.Text("Key解析失败，请联系管理员处理"))
		return
	}

	var msgInfoBuilder strings.Builder
	msgInfoBuilder.WriteString(fmt.Sprintf(msgInfoFormat, "标题", topSetu.Title))
	msgInfoBuilder.WriteString(fmt.Sprintf(msgInfoFormat, "作者", topSetu.Author))
	msgInfoBuilder.WriteString(fmt.Sprintf(msgInfoFormat, "标签", strings.Join(topSetu.Tags, ",")))

	ctx.Send(message.Message{
		message.At(ctx.Event.UserID),
		message.Text(msgInfoBuilder.String()),
		message.Image(picUrl),
	})
}
func getSetu(keyword string, r18 bool) []*Data {

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, baseUrl, nil)
	if err != nil {
		log.Errorf("fail getting setu: %v", err)
		return nil
	}
	query := req.URL.Query()
	query.Add("keyword", keyword)
	if r18 {
		query.Add("r18", "1")
	} else {
		query.Add("r18", "0")
	}
	req.URL.RawQuery = query.Encode()

	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("fail getting setu: %v", err)
		return nil
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	respBody := Resp{}
	err = json.Unmarshal(respBytes, &respBody)
	if err != nil {
		log.Errorf("fail unmarshaling setu: %v", err)
		return nil
	}

	if respBody.Err != "" {
		log.Errorf("fail getting setu: %s", respBody.Err)
		return nil
	}
	return respBody.Data
}
